package executor

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"

	"github.com/instacli/instacli/internal/installers"
)

// LocalExecutor executes commands on the local machine
type LocalExecutor struct {
	os installers.OS
	pm installers.PackageManager
}

// NewLocalExecutor creates a new local executor
func NewLocalExecutor() *LocalExecutor {
	e := &LocalExecutor{}
	e.detectOS()
	e.detectPackageManager()
	return e
}

func (e *LocalExecutor) detectOS() {
	switch runtime.GOOS {
	case "linux":
		e.os = installers.OSLinux
	case "darwin":
		e.os = installers.OSMacOS
	case "windows":
		e.os = installers.OSWindows
	default:
		e.os = installers.OSLinux
	}
}

func (e *LocalExecutor) detectPackageManager() {
	switch e.os {
	case installers.OSLinux:
		// Check for apt
		if _, err := exec.LookPath("apt"); err == nil {
			e.pm = installers.PMApt
			return
		}
		// Check for dnf
		if _, err := exec.LookPath("dnf"); err == nil {
			e.pm = installers.PMDnf
			return
		}
		// Check for yum
		if _, err := exec.LookPath("yum"); err == nil {
			e.pm = installers.PMYum
			return
		}
		// Check for pacman
		if _, err := exec.LookPath("pacman"); err == nil {
			e.pm = installers.PMPacman
			return
		}
		e.pm = installers.PMNone

	case installers.OSMacOS:
		e.pm = installers.PMBrew

	case installers.OSWindows:
		e.pm = installers.PMChoco
	}
}

// Run executes a command and returns the output
func (e *LocalExecutor) Run(cmd string) (string, error) {
	var shell, flag string

	switch e.os {
	case installers.OSWindows:
		shell = "powershell"
		flag = "-Command"
	default:
		shell = "bash"
		flag = "-c"
	}

	command := exec.Command(shell, flag, cmd)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	output := stdout.String()
	if output == "" {
		output = stderr.String()
	}

	return strings.TrimSpace(output), err
}

// RunWithProgress runs a command with real-time output
func (e *LocalExecutor) RunWithProgress(cmd string, onOutput func(line string)) error {
	var shell, flag string

	switch e.os {
	case installers.OSWindows:
		shell = "powershell"
		flag = "-Command"
	default:
		shell = "bash"
		flag = "-c"
	}

	command := exec.Command(shell, flag, cmd)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		return err
	}

	// Read stdout
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

	// Read stderr
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

	return command.Wait()
}

// GetOS returns the detected OS
func (e *LocalExecutor) GetOS() installers.OS {
	return e.os
}

// GetPackageManager returns the detected package manager
func (e *LocalExecutor) GetPackageManager() installers.PackageManager {
	return e.pm
}

// IsRoot checks if running with elevated privileges
func (e *LocalExecutor) IsRoot() bool {
	switch e.os {
	case installers.OSWindows:
		// Check for admin on Windows
		_, err := exec.Command("net", "session").Output()
		return err == nil
	default:
		// Check for root on Unix
		output, _ := e.Run("id -u")
		return strings.TrimSpace(output) == "0"
	}
}

// Verify interface implementation
var _ installers.Executor = (*LocalExecutor)(nil)
