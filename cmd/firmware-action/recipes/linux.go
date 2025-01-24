// SPDX-License-Identifier: MIT

// Package recipes / linux
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/container"
	"github.com/9elements/firmware-action/logging"
)

var errUnknownArchCrossCompile = errors.New("unknown architecture for cross-compilation")

// ANCHOR: LinuxSpecific

// LinuxSpecific is used to store data specific to linux
type LinuxSpecific struct {
	// TODO: either use or remove
	GccVersion string `json:"gcc_version"`
}

// ANCHOR_END: LinuxSpecific

// ANCHOR: LinuxOpts

// LinuxOpts is used to store all data needed to build linux
type LinuxOpts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleEdk2"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Specifies target architecture, such as 'x86' or 'arm64'.
	// Supported options:
	//   - 'i386'
	//   - 'amd64'
	//   - 'arm'
	//   - 'arm64'
	Arch string `json:"arch"`

	// Gives the (relative) path to the defconfig that should be used to build the target.
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`

	// Linux specific options
	LinuxSpecific
}

// ANCHOR_END: LinuxOpts

// GetDepends is used to return list of dependencies
func (opts LinuxOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts LinuxOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// buildFirmware builds linux kernel
//
//	docs: https://www.kernel.org/doc/html/latest/kbuild/index.html
func (opts LinuxOpts) buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) error {
	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      opts.SdkURL,
		MountContainerDir: ContainerWorkDir,
		MountHostDir:      opts.RepoPath,
		WorkdirContainer:  ContainerWorkDir,
		ContainerInputDir: opts.ContainerInputDir,
		InputDirs:         opts.InputDirs,
		InputFiles:        opts.InputFiles,
	}
	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
		slog.Error(
			"Failed to start a container",
			slog.Any("error", err),
		)
		return err
	}

	// Copy over the defconfig file
	defconfigBasename := filepath.Base(opts.DefconfigPath)
	if strings.Contains(defconfigBasename, ".defconfig") {
		// 'make $defconfigBasename' will fail for Linux kernel if the $defconfigBasename
		// contains '.defconfig' string ...
		// it will just fail with generic error (defconfigBasename="linux.defconfig"):
		//   make[1]: *** No rule to make target 'linux.defconfig'.  Stop.
		//   make: *** [Makefile:704: linux.defconfig] Error 2
		// but defconfigBasename="linux_defconfig" works fine
		// Don't know why, just return error and let user deal with it.
		return fmt.Errorf(
			"filename '%s' specified by defconfig_path must not contain '.defconfig'",
			defconfigBasename,
		)
	}
	defconfigRegex := regexp.MustCompile(`.*defconfig$`)
	if !defconfigRegex.MatchString(defconfigBasename) {
		// 'make $defconfigBasename' will fail for Linux kernel if the file
		// does not end with 'defconfig'
		return fmt.Errorf(
			"filename '%s' specified by defconfig_path must end with 'defconfig'",
			defconfigBasename,
		)
	}
	//   not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(ContainerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, opts.DefconfigPath)),
	)

	// Setup environment variables in the container
	//   Handle cross-compilation: Map architecture to cross-compiler
	crossCompile := map[string]string{
		"i386":  "i686-linux-gnu-",
		"amd64": "x86-64-linux-gnu-",
		"arm":   "arm-linux-gnueabi-",
		"arm64": "aarch64-linux-gnu-",
	}
	envVars := map[string]string{
		"ARCH": NormalizeArchitectureForLinux(opts.Arch),
	}

	// Check if cross-compilation is needed
	if NormalizeArchitecture(runtime.GOARCH) != NormalizeArchitecture(opts.Arch) {
		val, ok := crossCompile[opts.Arch]
		if !ok {
			err = errUnknownArchCrossCompile
			slog.Error(
				"Selected unknown cross compilation target architecture",
				slog.String("system_architecture", runtime.GOARCH),
				slog.String("target_architecture", opts.Arch),
				slog.Any("error", err),
			)
			return err
		}
		if val != "" {
			envVars["CROSS_COMPILE"] = val
		}
	}

	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Assemble commands to build
	// TODO: make independent on OS
	buildSteps := [][]string{
		// remove existing config if exists
		//   -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// the symlink simplifies this command
		{"mkdir", "-p", fmt.Sprintf("arch/%s/configs/", NormalizeArchitectureForLinux(opts.Arch))},
		{"mv", defconfigBasename, fmt.Sprintf("arch/%s/configs/%s", NormalizeArchitectureForLinux(opts.Arch), defconfigBasename)},
		// generate dotconfig from defconfig
		{"make", defconfigBasename},
		// compile
		{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		{"make", "savedefconfig"},
	}

	// Execute build commands
	for step := range buildSteps {
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build linux",
				slog.Any("error", err),
			)
			return fmt.Errorf("linux build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, opts.GetArtifacts())
}
