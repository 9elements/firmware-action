// SPDX-License-Identifier: MIT

// Package filesystem implements things like moving files, copying them, etc.
package filesystem

import (
	"fmt"
	"io"
	"os"
)

// CheckFileExists checks if file exists at PATH
func CheckFileExists(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}

	fileInfo, err := os.Stat(path)
	if err == nil {
		// path exists
		if fileInfo.IsDir() {
			return fmt.Errorf("path '%s' is a directory", path)
		}
		return os.ErrExist
	}
	return err
}

// checkBeforeCopyOrMove runs some tests needed for both CopyFile and MoveFile
func checkBeforeCopyOrMove(pathSource, pathDestination string) error {
	if err := CheckFileExists(pathSource); os.IsNotExist(err) {
		return err
	}
	if err := CheckFileExists(pathDestination); os.IsExist(err) {
		return fmt.Errorf("destination '%s' %w", pathDestination, err)
	}
	pathSourceStat, err := os.Stat(pathSource)
	if err != nil {
		return err
	}
	if !pathSourceStat.Mode().IsRegular() {
		return fmt.Errorf("file %s is not a regular file", pathSource)
	}
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

// MoveFile moves file from SOURCE to DESTINATION
func MoveFile(pathSource, pathDestination string) error {
	// Checks and tests
	if err := checkBeforeCopyOrMove(pathSource, pathDestination); err != nil {
		return err
	}

	return os.Rename(pathSource, pathDestination)
}
