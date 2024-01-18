package cli

import (
	"fmt"
	"myapp/sandbox"

	"github.com/akamensky/argparse"
)

type RemoveCommandArguments struct {
	sandboxId *string
}

func GetRemoveCommandParser(parser *argparse.Parser) (statusCommand *argparse.Command, statusCommandArgs StatusCommandArguments) {
	statusCommand = parser.NewCommand("rm", "remove a sandbox form disk,")
	return statusCommand, StatusCommandArguments{
		sandboxId: statusCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to remove use * for all"}),
	}
}

func ExecuteRemoveCommand(statusCommandArgs StatusCommandArguments) error {
	if *statusCommandArgs.sandboxId == "*" {
		return removeAllSandboxes()
	} else {
		return removeSandbox(*statusCommandArgs.sandboxId)
	}
}

func removeSandbox(sandboxId string) error {
	all, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}
	for _, sandbox := range all {
		if sandbox.SandboxId == sandboxId {
			fmt.Printf("Removing sandbox %s\n", sandboxId)
			return sandbox.Remove()
		}
	}

	fmt.Printf("Sandbox %s not found\n", sandboxId)
	return nil
}

func removeAllSandboxes() error {
	all, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}
	for _, sandbox := range all {
		if err := sandbox.Remove(); err != nil {
			return err
		}
	}
	return nil
}
