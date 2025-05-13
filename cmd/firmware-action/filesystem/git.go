// SPDX-License-Identifier: MIT

// Package filesystem / git implements git-related commands
package filesystem

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ErrNotGitRepository is returned when a git command is executed outside of git repository
var ErrNotGitRepository = errors.New("not a git repository")

// gitRun is generic function to execute any git command in sub-directory
func gitRun(subdir string, command []string) (string, error) {
	// subdir is relative to current working directory

	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Change current working directory into the repository / submodule
	defer os.Chdir(pwd) //nolint:errcheck
	err = os.Chdir(subdir)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Failed to change current working directory to '%s'", subdir),
			slog.Any("error", err),
		)
		return "", err
	}

	// Run git describe
	cmd := exec.Command(command[0], command[1:]...)
	stdout, err := cmd.CombinedOutput()
	stdoutStr := string(stdout)
	if err != nil {
		pattern := regexp.MustCompile(`^fatal: not a git repository.*`)
		match := pattern.MatchString(stdoutStr)
		if match {
			err = ErrNotGitRepository
		}

		slog.Error(
			fmt.Sprintf("Failed to run git command in '%s'", subdir),
			slog.String("stdout and stderr", stdoutStr),
			slog.Any("error", err),
		)
		return stdoutStr, err
	}

	return stdoutStr, nil
}

type describe struct {
	abbrev int
	dirty  bool
	always bool
}

// GitDescribeCoreboot is coreboot-specific git describe command (do not touch)
func GitDescribeCoreboot(repoPath string) (string, error) {
	cfg := describe{
		abbrev: 12,
		dirty:  true,
		always: true,
	}

	// Check validity of the returned string
	hash, err := gitDescribe(repoPath, &cfg)

	pattern := regexp.MustCompile(`[\d\w]{13}(\-dirty)?`)
	valid := pattern.MatchString(hash)
	if !valid {
		slog.Warn(
			fmt.Sprintf("Output of 'git describe' for '%s' seems to be invalid, output is: '%s'", repoPath, hash),
		)
	}

	return hash, err
}

// GitDescribe is a generic git describe command, with some sane defaults
func GitDescribe(repoPath string) (string, error) {
	cfg := describe{
		abbrev: 8,
		dirty:  true,
		always: true,
	}
	return gitDescribe(repoPath, &cfg)
}

func gitDescribe(repoPath string, cfg *describe) (string, error) {
	cmd := []string{"git", "describe"}
	if cfg.abbrev > 0 {
		cmd = append(cmd, fmt.Sprintf("--abbrev=%d", cfg.abbrev))
	}
	if cfg.dirty {
		cmd = append(cmd, "--dirty")
	}
	if cfg.always {
		cmd = append(cmd, "--always")
	}

	describe, err := gitRun(repoPath, cmd)
	if err != nil {
		return "", err
	}

	// Remove trailing newline
	result := strings.TrimSpace(describe)

	// Check validity of the returned string
	pattern := regexp.MustCompile(`[\d\w]{13}(\-dirty)?`)
	valid := pattern.MatchString(result)
	if !valid {
		slog.Warn(
			fmt.Sprintf("Output of 'git describe' for '%s' seems to be invalid commit hash, output is '%s'. Will carry on.", repoPath, result),
		)
	}

	return result, err
}
