// SPDX-License-Identifier: MIT

// Package filesystem implements things like moving files, copying them, etc.
package filesystem

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/plus3it/gorecurcopy"
)

var (
	// ErrEmptyPath is returned when function is called with empty path parameter
	ErrEmptyPath = errors.New("provided path is empty")
	// ErrPathIsDirectory is returned when path exists, but is a directory and not a file
	ErrPathIsDirectory = fmt.Errorf("provided path is directory: %w", os.ErrExist)
	// ErrFileNotRegular is returned when path exists, but is not a regular file
	ErrFileNotRegular = errors.New("file is not regular file")
	// ErrFileModified is returned when a file in given path was modified
	ErrFileModified = errors.New("file has been modified since")
)

// CheckFileExists checks if file exists at PATH
func CheckFileExists(path string) error {
	// Possible returns:
	//   ErrEmptyPath		if path parameter is empty
	//   os.ErrNotExist		if path does not exists
	//   os.ErrExist		if path exists and is file
	//   ErrPathIsDirectory		if path exists and is directory
	//   ErrFileNotRegular		if path exists and is not regular file or directory
	if path == "" {
		return fmt.Errorf("%w: %s", ErrEmptyPath, path)
	}

	fileInfo, err := os.Stat(path)
	if err == nil {
		// path exists
		if fileInfo.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathIsDirectory, path)
		}
		return os.ErrExist
	}
	return err
}

// checkBeforeCopyOrMove runs some tests needed for both CopyFile and MoveFile
func checkBeforeCopyOrMove(pathSource, pathDestination string) error {
	// Source must exists (be it file or directory)
	err := CheckFileExists(pathSource)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Destination must not exists
	if err := CheckFileExists(pathDestination); errors.Is(err, os.ErrExist) {
		return fmt.Errorf("destination '%s' %w", pathDestination, err)
	}

	// Source must be a regular file or directory
	pathSourceStat, err := os.Stat(pathSource)
	if err != nil {
		return err
	}
	if !pathSourceStat.Mode().IsRegular() && !pathSourceStat.IsDir() {
		return fmt.Errorf("%w: %s", ErrFileNotRegular, pathSource)
	}

	// No problems found
	return nil
}

// CopyFile copies file from SOURCE to DESTINATION
func CopyFile(pathSource, pathDestination string) error {
	// Checks and tests
	if err := checkBeforeCopyOrMove(pathSource, pathDestination); err != nil {
		return err
	}

	// Open source and destination
	src, err := os.Open(pathSource)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(pathDestination)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	_, err = io.Copy(dst, src)
	return err
}

// CopyDir copies directory recursively from SOURCE to DESTINATION
func CopyDir(pathSource, pathDestination string) error {
	// Checks and tests
	if err := checkBeforeCopyOrMove(pathSource, pathDestination); err != nil {
		return err
	}

	// Create a destination directory
	if err := os.MkdirAll(pathDestination, 0o750); err != nil {
		return err
	}

	// Copy
	return gorecurcopy.CopyDirectory(pathSource, pathDestination)
}

// MoveFile moves file from SOURCE to DESTINATION
func MoveFile(pathSource, pathDestination string) error {
	// Checks and tests
	if err := checkBeforeCopyOrMove(pathSource, pathDestination); err != nil {
		return err
	}

	return os.Rename(pathSource, pathDestination)
}

// DirTree is equivalent to "tree" command
func DirTree(root string) ([]string, error) {
	var files []string

	// WalkDir is faster than Walk
	// https://pkg.go.dev/path/filepath#Walk
	//   > Walk is less efficient than WalkDir, introduced in Go 1.16, which avoids
	//   > calling os.Lstat on every visited file or directory.
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, _ error) error {
		foundItem := path
		if info.IsDir() {
			foundItem = fmt.Sprintf("%s/", path)
		}
		files = append(files, foundItem)
		return nil
	})

	for _, file := range files {
		fmt.Println(file)
	}

	return files, err
}

// LoadLastRunTime loads time of the last execution from file
func LoadLastRunTime(pathLastRun string) (time.Time, error) {
	// Return zero time if file doesn't exist
	err := CheckFileExists(pathLastRun)
	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}, nil
	}

	content, err := os.ReadFile(pathLastRun)
	// Return zero and error on reading error
	if err != nil {
		slog.Warn(
			fmt.Sprintf("Error when reading file '%s'", pathLastRun),
			slog.Any("error", err),
		)
		return time.Time{}, err
	}

	// File should contain time in RFC3339Nano format
	lastRun, err := time.Parse(time.RFC3339Nano, string(content))
	// Return zero and error on parsing error
	if err != nil {
		slog.Warn(
			fmt.Sprintf("Error when parsing time-stamp from '%s'", pathLastRun),
			slog.Any("error", err),
		)
		return time.Time{}, err
	}
	return lastRun, nil
}

// SaveCurrentRunTime writes the current time into file
func SaveCurrentRunTime(pathLastRun string) error {
	// Create temporaryFilesDir

	// Create directory if needed
	dir := filepath.Dir(pathLastRun)
	err := CheckFileExists(dir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Write the current time into file
	return os.WriteFile(pathLastRun, []byte(time.Now().Format(time.RFC3339Nano)), 0o666)
}

// GetFileModTime returns modification time of a file
func GetFileModTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// AnyFileNewerThan checks recursively if any file in given path (can be directory or file) has
// modification time newer than the given time.
// Returns:
// - true if a newer file is found
// - false if no newer file is found or givenTime is zero
// Function is lazy, and returns on first positive occurrence.
func AnyFileNewerThan(path string, givenTime time.Time) (bool, error) {
	// If path does not exist
	err := CheckFileExists(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	// If given time is zero, assume up-to-date
	// This is handy especially for CI, where we can't assume that people will cache firmware-action
	//   timestamp directory, but they will likely cache the produced files
	if givenTime.Equal(time.Time{}) {
		return false, nil
	}

	// If path is directory
	if errors.Is(err, ErrPathIsDirectory) {
		errMod := filepath.WalkDir(path, func(path string, info os.DirEntry, _ error) error {
			// skip .git
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				fileInfo, err := info.Info()
				if err != nil {
					return err
				}
				if fileInfo.ModTime().After(givenTime) {
					return fmt.Errorf("file '%s' has been modified: %w", path, ErrFileModified)
				}
			}
			return nil
		})
		if errors.Is(errMod, ErrFileModified) {
			return true, nil
		}
		return false, nil
	}

	// If path is file
	if errors.Is(err, os.ErrExist) {
		modTime, errMod := GetFileModTime(path)
		if errMod != nil {
			slog.Warn(
				fmt.Sprintf("Encountered error when getting modification time of file '%s'", path),
				slog.Any("error", errMod),
			)
			return false, errMod
		}
		return modTime.After(givenTime), nil
	}

	// If path is neither file nor directory
	return false, err
}
