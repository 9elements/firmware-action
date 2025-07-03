// SPDX-License-Identifier: MIT

// Package recipes / uboot
package recipes

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/container"
	"github.com/9elements/firmware-action/cmd/firmware-action/logging"
)

// UBootSpecific is used to store data specific to u-root

// ANCHOR: UBootOpts

// UBootOpts is used to store all data needed to build u-root
type UBootOpts struct {
	// Common options like paths etc.
	CommonOpts

	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleEdk2"]
	Depends []string `json:"depends"`

	// Specifies target architecture, such as 'x86' or 'arm64'
	Arch string `json:"arch"`

	// Gives the (relative) path to the defconfig that should be used to build the target.
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`
}

// ANCHOR_END: UBootOpts

// GetDepends is used to return list of dependencies
func (opts UBootOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts UBootOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// buildFirmware builds u-root
func (opts UBootOpts) buildFirmware(ctx context.Context, client *dagger.Client) error {
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

	// U-Boot is closely related to Linux, so I assume similar requirements / problems
	// Copy over the defconfig file
	defconfigBasename := filepath.Base(opts.DefconfigPath)

	err = ValidateLinuxDefconfigFilename(opts.DefconfigPath)
	if err != nil {
		return err
	}
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

	// Setup environment variables in the container
	//   Handle cross-compilation: Map architecture to cross-compiler
	envVars, err := LinuxCrossCompilationArchMap(opts.Arch)
	if err != nil {
		return err
	}

	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Assemble commands to build
	// TODO: make independent on OS
	buildSteps := [][]string{
		// remove existing config if exists
		//   -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// generate dotconfig from defconfig
		{"mv", defconfigBasename, filepath.Join("configs", defconfigBasename)},
		{"make", defconfigBasename},
		// compile
		{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		{"make", "savedefconfig"},
	}

	// Execute build commands
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build u-boot",
				slog.Any("error", err),
			)

			return fmt.Errorf("u-boot build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, opts.GetArtifacts())
}
