package cli

import (
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
		Command:      tryCommand.StringPositional(&argparse.Options{Required: false, Default: helper.GetPrimaryShell(), Help: "The command to execute inside the sandbox, default is the primary shell"}),
	}
}

func ExecuteTryCommand(args TryCommandArguments) (int, error) {
	sb, err := sandbox.CreateSandbox()
	if err != nil {
		return 0, err
	}

	resultCode, err := sb.Execute("ls", sandbox.SandboxParams{
		AllowNetwork:      true,
		AllowEnv:          true,
		UserId:            1000,
		GroupId:           1000,
		AllowChangeUserId: true,
	})

	return resultCode, err
}
