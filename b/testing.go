package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func main() {
	base, err := os.MkdirTemp("", "sandbox*")
	if err != nil {
		panic(err)
	}
	sandboxDir := base
	SandboxRootFs := "rootfs"
	SandboxUpperDir := "upper"
	SandboxWorkDir := "workdir"

	println(sandboxDir)

	rootFsBasePath := path.Join(sandboxDir, SandboxRootFs)
	upperDirBasePath := path.Join(sandboxDir, SandboxUpperDir)
	workDirBasePath := path.Join(sandboxDir, SandboxWorkDir)

	if err := os.MkdirAll(rootFsBasePath, 0755); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(upperDirBasePath, 0755); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(workDirBasePath, 0755); err != nil {
		panic(err)
	}

	if err := syscall.Mount("overlay", rootFsBasePath, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", "/", upperDirBasePath, workDirBasePath)); err != nil {
		panic(err)
	}

	defer syscall.Unmount(rootFsBasePath, 0)

	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: rootFsBasePath,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()

	// pid
	fmt.Println(cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
