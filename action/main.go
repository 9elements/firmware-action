// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// Documentation [is hosted in GitHub pages](https://9elements.github.io/firmware-action/)
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/9elements/firmware-action/action/filesystem"
	"github.com/9elements/firmware-action/action/recipes"
	"github.com/alecthomas/kong"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// CLI (Command Line Interface) holds data from environment
var CLI struct {
	Config string `type:"path" required:"" default:"${config_file}" help:"Path to configuration file"`

	Build struct {
		Target    string `required:"" help:"Select which target to build, use ID from configuration file"`
		Recursive bool   `help:"Build recursively with all dependencies and payloads"`
	} `cmd:"build" help:"Build a target defined in configuration file."`

	GenerateConfig struct{} `cmd:"generate-config" help:"Generate empty configuration file."`
}

func run(ctx context.Context) error {
	err := getInputsFromEnvironment()
	if err != nil {
		return err
	}
	log.Printf("Inputs:\nConfig:    %s\nTarget:    %s\nRecursive: %t\n", CLI.Config, CLI.Build.Target, CLI.Build.Recursive)

	// Parse configuration file
	var myConfig recipes.Config
	myConfig, err = recipes.ReadConfig(CLI.Config)
	if err != nil {
		return err
	}

	// Lets build stuff
	_, err = recipes.Build(ctx, CLI.Build.Target, CLI.Build.Recursive, myConfig, recipes.Execute)
	return err
}

func getInputsFromEnvironment() error {
	// Check for GitHub
	// https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	_, exists := os.LookupEnv("GITHUB_ACTIONS")
	if exists {
		log.Print("Running in GitHub mode")
		return parseGithub()
	}

	// Check for GitLab, ... (possibly add other CIs)
	// TODO

	// Use command line interface
	return parseCli()
}

func parseCli() error {
	// Get inputs from command line options
	log.Print("Running in CLI mode")
	ctx := kong.Parse(
		&CLI,
		kong.Description("Utility to create firmware images for several open source firmware solutions"),
		kong.UsageOnError(),
		kong.Vars{
			"config_file": "firmware-action.json",
		},
	)

	switch ctx.Command() {
	case "build":
		// This is handled elsewhere
		return nil

	case "generate-config":
		// Check if config file exists
		err := filesystem.CheckFileExists(CLI.Config)
		if !errors.Is(err, os.ErrNotExist) {
			// The file exists, or is directory, or some other problem
			log.Printf("Can't generate configuration file at: %s", CLI.Config)
			return fmt.Errorf("%w: %s", err, CLI.Config)
		}

		// Create empty config
		myConfig := recipes.Config{
			Coreboot: map[string]recipes.CorebootOpts{"coreboot-example": {}},
			Linux:    map[string]recipes.LinuxOpts{"linux-example": {}},
			Edk2:     map[string]recipes.Edk2Opts{"edk2-example": {}},
		}

		// Convert to JSON
		jsonString, err := json.MarshalIndent(myConfig, "", "  ")
		if err != nil {
			fmt.Println("Unable to convert the config struct into a JSON string")
		}

		// Write to file
		log.Printf("Generating configuration file at: %s", CLI.Config)
		if err := os.WriteFile(CLI.Config, jsonString, 0o666); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	default:
		log.Fatal("Supplied unsupported command")
	}

	return nil
}

func parseGithub() error {
	// Get inputs from GitHub environment
	action := githubactions.New()
	regexTrue := regexp.MustCompile(`(?i)true`)

	CLI.Config = action.GetInput("config")
	CLI.Build.Target = action.GetInput("target")
	CLI.Build.Recursive = regexTrue.MatchString(action.GetInput("recursive"))

	return nil
}
