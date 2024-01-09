package main

//go:generate bash ./build.sh

/*
#cgo CFLAGS: -I/fuse-overlayfs
#cgo LDFLAGS: fuse-overlayfs/fuse-overlayfs.a fuse-overlayfs/lib/libgnu.a -static -lfuse3 -ldl -pthread -lrt -lz

extern int initialize_fuse_overlayfs(const char *mountpoint, const char *lowerdir, const char *upperdir, const char *workdir);
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
)

// InitializeFuseOverlayFS Go wrapper for the C function
func InitializeFuseOverlayFS(mountpoint, lowerdir, upperdir, workdir string) error {
	cMountpoint := C.CString(mountpoint)
	cLowerdir := C.CString(lowerdir)
	cUpperdir := C.CString(upperdir)
	cWorkdir := C.CString(workdir)
	// defer C.free(unsafe.Pointer(cMountpoint))
	// defer C.free(unsafe.Pointer(cLowerdir))
	// defer C.free(unsafe.Pointer(cUpperdir))
	// defer C.free(unsafe.Pointer(cWorkdir))

	ret := C.initialize_fuse_overlayfs(cMountpoint, cLowerdir, cUpperdir, cWorkdir)
	if ret != 0 {
		return fmt.Errorf("fuse-overlayfs failed with code: %d", ret)
	}
	return nil
}

func main() {
	// Add your main application logic here, or call InitializeFuseOverlayFS as required.
	println("Initializing fuse-overlayfs from Go")
	os.Mkdir("/home/k0in/merged", os.ModePerm)
	os.Mkdir("/home/k0in/rootfs", os.ModePerm)
	os.Mkdir("/home/k0in/upper", os.ModePerm)
	os.Mkdir("/home/k0in/work", os.ModePerm)

	err := InitializeFuseOverlayFS("/home/k0in/merged", "/home/k0in/rootfs", "/home/k0in/upper", "/home/k0in/work")
	if err != nil {
		panic(err)
	}
	println("fuse-overlayfs initialized")

	// wait for input
	var input string
	fmt.Scanln(&input)
	println("fuse-overlayfs unmounting")

	// unmount
	syscall.Unmount("/home/k0in/merged", 0)
}
