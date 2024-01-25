package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func MakeOverlay(lowerDir, upperDir, mountDir, workDir string) error {
	opts := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s,userxattr",
		lowerDir,
		upperDir,
		workDir,
	)
	return syscall.Mount("overlay", mountDir, "overlay", 0, opts)
}

func UnmountOverlay(mountDir string) error {
	return syscall.Unmount(mountDir, 0)
}

func MountDevices(rootFsPath string) error {
	devicesToMount := []string{
		"tty",
		"null",
		"zero",
		"full",
		"random",
		"urandom",
	}

	devPath := path.Join(rootFsPath, "dev")

	_ = os.RemoveAll(devPath)
	if err := os.MkdirAll(devPath, 0755); err != nil {
		return err
	}

	for _, deviceName := range devicesToMount {
		dst := path.Join(devPath, deviceName)

		_ = os.RemoveAll(dst)
		if _, err := os.Create(dst); err != nil {
			continue
		}

		cmd := exec.Command("mount", "-o", "bind", path.Join("/dev", deviceName), dst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			continue
		}
	}

	// mount /dev/pts
	// mount -t devpts pts /mnt/linux/dev/pts

	/// we create the mount point
	ptsPath := path.Join(rootFsPath, "dev", "pts")
	if err := os.MkdirAll(ptsPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("mount", "-t", "devpts", "devpts", ptsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func UnmountDevices(upperFs, rootFsPath string) error {
	_ = syscall.Unmount(path.Join(rootFsPath, "dev"), 0)
	return os.RemoveAll(path.Join(upperFs, "dev"))
}

func MountProc(rootFsPath string) error {
	procPath := path.Join(rootFsPath, "proc")
	if err := os.MkdirAll(procPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("mount", "-t", "proc", "proc", procPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func UnmountProc(upperFs string) error {
	_ = syscall.Unmount(path.Join(upperFs, "proc"), 0)
	return os.Remove(path.Join(upperFs, "proc"))
}
