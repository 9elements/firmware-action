// SPDX-License-Identifier: MIT

// Package recipes / edk2
package recipes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
)

var errUnknownArch = errors.New("unknown architecture")

// Used to store data from githubaction.Action
// For details see action.yml
type edk2Opts struct {
	platform    string
	releaseType string
	gccVersion  string
}

// edk2GetOpts is used to fill edk2Opts with data from githubaction.Action
func edk2GetOpts(getInputVar getValFunc, getEnvVar getValFunc) (edk2Opts, error) {
	opts := edk2Opts{
		platform:    getInputVar("edk2__platform"),
		releaseType: getInputVar("edk2__release_type"),
		gccVersion:  getEnvVar("USE_GCC_VERSION"),
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

// edk2 builds edk2
func edk2(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *edk2Opts, artifacts *[]container.Artifacts) error {
	envVars := map[string]string{
		"WORKSPACE":      common.containerWorkDir,
		"EDK_TOOLS_PATH": "/tools/Edk2/BaseTools",
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

	// Figure out target architectures
	architectures := map[string]string{
		"AARCH64": "-a AARCH64",
		"ARM":     "-a ARM",
		"IA32":    "-a IA32",
		"IA32X64": "-a IA32 -a X64",
		"X64":     "-a X64",
	}
	arch, ok := architectures[common.arch]
	if !ok {
		return fmt.Errorf("%w: %s", errUnknownArch, common.arch)
	}

	// Assemble build arguments
	//   and read content of the config file at "defconfig_path"
	buildArgs := fmt.Sprintf("%s -p %s -b %s -t GCC%s", arch, opts.platform, opts.releaseType, opts.gccVersion)
	defconfigFileArgs, err := os.ReadFile(common.defconfigPath)
	if err != nil {
		return err
	}

	// Assemble commands to build
	buildSteps := [][]string{
		{"bash", "-c", fmt.Sprintf("source ./edksetup.sh; build %s %s", buildArgs, string(defconfigFileArgs))},
	}

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
