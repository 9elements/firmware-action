// SPDX-License-Identifier: MIT

// Package recepies yay!
package recepies

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/sethvargo/go-githubactions"
)

//===================
// Universal options
//===================

var errRequiredOptionUndefined = errors.New("required option is undefined")

type (
	getValFunc func(string) string
)

// commontOpts is common to all targets
// Used to store data from githubaction.Action
// For details see action.yml
type commonOpts struct {
	target           string
	sdkVersion       string
	repoPath         string
	defconfigPath    string
	containerWorkDir string
	outputDir        string
}

// commonGetOpts is used to fill commontOpts with data from githubaction.Action
func commonGetOpts(get getValFunc) (commonOpts, error) {
	opts := commonOpts{
		target:           get("target"),
		sdkVersion:       get("sdk_version"),
		repoPath:         get("repo_path"),
		defconfigPath:    get("defconfig_path"),
		containerWorkDir: get("GITHUB_WORKSPACE"),
		outputDir:        get("output"),
	}

	// Check if required options are not empty
	missing := []string{}
	requiredOptions := map[string]string{
		"target":           opts.target,
		"sdk_version":      opts.sdkVersion,
		"repo_path":        opts.repoPath,
		"defconfig_path":   opts.defconfigPath,
		"containerWorkDir": opts.containerWorkDir,
		"output":           opts.outputDir,
	}
	for key, val := range requiredOptions {
		if val == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return opts, fmt.Errorf("%w: %s", errRequiredOptionUndefined, strings.Join(missing, ", "))
	}

	// Check if sdk_version is URL, if not make it into URL defaulting to our containers
	_, err := url.ParseRequestURI(opts.sdkVersion)
	if err != nil {
		opts.sdkVersion = path.Join("ghcr.io/9elements/firmware-action", opts.sdkVersion)
	}
	return opts, nil
}

//==========
// COREBOOT
//==========

// Used to store data from githubaction.Action
// For details see action.yml
type corebootOpts struct {
	payloadFilePath  string
	blobIntelIfdPath string
	blobIntelMePath  string
	blobIntelGbePath string
	fspBinaryPath    string
	fspHeaderPath    string
}

// commonGetOpts is used to fill corebootOpts with data from githubaction.Action
func corebootGetOpts(get getValFunc) (corebootOpts, error) {
	opts := corebootOpts{
		payloadFilePath:  get("coreboot__payload_file_path"),
		blobIntelIfdPath: get("coreboot__blob_intel_ifd_path"),
		blobIntelMePath:  get("coreboot__blob_intel_me_path"),
		blobIntelGbePath: get("coreboot__blob_intel_gbe_path"),
		fspBinaryPath:    get("coreboot__fsp_binary_path"),
		fspHeaderPath:    get("coreboot__fsp_header_path"),
	}

	// Check if required options are not empty
	// ... I don't think any of these are always required, might depend on provided defconfig
	return opts, nil
}

//=======
// LINUX
//=======

//======
// EDK2
//======

//=====================
// Universal Functions
//=====================

// Execute recepie
func Execute(ctx context.Context, client *dagger.Client, action *githubactions.Action) error {
	switch action.GetInput("target") {
	case "coreboot":
		common, err := commonGetOpts(action.GetInput)
		if err != nil {
			return err
		}
		opts, err := corebootGetOpts(action.GetInput)
		if err != nil {
			return err
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "build", "coreboot.rom"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
			},
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
			},
		}
		return coreboot(ctx, client, &common, &opts, &artifacts)
	/*
		case "linux":
			return linux(ctx, action, client)
		case "edk2":
			return edk2(ctx, action, client)
	*/
	case "":
		return fmt.Errorf("no target specified")
	default:
		return fmt.Errorf("unsupported target: %s", action.GetInput("target"))
	}
}
