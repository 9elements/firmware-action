// SPDX-License-Identifier: MIT

// Package recipes / edk2
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// ANCHOR: Edk2Specific

// Edk2Specific is used to store data specific to coreboot.
//
//	simplified because of issue #92
type Edk2Specific struct {
	// Specifies which build command to use
	// GCC version is exposed in the container container as USE_GCC_VERSION environment variable
	// Examples:
	//   "source ./edksetup.sh; build -t GCC5 -a IA32 -p UefiPayloadPkg/UefiPayloadPkg.dsc"
	//   "python UefiPayloadPkg/UniversalPayloadBuild.py"
	//   "Intel/AlderLakeFspPkg/BuildFv.sh"
	BuildCommand string `json:"build_command" validate:"required"`
}

// ANCHOR_END: Edk2Specific

// ANCHOR: Edk2Opts

// Edk2Opts is used to store all data needed to build edk2.
type Edk2Opts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleLinux"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Specifies target architecture, such as 'x86' or 'arm64'. Currently unused for coreboot.
	// Supported options:
	//   - 'AARCH64'
	//   - 'ARM'
	//   - 'IA32'
	//   - 'IA32X64'
	//   - 'X64'
	Arch string `json:"arch"`

	// Gives the (relative) path to the defconfig that should be used to build the target.
	// For EDK2 this is a one-line file containing the build arguments such as
	//   '-D BOOTLOADER=COREBOOT -D TPM_ENABLE=TRUE -D NETWORK_IPXE=TRUE'.
	DefconfigPath string `json:"defconfig_path" validate:"filepath"`

	// Coreboot specific options
	Edk2Specific `validate:"required"`
}

// ANCHOR_END: Edk2Opts

// GetDepends is used to return list of dependencies
func (opts Edk2Opts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts Edk2Opts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// buildFirmware builds edk2 or Intel FSP
func (opts Edk2Opts) buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) (*dagger.Container, error) {
	envVars := map[string]string{
		"WORKSPACE":      ContainerWorkDir,
		"EDK_TOOLS_PATH": "/tools/Edk2/BaseTools",
	}

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

	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
		slog.Error(
			"Failed to start a container",
			slog.Any("error", err),
		)
		return nil, err
	}

	// Setup environment variables in the container
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Assemble build arguments
	//   and read content of the config file at "defconfig_path"
	var defconfigFileArgs []byte
	if opts.DefconfigPath != "" {
		if _, err := os.Stat(opts.DefconfigPath); !errors.Is(err, os.ErrNotExist) {
			defconfigFileArgs, err = os.ReadFile(opts.DefconfigPath)
			if err != nil {
				return nil, err
			}
		} else {
			slog.Warn(
				fmt.Sprintf("Failed to read file '%s' as defconfig_path: file does not exist", opts.DefconfigPath),
				slog.String("suggestion", "Double check the path for defconfig"),
				slog.Any("error", err),
			)
		}
	}

	// Assemble commands to build
	buildSteps := [][]string{}
	if !(runtime.GOARCH == "386" || runtime.GOARCH == "amd64") {
		// On all non-x86 architectures we have to also build the BaseTools
		// Docs: https://go.dev/doc/install/source#environment
		buildSteps = append(buildSteps, []string{"bash", "-c", "cd ${TOOLSDIR}/Edk2/; make -C BaseTools/ -j $(nproc)"})
	}
	buildSteps = append(buildSteps, []string{"bash", "-c", fmt.Sprintf("%s %s", opts.BuildCommand, string(defconfigFileArgs))})

	// Build
	var myContainerPrevious *dagger.Container
	for step := range buildSteps {
		myContainerPrevious = myContainer
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build edk2",
				slog.Any("error", err),
			)
			return myContainerPrevious, fmt.Errorf("edk2 build failed: %w", err)
		}
	}

	// Extract artifacts
	return myContainer, container.GetArtifacts(ctx, myContainer, opts.GetArtifacts())
}
