package sandbox

import (
	"encoding/json"
	"os"
	"path"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type Sandbox struct {
	SandboxId string
}

type SandboxInfo struct {
	AddedFiles []string `json:"addedFiles"`
}

func getSandboxBaseDir() (baseDir string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".sandboxes"), nil
}

func CreateSandbox() (sandbox *Sandbox, err error) {
	sandboxId, err := nanoid.New(8)
	if err != nil {
		return nil, err
	}

	// append sb to the sandbox name
	sandboxId = "sb-" + sandboxId

	baseDir, err := getSandboxBaseDir()
	if err != nil {
		return nil, err
	}
	sandboxFolder := path.Join(baseDir, sandboxId)
	err = os.MkdirAll(sandboxFolder, 0755)
	if err != nil {
		return nil, err
	}

	sandboxConfigFilePath := path.Join(sandboxFolder, "config.json")
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
		SandboxId: sandboxId,
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
			SandboxId: file.Name(),
		}
		sandboxes = append(sandboxes, &sandbox)
	}

	return sandboxes, nil
}

func LoadSandboxById(sandboxId string) (sandbox *Sandbox, err error) {
	// check if the sandbox directory exists
	baseDir, err := getSandboxBaseDir()
	if err != nil {
		return nil, err
	}

	sandboxFolder := path.Join(baseDir, sandboxId)

	_, err = os.Stat(sandboxFolder)
	if err != nil {
		return nil, err
	}

	return &Sandbox{
		SandboxId: sandboxId,
	}, nil
}

func (sandbox *Sandbox) GetStatus() (status string, err error) {
	return "", nil
}

func (sandbox *Sandbox) Remove() error {
	return nil
}
