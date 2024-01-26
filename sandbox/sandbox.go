package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	sandboxDirName = ".sandboxes"
)

type Sandbox struct {
	overlayFs *OverlayFs
}

type SandboxParams struct {
	// todo
	AllowNetwork bool
	AllowEnv     bool
	// user mode
	UserId            uint32
	GroupId           uint32
	AllowChangeUserId bool
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

func LoadSandboxFrom(sandboxBaseDir string) (*Sandbox, error) {
	overlay, err := OpenOverlay(sandboxBaseDir)
	if err != nil {
		return nil, err
	}
	sandbox := &Sandbox{
		overlayFs: overlay,
	}
	return sandbox, nil
}

func (s *Sandbox) LoadSandbox(sandboxId string) (*Sandbox, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	sandboxDir := path.Join(userDir, sandboxDirName, sandboxId)
	return LoadSandboxFrom(sandboxDir)
}

func (s *Sandbox) Execute(command string, params SandboxParams) (returnCode int, err error) {
	// todo mount devices
	currentWorkingDir := "/"
	if workDir, err := os.Getwd(); err == nil {
		currentWorkingDir = workDir
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, command))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65536,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65536,
			},
		},
		Credential: &syscall.Credential{
			Uid: params.UserId,
			Gid: params.GroupId,
		},

		GidMappingsEnableSetgroups: params.AllowChangeUserId, // enable su command
		// Chroot:                     s.overlayFs.GetRootFsPath(),
	}

	if !params.AllowNetwork {
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
	}

	if params.AllowEnv {
		cmd.Env = os.Environ()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	return cmd.ProcessState.ExitCode(), nil
}

func (s *Sandbox) DeleteSandbox() error {
	_ = s.overlayFs.Destroy()
	return os.RemoveAll(s.overlayFs.BaseDir)
}

func (s *Sandbox) GetOverlay() *OverlayFs {
	return s.overlayFs
}
