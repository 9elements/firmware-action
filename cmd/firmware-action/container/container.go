// SPDX-License-Identifier: MIT

// Package container for dealing with containers via dagger
package container

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/logging"
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
	ContainerURL      string   // URL or name of docker container
	MountHostDir      string   // Directory from host to mount into container
	MountContainerDir string   // Where to mount ^^^ host directory inside container
	WorkdirContainer  string   // Workdir of the container, specified by GITHUB_WORKSPACE environment variable
	ContainerInputDir string   // Directory for input files
	InputDirs         []string // List of directories to copy into container
	InputFiles        []string // List of files to copy into container
}

// Validate the data in struct
func (opts SetupOpts) Validate() error {
	// None of the directories can be empty string
	var err error
	if opts.MountContainerDir == "" {
		err = errors.Join(err, errDirectoryNotSpecified)
		slog.Error(
			"Mountpoint path cannot be empty string",
			slog.String("suggestion", "Specify where the host directory should be mounted in the container"),
			slog.Any("error", err),
		)
	}
	if opts.MountHostDir == "" {
		err = errors.Join(err, errDirectoryNotSpecified)
		slog.Error(
			"Host directory path for mounting cannot be empty string",
			slog.String("suggestion", "Specify which host directory will be mounted into the container"),
			slog.Any("error", err),
		)
	}
	if opts.WorkdirContainer == "" {
		err = errors.Join(err, errDirectoryNotSpecified)
		slog.Error(
			"WORKDIR cannot be empty string",
			slog.String("suggestion", "Specify working directory for the container"),
			slog.Any("error", err),
		)
	}

	// The mount target directory in container must not be root
	if opts.MountContainerDir == "." || opts.MountContainerDir == "/" {
		err = errors.Join(err, errDirectoryInvalid)
		slog.Error(
			"Container mountpoint cannot be '.' or '/'",
			slog.String("suggestion", "Pick another directory, preferably absolute path"),
			slog.Any("error", err),
		)
	}

	// If any input file or directory specified, inputDir must be defined
	if (len(opts.InputDirs) > 0 || len(opts.InputFiles) > 0) && opts.ContainerInputDir == "" {
		err = errors.Join(err, errDirectoryNotSpecified)
		slog.Error(
			"Container InputDir cannot be empty string when using InputFiles and/or InputDirs",
			slog.String("suggestion", "Specify directory for input files and directories"),
			slog.Any("error", err),
		)
	}
	return err
}

// Setup for setting up a Docker container via dagger
func Setup(ctx context.Context, client *dagger.Client, opts *SetupOpts, dockerfileDirectoryPath string) (*dagger.Container, error) {
	// dockerfileDirectoryPath allows to use Dockerfile and build locally,
	//   which is handy for testing changes to said Dockerfile without the need to
	//   have the container uploaded into package registry

	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	// Setup container either from URL or build from Dockerfile
	var container *dagger.Container
	if dockerfileDirectoryPath == "" {
		// Use URL
		slog.Info("Container setup running in URL mode")

		// Make sure there is a non-empty URL or name provided
		if opts.ContainerURL == "" {
			slog.Error(
				"Container setup was provided with empty URL",
				slog.String("suggestion", "Provide URL or Dockerfile"),
				slog.Any("error", errEmptyURL),
			)
			return nil, errEmptyURL
		}

		// Pull docker container
		container = client.Container().From(opts.ContainerURL)
		imageRef, _ := container.ImageRef(ctx)
		slog.Info(
			"Container information",
			slog.String("Image reference", imageRef),
		)
	} else {
		// Use Dockerfile
		slog.Info("Container setup running in Dockerfile mode")

		container = client.Container().Build(
			client.Host().Directory(dockerfileDirectoryPath),
		)
	}

	// Mount repository into the container
	//   WithDirectory
	//     Copy files from host to container
	//     Creates directory tree if needed
	//   WithMountedDirectory
	//     Create a OverlayFS with bottom layer Read-Only
	//     Directory in container must exist
	container = container.
		WithExec([]string{"mkdir", "-p", opts.MountContainerDir}).
		WithMountedDirectory(
			opts.MountContainerDir,
			client.Host().Directory(opts.MountHostDir)).
		WithWorkdir(opts.WorkdirContainer)

	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return nil, err
	}

	// Make input directory
	inputDirPath := filepath.Join(opts.WorkdirContainer, opts.ContainerInputDir)
	container = container.WithExec([]string{"mkdir", "-p", inputDirPath})

	// Copy input directories into the container
	// We cannot do nested WithMountedDirectory, it silently breaks
	for _, val := range opts.InputDirs {
		container = container.
			WithExec([]string{"mkdir", "-p", filepath.Join(inputDirPath, filepath.Base(val))}).
			WithDirectory(
				filepath.Join(inputDirPath, filepath.Base(val)),
				client.Host().Directory(filepath.Join(pwd, val)),
			)
	}

	// Copy over input files
	for _, val := range opts.InputFiles {
		container = container.
			WithFile(
				filepath.Join(inputDirPath, filepath.Base(val)),
				client.Host().File(filepath.Join(pwd, val)),
			)
	}

	container, err = container.Sync(ctx)
	if err != nil {
		message := "Failed to spin up a container"
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error(
				message,
				slog.String("suggestion", "Your network configuration likely changed, try this: https://docs.dagger.io/troubleshooting#dagger-is-unable-to-resolve-host-names-after-network-configuration-changes"),
				slog.Any("error", err),
			)
		}
		if strings.Contains(err.Error(), "failed to do request") && strings.Contains(err.Error(), "i/o timeout") {
			slog.Error(
				message,
				slog.String("suggestion", "try this: https://archive.docs.dagger.io/0.9/235290/troubleshooting/#dagger-pipeline-is-unable-to-resolve-host-names-after-network-configuration-changes"),
				slog.Any("error", err),
			)
		}
		if strings.Contains(err.Error(), "timed out waiting for session params") && runtime.GOOS == "linux" {
			// On Linux, check if 'iptable_nat' kernel module is loaded
			content, err := os.ReadFile("/proc/modules")
			if err != nil {
				pattern := regexp.MustCompile(`^iptable_nat`)
				if pattern.FindString(string(content)) == "" {
					slog.Error(
						message,
						slog.String("suggestion", "dagger requires the 'iptable_nat' Linux kernel module in order to function properly, https://docs.dagger.io/troubleshooting#dagger-restarts-with-a-cni-setup-error"),
						slog.Any("error", err),
					)
				}
			}
		}
		slog.Error(
			message,
			slog.String("suggestion", "something is wrong with dagger, please check dagger troubleshooting guide at: https://docs.dagger.io/troubleshooting"),
			slog.Any("error", err),
		)
	}

	return container, err
}

// Artifacts is passes to GetArtifacts as argument, and specifies extraction of files
// form container at containerDir to host at hostDir
type Artifacts struct {
	ContainerPath string // Path inside container
	ContainerDir  bool   // Is ^^^ path directory?
	HostPath      string // Path inside host
	HostDir       bool   // Is ^^^ path directory?
}

// GetArtifacts extracts files from container to host
// Either both ContainerDir and HostDir must be directories, or both must be files
func GetArtifacts(ctx context.Context, container *dagger.Container, artifacts *[]Artifacts) error {
	for _, artifact := range *artifacts {
		if artifact.ContainerPath == "" || artifact.HostPath == "" {
			return errDirectoryNotSpecified
		}

		// Get reference to artifacts directory in the container
		var err error

		if artifact.HostDir {
			if err := os.MkdirAll(artifact.HostPath, 0o755); err != nil {
				return err
			}
		}

		// Export
		// If AllowParentDirPath is true, the path argument can be a directory path, in which case
		// the file will be created in that directory.
		if artifact.ContainerDir {
			// container side
			output := container.Directory(artifact.ContainerPath)
			// host side
			dirName := filepath.Base(artifact.ContainerPath)
			_, err = output.Export(ctx, filepath.Join(artifact.HostPath, dirName))
		} else {
			output := container.File(artifact.ContainerPath)
			_, err = output.Export(
				ctx,
				artifact.HostPath,
				dagger.FileExportOpts{AllowParentDirPath: true},
			)
		}

		// Copy contents of containers artifacts directory to host
		if err != nil {
			return fmt.Errorf("%w: %w: %s -> %s", errExportFailed, err, artifact.ContainerPath, artifact.HostPath)
		}
		slog.Debug(fmt.Sprintf("Artifact export: %s -> %s", artifact.ContainerPath, artifact.HostPath))
	}

	return nil
}
