// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/9elements/firmware-action/action/recepies"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	action := githubactions.New()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()
	return recepies.Execute(ctx, client, action)
}
