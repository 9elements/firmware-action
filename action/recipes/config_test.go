// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	commonDummy := CommonOpts{
		SdkURL:        "https://ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		Arch:          "dummy",
		RepoPath:      "dummy/dir/",
		DefconfigPath: "dummy",
		OutputDir:     "dummy/dir/",
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
				Coreboot: []CorebootOpts{},
			},
		},
		{
			name: "missing required coreboot opts",
			opts: Config{
				Coreboot: []CorebootOpts{
					{
						ID: "coreboot-A",
					},
				},
			},
			wantErr: ErrRequiredOptionUndefined,
		},
		{
			name:    "required coreboot opts present",
			wantErr: nil,
			opts: Config{
				Coreboot: []CorebootOpts{
					{
						ID:     "coreboot-A",
						Common: commonDummy,
					},
				},
			},
		},
		{
			name:    "required coreboot opts present 2",
			wantErr: nil,
			opts: Config{
				Coreboot: []CorebootOpts{
					{
						ID:       "coreboot-A",
						Common:   commonDummy,
						Specific: CorebootSpecific{},
					},
				},
			},
		},
		{
			name:    "missing common opts",
			wantErr: ErrRequiredOptionUndefined,
			opts: Config{
				Coreboot: []CorebootOpts{
					{
						ID:       "coreboot-A",
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
		Coreboot: []CorebootOpts{},
		Linux:    []LinuxOpts{},
		Edk2:     []Edk2Opts{},
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
	equal := cmp.Equal(&configOriginal, configNew)
	if !equal {
		fmt.Println(cmp.Diff(&configOriginal, configNew))
		assert.True(t, equal, "written and read configuration are now equal")
	}
}
