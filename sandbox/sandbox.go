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
	AddedFiles []string `json:"addedFiles"`
}

type SandboxStatus struct {
	Files []string `json:"addedFiles"`
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
		AddedFiles: []string{},
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

func (sandbox *Sandbox) GetStatus() (status *SandboxStatus, err error) {
	fmt.Printf("sandbox folder: %s\n", sandbox.SandboxBaseDir)
	// find all files inside the sandbox folder
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

	return &SandboxStatus{
		Files: files,
	}, nil
}

func (sandbox *Sandbox) Remove() error {
	return os.RemoveAll(sandbox.SandboxBaseDir)
}
