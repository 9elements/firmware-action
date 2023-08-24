// SPDX-License-Identifier: MIT

// Package container for dealing with containers via dagger
package container

import (
	"context"
	"errors"

	"dagger.io/dagger"
)

var (
	errEmptyURL              = errors.New("invalid docker URL")
	errDirectoryNotSpecified = errors.New("empty string for directory path was given")
	errDirectoryInvalid      = errors.New("host directory cannot be mounted into '/' or '.' in the container")
	errExportFailed          = errors.New("failed to export artifacts from container")
)

// SetupOpts congregates options for Setup function
// None of the values can be empty string, and mountContainerDir cannot be '.' or '/'
type SetupOpts struct {
	containerURL      string
	mountContainerDir string
	mountHostDir      string
	workdirContainer  string
}

// Setup for setting up a Docker container via dagger
func Setup(ctx context.Context, client *dagger.Client, opts *SetupOpts) (*dagger.Container, error) {
	// Make sure there is a non-empty URL or name provided
	if opts.containerURL == "" {
		return nil, errEmptyURL
	}

	// None of the directories can be empty string
	for _, val := range []string{opts.mountContainerDir, opts.mountHostDir, opts.workdirContainer} {
		if val == "" {
			return nil, errDirectoryNotSpecified
		}
	}

	// The mount target directory in container must not be root
	if opts.mountContainerDir == "." || opts.mountContainerDir == "/" {
		return nil, errDirectoryInvalid
	}

	// Pull docker container
	container := client.Container().From(opts.containerURL)

	// Mount repository into the container
	return container.
		WithExec([]string{"mkdir", "-p", opts.mountContainerDir}).
		WithMountedDirectory(
			opts.mountContainerDir,
			client.Host().Directory(opts.mountHostDir)).
		WithWorkdir(opts.workdirContainer).
		Sync(ctx)
	// WithDirectory
	//	Copy files from host to container
	//	Creates directory tree if needed
	// WithMountedDirectory
	//	Create a OverlayFS with bottom layer Read-Only
	//	Directory in container must exist
}

// Artifacts is passes to GetArtifacts as argument, and specifies extraction of files
// form container at containerDir to host at hostDir
type Artifacts struct {
	containerDir string
	hostDir      string
}

// GetArtifacts extracts files from container to host
func GetArtifacts(ctx context.Context, container *dagger.Container, artifacts *Artifacts) error {
	if artifacts.containerDir == "" || artifacts.hostDir == "" {
		return errDirectoryNotSpecified
	}

	// Get reference to artifacts directory in the container
	output := container.Directory(artifacts.containerDir)

	// Copy contents of containers artifacts directory to host
	success, err := output.Export(ctx, artifacts.hostDir)
	if err != nil || !success {
		return errExportFailed
	}

	return nil
}
