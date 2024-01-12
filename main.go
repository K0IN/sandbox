package main

import (
	"fmt"
	"io/ioutil"
	fuse_overlay_fs "myapp/fuse-overlay-fs"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func Mount(merge, upper, lower string) error {
	tmpDir, err := ioutil.TempDir("", "fuse-overlayfs")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	// defer os.RemoveAll(tmpDir)
	tmpBinPath := path.Join(tmpDir, "fuse-overlayfs-bin")
	err = ioutil.WriteFile(tmpBinPath, fuse_overlay_fs.FuseOverlayFSBin, 0755)
	if err != nil {
		return fmt.Errorf("failed to write fuse-overlayfs-bin: %w", err)
	}
	// uidmapping=0:0:1:1000:1000:1,gidmapping=0:0:1:1000:1000:1,
	mounts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lower, upper, merge)
	userMapping := "uidmapping=0:0:1:1000:1000:1,gidmapping=0:0:1:1000:1000:1"
	result, err := exec.Command(tmpBinPath, "-o", userMapping+","+mounts, merge).CombinedOutput()
	println(string(result))
	if err != nil {
		return fmt.Errorf("failed to mount overlayfs: %w", err)
	}
	return nil
}

func Unmount(merge string) error {
	return exec.Command("fusermount3", "-u", merge).Run()
}

func main() {
	println("you are ", os.Getgid(), os.Getuid())
	if len(os.Args) == 1 {
		// we fork ourselves with a new namespace
		// we need to be root to do that
		cmd := exec.Command(os.Args[0], "you are root")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
			Foreground: true,
			Noctty:     false,
			Credential: &syscall.Credential{Uid: 0, Gid: 0},
			UidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: os.Getgid(), Size: 1},
			},
			GidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: os.Getgid(), Size: 1},
			},
		}
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			panic(err)
		}

		if err := cmd.Wait(); err != nil {
			panic(err)
		}

	} else {
		sandboxDir, err := os.MkdirTemp("", "sandbox")
		if err != nil {
			panic(err)
		}

		defer println("cleanup", sandboxDir)
		defer os.RemoveAll(sandboxDir)

		upperDir := path.Join(sandboxDir, "upperdir")
		mergedDir := path.Join(sandboxDir, "merged")

		os.MkdirAll(upperDir, 0777)
		os.MkdirAll(mergedDir, 0777)
		// os.Chown(upperDir, os.Getuid(), os.Getgid())
		// os.Chown(mergedDir, os.Getuid(), os.Getgid())

		// err := mountDevices(sandboxDir, []string{"tty", "null", "zero", "full", "random", "urandom"})

		if err := Mount(mergedDir, upperDir, "/"); err != nil {
			panic(err)
		}
		defer Unmount(mergedDir)

		if err := syscall.Mount("tmpfs", path.Join(sandboxDir, "merged", "dev"), "tmpfs", 0, ""); err != nil {
			panic(err)
		}
		defer syscall.Unmount(path.Join(sandboxDir, "merged", "dev"), syscall.MNT_FORCE)

		if err := syscall.Mount("proc", path.Join(sandboxDir, "merged", "proc"), "proc", 0, ""); err != nil {
			panic(fmt.Errorf("cannot mount proc: %w", err))
		}
		defer syscall.Unmount(path.Join(sandboxDir, "merged", "proc"), syscall.MNT_FORCE)

		println("your new rootfs path is: ", mergedDir)
		cmd := exec.Command("/bin/sh", "-c", "whoami ; id ; /bin/bash")
		cmd.Dir = "/"
		err = syscall.Chroot(mergedDir)

		if err != nil {
			panic(fmt.Errorf("cannot chroot: %w", err))
		}

		err = syscall.Sethostname([]byte("sandbox"))
		if err != nil {
			fmt.Printf("cannot set hostname: %v\n", err)
		}
		// cmd.SysProcAttr = &syscall.SysProcAttr{
		// 	Chroot: mergedDir,
		// }
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			panic(fmt.Errorf("cannot start command: %w", err))
		}

		if err := cmd.Wait(); err != nil {
			panic(fmt.Errorf("cannot wait command: %w", err))
		}

		println("done")
		os.Exit(0)
	}
}
