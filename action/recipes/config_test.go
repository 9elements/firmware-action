// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:    "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		Arch:      "dummy",
		RepoPath:  "dummy/dir/",
		OutputDir: "dummy/dir/",
	}
	specificDummy := CorebootSpecific{
		DefconfigPath: "dummy",
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
			name: "missing required coreboot opts",
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {},
				},
			},
			wantErr: ErrRequiredOptionUndefined,
		},
		{
			name:    "required coreboot opts present",
			wantErr: nil,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						Common:   commonDummy,
						Specific: specificDummy,
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
						Common:   commonDummy,
						Specific: specificDummy,
					},
				},
			},
		},
		{
			name:    "missing common opts",
			wantErr: ErrRequiredOptionUndefined,
			opts: Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						Common:   CommonOpts{},
						Specific: CorebootSpecific{},
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
	err := WriteConfig(configFilepath, configOriginal)
	assert.NoError(t, err)

	// Read
	configNew, err := ReadConfig(configFilepath)
	assert.NoError(t, err)

	// Compare
	equal := cmp.Equal(configOriginal, configNew)
	if !equal {
		fmt.Println(cmp.Diff(configOriginal, configNew))
		assert.True(t, equal, "written and read configuration are not equal")
	}
}

func TestConfigEnvVars(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:    "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		Arch:      "dummy",
		RepoPath:  "dummy/dir/",
		OutputDir: "dummy/dir/",
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
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := Config{
				Coreboot: map[string]CorebootOpts{
					"coreboot-A": {
						Common: CommonOpts{
							SdkURL:    tc.url,
							Arch:      "dummy",
							RepoPath:  "dummy/dir/",
							OutputDir: "dummy/dir/",
						},
						Specific: CorebootSpecific{
							DefconfigPath: "dummy",
						},
					},
				},
			}
			for key, value := range tc.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
				fmt.Printf("Setting %s = %s\n", key, value)
			}

			// Write and read config
			// The read function handles the expansion of environment variables
			tmpDir := t.TempDir()
			configFilepath := filepath.Join(tmpDir, "test.json")
			// Write
			err := WriteConfig(configFilepath, opts)
			assert.NoError(t, err)
			// Read
			optsConverted, err := ReadConfig(configFilepath)
			assert.NoError(t, err)

			// err = ValidateConfig(optsConverted)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.urlExpected, optsConverted.Coreboot["coreboot-A"].Common.SdkURL)
			assert.Equal(t, tc.repoPathExpected, optsConverted.Coreboot["coreboot-A"].Common.RepoPath)
		})
	}
}
