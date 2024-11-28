// SPDX-License-Identifier: MIT

// Package filesystem implements things like moving files, copying them, etc.
package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/plus3it/gorecurcopy"
)

var (
	// ErrEmptyPath is returned when function is called with empty path parameter
	ErrEmptyPath = errors.New("provided path is empty")
	// ErrPathIsDirectory is returned when path exists, but is a directory and not a file
	ErrPathIsDirectory = fmt.Errorf("provided path is directory: %w", os.ErrExist)
	// ErrFileNotRegular is returned when path exists, but is not a regular file
	ErrFileNotRegular = errors.New("file is not regular file")
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

	err := filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {
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
