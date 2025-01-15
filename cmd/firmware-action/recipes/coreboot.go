// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/container"
	"github.com/9elements/firmware-action/environment"
	"github.com/9elements/firmware-action/filesystem"
	"github.com/9elements/firmware-action/logging"
)

// BlobDef is used to store information about a single blob.
// This structure is not exposed to the user, it is filled in automatically based on user input.
type BlobDef struct {
	// Path to the blob (either file or directory)
	Path string `validate:"required"`

	// Blobs get renamed when moved to this string
	DestinationFilename string `validate:"required"`

	// Kconfig key specifying the filepath to the blob in defconfig
	KconfigKey string `validate:"required"`

	// Is blob a directory? If blob is file, set to FALSE
	IsDirectory bool `validate:"required,boolean"`
}

// ANCHOR: CorebootBlobs

// CorebootBlobs is used to store data specific to coreboot.
type CorebootBlobs struct {
	// ** List of supported blobs **
	// NOTE: The blobs may not be added to the ROM, depends on provided defconfig.
	//
	// Gives the (relative) path to the payload.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(MAINBOARDDIR)/payload`.
	// The Kconfig `CONFIG_PAYLOAD_FILE` will be changed to point to the same path.
	PayloadFilePath string `json:"payload_file_path" type:"blob"`

	// Gives the (relative) path to the Intel Flash descriptor binary.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/descriptor.bin`.
	// The Kconfig `CONFIG_IFD_BIN_PATH` will be changed to point to the same path.
	IntelIfdPath string `json:"intel_ifd_path" type:"blob"`

	// Gives the (relative) path to the Intel Management engine binary.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/me.bin`.
	// The Kconfig `CONFIG_ME_BIN_PATH` will be changed to point to the same path.
	IntelMePath string `json:"intel_me_path" type:"blob"`

	// Gives the (relative) path to the Intel Gigabit Ethernet engine binary.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/gbe.bin`.
	// The Kconfig `CONFIG_GBE_BIN_PATH` will be changed to point to the same path.
	IntelGbePath string `json:"intel_gbe_path" type:"blob"`

	// Gives the (relative) path to the Intel 10 Gigabit Ethernet engine binary.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/10gbe0.bin`.
	// The Kconfig `CONFIG_10GBE_0_BIN_PATH` will be changed to point to the same path.
	Intel10Gbe0Path string `json:"intel_10gbe0_path" type:"blob"`

	// Gives the (relative) path to the Intel FSP binary.
	// In a 'coreboot' build, the file will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/Fsp.fd`.
	// The Kconfig `CONFIG_FSP_FD_PATH` will be changed to point to the same path.
	FspBinaryPath string `json:"fsp_binary_path" type:"blob"`

	// Gives the (relative) path to the Intel FSP header folder.
	// In a 'coreboot' build, the files will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/Include`.
	// The Kconfig `CONFIG_FSP_HEADER_PATH` will be changed to point to the same path.
	FspHeaderPath string `json:"fsp_header_path" type:"blob"`

	// Gives the (relative) path to the Video BIOS Table binary.
	// In a 'coreboot' build, the files will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/vbt.bin`.
	// The Kconfig `CONFIG_INTEL_GMA_VBT_FILE` will be changed to point to the same path.
	VbtPath string `json:"vbt_path" type:"blob"`

	// Gives the (relative) path to the Embedded Controller binary.
	// In a 'coreboot' build, the files will be placed at
	//   `3rdparty/blobs/mainboard/$(CONFIG_MAINBOARD_DIR)/ec.bin`.
	// The Kconfig `CONFIG_EC_BIN_PATH` will be changed to point to the same path.
	EcPath string `json:"ec_path" type:"blob"`
}

// ANCHOR_END: CorebootBlobs

// ANCHOR: CorebootOpts

// CorebootOpts is used to store all data needed to build coreboot.
type CorebootOpts struct {
	// List of IDs this instance depends on
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// Gives the (relative) path to the defconfig that should be used to build the target.
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`

	// Coreboot specific options
	Blobs CorebootBlobs `json:"blobs"`
}

// ANCHOR_END: CorebootOpts

// GetDepends is used to return list of dependencies
func (opts CorebootOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts CorebootOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// GetSources returns slice of paths to all sources which are used for build
func (opts CorebootOpts) GetSources() []string {
	sources := opts.CommonOpts.GetSources()

	// Add blobs to list of sources
	blobs, err := corebootProcessBlobs(opts.Blobs)
	if err != nil {
		slog.Error(
			"Failed to process all blobs",
			slog.Any("error", err),
		)
		return nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return nil
	}
	for blob := range blobs {
		// Path to local file on host
		src := filepath.Join(
			pwd,
			blobs[blob].Path,
		)
		sources = append(sources, src)
	}

	return sources
}

// corebootProcessBlobs is used to figure out blobs from provided data
func corebootProcessBlobs(opts CorebootBlobs) ([]BlobDef, error) {
	blobMap := map[string]BlobDef{
		// Payload
		// docs: https://doc.coreboot.org/payloads.html
		"payload_file_path": {
			DestinationFilename: "payload",
			KconfigKey:          "CONFIG_PAYLOAD_FILE",
			IsDirectory:         false,
		},
		// Intel IFD (Intel Flash Descriptor)
		// docs: https://doc.coreboot.org/util/ifdtool/layout.html
		"intel_ifd_path": {
			DestinationFilename: "descriptor.bin",
			KconfigKey:          "CONFIG_IFD_BIN_PATH",
			IsDirectory:         false,
		},
		// Intel ME (Intel Management Engine)
		"intel_me_path": {
			DestinationFilename: "me.bin",
			KconfigKey:          "CONFIG_ME_BIN_PATH",
			IsDirectory:         false,
		},
		// Intel GbE (Intel Gigabit Ethernet)
		"intel_gbe_path": {
			DestinationFilename: "gbe.bin",
			KconfigKey:          "CONFIG_GBE_BIN_PATH",
			IsDirectory:         false,
		},
		// Intel 10 GbE
		"intel_10gbe0_path": {
			DestinationFilename: "10gbe0.bin",
			KconfigKey:          "CONFIG_10GBE_0_BIN_PATH",
			IsDirectory:         false,
		},
		// Intel FSP binary (Intel Firmware Support Package)
		"fsp_binary_path": {
			DestinationFilename: "Fsp.fd",
			KconfigKey:          "CONFIG_FSP_FD_PATH",
			IsDirectory:         false,
		},
		// Intel FSP header (Intel Firmware Support Package)
		"fsp_header_path": {
			DestinationFilename: "Include",
			KconfigKey:          "CONFIG_FSP_HEADER_PATH",
			IsDirectory:         true,
		},
		// VBT (Video BIOS Table)
		"vbt_path": {
			DestinationFilename: "vbt.bin",
			KconfigKey:          "CONFIG_INTEL_GMA_VBT_FILE",
			IsDirectory:         false,
		},
		// EC (Embedded Controller)
		"ec_path": {
			DestinationFilename: "ec.bin",
			KconfigKey:          "CONFIG_EC_BIN_PATH",
			IsDirectory:         false,
		},
	}
	blobs := []BlobDef{}

	blob := reflect.ValueOf(opts)
	for i := 0; i < blob.Type().NumField(); i++ {
		t := blob.Type().Field(i)

		jsonTag := t.Tag.Get("json")
		jsonType := t.Tag.Get("type")
		if jsonTag != "" && jsonType == "blob" {
			newBlob := blobMap[jsonTag]
			newBlob.Path = blob.Field(i).Interface().(string)
			if newBlob.Path != "" {
				blobs = append(blobs, newBlob)
			}
		}
	}
	return blobs, nil
}

// buildFirmware builds coreboot with all blobs and stuff
func (opts CorebootOpts) buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) (*dagger.Container, error) {
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
		return nil, err
	}

	// Copy over the defconfig file
	defconfigBasename := filepath.Base(opts.DefconfigPath)
	//   not sure why, but without the 'pwd' I am getting different results between CI and 'go test'
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return nil, err
	}
	myContainer = myContainer.WithFile(
		filepath.Join(ContainerWorkDir, defconfigBasename),
		client.Host().File(filepath.Join(pwd, opts.DefconfigPath)),
	)

	// Get value of CONFIG_MAINBOARD_DIR / MAINBOARD_DIR variable from dotconfig
	//   to extract value of 'CONFIG_MAINBOARD_DIR', there must be '.config'
	generateDotConfigCmd := []string{"make", fmt.Sprintf("KBUILD_DEFCONFIG=%s", defconfigBasename), "defconfig"}
	myContainerPrevious := myContainer
	mainboardDir, err := myContainer.
		WithExec(generateDotConfigCmd).
		WithExec([]string{"./util/scripts/config", "-s", "CONFIG_MAINBOARD_DIR"}).
		Stdout(ctx)
	if err != nil {
		slog.Error(
			"Failed to get value of MAINBOARD_DIR from .config",
			slog.Any("error", err),
		)
		return myContainerPrevious, err
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
	blobs, err := corebootProcessBlobs(opts.Blobs)
	if err != nil {
		slog.Error(
			"Failed to process all blobs",
			slog.Any("error", err),
		)
		return nil, err
	}
	for blob := range blobs {
		// Path to local file on host
		src := filepath.Join(
			pwd,
			blobs[blob].Path,
		)
		// Path to file in container
		dst := filepath.Join(
			filepath.Join("3rdparty/blobs/mainboard", mainboardDir),
			blobs[blob].DestinationFilename,
		)

		// Copy into container
		if err = filesystem.CheckFileExists(src); !errors.Is(err, os.ErrExist) {
			slog.Error(
				fmt.Sprintf("Blob '%s' was not found", src),
				slog.String("suggestion", "blobs are copied into container separately from 'input_files' and 'input_dirs', the path should point to files on your host"),
				slog.Any("error", err),
			)
			return nil, err
		}
		if blobs[blob].IsDirectory {
			// Directory
			slog.Info(fmt.Sprintf("Copying directory '%s' to container at '%s'", src, dst))
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
			[]string{"./util/scripts/config", "--set-str", blobs[blob].KconfigKey, dst},
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
	// envVars := map[string]string{}
	envVars, err := corebootPassEnvVars(opts.RepoPath)
	if err != nil {
		slog.Error(
			"Failed to extract environment variables from current environment",
			slog.Any("error", err),
		)
		return myContainerPrevious, fmt.Errorf("coreboot build failed: %w", err)
	}
	for key, value := range envVars {
		myContainer = myContainer.WithEnvVariable(key, value)
	}

	// Build
	for step := range buildSteps {
		myContainerPrevious := myContainer
		myContainer, err = myContainer.
			WithExec(buildSteps[step]).
			Sync(ctx)
		if err != nil {
			slog.Error(
				"Failed to build coreboot",
				slog.Any("error", err),
			)
			return myContainerPrevious, fmt.Errorf("coreboot build failed: %w", err)
		}
	}

	// Extract artifacts
	return myContainer, container.GetArtifacts(ctx, myContainer, opts.CommonOpts.GetArtifacts())
}

func corebootPassEnvVars(repoPath string) (map[string]string, error) {
	passVariables := []string{"KERNELVERSION", "BUILD_TIMELESS"}
	envVariables := environment.FetchEnvVars(passVariables)

	// coreboot build system takes a version from:
	// - environment variable: KERNELVERSION
	// - shell command: git describe ...
	// - content of file: .coreboot-version

	// To check for coreboot version in compiled binary, run these commands:
	//   $ cbfstool build/coreboot.rom extract -n build_info -f /tmp/foo
	//   $ grep COREBOOT_VERSION /tmp/foo

	// coreboot make will fail to run 'git describe' because of
	//   missing '.git' directory once the content of repoPath
	//   is copied into the container
	// To fix this we need to run 'git describe' now and create a new
	//   environment variable to pass over into the container
	// This way, the compiled coreboot binary will not have unknown version

	// If KERNELVERSION is defined in current environment, do nothing
	if _, ok := envVariables["KERNELVERSION"]; ok {
		return envVariables, nil
	}

	// If .coreboot-version file exists in coreboot directory, do nothing
	corebootVersionPath := filepath.Join(repoPath, ".coreboot-version")
	err := filesystem.CheckFileExists(corebootVersionPath)
	if errors.Is(err, os.ErrExist) {
		return envVariables, nil
	}

	// At this point we checked that user did not define their own coreboot version
	// coreboot build system would at this point attempt to run git describe, which would fail
	// Define a new environment variable KERNELVERSION with value from git describe
	//   and then pass it into the container
	err = filesystem.CheckFileExists(repoPath)
	if errors.Is(err, filesystem.ErrPathIsDirectory) {
		describe, err := filesystem.GitDescribeCoreboot(repoPath)
		if err != nil {
			return nil, err
		}
		envVariables["KERNELVERSION"] = describe
	}

	return envVariables, nil
}
