// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently support are coreboot and Linux pipelines.
package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()
	if err := run(context.Background(), action); err != nil {
		action.Fatalf("%v", err)
	}
}

func run(ctx context.Context, action *githubactions.Action) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	switch action.GetInput("target") {
	case "coreboot":
		return coreboot(ctx, action, client)
	case "linux":
		return linux(ctx, action, client)
	case "":
		return fmt.Errorf("no target specified")
	default:
		return fmt.Errorf("unsupported target: %s", action.GetInput("target"))
	}
}
