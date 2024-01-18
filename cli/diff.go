package cli

import (
	"github.com/akamensky/argparse"
)

type DiffCommandArguments struct {
	sandboxId *string
}

func GetDiffCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command, statusCommandArgs DiffCommandArguments) {
	diffCommand = parser.NewCommand("diff", "Show all the changes in a sandbox as a diff")
	return diffCommand, DiffCommandArguments{
		sandboxId: diffCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to commit"}),
	}
}

func ExecuteDiffCommand(statusCommandArgs DiffCommandArguments) error {
	// todo
	return nil
}
