package main

import (
	"fmt"
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

	// cmd := exec.Command("fuse-overlayfs", "-t", "overlay", "overlay", "-o", "userxattr", "-o", "squash_to_root", "-o", opts, mountDir) //
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", opts, mountDir) //
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
	// cmd := exec.Command(os.Args[0], "test")
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
	// 	UidMappings: []syscall.SysProcIDMap{
	// 		{
	// 			ContainerID: 0,
	// 			HostID:      os.Getuid(),
	// 			Size:        1,
	// 		},
	// 		{
	// 			ContainerID: 1,
	// 			HostID:      1,
	// 			Size:        65534,
	// 		},
	// 	},
	// 	GidMappings: []syscall.SysProcIDMap{
	// 		{
	// 			ContainerID: 0,
	// 			HostID:      os.Getgid(),
	// 			Size:        1,
	// 		},
	// 		{
	// 			ContainerID: 1,
	// 			HostID:      1,
	// 			Size:        65534,
	// 		},
	// 	},
	// 	Credential: &syscall.Credential{
	// 		Uid: 0,
	// 		Gid: 0,
	// 	},
	// }

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
	cmd := exec.Command("find", "/", "-maxdepth", "1")
	allRootMountsRaw, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("findmnt", "--real", "-r", "-o", "target", "-n")
	allMountsRaw, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	allMounts := strings.Split(string(allMountsRaw), "\n")
	allRootMounts := strings.Split(string(allRootMountsRaw), "\n")
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

func createSandboxInsideNamespace() {
	rootFs := createRootFs()
	// defer os.RemoveAll(rootFs)

	println("Created rootfs at", rootFs)

	if err := mountDevices(rootFs); err != nil {
		panic(err)
	}

	if err := mountProc(rootFs); err != nil {
		panic(err)
	}

	if err := syscall.Sethostname([]byte("sandbox")); err != nil {
		panic(err)
	}

	// if err := syscall.Chroot(rootFs); err != nil {
	// 	panic(err)
	// }
	//
	// // list root dir
	// allFiles, err := os.ReadDir("/")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, file := range allFiles {
	// 	info, _ := file.Info()
	// 	println("Found file:", file.Name(), "is dir:", file.IsDir(), info.Mode().Perm())
	// }

	// check if /bin/sh exists
	//if _, err := os.Stat(path.Join(rootFs, "/bin/sh")); err != nil {
	//	panic(fmt.Errorf("error checking if /bin/sh exists: %v", err))
	//}

	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//	Chroot: rootFs,
	}
	cmd.Env = os.Environ()
	cmd.Dir = "/" // path.Join(rootFs, "/bin")
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
		createSandboxInsideNamespace()
	}
}
