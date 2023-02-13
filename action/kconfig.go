// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently support are coreboot and Linux pipelines.
package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/sethvargo/go-githubactions"
)

// Kconfig holds the same information as a linux 'defconfig' or '.config' file
type Kconfig map[string]string

const isNotSet string = "is not set"

// NewKconfig generates a Kconfig object from text file
// Keys that are "not set" have the '# ' prefix and the value 'isNotSet'
func NewKconfig(c string) (*Kconfig, error) {
	scanner := bufio.NewScanner(strings.NewReader(c))

	out := Kconfig{}
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		if strings.HasPrefix(scanner.Text(), "# ") && strings.HasSuffix(scanner.Text(), " is not set") {
			val := strings.ReplaceAll(scanner.Text(), " is not set", "")
			out[val] = isNotSet
		} else if strings.Contains(scanner.Text(), "=") {
			arr := strings.Split(scanner.Text(), "=")
			out[arr[0]] = arr[1]
		}
	}

	return &out, scanner.Err()
}

// String generates a multiline representation as this would be
// a defconfig or .config file
func (c *Kconfig) String() string {
	s := ""
	for k, v := range *c {
		if strings.HasPrefix(k, "# ") && v == isNotSet {
			s += fmt.Sprintf("%s %s\n", k, v)
		} else {
			s += fmt.Sprintf("%s=%s\n", k, v)
		}
	}
	return s
}

// KeyIsNotSet returns true if the key "is not set"
// It returns false if the key is set or doesn't exist
func (c Kconfig) KeyIsNotSet(k string) bool {
	val, ok := c["# "+k]
	if !ok {
		return false
	}
	return val == isNotSet
}

// UnsetKey marks the key as "is not set"
func (c *Kconfig) UnsetKey(k string) {
	(*c)["# "+k] = isNotSet
}

// EvalPath replaces "$key" and "$(key)" in the input argument with
// the corresponding value
func (c *Kconfig) EvalPath(p string) string {
	for k, v := range *c {
		if strings.HasPrefix(k, "#") {
			continue
		}
		if strings.Contains(p, "$("+k+")") {
			p = strings.ReplaceAll(p, "$("+k+")", v)
		}
		if strings.Contains(p, "$"+k) {
			p = strings.ReplaceAll(p, "$"+k, v)
		}
	}

	return p
}

func generateDotConfigFromDefconfig(ctx context.Context, action *githubactions.Action, client *dagger.Client, defconfig string) (*Kconfig, error) {
	corebootContainer, err := setupCorebootContainer(ctx, action, client, defconfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .config: %v", err)
	}

	// generate .config
	corebootContainer = corebootContainer.WithExec([]string{"rm", ".config"})
	exitcode, err := corebootContainer.ExitCode(ctx)
	if exitcode != 0 {
		return nil, fmt.Errorf("Non zero exit code %d: %v", exitcode, err)
	}

	o, err := corebootContainer.Stdout(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	o, err = corebootContainer.Stderr(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	corebootContainer = corebootContainer.WithExec([]string{"make", "defconfig", "KBUILD_DEFCONFIG=defconfig"})
	exitcode, err = corebootContainer.ExitCode(ctx)
	if exitcode != 0 {
		return nil, fmt.Errorf("Non zero exit code %d: %v", exitcode, err)
	}
	o, err = corebootContainer.Stdout(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	o, err = corebootContainer.Stderr(ctx)
	if o != "" {
		fmt.Println(o)
	}
	if err != nil {
		return nil, err
	}

	// Extract .config
	rom := corebootContainer.File(".config")
	dotconfigRaw, err := rom.Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .config: %v", err)
	}

	var dotConfig *Kconfig
	dotConfig, err = NewKconfig(dotconfigRaw)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert .config: %v", err)
	}
	return dotConfig, nil
}
