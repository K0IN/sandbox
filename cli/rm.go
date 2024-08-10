package cli

import (
	"fmt"
	"myapp/sandbox"
	"os"

	"github.com/akamensky/argparse"
)

type RemoveCommandArguments struct {
	sandboxId *[]string
}

func GetRemoveCommandParser(parser *argparse.Parser) (removeCommand *argparse.Command, statusCommandArgs RemoveCommandArguments) {
	removeCommand = parser.NewCommand("rm", "remove a sandbox form disk,")
	return removeCommand, RemoveCommandArguments{
		sandboxId: removeCommand.StringList("i", "id", &argparse.Options{Required: true, Help: "The id of the sandbox to remove"}),
	}
}

func ExecuteRemoveCommand(statusCommandArgs RemoveCommandArguments) error {
	for _, id := range *statusCommandArgs.sandboxId {
		path, err := sandbox.GetSandboxPathForId(id)
		if err != nil {
			fmt.Printf("failed to get sandbox path for id %s: %v\n", id, err)
		}
		if res, err := os.Stat(path); res == nil || err != nil {
			fmt.Printf("sandbox %s does not exist\n", id)
			continue
		}
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("failed to remove sandbox files %s: %v\n", id, err)
		}
		// fmt.Printf("sandbox %s removed\n", id)
	}
	return nil
}
