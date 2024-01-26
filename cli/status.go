package cli

import (
	"github.com/akamensky/argparse"
)

type StatusCommandArguments struct {
	sandboxId *string
}

func GetStatusCommandParser(parser *argparse.Parser) (statusCommand *argparse.Command, statusCommandArgs StatusCommandArguments) {
	statusCommand = parser.NewCommand("status", "Show the status of the sandbox")

	return statusCommand, StatusCommandArguments{
		sandboxId: statusCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to show the status of"}),
	}
}

func ExecuteStatusCommand(statusCommandArgs StatusCommandArguments) error {
	return nil // todo
}
