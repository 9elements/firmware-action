// SPDX-License-Identifier: MIT
package filesystem

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

	// Non-existing directory
	path = filepath.Join(tmpDir, "directory")
	assert.ErrorIs(t, CheckFileExists(path), os.ErrNotExist)

	// Existing directory
	assert.NoError(t, os.MkdirAll(path, 0o775))
	assert.ErrorIs(t, CheckFileExists(path), ErrPathIsDirectory)
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

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()
	pathSrc := filepath.Join(tmpDir, "test_dir")
	pathSrcFile := filepath.Join(pathSrc, "textfile.txt")
	pathDest := filepath.Join(tmpDir, "test_dir_new")
	pathDestFile := filepath.Join(pathDest, "textfile.txt")

	assert.NoError(t, os.MkdirAll(pathSrc, 0o775))
	assert.NoError(t, os.WriteFile(pathSrcFile, []byte(""), 0o666))
	assert.ErrorIs(t, CheckFileExists(pathSrcFile), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDestFile), os.ErrNotExist)

	// Copy it
	assert.NoError(t, CopyDir(pathSrc, pathDest))

	assert.ErrorIs(t, CheckFileExists(pathSrcFile), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDestFile), os.ErrExist)

	// Copy it again (should fail)
	assert.ErrorIs(t, CopyDir(pathSrc, pathDest), os.ErrExist)

	assert.ErrorIs(t, CheckFileExists(pathSrcFile), os.ErrExist)
	assert.ErrorIs(t, CheckFileExists(pathDestFile), os.ErrExist)
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

func TestLastSaveRunTime(t *testing.T) {
	currentTime := time.Now()

	tmpDir := t.TempDir()
	pathTimeFile := filepath.Join(tmpDir, "last_run_time.txt")

	// Load - should fallback because no file exists, but no error
	loadTime, err := LoadLastRunTime(pathTimeFile)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, time.Time{}, loadTime)
	assert.ErrorIs(t, CheckFileExists(pathTimeFile), os.ErrNotExist)

	// Save
	err = SaveCurrentRunTime(pathTimeFile)
	assert.NoError(t, err)
	// file should now exist
	assert.ErrorIs(t, CheckFileExists(pathTimeFile), os.ErrExist)

	// Load again - should now work since file exists
	loadTime, err = LoadLastRunTime(pathTimeFile)
	assert.NoError(t, err)
	assert.True(t, loadTime.After(currentTime))
	assert.True(t, time.Now().After(loadTime))
}

func TestGetFileModTime(t *testing.T) {
	tmpDir := t.TempDir()
	pathFile := filepath.Join(tmpDir, "test.txt")

	// Missing file - should fail
	modTime, err := GetFileModTime(pathFile)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, time.Time{}, modTime)
	assert.ErrorIs(t, CheckFileExists(pathFile), os.ErrNotExist)

	// Make file
	err = os.WriteFile(pathFile, []byte{}, 0o666)
	assert.NoError(t, err)
	assert.ErrorIs(t, CheckFileExists(pathFile), os.ErrExist)

	// Should work
	_, err = GetFileModTime(pathFile)
	assert.NoError(t, err)
}

func TestAnyFileNewerThan(t *testing.T) {
	tmpDir := t.TempDir()
	pathFile := filepath.Join(tmpDir, "test.txt")

	// Call on missing file - should fail
	mod, err := AnyFileNewerThan(pathFile, time.Now())
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.False(t, mod)

	// Call on existing file
	// - Make file
	err = os.WriteFile(pathFile, []byte{}, 0o666)
	assert.NoError(t, err)
	assert.ErrorIs(t, CheckFileExists(pathFile), os.ErrExist)
	// - Should work - is file newer than last year? (true)
	mod, err = AnyFileNewerThan(pathFile, time.Now().AddDate(-1, 0, 0))
	assert.NoError(t, err)
	assert.True(t, mod)
	// - Should work - is file newer than next year? (false)
	mod, err = AnyFileNewerThan(pathFile, time.Now().AddDate(1, 0, 0))
	assert.NoError(t, err)
	assert.False(t, mod)

	// Call on nested directory
	// - Make directory tree
	subDirRoot := filepath.Join(tmpDir, "test")
	subSubDir := filepath.Join(subDirRoot, "deep_test/even_deeper")
	err = os.MkdirAll(subSubDir, os.ModePerm)
	assert.NoError(t, err)
	// - Make file
	deepFile := filepath.Join(subSubDir, "test.txt")
	err = os.WriteFile(deepFile, []byte{}, 0o666)
	assert.NoError(t, err)
	assert.ErrorIs(t, CheckFileExists(deepFile), os.ErrExist)
	// - Should work - older
	mod, err = AnyFileNewerThan(subDirRoot, time.Now().AddDate(-1, 0, 0))
	assert.NoError(t, err)
	assert.True(t, mod)
	// - Should work - newer
	mod, err = AnyFileNewerThan(subDirRoot, time.Now().AddDate(1, 0, 0))
	assert.NoError(t, err)
	assert.False(t, mod)
}

func TestFilenamify(t *testing.T) {
	testCases := []struct {
		name           string
		inputName      string
		inputExtension string
		output         string
	}{
		{
			name:           "empty strings",
			inputName:      "",
			inputExtension: "",
			output:         ".",
		},
		{
			name:           "unicode control characters",
			inputName:      "foo\u0000bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo<bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo>bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo:bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo\"bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo/bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo\\bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo\bar",
			inputExtension: "",
			output:         "foo_ar.",
		},
		{
			inputName:      "foo|bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo?bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo*bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo/bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo!bar",
			inputExtension: "",
			output:         "foo_bar.",
		},
		{
			inputName:      "foo//bar",
			inputExtension: "",
			output:         "foo__bar.",
		},
		{
			inputName:      "//foo//bar//",
			inputExtension: "",
			output:         "__foo__bar__.",
		},
		{
			inputName:      "foo\\\\\\bar",
			inputExtension: "",
			output:         "foo___bar.",
		},
		{
			inputName:      "foo[*]bar",
			inputExtension: "",
			output:         "foo___bar.",
		},
		{
			inputName:      "foo bar",
			inputExtension: "01234567890123456789",
			output:         "foo_bar.01234567890123",
		},
		{
			inputName:      "foo\nbar",
			inputExtension: "txt",
			output:         "foo_bar.txt",
		},
		{
			inputName:      "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789junk",
			inputExtension: "txt",
			output:         "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789.txt",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.inputName, func(t *testing.T) {
			result := Filenamify(tc.inputName, tc.inputExtension)
			assert.Equal(t, tc.output, result)
			assert.LessOrEqual(t, len(result), 255)
		})
	}
}
