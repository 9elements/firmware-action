// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/9elements/firmware-action/action/filesystem"
)

// Used to store information about a single blob
type blobDef struct {
	actionInput         string
	destinationFilename string
	kconfigKey          string
	isDirectory         bool
}

// Used to store data from githubaction.Action
//
//	For details see action.yml
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
			{
				// VBT (Video BIOS Table)
				actionInput:         get("coreboot__vbt_path"),
				destinationFilename: "vbt.bin",
				kconfigKey:          "CONFIG_INTEL_GMA_VBT_FILE",
				isDirectory:         false,
			},
			{
				// EC (Embedded Controller)
				actionInput:         get("coreboot__ec_path"),
				destinationFilename: "ec.bin",
				kconfigKey:          "CONFIG_EC_BIN_PATH",
				isDirectory:         false,
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

// coreboot builds coreboot with all blobs and stuff
func coreboot(ctx context.Context, client *dagger.Client, common *commonOpts, dockerfileDirectoryPath string, opts *corebootOpts, artifacts *[]container.Artifacts) error {
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
	//   not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(common.containerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, common.defconfigPath)),
	)

	// Get value of CONFIG_MAINBOARD_DIR / MAINBOARD_DIR variable from dotconfig
	//   to extract value of 'CONFIG_MAINBOARD_DIR', there must be '.config'
	generateDotConfigCmd := []string{"make", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename), "defconfig"}
	mainboardDir, err := myContainer.
		WithExec(generateDotConfigCmd).
		WithExec([]string{"./util/scripts/config", "-s", "CONFIG_MAINBOARD_DIR"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	//   strip newline from mainboardDir
	mainboardDir = strings.Replace(mainboardDir, "\n", "", -1)

	// Assemble commands to build
	buildSteps := [][]string{
		// remove existing config if exists
		// -f: ignore nonexistent files
		{"rm", "-f", ".config"},
		// generate dotconfig from defconfig
		generateDotConfigCmd,
	}

	// Handle blobs
	// Firstly copy all the blobs into building container.
	// Then use './util/scripts/config' script in coreboot repository to update configuration
	//   options for said blobs (this must run inside container).
	for blob := range opts.blobs {
		// Path to local file on host
		src := filepath.Join(
			pwd,
			opts.blobs[blob].actionInput,
		)
		// Path to file in container
		dst := filepath.Join(
			filepath.Join("3rdparty/blobs/mainboard", mainboardDir),
			opts.blobs[blob].destinationFilename,
		)

		// Copy into container
		if err = filesystem.CheckFileExists(src); !errors.Is(err, os.ErrExist) {
			return err
		}
		if opts.blobs[blob].isDirectory {
			// Directory
			log.Printf("Copying directory '%s' to container at '%s'", src, dst)
			myContainer = myContainer.WithExec([]string{"mkdir", "-p", dst})
			// myContainer = myContainer.WithMountedDirectory(
			// can't use WithMountedDirectory because the repo (aka working directory)
			//   is already mounted with WithMountedDirectory
			//   this nesting causes problems
			myContainer = myContainer.WithDirectory(
				dst,
				client.Host().Directory(src),
			)
		} else {
			// File
			myContainer = myContainer.WithFile(
				dst,
				client.Host().File(src),
			)
		}

		// Fix defconfig
		buildSteps = append(
			buildSteps,
			// update coreboot config value related to blob to actual path of the blob
			[]string{"./util/scripts/config", "--set-str", opts.blobs[blob].kconfigKey, dst},
		)
	}

	buildSteps = append(
		buildSteps,
		// compile
		[]string{"make", "-j", fmt.Sprintf("%d", runtime.NumCPU())},
		// for documenting purposes
		[]string{"make", "savedefconfig"},
	)

	// Setup environment variables in the container
	envVars := map[string]string{}
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
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
