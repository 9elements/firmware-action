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
	"reflect"
	"regexp"
	"strings"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/container"
	"github.com/9elements/firmware-action/cmd/firmware-action/logging"
	"github.com/go-playground/validator/v10"
)

var (
	// ErrVerboseJSON is raised when JSONVerboseError can't find location of problem in JSON configuration file
	ErrVerboseJSON = errors.New("unable to pinpoint the problem in JSON file")
	// ErrEnvVarUndefined is raised when undefined environment variable is found in JSON configuration file
	ErrEnvVarUndefined = errors.New("environment variable used in JSON file is not present in the environment")
	// ErrNestedOutputDirs is raised when one module's output directory is a subdirectory of another module's output directory
	ErrNestedOutputDirs = errors.New("nested output directories detected")
	// ErrDuplicateOutputDirs is raised when multiple modules use the same output directory
	ErrDuplicateOutputDirs = errors.New("duplicate output directories detected")
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
	// Can also be a absolute or relative path to Dockerfile to build the image on the fly.
	// NOTE: Updating the sdk_version might result in different binaries using the
	//   same source code.
	// ANCHOR: CommonOptsSdkURLExamples
	// Examples:
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:main
	//   https://ghcr.io/9elements/firmware-action/coreboot_4.19:latest
	//   https://ghcr.io/9elements/firmware-action/edk2-stable202111:latest
	//   file://./my-image/Dockerfile
	//   file://./my-image/
	//   file://my-image/Dockerfile
	//   file:///home/user/my-image/Dockerfile
	//   file:///home/user/my-image/
	// ANCHOR_END: CommonOptsSdkURLExamples
	// NOTE:
	//   'file://' path cannot contain '..'
	// See https://github.com/orgs/9elements/packages
	SdkURL string `json:"sdk_url" validate:"required"`

	// Gives the (relative) path to the target (firmware) repository.
	// If the current repository contains the selected target, specify: '.'
	// Otherwise the path should point to the target (firmware) repository submodule that
	//   had been previously checked out.
	RepoPath string `json:"repo_path" validate:"required,dirpath"`

	// Specifies the (relative) paths to directories where are produced files (inside Container).
	ContainerOutputDirs []string `json:"container_output_dirs" validate:"dive,filepath|dirpath"`

	// Specifies the (relative) paths to produced files (inside Container).
	ContainerOutputFiles []string `json:"container_output_files" validate:"dive,filepath|dirpath"`

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
	OutputDir string `json:"output_dir" validate:"required,filepath|dirpath"`

	// Specifies the (relative) paths to directories which should be copied into the container.
	InputDirs []string `json:"input_dirs" validate:"dive,filepath|dirpath"`

	// Specifies the (relative) paths to file which should be copied into the container.
	InputFiles []string `json:"input_files" validate:"dive,filepath|dirpath"`

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
	ContainerInputDir string `json:"container_input_dir" validate:"filepath|dirpath"`

	// Overview:
	//   NOTE: $PWD in the container is /workdir
	//   defined in recipes.go with "ContainerWorkDir"
	//
	// | Configuration option   | Host side              | Direction            | Container side                   |
	// |:-----------------------|:-----------------------|:--------------------:|:---------------------------------|
	// | RepoPath               | $RepoPath              | Host  --> Container  | /workdir                         |
	// |                        |                        |                      |                                  |
	// | OutputDir              | $(pwd)/$OutputDir      | Host <--  Container  | N/A                              |
	// | ContainerOutputDirs    | $(pwd)/$OutputDir/...  | Host <--  Container  | $ContainerOutputDirs             |
	// | ContainerOutputFiles   | $(pwd)/$OutputDir/...  | Host <--  Container  | $ContainerOutputFiles            |
	// |                        |                        |                      |                                  |
	// | ContainerInputDir      | N/A                    | Host  --> Container  | /workdir/$ContainerInputDir      |
	// | InputDirs              | $InputDirs             | Host  --> Container  | /workdir/$ContainerInputDir/...  |
	// | InputFiles             | $InputFiles            | Host  --> Container  | /workdir/$ContainerInputDir/...  |
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

// ANCHOR: CommonOptsGetSources

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

// ANCHOR_END: CommonOptsGetSources

// GetRepoPath returns Repository path
func (opts CommonOpts) GetRepoPath() string {
	return opts.RepoPath
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

	// defined in uboot.go
	UBoot map[string]UBootOpts `json:"u-boot" validate:"dive"`
}

// AllModules method returns slice with all modules
func (c Config) AllModules() map[string]FirmwareModule {
	modules := make(map[string]FirmwareModule)

	configValue := reflect.ValueOf(c)

	for i := range configValue.Type().NumField() {
		fieldValue := configValue.Field(i)

		// Check if the field is a map.
		if fieldValue.Kind() == reflect.Map {
			// Iterate over the keys in the map.
			for _, key := range fieldValue.MapKeys() {
				value := fieldValue.MapIndex(key)

				// Type-assert the value to FirmwareModule.
				if module, ok := value.Interface().(FirmwareModule); ok {
					modules[key.String()] = module
				} else {
					slog.Error(
						fmt.Sprintf("Value for key '%v' in config does not implement FirmwareModule", key),
						slog.String("suggestion", logging.ThisShouldNotHappenMessage),
					)
				}
			}
		}
	}

	return modules
}

// Merge method will take other Config instance and adopt all of its modules
func (c Config) Merge(other Config) (Config, error) {
	merged := Config{}

	// Use reflection on the merged instance.
	vMerged := reflect.ValueOf(&merged).Elem()
	vC := reflect.ValueOf(c)
	vOther := reflect.ValueOf(other)
	t := vMerged.Type()

	// Iterate over all fields of the struct.
	for i := range t.NumField() {
		fieldType := t.Field(i)
		// Process only map fields.
		if fieldType.Type.Kind() == reflect.Map {
			// Create a new map for the merged result.
			mergedMap := reflect.MakeMap(fieldType.Type)

			// Get the map from c (receiver) and copy its key/value pairs.
			mapC := vC.Field(i)
			if mapC.IsValid() && !mapC.IsNil() {
				for _, key := range mapC.MapKeys() {
					mergedMap.SetMapIndex(key, mapC.MapIndex(key))
				}
			}

			// Get the map from other and merge its entries.
			mapOther := vOther.Field(i)
			if mapOther.IsValid() && !mapOther.IsNil() {
				for _, key := range mapOther.MapKeys() {
					// If the key already exists, print a warning.
					if existing := mergedMap.MapIndex(key); existing.IsValid() {
						fmt.Printf("Warning: overriding key %v in field %s\n", key, fieldType.Name)
					}
					mergedMap.SetMapIndex(key, mapOther.MapIndex(key))
				}
			}
			// Set the merged map into the new struct.
			vMerged.Field(i).Set(mergedMap)
		} else {
			// For non-map fields, just copy the value from c.
			vMerged.Field(i).Set(vC.Field(i))
		}
	}

	return merged, nil
}

// FirmwareModule interface
type FirmwareModule interface {
	GetDepends() []string
	GetArtifacts() *[]container.Artifacts
	GetContainerOutputDirs() []string
	GetContainerOutputFiles() []string
	GetOutputDir() string
	GetSources() []string
	buildFirmware(ctx context.Context, client *dagger.Client) error
	GetRepoPath() string
}

// ======================
//  Functions for Config
// ======================

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

	// Check for nested/duplicate output directories
	return validateOutputDirectories(conf.AllModules())
}

// validateOutputDirectories checks for nested or duplicate output directories
func validateOutputDirectories(modules map[string]FirmwareModule) error {
	// Check for nested output directories
	//
	// Here are few examples:
	//   VALID:
	//     .
	//     ├── output-linux
	//     │   └── linux.bin
	//     └── output-uroot
	//         └── uroot.bin
	//   VALID:
	//     .
	//     └── output
	//         ├── output-linux
	//         │   └── linux.bin
	//         └── output-uroot
	//             └── uroot.bin
	//   INVALID:
	//     .
	//     └── output
	//         ├── output-linux
	//         │   └── linux.bin
	//         └── uroot.bin

	outputDirs := make(map[string]string) // Map of output dir -> module name
	var issues []error

	// Check for duplicate output directories
	for moduleName, module := range modules {
		outputDir := filepath.Clean(module.GetOutputDir())

		// Check if this output directory already exists in our map
		if existingModule, exists := outputDirs[outputDir]; exists {
			// We found a duplicate output directory
			errMsg := fmt.Sprintf("modules '%s' and '%s' have the same output directory '%s'",
				existingModule, moduleName, outputDir)
			err := fmt.Errorf("%w: %s", ErrDuplicateOutputDirs, errMsg)
			slog.Error(
				"Detected duplicate output directories",
				slog.String("suggestion", "Please make sure that each module has its own unique output directory. Each module needs exclusive control over its output directory."),
				slog.Any("error", err),
			)
			issues = append(issues, err)
		}

		// Add this output directory to our map
		outputDirs[outputDir] = moduleName
	}

	// Check if any output directory is a subdirectory of another
	for dir1, module1 := range outputDirs {
		for dir2, module2 := range outputDirs {
			// Skip comparing to itself
			if dir1 == dir2 {
				continue
			}

			// Check for nesting
			// adding `filepath.Separator` is necessary because of the `strings.HasPrefix` to avoid false positives
			dirSep := string(filepath.Separator)
			if strings.HasPrefix(dir1+dirSep, dir2+dirSep) {
				errMsg := fmt.Sprintf("output directory '%s' of module '%s' is a subdirectory of '%s' from module '%s'",
					dir1, module1, dir2, module2)
				err := fmt.Errorf("%w: %s", ErrNestedOutputDirs, errMsg)
				slog.Error(
					"Detected nested output directories",
					slog.String("suggestion", "Please make sure that each module has its own unique output directory. Each module needs exclusive control over its output directory. This directory is deleted when changes are detected and re-build is required."),
					slog.Any("error", err),
				)
				issues = append(issues, err)
			}
		}
	}

	// If we found any issues, return a combined error
	if len(issues) > 0 {
		var combinedErr error
		for _, err := range issues {
			combinedErr = errors.Join(combinedErr, err)
		}
		return combinedErr
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

// ReadConfigs is for reading and parsing multiple JSON configuration files into single Config struct
func ReadConfigs(filepaths []string) (*Config, error) {
	var allConfigs Config
	for _, filepath := range filepaths {
		trimmedFilepath := strings.TrimSpace(filepath)
		slog.Debug("Reading config",
			slog.String("path", trimmedFilepath),
		)
		payload, err := ReadConfig(trimmedFilepath)
		if err != nil {
			return nil, err
		}
		allConfigs, err = allConfigs.Merge(*payload)
		if err != nil {
			return nil, err
		}
	}
	return &allConfigs, nil
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
