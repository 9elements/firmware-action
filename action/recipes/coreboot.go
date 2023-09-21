// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/9elements/firmware-action/action/filesystem"
)

type ctxKey string

const (
	blobsLocation                   = "3rdparty/blobs/mainboard/${MAINBOARD_DIR}"
	blobsLocationCtxKey      ctxKey = "BLOBS_LOCATION"
	additionalCommandsCtxKey ctxKey = "ADDITIONAL_COMMANDS_TO_RUN"
)

// coreboot builds coreboot with all blobs and stuff
func coreboot(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *corebootOpts, artifacts *[]container.Artifacts) error {
	envVars := map[string]string{}

	additionalCommandsToRun := [][]string{} // Used to modify defconfig inside container

	// Handle blobs
	//   To keep using 'buildWithKernelBuildSystem' without breaking Linux kernel build
	//   there has to be a bit hackish solution ... :(
	//   Biggest problem is that most preparation should be done here, but the container
	//   is not available yet. So this solution exploits the context to pass over few things.
	// Firstly copy all the blobs into temporary location, which will then be mounted into
	//   building container.
	// Then use './util/scripts/config' script in coreboot repository to update configuration
	//   options for said blobs (this must run inside container).

	//  Collect all blobs into temporary directory
	tmpDir, err := os.MkdirTemp("", "blobs")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	for blob := range opts.blobs {
		src := opts.blobs[blob].actionInput
		dst := filepath.Join(
			blobsLocation,
			opts.blobs[blob].destinationFilename,
		)
		tmpDst := filepath.Join(
			tmpDir,
			opts.blobs[blob].destinationFilename,
		)

		// Copy into proper places
		if opts.blobs[blob].isDirectory {
			// Directory
			if err := filesystem.CopyDir(src, tmpDst); err != nil {
				return err
			}
		} else {
			// File
			if err := filesystem.CopyFile(src, tmpDst); err != nil {
				return err
			}
		}

		// Fix defconfig
		additionalCommandsToRun = append(
			additionalCommandsToRun,
			// The '"sh", "-c"' is needed to get access to environment variables
			//   meaning replace '${MAINBOARD_DIR}' in 'blobsLocation' with actual path
			[]string{"sh", "-c", fmt.Sprintf("./util/scripts/config --set-str %s \"%s\"", opts.blobs[blob].kconfigKey, dst)},
		)
	}

	// Add additionalCommandsToRun into context
	ctx = context.WithValue(ctx, additionalCommandsCtxKey, additionalCommandsToRun)
	ctx = context.WithValue(ctx, blobsLocationCtxKey, tmpDir)

	// Build
	return buildWithKernelBuildSystem(ctx, client, common, dockerfileDirectoryPath, envVars, artifacts)
}
