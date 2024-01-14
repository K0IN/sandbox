package main

import (
	"fmt"
	fuse_overlay_fs "myapp/fuse-overlay-fs"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"syscall"
)

func makeOverlay(lowerDir, upperDir, mountDir, workDir string) error {
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

func mountDevices(rootFsPath string) error {
	devicesToMount := []string{
		"tty",
		"null",
		"zero",
		"full",
		"random",
		"urandom",
	}
	devPath := path.Join(rootFsPath, "dev")
	if err := os.MkdirAll(devPath, 0755); err != nil {
		return err
	}

	for _, deviceName := range devicesToMount {
		dst := path.Join(devPath, deviceName)
		if _, err := os.Create(dst); err != nil {
			//return err
			continue
		}

		cmd := exec.Command("mount", "-o", "bind", "/dev/"+deviceName, dst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			// return err
			continue
		}
	}
	return nil
}

func mountProc(rootFsPath string) error {
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

func forkSelfIntoNewNamespace() {
	cmd := exec.Command("unshare", "--mount", "--user", "--map-root-user", "--pid", "--fork", "--uts", os.Args[0], "test")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	os.Exit(cmd.ProcessState.ExitCode())
}

func listAllMounts() []string {
	allRootFiles, err := os.ReadDir("/")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("findmnt", "--real", "-r", "-o", "target", "-n")
	allMountsRaw, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	allMounts := strings.Split(string(allMountsRaw), "\n")
	allRootMounts := make([]string, 0, len(allRootFiles))
	for _, rootFile := range allRootFiles {
		allRootMounts = append(allRootMounts, "/"+rootFile.Name())
	}
	allMounts = append(allMounts, allRootMounts...)

	uniqueMountPoints := make(map[string]bool)
	for _, mp := range allMounts {
		uniqueMountPoints[mp] = true
	}

	result := make([]string, 0, len(uniqueMountPoints))
	for mountPoint := range uniqueMountPoints {
		result = append(result, strings.TrimSpace(mountPoint))
	}

	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) <= len(result[j])
	})

	return result
}

func createRootFs() string {
	sandboxDir := "/tmp/sandbox"

	// remove the sandbox dir if it exists
	if err := os.RemoveAll(sandboxDir); err != nil {
		//panic(err)
	}

	// lets first create all the directories we need
	rootFsBasePath := path.Join(sandboxDir, "rootfs")
	upperDirBasePath := path.Join(sandboxDir, "upperdir")
	workDirBasePath := path.Join(sandboxDir, "workdir")

	if err := os.MkdirAll(rootFsBasePath, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(upperDirBasePath, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(workDirBasePath, 0755); err != nil {
		panic(err)
	}

	allMountPoints := listAllMounts()
	for _, mountPoint := range allMountPoints {
		if mountPoint == "/" || mountPoint == "/dev" || mountPoint == "/proc" || mountPoint == "" {
			continue // Skip these special mount points
		}
		// check if the mountpoint is a directory
		fileInfo, err := os.Stat(mountPoint)
		if err != nil || !fileInfo.IsDir() {
			continue
		}

		//println("Creating overlayfs for mountpoint:", mountPoint)
		rootFsPath := path.Join(rootFsBasePath, mountPoint)

		upperDir := path.Join(upperDirBasePath, mountPoint)
		workDir := path.Join(workDirBasePath, mountPoint)

		// now we can create the overlayfs
		if err := makeOverlay(mountPoint, upperDir, rootFsPath, workDir); err != nil {
			fmt.Printf("Error creating overlayfs for mountpoint %s: %v\n", mountPoint, err)
			continue
		}
	}

	return rootFsBasePath
}

func createSandboxInsideNamespace(entryCommand string) {
	rootFs := createRootFs()
	// defer os.RemoveAll(rootFs)

	println("Created rootfs at", rootFs)

	if err := mountDevices(rootFs); err != nil {
		panic(err)
	}

	if err := mountProc(rootFs); err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	newHostname := fmt.Sprintf("sandbox@%s", hostname)

	if err := syscall.Sethostname([]byte(newHostname)); err != nil {
		panic(err)
	}

	// current dir
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		currentWorkingDir = "/"
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, entryCommand))
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: rootFs,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	println("Starting sandbox as user", os.Getuid(), "and group", os.Getgid())
	// we fork ourselves into a new namespace
	if len(os.Args) == 1 {
		forkSelfIntoNewNamespace()
	} else {
		if os.Getuid() != 0 || os.Getgid() != 0 {
			panic("we need to be root to mount the overlayfs")
		}
		// we are inside the new namespace
		createSandboxInsideNamespace("/bin/fish")
	}
}
