// SPDX-License-Identifier: MIT

// Package environment
package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name              string
		envCreate         map[string]string
		envSearchFor      []string
		envExpectedResult map[string]string
	}{
		{
			name:              "no env vars created, no found",
			envCreate:         map[string]string{},
			envSearchFor:      []string{"MY_VAR"},
			envExpectedResult: map[string]string{},
		},
		{
			name: "one created, one found",
			envCreate: map[string]string{
				"MY_VAR": "my_val",
			},
			envSearchFor: []string{"MY_VAR"},
			envExpectedResult: map[string]string{
				"MY_VAR": "my_val",
			},
		},
		{
			name: "some created, some found",
			envCreate: map[string]string{
				"MY_VAR":  "my_val",
				"MY_VAR2": "my_val2",
			},
			envSearchFor: []string{"MY_VAR"},
			envExpectedResult: map[string]string{
				"MY_VAR": "my_val",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envCreate {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
				t.Logf("Setting %s = %s\n", key, value)
			}

			result := FetchEnvVars(tc.envSearchFor)

			assert.Equal(t, tc.envExpectedResult, result)
		})
	}
}
