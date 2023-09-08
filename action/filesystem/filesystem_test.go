// SPDX-License-Identifier: MIT
package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "textfile.txt")

	// Non-existing file
	assert.ErrorIs(t, CheckFileExists(path), os.ErrNotExist)

	// Existing file
	assert.NoError(t, os.WriteFile(path, []byte(""), 0o666))
	assert.ErrorIs(t, CheckFileExists(path), os.ErrExist)
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	pathSrc := filepath.Join(tmpDir, "textfile.txt")
	pathDest := filepath.Join(tmpDir, "textfile_new.txt")

	assert.NoError(t, os.WriteFile(pathSrc, []byte(""), 0o666))
	assert.ErrorIs(t, CheckFileExists(pathSrc), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDest), os.ErrNotExist)

	// Copy it
	assert.NoError(t, CopyFile(pathSrc, pathDest))

	assert.ErrorIs(t, CheckFileExists(pathSrc), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDest), os.ErrExist)

	// Copy it again (should fail)
	assert.ErrorIs(t, CopyFile(pathSrc, pathDest), os.ErrExist)

	assert.ErrorIs(t, CheckFileExists(pathSrc), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDest), os.ErrExist)
}

func TestMoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	pathSrc := filepath.Join(tmpDir, "textfile.txt")
	pathDest := filepath.Join(tmpDir, "textfile_new.txt")

	assert.NoError(t, os.WriteFile(pathSrc, []byte(""), 0o666))
	assert.ErrorIs(t, CheckFileExists(pathSrc), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDest), os.ErrNotExist)

	// Move it
	assert.NoError(t, MoveFile(pathSrc, pathDest))

	assert.ErrorIs(t, CheckFileExists(pathSrc), os.ErrNotExist)
	assert.ErrorIs(t, CheckFileExists(pathDest), os.ErrExist)
}

func TestDirTree(t *testing.T) {
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	files, err := DirTree(pwd)
	assert.NoError(t, err)
	assert.True(t, len(files) > 0, "found no files or directories")
}
