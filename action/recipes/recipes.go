// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
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

//==========
// COREBOOT
//==========

// Used to store information about a single blob
type blobDef struct {
	actionInput         string
	destinationFilename string
	kconfigKey          string
	isDirectory         bool
}

// Used to store data from githubaction.Action
// For details see action.yml
type corebootOpts struct {
	blobs []blobDef
}

// commonGetOpts is used to fill corebootOpts with data from githubaction.Action
func corebootGetOpts(get getValFunc) (corebootOpts, error) {
	// 'allOpts' most importantly contains definitions of all possible (supported) blobs
	allOpts := corebootOpts{
		blobs: []blobDef{
			{
				// Payload
				// docs: https://doc.coreboot.org/payloads.html
				actionInput:         get("coreboot__payload_file_path"),
				destinationFilename: "payload",
				kconfigKey:          "CONFIG_PAYLOAD_FILE",
				isDirectory:         false,
			},
			{
				// Intel IFD (Intel Flash Descriptor)
				// docs: https://doc.coreboot.org/util/ifdtool/layout.html
				actionInput:         get("coreboot__intel_ifd_path"),
				destinationFilename: "descriptor.bin",
				kconfigKey:          "CONFIG_IFD_BIN_PATH",
				isDirectory:         false,
			},
			{
				// Intel ME (Intel Management Engine)
				actionInput:         get("coreboot__intel_me_path"),
				destinationFilename: "me.bin",
				kconfigKey:          "CONFIG_ME_BIN_PATH",
				isDirectory:         false,
			},
			{
				// Intel GbE (Intel Gigabit Ethernet)
				actionInput:         get("coreboot__intel_gbe_path"),
				destinationFilename: "gbe.bin",
				kconfigKey:          "CONFIG_GBE_BIN_PATH",
				isDirectory:         false,
			},
			{
				// Intel FSP binary (Intel Firmware Support Package)
				actionInput:         get("coreboot__fsp_binary_path"),
				destinationFilename: "Fsp.fd",
				kconfigKey:          "CONFIG_FSP_FD_PATH",
				isDirectory:         false,
			},
			{
				// Intel FSP header (Intel Firmware Support Package)
				actionInput:         get("coreboot__fsp_header_path"),
				destinationFilename: "Include",
				kconfigKey:          "CONFIG_FSP_HEADER_PATH",
				isDirectory:         true,
			},
		},
	}

	// If any of blobs defined in 'allOpts' is passed into the action as input, append it to 'opts'
	opts := corebootOpts{}
	for blob := range allOpts.blobs {
		if allOpts.blobs[blob].actionInput != "" {
			opts.blobs = append(opts.blobs, allOpts.blobs[blob])
		}
	}

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

// Used to store data from githubaction.Action
// For details see action.yml
type edk2Opts struct {
	platform    string
	releaseType string
}

// edk2GetOpts is used to fill edk2Opts with data from githubaction.Action
func edk2GetOpts(get getValFunc) (edk2Opts, error) {
	opts := edk2Opts{
		platform:    get("edk2__platform"),
		releaseType: get("edk2__release_type"),
	}

	// Check if required options are not empty
	missing := []string{}
	requiredOptions := map[string]string{
		"edk2__platform":     opts.platform,
		"edk2__release_type": opts.releaseType,
	}
	for key, val := range requiredOptions {
		if val == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return opts, fmt.Errorf("%w: %s", errRequiredOptionUndefined, strings.Join(missing, ", "))
	}

	return opts, nil
}

//=====================
// Universal Functions
//=====================

// Execute recipe
func Execute(ctx context.Context, client *dagger.Client, action *githubactions.Action) error {
	common, err := commonGetOpts(action.GetInput, action.Getenv)
	if err != nil {
		return err
	}

	switch common.target {
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
		opts, err := linuxGetOpts(action.GetInput)
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
		opts, err := edk2GetOpts(action.GetInput)
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
	// Not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(common.containerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, common.defconfigPath)),
	)

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
	generateDotConfigCmd := []string{"make", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename), "defconfig"}
	switch common.target {
	case "coreboot":
		// this thing works because of some black magic script in coreboot build chain
		// does not work for linux kernel
		buildSteps = append(
			buildSteps,
			generateDotConfigCmd,
		)
		// Everything that follows is just to get blobs into proper place
		//   Add additional commands to fix the defconfig
		buildSteps = append(
			buildSteps,
			ctx.Value(additionalCommandsCtxKey).([][]string)[:]...,
		)
		//   Add CONFIG_MAINBOARD_DIR from defconfig as environment variable MAINBOARD_DIR
		mainboardDir, err := myContainer.
			WithExec(generateDotConfigCmd).
			WithExec([]string{"./util/scripts/config", "-s", "CONFIG_MAINBOARD_DIR"}).
			Stdout(ctx)
			// To extract value of 'CONFIG_MAINBOARD_DIR', there must be '.config'
		if err != nil {
			return err
		}
		//   Strip newline from mainboardDir
		mainboardDir = strings.Replace(mainboardDir, "\n", "", -1)
		myContainer = myContainer.WithEnvVariable("MAINBOARD_DIR", mainboardDir)
		//   Add directory with blobs to container
		myContainer = myContainer.
			WithMountedDirectory(
				filepath.Join(common.containerWorkDir, common.repoPath, blobsLocation),
				client.Host().Directory(ctx.Value(blobsLocationCtxKey).(string)),
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
