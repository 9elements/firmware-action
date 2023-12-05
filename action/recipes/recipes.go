// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// ErrRequiredOptionUndefined is raised when required option is empty or undefined
var (
	ErrRequiredOptionUndefined = errors.New("required option is undefined")
	ErrTargetMissing           = errors.New("no target specified")
	ErrTargetInvalid           = errors.New("unsupported target")
)

// ContainerWorkDir specifies directory in container used as work directory
var ContainerWorkDir = "/workdir"

// Execute recipe
func Execute(ctx context.Context, target string, client *dagger.Client) error {
	common := CommonOpts{
		SdkURL:        "https://ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		Arch:          "dummy",
		RepoPath:      "dummy/dir/",
		DefconfigPath: "dummy",
		OutputDir:     "dummy/dir/",
	}

	switch target {
	case "coreboot":
		opts := CorebootOpts{
			Common:   common,
			Specific: CorebootSpecific{},
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "build", "coreboot.rom"),
				ContainerDir:  false,
				HostPath:      common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.OutputDir,
				HostDir:       true,
			},
		}
		return coreboot(ctx, client, &opts, "", &artifacts)
	case "linux":
		opts := LinuxOpts{
			Common:   common,
			Specific: LinuxSpecific{},
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "vmlinux"),
				ContainerDir:  false,
				HostPath:      common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.OutputDir,
				HostDir:       true,
			},
		}
		return linux(ctx, client, &opts, "", &artifacts)
	case "edk2":
		opts := Edk2Opts{
			Common:   common,
			Specific: Edk2Specific{},
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "Build"),
				ContainerDir:  true,
				HostPath:      common.OutputDir,
				HostDir:       true,
			},
		}
		return edk2(ctx, client, &opts, "", &artifacts)
	case "":
		return ErrTargetMissing
		// return fmt.Errorf("no target specified")
	default:
		return fmt.Errorf("%w: %s", ErrTargetInvalid, target)
	}
}
