// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestAllModules(t *testing.T) {
	testCases := []struct {
		name        string
		opts        Config
		wantModules map[string]FirmwareModule
	}{
		{
			name:        "empty",
			opts:        Config{},
			wantModules: map[string]FirmwareModule{},
		},
		{
			name: "simple",
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {},
				},
			},
			wantModules: map[string]FirmwareModule{
				"coreboot-A": CorebootOpts{},
			},
		},
		{
			name: "more complex",
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						DefconfigPath: "dummy",
					},
				},
			},
			wantModules: map[string]FirmwareModule{
				"coreboot-A": CorebootOpts{
					DefconfigPath: "dummy",
				},
				"linux-A": LinuxOpts{
					DefconfigPath: "dummy",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allModules := tc.opts.AllModules()
			assert.True(t, reflect.DeepEqual(tc.wantModules, allModules))
		})
	}
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		name       string
		optsA      Config
		optsB      Config
		wantConfig Config
		wantErr    error
	}{
		{
			name:       "empty",
			optsA:      Config{},
			optsB:      Config{},
			wantConfig: Config{},
			wantErr:    nil,
		},
		{
			name: "simple",
			optsA: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {},
				},
			},
			optsB: Config{},
			wantConfig: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {},
				},
			},
			wantErr: nil,
		},
		{
			name: "more complex",
			optsA: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						DefconfigPath: "dummy",
					},
				},
			},
			optsB: Config{
				Linux: map[string]LinuxOpts{
					"linux-A": {
						DefconfigPath: "dummy",
					},
				},
			},
			wantConfig: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						DefconfigPath: "dummy",
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bigConfig, err := tc.optsA.Merge(tc.optsB)

			// Initialize empty maps in wantConfig so that deep equal works correctly
			wantConfigValue := reflect.ValueOf(&tc.wantConfig).Elem()
			for i := 0; i < wantConfigValue.NumField(); i++ {
				fieldValue := wantConfigValue.Field(i)
				if fieldValue.Kind() == reflect.Map && fieldValue.IsNil() {
					fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
				}
			}

			// Compare the merged config with the expected one
			assert.Equal(t, tc.wantConfig, bigConfig)
			assert.True(t, reflect.DeepEqual(tc.wantConfig, bigConfig))
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:            "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		RepoPath:          "dummy/dir/",
		OutputDir:         "dummy/dir/",
		ContainerInputDir: "inputs/",
	}

	testCases := []struct {
		name    string
		wantErr error
		opts    Config
	}{
		{
			name:    "completely empty",
			wantErr: nil,
			opts:    Config{},
		},
		{
			name:    "empty coreboot opts",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{},
			},
		},
		{
			name:    "missing required coreboot opts",
			wantErr: ErrFailedValidation,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {},
				},
			},
		},
		{
			name:    "required coreboot opts present",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts:    commonDummy,
						DefconfigPath: "dummy",
					},
				},
			},
		},
		{
			name:    "required coreboot opts present 2",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts:    commonDummy,
						DefconfigPath: "dummy",
						Blobs:         map[string]string{},
					},
				},
			},
		},
		{
			name:    "missing common opts",
			wantErr: ErrFailedValidation,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts:    CommonOpts{},
						DefconfigPath: "dummy",
						Blobs:         map[string]string{},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.opts)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestValidateConfigNestedOutputDirs(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:            "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		RepoPath:          "dummy/dir/",
		ContainerInputDir: "inputs/",
	}

	testCases := []struct {
		name    string
		wantErr error
		opts    Config
	}{
		{
			name:    "valid - different output dirs",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output-coreboot/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output-linux/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
			},
		},
		{
			name:    "valid - sibling output dirs",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/coreboot/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/linux/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
			},
		},
		{
			name:    "invalid - nested output dirs",
			wantErr: ErrNestedOutputDirs,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/linux/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
			},
		},
		{
			name:    "invalid - duplicate output dirs",
			wantErr: ErrDuplicateOutputDirs,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
			},
		},
		{
			name:    "valid - similar but not nested dirs",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output-dir/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
				Linux: map[string]LinuxOpts{
					"linux-A": {
						CommonOpts: CommonOpts{
							SdkURL:            commonDummy.SdkURL,
							RepoPath:          commonDummy.RepoPath,
							OutputDir:         "output-dir-extra/",
							ContainerInputDir: commonDummy.ContainerInputDir,
						},
						DefconfigPath: "dummy",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.opts)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestConfigReadAndWrite(t *testing.T) {
	configOriginal := Config{
		Coreboot: map[string]CorebootOpts{},
		Linux:    map[string]LinuxOpts{},
		Edk2:     map[string]Edk2Opts{},
	}

	tmpDir := t.TempDir()
	configFilepath := filepath.Join(tmpDir, "test.json")

	// Write
	err := WriteConfig(configFilepath, &configOriginal)
	assert.NoError(t, err)

	// Read
	configNew, err := ReadConfig(configFilepath)
	assert.NoError(t, err)

	// Compare
	equal := cmp.Equal(&configOriginal, configNew)
	if !equal {
		t.Log(cmp.Diff(configOriginal, configNew))
		assert.True(t, equal, "written and read configuration are not equal")
	}
}

func TestFindAllEnvVars(t *testing.T) {
	testCases := []struct {
		name            string
		text            string
		expectedEnvVars []string
	}{
		{
			name:            "no env vars",
			text:            "dummy string",
			expectedEnvVars: []string{},
		},
		{
			name:            "one env var",
			text:            "dummy string with $MY_VAR",
			expectedEnvVars: []string{"MY_VAR"},
		},
		{
			name:            "one env var with brackets",
			text:            "dummy string with ${MY_VAR}",
			expectedEnvVars: []string{"MY_VAR"},
		},
		{
			name:            "two env vars",
			text:            "dummy string with $MY_VAR and ${MY_VAR}",
			expectedEnvVars: []string{"MY_VAR", "MY_VAR"},
		},
		{
			name:            "two env vars with numbers",
			text:            "dummy string with $MY_VAR1 and ${MY_VAR2}",
			expectedEnvVars: []string{"MY_VAR1", "MY_VAR2"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundVars := FindAllEnvVars(tc.text)
			t.Log(foundVars)
			t.Log(tc.expectedEnvVars)

			assert.Len(t, foundVars, len(tc.expectedEnvVars))
			// If both slices are of zero length, then the comparison fails for whatever reason
			if len(tc.expectedEnvVars) > 0 {
				assert.True(t, reflect.DeepEqual(tc.expectedEnvVars, foundVars))
			}
		})
	}
}

func TestConfigEnvVars(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:   "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		RepoPath: "dummy/dir/",
	}

	testCases := []struct {
		name             string
		wantErr          error
		url              string
		urlExpected      string
		repoPath         string
		repoPathExpected string
		envVars          map[string]string
	}{
		{
			name:             "no env vars",
			wantErr:          nil,
			url:              commonDummy.SdkURL,
			urlExpected:      commonDummy.SdkURL,
			repoPath:         commonDummy.RepoPath,
			repoPathExpected: commonDummy.RepoPath,
			envVars:          map[string]string{},
		},
		{
			name:             "env vars",
			wantErr:          nil,
			url:              "ghcr.io/$TEST_ENV_VAR/coreboot_4.19:main",
			urlExpected:      commonDummy.SdkURL,
			repoPath:         commonDummy.RepoPath,
			repoPathExpected: commonDummy.RepoPath,
			envVars: map[string]string{
				"TEST_ENV_VAR": "9elements/firmware-action",
			},
		},
		{
			name:             "env vars with brackets",
			wantErr:          nil,
			url:              "ghcr.io/${TEST_ENV_VAR}/coreboot_4.19:main",
			urlExpected:      commonDummy.SdkURL,
			repoPath:         commonDummy.RepoPath,
			repoPathExpected: commonDummy.RepoPath,
			envVars: map[string]string{
				"TEST_ENV_VAR": "9elements/firmware-action",
			},
		},
		{
			name:             "multiple env vars in multiple entries",
			wantErr:          nil,
			url:              "ghcr.io/${TEST_ENV_VAR_PROJECT}/${TEST_ENV_VAR_SDK}:${TEST_ENV_VAR_VERSION}",
			urlExpected:      commonDummy.SdkURL,
			repoPath:         "${TEST_ENV_VAR_REPOPATH}",
			repoPathExpected: commonDummy.RepoPath,
			envVars: map[string]string{
				"TEST_ENV_VAR_PROJECT":  "9elements/firmware-action",
				"TEST_ENV_VAR_VERSION":  "main",
				"TEST_ENV_VAR_SDK":      "coreboot_4.19",
				"TEST_ENV_VAR_REPOPATH": "dummy/dir/",
			},
		},
		{
			name:             "undefined env var",
			wantErr:          ErrEnvVarUndefined,
			url:              "ghcr.io/$TEST_ENV_VAR/coreboot_4.19:main",
			urlExpected:      commonDummy.SdkURL,
			repoPath:         commonDummy.RepoPath,
			repoPathExpected: commonDummy.RepoPath,
			envVars:          map[string]string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						CommonOpts: CommonOpts{
							SdkURL:            tc.url,
							RepoPath:          "dummy/dir/",
							OutputDir:         "dummy/dir/",
							ContainerInputDir: "inputs/",
						},
						DefconfigPath: "dummy",
					},
				},
			}
			for key, value := range tc.envVars {
				t.Setenv(key, value)
				t.Logf("Setting %s = %s\n", key, value)
			}

			// Write and read config
			// The read function handles the expansion of environment variables
			tmpDir := t.TempDir()
			configFilepath := filepath.Join(tmpDir, "test.json")
			// Write
			err := WriteConfig(configFilepath, &opts)
			assert.NoError(t, err)
			// Read
			optsConverted, err := ReadConfig(configFilepath)

			// err = ValidateConfig(optsConverted)
			assert.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr == nil {
				assert.Equal(t, tc.urlExpected, optsConverted.Coreboot["coreboot-A"].SdkURL)
				assert.Equal(t, tc.repoPathExpected, optsConverted.Coreboot["coreboot-A"].RepoPath)
			}
		})
	}
}

func TestOffsetToLineNumber(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		offset    int
		line      int
		character int
		wantErr   error
	}{
		{
			name:    "empty string",
			input:   "",
			offset:  1,
			wantErr: ErrVerboseJSON,
		},
		{
			name:      "1 line, offset 0",
			input:     "dummy line",
			offset:    0,
			line:      1,
			character: 1,
			wantErr:   nil,
		},
		{
			name:      "1 line, offset 1",
			input:     "dummy line",
			offset:    1,
			line:      1,
			character: 2,
			wantErr:   nil,
		},
		{
			name: "2 lines, offset in 1st line",
			// offset:  012345 6789
			input:     "dummy\nline",
			offset:    1,
			line:      1,
			character: 2,
			wantErr:   nil,
		},
		{
			name: "2 lines, offset end of 1st line",
			// offset:  012345 6789
			input:     "dummy\nline",
			offset:    4,
			line:      1,
			character: 5,
			wantErr:   nil,
		},
		{
			name:      "2 lines, offset in 2nd line",
			input:     "dummy\nline",
			offset:    7,
			line:      2,
			character: 2,
			wantErr:   nil,
		},
		{
			name: "2 lines, offset in 2nd line",
			// offset:  0 1 2 3 4567
			input:     "\n\n\n\nline",
			offset:    7,
			line:      5,
			character: 4,
			wantErr:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line, character, err := offsetToLineNumber(tc.input, tc.offset)
			assert.ErrorIs(t, err, tc.wantErr)

			if err != nil {
				// no need to continue on error
				return
			}

			assert.Equal(t, tc.line, line, "line is wrong")
			assert.Equal(t, tc.character, character, "character is wrong")
		})
	}
}
