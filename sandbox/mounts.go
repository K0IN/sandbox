package sandbox

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

type SpecialMount struct {
	rootFsPath string

	MountProc bool
	MountDev  bool

	DevicesToMount []string
}

func CreateSpecialMounts(rootFsPath string) (*SpecialMount, error) {
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
		rootFsPath: rootFsPath,
	}, nil
}

func (s *SpecialMount) Mount() error {
	if s.MountProc {
		if err := MountProc(s.rootFsPath); err != nil {
			return err
		}
	}

	if s.MountDev {
		if err := mountDevices(s.rootFsPath, s.DevicesToMount); err != nil {
			return err
		}
	}

	return nil
}

func (s *SpecialMount) Unmount() error {
	if s.MountProc {
		if err := UnmountProc(s.rootFsPath); err != nil {
			return fmt.Errorf("failed to unmount proc: %w", err)
		}
	}

	if s.MountDev {
		if err := UnmountDevices(s.rootFsPath, s.DevicesToMount); err != nil {
			return fmt.Errorf("failed to unmount dev: %w", err)
		}
	}

	return nil
}

func mountDevices(rootFsPath string, devicesToMount []string) error {
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

		if err := syscall.Mount(path.Join("/dev", deviceName), dst, "", syscall.MS_BIND, ""); err != nil {
			continue
		}
	}

	ptsPath := path.Join(rootFsPath, "dev", "pts")
	if err := os.MkdirAll(ptsPath, 0755); err != nil {
		return err
	}

	if err := syscall.Mount("devpts", ptsPath, "devpts", 0, ""); err != nil {
		return err
	}

	return nil
}

func UnmountDevices(rootFsPath string, devicesToUnmount []string) error {
	if err := syscall.Unmount(path.Join(rootFsPath, "dev", "pts"), 0); err != nil {
		return err
	}

	for _, deviceName := range devicesToUnmount {
		if err := syscall.Unmount(path.Join(rootFsPath, "dev", deviceName), 0); err != nil {
			return err
		}
	}

	return os.RemoveAll(path.Join(rootFsPath, "dev"))
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
	_ = syscall.Unmount(procPath, 0)
	return os.RemoveAll(procPath)
}
