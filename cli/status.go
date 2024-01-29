package cli

import (
	"myapp/sandbox"

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
	sb, err := sandbox.LoadSandbox(*statusCommandArgs.sandboxId)
	if err != nil {
		return err
	}

	overlay := sb.GetOverlay()
	changedFiles, _ := overlay.GetChangedFiles()

	if len(changedFiles) == 0 {
		println("no changes")
		return nil
	}

	for _, changedFile := range changedFiles {
		if overlay.IsStaged(changedFile) {
			println("staged: " + changedFile)
		} else {
			println("not staged: " + changedFile)
		}
	}

	return nil
}
