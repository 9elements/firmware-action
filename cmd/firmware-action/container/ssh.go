// SPDX-License-Identifier: MIT

// Package container for dealing with containers via dagger
package container

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"dagger.io/dagger"
)

// ErrParseAddress is raised when containers IP address could not be parsed
var ErrParseAddress = errors.New("could not parse address string")

// OptsOpenSSH stores options for SSH tunnel for OpenSSH function
type OptsOpenSSH struct {
	WaitFunc  func()      // Waiting function holding container with SSH running
	Password  string      // Filled in by OpenSSH function
	IPv4      string      // Filled in by OpenSSH function
	Port      string      // Filled in by OpenSSH function
	MutexData *sync.Mutex // Mutex for modifying data in this struct

	// We could do with single channel here, but for clarity and less mental overhead there are 2
	TunnelClose chan (bool) // Channel to signal that SSH tunnel is ready
	TunnelReady chan (bool) // Channel to signal that SSH tunnel is not longer needed and can be closed
}

// Wait calls WaitFunc
func (s OptsOpenSSH) Wait() {
	s.WaitFunc()
}

// Address function parses provided string and populates IPv4 and Port (port defaults to 22 if not found)
func (s *OptsOpenSSH) Address(address string) error {
	s.MutexData.Lock()
	var err error

	sshAddressSplit := strings.Split(address, ":")
	switch len(sshAddressSplit) {
	case 1:
		// IP address but no port
		s.IPv4 = sshAddressSplit[0]
		s.Port = "22"
	case 2:
		// Possibly both IP address and port
		s.IPv4 = sshAddressSplit[0]
		if s.IPv4 == "" {
			err = fmt.Errorf("%w: '%s'", ErrParseAddress, address)
		}
		s.Port = sshAddressSplit[1]
		if s.Port == "" {
			s.Port = "22"
		}
	default:
		err = fmt.Errorf("%w: '%s'", ErrParseAddress, address)
	}

	s.MutexData.Unlock()
	return err
}

// SettingsSSH is for functional option pattern
type SettingsSSH func(*OptsOpenSSH)

// WithWaitPressEnter is one possible function to pass into OpenSSH
// It will wait until user presses ENTER key to shutdown the container
func WithWaitPressEnter() SettingsSSH {
	return func(s *OptsOpenSSH) {
		s.WaitFunc = func() {
			<-s.TunnelReady
			fmt.Print("Press ENTER to stop container ")
			fmt.Scanln() // nolint:errcheck
			s.TunnelClose <- true
		}
	}
}

// WithWaitNone is one possible function to pass into OpenSSH
// It will not wait
func WithWaitNone() SettingsSSH {
	return func(s *OptsOpenSSH) {
		s.WaitFunc = func() {
			fmt.Println("Skipping waiting")
		}
	}
}

// NewSettingsSSH returns a SettingsSSH
func NewSettingsSSH(opts ...SettingsSSH) *OptsOpenSSH {
	// Defaults
	var m sync.Mutex
	s := &OptsOpenSSH{
		MutexData:   &m,
		TunnelClose: make(chan (bool)),
		TunnelReady: make(chan (bool)),
	}
	WithWaitPressEnter()(s)

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// OpenSSH takes a container and starts SSH server with port exposed to the host
func OpenSSH(
	ctx context.Context,
	client *dagger.Client,
	container *dagger.Container,
	workdir string,
	opts *OptsOpenSSH,
) error {
	// Example in docs:
	//   https://docs.dagger.io/cookbook/#expose-service-containers-to-host
	// This feature is untested and instead relies on tears, blood and sweat produced during
	//   it's development to work.
	// UPDATE: After more tears, blood and sweat we also have some testing! Yippee!

	if container == nil {
		log.Println("skipping SSH because no container was given")
		return nil
	}

	if workdir == "" {
		workdir = "/"
	}

	// Generate a password for the root user
	opts.MutexData.Lock()
	opts.Password = generatePassword(16)
	opts.MutexData.Unlock()

	// Prepare the container
	container = container.
		WithExec([]string{"bash", "-c", fmt.Sprintf("echo 'root:%s' | chpasswd", opts.Password)}).
		WithExec([]string{"bash", "-c", fmt.Sprintf("echo 'cd %s' >> /root/.bashrc", workdir)}).
		WithExec([]string{"/usr/sbin/sshd", "-D"})

	// ANCHOR: ContainerAsService
	// Convert container to service with exposed SSH port
	const sshPort = 22
	sshServiceDoc := container.WithExposedPort(sshPort).AsService()

	// Expose the SSH server to the host
	sshServiceTunnel, err := client.Host().Tunnel(sshServiceDoc).Start(ctx)
	if err != nil {
		fmt.Println("Problem getting tunnel up")
		return err
	}
	defer sshServiceTunnel.Stop(ctx) // nolint:errcheck
	// ANCHOR_END: ContainerAsService

	// Get and print instructions on how to connect
	sshAddress, err := sshServiceTunnel.Endpoint(ctx)
	errAddr := opts.Address(sshAddress)
	if err != nil || errAddr != nil {
		fmt.Println("problem getting address")
		return errors.Join(err, errAddr)
	}
	fmt.Println("Container was reverted back into a state before failed command was executed")
	fmt.Printf("Connect into the container with:\n  ssh root@%s -p %s -o PreferredAuthentications=password\n", opts.IPv4, opts.Port)
	fmt.Printf("Password is:\n  %s\n", opts.Password)
	fmt.Println("SSH up and running")

	// Wait for user to press key
	go opts.Wait()
	opts.TunnelReady <- true
	<-opts.TunnelClose

	fmt.Println("DONE")
	return nil
}

func generatePassword(length int) string {
	// I suppose we could use crypto/rand, but this seems simpler
	// Also, it is meant only as temporary password for a temporary container which gets
	//   shut-down / removed afterwards. I think it is good enough.
	characters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	pass := make([]rune, length)
	for i := range pass {
		pass[i] = characters[rand.Intn(len(characters))]
	}
	return string(pass)
}
