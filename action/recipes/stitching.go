// SPDX-License-Identifier: MIT

// Package recipes / stitching
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/9elements/firmware-action/action/logging"
	"github.com/dustin/go-humanize"
)

var (
	errFailedToDetectRomSize = errors.New("failed to detect ROM size from IFD")
	errBaseFileBiggerThanIfd = errors.New("base_file is bigger than size defined in IFD")
)

const ifdtoolPath = "ifdtool"

// ANCHOR: IfdtoolEntry

// IfdtoolEntry is for injecting a file at `path` into region `TargetRegion`
type IfdtoolEntry struct {
	// Gives the (relative) path to the binary blob
	Path string `json:"path" validate:"required,filepath"`

	// Region where to inject the file
	// For supported options see `ifdtool --help`
	TargetRegion string `json:"target_region" validate:"required"`

	// Additional (optional) arguments and flags
	// For example:
	//   `--platform adl`
	// For supported options see `ifdtool --help`
	OptionalArguments []string `json:"optional_arguments"`

	// Ignore entry if the file is missing
	IgnoreIfMissing bool `json:"ignore_if_missing" type:"boolean"`

	// For internal use only - whether or not the blob should be injected
	// Firstly it is checked if the blob file exists, if not a if `IgnoreIfMissing` is set to `true`,
	//   then `Skip` is set to `true` to remove need for additional repetitive checks later in program
	Skip bool
}

// ANCHOR_END: IfdtoolEntry

// ANCHOR: FirmwareStitchingOpts

// FirmwareStitchingOpts is used to store all data needed to stitch firmware
type FirmwareStitchingOpts struct {
	// List of IDs this instance depends on
	Depends []string `json:"depends"`

	// Common options like paths etc.
	CommonOpts

	// BaseFile into which inject files.
	// !!! Must contain IFD !!!
	// Examples:
	//   - coreboot.rom
	//   - ifd.bin
	BaseFilePath string `json:"base_file_path" validate:"required,filepath"`

	// Platform - passed to all `ifdtool` calls with `--platform`
	Platform string `json:"platform"`

	// List of instructions for ifdtool
	IfdtoolEntries []IfdtoolEntry `json:"ifdtool_entries"`

	// List of instructions for cbfstool
	// TODO ???
}

// ANCHOR_END: FirmwareStitchingOpts

// GetDepends is used to return list of dependencies
func (opts FirmwareStitchingOpts) GetDepends() []string {
	return opts.Depends
}

// GetArtifacts returns list of wanted artifacts from container
func (opts FirmwareStitchingOpts) GetArtifacts() *[]container.Artifacts {
	return opts.CommonOpts.GetArtifacts()
}

// ExtractSizeFromString uses regex to find size of ROM in MB
func ExtractSizeFromString(text string) ([]uint64, error) {
	// Component 1 and 2 represent flash chips on motherboard
	// 1st is a must, 2nd is optional
	// Example:
	//   "  Component 2 Density:                 32MB"
	//   "  Component 1 Density:                 64MB"
	// FindSubmatch:
	//   "  Component 1 Density:                 64MB"
	//      ^-----------------^:^---------------^^--^
	//       %s                : \s*              (\w+)
	items := []string{
		"Component 1 Density",
		"Component 2 Density",
	}
	results := []uint64{}
	for _, item := range items {
		re := regexp.MustCompile(fmt.Sprintf("%s:\\s*(\\w+)", item))
		matches := re.FindSubmatch([]byte(text))
		if len(matches) >= 1 {
			size, err := StringToSizeMB(string(matches[1]))
			if err != nil {
				return []uint64{}, err
			}
			results = append(results, size)
		} else {
			return []uint64{}, fmt.Errorf("could not find '%s' in ifdtool dump: %w", item, errFailedToDetectRomSize)
		}
	}
	return results, nil
}

// StringToSizeMB parses string and returns size in MB
func StringToSizeMB(text string) (uint64, error) {
	// Check for UNUSED
	if strings.ToLower(text) == "unused" {
		return 0, nil
	}

	// Cleanup string
	re := regexp.MustCompile(`\s+`)
	text = string(re.ReplaceAll([]byte(text), []byte("")))

	// Parse integer
	reUnits := regexp.MustCompile(`([kMGT])B`)
	numberString := reUnits.ReplaceAll([]byte(text), []byte("${1}iB"))
	number, err := humanize.ParseBytes(string(numberString))
	if err != nil {
		return 0, errFailedToDetectRomSize
	}

	return number, nil
}

// assemble command for ifdtool
func ifdtoolCmd(platform string, arguments []string) []string {
	cmd := []string{ifdtoolPath}
	if platform != "" {
		// TODO: Wanted to expand this to --platform
		//   but ifdtool has a bug in this long flag
		//   https://review.coreboot.org/c/coreboot/+/80432
		cmd = append(cmd, []string{"-p", platform}[:]...)
	}
	cmd = append(cmd, arguments[:]...)
	return cmd
}

// buildFirmware builds coreboot with all blobs and stuff
func (opts FirmwareStitchingOpts) buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) (*dagger.Container, error) {
	// Check that all files have unique filenames (they are copied into the same dir)
	copiedFiles := map[string]string{}
	for _, entry := range opts.IfdtoolEntries {
		filename := filepath.Base(entry.Path)
		if _, ok := copiedFiles[filename]; ok {
			slog.Error(
				fmt.Sprintf("File '%s' and '%s' have the same filename", entry.Path, copiedFiles[filename]),
				slog.String("suggestion", "Each file must have a unique name because they get copied into single directory"),
				slog.Any("error", os.ErrExist),
			)
			return nil, os.ErrExist
		}
		copiedFiles[filename] = entry.Path
	}

	// Spin up container
	containerOpts := container.SetupOpts{
		ContainerURL:      opts.SdkURL,
		MountContainerDir: ContainerWorkDir,
		MountHostDir:      opts.RepoPath,
		WorkdirContainer:  ContainerWorkDir,
	}
	myContainer, err := container.Setup(ctx, client, &containerOpts, dockerfileDirectoryPath)
	if err != nil {
		slog.Error(
			"Failed to start a container",
			slog.Any("error", err),
		)
		return nil, err
	}

	// Copy all the files into container
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error(
			"Could not get working directory",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return nil, err
	}
	newBaseFilePath := filepath.Join(ContainerWorkDir, filepath.Base(opts.BaseFilePath))
	myContainer = myContainer.WithFile(
		newBaseFilePath,
		client.Host().File(filepath.Join(pwd, opts.BaseFilePath)),
	)
	oldBaseFilePath := opts.BaseFilePath
	opts.BaseFilePath = newBaseFilePath
	for entry := range opts.IfdtoolEntries {
		containerPath := filepath.Join(ContainerWorkDir, filepath.Base(opts.IfdtoolEntries[entry].Path))
		hostPath := filepath.Join(pwd, opts.IfdtoolEntries[entry].Path)
		hostFile := client.Host().File(hostPath)

		// Check if the file exists on host filesystem
		_, err := os.Stat(hostPath)
		if err == nil {
			myContainer = myContainer.WithFile(
				containerPath,
				hostFile,
			)
			opts.IfdtoolEntries[entry].Path = containerPath
		} else if opts.IfdtoolEntries[entry].IgnoreIfMissing {
			// We can ignore this missing file
			opts.IfdtoolEntries[entry].Skip = true
			slog.Warn(
				fmt.Sprintf("Can't copy file '%s' - does not exists, ignoring because 'ignore_if_missing' is set", opts.IfdtoolEntries[entry].Path),
			)
		} else {
			// We cannot ignore this missing file
			slog.Error(
				fmt.Sprintf("Can't copy file '%s' - does not exists", opts.IfdtoolEntries[entry].Path),
				slog.String("suggestion", "Double check provided path to file"),
				slog.Any("error", err),
			)
			return nil, err
		}
	}

	// Get the size of image (total size)
	cmd := ifdtoolCmd(opts.Platform, []string{"--dump", opts.BaseFilePath})
	myContainerPrevious := myContainer
	ifdtoolStdout, err := myContainer.WithExec(cmd).Stdout(ctx)
	if err != nil {
		slog.Error(
			"Failed to dump Intel Firmware Descriptor (IFD)",
			slog.Any("error", err),
		)
		return myContainerPrevious, err
	}
	size, err := ExtractSizeFromString(ifdtoolStdout)
	if err != nil {
		slog.Error(
			"Failed extract size from Intel Firmware Descriptor (IFD)",
			slog.Any("error", err),
		)
		return nil, err
	}
	var totalSize uint64
	for _, i := range size {
		totalSize += i
	}
	slog.Info(
		fmt.Sprintf("Intel Firmware Descriptor (IFD) detected size: %s B", humanize.Comma(int64(totalSize))),
	)

	// Read the base file
	baseFile, err := os.ReadFile(oldBaseFilePath)
	if err != nil {
		return nil, err
	}
	baseFileSize := uint64(len(baseFile))
	slog.Info(
		fmt.Sprintf("Size of '%s': %s B", filepath.Base(oldBaseFilePath), humanize.Comma(int64(baseFileSize))),
	)
	if baseFileSize > totalSize {
		err = errBaseFileBiggerThanIfd
		slog.Error(
			fmt.Sprintf("Provided base_file '%s' is bigger (%s B) than defined in IFD (%s B)",
				filepath.Base(oldBaseFilePath),
				humanize.Comma(int64(baseFileSize)),
				humanize.Comma(int64(totalSize)),
			),
			slog.Any("error", err),
		)
		return nil, err
	}

	// Take baseFile content and expand it to correct size
	//   fill the empty space with 0xFF
	blank := make([]byte, totalSize-baseFileSize)
	for i := range blank {
		blank[i] = 0xFF
	}
	firmwareImage := []byte{}
	firmwareImage = append(firmwareImage, baseFile[:]...)
	firmwareImage = append(firmwareImage, blank[:]...)

	imageFilename := fmt.Sprintf("new_%s", filepath.Base(opts.BaseFilePath))
	slog.Info(
		fmt.Sprintf(
			"File '%s' is being expanded to ROM size %s B as '%s'",
			filepath.Base(opts.BaseFilePath),
			humanize.Comma(int64(len(firmwareImage))),
			imageFilename,
		),
	)
	firmwareImageFile, err := os.Create(imageFilename)
	if err != nil {
		return nil, err
	}
	_, err = firmwareImageFile.Write(firmwareImage)
	if err != nil {
		return nil, err
	}
	firmwareImageFile.Close()
	myContainer = myContainer.WithFile(
		filepath.Join(ContainerWorkDir, imageFilename),
		client.Host().File(filepath.Join(pwd, imageFilename)),
	)

	// Populate regions with ifdtool
	for entry := range opts.IfdtoolEntries {
		slog.Info(
			fmt.Sprintf("Injecting '%s' into '%s' region in '%s'",
				opts.IfdtoolEntries[entry].Path,
				opts.IfdtoolEntries[entry].TargetRegion,
				imageFilename,
			),
		)

		// Check if file exists, and if missing file can be ignored
		if opts.IfdtoolEntries[entry].Skip {
			slog.Warn(
				fmt.Sprintf("Can't inject file '%s' - does not exists, ignoring because 'ignore_if_missing' is set", opts.IfdtoolEntries[entry].Path),
			)
			continue
		}

		// Inject binaries
		cmd := ifdtoolCmd(
			opts.Platform,
			[]string{
				"--inject",
				fmt.Sprintf("%s:%s",
					opts.IfdtoolEntries[entry].TargetRegion,
					opts.IfdtoolEntries[entry].Path),
				imageFilename,
			},
		)
		myContainerPrevious = myContainer
		myContainer, err = myContainer.WithExec(cmd).Sync(ctx)
		if err != nil {
			slog.Error("Failed to inject region")
			return myContainerPrevious, err
		}

		// ifdtool makes a new file '<filename>.new', so let's rename back to original name
		imageFilenameNew := fmt.Sprintf("%s.new", imageFilename)
		cmd = []string{"mv", "--force", imageFilenameNew, imageFilename}
		myContainerPrevious = myContainer
		myContainer, err = myContainer.WithExec(cmd).Sync(ctx)
		if err != nil {
			slog.Error(
				fmt.Sprintf("Failed to rename '%s' to '%s'", imageFilenameNew, imageFilename),
			)
			return myContainerPrevious, err
		}
	}

	// Extract artifacts
	return myContainer, container.GetArtifacts(ctx, myContainer, opts.CommonOpts.GetArtifacts())
}
