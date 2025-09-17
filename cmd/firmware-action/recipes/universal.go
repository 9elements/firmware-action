// SPDX-License-Identifier: MIT

// Package recipes / universal
package recipes

import (
	"context"
	"fmt"
	"log/slog"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/container"
)

// UniversalSpecific is used to store data specific to the universal command module
// ANCHOR: UniversalSpecific
type UniversalSpecific struct {
	// Specifies build commands to execute inside container
	BuildCommands []string `json:"build_commands" validate:"required"`
}

// ANCHOR_END: UniversalSpecific

// ANCHOR: UniversalOpts

// UniversalOpts is used to store all data needed to run universal commands
type UniversalOpts struct {
	// Common options like paths etc.
	CommonOpts

	// Universal specific options
	UniversalSpecific

	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleEdk2"]
	Depends []string `json:"depends"`
}

// ANCHOR_END: UniversalOpts

// GetDepends is used to return list of dependencies
func (opts UniversalOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts UniversalOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// buildFirmware builds (or rather executes) universal command module
func (opts UniversalOpts) buildFirmware(ctx context.Context, client *dagger.Client) error {
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
	buildSteps := [][]string{}
	for _, cmd := range opts.BuildCommands {
		buildSteps = append(
			buildSteps,
			[]string{"bash", "-c", cmd},
		)
	}

	// Execute build commands
	for step := range buildSteps {
		myContainer = myContainer.WithExec(buildSteps[step])
	}

	myContainer, err = myContainer.Sync(ctx)
	if err != nil {
		slog.Error(
			"Failed to build universal",
			slog.Any("error", err),
		)

		return fmt.Errorf("universal build failed: %w", err)
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, opts.GetArtifacts())
}
