// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCorebootProcessBlobs(t *testing.T) {
	testCases := []struct {
		name            string
		corebootOptions map[string]string
		expected        []BlobDef
		wantErr         error
	}{
		{
			name:            "empty",
			corebootOptions: map[string]string{},
			expected:        []BlobDef{},
			wantErr:         nil,
		},
		{
			name: "payload exists",
			corebootOptions: map[string]string{
				"CONFIG_PAYLOAD_FILE": "dummy/path/to/payload.bin",
			},
			expected: []BlobDef{
				{
					Path:                "dummy/path/to/payload.bin",
					DestinationFilename: "payload.bin",
					KconfigKey:          "CONFIG_PAYLOAD_FILE",
				},
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pwd, err := os.Getwd()
			assert.NoError(t, err)

			tmpDir := t.TempDir()
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)
			defer os.Chdir(pwd) // nolint:errcheck

			for i := range tc.expected {
				// If we do not want error
				payloadFile := filepath.Join(tmpDir, tc.expected[i].Path)
				payloadDir := filepath.Dir(payloadFile)
				if tc.wantErr == nil {
					// Create the temp directory
					err := os.MkdirAll(payloadDir, 0o750)
					assert.NoError(t, err)
					// Create the temp file to act as payload
					err = os.WriteFile(payloadFile, []byte{}, 0o666)
					assert.NoError(t, err)
				}
			}

			opts := CorebootOpts{
				Blobs: tc.corebootOptions,
			}
			output, err := opts.ProcessBlobs()
			assert.ErrorIs(t, err, tc.wantErr)
			if err == nil {
				equal := cmp.Equal(tc.expected, output)
				if !equal {
					t.Log(cmp.Diff(tc.expected, output))
					assert.True(t, equal, "processing blob parameters failed")
				}
			}
		})
	}
}

type gitCloneOpts struct {
	projectName string
	dirName     string
	destination string
	branch      string
	tag         string
	depth       int
	url         string
	fetch       bool
}

func gitCloneWithCache(tb testing.TB, opts *gitCloneOpts) {
	tb.Helper()

	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		tb.Error(err)
	}

	// Make directory for temporary testing files
	tmpFiles := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__")
	err = os.MkdirAll(tmpFiles, 0o750)
	if err != nil {
		tb.Errorf("failed to create TMP dir: %s", err.Error())
	}

	repoPath := filepath.Join(tmpFiles, opts.dirName)

	// Clone repository into cache if not done yet
	if errors.Is(filesystem.CheckFileExists(repoPath), os.ErrNotExist) {
		err = os.Chdir(tmpFiles)
		if err != nil {
			tb.Errorf("failed to change directory to '%s': %s", tmpFiles, err.Error())
		}

		command := []string{"git", "clone"}
		if opts.branch != "" {
			command = append(command, "--branch", opts.branch)
		}
		if opts.depth != 0 {
			command = append(command, "--depth", strconv.Itoa(opts.depth))
		}
		command = append(command, opts.url, opts.dirName)

		// Clone
		cmd := exec.Command(command[0], command[1:]...)
		err = cmd.Run()
		if err != nil {
			tb.Errorf("failed to 'git clone': %s", err.Error())
		}

		// Change to repository
		err = os.Chdir(repoPath)
		if err != nil {
			tb.Errorf("failed to change directory to '%s': %s", repoPath, err.Error())
		}

		if opts.fetch || opts.tag != "" {
			// Fetch
			cmds := [][]string{
				{"git", "fetch", "-a"},
				{"git", "fetch", "-t"},
			}
			for _, cmd := range cmds {
				command := exec.Command(cmd[0], cmd[1:]...)
				err = command.Run()
				if err != nil {
					tb.Errorf("failed to 'git fetch': %s", err.Error())
				}
			}
		}

		if opts.tag != "" {
			// Checkout a tag
			cmd = exec.Command("git", "checkout", opts.tag)
			err = cmd.Run()
			if err != nil {
				tb.Errorf("failed to 'git checkout %s': %s", opts.tag, err.Error())
			}
		}

		// Init git submodules
		cmd = exec.Command("git", "submodule", "update", "--init", "--checkout")
		err = cmd.Run()
		if err != nil {
			tb.Errorf("failed to init git submodules: %s", err.Error())
		}
	}

	// Copy repository into destination
	err = filesystem.CopyDir(repoPath, opts.destination)
	if err != nil {
		tb.Errorf("failed to copy git repository from cache: %s", err.Error())
	}

	err = os.Chdir(pwd)
	if err != nil {
		tb.Error(err)
	}
}

func TestCorebootBuild(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	common := CommonOpts{
		OutputDir: "output",
		ContainerOutputFiles: []string{
			"build/coreboot.rom",
			"defconfig",
		},
	}
	// The universal module is used in this test to check version of compiled coreboot binary
	optionsUniversal := UniversalOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output-universal",
			ContainerOutputFiles: []string{
				"build_info.txt",
			},
			ContainerInputDir: "input",
		},
		UniversalSpecific: UniversalSpecific{
			BuildCommands: []string{
				"cbfstool coreboot.rom extract -n build_info -f build_info.txt",
			},
		},
	}

	testCases := []struct {
		name             string
		corebootVersion  string
		corebootOptions  CorebootOpts
		universalOptions UniversalOpts
		cmds             [][]string
		versionRegex     string
		wantErr          error
	}{
		{
			name:            "normal build for QEMU",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			universalOptions: optionsUniversal,
			versionRegex:     `4\.19`,
			wantErr:          nil,
		},
		{
			name:            "binary payload - file does not exists",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
				Blobs: map[string]string{
					"CONFIG_PAYLOAD_FILE": "my_payload",
				},
			},
			universalOptions: optionsUniversal,
			wantErr:          os.ErrNotExist,
		},
		{
			name:            "binary payload - file exists but empty",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
				Blobs: map[string]string{
					"CONFIG_ME_BIN_PATH": "intel_me.bin",
				},
			},
			universalOptions: optionsUniversal,
			cmds: [][]string{
				{"touch", "intel_me.bin"},
			},
			versionRegex: `4\.19`,
			wantErr:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)
			defer client.Close()

			// Prepare options
			tmpDir := t.TempDir()
			// Prepare options - coreboot
			tc.corebootOptions.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/coreboot_%s:main", tc.corebootVersion)
			tc.corebootOptions.RepoPath = filepath.Join(tmpDir, "coreboot")
			// Prepare options - universal module
			tc.universalOptions.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/coreboot_%s:main", tc.corebootVersion)
			tc.universalOptions.RepoPath = filepath.Join(tmpDir, tc.corebootOptions.OutputDir)

			// Change current working directory
			pwd, err := os.Getwd()
			defer os.Chdir(pwd) // nolint:errcheck
			assert.NoError(t, err)
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Clone coreboot repo
			opts := gitCloneOpts{
				dirName:     fmt.Sprintf("coreboot-%s", tc.corebootVersion),
				destination: filepath.Join(tmpDir, "coreboot"),
				tag:         tc.corebootVersion,
				url:         "https://review.coreboot.org/coreboot",
				fetch:       true,
			}
			gitCloneWithCache(t, &opts)

			// Copy over defconfig file into tmpDir
			repoRootPath, err := filepath.Abs(filepath.Join(pwd, "../../.."))
			assert.NoError(t, err)
			//   repoPath = path to end user repository (in this case somewhere in /tmp)
			//   repoRootPath    = path to our repository with this code (contains configuration files for testing)
			err = filesystem.CopyFile(
				filepath.Join(repoRootPath, fmt.Sprintf("tests/coreboot_%s/seabios.defconfig", tc.corebootVersion)),
				filepath.Join(tmpDir, tc.corebootOptions.DefconfigPath),
			)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, tc.corebootOptions.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			tc.corebootOptions.OutputDir = outputPath
			outputPathUniversal := filepath.Join(tmpDir, tc.universalOptions.OutputDir)
			tc.universalOptions.OutputDir = outputPathUniversal

			// Prep
			for cmd := range tc.cmds {
				err = exec.Command(tc.cmds[cmd][0], tc.cmds[cmd][1:]...).Run()
				assert.NoError(t, err)
			}
			// Try to build coreboot
			err = tc.corebootOptions.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}

			if tc.wantErr == nil {
				// Check coreboot version
				err = tc.universalOptions.buildFirmware(ctx, client)
				assert.ErrorIs(t, err, tc.wantErr)

				// Check file with coreboot version exists
				corebootVersionFile := filepath.Join(outputPathUniversal, "build_info.txt")
				assert.ErrorIs(t, filesystem.CheckFileExists(corebootVersionFile), os.ErrExist)

				// Find the coreboot version
				corebootVersionFileContent, err := os.ReadFile(corebootVersionFile)
				assert.NoError(t, err)
				versionEntryPatter := regexp.MustCompile(`COREBOOT_VERSION: (.*)`)
				version := string(versionEntryPatter.FindSubmatch(corebootVersionFileContent)[1])

				versionPattern := regexp.MustCompile(tc.versionRegex)
				assert.True(t, versionPattern.MatchString(version), fmt.Sprintf("found version '%s' does not match expected regex '%s'", version, tc.versionRegex))
			}
		})
	}
}

func gitCloneAsSubmoduleWithCache(tb testing.TB, opts *gitCloneOpts) {
	tb.Helper()

	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		tb.Error(err)
	}

	// Make directory for temporary testing files
	tmpFiles := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__")
	err = os.MkdirAll(tmpFiles, 0o750)
	if err != nil {
		tb.Errorf("failed to create TMP dir: %s", err.Error())
	}
	repoPath := filepath.Join(tmpFiles, opts.projectName)
	err = os.MkdirAll(repoPath, 0o750)
	if err != nil {
		tb.Errorf("failed to create TMP dir: %s", err.Error())
	}

	// Clone repository into cache if not done yet
	if errors.Is(filesystem.CheckFileExists(filepath.Join(repoPath, ".git")), os.ErrNotExist) {
		err = os.Chdir(repoPath)
		if err != nil {
			tb.Errorf("failed to change directory to '%s': %s", repoPath, err.Error())
		}

		// Make empty repository and add coreboot as git submodule
		cmds := [][]string{
			{"git", "init"},
			{"git", "submodule", "add", opts.url, opts.dirName},
			{"git", "submodule", "update", "--init", "--checkout"},
		}
		for _, cmd := range cmds {
			command := exec.Command(cmd[0], cmd[1:]...)
			err = command.Run()
			if err != nil {
				tb.Errorf("failed to run command: '%v': %s", cmd, err.Error())
			}
		}

		// Change to coreboot submodule
		err = os.Chdir(filepath.Join(repoPath, "coreboot"))
		if err != nil {
			tb.Errorf("failed to change directory to '%s': %s", repoPath, err.Error())
		}

		if opts.fetch || opts.tag != "" {
			// Fetch
			cmds := [][]string{
				{"git", "fetch", "-a"},
				{"git", "fetch", "-t"},
			}
			for _, cmd := range cmds {
				command := exec.Command(cmd[0], cmd[1:]...)
				err = command.Run()
				if err != nil {
					tb.Errorf("failed to 'git fetch': %s", err.Error())
				}
			}
		}

		if opts.tag != "" {
			// Checkout tag
			cmd := exec.Command("git", "checkout", opts.tag)
			err = cmd.Run()
			if err != nil {
				tb.Errorf("failed to 'git checkout %s': %s", opts.tag, err.Error())
			}
		}

		// Init git submodules
		cmd := exec.Command("git", "submodule", "update", "--init", "--checkout")
		err = cmd.Run()
		if err != nil {
			tb.Errorf("failed to init git submodules: %s", err.Error())
		}
	}

	if errors.Is(filesystem.CheckFileExists(repoPath), os.ErrNotExist) {
		tb.Errorf("dir does not exists '%s'", repoPath)
	}
	if errors.Is(filesystem.CheckFileExists(filepath.Join(repoPath, ".git")), os.ErrNotExist) {
		tb.Errorf("dir does not exists '%s/%s'", repoPath, ".git")
	}
	if errors.Is(filesystem.CheckFileExists(filepath.Join(repoPath, "coreboot")), os.ErrNotExist) {
		tb.Errorf("dir does not exists '%s/%s'", repoPath, "coreboot")
	}

	// Copy repository into destination
	err = filesystem.CopyDir(repoPath, opts.destination)
	if err != nil {
		tb.Errorf("failed to copy git repository from cache: %s", err.Error())
	}

	err = os.Chdir(pwd)
	if err != nil {
		tb.Error(err)
	}
}

func TestCorebootSubmodule(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	common := CommonOpts{
		OutputDir: "output",
		ContainerOutputFiles: []string{
			"build/coreboot.rom",
			"defconfig",
		},
	}
	// The universal module is used in this test to check version of compiled coreboot binary
	optionsUniversal := UniversalOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output-universal",
			ContainerOutputFiles: []string{
				"build_info.txt",
			},
			ContainerInputDir: "input",
		},
		UniversalSpecific: UniversalSpecific{
			BuildCommands: []string{
				"cbfstool coreboot.rom extract -n build_info -f build_info.txt",
			},
		},
	}

	testCases := []struct {
		name               string
		corebootVersion    string
		corebootOptions    CorebootOpts
		universalOptions   UniversalOpts
		envVars            map[string]string
		versionFileContent string
		versionRegex       string
		wantErr            error
	}{
		{
			name:            "normal build for QEMU with user-defined KERNELVERSION",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			universalOptions: optionsUniversal,
			envVars: map[string]string{
				"KERNELVERSION": "0.1.2",
			},
			versionFileContent: "",
			versionRegex:       `0\.1\.2`,
			wantErr:            nil,
		},
		{
			name:            "normal build for QEMU with user-created .coreboot-version",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			universalOptions:   optionsUniversal,
			versionFileContent: "0.1.3",
			versionRegex:       `0\.1\.3`,
			wantErr:            nil,
		},
		{
			name:            "normal build for QEMU with auto-generated KERNELVERSION",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			universalOptions:   optionsUniversal,
			versionFileContent: "",
			versionRegex:       `4\.19`,
			wantErr:            nil,
		},
		{
			name:            "normal build for QEMU with auto-generated dirty KERNELVERSION",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			universalOptions:   optionsUniversal,
			versionFileContent: "",
			versionRegex:       `4\.19\-dirty`,
			wantErr:            nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)
			defer client.Close()

			// Prepare options
			tmpDir := t.TempDir()
			projectName := fmt.Sprintf("coreboot-as-submodule-%s", tc.corebootVersion)
			dirName := "coreboot"
			// Prepare options - coreboot
			tc.corebootOptions.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/coreboot_%s:main", tc.corebootVersion)
			tc.corebootOptions.RepoPath = filepath.Join(tmpDir, projectName, dirName)
			// Prepare options - universal module
			tc.universalOptions.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/coreboot_%s:main", tc.corebootVersion)
			tc.universalOptions.RepoPath = filepath.Join(tmpDir, tc.corebootOptions.OutputDir)

			// Change current working directory
			pwd, err := os.Getwd()
			defer os.Chdir(pwd) // nolint:errcheck
			assert.NoError(t, err)
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Clone coreboot repo
			opts := gitCloneOpts{
				projectName: projectName,
				dirName:     dirName,
				destination: filepath.Join(tmpDir, projectName),
				tag:         tc.corebootVersion,
				url:         "https://review.coreboot.org/coreboot",
				fetch:       true,
			}
			gitCloneAsSubmoduleWithCache(t, &opts)

			// Copy over defconfig file into tmpDir
			repoRootPath, err := filepath.Abs(filepath.Join(pwd, "../../.."))
			assert.NoError(t, err)
			//   repoPath = path to end user repository (in this case somewhere in /tmp)
			//   repoRootPath    = path to our repository with this code (contains configuration files for testing)
			err = filesystem.CopyFile(
				filepath.Join(repoRootPath, fmt.Sprintf("tests/coreboot_%s/seabios.defconfig", tc.corebootVersion)),
				filepath.Join(tmpDir, tc.corebootOptions.DefconfigPath),
			)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, tc.corebootOptions.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			tc.corebootOptions.OutputDir = outputPath
			outputPathUniversal := filepath.Join(tmpDir, tc.universalOptions.OutputDir)
			tc.universalOptions.OutputDir = outputPathUniversal

			// Prep - environment variables
			for key, value := range tc.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
				t.Logf("Setting %s = %s\n", key, value)
			}

			// Make version file if required
			if tc.versionFileContent != "" {
				assert.NoError(t, os.WriteFile(filepath.Join(projectName, dirName, ".coreboot-version"), []byte(tc.versionFileContent), 0o666))
			}

			// Try to build coreboot
			err = tc.corebootOptions.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}

			// Check coreboot version
			err = tc.universalOptions.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check file with coreboot version exists
			corebootVersionFile := filepath.Join(outputPathUniversal, "build_info.txt")
			assert.ErrorIs(t, filesystem.CheckFileExists(corebootVersionFile), os.ErrExist)

			// Find the coreboot version
			corebootVersionFileContent, err := os.ReadFile(corebootVersionFile)
			assert.NoError(t, err)
			versionEntryPatter := regexp.MustCompile(`COREBOOT_VERSION: (.*)`)
			version := string(versionEntryPatter.FindSubmatch(corebootVersionFileContent)[1])

			versionPattern := regexp.MustCompile(tc.versionRegex)
			assert.True(t, versionPattern.MatchString(version), fmt.Sprintf("found version '%s' does not match expected regex '%s'", version, tc.versionRegex))
		})
	}
}
