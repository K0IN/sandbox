package fuse_overlay_fs

//go:generate bash ./build.sh

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed fuse-overlayfs-bin
var fuseOverlayFSBin []byte
var execPath *string

func GetExecPath() (string, error) {
	if execPath == nil {
		if installedCommandPath, err := findInstalledOverlayFsBin(); err == nil {
			execPath = &installedCommandPath
			return installedCommandPath, nil
		}

		if tmpPath, err := unpackFuseOverlayFSBin(); err == nil {
			execPath = &tmpPath
			return tmpPath, nil
		}

		return "", fmt.Errorf("could not find or unpack fuse-overlayfs binary")
	}
	return *execPath, nil
}

func findInstalledOverlayFsBin() (string, error) {
	return exec.LookPath("fuse-overlayfs")
}

func unpackFuseOverlayFSBin() (string, error) {
	tmpFile, err := os.CreateTemp("", "fuse-overlayfs-bin")
	if err != nil {
		return "", err
	}
	if _, err := tmpFile.Write(fuseOverlayFSBin); err != nil {
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	tmpBinPath := tmpFile.Name()
	if err := os.Chmod(tmpBinPath, 0755); err != nil {
		return "", err
	}
	return tmpBinPath, nil
}
