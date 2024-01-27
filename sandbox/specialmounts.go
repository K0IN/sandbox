package sandbox

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

type SpecialMount struct {
	mountPath string

	MountProc bool
	MountDev  bool

	DevicesToMount []string
}

func CreateSpecialMounts(mountPath string) (*SpecialMount, error) {
	return &SpecialMount{
		MountProc: true,
		MountDev:  true,
		DevicesToMount: []string{
			"tty",
			"null",
			"zero",
			"full",
			"random",
			"urandom",
		},
		mountPath: mountPath,
	}, nil
}

func (s *SpecialMount) Mount() error {
	if s.MountProc {
		if err := MountProc(s.mountPath); err != nil {
			return err
		}
	}

	if s.MountDev {
		if err := mountDevices(s.mountPath, s.DevicesToMount); err != nil {
			return err
		}
	}

	return nil
}

func (s *SpecialMount) Unmount() error {
	if s.MountProc {
		if err := UnmountProc(s.mountPath); err != nil {
			return fmt.Errorf("failed to unmount proc: %w", err)
		}
	}

	if s.MountDev {
		if err := UnmountDevices(s.mountPath, s.DevicesToMount); err != nil {
			return fmt.Errorf("failed to unmount dev: %w", err)
		}
	}
	return nil
}

func mountDevices(rootFsPath string, devicesToMount []string) error {
	devPath := path.Join(rootFsPath, "dev")

	_ = os.RemoveAll(devPath)
	if err := os.MkdirAll(devPath, 0666); err != nil {
		return err
	}

	if err := syscall.Mount("devtmpfs", devPath, "devtmpfs", 0, ""); err != nil {
		return err
	}

	ptsPath := path.Join(rootFsPath, "dev", "pts")
	if err := os.MkdirAll(ptsPath, 0666); err != nil {
		return err
	}

	if err := syscall.Mount("devpts", ptsPath, "devpts", 0, ""); err != nil {
		return err
	}

	return nil
}

func UnmountDevices(rootFsPath string, devicesToUnmount []string) error {
	if err := syscall.Unmount(path.Join(rootFsPath, "dev", "pts"), 0); err != nil {
		fmt.Printf("failed to unmount devpts: %s\n", err)
	}

	if err := syscall.Unmount(path.Join(rootFsPath, "dev"), 0); err != nil {
		fmt.Printf("failed to unmount devtmpfs: %s\n", err)
	}

	if err := os.RemoveAll(path.Join(rootFsPath, "dev")); err != nil {
		fmt.Printf("failed to remove dev: %s\n", err)
	}
	return nil
}

func MountProc(rootFsPath string) error {
	procPath := path.Join(rootFsPath, "proc")
	if err := os.MkdirAll(procPath, 0755); err != nil {
		return err
	}
	return syscall.Mount("proc", procPath, "proc", 0, "")
}

func UnmountProc(rootFsPath string) error {
	procPath := path.Join(rootFsPath, "proc")
	err := syscall.Unmount(procPath, 0)
	if err != nil {
		fmt.Printf("failed to unmount proc: %s\n", err)
	}

	if err := os.RemoveAll(procPath); err != nil {
		fmt.Printf("failed to remove proc: %s\n", err)
	}

	return nil
}
