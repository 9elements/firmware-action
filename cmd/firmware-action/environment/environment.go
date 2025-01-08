// SPDX-License-Identifier: MIT

// Package environment for is interacting with environment variables
package environment

import (
	"os"
)

// FetchEnvVars when provided with list of environment variables is to
// return a map of variables and values for those that exist in the environment
func FetchEnvVars(variables []string) map[string]string {
	result := make(map[string]string)

	for _, variable := range variables {
		value, exists := os.LookupEnv(variable)
		if exists {
			result[variable] = value
		}
	}

	return result
}
