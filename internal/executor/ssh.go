package executor

import (
	"fmt"
	"strings"
	"time"

	"github.com/instacli/instacli/internal/installers"
	"golang.org/x/crypto/ssh"
)

// SSHConfig holds SSH connection details
type SSHConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	PrivateKey string // Path to private key
}

// SSHExecutor executes commands over SSH
type SSHExecutor struct {
	config SSHConfig
	client *ssh.Client
	os     installers.OS
	pm     installers.PackageManager
}

// NewSSHExecutor creates a new SSH executor
func NewSSHExecutor(config SSHConfig) *SSHExecutor {
	if config.Port == 0 {
		config.Port = 22
	}
	return &SSHExecutor{
		config: config,
		os:     installers.OSLinux, // SSH is typically to Linux servers
	}
}

// Connect establishes SSH connection
func (e *SSHExecutor) Connect() error {
	var authMethods []ssh.AuthMethod

	// Password authentication
	if e.config.Password != "" {
		authMethods = append(authMethods, ssh.Password(e.config.Password))
	}

	// Key-based authentication would be added here
	// if e.config.PrivateKey != "" { ... }

	sshConfig := &ssh.ClientConfig{
		User:            e.config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Proper host key verification
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	e.client = client
	e.detectPackageManager()
	return nil
}

// Close closes the SSH connection
func (e *SSHExecutor) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

func (e *SSHExecutor) detectPackageManager() {
	// Check for apt
	if output, _ := e.Run("which apt"); strings.Contains(output, "apt") {
		e.pm = installers.PMApt
		return
	}
	// Check for dnf
	if output, _ := e.Run("which dnf"); strings.Contains(output, "dnf") {
		e.pm = installers.PMDnf
		return
	}
	// Check for yum
	if output, _ := e.Run("which yum"); strings.Contains(output, "yum") {
		e.pm = installers.PMYum
		return
	}
	e.pm = installers.PMNone
}

// Run executes a command over SSH
func (e *SSHExecutor) Run(cmd string) (string, error) {
	if e.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := e.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return strings.TrimSpace(string(output)), err
}

// RunWithProgress runs a command with streaming output
func (e *SSHExecutor) RunWithProgress(cmd string, onOutput func(line string)) error {
	if e.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := e.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}

	if err := session.Start(cmd); err != nil {
		return err
	}

	// Stream stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				onOutput(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// Stream stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				onOutput(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	return session.Wait()
}

// GetOS returns the target OS
func (e *SSHExecutor) GetOS() installers.OS {
	return e.os
}

// GetPackageManager returns the detected package manager
func (e *SSHExecutor) GetPackageManager() installers.PackageManager {
	return e.pm
}

// IsRoot checks if connected as root
func (e *SSHExecutor) IsRoot() bool {
	output, _ := e.Run("id -u")
	return strings.TrimSpace(output) == "0"
}

// Verify interface implementation
var _ installers.Executor = (*SSHExecutor)(nil)
