// SPDX-License-Identifier: MIT

// Package recepies / coreboot
package recepies

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// coreboot builds coreboot with all blobs and stuff
func coreboot(ctx context.Context, client *dagger.Client, common *commonOpts, opts *corebootOpts, artifacts *[]container.Artifacts) error {
	// TODO: get blobs in place!

	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      common.sdkVersion,
		MountContainerDir: common.containerWorkDir,
		MountHostDir:      common.repoPath,
		WorkdirContainer:  common.containerWorkDir,
	}
	corebootContainer, err := container.Setup(ctx, client, &containerOpts)
	if err != nil {
		return err
	}

	// Copy over the defconfig file
	defconfigBasename := filepath.Base(common.defconfigPath)
	corebootContainer = corebootContainer.WithFile(
		common.defconfigPath,
		corebootContainer.File(
			filepath.Join(common.containerWorkDir, defconfigBasename),
		),
	)

	// Build
	buildSteps := [][]string{
		// remove existing config if exists
		// -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// generate .config
		{"make", "defconfig", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename)},
		// compile
		{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		{"make", "savedefconfig"},
	}
	for step := range buildSteps {
		corebootContainer, err = corebootContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			return fmt.Errorf("coreboot build failed: %w", err)
		}
	}

	// Extract artifacts
	err = container.GetArtifacts(ctx, corebootContainer, artifacts)
	if err != nil {
		return err
	}

	return nil
}
