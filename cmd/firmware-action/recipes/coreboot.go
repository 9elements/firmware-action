// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/container"
	"github.com/9elements/firmware-action/cmd/firmware-action/environment"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/9elements/firmware-action/cmd/firmware-action/logging"
)

// BlobDef is used to store information about a single blob.
// This structure is not exposed to the user, it is filled in automatically based on user input.
type BlobDef struct {
	// Path to the blob (either file or directory)
	Path string

	// Blobs get renamed when moved to this string
	DestinationFilename string

	// Kconfig key specifying the filepath to the blob in defconfig
	KconfigKey string
}

// ANCHOR: CorebootOpts

// CorebootOpts is used to store all data needed to build coreboot.
type CorebootOpts struct {
	// List of IDs this instance depends on
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Gives the (relative) path to the defconfig that should be used to build the target.
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`

	// Blobs
	// The blobs will be copied into the container into directory:
	//   3rdparty/blobs/mainboard/${CONFIG_MAINBOARD_DIR}/
	// And the blobs will remain their name
	// NOTE: The blobs may not be added to the ROM, depends on provided defconfig.
	// Example:
	//   Config:
	//     "CONFIG_PAYLOAD_FILE": "./my-payload.bin"
	//   Will result in blob "my-payload.bin" at
	//     "3rdparty/blobs/mainboard/${CONFIG_MAINBOARD_DIR}/my-payload.bin"
	Blobs map[string]string `json:"blobs"`
}

// ANCHOR_END: CorebootOpts

// GetDepends is used to return list of dependencies
func (opts CorebootOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts CorebootOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// ANCHOR: CorebootOptsGetSources

// GetSources returns slice of paths to all sources which are used for build
func (opts CorebootOpts) GetSources() []string {
	sources := opts.CommonOpts.GetSources()

	// Add DefconfigPath to list of sources
	sources = append(sources, opts.DefconfigPath)

	// Add blobs to list of sources
	blobs, err := opts.ProcessBlobs()
	if err != nil {
		slog.Error(
			"Failed to process all blobs",
			slog.Any("error", err),
		)
		return nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return nil
	}
	for blob := range blobs {
		// Path to local file on host
		src := filepath.Join(
			pwd,
			blobs[blob].Path,
		)
		sources = append(sources, src)
	}

	return sources
}

// ANCHOR_END: CorebootOptsGetSources

// ProcessBlobs is used to figure out blobs from provided data
func (opts CorebootOpts) ProcessBlobs() ([]BlobDef, error) {
	blobs := []BlobDef{}

	for key, value := range opts.Blobs {
		if key != "" && value != "" {
			newBlob := BlobDef{
				KconfigKey: key,
				Path:       value,
				// Blobs get renamed when moved to this string
				DestinationFilename: filepath.Base(value),
			}

			blobs = append(blobs, newBlob)
		}
	}
	return blobs, nil
}

// buildFirmware builds coreboot with all blobs and stuff
func (opts CorebootOpts) buildFirmware(ctx context.Context, client *dagger.Client) error {
	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      opts.SdkURL,
		MountContainerDir: ContainerWorkDir,
		MountHostDir:      opts.RepoPath,
		WorkdirContainer:  ContainerWorkDir,
		ContainerInputDir: opts.ContainerInputDir,
		InputDirs:         opts.InputDirs,
		InputFiles:        opts.InputFiles,
	}
	myContainer, err := container.Setup(ctx, client, &containerOpts)
	if err != nil {
		slog.Error(
			"Failed to start a container",
			slog.Any("error", err),
		)
		return err
	}

	// Copy over the defconfig file
	defconfigBasename := filepath.Base(opts.DefconfigPath)
	//   not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(ContainerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, opts.DefconfigPath)),
	)

	// Get value of CONFIG_MAINBOARD_DIR / MAINBOARD_DIR variable from dotconfig
	//   to extract value of 'CONFIG_MAINBOARD_DIR', there must be '.config'
	generateDotConfigCmd := []string{"make", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename), "defconfig"}
	mainboardDir, err := myContainer.
		WithExec(generateDotConfigCmd).
		WithExec([]string{"./util/scripts/config", "-s", "CONFIG_MAINBOARD_DIR"}).
		Stdout(ctx)
	if err != nil {
		slog.Error(
			"Failed to get value of MAINBOARD_DIR from .config",
			slog.Any("error", err),
		)
		return err
	}
	//   strip newline from mainboardDir
	mainboardDir = strings.ReplaceAll(mainboardDir, "\n", "")

	// Assemble commands to build
	buildSteps := [][]string{
		// remove existing config if exists
		// -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// generate dotconfig from defconfig
		generateDotConfigCmd,
	}

	// Handle blobs
	// Firstly copy all the blobs into building container.
	// Then use './util/scripts/config' script in coreboot repository to update configuration
	//   options for said blobs (this must run inside container).
	blobs, err := opts.ProcessBlobs()
	if err != nil {
		slog.Error(
			"Failed to process all blobs",
			slog.Any("error", err),
		)
		return err
	}
	for blob := range blobs {
		// Path to local file on host
		src := filepath.Join(
			pwd,
			blobs[blob].Path,
		)
		// Path to file in container
		dst := filepath.Join(
			filepath.Join("3rdparty/blobs/mainboard", mainboardDir),
			blobs[blob].DestinationFilename,
		)

		// Copy into container
		err = filesystem.CheckFileExists(src)
		if !errors.Is(err, os.ErrExist) {
			slog.Error(
				fmt.Sprintf("Blob '%s' was not found", src),
				slog.String("suggestion", "blobs are copied into container separately from 'input_files' and 'input_dirs', the path should point to files on your host"),
				slog.Any("error", err),
			)
			return err
		}
		if errors.Is(err, filesystem.ErrPathIsDirectory) {
			// Directory
			slog.Info(fmt.Sprintf("Copying directory '%s' to container at '%s'", src, dst))
			myContainer = myContainer.WithExec([]string{"mkdir", "-p", dst})
			// myContainer = myContainer.WithMountedDirectory(
			// can't use WithMountedDirectory because the repo (aka working directory)
			//   is already mounted with WithMountedDirectory
			//   this nesting causes problems
			myContainer = myContainer.WithDirectory(
				dst,
				client.Host().Directory(src),
			)
		} else if errors.Is(err, os.ErrExist) {
			// File
			myContainer = myContainer.WithFile(
				dst,
				client.Host().File(src),
			)
		} else {
			// Something is wrong, blob is not file nor directory
			slog.Error(
				"Failed to process a blob in coreboot configuration",
				slog.String("suggestion", "Please double-check blobs, and make sure the paths point to actual file or directory"),
				slog.String("blob_path", blobs[blob].Path),
				slog.String("blob_kconfig", blobs[blob].KconfigKey),
				slog.Any("error", err),
			)
			return err
		}

		// Fix defconfig
		buildSteps = append(
			buildSteps,
			// update coreboot config value related to blob to actual path of the blob
			[]string{"./util/scripts/config", "--set-str", blobs[blob].KconfigKey, dst},
		)
	}

	buildSteps = append(
		buildSteps,
		// compile
		[]string{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		[]string{"make", "savedefconfig"},
	)

	// Setup environment variables in the container
	// envVars := map[string]string{}
	envVars, err := corebootPassEnvVars(opts.RepoPath)
	if err != nil {
		slog.Error(
			"Failed to extract environment variables from current environment",
			slog.Any("error", err),
		)
		return fmt.Errorf("coreboot build failed: %w", err)
	}
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Build
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build coreboot",
				slog.Any("error", err),
			)
			return fmt.Errorf("coreboot build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, opts.CommonOpts.GetArtifacts())
}

func corebootPassEnvVars(repoPath string) (map[string]string, error) {
	passVariables := []string{"KERNELVERSION", "BUILD_TIMELESS"}
	envVariables := environment.FetchEnvVars(passVariables)

	// coreboot build system takes a version from:
	// - environment variable: KERNELVERSION
	// - shell command: git describe ...
	// - content of file: .coreboot-version

	// To check for coreboot version in compiled binary, run these commands:
	//   $ cbfstool build/coreboot.rom extract -n build_info -f /tmp/foo
	//   $ grep COREBOOT_VERSION /tmp/foo

	// coreboot make will fail to run 'git describe' because of
	//   missing '.git' directory once the content of repoPath
	//   is copied into the container
	// To fix this we need to run 'git describe' now and create a new
	//   environment variable to pass over into the container
	// This way, the compiled coreboot binary will not have unknown version

	// If KERNELVERSION is defined in current environment, do nothing
	if _, ok := envVariables["KERNELVERSION"]; ok {
		return envVariables, nil
	}

	// If .coreboot-version file exists in coreboot directory, do nothing
	corebootVersionPath := filepath.Join(repoPath, ".coreboot-version")
	err := filesystem.CheckFileExists(corebootVersionPath)
	if errors.Is(err, os.ErrExist) {
		return envVariables, nil
	}

	// At this point we checked that user did not define their own coreboot version
	// coreboot build system would at this point attempt to run git describe, which would fail
	// Define a new environment variable KERNELVERSION with value from git describe
	//   and then pass it into the container
	err = filesystem.CheckFileExists(repoPath)
	if errors.Is(err, filesystem.ErrPathIsDirectory) {
		describe, err := filesystem.GitDescribeCoreboot(repoPath)
		if err != nil {
			return nil, err
		}
		envVariables["KERNELVERSION"] = describe
	}

	return envVariables, nil
}
