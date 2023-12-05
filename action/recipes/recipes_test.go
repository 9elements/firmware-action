// SPDX-License-Identifier: MIT
package recipes

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name    string
		wantErr error
		target  string
	}{
		{
			name:    "empty target string",
			wantErr: ErrTargetMissing,
			target:  "",
		},
		{
			name:    "invalid target",
			wantErr: ErrTargetInvalid,
			target:  "dummy",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = Execute(ctx, tc.target, client)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
