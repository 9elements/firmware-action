// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/heimdalr/dag"
)

// ErrRequiredOptionUndefined is raised when required option is empty or undefined
var (
	ErrRequiredOptionUndefined   = errors.New("required option is undefined")
	ErrTargetMissing             = errors.New("no target specified")
	ErrTargetInvalid             = errors.New("unsupported target")
	ErrBuildFailed               = errors.New("build failed")
	ErrDependencyTreeUndefDep    = errors.New("module has invalid dependency")
	ErrDependencyTreeUnderTarget = errors.New("target not found in dependency tree")
)

// ContainerWorkDir specifies directory in container used as work directory
var ContainerWorkDir = "/workdir"

type firmwareType interface {
	Dep() []string
}

// Dep to get dependencies
func (c CorebootOpts) Dep() []string {
	return c.Depends
}

// Dep to get dependencies
func (c LinuxOpts) Dep() []string {
	return c.Depends
}

// Dep to get dependencies
func (c Edk2Opts) Dep() []string {
	return c.Depends
}

func forestAddVertex(forest *dag.DAG, key string, value firmwareType, dependencies [][]string) ([][]string, error) {
	err := forest.AddVertexByID(key, key)
	if err != nil {
		return nil, err
	}
	for _, dep := range value.Dep() {
		dependencies = append(dependencies, []string{key, dep})
	}
	return dependencies, nil
}

// Build recipes, possibly recursively
func Build(ctx context.Context, target string, recursive bool, config Config, executor func(context.Context, string, Config) error) ([]string, error) {
	dependencyForest := dag.NewDAG()
	dependencies := [][]string{}
	var err error

	// Create the forest (forest = multiple independent trees)
	//   Add all items as vertexes into the tree
	//   There should be better way to do this other than doing 3x the same thing
	// -- coreboot --
	for key, value := range config.Coreboot {
		dependencies, err = forestAddVertex(dependencyForest, key, value, dependencies)
		if err != nil {
			return nil, err
		}
	}
	// -- linux --
	for key, value := range config.Linux {
		dependencies, err = forestAddVertex(dependencyForest, key, value, dependencies)
		if err != nil {
			return nil, err
		}
	}
	// -- edk2 --
	for key, value := range config.Edk2 {
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
	flowCallback := func(d *dag.DAG, id string, parentResults []dag.FlowResult) (interface{}, error) {
		v, err := d.GetVertex(id)
		if err != nil {
			return nil, err
		}
		queue = append(queue, v.(string))
		return nil, nil
	}
	_, err = dependencyForest.DescendantsFlow(target, nil, flowCallback)
	if err != nil {
		return nil, err
	}
	slices.Reverse(queue)

	// Build each item in queue (if recursive)
	log.Printf("building queue: %v", queue)
	if recursive {
		builds := []string{}
		log.Printf("building '%s' recursively", target)
		for _, item := range queue {
			log.Printf("- building %s", item)
			err = executor(ctx, item, config)
			if err != nil {
				return nil, err
			}
			builds = append(builds, item)
		}
		return builds, nil
	}
	// else build only the target
	log.Printf("building '%s' NOT recursively", target)
	return []string{target}, executor(ctx, target, config)
}

// Execute a build step
func Execute(ctx context.Context, target string, config Config) error {
	// Setup dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Find requested target
	if _, ok := config.Coreboot[target]; ok {
		// Coreboot
		opts := config.Coreboot[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "build", "coreboot.rom"),
				ContainerDir:  false,
				HostPath:      config.Coreboot[target].Common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      config.Coreboot[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return coreboot(ctx, client, &opts, "", &artifacts)
	} else if _, ok = config.Linux[target]; ok {
		// Linux
		opts := config.Linux[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "vmlinux"),
				ContainerDir:  false,
				HostPath:      config.Linux[target].Common.OutputDir,
				HostDir:       true,
			},
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "defconfig"),
				ContainerDir:  false,
				HostPath:      config.Linux[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return linux(ctx, client, &opts, "", &artifacts)
	} else if _, ok = config.Edk2[target]; ok {
		// Edk2
		opts := config.Edk2[target]
		artifacts := []container.Artifacts{
			{
				ContainerPath: filepath.Join(ContainerWorkDir, "Build"),
				ContainerDir:  true,
				HostPath:      config.Edk2[target].Common.OutputDir,
				HostDir:       true,
			},
		}
		return edk2(ctx, client, &opts, "", &artifacts)
	}
	return ErrTargetMissing
}
