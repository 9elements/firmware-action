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

// Used to store data from githubaction.Action
//
//	For details see action.yml
type linuxOpts struct {
	gccVersion string
}

// linuxGetOpts is used to fill linuxOpts with data from githubaction.Action
//
//	at the moment, there are no linux-specific options needed
//	It is here to keep the same structure to other recipes
func linuxGetOpts(_ getValFunc, getEnvVar getValFunc) (linuxOpts, error) {
	opts := linuxOpts{
		gccVersion: getEnvVar("GCC_VERSION"),
	}
	return opts, nil
}

// linux builds linux kernel
//
//	docs: https://www.kernel.org/doc/html/latest/kbuild/index.html
func linux(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *linuxOpts, artifacts *[]container.Artifacts) error {
	// No linuxOpts are needed nor used at the moment
	_ = opts

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
		filepath.Join(common.containerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, common.defconfigPath)),
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
		"ARCH": common.arch,
	}

	val, ok := crossCompile[common.arch]
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
	buildSteps := [][]string{
		// remove existing config if exists
		//   -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// x86_64 reuses x86
		{"ln", "--symbolic", "--relative", "arch/x86", "arch/x86_64"},
		// the symlink simplifies this command
		{"cp", defconfigBasename, fmt.Sprintf("arch/%s/configs/%s", common.arch, defconfigBasename)},
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
			return fmt.Errorf("%s build failed: %w", common.target, err)
		}
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, artifacts)
}
