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
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/sethvargo/go-githubactions"
)

//===================
// Universal options
//===================

var (
	errRequiredOptionUndefined = errors.New("required option is undefined")
	errArchUndefined           = errors.New("environment variable 'ARCH' is not defined")
)

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
func commonGetOpts(get getValFunc) (commonOpts, error) {
	opts := commonOpts{
		target:           get("target"),
		sdkVersion:       get("sdk_version"),
		arch:             get("architecture"),
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

// Used to store data from githubaction.Action
// For details see action.yml
type linuxOpts struct{}

// linuxGetOpts is used to fill linuxOpts with data from githubaction.Action
func linuxGetOpts(_ getValFunc) (linuxOpts, error) {
	opts := linuxOpts{}
	return opts, nil
}

//======
// EDK2
//======

//=====================
// Universal Functions
//=====================

// Execute recipe
func Execute(ctx context.Context, client *dagger.Client, action *githubactions.Action) error {
	common, err := commonGetOpts(action.GetInput)
	if err != nil {
		return err
	}

	switch action.GetInput("target") {
	case "coreboot":
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
		return coreboot(ctx, client, &common, "", &opts, &artifacts)
	case "linux":
		opts, err := linuxGetOpts(action.GetInput)
		if err != nil {
			return err
		}
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "vmlinux"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
			},
			{
				ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      common.outputDir,
			},
		}
		return linux(ctx, client, &common, "", &opts, &artifacts)
	/*
		case "edk2":
			return edk2(ctx, action, client)
	*/
	case "":
		return fmt.Errorf("no target specified")
	default:
		return fmt.Errorf("unsupported target: %s", action.GetInput("target"))
	}
}

// buildWithKernelBuildSystem is a generic function to build stuff with Kernel Build System
// usable for linux kernel and coreboot
// https://www.kernel.org/doc/html/latest/kbuild/index.html
func buildWithKernelBuildSystem(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, envVars map[string]string, artifacts *[]container.Artifacts) error {
	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      common.sdkVersion,
		MountContainerDir: common.containerWorkDir,
		MountHostDir:      common.repoPath,
		WorkdirContainer:  common.containerWorkDir,
	}
	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
		return err
	}

	// Copy over the defconfig file
	defconfigBasename := filepath.Base(common.defconfigPath)
	if strings.Contains(defconfigBasename, ".defconfig") {
		// 'make $defconfigBasename' will fail for Linux kernel if the $defconfigBasename
		// contains '.defconfig' string ...
		// it will just fail with generic error (defconfigBasename="linux.defconfig"):
		//   make[1]: *** No rule to make target 'linux.defconfig'.  Stop.
		//   make: *** [Makefile:704: linux.defconfig] Error 2
		// but defconfigBasename="linux_defconfig" works fine
		return fmt.Errorf("filename '%s' specified by defconfig_path must not contain '.defconfig' in the name", defconfigBasename)
	}
	myContainer = myContainer.WithFile(
		common.defconfigPath,
		myContainer.File(
			filepath.Join(common.containerWorkDir, defconfigBasename),
		))

	// Setup environment variables in the container
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Assemble commands to build
	buildSteps := [][]string{
		// remove existing config if exists
		// -f: ignore nonexistent files
		{"rm", "-f", ".config"},
	}
	switch common.target {
	case "coreboot":
		// this thing works because of some black magic script in coreboot build chain
		// does not work for linux kernel
		buildSteps = append(
			buildSteps,
			[]string{"make", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename), "defconfig"},
		)
	case "linux":
		// results should be: make ARCH=x86 custom_defconfig
		arch, ok := envVars["ARCH"]
		if !ok {
			return errArchUndefined
		}
		buildSteps = append(
			buildSteps,
			[]string{"ln", "--symbolic", "--relative", "arch/x86", "arch/x86_64"},
			[]string{"cp", defconfigBasename, fmt.Sprintf("arch/%s/configs/%s", arch, defconfigBasename)},
			[]string{"make", defconfigBasename},
		)
	}
	buildSteps = append(
		buildSteps,
		[]string{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		[]string{"make", "savedefconfig"},
	)

	// Build
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			return fmt.Errorf("%s build failed: %w", common.target, err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, artifacts)
}
