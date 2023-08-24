// SPDX-License-Identifier: MIT

// Package container
package container

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	// This test is rather slow (between 10s and 20s)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name          string
		opts          SetupOpts
		wantErr       error
		lsContains    string
		lsNotContains string
	}{
		{
			name: "empty URL",
			opts: SetupOpts{
				containerURL:      "",
				mountContainerDir: "",
				mountHostDir:      "",
				workdirContainer:  "",
			},
			wantErr:       errEmptyURL,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "empty directory strings",
			opts: SetupOpts{
				containerURL:      "ubuntu:latest",
				mountContainerDir: "",
				mountHostDir:      "",
				workdirContainer:  "",
			},
			wantErr:       errDirectoryNotSpecified,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "invalid directory strings: .",
			opts: SetupOpts{
				containerURL:      "ubuntu:latest",
				mountContainerDir: ".",
				mountHostDir:      ".",
				workdirContainer:  ".",
			},
			wantErr:       errDirectoryInvalid,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "invalid directory strings: /",
			opts: SetupOpts{
				containerURL:      "ubuntu:latest",
				mountContainerDir: "/",
				mountHostDir:      ".",
				workdirContainer:  ".",
			},
			wantErr:       errDirectoryInvalid,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "valid inputs",
			opts: SetupOpts{
				containerURL:      "ubuntu:latest",
				mountContainerDir: "/src",
				mountHostDir:      ".",
				workdirContainer:  "/src",
			},
			wantErr:       nil,
			lsContains:    "container_test.go",
			lsNotContains: "tmp",
		},
	}
	for _, tc := range testCases {
		container, err := Setup(ctx, client, &tc.opts)
		assert.ErrorIs(t, err, tc.wantErr)
		if err != nil {
			// No need to continue on err
			continue
		}

		// Get contents of current workign directory in the container
		stdout, err := container.WithExec([]string{"ls", "."}).
			Stdout(ctx)
		assert.NoError(t, err)

		// Check the directory contents
		ls := strings.Split(stdout, "\n")
		assert.True(t, slices.Contains(ls, tc.lsContains))
		if tc.lsNotContains != "" {
			assert.True(t, !slices.Contains(ls, tc.lsNotContains))
		}
	}
}

func TestGetArtifacts(t *testing.T) {
	// This test is rather slow (between 10s and 20s)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tmpDir := t.TempDir()
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()
	filename := "my_file.txt"
	prefix := "my_directory"

	opts := SetupOpts{
		containerURL:      "ubuntu:latest",
		mountContainerDir: "/src",
		mountHostDir:      ".",
		workdirContainer:  "/src",
	}

	testCases := []struct {
		name           string
		cmdToRun       [][]string
		artifacts      Artifacts
		wantErrExport  error
		filepathToTest string
		wantErrFile    error
	}{
		{
			name:     "empty directory string",
			cmdToRun: [][]string{{"true"}},
			artifacts: Artifacts{
				containerDir: "",
				hostDir:      tmpDir,
			},
			wantErrExport:  errDirectoryNotSpecified,
			filepathToTest: filename,
			wantErrFile:    os.ErrNotExist,
		},
		{
			name:     "single file",
			cmdToRun: [][]string{{"touch", filename}},
			artifacts: Artifacts{
				containerDir: "/src",
				hostDir:      tmpDir,
			},
			wantErrExport:  nil,
			filepathToTest: filename,
			wantErrFile:    os.ErrExist,
		},
		{
			name: "entire directory",
			cmdToRun: [][]string{
				{"mkdir", prefix},
				{"touch", filepath.Join(prefix, filename)},
			},
			artifacts: Artifacts{
				containerDir: "/src",
				hostDir:      tmpDir,
			},
			wantErrExport:  nil,
			filepathToTest: filepath.Join(prefix, filename),
			wantErrFile:    os.ErrExist,
		},
		{
			name:     "export non-existing directory",
			cmdToRun: [][]string{{"true"}},
			artifacts: Artifacts{
				containerDir: "/some/non-existing/directory/path",
				hostDir:      tmpDir,
			},
			wantErrExport:  errExportFailed,
			filepathToTest: filepath.Join(prefix, filename),
			wantErrFile:    os.ErrNotExist,
		},
	}
	for _, tc := range testCases {
		container, err := Setup(ctx, client, &opts)
		assert.NoError(t, err)

		// Run commands in container
		for _, cmd := range tc.cmdToRun {
			container, err = container.WithExec(cmd).
				Sync(ctx)
			assert.NoError(t, err)
		}

		// Extract artifacts
		err = GetArtifacts(ctx, container, &tc.artifacts)
		assert.ErrorIs(t, err, tc.wantErrExport)
		if err != nil {
			// No need to continue on err
			continue
		}
		assert.ErrorIs(t,
			filesystem.CheckFileExists(filepath.Join(tmpDir, tc.filepathToTest)),
			tc.wantErrFile)
	}
}
