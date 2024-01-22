package sandbox

import (
	"encoding/json"
	"fmt"
	"myapp/helper"
	"os"
	"path"
	"path/filepath"
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	SandboxDirName        = ".sandboxes"
	SandboxConfigFileName = "config.json"
	SandboxUpperDir       = "upper"
	SandboxWorkDir        = "workdir"
	SandboxRootFs         = "rootfs"
)

type SandboxDirectories struct {
	SandboxDir       string
	RootFsBasePath   string
	UpperDirBasePath string
	WorkDirBasePath  string
}

type Sandbox struct {
	SandboxId  string
	SandboxDir string
}

type SandboxInfo struct {
	StagedFiles  []string `json:"stagedFiles"`
	ChangedFiles []string `json:"changedFiles"`
}

func createSandboxDirs(sandboxDir string) (*SandboxDirectories, error) {
	err := os.MkdirAll(sandboxDir, 0755) // ensure sandbox dir exists
	if err != nil {
		return nil, err
	}

	// lets first create all the directories we need
	rootFsBasePath := path.Join(sandboxDir, SandboxRootFs)
	upperDirBasePath := path.Join(sandboxDir, SandboxUpperDir)
	workDirBasePath := path.Join(sandboxDir, SandboxWorkDir)

	if err := os.MkdirAll(rootFsBasePath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(upperDirBasePath, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(workDirBasePath, 0755); err != nil {
		return nil, err
	}

	return &SandboxDirectories{
		SandboxDir:       sandboxDir,
		RootFsBasePath:   rootFsBasePath,
		UpperDirBasePath: upperDirBasePath,
		WorkDirBasePath:  workDirBasePath,
	}, nil
}

func getSandboxBaseDir() (baseDir string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, SandboxDirName), nil
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

	if _, err = createSandboxDirs(sandboxFolder); err != nil {
		return nil, err
	}

	return &Sandbox{
		SandboxId:  sandboxId,
		SandboxDir: sandboxFolder,
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
			SandboxId:  file.Name(),
			SandboxDir: path.Join(baseDir, file.Name()),
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
		SandboxId:  sandboxId,
		SandboxDir: sandboxFolder,
	}, nil
}

func (sandbox *Sandbox) GetStatus() (status *SandboxInfo, err error) {
	upperDir := path.Join(sandbox.SandboxDir, SandboxUpperDir)
	// loop through all files in the upper directory
	var files []string
	err = filepath.Walk(upperDir, func(path string, info os.FileInfo, _err error) error {
		relativePath, err := filepath.Rel(upperDir, path)
		if err != nil {
			return err
		}

		// we ignore the tmp directory
		if strings.HasPrefix(relativePath, "tmp/") {
			return nil
		}

		files = append(files, relativePath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sandboxConfigFilePath := path.Join(sandbox.SandboxDir, SandboxConfigFileName)
	configFile, err := os.OpenFile(sandboxConfigFilePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var sandboxInfo SandboxInfo
	err = json.NewDecoder(configFile).Decode(&sandboxInfo)
	if err != nil {
		return nil, err
	}

	return &SandboxInfo{
		StagedFiles:  sandboxInfo.StagedFiles,
		ChangedFiles: files,
	}, nil
}

func (sandbox *Sandbox) AddStagedFile(file string) error {
	// add the staged files to the config file
	sandboxConfigFilePath := path.Join(sandbox.SandboxDir, SandboxConfigFileName)
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
	sandboxConfigFilePath := path.Join(sandbox.SandboxDir, SandboxConfigFileName)
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
	sandboxConfigFilePath := path.Join(sandbox.SandboxDir, SandboxConfigFileName)
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
	// get all staged files
	sandboxStatus, err := sandbox.GetStatus()
	if err != nil {
		return err
	}

	print(len(sandboxStatus.StagedFiles), " staged files", len(sandboxStatus.ChangedFiles), " changed files\n")

	if len(sandboxStatus.StagedFiles) == 0 {
		fmt.Println("No staged files to commit")
		return nil
	}

	// copy all staged files to the upper directory
	for _, stagedFile := range sandboxStatus.StagedFiles {
		// check if the file exists in the upper directory
		upperFilePath := path.Join(sandbox.SandboxDir, SandboxUpperDir, stagedFile)
		originalFilePath := path.Join("/", stagedFile)
		// copy the file to the upper directory
		os.MkdirAll(path.Dir(originalFilePath), 0755)
		// commit the file from the upper directory to the rootfs

		fmt.Printf("Committing file %s to %s\n", upperFilePath, originalFilePath)
		if err := helper.CopyFile(upperFilePath, originalFilePath); err != nil {
			return err
		}
	}

	return nil
}

func (sandbox *Sandbox) Remove() error {
	return os.RemoveAll(sandbox.SandboxDir)
}
