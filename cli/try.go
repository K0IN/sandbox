package cli

import (
	"fmt"
	"myapp/helper"
	"myapp/sandbox"
	"os"

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
		Command:      tryCommand.StringPositional(&argparse.Options{Required: false, Default: helper.GetPrimaryShell(), Help: "The command to execute inside the sandbox, default is the primary shell. If your command contains spaces, wrap it in quotes"}),
	}
}

func ExecuteTryCommand(args TryCommandArguments) (int, error) {
	sb, err := getSandbox(args)
	if err != nil {
		return 0, err
	}

	command := getCommand(args)
	resultCode, err := sb.Execute(command, sandbox.SandboxParams{
		AllowNetwork:      *args.AllowNetwork,
		AllowEnv:          *args.AllowEnv,
		UserId:            uint32(os.Getuid()),
		GroupId:           uint32(os.Getegid()),
		AllowChangeUserId: true,
	})

	if err == nil && *args.Persist {
		path := sb.GetPath()
		fmt.Printf("Sandbox created at %s\n", path)
	}

	if !*args.Persist && *args.SandboxId != "" {
		sb.Delete()
	}

	return resultCode, err
}

func getSandbox(args TryCommandArguments) (*sandbox.Sandbox, error) {
	if *args.SandboxId != "" {
		return sandbox.LoadSandbox(*args.SandboxId)
	}
	return sandbox.CreateSandbox()
}

func getCommand(args TryCommandArguments) string {
	if *args.Command != "" {
		return *args.Command
	}
	return helper.GetPrimaryShell()
}
