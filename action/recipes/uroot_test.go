// SPDX-License-Identifier: MIT

// Package recipes / uroot
package recipes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestURoot(t *testing.T) {
	// This test is a slow (like 200 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(pwd) // nolint:errcheck

	URootOpts := URootOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output",
			ContainerOutputFiles: []string{
				"initramfs.cpio",
			},
		},
	}

	testCases := []struct {
		name          string
		uRootVersion  string
		golangVersion string
		arch          string
		wantErr       error
	}{
		{
			name:          "normal build v0.14 in v1.x",
			uRootVersion:  "v0.14.0",
			golangVersion: "1",
			arch:          "amd64",
			wantErr:       nil,
		},
		{
			name:          "normal build v0.13.1 in v1.x",
			uRootVersion:  "v0.13.1",
			golangVersion: "1",
			arch:          "amd64",
			wantErr:       nil,
		},
		{
			name:          "normal build v0.12 in v1.x",
			uRootVersion:  "v0.12.0",
			golangVersion: "1",
			arch:          "amd64",
			wantErr:       nil,
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

			myURootOpts := URootOpts
			myURootOpts.SdkURL = fmt.Sprintf("golang:%s", tc.golangVersion)
			myURootOpts.BuildCommand = fmt.Sprintf("go build; GOARCH=%s ./u-root -o initramfs.cpio core boot", tc.arch)
			myURootOpts.RepoPath = filepath.Join(tmpDir, "u-root")

			// Change current working directory
			pwd, err := os.Getwd()
			defer os.Chdir(pwd) // nolint:errcheck
			assert.NoError(t, err)
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Clone coreboot repo
			cmd := exec.Command("git", "clone", "--branch", tc.uRootVersion, "--depth", "1", "https://github.com/u-root/u-root.git")
			err = cmd.Run()
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, myURootOpts.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			myURootOpts.OutputDir = outputPath

			// Try to build u-root initramfs
			_, err = myURootOpts.buildFirmware(ctx, client, "")
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "initramfs.cpio")), os.ErrExist)
		})
	}
	assert.NoError(t, os.Chdir(pwd)) // just to make sure
}
