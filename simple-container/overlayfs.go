package container

import (
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type OverlayFs struct {
	upperDir  string
	workDir   string
	MergedDir string
}

func CreateOverlayFs(tmpFolder string) (*OverlayFs, error) {
	o := &OverlayFs{}
	o.MergedDir = path.Join(tmpFolder, "merged")
	o.upperDir = path.Join(tmpFolder, "upper")
	o.workDir = path.Join(tmpFolder, "work")
	return o, nil
}

func (o *OverlayFs) Mount() error {
	err := os.MkdirAll(o.upperDir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(o.workDir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(o.MergedDir, 0755)
	if err != nil {
		panic(err)
	}

	mountArgs := fmt.Sprintf("lowerdir=/,upperdir=%s,workdir=%s", o.upperDir, o.workDir)
	return syscall.Mount("overlay", o.MergedDir, "overlay", 0, mountArgs)
}

func (o *OverlayFs) Unmount() error {
	return syscall.Unmount(o.MergedDir, syscall.MNT_FORCE)
}

func showFile(upperDir, rootFsPath string) {
	dmp := diffmatchpatch.New()

	upperFilePath := path.Join(upperDir, rootFsPath)
	rootFsFilePath := path.Join("/", rootFsPath)

	text1, err := os.ReadFile(upperFilePath)
	if err != nil {
		panic(err)
	}

	text2, err := os.ReadFile(rootFsFilePath)
	if err != nil {
		panic(err)
	}

	diffs := dmp.DiffMain(string(text1), string(text2), false)
	fmt.Println(dmp.DiffPrettyText(diffs))
}

func (o *OverlayFs) ShowDiff() {
	showFile(o.upperDir, "/etc/hosts")
	showFile(o.upperDir, "/etc/resolv.conf")
	showFile(o.upperDir, "/etc/hostname")
}

func (o *OverlayFs) GetRootDir() string {
	return o.MergedDir
}
