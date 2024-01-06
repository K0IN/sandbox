package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type ExecConfig struct {
	Env            []string
	WorkDir        string
	Rootfs         string
	NameSpaceFlags uintptr
}

func ExecuteCommand(command string, execConfig *ExecConfig) error {
	args := fmt.Sprintf("cd %s && %s", execConfig.WorkDir, command)
	// println("Executing command:", args)
	cmd := exec.Command("/bin/sh", "-c", args)
	cmd.Env = execConfig.Env
	cmd.Dir = "/"
	attr := syscall.SysProcAttr{
		Chroot:     execConfig.Rootfs,
		Cloneflags: execConfig.NameSpaceFlags,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getgid(),
				Size:        1,
			},
		},
	}

	// if err := syscall.Mount("proc", chroot+"/proc", "proc", 0, ""); err != nil {
	// 	return fmt.Errorf("failed to mount proc: %w", err)
	// }
	// defer syscall.Unmount(chroot+"/proc", 0)

	cmd.SysProcAttr = &attr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}
