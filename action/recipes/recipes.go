// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// ErrRequiredOptionUndefined is raised when required option is empty or undefined
var (
	ErrRequiredOptionUndefined = errors.New("required option is undefined")
	ErrTargetMissing           = errors.New("no target specified")
	ErrTargetInvalid           = errors.New("unsupported target")
	ErrBuildFailed             = errors.New("build failed")
)

// ContainerWorkDir specifies directory in container used as work directory
var ContainerWorkDir = "/workdir"

// Build recipes, possibly recursively
func Build(ctx context.Context, target string, recursive bool, config Config) error {
	// TODO: This function should be able to recursively build
	//   targets and all their dependencies
	if recursive {
		log.Print("build recursively\n")
	}

	if _, ok := config.Coreboot[target]; ok {
		// Coreboot
		return Execute(ctx, target, "coreboot", config)
	} else if _, ok = config.Linux[target]; ok {
		// Linux
		return Execute(ctx, target, "linux", config)
	} else if _, ok = config.Edk2[target]; ok {
		// Edk2
		return Execute(ctx, target, "edk2", config)
	}

	log.Fatal("Target not found")
	return ErrBuildFailed
}

// Execute recipe
func Execute(ctx context.Context, target string, targetType string, config Config) error {
	log.Printf("building '%s' (%s)", target, targetType)

	// Setup dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Find requested target
	switch targetType {
	case "coreboot":
		log.Printf("building %s", targetType)
		opts := config.Coreboot[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "build", "coreboot.rom"),
				ContainerDir:  false,
				HostPath:      config.Coreboot[target].Common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      config.Coreboot[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return coreboot(ctx, client, &opts, "", &artifacts)
	case "linux":
		log.Printf("building %s", targetType)
		opts := config.Linux[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "vmlinux"),
				ContainerDir:  false,
				HostPath:      config.Linux[target].Common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      config.Linux[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return linux(ctx, client, &opts, "", &artifacts)
	case "edk2":
		log.Printf("building %s", targetType)
		opts := config.Edk2[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "Build"),
				ContainerDir:  true,
				HostPath:      config.Edk2[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return edk2(ctx, client, &opts, "", &artifacts)
	case "":
		return ErrTargetMissing
	default:
		return fmt.Errorf("%w: %s", ErrTargetInvalid, target)
	}
}
