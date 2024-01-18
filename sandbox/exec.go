package sandbox

import (
	"fmt"
	"myapp/helper"
	"os"
	"os/exec"
)

type SandboxConfig struct {
	AllowNetwork bool
	AllowEnv     bool
	SandboxId    string
	Command      string
	HostDir      string
}

/* this is the side OUTSIDE the namespace which will start the namespace part later */

func ForkSelfIntoNewNamespace(config SandboxConfig) int {
	// todo use cmd.SysProcAttr if unshare is not available

	unshareArguments := []string{
		"--mount",
		"--user",
		"--map-root-user",
		"--pid",
		"--fork",
		"--uts",
	}

	if config.AllowNetwork {
		unshareArguments = append(unshareArguments, "--net")
	}

	// fork arguments
	forkArguments := []string{
		os.Args[0],
		"sandbox-entry",
		"--hostname", config.SandboxId,
		"--sandboxdir", config.HostDir,
		"--command", config.Command,
	}
	arguments := append(unshareArguments, forkArguments...)

	cmd := exec.Command("unshare", arguments...)
	if config.AllowEnv {
		cmd.Env = os.Environ()
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()

	if piped, err := helper.IsOutputPiped(); !piped && err == nil {
		fmt.Printf("Sandbox: %s at %s exited with code: %d\n", config.SandboxId, config.HostDir, cmd.ProcessState.ExitCode())
	} else {
		println(config.SandboxId)
	}

	return cmd.ProcessState.ExitCode()
}
