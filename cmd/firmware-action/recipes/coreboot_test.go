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
	"github.com/9elements/firmware-action/filesystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCorebootProcessBlobs(t *testing.T) {
	testCases := []struct {
		name            string
		corebootOptions CorebootBlobs
		expected        []BlobDef
	}{
		{
			name:            "empty",
			corebootOptions: CorebootBlobs{},
			expected:        []BlobDef{},
		},
		{
			name: "payload",
			corebootOptions: CorebootBlobs{
				PayloadFilePath: "dummy/path/to/payload.bin",
			},
			expected: []BlobDef{
				{
					Path:                "dummy/path/to/payload.bin",
					DestinationFilename: "payload",
					KconfigKey:          "CONFIG_PAYLOAD_FILE",
					IsDirectory:         false,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := corebootProcessBlobs(tc.corebootOptions)
			assert.NoError(t, err)

			equal := cmp.Equal(tc.expected, output)
			if !equal {
				fmt.Println(cmp.Diff(tc.expected, output))
				assert.True(t, equal, "processing blob parameters failed")
			}
		})
	}
}

func gitCloneWithCache(dirName string, destination string, branch string, tag string, depth int, url string, fetch bool) error {
	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Make directory for temporary testing files
	tmpFiles := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__")
	err = os.MkdirAll(tmpFiles, 0o750)
	if err != nil {
		return fmt.Errorf("%w: failed to create TMP dir", err)
	}

	repoPath := filepath.Join(tmpFiles, dirName)

	// Clone repository into cache if not done yet
	if errors.Is(filesystem.CheckFileExists(repoPath), os.ErrNotExist) {
		err = os.Chdir(tmpFiles)
		if err != nil {
			return fmt.Errorf("%w: failed to change directory to '%s'", err, tmpFiles)
		}

		command := []string{"git", "clone"}
		if branch != "" {
			command = append(command, "--branch", branch)
		}
		if depth != 0 {
			command = append(command, "--depth", strconv.Itoa(depth))
		}
		command = append(command, url, dirName)

		// Clone
		cmd := exec.Command(command[0], command[1:]...)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("%w: failed to 'git clone'", err)
		}

		// Change to repository
		err = os.Chdir(repoPath)
		if err != nil {
			return fmt.Errorf("%w: failed to change directory to '%s'", err, repoPath)
		}

		if fetch || tag != "" {
			// Fetch
			cmds := [][]string{
				{"git", "fetch", "-a"},
				{"git", "fetch", "-t"},
			}
			for _, cmd := range cmds {
				command := exec.Command(cmd[0], cmd[1:]...)
				err = command.Run()
				if err != nil {
					return fmt.Errorf("%w: failed to 'git fetch'", err)
				}
			}
		}

		if tag != "" {
			// Checkout a tag
			cmd = exec.Command("git", "checkout", tag)
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("%w: failed to 'git checkout %s'", err, tag)
			}
		}

		// Init git submodules
		cmd = exec.Command("git", "submodule", "update", "--init", "--checkout")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("%w: failed to init git submodules", err)
		}
	}

	// Copy repository into destination
	err = filesystem.CopyDir(repoPath, destination)
	if err != nil {
		return fmt.Errorf("%w: failed to copy git repository from cache", err)
	}

	err = os.Chdir(pwd)
	if err != nil {
		return err
	}

	return nil
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
				Blobs: CorebootBlobs{
					PayloadFilePath: "my_payload",
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
				Blobs: CorebootBlobs{
					IntelMePath: "intel_me.bin",
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
			err = gitCloneWithCache(fmt.Sprintf("coreboot-%s", tc.corebootVersion), filepath.Join(tmpDir, "coreboot"), "", tc.corebootVersion, 0, "https://review.coreboot.org/coreboot", true)
			assert.NoError(t, err)

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
			_, err = tc.corebootOptions.buildFirmware(ctx, client, "")
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}

			if tc.wantErr == nil {
				// Check coreboot version
				_, err = tc.universalOptions.buildFirmware(ctx, client, "")
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
