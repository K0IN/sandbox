package cli

import (
	"fmt"
	"myapp/sandbox"
	"os"

	"github.com/akamensky/argparse"
)

type PruneCommandArguments struct {
	force *bool
}

func GetPruneCommandParser(parser *argparse.Parser) (removeCommand *argparse.Command, statusCommandArgs PruneCommandArguments) {
	removeCommand = parser.NewCommand("prune", "remove all sandboxes")
	return removeCommand, PruneCommandArguments{
		force: removeCommand.Flag("f", "force", &argparse.Options{Required: false, Help: "Dont ask for confirmation"}),
	}
}

func ExecutePruneCommand(statusCommandArgs PruneCommandArguments) error {
	allSandboxIds, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}

	if !*statusCommandArgs.force {
		fmt.Println("Are you sure you want to remove all sandboxes? [y/N]")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			fmt.Println("Aborted")
			return nil
		}
	}

	for _, id := range allSandboxIds {
		path, err := sandbox.GetSandboxPathForId(id)
		if err != nil {
			fmt.Printf("failed to get sandbox path for id %s: %v\n", id, err)
		}

		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("failed to remove sandbox files %s: %v\n", id, err)
		}
		// fmt.Printf("sandbox %s removed\n", id)

	}
	return nil
}
