// SPDX-License-Identifier: MIT

// Package recipes / edk2
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// Edk2Specific is used to store data specific to coreboot.
/* TODO: removed because of issue #92
type Edk2Specific struct {
	// Gives the (relative) path to the defconfig that should be used to build the target.
	// For EDK2 this is a one-line file containing the build arguments such as
	//   '-D BOOTLOADER=COREBOOT -D TPM_ENABLE=TRUE -D NETWORK_IPXE=TRUE'.
	//   Some arguments will be added automatically:
	//     '-a <architecture>'
	//     '-p <edk2__platform>'
	//     '-b <edk2__release_type>'
	//     '-t <GCC version>' (defined as part of docker toolchain, selected by SdkURL)
	DefconfigPath string `json:"defconfig_path" validate:"filepath"`

	// Specifies the DSC to use when building EDK2
	// Example: UefiPayloadPkg/UefiPayloadPkg.dsc
	Platform string `json:"platform" validate:"filepath"`

	// Specifies the build type to use when building EDK2
	// Supported options: DEBUG, RELEASE
	ReleaseType string `json:"release_type" validate:"required"`

	// Specifies which build command to use
	// Examples:
	//   "source ./edksetup.sh; build"
	//   "python UefiPayloadPkg/UniversalPayloadBuild.py"
	//   "Intel/AlderLakeFspPkg/BuildFv.sh"
	BuildCommand string `json:"build_command" validate:"required"`
}
*/
// ANCHOR: Edk2Specific
// Edk2Specific is used to store data specific to coreboot.
type Edk2Specific struct {
	// Specifies which build command to use
	// GCC version is exposed in the docker container as USE_GCC_VERSION environment variable
	// Examples:
	//   "source ./edksetup.sh; build -t GCC5 -a IA32 -p UefiPayloadPkg/UefiPayloadPkg.dsc"
	//   "python UefiPayloadPkg/UniversalPayloadBuild.py"
	//   "Intel/AlderLakeFspPkg/BuildFv.sh"
	BuildCommand string `json:"build_command" validate:"required"`
}

// ANCHOR_END: Edk2Specific

// Edk2Opts is used to store all data needed to build edk2.
type Edk2Opts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleLinux"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Gives the (relative) path to the defconfig that should be used to build the target.
	// For EDK2 this is a one-line file containing the build arguments such as
	//   '-D BOOTLOADER=COREBOOT -D TPM_ENABLE=TRUE -D NETWORK_IPXE=TRUE'.
	DefconfigPath string `json:"defconfig_path" validate:"filepath"`

	// Coreboot specific options
	Edk2Specific `validate:"required"`
}

// GetDepends is used to return list of dependencies
func (opts Edk2Opts) GetDepends() []string {
	return opts.Depends
}

// edk2 builds edk2
func edk2(ctx context.Context, client *dagger.Client, opts *Edk2Opts, dockerfileDirectoryPath string, artifacts *[]container.Artifacts) error {
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
	}

	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
		return err
	}

	// Setup environment variables in the container
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Get GCC version from environment variable
	/* TODO: removed because of issue #92
	gccVersion, err := myContainer.EnvVariable(ctx, "USE_GCC_VERSION")
	if err != nil {
		return err
	}
	*/

	// Figure out target architectures
	/* TODO: removed because of issue #92
	architectures := map[string]string{
		"AARCH64": "-a AARCH64",
		"ARM":     "-a ARM",
		"IA32":    "-a IA32",
		"IA32X64": "-a IA32 -a X64",
		"X64":     "-a X64",
	}
	arch, ok := architectures[opts.Arch]
	if !ok {
		return fmt.Errorf("%w: %s", errUnknownArch, opts.Arch)
	}
	*/

	// Assemble build arguments
	//   and read content of the config file at "defconfig_path"
	// NOTE: removed because of issue #92
	// buildArgs := fmt.Sprintf("%s -p %s -b %s -t GCC%s", arch, opts.Specific.Platform, opts.Specific.ReleaseType, gccVersion)
	var defconfigFileArgs []byte
	if opts.DefconfigPath != "" {
		if _, err := os.Stat(opts.DefconfigPath); !errors.Is(err, os.ErrNotExist) {
			defconfigFileArgs, err = os.ReadFile(opts.DefconfigPath)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Failed to read file '%s' as defconfig_path: file does not exist", opts.DefconfigPath)
		}
	}

	// Assemble commands to build
	buildSteps := [][]string{
		//{"bash", "-c", fmt.Sprintf("source ./edksetup.sh; build %s %s", buildArgs, string(defconfigFileArgs))},
		{"bash", "-c", fmt.Sprintf("%s %s", opts.BuildCommand, string(defconfigFileArgs))},
	}

	// Build
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			return fmt.Errorf("edk2 build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, artifacts)
}
