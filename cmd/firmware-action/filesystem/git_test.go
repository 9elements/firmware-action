// SPDX-License-Identifier: MIT

//go:build go1.24

package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func gitRepoPrepare(t *testing.T, tmpDir string) {
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

func TestGitRun(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	// Test git status ordinary directory (not git repo)
	_, err := gitRun("./", []string{"git", "status"})
	assert.ErrorIs(t, err, ErrNotGitRepository)

	gitRepoPrepare(t, tmpDir)

	// Test git status
	stdout, err := gitRun("./", []string{"git", "status"})
	assert.NoError(t, err)
	assert.Equal(t, "On branch master\nnothing to commit, working tree clean\n", stdout)
}

func TestGitDescribeCoreboot(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	gitRepoPrepare(t, tmpDir)

	// Test git status
	describe, err := GitDescribeCoreboot("./")
	assert.NoError(t, err)
	assert.Equal(t, "4eeb1eaf0c81", describe)
	// This magic value comes from manual execution of the test
	// Since the content, author and time of the commit are hard-coded,
	//   the commit hash is always the same
}
