// SPDX-License-Identifier: MIT

// Package recipes / edk2
package recipes

import (
	"context"
	"errors"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

var errUnknownArch = errors.New("unknown architecture")

// Edk2Specific is used to store data specific to coreboot.
// ANCHOR: Edk2Specific
type Edk2Specific struct {
	// Specifies the DSC to use when building EDK2
	// Example: UefiPayloadPkg/UefiPayloadPkg.dsc
	Platform string `json:"platform" validate:"required"`

	// Specifies the build type to use when building EDK2
	// Supported options: DEBUG, RELEASE
	ReleaseType string `json:"release_type" validate:"required"`
}
// ANCHOR_END: Edk2Specific

// Edk2Opts is used to store all data needed to build edk2.
type Edk2Opts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleLinux"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	Common CommonOpts `json:"common" validate:"required"`

	// Coreboot specific options
	Specific Edk2Specific `json:"specific" validate:"required"`
}

// edk2 builds edk2
func edk2(ctx context.Context, client *dagger.Client, opts *Edk2Opts, dockerfileDirectoryPath string, artifacts *[]container.Artifacts) error {
	envVars := map[string]string{
		"WORKSPACE":      ContainerWorkDir,
		"EDK_TOOLS_PATH": "/tools/Edk2/BaseTools",
	}

	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      opts.Common.SdkURL,
		MountContainerDir: ContainerWorkDir,
		MountHostDir:      opts.Common.RepoPath,
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
	gccVersion, err := myContainer.EnvVariable(ctx, "USE_GCC_VERSION")
	if err != nil {
		return err
	}

	// Figure out target architectures
	architectures := map[string]string{
		"AARCH64": "-a AARCH64",
		"ARM":     "-a ARM",
		"IA32":    "-a IA32",
		"IA32X64": "-a IA32 -a X64",
		"X64":     "-a X64",
	}
	arch, ok := architectures[opts.Common.Arch]
	if !ok {
		return fmt.Errorf("%w: %s", errUnknownArch, opts.Common.Arch)
	}

	// Assemble build arguments
	//   and read content of the config file at "defconfig_path"
	buildArgs := fmt.Sprintf("%s -p %s -b %s -t GCC%s", arch, opts.Specific.Platform, opts.Specific.ReleaseType, gccVersion)
	defconfigFileArgs, err := os.ReadFile(opts.Common.DefconfigPath)
	if err != nil {
		return err
	}

	// Assemble commands to build
	buildSteps := [][]string{
		{"bash", "-c", fmt.Sprintf("source ./edksetup.sh; build %s %s", buildArgs, string(defconfigFileArgs))},
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
