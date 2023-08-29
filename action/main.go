// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/recipes"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// Setup dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Setup GitHub action SDK
	action := githubactions.New()

	return recipes.Execute(ctx, client, action)
}
