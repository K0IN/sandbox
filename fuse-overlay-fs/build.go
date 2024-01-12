package fuse_overlay_fs

//go:generate bash ./build.sh

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

//go:embed fuse-overlayfs-bin
var FuseOverlayFSBin []byte

type OverlayFs struct {
	upperDir   string
	workDir    string
	MergedDir  string
	binaryPath string
}



func CreateOverlayFs(tmpFolder string) (*OverlayFs, error) {
	o := &OverlayFs{}
	tmpDir, err := ioutil.TempDir("", "fuse-overlayfs")
	if err != nil {
		panic(err)
	}
	tmpBinPath := path.Join(tmpDir, "fuse-overlayfs-bin")
	err = ioutil.WriteFile(tmpBinPath, FuseOverlayFSBin, 0755)
	if err != nil {
		panic(err)
	}
	o.binaryPath = tmpBinPath

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

	mounts := fmt.Sprintf("uidmapping=0:10:100:100:10000:2000,gidmapping=0:10:100:100:10000:2000,lowerdir=%s,upperdir=%s,workdir=%s", "/", o.upperDir, o.workDir)
	result, err := exec.Command(o.binaryPath, "-o", mounts, o.MergedDir).CombinedOutput()
	if err != nil {
		return err
	}
	println(string(result))
	return nil
}

func (o *OverlayFs) Unmount() error {
	exec.Command("fusermount", "-u", "/home/k0in/merged").Run()
	return nil
}

// func showFile(upperDir, rootFsPath string) {
// 	dmp := diffmatchpatch.New()
//
// 	upperFilePath := path.Join(upperDir, rootFsPath)
// 	rootFsFilePath := path.Join("/", rootFsPath)
//
// 	text1, err := os.ReadFile(upperFilePath)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	text2, err := os.ReadFile(rootFsFilePath)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	diffs := dmp.DiffMain(string(text1), string(text2), false)
// 	fmt.Println(dmp.DiffPrettyText(diffs))
// }

func (o *OverlayFs) ShowDiff() {
	// showFile(o.upperDir, "/etc/hosts")
	// showFile(o.upperDir, "/etc/resolv.conf")
	// showFile(o.upperDir, "/etc/hostname")
}

func (o *OverlayFs) GetRootDir() string {
	return o.MergedDir
}
