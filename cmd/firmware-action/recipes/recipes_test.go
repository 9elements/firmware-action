// SPDX-License-Identifier: MIT
package recipes

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	const interactive = false
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name       string
		wantErr    error
		target     string
		targetType string
		config     Config
	}{
		{
			name:    "empty target string",
			wantErr: ErrTargetMissing,
			target:  "",
			config:  Config{},
		},
		{
			name:    "invalid target",
			wantErr: ErrTargetMissing,
			target:  "dummy",
			config:  Config{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = Execute(ctx, tc.target, &tc.config, interactive)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestExecuteSkipAndMissing(t *testing.T) {
	const interactive = false
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	// Change current working directory
	pwd, err := os.Getwd()
	defer os.Chdir(pwd) // nolint:errcheck
	assert.NoError(t, err)
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	// Create configuration
	const target = "dummy"
	const outputDir = "output-coreboot/"
	const depends = "pre-dummy"
	const outputDir2 = "output-coreboot2/"
	myConfig := Config{
		Coreboot: map[string]CorebootOpts{
			target: {
				Depends: []string{depends},
				CommonOpts: CommonOpts{
					OutputDir: outputDir,
					ContainerOutputFiles: []string{
						"build/coreboot.rom",
						"defconfig",
					},
				},
			},
			depends: {
				CommonOpts: CommonOpts{
					OutputDir: outputDir2,
					ContainerOutputFiles: []string{
						"build/coreboot.rom",
						"defconfig",
					},
				},
			},
		},
	}

	// Files from the 2nd modules are missing
	// This should fail since the 2nd module is in Depends
	err = Execute(ctx, target, &myConfig, interactive)
	assert.ErrorIs(t, err, ErrDependencyOutputMissing)

	// Create the output directory
	// Should build because the directory is empty
	err = os.Mkdir(outputDir, os.ModePerm)
	assert.NoError(t, err)
	err = Execute(ctx, target, &myConfig, interactive)
	assert.ErrorIs(t, err, ErrDependencyOutputMissing)
}

func executeDummy(_ context.Context, _ string, _ *Config, _ bool) error {
	return nil
}

func TestBuild(t *testing.T) {
	ctx := context.Background()

	testConfig := Config{
		Coreboot: map[string]CorebootOpts{
			"coreboot-0": {Depends: []string{}},
			"coreboot-A": {Depends: []string{"linux-A"}},
			"coreboot-B": {Depends: []string{"edk2-B"}},
			"coreboot-C": {Depends: []string{"linux-C", "edk2-C"}},
		},
		Linux: map[string]LinuxOpts{
			"linux-A": {Depends: []string{}},
			"linux-C": {Depends: []string{}},
		},
		Edk2: map[string]Edk2Opts{
			"edk2-B": {Depends: []string{}},
			"edk2-C": {Depends: []string{}},
		},
	}

	testConfigDependencyHell := Config{
		// Please keep everything in coreboot for simplicity sake
		// There is a test which checks order of builds, and it would explode in complexity
		Coreboot: map[string]CorebootOpts{
			"pizza":  {Depends: []string{"dough", "cheese"}},
			"dough":  {Depends: []string{"flour", "water"}},
			"cheese": {Depends: []string{"milk"}},
			"flour":  {Depends: []string{}},
			"water":  {Depends: []string{}},
			"milk":   {Depends: []string{"water"}},
		},
		Linux: map[string]LinuxOpts{},
		Edk2:  map[string]Edk2Opts{},
	}
	//  Pizza
	//    |
	//    +------+
	//    |      |
	//  dough  cheese
	//    |      |
	//    +---+  |
	//    |   |  |
	//  flour | milk
	//        |  |
	//       water

	testCases := []struct {
		name      string
		wantErr   error
		target    string
		recursive bool
		config    Config
	}{
		{
			name:      "unknown dependency",
			wantErr:   ErrDependencyTreeUndefDep,
			target:    "",
			recursive: false,
			config: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {Depends: []string{"dummy"}},
				},
			},
		},
		{
			name:      "circular self-dependency",
			wantErr:   ErrDependencyTreeUndefDep,
			target:    "",
			recursive: false,
			config: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {Depends: []string{"coreboot-A"}},
				},
			},
		},
		{
			name:      "circular dependency",
			wantErr:   ErrDependencyTreeUndefDep,
			target:    "",
			recursive: false,
			config: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {Depends: []string{"coreboot-B"}},
					"coreboot-B": {Depends: []string{"coreboot-A"}},
				},
			},
		},
		{
			name:      "unknown target",
			wantErr:   ErrDependencyTreeUnderTarget,
			target:    "",
			recursive: false,
			config:    testConfig,
		},
		{
			name:      "dependency clusterfuck",
			wantErr:   nil,
			target:    "pizza",
			recursive: false,
			config:    testConfigDependencyHell,
		},
		{
			name:      "dependency clusterfuck - middle",
			wantErr:   nil,
			target:    "milk",
			recursive: false,
			config:    testConfigDependencyHell,
		},
		{
			name:      "two leaves and one root",
			wantErr:   nil,
			target:    "stitch",
			recursive: false,
			config: Config{
				Edk2: map[string]Edk2Opts{
					"edk2-build-a": {Depends: []string{}},
					"edk2-build-b": {Depends: []string{}},
				},
				FirmwareStitching: map[string]FirmwareStitchingOpts{
					"stitch": {Depends: []string{"edk2-build-a"}},
				},
			},
		},
		{
			name:      "one root and two leaves",
			wantErr:   nil,
			target:    "stitch-a",
			recursive: false,
			config: Config{
				Edk2: map[string]Edk2Opts{
					"edk2-build": {Depends: []string{}},
				},
				FirmwareStitching: map[string]FirmwareStitchingOpts{
					"stitch-a": {Depends: []string{"edk2-build"}},
					"stitch-b": {Depends: []string{"edk2-build"}},
				},
			},
		},
	}

	const interactive = false
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Build(
				ctx,
				tc.target,
				tc.recursive,
				interactive,
				&tc.config,
				executeDummy,
			)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
	const recursive = true
	t.Run("recursive", func(t *testing.T) {
		builds, err := Build(
			ctx,
			"pizza",
			recursive,
			interactive,
			&testConfigDependencyHell,
			executeDummy,
		)
		assert.ErrorIs(t, err, nil)

		// Check for length
		assert.Equal(t, len(testConfigDependencyHell.Coreboot), len(builds))

		// Go though 'builds' and check if for each builds, the dependencies are already complete
		done := []string{}
		for _, item := range builds {
			for _, i := range testConfigDependencyHell.Coreboot[item.Name].Depends {
				assert.Contains(t, done, i)
			}
			done = append(done, item.Name)
		}
	})
}
