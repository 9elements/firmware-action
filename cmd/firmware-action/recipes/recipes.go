// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/heimdalr/dag"
)

// Errors for recipes
var (
	ErrBuildFailed               = errors.New("build failed")
	ErrBuildSkipped              = errors.New("build skipped")
	ErrDependencyTreeUndefDep    = errors.New("module has invalid dependency")
	ErrDependencyTreeUnderTarget = errors.New("target not found in dependency tree")
	ErrDependencyOutputMissing   = errors.New("output of one or more dependencies is missing")
	ErrFailedValidation          = errors.New("config failed validation")
	ErrTargetInvalid             = errors.New("unsupported target")
	ErrTargetMissing             = errors.New("no target specified")
)

// ContainerWorkDir specifies directory in container used as work directory
var ContainerWorkDir = "/workdir"

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
	_, err = dependencyForest.DescendantsFlow(target, nil, flowCallback)
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

			if err != nil && !errors.Is(err, ErrBuildSkipped) {
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
		if item.BuildResult != nil && !errors.Is(item.BuildResult, ErrBuildSkipped) {
			err = item.BuildResult
		}
	}

	return builds, err
}

// Execute a build step
func Execute(ctx context.Context, target string, config *Config, interactive bool) error {
	// Setup dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Find requested target
	modules := config.AllModules()
	if _, ok := modules[target]; ok {
		// Check if output artifacts already exist
		for _, artifact := range *modules[target].GetArtifacts() {
			if _, err := os.Stat(artifact.HostPath); err == nil {
				slog.Warn(fmt.Sprintf("Output directory for '%s' already exists, skipping build", target))
				return ErrBuildSkipped
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

		// Build module
		myContainer, err := modules[target].buildFirmware(ctx, client, "")
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
