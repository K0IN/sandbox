package main

import (
	"fmt"
	sandbox "myapp/sandbox"
	"os"
	"os/exec"
	"syscall"
)

func setSandboxHostname() error {
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

func createSandboxInsideNamespace(entryCommand string) {
	sandboxDir, rootFs, _, err := sandbox.CreateSandboxDirectories()
	if err != nil {
		panic(err)
	}

	println("Created sandbox directories", sandboxDir)

	if err := sandbox.MountDevices(rootFs); err != nil {
		panic(err)
	}

	if err := sandbox.MountProc(rootFs); err != nil {
		panic(err)
	}

	_ = setSandboxHostname()

	// current dir
	currentWorkingDir := "/"
	if workDir, err := os.Getwd(); err == nil {
		currentWorkingDir = workDir
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, entryCommand))
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: rootFs,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		sandbox.ForkSelfIntoNewNamespace(os.Args) // this will call us again with an argument
	} else {
		if os.Getuid() != 0 || os.Getgid() != 0 {
			panic("started in namespace mode but not as root")
		}
		// we are inside the new namespace
		shell := sandbox.GetPrimaryShell()
		println("Starting sandbox with shell", shell)
		createSandboxInsideNamespace(shell)
	}
}
