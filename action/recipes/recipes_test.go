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

func TestExecuteSkip(t *testing.T) {
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
	const outputDir = "output-coreboot"
	myConfig := Config{
		Coreboot: map[string]CorebootOpts{
			target: {
				CommonOpts: CommonOpts{
					OutputDir: outputDir,
					ContainerOutputFiles: []string{
						"build/coreboot.rom",
						"defconfig",
					},
				},
			},
		},
	}

	// Create the output directory
	err = os.Mkdir(outputDir, os.ModePerm)
	assert.NoError(t, err)

	// Since there is now existing output directory, it should skip the build
	err = Execute(ctx, target, &myConfig, interactive)
	assert.ErrorIs(t, err, ErrBuildSkipped)
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
			for _, i := range testConfigDependencyHell.Coreboot[item].Depends {
				assert.Contains(t, done, i)
			}
			done = append(done, item)
		}
	})
}
