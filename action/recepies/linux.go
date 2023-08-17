// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package recepies

import (
	"context"

	"dagger.io/dagger"
	"github.com/sethvargo/go-githubactions"
)

func linux(ctx context.Context, action *githubactions.Action, client *dagger.Client) error {
	return nil
}
