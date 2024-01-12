package main

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

// a go sandbox that can run code on **your** machine without changing actual files
// we use namespaces to isolate a new bash process
// to protect your disk we mount a overlayfs on top of your rootfs
// then we chroot into the new rootfs
// after you exit the bash process we unmount the overlayfs show the diff and delete the overlayfs
// the sandbox is saved to /tmp/sandbox and is deleted after you exit the bash process
// inside the sandbox you are root and can do whatever you want
// the sandbox is not persistent and is deleted after you exit the bash process

func forkSelfIntoNewNamespace() {
	cmd := exec.Command(os.Args[0], "fork")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
		// we need to set the uid and gid to 0 to be root inside the new namespace
		// we also need to set the uid and gid to 0 to be able to mount the overlayfs
		UidMappings: []syscall.SysProcIDMap{

			// map all users from host
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{

			// map all groups from host
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
		// we need to set the gid to 0 to be able to mount the overlayfs
		// GidMappingsEnableSetgroups: false,
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		AmbientCaps: []uintptr{
			// we need CAP_SYS_ADMIN to mount the overlayfs as number
			37,
			// we  need CAP_SYS_CHROOT to chroot as number
			19,
		},
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

}

func createSandboxInsideNamespace() {
	sandboxDir := "/tmp/sandbox4"

	// we are inside the new namespace
	// we create a new directory for the sandbox
	// if the /tmp/sandbox directory already exists we delete it
	if _, err := os.Stat(sandboxDir); err == nil {
		err = os.RemoveAll(sandboxDir)
		if err != nil {
			panic(err)
		}
	}

	upperDir := path.Join(sandboxDir, "upper")
	mountDir := path.Join(sandboxDir, "rootfs")

	// we create the directories for the overlayfs
	if err := os.MkdirAll(upperDir, 0755); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(mountDir, 0755); err != nil {
		panic(err)
	}

	// we mount the overlayfs
	mntCmd := exec.Command("fuse-overlayfs", "overlay", mountDir, "-o", "lowerdir=/,upperdir="+upperDir+",workdir="+upperDir+",squash_to_root")
	result, err := mntCmd.CombinedOutput()
	println(string(result))

	if err != nil {
		panic(err)
	}

	defer func() {
		cmd := exec.Command("fusermount", "-u", mountDir)
		result, _ := cmd.CombinedOutput()
		println(string(result))

	}()

	// we chroot into the new rootfs
	err = syscall.Chroot(mountDir)
	if err != nil {
		panic(err)
	}

	// we change the working directory to /
	err = os.Chdir("/")
	if err != nil {
		panic(err)
	}

	// we set the hostname to sandbox
	err = syscall.Sethostname([]byte("sandbox"))
	if err != nil {
		panic(err)
	}

	// now we can run bash
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(0)
}

func main() {
	println("Starting sandbox as user", os.Getuid(), "and group", os.Getgid())
	// we fork ourselves into a new namespace
	if len(os.Args) == 1 {
		forkSelfIntoNewNamespace()
	} else {
		// we are inside the new namespace
		createSandboxInsideNamespace()
	}
}
