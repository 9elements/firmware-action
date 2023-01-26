// SPDX-License-Identifier: MIT

package action

import (
	"context"

	"dagger.io/dagger"
	"github.com/sethvargo/go-githubactions"
)

func linux(ctx context.Context, action *githubactions.Action, client *dagger.Client) error {
	return nil
}
