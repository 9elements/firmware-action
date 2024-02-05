// SPDX-License-Identifier: MIT

// Package recipes / config
package recipes

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/go-playground/validator/v10"
)

// =================
//  Data structures
// =================

// CommonOpts is common to all targets
// Used to store data from githubaction.Action
// For details see action.yml
// ANCHOR: CommonOpts
type CommonOpts struct {
	// Specifies the container toolchain tag to use when building the image.
	// This has an influence on the IASL, GCC and host GCC version that is used to build
	//   the target. You must match the source level and sdk_version.
	// NOTE: Updating the sdk_version might result in different binaries using the
	//   same source code.
	// Examples:
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:main
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:latest
	//   https://ghcr.io/9elements/firmware-action/edk2-stable202111:latest
	// See https://github.com/orgs/9elements/packages
	SdkURL string `json:"sdk_url" validate:"required"`

	// Gives the (relative) path to the target (firmware) repository.
	// If the current repository contains the selected target, specify: '.'
	// Otherwise the path should point to the target (firmware) repository submodule that
	//   had been previously checked out.
	RepoPath string `json:"repo_path" validate:"required,dirpath"`

	// Specifies the (relative) paths to directories where are produced files (inside Container).
	ContainerOutputDirs []string `json:"container_output_dirs" validate:"dive,dirpath"`

	// Specifies the (relative) paths to produced files (inside Container).
	ContainerOutputFiles []string `json:"container_output_files" validate:"dive,filepath"`

	// Specifies the (relative) path to directory into which place the produced files.
	//   Directories listed in ContainerOutputDirs and files listed in ContainerOutputFiles
	//   will be exported here.
	// Example:
	//   Following setting:
	//     ContainerOutputDirs = []string{"Build/"}
	//     ContainerOutputFiles = []string{"coreboot.rom", "defconfig"}
	//     OutputDir = "myOutput"
	//   Will result in:
	//     myOutput/
	//     ├── Build/
	//     ├── coreboot.rom
	//     └── defconfig
	OutputDir string `json:"output_dir" validate:"required,dirpath"`
}

// ANCHOR_END: CommonOpts

// GetArtifacts returns list of wanted artifacts from container
func (opts CommonOpts) GetArtifacts() *[]container.Artifacts {
	var artifacts []container.Artifacts

	// Directories
	for _, pathDir := range opts.ContainerOutputDirs {
		artifacts = append(artifacts, container.Artifacts{
			ContainerPath: filepath.Join(ContainerWorkDir, pathDir),
			ContainerDir:  true,
			HostPath:      opts.OutputDir,
			HostDir:       true,
		})
	}

	// Files
	for _, pathFile := range opts.ContainerOutputFiles {
		artifacts = append(artifacts, container.Artifacts{
			ContainerPath: filepath.Join(ContainerWorkDir, pathFile),
			ContainerDir:  false,
			HostPath:      opts.OutputDir,
			HostDir:       true,
		})
	}

	return &artifacts
}

// Config is for storing parsed configuration file
type Config struct {
	// defined in coreboot.go
	Coreboot map[string]CorebootOpts `json:"coreboot" validate:"dive"`

	// defined in linux.go
	Linux map[string]LinuxOpts `json:"linux" validate:"dive"`

	// defined in edk2.go
	Edk2 map[string]Edk2Opts `json:"edk2" validate:"dive"`
}

// AllModules method returns slice with all modules
func (c Config) AllModules() map[string]FirmwareModule {
	modules := make(map[string]FirmwareModule)
	for key, value := range c.Coreboot {
		modules[key] = value
	}
	for key, value := range c.Linux {
		modules[key] = value
	}
	for key, value := range c.Edk2 {
		modules[key] = value
	}
	return modules
}

// FirmwareModule interface
type FirmwareModule interface {
	GetDepends() []string
	GetArtifacts() *[]container.Artifacts
	buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) error
}

// ===========
//  Functions
// ===========

// ValidateConfig is used to validate the configuration struct read out of JSON file
func ValidateConfig(conf Config) error {
	// https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(conf)
	if err != nil {
		log.Print(err)
		return ErrRequiredOptionUndefined
	}
	return nil
}

// ReadConfig is for reading and parsing JSON configuration file into Config struct
func ReadConfig(filepath string) (Config, error) {
	// Read JSON file
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return Config{}, err
	}

	// Expand environment variables
	contentStr := string(content)
	contentStr = os.ExpandEnv(contentStr)
	content = []byte(contentStr)

	// Decode JSON
	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return Config{}, err
	}

	// Validate config
	err = ValidateConfig(payload)
	if err != nil {
		log.Print("Provided JSON configuration file failed validation")
		return Config{}, err
	}

	return payload, nil
}

// WriteConfig is for writing Config struct into JSON configuration file
func WriteConfig(filepath string, config Config) error {
	// Generate JSON
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal("Unable to convert the struct to a JSON string")
	}

	// Write JSON to file
	if err := os.WriteFile(filepath, b, 0o666); err != nil {
		log.Fatal(err)
	}

	return nil
}
