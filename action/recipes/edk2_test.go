// SPDX-License-Identifier: MIT

// Package recipes / edk2
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
	"github.com/stretchr/testify/assert"
)

func TestEdk2(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(pwd) // nolint:errcheck

	// Use "" if you want to test containers from github package registry
	// Use "../../docker/edk2" if you want to test containers built fresh from Dockerfile
	dockerfilePath := ""
	if false {
		dockerfilePath, err = filepath.Abs("../../docker/edk2")
		assert.NoError(t, err)
	}

	testCases := []struct {
		name        string
		edk2Version string
		platform    string
		arch        string
		release     string
		wantErr     error
	}{
		{
			name:        "normal build",
			edk2Version: "edk2-stable202105",
			platform:    "UefiPayloadPkg/UefiPayloadPkg.dsc",
			arch:        "X64",
			release:     "DEBUG",
			wantErr:     nil,
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
			opts := map[string]string{
				"target":           "edk2",
				"sdk_version":      fmt.Sprintf("%s:main", tc.edk2Version),
				"architecture":     tc.arch,
				"repo_path":        filepath.Join(tmpDir, "Edk2"),
				"defconfig_path":   "defconfig",
				"containerWorkDir": "/Edk2",
				"GITHUB_WORKSPACE": "/Edk2",
				"output":           "output",
			}
			getFunc := func(key string) string {
				return opts[key]
			}
			common, err := commonGetOpts(getFunc, getFunc)
			assert.NoError(t, err)
			edk2Opts := edk2Opts{
				platform:    tc.platform,
				releaseType: tc.release,
			}

			// Change current working directory
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)
			defer os.Chdir(pwd) // nolint:errcheck

			// Clone edk2 repo
			cmd := exec.Command("git", "clone", "--recurse-submodules", "--branch", tc.edk2Version, "--depth", "1", "https://github.com/tianocore/edk2.git", "Edk2")
			err = cmd.Run()
			assert.NoError(t, err)
			err = os.Chdir(common.repoPath)
			assert.NoError(t, err)

			// Create "defconfig_path" file
			err = os.WriteFile(common.defconfigPath, []byte("-D BOOTLOADER=COREBOOT -a IA32 -t GCC5"), 0o644)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, common.outputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			artifacts := []container.Artifacts{
				{
					ContainerPath: filepath.Join(common.containerWorkDir, "Build"),
					ContainerDir:  true,
					HostPath:      outputPath,
					HostDir:       true,
				},
			}

			// Try to build edk2
			err = edk2(ctx, client, &common, dockerfilePath, &edk2Opts, &artifacts)
			assert.NoError(t, err)

			// Check artifacts
			fileInfo, err := os.Stat(filepath.Join(outputPath, "Build"))
			assert.NoError(t, err)
			assert.True(t, fileInfo.IsDir())
		})
	}
}
