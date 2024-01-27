package cli

import (
	"fmt"
	"myapp/helper"
	"myapp/sandbox"

	"github.com/akamensky/argparse"
)

type TryCommandArguments struct {
	SandboxId    *string
	Command      *string
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
		// AsRootUser:   tryCommand.Flag("r", "root", &argparse.Options{Required: false, Default: false, Help: "Run the command as root user"}),
		Command: tryCommand.StringPositional(&argparse.Options{Required: false, Default: helper.GetPrimaryShell(), Help: "The command to execute inside the sandbox, default is the primary shell"}),
	}
}

func ExecuteTryCommand(args TryCommandArguments) (int, error) {
	sb, err := sandbox.CreateSandbox()
	if err != nil {
		return 0, err
	}

	resultCode, path, err := sb.Execute("/bin/bash", sandbox.SandboxParams{
		AllowNetwork:      true,
		AllowEnv:          true,
		UserId:            0,
		GroupId:           0,
		AllowChangeUserId: true,
	})

	if err == nil {
		fmt.Printf("Sandbox created at %s\n", path)
	}

	return resultCode, err
}
