package sandbox

import (
	"os"
	"os/exec"
)

type SandboxConfig struct {
	AllowNetwork bool
	AllowEnv     bool
	Arguments    []string
}

func ForkSelfIntoNewNamespace(config SandboxConfig) {
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
	forkArguments := append([]string{os.Args[0], "sandbox-secret"}, config.Arguments...)
	arguments := append(unshareArguments, forkArguments...)

	cmd := exec.Command("unshare", arguments...)
	if config.AllowEnv {
		cmd.Env = os.Environ()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// panic(err)
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
