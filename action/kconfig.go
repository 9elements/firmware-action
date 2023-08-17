// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/kconfig"
	"github.com/sethvargo/go-githubactions"
)

func generateDotConfigFromDefconfig(ctx context.Context, action *githubactions.Action, client *dagger.Client, defconfig string) (*kconfig.Kconfig, error) {
	corebootContainer, err := setupCorebootContainer(ctx, action, client, defconfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .config: %v", err)
	}

	// generate .config
	corebootContainer = corebootContainer.WithExec([]string{"rm", ".config"})
	corebootContainer, err = corebootContainer.Sync(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error during execution: %v", err)
	}

	o, err := corebootContainer.Stdout(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	o, err = corebootContainer.Stderr(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	corebootContainer = corebootContainer.WithExec([]string{"make", "defconfig", "KBUILD_DEFCONFIG=defconfig"})
	corebootContainer, err = corebootContainer.Sync(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error during execution: %v", err)
	}

	o, err = corebootContainer.Stdout(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	o, err = corebootContainer.Stderr(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	// Extract .config
	rom := corebootContainer.File(".config")
	dotconfigRaw, err := rom.Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .config: %v", err)
	}

	var dotConfig *kconfig.Kconfig
	dotConfig, err = kconfig.NewKconfig(dotconfigRaw)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert .config: %v", err)
	}

	return dotConfig, nil
}
