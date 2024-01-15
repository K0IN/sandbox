package sandbox

import (
	"fmt"
	fuse_overlay_fs "myapp/fuse-overlay-fs"
	"os"
	"os/exec"
	"path"
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

func CreateSandboxDirectories() (string, string, string, error) {
	sandboxDir, err := os.MkdirTemp("", "sandbox")
	if err != nil {
		return "", "", "", err
	}

	// lets first create all the directories we need
	rootFsBasePath := path.Join(sandboxDir, "rootfs")
	upperDirBasePath := path.Join(sandboxDir, "upperdir")
	workDirBasePath := path.Join(sandboxDir, "workdir")

	if err := os.MkdirAll(rootFsBasePath, 0755); err != nil {
		return "", "", "", err
	}
	if err := os.MkdirAll(upperDirBasePath, 0755); err != nil {
		return "", "", "", err
	}
	if err := os.MkdirAll(workDirBasePath, 0755); err != nil {
		return "", "", "", err
	}

	if err := MakeOverlay("/", upperDirBasePath, rootFsBasePath, workDirBasePath); err != nil {
		return "", "", "", err
	}

	return sandboxDir, rootFsBasePath, upperDirBasePath, nil
}