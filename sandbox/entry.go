package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

/* this is the side INSIDE the namespace */

func CreateSandboxInsideNamespace(entryCommand, hostname, hostPath string) int {
	sandboxPaths, err := CreateSandboxDirectories(hostPath)
	if err != nil {
		panic(err)
	}

	if err := MountDevices(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}
	defer UnmountDevices(sandboxPaths.RootFsBasePath)

	if err := MountProc(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}
	defer UnmountProc(sandboxPaths.RootFsBasePath)

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
