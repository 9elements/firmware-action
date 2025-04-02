// SPDX-License-Identifier: MIT

// Package environment for is interacting with environment variables
package environment

import (
	"fmt"
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

// DetectGithub function returns True when the execution environment is detected to be GitHub CI
func DetectGithub() bool {
	// Check for GitHub
	// https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	_, exists := os.LookupEnv("GITHUB_ACTIONS")
	return exists
}

// LogGroupStart is used to group log lines
func LogGroupStart(groupName string) {
	// GitHub
	if DetectGithub() {
		// Docs:
		// https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions#grouping-log-lines
		fmt.Println("::group::" + groupName)
	}

	// TODO: GitLab
	// Docs: https://docs.gitlab.com/ci/jobs/job_logs/#custom-collapsible-sections
}

// LogGroupStop is used to group log lines
func LogGroupStop(groupName string) {
	// GitHub
	if DetectGithub() {
		fmt.Println("::endgroup::")
	}

	// TODO: GitLab (needs 'groupName')
	// Docs: https://docs.gitlab.com/ci/jobs/job_logs/#custom-collapsible-sections
	_ = groupName
}
