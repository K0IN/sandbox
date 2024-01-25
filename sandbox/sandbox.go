package sandbox

import (
	"os"
	"path"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	sandboxDirName = ".sandboxes"
)

type Sandbox struct {
	overlayFs *OverlayFs
}

func CreateSandboxAt(sandboxBaseDir string) (*Sandbox, error) {
	overlay, err := CreateOverlay(sandboxBaseDir)
	if err != nil {
		return nil, err
	}
	return &Sandbox{
		overlayFs: overlay,
	}, nil
}

func CreateSandbox() (*Sandbox, error) {
	sandboxId, err := nanoid.New(6)
	if err != nil {
		return nil, err
	}
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	sandboxDir := path.Join(userDir, sandboxDirName, sandboxId)
	return CreateSandboxAt(sandboxDir)
}

func LoadSandbox(sandboxId string) (*Sandbox, error) {
	// todo
	return nil, nil
}

func (s *Sandbox) Execute(command string) error {
	// todo
	return nil
}

func (s *Sandbox) DeleteSandbox() error {
	// todo
	return nil
}

func (s *Sandbox) GetOverlay() (*OverlayFs, error) {
	// todo
	return nil, nil
}
