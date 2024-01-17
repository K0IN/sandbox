package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func CreateSandboxInsideNamespace(entryCommand, hostname string) int {
	sandboxPaths, err := CreateSandboxDirectories()
	if err != nil {
		panic(err)
	}

	println("Created sandbox directories", sandboxPaths.SandboxDir)

	if err := MountDevices(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}

	if err := MountProc(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}

	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fmt.Println("Failed to set hostname", err)
	}

	// current dir
	currentWorkingDir := "/"
	if workDir, err := os.Getwd(); err == nil {
		currentWorkingDir = workDir
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, entryCommand))
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: sandboxPaths.RootFsBasePath,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return cmd.ProcessState.ExitCode()
}
