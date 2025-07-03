// SPDX-License-Identifier: MIT

// Package recipes / uroot
package recipes

import (
	"context"
	"fmt"
	"log/slog"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/container"
)

// URootSpecific is used to store data specific to u-root
// ANCHOR: URootSpecific
type URootSpecific struct {
	// Specifies build command to use
	BuildCommand string `json:"build_command" validate:"required"`
}

// ANCHOR_END: URootSpecific

// ANCHOR: URootOpts

// URootOpts is used to store all data needed to build u-root
type URootOpts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleEdk2"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// u-root specific options
	URootSpecific
}

// ANCHOR_END: URootOpts

// GetDepends is used to return list of dependencies
func (opts URootOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts URootOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// buildFirmware builds u-root
func (opts URootOpts) buildFirmware(ctx context.Context, client *dagger.Client) error {
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

	// Assemble commands to build
	buildSteps := [][]string{
		// run user-defined build command
		{"bash", "-c", opts.BuildCommand},
	}

	// Execute build commands
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build u-root",
				slog.Any("error", err),
			)

			return fmt.Errorf("u-root build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, opts.GetArtifacts())
}
