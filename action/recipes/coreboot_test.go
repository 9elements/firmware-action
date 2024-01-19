// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/9elements/firmware-action/action/filesystem"
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

func TestCoreboot(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	common := CommonOpts{
		Arch:      "x86",
		OutputDir: "output",
	}

	testCases := []struct {
		name            string
		corebootVersion string
		corebootOptions CorebootOpts
		cmds            [][]string
		wantErr         error
	}{
		{
			name:            "normal build for QEMU",
			corebootVersion: "4.19",
			corebootOptions: CorebootOpts{
				CommonOpts:    common,
				DefconfigPath: "seabios_defconfig",
			},
			wantErr: nil,
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
			wantErr: os.ErrNotExist,
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
			cmds: [][]string{
				{"touch", "intel_me.bin"},
			},
			wantErr: nil,
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
			tc.corebootOptions.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/coreboot_%s:main", tc.corebootVersion)
			tc.corebootOptions.RepoPath = filepath.Join(tmpDir, "coreboot")

			// Change current working directory
			pwd, err := os.Getwd()
			defer os.Chdir(pwd) // nolint:errcheck
			assert.NoError(t, err)
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Clone coreboot repo
			cmd := exec.Command("git", "clone", "--branch", tc.corebootVersion, "--depth", "1", "https://review.coreboot.org/coreboot")
			err = cmd.Run()
			assert.NoError(t, err)

			// Copy over defconfig file into tmpDir
			repoRootPath, err := filepath.Abs(filepath.Join(pwd, "../.."))
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
			artifacts := []container.Artifacts{
				{
					ContainerPath: filepath.Join(ContainerWorkDir, "build", "coreboot.rom"),
					ContainerDir:  false,
					HostPath:      outputPath,
					HostDir:       true,
				},
				{
					ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
					ContainerDir:  false,
					HostPath:      outputPath,
					HostDir:       true,
				},
			}

			// Prep
			for cmd := range tc.cmds {
				err = exec.Command(tc.cmds[cmd][0], tc.cmds[cmd][1:]...).Run()
				assert.NoError(t, err)
			}
			// Try to build coreboot
			err = coreboot(ctx, client, &tc.corebootOptions, "", &artifacts)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}
		})
	}
}
