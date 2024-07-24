// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// Documentation [is hosted in GitHub pages](https://9elements.github.io/firmware-action/)
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/9elements/firmware-action/action/logging"
	"github.com/9elements/firmware-action/action/recipes"
	"github.com/alecthomas/kong"
	"github.com/go-git/go-git/v5"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	logging.InitLogger(slog.LevelInfo)

	if err := run(context.Background()); err != nil {
		slog.Error(
			"firmware-action failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}

const firmwareActionVersion = "v0.3.2"

// CLI (Command Line Interface) holds data from environment
var CLI struct {
	JSON   bool `default:"false" help:"switch to JSON stdout and stderr output"`
	Indent bool `default:"false" help:"enable indentation for JSON output"`
	Debug  bool `default:"false" help:"increase verbosity"`

	Config string `type:"path" required:"" default:"${config_file}" help:"Path to configuration file"`

	Build struct {
		Target      string `required:"" help:"Select which target to build, use ID from configuration file"`
		Recursive   bool   `help:"Build recursively with all dependencies and payloads"`
		Interactive bool   `help:"Open interactive SSH into container if build fails"`
	} `cmd:"build" help:"Build a target defined in configuration file"`

	GenerateConfig struct{} `cmd:"generate-config" help:"Generate empty configuration file"`
	Version        struct{} `cmd:"version" help:"Print version and exit"`
}

func run(ctx context.Context) error {
	// Get arguments
	mode, err := getInputsFromEnvironment()
	if err != nil {
		return err
	}
	if mode == "" {
		// Exit on "version" or "generate-config"
		return nil
	}

	// Properly initialize logging
	level := slog.LevelInfo
	if CLI.Debug {
		level = slog.LevelDebug
	}
	logging.InitLogger(
		level,
		logging.WithJSON(CLI.JSON),
		logging.WithIndent(CLI.Indent),
	)
	slog.Info(
		fmt.Sprintf("Running in %s mode", mode),
		slog.String("input/config", CLI.Config),
		slog.String("input/target", CLI.Build.Target),
		slog.Bool("input/recursive", CLI.Build.Recursive),
		slog.Bool("input/interactive", CLI.Build.Interactive),
	)

	// Check if submodules were initialized
	// If they are uninitialized, just print warning to the user, ignore all errors
	// Docs: https://git-scm.com/docs/git-submodule#Documentation/git-submodule.txt-status--cached--recursive--ltpathgt82308203
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repo, err := git.PlainOpen(pwd)
	if err == nil {
		worktree, err := repo.Worktree()
		if err != nil {
			goto submodule_out
		}
		submodules, err := worktree.Submodules()
		if err != nil {
			goto submodule_out
		}
		status, err := submodules.Status()
		if err != nil {
			goto submodule_out
		}
		pattern := regexp.MustCompile(`(?m)^\-.*`)
		matches := pattern.FindAllString(status.String(), -1)
		for _, v := range matches {
			patterSub := regexp.MustCompile(`^\-[\d\w]* `)
			slog.Warn(
				"Git submodule seems to be uninitialized",
				slog.String("suggestion", "run 'git submodule update --init --recursive'"),
				slog.String("offending_submodule", patterSub.ReplaceAllString(v, "")),
			)
		}
	}
submodule_out:

	// Parse configuration file
	var myConfig *recipes.Config
	myConfig, err = recipes.ReadConfig(CLI.Config)
	if err != nil {
		return err
	}

	// Lets build stuff
	_, err = recipes.Build(
		ctx,
		CLI.Build.Target,
		CLI.Build.Recursive,
		CLI.Build.Interactive,
		myConfig,
		recipes.Execute,
	)
	if err == nil || errors.Is(err, recipes.ErrBuildSkipped) {
		slog.Info("Build finished successfully")
		return nil
	}
	return err
}

func getInputsFromEnvironment() (string, error) {
	// Check for GitHub
	// https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	_, exists := os.LookupEnv("GITHUB_ACTIONS")
	if exists {
		return parseGithub()
	}

	// Check for GitLab, ... (possibly add other CIs)
	// TODO

	// Use command line interface
	return parseCli()
}

func parseCli() (string, error) {
	// Get inputs from command line options
	ctx := kong.Parse(
		&CLI,
		kong.Description("Utility to create firmware images for several open source firmware solutions"),
		kong.UsageOnError(),
		kong.Vars{
			"config_file": "firmware-action.json",
		},
	)
	mode := "CLI"

	switch ctx.Command() {
	case "build":
		// This is handled elsewhere
		return mode, nil

	case "generate-config":
		// Check if config file exists
		err := filesystem.CheckFileExists(CLI.Config)
		if !errors.Is(err, os.ErrNotExist) {
			// The file exists, or is directory, or some other problem
			slog.Error(
				fmt.Sprintf("Can't generate configuration file at: %s", CLI.Config),
				slog.Any("error", err),
			)
			return "", err
		}

		// Create empty config
		myConfig := recipes.Config{
			Coreboot:          map[string]recipes.CorebootOpts{"coreboot-example": {}},
			Linux:             map[string]recipes.LinuxOpts{"linux-example": {}},
			Edk2:              map[string]recipes.Edk2Opts{"edk2-example": {}},
			FirmwareStitching: map[string]recipes.FirmwareStitchingOpts{"stitching-example": {}},
		}

		// Convert to JSON
		jsonString, err := json.MarshalIndent(myConfig, "", "  ")
		if err != nil {
			slog.Error(
				"Unable to convert the config struct into a JSON string",
				slog.String("suggestion", logging.ThisShouldNotHappenMessage),
				slog.Any("error", err),
			)
			return "", err
		}

		// Write to file
		slog.Info(fmt.Sprintf("Generating configuration file at: %s", CLI.Config))
		if err := os.WriteFile(CLI.Config, jsonString, 0o666); err != nil {
			slog.Error(
				"Unable to write generated configuration into file",
				slog.Any("error", err),
			)
			return "", err
		}
		return "", nil

	case "version":
		// Print version and exit
		fmt.Println(firmwareActionVersion)
		return "", nil

	default:
		// This should not happen
		err := errors.New("unsupported command")
		slog.Error(
			"Supplied unsupported command",
			slog.String("suggestion", logging.ThisShouldNotHappenMessage),
			slog.Any("error", err),
		)
		return mode, err
	}
}

func parseGithub() (string, error) {
	// Get inputs from GitHub environment
	action := githubactions.New()
	regexTrue := regexp.MustCompile(`(?i)true`)

	CLI.Config = action.GetInput("config")
	CLI.Build.Target = action.GetInput("target")
	CLI.Build.Recursive = regexTrue.MatchString(action.GetInput("recursive"))
	CLI.JSON = regexTrue.MatchString(action.GetInput("json"))

	return "GitHub", nil
}
