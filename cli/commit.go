package cli

import (
	"github.com/akamensky/argparse"
)

type CommitCommandArguments struct {
	sandboxId *string
	yes       *bool
}

func GetCommitCommandParser(parser *argparse.Parser) (commitCommand *argparse.Command, statusCommandArgs CommitCommandArguments) {
	commitCommand = parser.NewCommand("commit", "commit a sandbox to your host machine")
	return commitCommand, CommitCommandArguments{
		sandboxId: commitCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to commit"}),
		yes:       commitCommand.Flag("y", "yes", &argparse.Options{Required: false, Default: false, Help: "skip confirmation"}),
	}
}

func ExecuteCommitCommand(statusCommandArgs CommitCommandArguments) error {
	return nil // todo
}
