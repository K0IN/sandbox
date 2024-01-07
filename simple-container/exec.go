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

func prepareCommand(command string, execConfig *ExecConfig) string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}

	return fmt.Sprintf("hostname \"%s-%s\" && cd \"%s\" && %s", hostName, "sandbox", execConfig.WorkDir, command)

}

func ExecuteCommand(command string, execConfig *ExecConfig) error {
	args := prepareCommand(command, execConfig)
	// println("Executing command:", args)
	cmd := exec.Command("/bin/sh", "-c", args)
	cmd.Env = execConfig.Env
	cmd.Dir = "/"
	attr := syscall.SysProcAttr{
		Chroot:     execConfig.Rootfs,
		Cloneflags: execConfig.NameSpaceFlags,
		// Unshareflags: ,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      1000,
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

// https://blog.scottlowe.org/2013/09/04/introducing-linux-network-namespaces/ and https://lwn.net/Articles/580893/
func AddNetwork(name string, ip string) error {
	// https://github.com/vishvananda/netlink

	// ip link add veth0 type veth peer name veth1
	// ip link set veth1 netns <pid>
	// ip netns exec <pid> ip addr add
	// ip netns exec <pid> ip link set veth1 up
	// ip netns exec <pid> ip route add default via <ip>
	return nil
}

func AddMount() error {
	//sudo mount -o ro,noload /dev/sda1 /media/2tb
	//  syscall.Mount("/dev/sda1", "/media/2tb", "ext4", syscall.MS_RDONLY, "")
	return nil
}
