package sandbox

import (
	"os"
	"os/exec"
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
