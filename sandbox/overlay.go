package sandbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"syscall"
)

const (
	sandboxConfigFileName = "config.json"
	sandboxUpperDir       = "upper"
	sandboxWorkDir        = "workdir"
	sandboxMountPointDir  = "rootfs"
	sandboxLowerDir       = "/"
)

type OverlayFsInfo struct {
	StagedFiles  []string `json:"stagedFiles"`
	ChangedFiles []string `json:"changedFiles"`
}

type OverlayFs struct {
	BaseDir  string
	FileInfo OverlayFsInfo
}

func createDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

func writeJsonToFile(data OverlayFsInfo, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

func readJsonFromFile(filePath string) OverlayFsInfo {
	file, err := os.Open(filePath)
	if err != nil {
		return OverlayFsInfo{}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var fileInfo OverlayFsInfo
	if err := decoder.Decode(&fileInfo); err != nil {
		return OverlayFsInfo{}
	}
	return fileInfo
}

func OpenOverlay(sandboxDir string) (*OverlayFs, error) {
	fileInfo := readJsonFromFile(path.Join(sandboxDir, sandboxConfigFileName))
	return &OverlayFs{
		BaseDir:  sandboxDir,
		FileInfo: fileInfo,
	}, nil
}

func CreateOverlay(sandboxDir string) (*OverlayFs, error) {
	if err := createDirIfNotExists(sandboxDir); err != nil {
		return nil, err
	}

	rootFsPath := path.Join(sandboxDir, sandboxMountPointDir)
	upperDirPath := path.Join(sandboxDir, sandboxUpperDir)
	workDirPath := path.Join(sandboxDir, sandboxWorkDir)

	if err := createDirIfNotExists(rootFsPath); err != nil {
		return nil, err
	}

	if err := createDirIfNotExists(upperDirPath); err != nil {
		return nil, err
	}

	if err := createDirIfNotExists(workDirPath); err != nil {
		return nil, err
	}

	// write a default config file
	config := OverlayFsInfo{
		StagedFiles:  []string{},
		ChangedFiles: []string{},
	}

	configFilePath := path.Join(sandboxDir, sandboxConfigFileName)
	if err := writeJsonToFile(config, configFilePath); err != nil {
		return nil, err
	}
	return OpenOverlay(sandboxDir)
}

func UnmountOverlay(mountDir string) error {
	return syscall.Unmount(mountDir, 0)
}

func (s *OverlayFs) Mount() error {
	opts := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s,userxattr",
		sandboxLowerDir,
		path.Join(s.BaseDir, sandboxUpperDir),
		path.Join(s.BaseDir, sandboxWorkDir),
	)
	return syscall.Mount("overlay", s.GetMountPath(), "overlay", 0, opts)
}

func (s *OverlayFs) UnMount() error {
	return syscall.Unmount(s.GetMountPath(), 0)
}

func (s *OverlayFs) CommitToDisk() error {
	return nil
}

func (s *OverlayFs) StageFile(filePath string) error {
	return nil
}

func (s *OverlayFs) UnstageFile(filePath string) error {
	return nil
}

func (s *OverlayFs) GetStagedFiles() ([]string, error) {
	return nil, nil
}

func (s *OverlayFs) GetChangedFiles() ([]string, error) {
	return nil, nil
}

func (s *OverlayFs) GetMountPath() string {
	return path.Join(s.BaseDir, sandboxMountPointDir)
}