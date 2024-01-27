package helper

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type MountInfoHeader struct {
	FileSystems []MountInfo `json:"filesystems"`
}

type MountInfo struct {
	Source   string       `json:"source"`
	Target   string       `json:"target"`
	Options  string       `json:"options"`
	Type     string       `json:"fstype"`
	Children *[]MountInfo `json:"children,omitempty"`
}

func FindAllMounts() ([]MountInfo, error) {
	//  findmnt -J --real -n --all
	_, err := exec.LookPath("findmnt")
	if err != nil {
		return nil, fmt.Errorf("command findmnt not found: %w", err)
	}

	cmd := exec.Command("findmnt", "-J", "--real", "-n", "--all")
	outRaw, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var mounts MountInfoHeader
	if err := json.Unmarshal(outRaw, &mounts); err != nil {
		return nil, err
	}

	return mounts.FileSystems, nil
}
