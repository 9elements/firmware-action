// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

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

var errRequiredOptionUndefined = errors.New("required option is undefined")

type getValFunc func(string) string

// commonOpts is common to all targets
// Used to store data from githubaction.Action
// For details see action.yml
type commonOpts struct {
	target           string
	sdkVersion       string
	arch             string
	repoPath         string
	defconfigPath    string
	containerWorkDir string
	outputDir        string
}

// commonGetOpts is used to fill commonOpts with data from githubaction.Action
func commonGetOpts(getInputVar getValFunc, getEnvVar getValFunc) (commonOpts, error) {
	opts := commonOpts{
		target:           getInputVar("target"),
		sdkVersion:       getInputVar("sdk_version"),
		arch:             getInputVar("architecture"),
		repoPath:         getInputVar("repo_path"),
		defconfigPath:    getInputVar("defconfig_path"),
		containerWorkDir: getEnvVar("GITHUB_WORKSPACE"),
		outputDir:        getInputVar("output"),
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

	// Check if sdk_version is URL to a container in some container registry
	//   (for example "docker.io/library/ubuntu:latest")
	// If sdk_version is not a URL, assume it is a name of container and make it into URL
	//   pointing to our container registry at "ghcr.io/9elements/firmware-action"
	// WARNING:
	//   For url.ParseRequestURI string "edk2-stable202105:main" is a valid URL (RFC 3986)
	//     so checking err alone is not enough.
	//   Valid URL should contain Fully Qualified Domain Name (FQDN) and so checking for empty
	//     parsedUrl.Hostname seems to do the trick.
	if parsedURL, err := url.ParseRequestURI(opts.sdkVersion); err != nil || parsedURL.Hostname() == "" {
		// opts.sdkVersion is not URL
		opts.sdkVersion = path.Join("ghcr.io/9elements/firmware-action", opts.sdkVersion)
	}
	return opts, nil
}

// Execute recipe
func Execute(ctx context.Context, client *dagger.Client, action *githubactions.Action) error {
	common, err := commonGetOpts(action.GetInput, action.Getenv)
	if err != nil {
		return err
	}

	switch common.target {
	case "coreboot":
		opts, err := corebootGetOpts(action.GetInput, action.Getenv)
		if err != nil {
			return err
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "build", "coreboot.rom"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
				HostDir:       true,
			},
		}
		return coreboot(ctx, client, &common, "", &opts, &artifacts)
	case "linux":
		opts, err := linuxGetOpts(action.GetInput, action.Getenv)
		if err != nil {
			return err
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "vmlinux"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
				HostDir:       true,
			},
		}
		return linux(ctx, client, &common, "", &opts, &artifacts)
	case "edk2":
		opts, err := edk2GetOpts(action.GetInput, action.Getenv)
		if err != nil {
			return err
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "Build"),
				ContainerDir:  true,
				HostPath:      common.outputDir,
				HostDir:       true,
			},
		}
		return edk2(ctx, client, &common, "", &opts, &artifacts)
	case "":
		return fmt.Errorf("no target specified")
	default:
		return fmt.Errorf("unsupported target: %s", common.target)
	}
}
