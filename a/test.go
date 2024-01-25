package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
)

// mount -t overlay overlay -o lowerdir=/,upperdir=/tmp/sandbox339789511/upper,workdir=/tmp/sandbox339789511/workdir,userxattr /tmp/sandbox339789511/rootfs

func mountOverlay(lowerDir, upperDir, mountDir, workDir string) error {
	opts := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s",
		lowerDir,
		upperDir,
		workDir,
	)

	return syscall.Mount("overlay", mountDir, "overlay", 0, opts)

}
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

	os.Chown(hostnamePath, 0, 0)

	return nil
}
func unmountOverlay(mountDir string) error {
	return syscall.Unmount(mountDir, 0)
}

func mountSpecial(rootFsBasePath string) {
	procPath := path.Join(rootFsBasePath, "proc")
	_ = os.MkdirAll(procPath, 0755)
	if err := syscall.Mount("proc", procPath, "proc", 0, ""); err != nil {
		panic(err)
	}

	sysPath := path.Join(rootFsBasePath, "sys")
	_ = os.MkdirAll(sysPath, 0755)
	println(sysPath)
	if err := syscall.Mount("sysfs", sysPath, "sysfs", 0, ""); err != nil {
		panic(err)
	}

	devPath := path.Join(rootFsBasePath, "dev")
	_ = os.MkdirAll(devPath, 0755)
	if err := syscall.Mount("devtmpfs", devPath, "devtmpfs", 0, ""); err != nil {
		panic(err)
	}

	ptsPath := path.Join(rootFsBasePath, "dev", "pts")
	_ = os.MkdirAll(ptsPath, 0755)
	if err := syscall.Mount("devpts", ptsPath, "devpts", 0, ""); err != nil {
		panic(err)
	}
}

func unmountSpecial(rootFsBasePath string) error {
	// unmount everything in reverse order
	ptsPath := path.Join(rootFsBasePath, "dev", "pts")
	if err := syscall.Unmount(ptsPath, 0); err != nil {
		return err
	}

	devPath := path.Join(rootFsBasePath, "dev")
	if err := syscall.Unmount(devPath, 0); err != nil {
		return err
	}

	sysPath := path.Join(rootFsBasePath, "sys")
	if err := syscall.Unmount(sysPath, 0); err != nil {
		return err
	}

	procPath := path.Join(rootFsBasePath, "proc")
	if err := syscall.Unmount(procPath, 0); err != nil {
		return err
	}

	os.RemoveAll(procPath)
	os.RemoveAll(sysPath)
	os.RemoveAll(devPath)
	os.RemoveAll(ptsPath)

	return nil
}

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

	if err := mountOverlay("/", upperDirBasePath, rootFsBasePath, workDirBasePath); err != nil {
		panic(err)
	}
	defer func() {
		err := unmountOverlay(rootFsBasePath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	mountSpecial(rootFsBasePath)

	defer func() {
		err := unmountSpecial(rootFsBasePath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	changeHostname(upperDirBasePath, "test")

	cmd := exec.Command("/bin/fish")
	cmd.Env = os.Environ()
	cmd.Dir = "/"
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
			Uid: 1000,
			Gid: 1000,
		},

		GidMappingsEnableSetgroups: true,
		Chroot:                     rootFsBasePath,
		//	Setpgid:                    true,
		//	Setsid:                     true,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Start()

	time.Sleep(1 * time.Second)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
