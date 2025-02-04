// SPDX-License-Identifier: MIT

// Package recipes / universal
package recipes

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestUniversal(t *testing.T) {
	// This test is a slow (like 200 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(pwd) // nolint:errcheck

	UniversalOpts := UniversalOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output",
			ContainerOutputFiles: []string{
				"hello.txt",
			},
		},
	}

	testCases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "normal build",
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

			myUniversalOpts := UniversalOpts
			myUniversalOpts.SdkURL = "ghcr.io/9elements/firmware-action/coreboot_4.19:main"
			myUniversalOpts.BuildCommands = []string{"echo 'Hello World!'", "touch hello.txt"}
			myUniversalOpts.RepoPath = filepath.Join(tmpDir, "dummy-repo")
			err = os.Mkdir(myUniversalOpts.RepoPath, 0o755)
			assert.NoError(t, err)

			// Change current working directory
			pwd, err := os.Getwd()
			defer os.Chdir(pwd) // nolint:errcheck
			assert.NoError(t, err)
			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, myUniversalOpts.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			myUniversalOpts.OutputDir = outputPath

			// Try to build universal
			err = myUniversalOpts.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "hello.txt")), os.ErrExist)
		})
	}
	assert.NoError(t, os.Chdir(pwd)) // just to make sure
}
