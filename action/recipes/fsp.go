// SPDX-License-Identifier: MIT

// Package recipes / fsp
package recipes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

// fsp builds fsp
func fsp(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *fspOpts, artifacts *[]container.Artifacts) error {
	// Prepare environment variables
	//   Read content of the config file at "defconfig_path"
	//   and use it to set FSP_BUILD_OPTION_PCD
	defconfigFileArgs, err := os.ReadFile(common.defconfigPath)
	if err != nil {
		return err
	}
	envVars := map[string]string{
		"WORKSPACE_SILICON":    filepath.Join(common.containerWorkDir, "Intel"),
		"WORKSPACE_COMMON":     filepath.Join(common.containerWorkDir, "Intel"),
		"FSP_BUILD_OPTION_PCD": string(defconfigFileArgs),
	}

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

	// Setup environment variables in the container
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Assemble commands to build
	buildSteps := []string{
		"cd \"${WORKSPACE_CORE}\"",
		"source ./edksetup.sh",
		fmt.Sprintf("export WORKSPACE=%s", common.containerWorkDir),
		// The value must change here and now, otherwise the build will fail
		// Message to the person who created this build system: insert "you had one job" meme
		"cd \"${WORKSPACE}\"",
		fmt.Sprintf("%s -%s", opts.buildCmd, opts.releaseType),
	}
	cmd := []string{"bash", "-c", strings.Join(buildSteps, "; ")}

	// Build
	myContainer, err = myContainer.
		WithExec(cmd).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("%s build failed: %w", common.target, err)
	}

	// Extract artifacts
	return container.GetArtifacts(ctx, myContainer, artifacts)
}
