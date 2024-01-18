package cli

import (
	"github.com/akamensky/argparse"
)

type CommitCommandArguments struct {
	sandboxId *string
}

func GetCommitCommandParser(parser *argparse.Parser) (commitCommand *argparse.Command, statusCommandArgs CommitCommandArguments) {
	commitCommand = parser.NewCommand("commit", "commit a sandbox to your host machine")
	return commitCommand, CommitCommandArguments{
		sandboxId: commitCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to commit"}),
	}
}

func ExecuteCommitCommand(statusCommandArgs CommitCommandArguments) error {
	// todo
	return nil
}
