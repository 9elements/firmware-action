// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/heimdalr/dag"
)

// Errors for recipes
var (
	ErrFailedValidation          = errors.New("config failed validation")
	ErrTargetMissing             = errors.New("no target specified")
	ErrTargetInvalid             = errors.New("unsupported target")
	ErrBuildFailed               = errors.New("build failed")
	ErrDependencyTreeUndefDep    = errors.New("module has invalid dependency")
	ErrDependencyTreeUnderTarget = errors.New("target not found in dependency tree")
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

// Build recipes, possibly recursively
func Build(
	ctx context.Context,
	target string,
	recursive bool,
	interactive bool,
	config *Config,
	executor func(context.Context, string, *Config, bool) error,
) ([]string, error) {
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
	if recursive {
		builds := []string{}
		slog.Info(fmt.Sprintf("Building '%s' recursively", target))
		for _, item := range queue {
			slog.Info(fmt.Sprintf("Building: %s", item))
			err = executor(ctx, item, config, interactive)
			if err != nil {
				return nil, err
			}
			builds = append(builds, item)
		}
		return builds, nil
	}
	// else build only the target
	slog.Info(fmt.Sprintf("Building '%s' NOT recursively", target))
	return []string{target}, executor(ctx, target, config, interactive)
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
