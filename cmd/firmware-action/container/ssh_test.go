// SPDX-License-Identifier: MIT

// Package container
package container

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestAddress(t *testing.T) {
	defaultIPv4 := "192.168.0.1"
	defaultPort := "22"

	testCases := []struct {
		name    string
		address string
		ipv4    string
		port    string
		wantErr error
	}{
		{
			name:    "classic",
			address: fmt.Sprintf("%s:%s", defaultIPv4, defaultPort),
			ipv4:    defaultIPv4,
			port:    defaultPort,
			wantErr: nil,
		},
		{
			name:    "no port",
			address: defaultIPv4,
			ipv4:    defaultIPv4,
			port:    defaultPort,
			wantErr: nil,
		},
		{
			name:    "no port but with colon",
			address: fmt.Sprintf("%s:", defaultIPv4),
			ipv4:    defaultIPv4,
			port:    defaultPort,
			wantErr: nil,
		},
		{
			name:    "no IP but port",
			address: fmt.Sprintf(":%s", defaultPort),
			ipv4:    defaultIPv4,
			port:    defaultPort,
			wantErr: ErrParseAddress,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewSettingsSSH()
			err := opts.Address(tc.address)
			assert.ErrorIs(t, err, tc.wantErr)
			if err != nil {
				// no need to check if errored
				return
			}
			assert.Equal(t, tc.ipv4, opts.IPv4)
			assert.Equal(t, tc.port, opts.Port)
		})
	}
}

func TestOpenSSH(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	assert.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name    string
		opts    SetupOpts
		optsSSH *OptsOpenSSH
	}{
		{
			name: "hello world",
			opts: SetupOpts{
				ContainerURL:      "ghcr.io/9elements/firmware-action/linux_6.1.45:main",
				MountContainerDir: "/src",
				MountHostDir:      ".",
				WorkdirContainer:  "/src",
			},
			optsSSH: NewSettingsSSH(WithWaitNone()),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			myContainer, err := Setup(ctx, client, &tc.opts, "")
			assert.NoError(t, err)

			// Open the SSH
			go func() {
				err := OpenSSH(ctx, client, myContainer, tc.opts.WorkdirContainer, tc.optsSSH)
				if err != nil {
					tc.optsSSH.TunnelClose <- true
				}
				assert.NoError(t, err)
			}()
			// Wait until SSH server is ready
			<-tc.optsSSH.TunnelReady

			// Connect with client
			config := &ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.Password(tc.optsSSH.Password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", tc.optsSSH.IPv4, tc.optsSSH.Port), config)
			assert.NoError(t, err)
			defer client.Close()

			// Open session
			session, err := client.NewSession()
			assert.NoError(t, err)
			defer session.Close()

			// Run simple command to test functionality
			var b bytes.Buffer
			session.Stdout = &b
			err = session.Run("/usr/bin/whoami")
			assert.NoError(t, err)
			assert.Equal(t, "root\n", b.String())

			tc.optsSSH.TunnelClose <- true
		})
	}
}
