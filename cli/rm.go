package cli

import (
	"fmt"
	"myapp/sandbox"
	"path"

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
	all, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}

	found := false
	for _, sandbox := range all {
		if match, err := path.Match(*statusCommandArgs.sandboxId, sandbox.SandboxId); err == nil && match {
			found = true
			fmt.Printf("Removing sandbox %s\n", sandbox.SandboxId)
			err := sandbox.Remove()
			if err != nil {
				return fmt.Errorf("Error removing sandbox %s: %s", sandbox.SandboxId, err)
			}
		}
	}

	if !found && len(all) > 0 {
		return fmt.Errorf("Sandbox not found")
	}
	return nil
}
