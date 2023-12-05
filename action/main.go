// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/recipes"
	"github.com/alecthomas/kong"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// CLI = Command line interface
var CLI struct {
	Config string `type:"path" required:"" default:"${config_file}" help:"Path to configuration file"`
	Build  struct {
		Recursive bool   `help:"Build recursively with all dependencies and payloads"`
		Target    string `required:"" help:"Select which target to build"`
	} `cmd:"" help:"Build firmware image"`
	GenerateConfig struct{} `cmd:"" help:"Generate example configuration file"`
}

func run(ctx context.Context) error {
	// Parse command line options
	cliCtx := kong.Parse(
		&CLI,
		kong.Description("Utility to create firmware images for several open source firmware solutions"),
		kong.UsageOnError(),
		kong.Vars{
			"config_file": "firmware-action.json",
		},
	)
	switch cliCtx.Command() {
	case "build":
		fmt.Printf("building\n")
	case "generate-config":
		fmt.Printf("generating config\n")
	default:
		panic(cliCtx.Command())
	}

	// Parse configuration file
	//configCtx := kong.Parse(&Config, kong.Configuration(kong.JSON, CLI.Config))
	//_ = configCtx

	myConfig := recipes.Config{}
	jsonString, err := json.Marshal(myConfig)
	if err != nil {
		fmt.Println("Unable to convert the struct to a JSON string")
	}
	fmt.Println(string(jsonString))
	return nil

	// Setup dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Setup GitHub action SDK
	// TODO: add JSON config loading based on action input
	action := githubactions.New()
	_ = action

	return recipes.Execute(ctx, "coreboot", client)
}
