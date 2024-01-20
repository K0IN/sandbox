package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

/* this is the side INSIDE the namespace */

func CreateSandboxInsideNamespace(entryCommand, hostname, sandboxDir string) int {
	rootFsBasePath := path.Join(sandboxDir, SandboxRootFs)
	upperDirBasePath := path.Join(sandboxDir, SandboxUpperDir)
	workDirBasePath := path.Join(sandboxDir, SandboxWorkDir)

	if err := MakeOverlay("/", upperDirBasePath, rootFsBasePath, workDirBasePath); err != nil {
		panic(err)
	}
	defer UnmountOverlay(rootFsBasePath)

	if err := MountDevices(rootFsBasePath); err != nil {
		panic(err)
	}
	defer UnmountDevices(upperDirBasePath, rootFsBasePath)

	if err := MountProc(rootFsBasePath); err != nil {
		panic(err)
	}
	defer UnmountProc(rootFsBasePath)

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
		Chroot: rootFsBasePath,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return cmd.ProcessState.ExitCode()
}
