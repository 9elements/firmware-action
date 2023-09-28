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
	"github.com/stretchr/testify/assert"
)

func TestCoreboot(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	testCases := []struct {
		name            string
		corebootVersion string
		corebootOptions corebootOpts
		cmds            [][]string
		wantErr         error
	}{
		{
			name:            "normal build for QEMU",
			corebootVersion: "4.19",
			corebootOptions: corebootOpts{},
			wantErr:         nil,
		},
		{
			name:            "binary payload - file does not exists",
			corebootVersion: "4.19",
			corebootOptions: corebootOpts{
				blobs: []blobDef{
					{
						actionInput:         "my_payload",
						destinationFilename: "payload",
						kconfigKey:          "CONFIG_PAYLOAD_FILE",
						isDirectory:         false,
					},
				},
			},
			wantErr: os.ErrNotExist,
		},
		{
			name:            "binary payload - file exists but empty",
			corebootVersion: "4.19",
			corebootOptions: corebootOpts{
				blobs: []blobDef{
					{
						actionInput:         "intel_me.bin",
						destinationFilename: "me.bin",
						kconfigKey:          "CONFIG_ME_BIN_PATH",
						isDirectory:         false,
					},
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
			whatever(t)
			t.Cleanup(whatever2)

			ctx := context.Background()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)
			defer client.Close()

			// Prepare options
			tmpDir := t.TempDir()
			opts := map[string]string{
				"target":           "coreboot",
				"sdk_version":      fmt.Sprintf("coreboot_%s:main", tc.corebootVersion),
				"architecture":     "x86",
				"repo_path":        filepath.Join(tmpDir, "coreboot"),
				"defconfig_path":   "seabios_defconfig",
				"containerWorkDir": "/coreboot",
				"GITHUB_WORKSPACE": "/coreboot",
				"output":           "output",
			}
			getFunc := func(key string) string {
				return opts[key]
			}
			common, err := commonGetOpts(getFunc, getFunc)
			assert.NoError(t, err)
			corebootOpts := tc.corebootOptions

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
			//   common.repoPath = path to end user repository (in this case somewhere in /tmp)
			//   repoRootPath    = path to our repository with this code (contains configuration files for testing)
			err = filesystem.CopyFile(
				filepath.Join(repoRootPath, fmt.Sprintf("tests/coreboot_%s/seabios.defconfig", tc.corebootVersion)),
				filepath.Join(tmpDir, common.defconfigPath),
			)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, common.outputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			artifacts := []container.Artifacts{
				{
					ContainerPath: filepath.Join(common.containerWorkDir, "build", "coreboot.rom"),
					ContainerDir:  false,
					HostPath:      outputPath,
					HostDir:       true,
				},
				{
					ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
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
			err = coreboot(ctx, client, &common, "", &corebootOpts, &artifacts)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}
		})
	}
}
