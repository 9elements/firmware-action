// SPDX-License-Identifier: MIT

// Package recipes / linux
package recipes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

var errUnknownArchCrossCompile = errors.New("unknown architecture for cross-compilation")

// LinuxSpecific is used to store data specific to linux
// ANCHOR: LinuxSpecific
type LinuxSpecific struct {
	// TODO: either use or remove
	GccVersion string `json:"gcc_version"`
}

// ANCHOR_END: LinuxSpecific

// LinuxOpts is used to store all data needed to build linux
type LinuxOpts struct {
	// List of IDs this instance depends on
	// Example: [ "MyLittleCoreboot", "MyLittleEdk2"]
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Gives the (relative) path to the defconfig that should be used to build the target.
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`

	// Coreboot specific options
	LinuxSpecific
}

// GetDepends is used to return list of dependencies
func (opts LinuxOpts) GetDepends() []string {
	return opts.Depends
}

// linux builds linux kernel
//
//	docs: https://www.kernel.org/doc/html/latest/kbuild/index.html
func linux(ctx context.Context, client *dagger.Client, opts *LinuxOpts, dockerfileDirectoryPath string, artifacts *[]container.Artifacts) error {
	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      opts.SdkURL,
		MountContainerDir: ContainerWorkDir,
		MountHostDir:      opts.RepoPath,
		WorkdirContainer:  ContainerWorkDir,
	}
	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
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
			"filename '%s' specified by defconfig_path must not contain '.defconfig' in the name",
			defconfigBasename,
		)
	}
	//   not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(ContainerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, opts.DefconfigPath)),
	)

	// Setup environment variables in the container
	//   Handle cross-compilation: Map architecture to cross-compiler
	crossCompile := map[string]string{
		"x86":    "i686-linux-gnu-",
		"x86_64": "",
		"arm":    "arm-linux-gnueabi-",
		"arm64":  "aarch64-linux-gnu-",
	}
	envVars := map[string]string{
		"ARCH": opts.Arch,
	}

	val, ok := crossCompile[opts.Arch]
	if !ok {
		return errUnknownArchCrossCompile
	}
	if val != "" {
		envVars["CROSS_COMPILE"] = val
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
		// x86_64 reuses x86
		{"ln", "--symbolic", "--relative", "arch/x86", "arch/x86_64"},
		// the symlink simplifies this command
		{"cp", defconfigBasename, fmt.Sprintf("arch/%s/configs/%s", opts.Arch, defconfigBasename)},
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
			return fmt.Errorf("linux build failed: %w", err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, artifacts)
}
