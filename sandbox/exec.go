package sandbox

import (
	"os"
	"os/exec"
)

type SandboxConfig struct {
	AllowNetwork bool
	AllowEnv     bool
	Hostname     string
	Arguments    []string
}

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
	forkArguments := append([]string{os.Args[0], "sandbox-entry", "--hostname", config.Hostname}, config.Arguments...)
	arguments := append(unshareArguments, forkArguments...)

	cmd := exec.Command("unshare", arguments...)
	if config.AllowEnv {
		cmd.Env = os.Environ()
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println("Running sandbox")
	_ = cmd.Run()
	return cmd.ProcessState.ExitCode()
}
