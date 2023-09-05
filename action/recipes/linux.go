// SPDX-License-Identifier: MIT

// Package recipes / linux
package recipes

import (
	"context"
	"errors"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

var errUnknownArchCrossCompile = errors.New("unknown architecture for cross-compilation")

// linux builds linux kernel
func linux(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *linuxOpts, artifacts *[]container.Artifacts) error {
	// Not sure if there will be any linuxOpts
	_ = opts

	crossCompile := map[string]string{
		"x86":    "i686-linux-gnu-",
		"x86_64": "",
		"arm":    "arm-linux-gnueabi-",
		"arm64":  "aarch64-linux-gnu-",
	}

	envVars := map[string]string{
		"ARCH": common.arch,
	}

	val, ok := crossCompile[common.arch]
	if !ok {
		return errUnknownArchCrossCompile
	}
	if val != "" {
		envVars["CROSS_COMPILE"] = val
	}

	return buildWithKernelBuildSystem(ctx, client, common, dockerfileDirectoryPath, envVars, artifacts)
}
