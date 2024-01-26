package cli

import (
	"github.com/akamensky/argparse"
)

type RemoveCommandArguments struct {
	sandboxId *string
}

func GetRemoveCommandParser(parser *argparse.Parser) (removeCommand *argparse.Command, statusCommandArgs RemoveCommandArguments) {
	removeCommand = parser.NewCommand("rm", "remove a sandbox form disk,")
	return removeCommand, RemoveCommandArguments{
		sandboxId: removeCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to remove you can also use glob patterns"}),
	}
}

func ExecuteRemoveCommand(statusCommandArgs RemoveCommandArguments) error {
	return nil // todo
}
