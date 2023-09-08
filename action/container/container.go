// SPDX-License-Identifier: MIT

// Package container for dealing with containers via dagger
package container

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

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
	ContainerURL      string // URL or name of docker container (name-only will try to look for containers in ghcr.io/9elements/firmware-action
	MountHostDir      string // Directory from host to mount into container
	MountContainerDir string // Where to mount ^^^ host directory inside container
	WorkdirContainer  string // Workdir of the container, specified by GITHUB_WORKSPACE environment variable
}

// Setup for setting up a Docker container via dagger
func Setup(ctx context.Context, client *dagger.Client, opts *SetupOpts, dockerfileDirectoryPath string) (*dagger.Container, error) {
	// dockerfileDirectoryPath allows to use Dockerfile and build locally,
	//   which is handy for testing changes to said Dockerfile without the need to
	//   have the container uploaded into package registry

	// None of the directories can be empty string
	for _, val := range []string{opts.MountContainerDir, opts.MountHostDir, opts.WorkdirContainer} {
		if val == "" {
			return nil, errDirectoryNotSpecified
		}
	}

	// The mount target directory in container must not be root
	if opts.MountContainerDir == "." || opts.MountContainerDir == "/" {
		return nil, errDirectoryInvalid
	}

	// Setup container either from URL or build from Dockerfile
	var container *dagger.Container
	if dockerfileDirectoryPath == "" {
		// Use URL
		fmt.Println("Container setup: URL mode")

		// Make sure there is a non-empty URL or name provided
		if opts.ContainerURL == "" {
			return nil, errEmptyURL
		}

		// Pull docker container
		container = client.Container().From(opts.ContainerURL)
	} else {
		// Use Dockerfile
		fmt.Println("Container setup: Dockerfile mode")

		container = client.Container().Build(
			client.Host().Directory(dockerfileDirectoryPath),
		)
	}

	// Mount repository into the container
	return container.
		WithExec([]string{"mkdir", "-p", opts.MountContainerDir}).
		WithMountedDirectory(
			opts.MountContainerDir,
			client.Host().Directory(opts.MountHostDir)).
		WithWorkdir(opts.WorkdirContainer).
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
	ContainerPath string // Path inside container
	ContainerDir  bool   // Is ^^^ path directory?
	HostPath      string // Path inside host
}

// GetArtifacts extracts files from container to host
// Either both ContainerDir and HostDir must be directories, or both must be files
func GetArtifacts(ctx context.Context, container *dagger.Container, artifacts *[]Artifacts) error {
	for _, artifact := range *artifacts {
		if artifact.ContainerPath == "" || artifact.HostPath == "" {
			return errDirectoryNotSpecified
		}

		// Get reference to artifacts directory in the container
		var success bool
		var err error

		// Export
		// If allowParentDirPath is true, the path argument can be a directory path, in which case
		// the file will be created in that directory.
		if artifact.ContainerDir {
			// container side
			output := container.Directory(artifact.ContainerPath)
			// host side
			dirName := filepath.Base(artifact.ContainerPath)
			success, err = output.Export(ctx, filepath.Join(artifact.HostPath, dirName))
		} else {
			output := container.File(artifact.ContainerPath)
			success, err = output.Export(
				ctx,
				artifact.HostPath,
				dagger.FileExportOpts{AllowParentDirPath: true},
			)
		}

		// Copy contents of containers artifacts directory to host
		if err != nil || !success {
			return fmt.Errorf("%w: %w: %s -> %s", errExportFailed, err, artifact.ContainerPath, artifact.HostPath)
		}
	}

	return nil
}
