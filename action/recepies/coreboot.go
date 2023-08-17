// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package recepies

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/kconfig"
	"github.com/plus3it/gorecurcopy"
	"github.com/sethvargo/go-githubactions"
)

const mainboardBlobsPath = "3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/"

type fileMoveIn struct {
	Name       string
	DestName   string
	KconfigKey string
	Directory  bool
}

var filesMoveIn = []fileMoveIn{
	{"payload_file", "payload", "CONFIG_PAYLOAD_FILE", false},
	{"blob_intel_ifd", "descriptor.bin", "CONFIG_IFD_BIN_PATH", false},
	{"blob_intel_gbe", "gbe.bin", "CONFIG_GBE_BIN_PATH", false},
	{"blob_intel_me", "me.bin", "CONFIG_ME_BIN_PATH", false},
	{"fsp_binary_path", "Fsp.fd", "CONFIG_FSP_FD_PATH", false},
	{"fsp_header_path", "Include", "CONFIG_FSP_HEADER_PATH", true},
}

func getPath(action *githubactions.Action, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("no identifier given")
	}
	input := action.GetInput(name)
	if input == "" {
		return "", fmt.Errorf("'%s' is not set", name)
	}
	path, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("Error converting path specified in '%s': %v", name, err)
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Path '%s' specified in '%s' : %v", path, name, err)
	}
	return path, nil
}

func copyFile(src string, dst string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, s.Mode())
	return err
}

func setupCorebootContainer(ctx context.Context, action *githubactions.Action, client *dagger.Client, defconfig string) (*dagger.Container, error) {

	// set custom pipeline name
	client = client.Pipeline("coreboot")
	if action.GetInput("sdk_version") == "" {
		return nil, fmt.Errorf("sdk_version not set")
	}

	// get coreboot container
	url := "ghcr.io/9elements/" + action.GetInput("sdk_version")
	fmt.Printf("Using docker image %s\n", url)
	corebootContainer := client.Container().From(url)

	// mount source directory into coreboot container
	codebase, err := getPath(action, "path")
	if err != nil {
		return nil, err
	}

	// get reference to source code
	srcMount := client.Host().Directory(codebase)

	srcMount = srcMount.WithNewFile("defconfig", defconfig)

	// mount it
	corebootContainer = corebootContainer.WithMountedDirectory(codebase, srcMount).WithWorkdir(codebase)
	return corebootContainer, nil
}

func preCheck(ctx context.Context, action *githubactions.Action) error {
	var err error
	var githubCtx *githubactions.GitHubContext
	args := []string{"defconfig", "path"}
	for i := range args {
		_, err = getPath(action, args[i])
		if err != nil {
			return fmt.Errorf("Error processing input argument '%s': %v", args[i], err)
		}
	}
	for i := range filesMoveIn {
		if action.GetInput(filesMoveIn[i].Name) == "" {
			continue
		}

		_, err = getPath(action, filesMoveIn[i].Name)
		if err != nil {
			return fmt.Errorf("Error processing input argument '%s': %v", filesMoveIn[i].Name, err)
		}
	}
	githubCtx, err = action.Context()
	if err != nil {
		return err
	}
	if githubCtx.Workspace == "" {
		return fmt.Errorf("GITHUB_WORKSPACE not set")
	}

	if action.GetInput("sdk_version") == "" {
		return fmt.Errorf("sdk_version not set")
	}
	return nil
}

func coreboot(ctx context.Context, action *githubactions.Action, client *dagger.Client) error { //nolint:gocyclo

	if err := preCheck(ctx, action); err != nil {
		return err
	}

	// load defconfig
	cfg, err := getPath(action, "defconfig")
	if err != nil {
		return err
	}

	data, err := os.ReadFile(cfg)
	if err != nil {
		return fmt.Errorf("Failed to read defconfig: %v", err)
	}

	defconfigConfig, err := kconfig.NewKconfig(string(data))
	if err != nil {
		return err
	}

	dotConfig, err := generateDotConfigFromDefconfig(ctx, action, client, defconfigConfig.String())
	if err != nil {
		return fmt.Errorf("Failed to extract .config: %v", err)
	}

	// Prepare blobs folder (it also holds precompiled payload and FSP binaries)

	tempDir, err := os.MkdirTemp("", "dagger")
	if err != nil {
		return fmt.Errorf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := range filesMoveIn {
		var blobFile string
		if action.GetInput(filesMoveIn[i].Name) == "" {
			fmt.Printf("Optional argument '%s' not set, skipping\n", filesMoveIn[i].Name)
			continue
		}
		blobFile, err = getPath(action, filesMoveIn[i].Name)
		if err != nil {
			return err
		}
		if !filesMoveIn[i].Directory {
			err = copyFile(blobFile, filepath.Join(tempDir, filesMoveIn[i].DestName))
			if err != nil {
				return err
			}
		} else {
			err = gorecurcopy.CopyDirectory(blobFile, filepath.Join(tempDir, filesMoveIn[i].DestName))
			if err != nil {
				return err
			}
		}
	}

	// Generate proper path
	var mainboardPath = dotConfig.EvalPath(mainboardBlobsPath)
	mainboardPath = strings.ReplaceAll(mainboardPath, "\"", "")

	// Fix defconfig to point to input files
	for i := range filesMoveIn {
		if action.GetInput(filesMoveIn[i].Name) == "" {
			continue
		}
		newPath := filepath.Join(mainboardPath, filesMoveIn[i].DestName)
		fmt.Printf("Updating defconfig '%s' to %s\n", filesMoveIn[i].KconfigKey, newPath)
		(*defconfigConfig)[filesMoveIn[i].KconfigKey] = newPath
	}

	corebootContainer, err := setupCorebootContainer(ctx, action, client, defconfigConfig.String())
	if err != nil {
		return fmt.Errorf("Failed to set up coreboot container: %v", err)
	}

	// get reference to blobs directory
	blobsMount := client.Host().Directory(tempDir)
	// mount it
	codebase, err := getPath(action, "path")
	if err != nil {
		return err
	}
	// Mount blobs
	corebootContainer = corebootContainer.WithMountedDirectory(codebase+"/"+mainboardPath, blobsMount)

	// compile coreboot rom file
	corebootContainer = corebootContainer.WithExec([]string{"make", "distclean"})
	corebootContainer = corebootContainer.WithExec([]string{"make", "defconfig", "KBUILD_DEFCONFIG=defconfig"})
	corebootContainer = corebootContainer.WithExec([]string{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())})

	corebootContainer = corebootContainer.WithExec([]string{"make", "savedefconfig"})
	// retrieve all action related context from environment
	githubCtx, err := action.Context()
	if err != nil {
		return err
	}
	// extract working directory from GitHub context
	workspace := githubCtx.Workspace
	if workspace == "" {
		return fmt.Errorf("GITHUB_WORKSPACE not set")
	}

	outDir := filepath.Join(workspace, "build")
	if err = os.Mkdir(outDir, os.ModePerm); err != nil {
		return err
	}

	// export rom file from build container back to host
	if _, err = corebootContainer.File("build/coreboot.rom").Export(ctx, filepath.Join(outDir, "coreboot.rom")); err != nil {
		return err
	}
	// export defconfig file from build container back to host
	if _, err = corebootContainer.File("defconfig").Export(ctx, filepath.Join(outDir, "defconfig")); err != nil {
		return err
	}
	// export config file from build container back to host
	if _, err = corebootContainer.File(".config").Export(ctx, filepath.Join(outDir, "config")); err != nil {
		return err
	}

	if _, err := corebootContainer.Sync(ctx); err != nil {
		return fmt.Errorf("Error during execution: %v", err)
	}

	return nil
}
