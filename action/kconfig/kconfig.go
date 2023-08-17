// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package kconfig

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Kconfig holds the same information as a linux 'defconfig' or '.config' file
type Kconfig map[string]string

const isNotSet string = "is not set"

func NewKconfigFromIoReader(r io.Reader) (*Kconfig, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewKconfig(string(buf))
}

// NewKconfig generates a Kconfig object from text file
// Keys that are "not set" have the '# ' prefix and the value 'isNotSet'
func NewKconfig(c string) (*Kconfig, error) {
	scanner := bufio.NewScanner(strings.NewReader(c))

	out := Kconfig{}
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		if strings.HasPrefix(scanner.Text(), "# ") && strings.HasSuffix(scanner.Text(), " "+isNotSet) {
			val := strings.ReplaceAll(scanner.Text(), " "+isNotSet, "")
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
	keys := make([]string, 0)
	for k, _ := range *c {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.HasPrefix(k, "# ") && c.KeyIsNotSet(k[2:]) {
			s += fmt.Sprintf("%s %s\n", k, isNotSet)
		} else {
			val, _ := c.Value(k)
			s += fmt.Sprintf("%s=%s\n", k, val)
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

// Value returns the value of key. If key isn't found or the entry is "is not set"
// this functions returns an error
func (c Kconfig) Value(k string) (string, error) {
	if c.KeyIsNotSet(k) {
		return "", fmt.Errorf("key is not set")
	}
	val, ok := c[k]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	return val, nil
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
