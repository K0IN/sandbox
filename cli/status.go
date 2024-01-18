package cli

import (
	"fmt"
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
	sandbox, err := sandbox.LoadSandboxById(*statusCommandArgs.sandboxId)
	if err != nil {
		return err
	}
	status, err := sandbox.GetStatus()
	if err != nil {
		return err
	}

	// todo: add a status icon for each file, to show if it is staged or not
	fmt.Printf("Status of sandbox %s\n", *statusCommandArgs.sandboxId)
	fmt.Printf("Files: %d\n", len(status.Files))
	for _, file := range status.Files {
		fmt.Printf("  %s\n", file)
	}

	return nil
}
