package cli

import (
	"fmt"
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
	sb, err := sandbox.LoadSandbox(*statusCommandArgs.sandboxId)
	if err != nil {
		return fmt.Errorf("failed to load sandbox: %w", err)
	}

	overlayFs := sb.GetOverlay()
	if *statusCommandArgs.yes || helper.Confirm("Are you sure you want to commit?") {
		return overlayFs.CommitToDisk()
	}
	return fmt.Errorf("commit aborted")
}
