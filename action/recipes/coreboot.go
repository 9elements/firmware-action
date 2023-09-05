// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// coreboot builds coreboot with all blobs and stuff
func coreboot(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *corebootOpts, artifacts *[]container.Artifacts) error {
	// TODO: get blobs in place!
	_ = opts
	envVars := map[string]string{}

	return buildWithKernelBuildSystem(ctx, client, common, dockerfileDirectoryPath, envVars, artifacts)
}
