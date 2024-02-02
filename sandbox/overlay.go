package sandbox

import (
	"encoding/json"
	"fmt"
	"myapp/helper"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
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
	StagedFiles []string `json:"stagedFiles"`
}

type OverlayFs struct {
	BaseDir  string
	workDir  string
	upperDir string
	mountDir string

	mountedPaths []string
}

func writeOverlayInformation(data OverlayFsInfo, filePath string) error {
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

func readOverlayInformation(filePath string) (*OverlayFsInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var fileInfo OverlayFsInfo
	if err := decoder.Decode(&fileInfo); err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

func OpenOverlay(sandboxDir string) (*OverlayFs, error) {
	return &OverlayFs{
		BaseDir:  sandboxDir,
		workDir:  path.Join(sandboxDir, sandboxWorkDir),
		upperDir: path.Join(sandboxDir, sandboxUpperDir),
		mountDir: path.Join(sandboxDir, sandboxMountPointDir),
	}, nil
}

func CreateOverlay(sandboxDir string) (*OverlayFs, error) {
	if err := helper.CreateDirIfNotExists(sandboxDir); err != nil {
		return nil, err
	}

	// write a default config file
	config := OverlayFsInfo{
		StagedFiles: []string{},
	}

	configFilePath := path.Join(sandboxDir, sandboxConfigFileName)
	if err := writeOverlayInformation(config, configFilePath); err != nil {
		return nil, err
	}
	return OpenOverlay(sandboxDir)
}

func (m *OverlayFs) mountOverlayFs(lower string) error {
	opts := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s,userxattr",
		lower,
		path.Join(m.upperDir, lower),
		path.Join(m.workDir, lower),
	)

	mountPoint := path.Join(m.GetMountPath(), lower)
	_ = os.MkdirAll(mountPoint, 0755)
	_ = os.MkdirAll(path.Join(m.upperDir, lower), 0755)
	_ = os.MkdirAll(path.Join(m.workDir, lower), 0755)

	if err := syscall.Mount("overlay", mountPoint, "overlay", 0, opts); err != nil {
		return fmt.Errorf("failed to mount %s to %s: %s", lower, mountPoint, err)
	}

	m.mountedPaths = append(m.mountedPaths, mountPoint)
	return nil
}

func (m *OverlayFs) mountRecursive(mounts []helper.MountInfo) error {
	for _, mount := range mounts {
		if err := m.mountOverlayFs(mount.Target); err != nil {
			fmt.Printf("failed to mount overlay fs: %s\n", err)
		}

		if mount.Children != nil {
			if err := m.mountRecursive(*mount.Children); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *OverlayFs) Mount() error {
	if err := helper.CreateDirIfNotExists(s.mountDir); err != nil {
		return fmt.Errorf("failed to create mount dir: %w", err)
	}

	if err := helper.CreateDirIfNotExists(s.upperDir); err != nil {
		return fmt.Errorf("failed to create upper dir: %w", err)
	}

	if err := helper.CreateDirIfNotExists(s.workDir); err != nil {
		return fmt.Errorf("failed to create work dir: %w", err)
	}
	allMounts, err := helper.FindAllMounts()
	if err != nil {
		return fmt.Errorf("failed to find all mounts: %w", err)
	}
	return s.mountRecursive(allMounts)
}

func (s *OverlayFs) UnMount() error {
	// sort the mounted paths by length, so that we unmount the deepest path first
	sort.Slice(s.mountedPaths, func(i, j int) bool {
		iPath := s.mountedPaths[i]
		jPath := s.mountedPaths[j]
		iParts := strings.Split(iPath, string(os.PathSeparator))
		jParts := strings.Split(jPath, string(os.PathSeparator))
		return len(iParts) > len(jParts)
	})

	for _, mountPoint := range s.mountedPaths {
		if err := syscall.Unmount(mountPoint, 0); err != nil {
			fmt.Printf("failed to unmount %s: %s\n", mountPoint, err)
		}
	}

	if err := os.RemoveAll(s.mountDir); err != nil {
		return err
	}
	return nil
}

func (s *OverlayFs) CommitToDisk() error {
	stagedFiles, err := s.GetStagedFiles()
	if err != nil {
		return err
	}

	for _, stagedFile := range stagedFiles {
		upperFilePath := path.Join(s.upperDir, stagedFile)
		lowerFilePath := path.Join(sandboxLowerDir, stagedFile)

		// copy the file from upper to lower
		if err := helper.CopyFile(upperFilePath, lowerFilePath); err != nil {
			return err
		}
	}

	return nil
}

func (s *OverlayFs) StageFile(filePath string) error {
	stagedFiles, err := s.GetStagedFiles()
	if err != nil {
		return err
	}

	// check if the file is already staged
	if s.IsStaged(filePath) {
		return nil
	}

	stagedFiles = append(stagedFiles, filePath)

	return writeOverlayInformation(OverlayFsInfo{StagedFiles: stagedFiles}, path.Join(s.BaseDir, sandboxConfigFileName))
}

func (s *OverlayFs) UnstageFile(filePath string) error {
	// remove the file from the staged files
	stagedFiles, err := s.GetStagedFiles()
	if err != nil {
		return err
	}

	newStagedFiles := []string{}
	for _, stagedFile := range stagedFiles {
		if stagedFile != filePath {
			newStagedFiles = append(newStagedFiles, stagedFile)
		}
	}

	fileInfo := &OverlayFsInfo{
		StagedFiles: newStagedFiles,
	}
	configFilePath := path.Join(s.BaseDir, sandboxConfigFileName)
	return writeOverlayInformation(*fileInfo, configFilePath)
}

func (s *OverlayFs) GetStagedFiles() ([]string, error) {
	fileInfo, err := readOverlayInformation(path.Join(s.BaseDir, sandboxConfigFileName))
	if err != nil {
		return nil, err
	}
	return fileInfo.StagedFiles, nil
}

func (s *OverlayFs) GetChangedFiles() ([]string, error) {
	allChangedFiles := []string{}
	// walk the upper dir and find all files that are not in the lower dir
	err := filepath.Walk(s.upperDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(s.upperDir, path)
			if err != nil {
				return err
			}

			allChangedFiles = append(allChangedFiles, relativePath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allChangedFiles, nil
}

func (s *OverlayFs) IsStaged(filePath string) bool {
	stagedFiles, err := s.GetStagedFiles()
	if err != nil {
		return false
	}
	for _, stagedFile := range stagedFiles {
		if filePath == stagedFile {
			return true
		}
	}
	return false
}

func (s *OverlayFs) GetMountPath() string {
	return path.Join(s.BaseDir, sandboxMountPointDir)
}
