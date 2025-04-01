// SPDX-License-Identifier: MIT

// Package container for dealing with containers via dagger
package container

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/logging"
)

var (
	errEmptyURL              = errors.New("invalid docker URL")
	errDirectoryNotSpecified = errors.New("empty string for directory path was given")
	errDirectoryInvalid      = errors.New("host directory cannot be mounted into '/' or '.' in the container")
	errExportFailed          = errors.New("failed to export artifacts from container")
	errContainerDiscontinued = errors.New("the used container is discontinued")
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
func Setup(ctx context.Context, client *dagger.Client, opts *SetupOpts) (*dagger.Container, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	// dockerfileDirectoryPath allows to use Dockerfile and build locally,
	//   which is handy for testing changes to said Dockerfile without the need to
	//   have the container uploaded into package registry
	dockerfileDirectoryPath := ""
	dockerfilePathPattern := regexp.MustCompile(`^file:\/\/.*`)
	if dockerfilePathPattern.MatchString(opts.ContainerURL) {
		// opts.ContainerURL is actually filepath
		dockerfileDockerfilePattern := regexp.MustCompile(`.*\/Dockerfile$`)
		pathPattern := regexp.MustCompile(`^file:\/\/`)
		dockerfileDirectoryPath = pathPattern.ReplaceAllString(opts.ContainerURL, "")

		// Docker requires to use directory, if path contains also Dockerfile as last element, remove it
		// to get the base directory
		if dockerfileDockerfilePattern.MatchString(opts.ContainerURL) {
			dockerfileDirectoryPath = filepath.Dir(dockerfileDirectoryPath)
		}
	} else {
		_ = CheckIfDiscontinued(opts.ContainerURL)
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

// CleanupAfterContainer performs cleanup operations after container use
func CleanupAfterContainer(ctx context.Context) error {
	// Unfortunately it is not possible to only remove the container used for building the module.
	// Dagger Engine somehow absorbs the other containers into itself (possibly into it's volume, not sure).
	// So to actually free up a disk space by deleting a container we have to delete the whole dagger engine
	//   container and it's volume.
	//
	// This function is used to free up disk space on constrained environments like GitHub Actions.
	//   GitHub-hosted public runners have only 14GB of disk space available.
	// If user wants to build complex firmware stacks in single job recursively, they will easily run
	//   out of disk space.
	//
	// WARNING: This will completely stop the Dagger engine. Any subsequent Dagger
	//   operations will need to reinitialize the Dagger client.

	slog.Info("Cleaning up Dagger container resources")

	// Step 1: Find the Dagger engine container
	findCmd := exec.CommandContext(ctx, "docker", "container", "ls", "--filter", "name=dagger-engine", "--format", "{{.ID}}")
	containerID, err := findCmd.Output()
	if err != nil {
		slog.Error(
			"Failed to find Dagger engine container",
			slog.Any("error", err),
		)
		return err
	}

	containerIDStr := strings.TrimSpace(string(containerID))
	if containerIDStr == "" {
		slog.Info("No Dagger engine container found to clean up")
		return nil
	}

	// Step 2: Stop the Dagger engine container
	slog.Debug(
		"Stopping Dagger engine container",
		slog.String("containerID", containerIDStr),
	)
	stopCmd := exec.CommandContext(ctx, "docker", "container", "stop", containerIDStr)
	stopOutput, err := stopCmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"Failed to stop Dagger engine container",
			slog.String("output", strings.TrimSpace(string(stopOutput))),
			slog.Any("error", err),
		)
		return err
	}

	// Step 3: Remove the Dagger engine container
	slog.Debug(
		"Removing Dagger engine container",
		slog.String("containerID", containerIDStr),
	)
	rmCmd := exec.CommandContext(ctx, "docker", "container", "rm", containerIDStr)
	rmOutput, err := rmCmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"Failed to remove Dagger engine container",
			slog.String("output", strings.TrimSpace(string(rmOutput))),
			slog.Any("error", err),
		)
		return err
	}

	// Step 4: Find and remove Dagger volumes
	volCmd := exec.CommandContext(ctx, "docker", "volume", "ls", "--filter", "dangling=true", "--format", "{{.Name}}")
	volumes, err := volCmd.Output()
	if err != nil {
		slog.Warn(
			"Failed to list Docker volumes",
			slog.Any("error", err),
		)
		// Continue even if this fails
	} else {
		volumeList := strings.Split(strings.TrimSpace(string(volumes)), "\n")
		for _, vol := range volumeList {
			if vol == "" {
				continue
			}
			slog.Debug(
				"Removing Docker volume",
				slog.String("volume", vol),
			)
			rmVolCmd := exec.CommandContext(ctx, "docker", "volume", "rm", vol)
			rmVolOutput, err := rmVolCmd.CombinedOutput()
			if err != nil {
				slog.Warn(
					"Failed to remove Docker volume",
					slog.String("volume", vol),
					slog.String("output", strings.TrimSpace(string(rmVolOutput))),
					slog.Any("error", err),
				)
				// Continue with other volumes even if one fails
			}
		}
	}

	// Step 5: Run system prune to clean up any remaining resources
	pruneCmd := exec.CommandContext(ctx, "docker", "system", "prune", "-f")
	pruneOutput, err := pruneCmd.CombinedOutput()
	slog.Debug(
		"Docker system prune output",
		slog.String("command", "docker system prune -f"),
		slog.String("output", strings.TrimSpace(string(pruneOutput))),
	)
	if err != nil {
		slog.Error(
			"Failed to prune Docker system",
			slog.String("output", strings.TrimSpace(string(pruneOutput))),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("Dagger container resources cleaned up successfully")
	return nil
}

// CheckIfDiscontinued prints a warning if the used container is discontinued
func CheckIfDiscontinued(containerURL string) error {
	// Check for discontinued containers
	listDiscontinued := []string{
		`.*ghcr\.io\/9elements\/coreboot:4\.19.*`,
		`.*ghcr\.io\/9elements\/uefi:edk\-stable202208.*`,
	}

	// Patterns to use with 'containerHub'
	listDiscontinuedPattern := []string{
		`coreboot_24\.02(:.*)?$`,
		`coreboot_4\.20(:.*)?$`,
		`coreboot_4\.22(:.*)?$`,
		`edk2\-stable202408(:.*)?$`,
		`linux_6\.11(:.*)?$`,
		`linux_6\.1\.11(:.*)?$`,
		`linux_6\.1\.45(:.*)?$`,
		`linux_6\.6\.52(:.*)?$`,
		`linux_6\.9\.9(:.*)?$`,
	}
	// The containers can be in multiple container hubs (GitHub, DockerHub, ...)
	containerHub := []string{
		`.*ghcr\.io\/9elements\/firmware\-action\/`,
		`.*docker\.io\/9elementscyberops\/`,
		`.*9elementscyberops\/`,
		// '9elementscyberops' is short for 'docker.io/9elementscyberops'
	}
	// For each hub, add each pattern
	for _, hub := range containerHub {
		for _, pattern := range listDiscontinuedPattern {
			listDiscontinued = append(listDiscontinued, hub+pattern)
		}
	}

	// Iterate over all patterns and check to match
	for _, discontinued := range listDiscontinued {
		pattern := regexp.MustCompile(discontinued)
		if pattern.MatchString(containerURL) {
			slog.Warn(
				"Using discontinued container",
				slog.String("suggestion", "The container will remain available, but will no longer receive any bug-fixes or updates. If you want maintained and up-to-date container, look at https://github.com/9elements/firmware-action#containers"),
			)
			return errContainerDiscontinued
		}
	}

	return nil
}
