// SPDX-License-Identifier: MIT

//go:build go1.24

// Package recipes / linux
package recipes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestExtractSizeFromString(t *testing.T) {
	testCases := []struct {
		name     string
		stdout   string
		expected []uint64
		wantErr  error
	}{
		{
			name:     "empty string",
			stdout:   "",
			expected: []uint64{},
			wantErr:  errFailedToDetectRomSize,
		},
		{
			name:     "component 1; 2 unused",
			stdout:   "Component 2 Density:   UNUSED\nComponent 1 Density:   16MB",
			expected: []uint64{16 * 1024 * 1024, 0},
			wantErr:  nil,
		},
		{
			name:     "component 1, nothing about 2",
			stdout:   "Component 1 Density:   16MB",
			expected: []uint64{},
			wantErr:  errFailedToDetectRomSize,
		},
		{
			name:     "component 1 and 2",
			stdout:   "Component 2 Density:   4MB\nComponent 1 Density:   8MB",
			expected: []uint64{8 * 1024 * 1024, 4 * 1024 * 1024},
			wantErr:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractSizeFromString(tc.stdout)

			equal := cmp.Equal(tc.expected, result)
			if !equal {
				t.Log(cmp.Diff(tc.expected, result))
				assert.True(t, equal, "failed to extract size of ROM from string")
			}

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestStringToSizeMB(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected uint64
		wantErr  error
	}{
		{
			name:     "empty string",
			text:     "",
			expected: 0,
			wantErr:  errFailedToDetectRomSize,
		},
		{
			name:     "UNUSED",
			text:     "UNUSED",
			expected: 0,
			wantErr:  nil,
		},
		{
			name:     "unused",
			text:     "unused",
			expected: 0,
			wantErr:  nil,
		},
		{
			name:     "4MB",
			text:     "4MB",
			expected: 4 * 1024 * 1024,
			wantErr:  nil,
		},
		{
			name:     "64MB",
			text:     "64MB",
			expected: 64 * 1024 * 1024,
			wantErr:  nil,
		},
		{
			name:     "64MB with white space",
			text:     "  64MB  ",
			expected: 64 * 1024 * 1024,
			wantErr:  nil,
		},
		{
			name:     "64MB with whitespace and newlines",
			text:     "\n  64MB  \n",
			expected: 64 * 1024 * 1024,
			wantErr:  nil,
		},
		{
			name:     "bogus string",
			text:     "bogus string",
			expected: 0,
			wantErr:  errFailedToDetectRomSize,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToSizeMB(tc.text)

			equal := cmp.Equal(tc.expected, result)
			if !equal {
				t.Log(cmp.Diff(tc.expected, result))
				assert.True(t, equal, "failed to decipher size")
			}

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

type makeFile struct {
	Path       string
	Content    string
	SourcePath string
}

func (base makeFile) MakeMe() error {
	log.Printf("Path:       %s", base.Path)
	log.Printf("Content:    %s", base.Content)
	log.Printf("SourcePath: %s", base.SourcePath)

	pwd, _ := os.Getwd()
	log.Printf("PWD:        %s", pwd)

	// If file does not exist, make it
	if _, err := os.Stat(base.Path); os.IsNotExist(err) {
		if base.Content != "" {
			// Create file with content
			file, err := os.Create(base.Path)
			if err != nil {
				return err
			}

			_, err = file.Write([]byte(base.Content))
			if err != nil {
				return err
			}
		} else {
			// Copy file from somewhere
			if _, err := os.Stat(base.SourcePath); os.IsNotExist(err) {
				log.Printf("[Mock MakeMe] file '%s' does not exists", base.SourcePath)
				return os.ErrNotExist
			}

			log.Print("Success")

			return filesystem.CopyFile(base.SourcePath, base.Path)
		}
	}

	log.Print("Success")

	return nil
}

func TestStitching(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Define common variables and values
	baseFileName := "base.img"
	common := CommonOpts{
		SdkURL:    "ghcr.io/9elements/firmware-action/coreboot_4.19:main",
		OutputDir: "output",
		ContainerOutputFiles: []string{
			fmt.Sprintf("new_%s", baseFileName),
		},
	}

	// Download blobs (contains example IFD.bin)
	blobDir := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__/blobs")
	if _, err := os.Stat(blobDir); os.IsNotExist(err) {
		err := exec.Command("git", "clone", "https://review.coreboot.org/blobs.git", blobDir).Run()
		assert.NoError(t, err)
	}

	testCases := []struct {
		name           string
		stitchingOpts  FirmwareStitchingOpts
		files          []makeFile
		expectedSha256 string
		wantErr        error
	}{
		{
			name: "real test - inject ME into IFD",
			stitchingOpts: FirmwareStitchingOpts{
				CommonOpts:   common,
				BaseFilePath: baseFileName,
				IfdtoolEntries: []IfdtoolEntry{
					{
						Path:         "me.bin",
						TargetRegion: "ME",
					},
				},
			},
			files: []makeFile{
				{
					Path:       baseFileName,
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/descriptor.bin"),
				},
				{
					Path:       "me.bin",
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/me.bin"),
				},
			},
			expectedSha256: "a09cf57dae3062b18ae84f6695a22c5e1e61e3a84a9c9de69af40a0e54b658d4",
			// this magic value was obtained by doing all steps manually
			wantErr: nil,
		},
		// Test behavior of IgnoreIfMissing on IfdtoolEntry
		{
			name: "   file and no ignore",
			stitchingOpts: FirmwareStitchingOpts{
				CommonOpts:   common,
				BaseFilePath: baseFileName,
				IfdtoolEntries: []IfdtoolEntry{
					{
						Path:         "me.bin",
						TargetRegion: "ME",
					},
				},
			},

			files: []makeFile{
				{
					Path:       baseFileName,
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/descriptor.bin"),
				},
				{
					Path:       "me.bin",
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/me.bin"),
				},
			},
			expectedSha256: "a09cf57dae3062b18ae84f6695a22c5e1e61e3a84a9c9de69af40a0e54b658d4",
			// Should stitch the files without complains
			// this magic value was obtained by doing all steps manually
			wantErr: nil,
		},
		{
			name: "   file and    ignore",
			stitchingOpts: FirmwareStitchingOpts{
				CommonOpts:   common,
				BaseFilePath: baseFileName,
				IfdtoolEntries: []IfdtoolEntry{
					{
						Path:            "me.bin",
						TargetRegion:    "ME",
						IgnoreIfMissing: true,
					},
				},
			},
			files: []makeFile{
				{
					Path:       baseFileName,
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/descriptor.bin"),
				},
				{
					Path:       "me.bin",
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/me.bin"),
				},
			},
			expectedSha256: "a09cf57dae3062b18ae84f6695a22c5e1e61e3a84a9c9de69af40a0e54b658d4",
			// Should stitch the files without complains
			// this magic value was obtained by doing all steps manually
			wantErr: nil,
		},
		{
			name: "no file and no ignore",
			stitchingOpts: FirmwareStitchingOpts{
				CommonOpts:   common,
				BaseFilePath: baseFileName,
				IfdtoolEntries: []IfdtoolEntry{
					{
						Path:         "me.bin",
						TargetRegion: "ME",
					},
				},
			},
			files: []makeFile{
				{
					Path:       baseFileName,
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/descriptor.bin"),
				},
			},
			// Should complain about missing file
			wantErr: os.ErrNotExist,
		},
		{
			name: "no file and    ignore",
			stitchingOpts: FirmwareStitchingOpts{
				CommonOpts:   common,
				BaseFilePath: baseFileName,
				IfdtoolEntries: []IfdtoolEntry{
					{
						Path:            "me.bin",
						TargetRegion:    "ME",
						IgnoreIfMissing: true,
					},
				},
			},
			files: []makeFile{
				{
					Path:       baseFileName,
					SourcePath: filepath.Join(blobDir, "mainboard/intel/emeraldlake2/descriptor.bin"),
				},
			},
			expectedSha256: "5c8283b8c668e6735afe3b4209dce64924f7d5f5da771d51b05d37d52dc48331",
			// The image file should only inflate to 16MB (filled with 0xFF) but not change in content
			// this magic value was obtained by running sha256sum on original file
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prep
			//   vars
			tmpDir := t.TempDir()
			tc.stitchingOpts.RepoPath = filepath.Join(tmpDir, "stitch")
			//   make repo dir
			err := os.Mkdir(tc.stitchingOpts.RepoPath, os.ModePerm)
			assert.NoError(t, err)

			outputPath := filepath.Join(tmpDir, tc.stitchingOpts.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)

			// Change current working directory
			t.Chdir(tmpDir)

			// Move files
			for i := range tc.files {
				err = tc.files[i].MakeMe()
				assert.NoError(t, err)
			}

			// Stitch
			ctx := t.Context()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)

			defer client.Close()

			err = tc.stitchingOpts.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			if tc.wantErr != nil {
				return
			}

			// Check artifacts
			finalImageFile := filepath.Join(
				outputPath,
				tc.stitchingOpts.ContainerOutputFiles[0],
			)
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(finalImageFile), os.ErrExist)
			}

			// Compare
			newContent, err := os.ReadFile(finalImageFile)
			assert.NoError(t, err)

			hash := sha256.New()
			hash.Write(newContent)
			hashHex := hex.EncodeToString(hash.Sum(nil))
			// TODO: fix expected vs actual values in all of these
			assert.Equal(t, tc.expectedSha256, hashHex)
		})
	}
}
