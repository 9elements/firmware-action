// SPDX-License-Identifier: MIT

// Package recipes yay!
package recipes

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/9elements/firmware-action/cmd/firmware-action/filesystem"
	"github.com/9elements/firmware-action/cmd/firmware-action/logging"
	"github.com/google/go-cmp/cmp"
)

// Change is generic struct to hold commonly needed variables for all change detection modules
type Change struct {
	// ResultFile stores path to a temporary file where information needed to detect changes is stored
	ResultFile string
	// ChangesDetected stores if changes were detected
	ChangesDetected bool
}

//========================
// Time Stamps

// ChangeTimeStamp is for detecting any change in source files based on time-stamps
type ChangeTimeStamp struct {
	Change
	Sources []string
}

// DetectChanges is a method for detecting changes based on Time-Stamp
func (c *ChangeTimeStamp) DetectChanges() bool {
	lastRun, err := filesystem.LoadLastRunTime(c.ResultFile)
	if err != nil {
		// If file does not exist, or can't be read, etc
		c.ChangesDetected = false
		return false
	}
	for _, source := range c.Sources {
		// Either returns time, or zero time and error
		//   zero time means there was no previous run
		changes, _ := filesystem.AnyFileNewerThan(source, lastRun)
		if changes {
			c.ChangesDetected = true
			return true
		}
	}
	c.ChangesDetected = false
	return false
}

// SaveCheckpoint is a method for saving checkpoint file for future change detection
func (c *ChangeTimeStamp) SaveCheckpoint(override bool) {
	// On success update the timestamp
	err := filesystem.CheckFileExists(c.ResultFile)
	if errors.Is(err, os.ErrNotExist) || override {
		slog.Debug("Saving timestamp checkpoint")
		_ = filesystem.SaveCurrentRunTime(c.ResultFile)
	}
}

//========================
// Configuration

// ChangeConfig is for detecting any change in firmware-action configuration
type ChangeConfig struct {
	Change
	Config *Config
}

// DetectChanges is a method for detecting changes based on Configuration file
func (c *ChangeConfig) DetectChanges(target string) bool {
	// I did consider to save only the small struct related to each module, but it was
	//   proving to be far too much work. Instead we save the whole configuration file (for each module
	//   separately) and only compare the relevant modules between these two configurations

	err := filesystem.CheckFileExists(c.ResultFile)
	if errors.Is(err, os.ErrExist) {
		oldConfig, err := ReadConfig(c.ResultFile)
		// The config might be old / obsolete
		// If the config is old / obsolete and no longer valid, it should just be ignored
		// and it should be assumed that re-build is needed
		if err != nil {
			slog.Warn(
				fmt.Sprintf("The configuration used for previous build, stored in '%s', is not valid and will be assumed obsolete", CompiledConfigsDir),
			)
			c.ChangesDetected = true
			return true
		}

		oldModules := oldConfig.AllModules()
		modules := c.Config.AllModules()
		c.ChangesDetected = !cmp.Equal(modules[target], oldModules[target])
		return c.ChangesDetected
	}

	// The file might be missing (user deleted it, CI did not cache it, etc.)
	// If the old config file is missing, just return false since no changes can be detected
	c.ChangesDetected = false
	return false
}

// SaveCheckpoint is a method for saving checkpoint file for future change detection
func (c *ChangeConfig) SaveCheckpoint(override bool) {
	err := filesystem.CheckFileExists(c.ResultFile)
	if errors.Is(err, os.ErrNotExist) || override {
		slog.Debug("Saving config checkpoint")

		dir := filepath.Dir(c.ResultFile)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			slog.Error("Cannot create directory for files to aid in change detection",
				slog.Any("error", err),
			)
		}

		err = WriteConfig(c.ResultFile, c.Config)
		if err != nil {
			slog.Warn("Failed to create a snapshot of configuration for detecting future changes",
				slog.Any("error", err),
			)
		}
	}
}

//========================
// Git commit hash

// ChangeGitHash is for detecting any change in git current commit hash
type ChangeGitHash struct {
	Change
	RepoPath           string
	currentGitDescribe string
}

func (c *ChangeGitHash) gitDescribe() {
	var err error
	c.currentGitDescribe, err = filesystem.GitDescribe(c.RepoPath)
	if err != nil {
		slog.Warn(
			"Failed to get git version",
			slog.String("git repository path", c.RepoPath),
			slog.String("suggestion", "Double check the RepoPath is correct, check that all git submodules are initialized"),
			slog.Any("error", err),
		)
		c.currentGitDescribe = ""
	}
}

// DetectChanges is a method for detecting changes based on Git commit hash
func (c *ChangeGitHash) DetectChanges() bool {
	// Update git describe
	c.gitDescribe()

	if c.currentGitDescribe != "" {
		content, err := os.ReadFile(c.ResultFile)
		if err == nil {
			// File exists and we can read it
			lastGitVersion := strings.TrimSpace(string(content))
			if c.currentGitDescribe != lastGitVersion {
				c.ChangesDetected = true
				return true
			}
		}
	}
	c.ChangesDetected = false
	return false
}

// SaveCheckpoint is a method for saving checkpoint file for future change detection
func (c *ChangeGitHash) SaveCheckpoint(target string, override bool) {
	// Update git describe
	c.gitDescribe()

	err := filesystem.CheckFileExists(c.ResultFile)
	if c.currentGitDescribe != "" && (errors.Is(err, os.ErrNotExist) || override) {
		slog.Debug("Saving git commit hash checkpoint")

		dir := filepath.Dir(c.ResultFile)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			slog.Error("Cannot create directory for files to aid in change detection",
				slog.Any("error", err),
			)
		}

		err := os.WriteFile(c.ResultFile, []byte(c.currentGitDescribe), 0o644)
		if err != nil {
			slog.Warn(
				fmt.Sprintf("Failed to save git version for '%s'", target),
				slog.Any("error", err),
			)
		}
	}
}

//========================
// Congregation of all change detection methods

// AllChanges is congregation of all methods of change detection
type AllChanges struct {
	TimeStamp     ChangeTimeStamp
	Configuration ChangeConfig
	GitHash       ChangeGitHash
}

// DetectChanges is a method for detecting changes based on combination of multiple methods
func (c *AllChanges) DetectChanges(target string) bool {
	c.TimeStamp.DetectChanges()
	c.Configuration.DetectChanges(target)
	c.GitHash.DetectChanges()

	result := c.TimeStamp.ChangesDetected || c.Configuration.ChangesDetected || c.GitHash.ChangesDetected

	// Debug output
	slog.Debug("Detected changes",
		slog.Bool("time-stamp", c.TimeStamp.ChangesDetected),
		slog.Bool("config", c.Configuration.ChangesDetected),
		slog.Bool("git-hash", c.GitHash.ChangesDetected),
		slog.Bool("conclusion", result),
	)

	return result
}

// SaveCheckpoint is a method for saving checkpoint files for future change detection
func (c *AllChanges) SaveCheckpoint(target string, override bool) {
	slog.Debug("Saving change detection checkpoints")
	c.TimeStamp.SaveCheckpoint(override)
	c.Configuration.SaveCheckpoint(override)
	c.GitHash.SaveCheckpoint(target, override)
}
