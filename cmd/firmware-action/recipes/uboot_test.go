// SPDX-License-Identifier: MIT

//go:build go1.24

// Package recipes / uboot
package recipes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestUBoot(t *testing.T) {
	// This test is a slow (like 200 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)

	UBootOpts := UBootOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output",
			ContainerOutputFiles: []string{
				"u-boot",
			},
		},
		DefconfigPath: "uboot_defconfig",
		Arch:          "arm64",
	}

	testCases := []struct {
		name          string
		uBootVersion  string
		golangVersion string
		wantErr       error
	}{
		{
			name:         "normal build v2025.01",
			uBootVersion: "2025.01",
			wantErr:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)

			defer client.Close()

			// Prepare options
			tmpDir := t.TempDir()

			myUBootOpts := UBootOpts
			myUBootOpts.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/uboot_%s:main", tc.uBootVersion)
			myUBootOpts.RepoPath = filepath.Join(tmpDir, "u-boot")

			// Change current working directory
			t.Chdir(tmpDir)

			// Clone coreboot repo
			cmd := exec.Command("bash", "-c", fmt.Sprintf("git clone https://source.denx.de/u-boot/u-boot.git; cd u-boot; git fetch -a; git checkout v%s", tc.uBootVersion))
			err = cmd.Run()
			assert.NoError(t, err)

			// Copy over defconfig file into tmpDir/linux
			defconfigPath := filepath.Join(tmpDir, myUBootOpts.DefconfigPath)
			repoRootPath, err := filepath.Abs(filepath.Join(pwd, "../../.."))
			assert.NoError(t, err)
			//   common.RepoPath = path to end user repository (in this case somewhere in /tmp)
			//   repoRootPath    = path to our repository with this code (contains configuration files for testing)
			defconfigLocalPath, err := filepath.Abs(filepath.Join(
				repoRootPath,
				fmt.Sprintf("tests/uboot_%s/uboot.defconfig", tc.uBootVersion),
			))
			assert.NoErrorf(t, err, "encountered issue with missing files, is '%s' the root of the repo?", repoRootPath)
			err = filesystem.CopyFile(
				defconfigLocalPath,
				defconfigPath,
			)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, myUBootOpts.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)

			myUBootOpts.OutputDir = outputPath

			// Try to build u-boot initramfs
			err = myUBootOpts.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "u-boot")), os.ErrExist)
		})
	}
}
