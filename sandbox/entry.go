package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

/* this is the side INSIDE the namespace */

func changeHostname(upperDir, hostname string) error {
	// write to /etc/hostname
	hostnamePath := path.Join(upperDir, "etc", "hostname")
	if err := os.MkdirAll(path.Dir(hostnamePath), 0755); err != nil {
		return err
	}

	f, err := os.Create(hostnamePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(hostname); err != nil {
		return err
	}

	return nil
}

func CreateSandboxInsideNamespace(entryCommand, hostname, sandboxDir string) int {
	rootFsBasePath := path.Join(sandboxDir, SandboxMountPointDir)
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

	if err := changeHostname(upperDirBasePath, hostname); err != nil {
		panic(err)
	}

	currentWorkingDir := "/"
	if workDir, err := os.Getwd(); err == nil {
		currentWorkingDir = workDir
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, entryCommand))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65536,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65536,
			},
		},
		Credential: &syscall.Credential{
			Uid: uint32(os.Getuid()), // todo: get uid from user
			Gid: uint32(os.Getgid()), // todo: get gid from user
		},

		GidMappingsEnableSetgroups: true, // enable su command
		Chroot:                     rootFsBasePath,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return cmd.ProcessState.ExitCode()
}
