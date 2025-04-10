// SPDX-License-Identifier: MIT

//go:build go1.24

// Package container
package container

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
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
		name              string
		opts              SetupOpts
		wantErr           error
		TestFiles         []string
		TestDirs          []string
		InputDirsPopulate []string
	}{
		{
			name: "empty URL",
			opts: SetupOpts{
				ContainerURL:      "",
				MountContainerDir: "/src",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  "/src",
			},
			wantErr:           errEmptyURL,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "empty directory strings",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "",
				MountHostDir:      "",
				WorkdirContainer:  "",
			},
			wantErr:           errDirectoryNotSpecified,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "invalid directory strings: .",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: ".",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  ".",
			},
			wantErr:           errDirectoryInvalid,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "invalid directory strings: /",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  ".",
			},
			wantErr:           errDirectoryInvalid,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "InputFiles without InputDir",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/src",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  "/src",
				InputFiles:        []string{"test.img"},
			},
			wantErr:           errDirectoryNotSpecified,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "InputFiles without InputDir",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/src",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  "/src",
				InputDirs:         []string{"test/"},
			},
			wantErr:           errDirectoryNotSpecified,
			TestDirs:          []string{},
			TestFiles:         []string{},
			InputDirsPopulate: []string{},
		},
		{
			name: "valid inputs",
			opts: SetupOpts{
				ContainerURL:      "ubuntu:latest",
				MountContainerDir: "/src",
				MountHostDir:      t.TempDir(),
				WorkdirContainer:  "/src",
				InputDirs: []string{
					"test-dir/",
				},
				InputFiles: []string{
					"test-file.img",
					"test-file.txt",
				},
				ContainerInputDir: "inputs",
			},
			wantErr: nil,
			TestDirs: []string{
				"inputs",
			},
			TestFiles: []string{
				"inputs/test-file.img",
				"inputs/test-file.txt",
				"inputs/test-dir/test-in-dir-file.txt",
			},
			InputDirsPopulate: []string{
				"test-dir/test-in-dir-file.txt",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test files
			tmpDir := t.TempDir()
			t.Chdir(tmpDir)

			// Create InputDirs
			for _, val := range tc.opts.InputDirs {
				err = os.MkdirAll(filepath.Join(tmpDir, val), os.ModePerm)
				assert.NoError(t, err)
			}
			// Create InputFiles
			for _, val := range tc.opts.InputFiles {
				f, err := os.Create(filepath.Join(tmpDir, val))
				assert.NoError(t, err)
				f.Close()
			}
			// Populate InputDirs
			for _, val := range tc.InputDirsPopulate {
				err = os.MkdirAll(filepath.Join(tmpDir, filepath.Dir(val)), os.ModePerm)
				assert.NoError(t, err)

				f, err := os.Create(filepath.Join(tmpDir, val))
				assert.NoError(t, err)
				f.Close()
			}

			// Spin up container
			container, err := Setup(ctx, client, &tc.opts)
			assert.ErrorIs(t, err, tc.wantErr)
			if err != nil {
				// No need to continue on err
				return
			}

			// Check the directory contents
			for _, val := range tc.TestDirs {
				_, err = container.WithExec([]string{"bash", "-c", fmt.Sprintf("[ -d %s ]", val)}).
					Sync(ctx)
				assert.NoError(t, err, fmt.Sprintf("Directory '%s' does not exists", val))
			}
			for _, val := range tc.TestFiles {
				_, err = container.WithExec([]string{"bash", "-c", fmt.Sprintf("[ -f %s ]", val)}).
					Sync(ctx)
				assert.NoError(t, err, fmt.Sprintf("File '%s' does not exists", val))
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

			container, err := Setup(ctx, client, &opts)
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

func TestCheckIfDiscontinued(t *testing.T) {
	testCases := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "not discontinued",
			url:     "ghcr.io/9elements/firmware-action/edk2-stable202408.01:main",
			wantErr: nil,
		},
		{
			name:    "discontinued main github",
			url:     "ghcr.io/9elements/firmware-action/edk2-stable202408:main",
			wantErr: errContainerDiscontinued,
		},
		{
			name:    "discontinued main dockerhub",
			url:     "docker.io/9elementscyberops/edk2-stable202408:main",
			wantErr: errContainerDiscontinued,
		},
		{
			name:    "discontinued latest",
			url:     "docker.io/9elementscyberops/edk2-stable202408:latest",
			wantErr: errContainerDiscontinued,
		},
		{
			name:    "discontinued tagged",
			url:     "docker.io/9elementscyberops/edk2-stable202408:v0.15.0",
			wantErr: errContainerDiscontinued,
		},
		{
			name:    "discontinued short",
			url:     "9elementscyberops/edk2-stable202408:v0.15.0",
			wantErr: errContainerDiscontinued,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, CheckIfDiscontinued(tc.url), tc.wantErr)
		})
	}
}
