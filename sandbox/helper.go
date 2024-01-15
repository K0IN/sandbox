package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func GetPrimaryShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return shell
	}
	if path, err := exec.LookPath("bash"); err == nil {
		return path
	}
	return "/bin/sh"
}

func SetSandboxHostname() error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	newHostname := fmt.Sprintf("sandbox@%s", hostname)

	if err := syscall.Sethostname([]byte(newHostname)); err != nil {
		return err
	}
	return nil
}
