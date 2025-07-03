// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestChangeTimeStamp(t *testing.T) {
	tmpDir := t.TempDir()

	const (
		StatusDir = ".firmware-action"
		target    = "dummy"
	)

	timestampsDir := filepath.Join(tmpDir, StatusDir, "timestamps")
	repoPath := filepath.Join(tmpDir, "repo")
	assert.NoError(t, os.MkdirAll(repoPath, os.ModePerm))

	resultFile := filepath.Join(timestampsDir, filesystem.Filenamify(target, "txt"))

	myTimeStamp := ChangeTimeStamp{
		Change: Change{
			ResultFile: resultFile,
		},
		Sources: []string{repoPath},
	}

	// Also quick and dirty test for AllChanges
	myAllChanges := AllChanges{
		TimeStamp: myTimeStamp,
	}

	// No file, nothing
	assert.False(t, myTimeStamp.DetectChanges())

	// Save a checkpoint file
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrNotExist)
	myTimeStamp.SaveCheckpoint(false)
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrExist)
	assert.False(t, myTimeStamp.DetectChanges())
	// Also quick and dirty test for AllChanges
	assert.False(t, myAllChanges.DetectChanges(target))

	// make a new file to trigger change detection
	// we need a short sleep to actually detect the change
	time.Sleep(5 * time.Millisecond)
	assert.NoError(t, os.WriteFile(filepath.Join(repoPath, "file.rom"), []byte("test"), 0o666))
	assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(repoPath, "file.rom")), os.ErrExist)
	assert.True(t, myTimeStamp.DetectChanges())
	// Also quick and dirty test for AllChanges
	assert.True(t, myAllChanges.DetectChanges(target))
}

func TestChangeConfig(t *testing.T) {
	tmpDir := t.TempDir()

	const (
		StatusDir = ".firmware-action"
		target    = "dummy"
	)

	CompiledConfigsDir = filepath.Join(tmpDir, StatusDir, "configs")
	repoPath := filepath.Join(tmpDir, "repo/")
	assert.NoError(t, os.MkdirAll(repoPath, os.ModePerm))

	resultFile := filepath.Join(CompiledConfigsDir, filesystem.Filenamify(target, "json"))

	const (
		outputDir  = "output-universal/"
		depends    = "pre-dummy"
		outputDir2 = "output-universal2/"
	)

	config := Config{
		Universal: map[string]UniversalOpts{
			target: {
				Depends: []string{depends},
				CommonOpts: CommonOpts{
					SdkURL:            "whatever",
					RepoPath:          repoPath,
					OutputDir:         outputDir,
					ContainerInputDir: "inputs/",
					ContainerOutputFiles: []string{
						"file.rom",
					},
				},
				UniversalSpecific: UniversalSpecific{
					BuildCommands: []string{"false"},
				},
			},
			depends: {
				CommonOpts: CommonOpts{
					SdkURL:            "whatever",
					RepoPath:          repoPath,
					OutputDir:         outputDir2,
					ContainerInputDir: "inputs/",
					ContainerOutputFiles: []string{
						"file.rom",
					},
				},
				UniversalSpecific: UniversalSpecific{
					BuildCommands: []string{"false"},
				},
			},
		},
	}

	myChangeConfig := ChangeConfig{
		Change: Change{
			ResultFile: resultFile,
		},
		Config: &config,
	}

	// No file, nothing
	assert.False(t, myChangeConfig.DetectChanges(target))

	// Save a checkpoint file
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrNotExist)
	myChangeConfig.SaveCheckpoint(false)
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrExist)
	assert.False(t, myChangeConfig.DetectChanges(target))

	// make a change to the configuration
	opts := myChangeConfig.Config.Universal[target]
	opts.BuildCommands = []string{}
	myChangeConfig.Config.Universal[target] = opts
	assert.True(t, myChangeConfig.DetectChanges(target))
}

func gitRepoPrepare(t *testing.T, tmpDir string) {
	// Copied from git_test.go

	// Create empty git repository
	repo, err := git.PlainInit(tmpDir, false)
	assert.NoError(t, err)

	// Create a worktree to interact with the repo
	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	// Create a simple text file for testing
	filename := "README.md"
	pathFile := filepath.Join(tmpDir, filename)
	err = os.WriteFile(pathFile, []byte{}, 0o666)
	assert.NoError(t, err)

	// Commit the file (must be relative path)
	_, err = worktree.Add(filename)
	assert.NoError(t, err)

	commitOpts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "john doe",
			Email: "john.doe@example.com",
			When:  time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC),
		},
	}
	_, err = worktree.Commit("Initial commit", commitOpts)
	assert.NoError(t, err)
}

func gitRepoUpdateReadme(t *testing.T, tmpDir string) {
	// Create empty git repository
	repo, err := git.PlainOpen(tmpDir)
	assert.NoError(t, err)

	// Create a worktree to interact with the repo
	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	// Create a simple text file for testing
	filename := "README.md"
	pathFile := filepath.Join(tmpDir, filename)
	err = os.WriteFile(pathFile, []byte("hello there"), 0o666)
	assert.NoError(t, err)

	// Commit the file (must be relative path)
	_, err = worktree.Add(filename)
	assert.NoError(t, err)

	commitOpts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "john doe",
			Email: "john.doe@example.com",
			When:  time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC),
		},
	}
	_, err = worktree.Commit("Additional commit", commitOpts)
	assert.NoError(t, err)
}

func TestChangeGitHash(t *testing.T) {
	tmpDir := t.TempDir()

	const (
		StatusDir = ".firmware-action"
		target    = "dummy"
	)

	gitRepoHashDir := filepath.Join(tmpDir, StatusDir, "git-hashes")
	repoPath := filepath.Join(tmpDir, "repo")
	assert.NoError(t, os.MkdirAll(repoPath, os.ModePerm))

	resultFile := filepath.Join(gitRepoHashDir, filesystem.Filenamify(target, "txt"))

	myChangeGitHash := ChangeGitHash{
		Change: Change{
			ResultFile: resultFile,
		},
		RepoPath: repoPath,
	}

	// No file, nothing
	assert.False(t, myChangeGitHash.DetectChanges())

	// Try to save a checkpoint file, but should fail because RepoPath is not Git repo
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrNotExist)
	myChangeGitHash.SaveCheckpoint(target, false)
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrNotExist)

	// Make the RepoPath into Git repo
	gitRepoPrepare(t, repoPath)
	assert.ErrorIs(t, filesystem.CheckFileExists(filepath.Join(repoPath, ".git")), filesystem.ErrPathIsDirectory)

	// Save a checkpoint file
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrNotExist)
	myChangeGitHash.SaveCheckpoint(target, false)
	assert.ErrorIs(t, filesystem.CheckFileExists(resultFile), os.ErrExist)
	assert.False(t, myChangeGitHash.DetectChanges())

	// make a new commit
	gitRepoUpdateReadme(t, repoPath)
	assert.True(t, myChangeGitHash.DetectChanges())
}
