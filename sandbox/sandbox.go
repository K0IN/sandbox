package sandbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	SandboxConfigFileName = "config.json"
	SandboxUpperDir       = "upper"
	SandboxWorkDir        = "workdir"
	SandboxRootFs         = "rootfs"
)

type Sandbox struct {
	SandboxId      string
	SandboxBaseDir string
}

type SandboxInfo struct {
	StagedFiles  []string `json:"stagedFiles"`
	ChangedFiles []string `json:"changedFiles"`
}

func getSandboxBaseDir() (baseDir string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".sandboxes"), nil
}

func getDirForSandbox(sandboxId string) (sandboxDir string, err error) {
	baseDir, err := getSandboxBaseDir()
	if err != nil {
		return "", err
	}
	return path.Join(baseDir, sandboxId), nil
}

func CreateSandbox() (sandbox *Sandbox, err error) {
	sandboxId, err := nanoid.New(8)
	if err != nil {
		return nil, err
	}

	// append sb to the sandbox name
	sandboxId = "sb-" + sandboxId

	sandboxFolder, err := getDirForSandbox(sandboxId)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(sandboxFolder, 0755)
	if err != nil {
		return nil, err
	}

	sandboxConfigFilePath := path.Join(sandboxFolder, SandboxConfigFileName)
	configFile, err := os.Create(sandboxConfigFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	// write a empty config file
	json.NewEncoder(configFile).Encode(SandboxInfo{
		StagedFiles: []string{},
	})

	return &Sandbox{
		SandboxId:      sandboxId,
		SandboxBaseDir: sandboxFolder,
	}, nil
}

func ListSandboxes() (sandboxes []*Sandbox, err error) {
	baseDir, err := getSandboxBaseDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		sandbox := Sandbox{
			SandboxId:      file.Name(),
			SandboxBaseDir: path.Join(baseDir, file.Name()),
		}
		sandboxes = append(sandboxes, &sandbox)
	}

	return sandboxes, nil
}

func LoadSandboxById(sandboxId string) (sandbox *Sandbox, err error) {
	// check if the sandbox directory exists
	sandboxFolder, err := getDirForSandbox(sandboxId)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(sandboxFolder)
	if err != nil {
		return nil, fmt.Errorf("sandbox not found with id: %s (%s)", sandboxId, sandboxFolder)
	}

	return &Sandbox{
		SandboxId:      sandboxId,
		SandboxBaseDir: sandboxFolder,
	}, nil
}

func (sandbox *Sandbox) GetStatus() (status *SandboxInfo, err error) {
	upperDir := path.Join(sandbox.SandboxBaseDir, SandboxUpperDir)
	// loop through all files in the upper directory
	var files []string
	err = filepath.Walk(upperDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &SandboxInfo{
		StagedFiles:  []string{},
		ChangedFiles: files,
	}, nil
}

func (sandbox *Sandbox) AddStagedFile(file string) error {
	// add the staged files to the config file
	sandboxConfigFilePath := path.Join(sandbox.SandboxBaseDir, SandboxConfigFileName)
	configFile, err := os.OpenFile(sandboxConfigFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	var sandboxInfo SandboxInfo
	err = json.NewDecoder(configFile).Decode(&sandboxInfo)
	if err != nil {
		return err
	}

	sandboxInfo.StagedFiles = append(sandboxInfo.StagedFiles, file)

	configFile.Seek(0, 0)
	err = json.NewEncoder(configFile).Encode(sandboxInfo)
	if err != nil {
		return err
	}

	return nil
}

func (sandbox *Sandbox) RemoveStagedFile(file string) error {
	// add the staged files to the config file
	sandboxConfigFilePath := path.Join(sandbox.SandboxBaseDir, SandboxConfigFileName)
	configFile, err := os.OpenFile(sandboxConfigFilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	var sandboxInfo SandboxInfo
	err = json.NewDecoder(configFile).Decode(&sandboxInfo)
	if err != nil {
		return err
	}

	for i, stagedFile := range sandboxInfo.StagedFiles {
		if stagedFile == file {
			sandboxInfo.StagedFiles = append(sandboxInfo.StagedFiles[:i], sandboxInfo.StagedFiles[i+1:]...)
			break
		}
	}

	configFile.Seek(0, 0)
	err = json.NewEncoder(configFile).Encode(sandboxInfo)
	if err != nil {
		return err
	}

	return nil
}

func (sandbox *Sandbox) IsStaged(file string) (bool, error) {
	// add the staged files to the config file
	sandboxConfigFilePath := path.Join(sandbox.SandboxBaseDir, SandboxConfigFileName)
	configFile, err := os.OpenFile(sandboxConfigFilePath, os.O_RDWR, 0644)
	if err != nil {
		return false, err
	}
	defer configFile.Close()

	var sandboxInfo SandboxInfo
	err = json.NewDecoder(configFile).Decode(&sandboxInfo)
	if err != nil {
		return false, err
	}

	for _, stagedFile := range sandboxInfo.StagedFiles {
		if stagedFile == file {
			return true, nil
		}
	}

	return false, nil
}

func (sandbox *Sandbox) Commit() error {
	return nil
}

func (sandbox *Sandbox) Remove() error {
	return os.RemoveAll(sandbox.SandboxBaseDir)
}
