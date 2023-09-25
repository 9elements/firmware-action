// SPDX-License-Identifier: MIT

// Package container
package container

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
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
				ContainerURL:      "",
				MountContainerDir: "/src",
				MountHostDir:      ".",
				WorkdirContainer:  "/src",
			},
			wantErr:       errEmptyURL,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "empty directory strings",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "",
				MountHostDir:      "",
				WorkdirContainer:  "",
			},
			wantErr:       errDirectoryNotSpecified,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "invalid directory strings: .",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: ".",
				MountHostDir:      ".",
				WorkdirContainer:  ".",
			},
			wantErr:       errDirectoryInvalid,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "invalid directory strings: /",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/",
				MountHostDir:      ".",
				WorkdirContainer:  ".",
			},
			wantErr:       errDirectoryInvalid,
			lsContains:    "",
			lsNotContains: "",
		},
		{
			name: "valid inputs",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/src",
				MountHostDir:      ".",
				WorkdirContainer:  "/src",
			},
			wantErr:       nil,
			lsContains:    "container_test.go",
			lsNotContains: "tmp",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container, err := Setup(ctx, client, &tc.opts, "")
			assert.ErrorIs(t, err, tc.wantErr)
			if err != nil {
				// No need to continue on err
				return
			}

			// Get contents of current working directory in the container
			stdout, err := container.WithExec([]string{"ls", "."}).
				Stdout(ctx)
			assert.NoError(t, err)

			// Check the directory contents
			ls := strings.Split(stdout, "\n")
			assert.True(t, slices.Contains(ls, tc.lsContains))
			if tc.lsNotContains != "" {
				assert.True(t, !slices.Contains(ls, tc.lsNotContains))
			}
		})
	}
}

func TestGetArtifacts(t *testing.T) {
	// This test is rather slow (between 10s and 20s)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()
	filename := "my_file.txt"
	filename2 := "my_file2.txt"
	prefix := "my_directory"
	hostArtifactsDir := "artifacts"

	opts := SetupOpts{
		ContainerURL:      "ubuntu:latest",
		MountContainerDir: "/src",
		MountHostDir:      ".",
		WorkdirContainer:  "/src",
	}

	testCases := []struct {
		name            string
		cmdToRun        [][]string
		artifacts       []Artifacts
		wantErrExport   error
		filepathsToTest []string
		wantErrFile     error
	}{
		{
			name: "empty directory string",
			cmdToRun: [][]string{
				{"true"},
			},
			artifacts: []Artifacts{{
				ContainerPath: "",
				ContainerDir:  true,
				HostPath:      hostArtifactsDir,
				HostDir:       true,
			}},
			wantErrExport:   errDirectoryNotSpecified,
			filepathsToTest: []string{},
			wantErrFile:     os.ErrNotExist, // Does not get tested
		},
		{
			name: "file -> file",
			//	/src/my_file.txt -> /tmp/.../artifacts/my_file.txt
			cmdToRun: [][]string{
				{"touch", filename},
			},
			artifacts: []Artifacts{{
				ContainerPath: filepath.Join("/src", filename),
				ContainerDir:  false,
				HostPath:      filepath.Join(hostArtifactsDir, filename),
				HostDir:       false,
			}},
			wantErrExport: nil,
			filepathsToTest: []string{
				filepath.Join(hostArtifactsDir, filename),
			},
			wantErrFile: os.ErrExist,
		},
		{
			name: "file -> dir",
			//	/src/my_file.txt -> /tmp/.../artifacts/
			cmdToRun: [][]string{
				{"touch", filename},
			},
			artifacts: []Artifacts{{
				ContainerPath: filepath.Join("/src", filename),
				ContainerDir:  false,
				HostPath:      hostArtifactsDir,
				HostDir:       true,
			}},
			wantErrExport: nil,
			filepathsToTest: []string{
				filepath.Join(hostArtifactsDir, filename),
			},
			wantErrFile: os.ErrExist,
		},
		{
			name: "2x file -> dir",
			//	/src/my_file.txt  -> /tmp/.../artifacts/
			//	/src/my_file2.txt -> /tmp/.../artifacts/
			cmdToRun: [][]string{
				{"touch", filename},
				{"touch", filename2},
				{"ls", "-a1"},
			},
			artifacts: []Artifacts{
				{
					ContainerPath: filepath.Join("/src", filename),
					ContainerDir:  false,
					HostPath:      hostArtifactsDir,
					HostDir:       true,
				},
				{
					ContainerPath: filepath.Join("/src", filename2),
					ContainerDir:  false,
					HostPath:      hostArtifactsDir,
					HostDir:       true,
				},
			},
			wantErrExport: nil,
			filepathsToTest: []string{
				filepath.Join(hostArtifactsDir, filename),
				filepath.Join(hostArtifactsDir, filename2),
			},
			wantErrFile: os.ErrExist,
		},
		{
			name: "dir -> dir",
			//	/src/ -> /tmp/.../artifacts/src
			cmdToRun: [][]string{
				{"mkdir", prefix},
				{"touch", filepath.Join(prefix, filename)},
			},
			artifacts: []Artifacts{{
				ContainerPath: "/src",
				ContainerDir:  true,
				HostPath:      hostArtifactsDir,
				HostDir:       true,
			}},
			wantErrExport: nil,
			filepathsToTest: []string{
				filepath.Join(hostArtifactsDir, "src", prefix, filename),
			},
			wantErrFile: os.ErrExist,
		},
		{
			name: "export non-existing directory",
			cmdToRun: [][]string{
				{"true"},
			},
			artifacts: []Artifacts{{
				ContainerPath: "/some/non-existing/directory/path",
				ContainerDir:  true,
				HostPath:      hostArtifactsDir,
				HostDir:       true,
			}},
			wantErrExport: errExportFailed,
			filepathsToTest: []string{
				filepath.Join(prefix, filename),
			},
			wantErrFile: os.ErrNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			container, err := Setup(ctx, client, &opts, "")
			assert.NoError(t, err)

			// Run commands in container
			for _, cmd := range tc.cmdToRun {
				container, err = container.WithExec(cmd).
					Sync(ctx)
				assert.NoError(t, err)
			}

			// Extract artifacts
			for key := range tc.artifacts {
				tc.artifacts[key].HostPath = filepath.Join(tmpDir, tc.artifacts[key].HostPath)
			}
			err = GetArtifacts(ctx, container, &tc.artifacts)

			assert.ErrorIs(t, err, tc.wantErrExport)
			if err != nil {
				// No need to continue on err
				return
			}
			for _, file := range tc.filepathsToTest {
				assert.ErrorIs(t,
					filesystem.CheckFileExists(filepath.Join(tmpDir, file)),
					tc.wantErrFile)
			}
		})
	}
}
