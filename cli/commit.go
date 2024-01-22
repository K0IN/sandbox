package cli

import (
	"myapp/helper"
	"myapp/sandbox"

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
	sandbox, err := sandbox.LoadSandboxById(*statusCommandArgs.sandboxId)
	if err != nil {
		return err
	}

	if !*statusCommandArgs.yes {
		manualConfirmation := helper.Confirm("Are you sure you want to commit this sandbox? This will overwrite any existing files with the same name in your host machine.")
		if !manualConfirmation {
			return nil
		}
	}
	return sandbox.Commit()
}
