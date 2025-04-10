// SPDX-License-Identifier: MIT

//go:build go1.24

// Package recipes / edk2
package recipes

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
)

func TestEdk2(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	common := CommonOpts{
		SdkURL:              "ghcr.io/9elements/firmware-action/edk2-stable202105:main",
		OutputDir:           "output",
		ContainerOutputDirs: []string{"Build/"},
	}

	testCases := []struct {
		name        string
		edk2Options Edk2Opts
		version     string
		gccVersion  string
		wantErr     error
	}{
		{
			name: "normal build",
			edk2Options: Edk2Opts{
				CommonOpts:    common,
				DefconfigPath: "defconfig",
				Edk2Specific: Edk2Specific{
					BuildCommand: "source ./edksetup.sh; build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG -t GCC5",
				},
			},
			version:    "edk2-stable202105",
			gccVersion: "GCC5",
			wantErr:    nil,
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
			tc.edk2Options.RepoPath = filepath.Join(tmpDir, "Edk2")

			// Create __tmp_files__ directory to store source-code
			// mostly useful for repeated local-run tests to save bandwidth and time
			tmpFiles := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__")
			err = os.MkdirAll(tmpFiles, 0o750)
			assert.NoError(t, err)
			t.Chdir(tmpFiles)

			// Clone edk2 repo
			_, err = os.Stat(tc.version)
			if err != nil {
				cmd := exec.Command("git", "clone", "--recurse-submodules", "--branch", tc.version, "--depth", "1", "https://github.com/tianocore/edk2.git", tc.version)
				err = cmd.Run()
				assert.NoError(t, err)
			}
			cmd := exec.Command("cp", "-R", tc.version, filepath.Join(tmpDir, "Edk2"))
			err = cmd.Run()
			assert.NoError(t, err)

			// Change current working directory
			t.Chdir(tmpDir)

			// Create "defconfig_path" file
			err = os.WriteFile(tc.edk2Options.DefconfigPath, []byte("-D BOOTLOADER=COREBOOT"), 0o644)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, tc.edk2Options.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			tc.edk2Options.OutputDir = outputPath

			// Try to build edk2
			err = tc.edk2Options.buildFirmware(ctx, client)
			assert.NoError(t, err)

			// Check artifacts
			_, err = os.Stat(outputPath)
			assert.NoError(t, err)
			fileInfo, err := os.Stat(filepath.Join(outputPath, "Build"))
			assert.NoError(t, err)
			assert.True(t, fileInfo.IsDir())
		})
	}
}
