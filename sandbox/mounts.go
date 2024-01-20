package sandbox

import (
	"fmt"
	fuse_overlay_fs "myapp/fuse-overlay-fs"
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

	if err := os.MkdirAll(mountDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(upperDir, 0755); err != nil {
		return err
	}

	fuseFsBin, err := fuse_overlay_fs.GetExecPath()
	if err != nil {
		return err
	}
	// i would really like to use the linux built in overlayfs, but i can't get it to work
	cmd := exec.Command(fuseFsBin, "-t", "overlay", "overlay", "-o", "userxattr", "-o", "squash_to_root", "-o", opts, mountDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error making overlay for opts %s %s: %v, output: %s", opts, mountDir, cmd.ProcessState.ExitCode(), output)
	}
	return nil
}

func UnmountOverlay(mountDir string) error {
	cmd := exec.Command("fusermount", "-u", mountDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error unmounting overlay %s: %v, output: %s", mountDir, cmd.ProcessState.ExitCode(), output)
	}
	return nil
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

type SandboxDirectories struct {
	SandboxDir       string
	RootFsBasePath   string
	UpperDirBasePath string
}

func CreateSandboxDirectories(sandboxDir string) (*SandboxDirectories, error) {
	err := os.MkdirAll(sandboxDir, 0755) // ensure sandbox dir exists
	if err != nil {
		return nil, err
	}

	// lets first create all the directories we need
	rootFsBasePath := path.Join(sandboxDir, SandboxRootFs)
	upperDirBasePath := path.Join(sandboxDir, SandboxUpperDir)
	workDirBasePath := path.Join(sandboxDir, SandboxWorkDir)

	if err := os.MkdirAll(rootFsBasePath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(upperDirBasePath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(workDirBasePath, 0755); err != nil {
		return nil, err
	}

	return &SandboxDirectories{
		SandboxDir:       sandboxDir,
		RootFsBasePath:   rootFsBasePath,
		UpperDirBasePath: upperDirBasePath,
	}, nil
}
