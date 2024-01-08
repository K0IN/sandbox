package fuse-fs

//go:generate bash ./build.sh

import (
	_ "embed"
	"io/ioutil"
	"os/exec"
)

//go:embed fuse-overlayfs-bin
var fuseOverlayFS []byte

func main() {
	// write the file to disk
	ioutil.WriteFile("fuse-overlayfs-bin-test", fuseOverlayFS, 0755)
	res, err := exec.Command("./fuse-overlayfs-bin-test", "-h").CombinedOutput()
	if err != nil {
		panic(err)
	}
	println(string(res))
}
