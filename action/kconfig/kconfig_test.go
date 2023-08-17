// SPDX-License-Identifier: MIT

// Package main implements the core logic of running composable Dagger pipelines
// via GitHub Actions. Currently supported are coreboot and Linux pipelines.
package kconfig

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewKconfig(t *testing.T) {
	type args struct {
		c string
	}
	tests := []struct {
		name    string
		args    args
		want    *Kconfig
		wantErr bool
	}{
		{
			name:    "Empty string",
			args:    args{c: ""},
			want:    &Kconfig{},
			wantErr: false,
		},
		{
			name: "Not set",
			args: args{c: "# CONFIG_FOO is not set"},
			want: &Kconfig{
				"# CONFIG_FOO": "is not set",
			},
			wantErr: false,
		},
		{
			name: "set y",
			args: args{c: "CONFIG_BAR=y"},
			want: &Kconfig{
				"CONFIG_BAR": "y",
			},
			wantErr: false,
		},
		{
			name: "set number",
			args: args{c: "\nCONFIG_BAR=100000"},
			want: &Kconfig{
				"CONFIG_BAR": "100000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKconfig(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKconfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKconfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKconfig_NewKconfigFromIoReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte("CONFIG_BAR=100000"))
	_, err := NewKconfigFromIoReader(buf)
	if err != nil {
		t.Errorf("NewKconfig() returned err %v", err)
	}
}

func TestKconfig_UnsetKey(t *testing.T) {
	k, err := NewKconfig("CONFIG_BAR=100000")
	if err != nil {
		t.Errorf("NewKconfig() returned err %v", err)
	}
	if k.KeyIsNotSet("CONFIG_BAR") {
		t.Errorf("Key CONFIG_BAR is not set, but should be")
	}
	k.UnsetKey("CONFIG_BAR")
	if !k.KeyIsNotSet("CONFIG_BAR") {
		t.Errorf("Key CONFIG_BAR is still set")
	}
}

func TestKconfig_EvalPath(t *testing.T) {
	k, err := NewKconfig("CONFIG_PATH=abc/cde\n# CONFIG_OTHER_PATH is not set")
	if err != nil {
		t.Errorf("NewKconfig() returned err %v", err)
	}
	if k.EvalPath("123$CONFIG_PATH456") != "123abc/cde456" {
		t.Errorf("Expected new path %s, but got %s", "123abc/cde456", k.EvalPath("123$CONFIG_PATH456"))
	}
	if k.EvalPath("123$(CONFIG_PATH)456") != "123abc/cde456" {
		t.Errorf("Expected new path %s, but got %s", "123abc/cde456", k.EvalPath("123$CONFIG_PATH456"))
	}
}

func TestKconfig_String(t *testing.T) {
	k, err := NewKconfig("CONFIG_PATH=abc/cde\n# CONFIG_OTHER_PATH is not set")
	if err != nil {
		t.Errorf("NewKconfig() returned err %v", err)
	}
	if k.String() != "# CONFIG_OTHER_PATH is not set\nCONFIG_PATH=abc/cde\n" {
		t.Errorf("String() returned %s, but should be %s", k.String(), "# CONFIG_OTHER_PATH is not set\nCONFIG_PATH=abc/cde\n")
	}
}

func TestKconfig_Value(t *testing.T) {
	k, err := NewKconfig("CONFIG_PATH=abc/cde\n# CONFIG_OTHER_PATH is not set")
	if err != nil {
		t.Errorf("NewKconfig() returned err %v", err)
	}
	_, err = k.Value("CONFIG_OTHER_PATH")
	if err == nil {
		t.Errorf("Expected error, but got none")
	}
	_, err = k.Value("CONFIG_DOES_NOT_EXIST")
	if err == nil {
		t.Errorf("Expected error, but got none")
	}
}
