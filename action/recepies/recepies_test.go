// SPDX-License-Identifier: MIT
package recepies

import (
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
)

func TestCommonGetOpts(t *testing.T) {
	// Try with actual github action
	action := githubactions.New()
	_, err := commonGetOpts(action.GetInput)
	assert.ErrorIs(t, err, errRequiredOptionUndefined)

	// Try few combinations of empty and non-empty values
	testCases := []struct {
		name    string
		opts    map[string]string
		wantErr error
	}{
		{
			name: "all empty",
			opts: map[string]string{
				"target":           "",
				"sdk_version":      "",
				"repo_path":        "",
				"defconfig_path":   "",
				"containerWorkDir": "",
				"output":           "",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: target",
			opts: map[string]string{
				"target":           "",
				"sdk_version":      "dummy",
				"repo_path":        "dummy",
				"defconfig_path":   "dummy",
				"containerWorkDir": "dummy",
				"GITHUB_WORKSPACE": "dummy",
				// containerWorkDir is assigned value from
				//   GITHUB_WORKSPACE environment variable in commonGetOpts func
				// So for this test to work if containerWorkDir is not empty,
				//   GITHUB_WORKSPACE must be defined in this map and also not empty
				"output": "dummy",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: sdk_version",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "",
				"repo_path":        "dummy",
				"defconfig_path":   "dummy",
				"containerWorkDir": "dummy",
				"GITHUB_WORKSPACE": "dummy",
				"output":           "dummy",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: path",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "dummy",
				"repo_path":        "",
				"defconfig_path":   "dummy",
				"containerWorkDir": "dummy",
				"GITHUB_WORKSPACE": "dummy",
				"output":           "dummy",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: defconfig_path",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "dummy",
				"repo_path":        "dummy",
				"defconfig_path":   "",
				"containerWorkDir": "dummy",
				"GITHUB_WORKSPACE": "dummy",
				"output":           "dummy",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: containerWorkDir",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "dummy",
				"repo_path":        "dummy",
				"defconfig_path":   "dummy",
				"containerWorkDir": "",
				"output":           "dummy",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "one empty: output",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "dummy",
				"repo_path":        "dummy",
				"defconfig_path":   "dummy",
				"containerWorkDir": "dummy",
				"output":           "",
			},
			wantErr: errRequiredOptionUndefined,
		},
		{
			name: "none empty",
			opts: map[string]string{
				"target":           "dummy",
				"sdk_version":      "dummy",
				"repo_path":        "dummy",
				"defconfig_path":   "dummy",
				"containerWorkDir": "dummy",
				"GITHUB_WORKSPACE": "dummy",
				"output":           "dummy",
			},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			getFunc := func(key string) string {
				return tc.opts[key]
			}
			_, err := commonGetOpts(getFunc)
			assert.ErrorIs(t, err, tc.wantErr)
			for key, val := range tc.opts {
				if val == "" {
					assert.ErrorContains(t, err, key)
				}
			}
		})
	}
}
