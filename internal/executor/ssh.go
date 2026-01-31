package executor

import (
	"fmt"
	"os"
	"path/filepath"
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
	PrivateKey string // Path to private key or the key content itself
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

	// Try SSH key authentication first (more common for servers)
	keyAuth := e.tryKeyAuth()
	if keyAuth != nil {
		authMethods = append(authMethods, keyAuth)
	}

	// Password authentication as fallback
	if e.config.Password != "" {
		authMethods = append(authMethods, ssh.Password(e.config.Password))
		// Also try keyboard-interactive for some servers
		authMethods = append(authMethods, ssh.KeyboardInteractive(
			func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range questions {
					answers[i] = e.config.Password
				}
				return answers, nil
			},
		))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication method provided (password or SSH key required)")
	}

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

// tryKeyAuth attempts to load SSH private key for authentication
func (e *SSHExecutor) tryKeyAuth() ssh.AuthMethod {
	// Check if PrivateKey is provided in config
	if e.config.PrivateKey != "" {
		// Check if it's a file path or key content
		if strings.HasPrefix(e.config.PrivateKey, "-----BEGIN") {
			// It's the key content itself
			signer, err := ssh.ParsePrivateKey([]byte(e.config.PrivateKey))
			if err == nil {
				return ssh.PublicKeys(signer)
			}
		} else {
			// It's a file path
			key, err := os.ReadFile(e.config.PrivateKey)
			if err == nil {
				signer, err := ssh.ParsePrivateKey(key)
				if err == nil {
					return ssh.PublicKeys(signer)
				}
			}
		}
	}

	// Try default SSH key locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	// Common SSH key locations
	keyPaths := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
		filepath.Join(homeDir, ".ssh", "id_dsa"),
	}

	for _, keyPath := range keyPaths {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			continue
		}

		return ssh.PublicKeys(signer)
	}

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
