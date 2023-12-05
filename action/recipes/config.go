// SPDX-License-Identifier: MIT

// Package recipes / config
package recipes

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

// =================
//  Data structures
// =================

// CommonOpts is common to all targets
// Used to store data from githubaction.Action
// For details see action.yml
type CommonOpts struct {
	// Specifies the docker toolchain tag to use when building the image.
	// This has an influence on the IASL, GCC and host GCC version that is used to build
	//   the target. You must match the source level and sdk_version.
	// NOTE: Updating the sdk_version might result in different binaries using the
	//   same source code.
	// Examples:
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:main
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:latest
	//   https://ghcr.io/9elements/firmware-action/edk2-stable202111:latest
	// See https://github.com/orgs/9elements/packages
	SdkURL string `json:"sdk_url" validate:"required,url"`

	// Specifies target architecture, such as 'x86' or 'arm64'. Currently unused for coreboot.
	// Supported options for linux:
	//   - 'x86'
	//   - 'x86_64'
	//   - 'arm'
	//   - 'arm64'
	// Supported options for edk2:
	//   - 'AARCH64'
	//   - 'ARM'
	//   - 'IA32'
	//   - 'IA32X64'
	//   - 'X64'
	Arch string `json:"arch"`

	// Gives the (relative) path to the target (firmware) repository.
	// If the current repository contains the selected target, specify: '.'
	// Otherwise the path should point to the target (firmware) repository submodule that
	//   had been previously checked out.
	RepoPath string `json:"repo_path" validate:"required,dirpath"`

	// Gives the (relative) path to the defconfig that should be used to build the target.
	// For coreboot and linux this is a defconfig.
	// For EDK2 this is a one-line file containing the build arguments such as
	//   '-D BOOTLOADER=COREBOOT -D TPM_ENABLE=TRUE -D NETWORK_IPXE=TRUE'.
	//   Some arguments will be added automatically:
	//     '-a <architecture>'
	//     '-p <edk2__platform>'
	//     '-b <edk2__release_type>'
	//     '-t <GCC version>' (defined as part of docker toolchain, selected by SdkURL)
	DefconfigPath string `json:"defconfig_path" validate:"required,filepath"`

	// Specifies the (relative) path to directory into which place the produced files.
	OutputDir string `json:"output_dir" validate:"required,dirpath"`
}

// Config is for storing parsed configuration file
type Config struct {
	// defined in coreboot.go
	Coreboot []CorebootOpts `json:"coreboot" validate:"dive"`

	// defined in linux.go
	Linux []LinuxOpts `json:"linux" validate:"dive"`

	// defined in edk2.go
	Edk2 []Edk2Opts `json:"edk2" validate:"dive"`
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
func ReadConfig(filepath string) (*Config, error) {
	// Read JSON file
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return nil, err
	}

	// Decode JSON
	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return nil, err
	}

	// Validate config
	err = ValidateConfig(payload)
	if err != nil {
		log.Fatal("Provided JSON configuration file failed validation")
		return nil, err
	}

	return &payload, nil
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
