package fuse_overlay_fs

//go:generate bash ./build.sh

import (
	_ "embed"
	"io/ioutil"
	"os/exec"
	"path"
)

//go:embed fuse-overlayfs-bin
var fuseOverlayFSBin []byte
var execPath *string

func GetExecPath() (string, error) {
	// maybe it is already known
	// find the "fuse-overlayfs" binary

	if execPath == nil {
		installedCommandPath, err := exec.LookPath("fuse-overlayfs")
		if err == nil {
			execPath = &installedCommandPath
			return installedCommandPath, nil
		}

		tmpPath, err := unpackFuseOverlayFSBin()
		if err == nil {
			execPath = &tmpPath
			return tmpPath, nil
		}

		return "", err
	}
	return *execPath, nil
}

func unpackFuseOverlayFSBin() (string, error) {
	tmpDir, err := ioutil.TempDir("", "fuse-overlayfs")
	if err != nil {
		return "", err
	}
	tmpBinPath := path.Join(tmpDir, "fuse-overlayfs-bin")
	err = ioutil.WriteFile(tmpBinPath, fuseOverlayFSBin, 0755)
	if err != nil {
		return "", err
	}
	return tmpBinPath, nil
}
