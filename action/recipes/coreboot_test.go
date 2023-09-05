// SPDX-License-Identifier: MIT

// Package recipes / coreboot
package recipes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/container"
	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestCoreboot(t *testing.T) {
	// This test is really slow (like 100 seconds)
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	const corebootVersion = "4.19"
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	// Prepare options
	tmpDir := t.TempDir()
	opts := map[string]string{
		"target":           "coreboot",
		"sdk_version":      fmt.Sprintf("coreboot_%s:main", corebootVersion),
		"architecture":     "x86",
		"repo_path":        filepath.Join(tmpDir, "coreboot"),
		"defconfig_path":   "defconfig",
		"containerWorkDir": "/coreboot",
		"GITHUB_WORKSPACE": "/coreboot",
		"output":           "output",
	}
	getFunc := func(key string) string {
		return opts[key]
	}
	common, err := commonGetOpts(getFunc)
	assert.NoError(t, err)
	corebootOpts := corebootOpts{}

	// Change current working directory
	pwd, err := os.Getwd()
	defer os.Chdir(pwd) // nolint:errcheck
	assert.NoError(t, err)
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	// Clone coreboot repo
	cmd := exec.Command("git", "clone", "--branch", corebootVersion, "--depth", "1", "https://review.coreboot.org/coreboot")
	err = cmd.Run()
	assert.NoError(t, err)
	err = os.Chdir(common.repoPath)
	assert.NoError(t, err)

	// Copy over defconfig file into tmpDir
	defconfigPath := filepath.Join(common.repoPath, "defconfig")
	err = filesystem.CopyFile(
		filepath.Join(pwd, fmt.Sprintf("../../tests/coreboot_%s/seabios.defconfig", corebootVersion)),
		defconfigPath,
	)
	//   ^^^ this relative path might be funky
	assert.NoError(t, err)

	// Artifacts
	outputPath := filepath.Join(tmpDir, common.outputDir)
	err = os.MkdirAll(outputPath, os.ModePerm)
	assert.NoError(t, err)
	artifacts := []container.Artifacts{
		{
			ContainerPath: filepath.Join(common.containerWorkDir, "build", "coreboot.rom"),
			ContainerDir:  false,
			HostPath:      outputPath,
		},
		{
			ContainerPath: filepath.Join(common.containerWorkDir, "defconfig"),
			ContainerDir:  false,
			HostPath:      outputPath,
		},
	}

	// Try to build coreboot
	err = coreboot(ctx, client, &common, "", &corebootOpts, &artifacts)
	assert.NoError(t, err)

	// Check artifacts
	assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "coreboot.rom")), os.ErrExist)
	assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(outputPath, "defconfig")), os.ErrExist)
}
