// SPDX-License-Identifier: MIT

package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/sethvargo/go-githubactions"
)

func coreboot(ctx context.Context, action *githubactions.Action, client *dagger.Client) error {
	// retrieve all action related context from environment
	githubCtx, err := action.Context()
	if err != nil {
		return err
	}
	// set custom pipeline name
	client = client.Pipeline("coreboot")
	// get coreboot container
	corebootContainer := client.Container().From("") //TODO: Get Dockerfile URL
	// extract working directory from GitHub context
	workspace := githubCtx.Workspace
	// get reference to source code
	src := client.Host().Directory(workspace)
	// mount source directory into coreboot container
	corebootContainer = corebootContainer.WithMountedDirectory("/src", src).WithWorkdir("/src")
	// compile coreboot rom file
	corebootContainer = corebootContainer.WithExec([]string{"make", "CPUS=$(nproc)"}) //TODO: configure target first
	// get reference to generated coreboot rom file
	romPath := "build/coreboot.rom"
	rom := corebootContainer.File(romPath)
	// export rom file from build container back to host
	_, err = rom.Export(ctx, romPath)

	return err
}
