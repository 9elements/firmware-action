// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"

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
	client, err := dagger.Connect(ctx)
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
