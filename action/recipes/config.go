// SPDX-License-Identifier: MIT

// Package recipes / config
package recipes

import (
	"context"
	"encoding/json"
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
	"github.com/go-playground/validator/v10"
)

var (
	// ErrVerboseJSON is raised when JSONVerboseError can't find location of problem in JSON configuration file
	ErrVerboseJSON = errors.New("unable to pinpoint the problem in JSON file")
	// ErrEnvVarUndefined is raised when undefined environment variable is found in JSON configuration file
	ErrEnvVarUndefined = errors.New("environment variable used in JSON file is not present in the environment")
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
	//   Will result in following structure being copied out of the container:
	//     myOutput/
	//     ├── Build/
	//     ├── coreboot.rom
	//     └── defconfig
	OutputDir string `json:"output_dir" validate:"required,dirpath"`

	// Specifies the (relative) paths to directories which should be copied into the container.
	InputDirs []string `json:"input_dirs" validate:"dive,dirpath"`

	// Specifies the (relative) paths to file which should be copied into the container.
	InputFiles []string `json:"input_files" validate:"dive,filepath"`

	// Specifies the path to directory where to place input files and directories inside container.
	//   Directories listed in ContainerInputDirs and files listed in ContainerInputFiles
	//   will be copied there.
	// Example:
	//   Following setting:
	//     InputDirs = []string{"config-files/"}
	//     InputFiles = []string{"README.md", "Taskfile.yml"}
	//     ContainerInputDir = "myInput"
	//   Will result in following structure being copied into the container:
	//     myInput/
	//     ├── config-files/
	//     ├── README.md
	//     └── Taskfile.yml
	ContainerInputDir string `json:"container_input_dir" validate:"dirpath"`

	// Overview:
	//
	// | Configuration option   | Host side              | Direction            | Container side                 |
	// |:-----------------------|:-----------------------|:--------------------:|:-------------------------------|
	// | RepoPath               | $RepoPath              | Host  --> Container  | $(pwd)                         |
	// |                        |                        |                      |                                |
	// | OutputDir              | $(pwd)/$OutputDir      | Host <--  Container  | N/A                            |
	// | ContainerOutputDirs    | $(pwd)/$OutputDir/...  | Host <--  Container  | $ContainerOutputDirs           |
	// | ContainerOutputFiles   | $(pwd)/$OutputDir/...  | Host <--  Container  | $ContainerOutputFiles          |
	// |                        |                        |                      |                                |
	// | ContainerInputDir      | N/A                    | Host  --> Container  | $(pwd)/$ContainerInputDir      |
	// | InputDirs              | $InputDirs             | Host  --> Container  | $(pwd)/$ContainerInputDir/...  |
	// | InputFiles             | $InputFiles            | Host  --> Container  | $(pwd)/$ContainerInputDir/...  |
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

// GetContainerOutputDirs returns list of output directories
func (opts CommonOpts) GetContainerOutputDirs() []string {
	return opts.ContainerOutputDirs
}

// GetContainerOutputFiles returns list of output directories
func (opts CommonOpts) GetContainerOutputFiles() []string {
	return opts.ContainerOutputFiles
}

// GetOutputDir returns output directory
func (opts CommonOpts) GetOutputDir() string {
	return opts.OutputDir
}

// GetSources returns slice of paths to all sources which are used for build
func (opts CommonOpts) GetSources() []string {
	sources := []string{}

	// Repository path
	sources = append(sources, opts.RepoPath)

	// Input files and directories
	sources = append(sources, opts.InputDirs[:]...)
	sources = append(sources, opts.InputFiles[:]...)

	return sources
}

// Config is for storing parsed configuration file
type Config struct {
	// defined in coreboot.go
	Coreboot map[string]CorebootOpts `json:"coreboot" validate:"dive"`

	// defined in linux.go
	Linux map[string]LinuxOpts `json:"linux" validate:"dive"`

	// defined in edk2.go
	Edk2 map[string]Edk2Opts `json:"edk2" validate:"dive"`

	// defined in stitching.go
	FirmwareStitching map[string]FirmwareStitchingOpts `json:"firmware_stitching" validate:"dive"`

	// defined in uroot.go
	URoot map[string]URootOpts `json:"u-root" validate:"dive"`

	// defined in universal.go
	Universal map[string]UniversalOpts `json:"universal" validate:"dive"`
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
	for key, value := range c.FirmwareStitching {
		modules[key] = value
	}
	for key, value := range c.URoot {
		modules[key] = value
	}
	for key, value := range c.Universal {
		modules[key] = value
	}
	return modules
}

// FirmwareModule interface
type FirmwareModule interface {
	GetDepends() []string
	GetArtifacts() *[]container.Artifacts
	GetContainerOutputDirs() []string
	GetContainerOutputFiles() []string
	GetOutputDir() string
	GetSources() []string
	buildFirmware(ctx context.Context, client *dagger.Client, dockerfileDirectoryPath string) (*dagger.Container, error)
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
		err = errors.Join(ErrFailedValidation, err)
		slog.Error(
			"Configuration file failed validation",
			slog.String("suggestion", "Double check the used configuration file"),
			slog.Any("error", err),
		)
		return err
	}
	return nil
}

// FindAllEnvVars returns all environment variables found in the provided string
func FindAllEnvVars(text string) []string {
	pattern := regexp.MustCompile(`\${?([a-zA-Z0-9_]+)}?`)
	result := pattern.FindAllString(text, -1)
	for index, value := range result {
		result[index] = pattern.ReplaceAllString(value, "$1")
	}
	return result
}

// ReadConfig is for reading and parsing JSON configuration file into Config struct
func ReadConfig(filepath string) (*Config, error) {
	// Read JSON file
	content, err := os.ReadFile(filepath)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Unable to open the configuration file '%s'", filepath),
			slog.Any("error", err),
		)
		return nil, err
	}

	contentStr := string(content)

	// Check if all environment variables are defined
	envVars := FindAllEnvVars(contentStr)
	undefinedVarFound := false
	for _, envVar := range envVars {
		_, found := os.LookupEnv(envVar)
		if !found {
			slog.Error(
				fmt.Sprintf("environment variable '%s' is undefined", envVar),
				slog.String("suggestion", "define the environment variable in the environment"),
				slog.Any("error", ErrEnvVarUndefined),
			)
			undefinedVarFound = true
		}
	}
	if undefinedVarFound {
		return nil, ErrEnvVarUndefined
	}

	// Expand environment variables
	contentStr = os.ExpandEnv(contentStr)

	// Decode JSON
	jsonDecoder := json.NewDecoder(strings.NewReader(contentStr))
	jsonDecoder.DisallowUnknownFields()
	// jsonDecoder will return error when contentStr has keys not matching fields in Config struct
	var payload Config
	err = jsonDecoder.Decode(&payload)
	if err != nil {
		JSONVerboseError(contentStr, err)
		return nil, err
	}

	// Validate config
	err = ValidateConfig(payload)
	if err != nil {
		// no slog.Error because already called in ValidateConfig
		return nil, err
	}

	return &payload, nil
}

// WriteConfig is for writing Config struct into JSON configuration file
func WriteConfig(filepath string, config *Config) error {
	// Generate JSON
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		slog.Error(
			"Unable to convert the configuration into a JSON string",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return err
	}

	// Write JSON to file
	if err := os.WriteFile(filepath, b, 0o666); err != nil {
		slog.Error(
			"Failed to write configuration into JSON file",
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

// JSONVerboseError is for getting more information out of json.Unmarshal() or Decoder.Decode()
//
//	Inspiration:
//	- https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors/
//	Docs:
//	- https://pkg.go.dev/encoding/json#Unmarshal
func JSONVerboseError(jsonString string, err error) {
	if jsonError, ok := err.(*json.SyntaxError); ok {
		// JSON-encoded data contain a syntax error
		line, character, _ := offsetToLineNumber(jsonString, int(jsonError.Offset))
		slog.Error(
			// https://pkg.go.dev/encoding/json#SyntaxError
			fmt.Sprintf("Syntax error at line %d, character %d", line, character),
			slog.Any("error", jsonError.Error()),
		)
		return
	}
	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		// JSON value is not appropriate for a given target type
		line, character, _ := offsetToLineNumber(jsonString, int(jsonError.Offset))
		slog.Error(
			fmt.Sprintf(
				"Expected type '%v', JSON contains field '%v' in struct '%s' instead (full path: %s), see line %d, character %d",
				// https://pkg.go.dev/encoding/json#UnmarshalTypeError
				jsonError.Type.Name(), // Go type
				jsonError.Value,       // JSON field type
				jsonError.Struct,      // Name of struct type containing the field
				jsonError.Field,       // the full path from root node to the field
				line,
				character,
			),
			slog.Any("error", jsonError.Error()),
		)
		return
	}
	slog.Error(
		"Sorry but could not pinpoint specific location of the problem in the JSON configuration file",
		slog.Any("error", err),
	)
}

func offsetToLineNumber(input string, offset int) (line int, character int, err error) {
	// NOTE: I do not take into account windows line endings
	//       I can't be bothered, the worst case is that with windows line-endings the character counter
	//       will be off by 1, which is a sacrifice I am willing to make

	if offset > len(input) || offset < 0 {
		err = fmt.Errorf("offset is out of bounds for given string: %w", ErrVerboseJSON)
		slog.Warn(
			"Failed to pinpoint exact location of error in JSON configuration file",
			slog.Any("error", err),
		)
		return 0, 0, err
	}

	line = 1
	character = 0
	for index, char := range input {
		if char == '\n' {
			line++
			character = 0
			continue
		}
		character++
		if index >= offset {
			break
		}
	}

	return
}
