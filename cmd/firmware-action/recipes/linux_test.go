// SPDX-License-Identifier: MIT

//go:build go1.24

// Package recipes / linux
package recipes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

func TestLinux(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	pwd, err := os.Getwd()
	assert.NoError(t, err)

	linuxOpts := LinuxOpts{
		CommonOpts: CommonOpts{
			OutputDir: "output",
			ContainerOutputFiles: []string{
				"vmlinux",
				"defconfig",
			},
		},
		DefconfigPath: "custom_defconfig",
	}

	testCases := []struct {
		name         string
		linuxVersion string
		arch         string
		wantErr      error
	}{
		{
			name:         "normal build for x86 64bit",
			linuxVersion: "6.1.127",
			arch:         "amd64",
			wantErr:      nil,
		},
		{
			name:         "normal build for x86 32bit",
			linuxVersion: "6.1.127",
			arch:         "i386",
			wantErr:      nil,
		},
		{
			name:         "normal build for arm64",
			linuxVersion: "6.1.127",
			arch:         "arm64",
			wantErr:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			linuxVersion, err := semver.NewVersion(tc.linuxVersion)
			assert.NoError(t, err)
			ctx := t.Context()
			client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
			assert.NoError(t, err)
			defer client.Close()

			// Prepare options
			tmpDir := t.TempDir()

			myLinuxOpts := linuxOpts
			myLinuxOpts.SdkURL = fmt.Sprintf("ghcr.io/9elements/firmware-action/linux_%d.%d:main", linuxVersion.Major(), linuxVersion.Minor())
			myLinuxOpts.Arch = tc.arch
			myLinuxOpts.RepoPath = filepath.Join(tmpDir, "linux")

			// Create __tmp_files__ directory to store source-code of Linux Kernel
			// mostly useful for repeated local-run tests to save bandwidth and time
			tmpFiles := filepath.Join(os.TempDir(), "__firmware-action_tmp_files__")
			err = os.MkdirAll(tmpFiles, 0o750)
			assert.NoError(t, err)
			t.Chdir(tmpFiles)

			// Download linux source code to __tmp_files__
			var commands [][]string
			// TODO: make these commands OS independent
			if errors.Is(filesystem.CheckFileExists(fmt.Sprintf("linux-%s", tc.linuxVersion)), os.ErrNotExist) {
				commands = [][]string{
					// Get Linux Kernel sources
					{"wget", "--quiet", fmt.Sprintf("https://cdn.kernel.org/pub/linux/kernel/v%d.x/linux-%s.tar.xz", linuxVersion.Major(), tc.linuxVersion)},
					{"wget", "--quiet", fmt.Sprintf("https://cdn.kernel.org/pub/linux/kernel/v%d.x/linux-%s.tar.sign", linuxVersion.Major(), tc.linuxVersion)},
					// un-xz
					{"unxz", "--keep", fmt.Sprintf("linux-%s.tar.xz", tc.linuxVersion)},
					// GPG verify
					{"gpg2", "--locate-keys", "torvalds@kernel.org", "gregkh@kernel.org"},
					{"gpg2", "--verify", fmt.Sprintf("linux-%s.tar.sign", tc.linuxVersion)},
					// un-tar
					{"tar", "-xvf", fmt.Sprintf("linux-%s.tar", tc.linuxVersion)},
				}
			}
			//   always copy from __tmp_files__ to tmpDir for each test
			commands = append(commands, []string{"cp", "-r", fmt.Sprintf("linux-%s", tc.linuxVersion), myLinuxOpts.RepoPath})
			for _, cmd := range commands {
				err = exec.Command(cmd[0], cmd[1:]...).Run()
				assert.NoError(t, err)
			}
			t.Chdir(myLinuxOpts.RepoPath)

			// Copy over defconfig file into tmpDir/linux
			defconfigPath := filepath.Join(myLinuxOpts.RepoPath, myLinuxOpts.DefconfigPath)
			repoRootPath, err := filepath.Abs(filepath.Join(pwd, "../../.."))
			assert.NoError(t, err)
			//   common.RepoPath = path to end user repository (in this case somewhere in /tmp)
			//   repoRootPath    = path to our repository with this code (contains configuration files for testing)
			defconfigLocalPath, err := filepath.Abs(filepath.Join(
				repoRootPath,
				fmt.Sprintf("tests/linux_%d.%d/linux.defconfig", linuxVersion.Major(), linuxVersion.Minor()),
			))
			assert.NoErrorf(t, err, "encountered issue with missing files, is '%s' the root of the repo?", repoRootPath)
			err = filesystem.CopyFile(
				defconfigLocalPath,
				defconfigPath,
			)
			assert.NoError(t, err)

			// Artifacts
			outputPath := filepath.Join(tmpDir, myLinuxOpts.OutputDir)
			err = os.MkdirAll(outputPath, os.ModePerm)
			assert.NoError(t, err)
			myLinuxOpts.OutputDir = outputPath

			// Try to build linux kernel
			err = myLinuxOpts.buildFirmware(ctx, client)
			assert.ErrorIs(t, err, tc.wantErr)

			// Check artifacts
			if tc.wantErr == nil {
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "vmlinux")), os.ErrExist)
				assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
			}
		})
	}
}
