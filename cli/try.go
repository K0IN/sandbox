package cli

import (
	"fmt"
	"myapp/sandbox"
	"os"

	"github.com/akamensky/argparse"
)

type TryCommandArguments struct {
	SandboxId    *string
	AllowNetwork *bool
	AllowEnv     *bool
	Persist      *bool
}

func GetTryCommandParser(parser *argparse.Parser) (tryCommand *argparse.Command, statusCommandArgs TryCommandArguments) {
	tryCommand = parser.NewCommand("try", "Execute a command inside a sandbox and review the changes")

	return tryCommand, TryCommandArguments{
		SandboxId:    tryCommand.String("i", "id", &argparse.Options{Required: false, Default: nil, Help: "Attach to a existing sandbox with the given in to run the command in"}),
		AllowNetwork: tryCommand.Flag("n", "network", &argparse.Options{Required: false, Default: false, Help: "Allow network"}),
		AllowEnv:     tryCommand.Flag("e", "env", &argparse.Options{Required: false, Default: false, Help: "Allow environment variables"}),
		Persist:      tryCommand.Flag("p", "persist", &argparse.Options{Required: false, Default: true, Help: "Persist the sandbox, else it will be deleted after the sandbox is exited"}),
	}
}

func ExecuteTryCommand(args TryCommandArguments) int {
	sandboxConfig := sandbox.SandboxConfig{
		AllowNetwork: *args.AllowNetwork,
		AllowEnv:     *args.AllowEnv,
		Arguments:    []string{},
	}

	if *args.SandboxId != "" {
		sb, err := sandbox.LoadSandboxById(*args.SandboxId)
		if err != nil {
			panic(fmt.Errorf("failed to load sandbox: %s %w", *args.SandboxId, err))
		}
		sandboxConfig.SandboxId = sb.SandboxId
		sandboxConfig.HostDir = sb.SandboxBaseDir
	} else if !*args.Persist {
		sb, err := sandbox.CreateSandbox()
		if err != nil {
			panic(fmt.Errorf("failed to create sandbox: %w", err))
		}
		sandboxConfig.SandboxId = sb.SandboxId
		sandboxConfig.HostDir = sb.SandboxBaseDir
	} else {
		sandboxConfig.SandboxId = "sandbox"
		sandboxDir, err := os.MkdirTemp("", "sandbox")
		if err != nil {
			panic(fmt.Errorf("failed to create sandbox tmp folder: %w", err))
		}
		sandboxConfig.HostDir = sandboxDir
	}

	returnCode := sandbox.ForkSelfIntoNewNamespace(sandboxConfig) // this will call us again with an argument
	return returnCode
}
