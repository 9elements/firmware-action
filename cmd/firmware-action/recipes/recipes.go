// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/container"
	"github.com/9elements/firmware-action/filesystem"
	"github.com/heimdalr/dag"
)

// Errors for recipes
var (
	ErrBuildFailed               = errors.New("build failed")
	ErrBuildUpToDate             = errors.New("build is up-to-date")
	ErrDependencyTreeUndefDep    = errors.New("module has invalid dependency")
	ErrDependencyTreeUnderTarget = errors.New("target not found in dependency tree")
	ErrDependencyOutputMissing   = errors.New("output of one or more dependencies is missing")
	ErrFailedValidation          = errors.New("config failed validation")
	ErrTargetInvalid             = errors.New("unsupported target")
	ErrTargetMissing             = errors.New("no target specified")
)

var (
	// ContainerWorkDir specifies directory in container used as work directory
	ContainerWorkDir = "/workdir"
	// TimestampsDir specifies directory for timestamps to detect changes in sources
	TimestampsDir = ".firmware-action/timestamps"
)

func forestAddVertex(forest *dag.DAG, key string, value FirmwareModule, dependencies [][]string) ([][]string, error) {
	err := forest.AddVertexByID(key, key)
	if err != nil {
		return nil, err
	}
	for _, dep := range value.GetDepends() {
		dependencies = append(dependencies, []string{key, dep})
	}
	return dependencies, nil
}

// BuildResults contains target name and result of its build
type BuildResults struct {
	Name        string
	BuildResult error
}

// Build recipes, possibly recursively
func Build(
	ctx context.Context,
	target string,
	recursive bool,
	interactive bool,
	config *Config,
	executor func(context.Context, string, *Config, bool) error,
) ([]BuildResults, error) {
	dependencyForest := dag.NewDAG()
	dependencies := [][]string{}
	var err error

	// Create the forest (forest = multiple independent trees)
	//   Add all items as vertexes into the tree
	for key, value := range config.AllModules() {
		dependencies, err = forestAddVertex(dependencyForest, key, value, dependencies)
		if err != nil {
			return nil, err
		}
	}

	// Add edges
	//   Edges must be added after all vertexes were are added
	for _, dep := range dependencies {
		err = dependencyForest.AddEdge(dep[0], dep[1])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrDependencyTreeUndefDep, err)
		}
	}

	// Check target is in Forest
	_, err = dependencyForest.GetVertex(target)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDependencyTreeUnderTarget, err)
	}

	// Create a queue in correct order (starting with leaves)
	queue := []string{}
	queueMutex := &sync.Mutex{} // Mutex to ensure concurrent access to queue is safe in the callback
	flowCallback := func(d *dag.DAG, id string, _ []dag.FlowResult) (interface{}, error) {
		v, err := d.GetVertex(id)
		if err != nil {
			return nil, err
		}
		queueMutex.Lock()
		queue = append(queue, v.(string))
		queueMutex.Unlock()
		return nil, nil
	}

	// Create a subgraph with target as the only root
	// Having multiple roots will result in miscalculation inside DescendantsFlow channel size calculation
	pruned, rootID, err := dependencyForest.GetDescendantsGraph(target)
	if err != nil {
		return nil, err
	}

	_, err = pruned.DescendantsFlow(rootID, nil, flowCallback)
	if err != nil {
		return nil, err
	}
	slices.Reverse(queue)

	// Build each item in queue (if recursive)
	slog.Info(fmt.Sprintf("Building queue: %v", queue))
	builds := []BuildResults{}
	if recursive {
		slog.Info(fmt.Sprintf("Building '%s' recursively", target))
		for _, item := range queue {
			slog.Info(fmt.Sprintf("Building: %s", item))

			err = executor(ctx, item, config, interactive)
			builds = append(builds, BuildResults{item, err})

			if err != nil && !errors.Is(err, ErrBuildUpToDate) {
				break
			}
		}
	} else {
		// else build only the target
		slog.Info(fmt.Sprintf("Building '%s' NOT recursively", target))

		err = executor(ctx, target, config, interactive)
		builds = append(builds, BuildResults{target, err})
	}

	// Check results
	err = nil
	for _, item := range builds {
		if item.BuildResult != nil && !errors.Is(item.BuildResult, ErrBuildUpToDate) {
			err = item.BuildResult
		}
	}

	return builds, err
}

// IsDirEmpty returns whether given directory is empty or not
func IsDirEmpty(path string) (bool, error) {
	// Source: https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// File.Readdirnames() take a parameter which is used to limit the number of returned values
	// It is enough to query only 1 child
	// File.Readdirnames() is faster than File.Readdir()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// Execute a build step
// func Execute(ctx context.Context, target string, config *Config, interactive bool, bulldozeMode bool) error {
func Execute(ctx context.Context, target string, config *Config, interactive bool) error {
	// Prep
	err := os.MkdirAll(TimestampsDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Find requested target
	modules := config.AllModules()
	if _, ok := modules[target]; ok {
		// Check if up-to-date
		// Either returns time, or zero time and error
		// zero time means there was no previous run
		timestampFile := filepath.Join(TimestampsDir, fmt.Sprintf("%s.txt", target))
		lastRun, _ := filesystem.LoadLastRunTime(timestampFile)

		sources := modules[target].GetSources()
		changesDetected := false
		for _, source := range sources {
			changes, _ := filesystem.AnyFileNewerThan(source, lastRun)
			if changes {
				changesDetected = true
				break
			}
		}

		// Check if output directory already exist
		// We want to skip build if the output directory exists and is not empty
		// If it is empty, then just continue with the building
		// If changes in sources were detected, re-build
		_, errExists := os.Stat(modules[target].GetOutputDir())
		empty, _ := IsDirEmpty(modules[target].GetOutputDir())
		if errExists == nil && !empty {
			if changesDetected {
				// If any of the sources changed, we need to rebuild
				os.RemoveAll(modules[target].GetOutputDir())
			} else {
				// Is already up-to-date
				slog.Warn(fmt.Sprintf("Target '%s' is up-to-date, skipping build", target))
				return ErrBuildUpToDate
			}
		}

		// Check if all outputs of required modules exist
		for _, prerequisite := range modules[target].GetDepends() {
			outputDir := modules[prerequisite].GetOutputDir()
			paths := modules[prerequisite].GetContainerOutputDirs()
			paths = append(paths, modules[prerequisite].GetContainerOutputFiles()...)

			for _, path := range paths {
				finalPath := filepath.Join(outputDir, filepath.Base(path))
				slog.Info(finalPath)
				if _, err := os.Stat(finalPath); os.IsNotExist(err) {
					slog.Error(
						"Missing output files and/or directories from one or more required module(s) defined in 'Depends'",
						slog.String("suggestion", "build needed modules or use '--recursive' build"),
						slog.Any("error", errors.Join(err, ErrDependencyOutputMissing)),
					)
					return ErrDependencyOutputMissing
				}
			}
		}

		// Setup dagger client
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		if err != nil {
			return err
		}
		defer client.Close()

		// Build the module
		myContainer, err := modules[target].buildFirmware(ctx, client, "")
		if err == nil {
			// On success update the timestamp
			_ = filesystem.SaveCurrentRunTime(timestampFile)
		}
		if err != nil && interactive {
			// If error, try to open SSH
			opts := container.NewSettingsSSH(container.WithWaitPressEnter())
			sshErr := container.OpenSSH(ctx, client, myContainer, ContainerWorkDir, opts)
			return errors.Join(err, sshErr)
		}
		return err
	}
	return ErrTargetMissing
}

// NormalizeArchitecture will translate various architecture strings into expected format
func NormalizeArchitecture(arch string) string {
	archMap := map[string]string{
		// x86 32-bit
		"IA-32":  "i386", // Intel
		"IA32":   "i386", // Intel
		"i686":   "i386", // common on Linux
		"386":    "i386", // GOARCH
		"x86":    "i386", // common on Windows
		"x86-32": "i386", // rare
		"x86_32": "i386", // rare
		// x86 64-bit
		"AMD64":  "amd64",
		"x64":    "amd64", // common on Windows
		"x86-64": "amd64",
		"x86_64": "amd64",
	}
	result, ok := archMap[arch]
	if result != "" && ok {
		return result
	}
	// fallback
	return arch
}

// NormalizeArchitectureForLinux will translate various architecture strings into format expected by Linux
func NormalizeArchitectureForLinux(arch string) string {
	normalArch := NormalizeArchitecture(arch)
	archMap := map[string]string{
		// x86 32-bit
		"i386": "x86",
		// x86 64-bit (x86_64 reuses x86)
		"amd64": "x86",
	}
	result, ok := archMap[normalArch]
	if result != "" && ok {
		return result
	}
	// fallback
	return arch
}
