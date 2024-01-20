package cli

import (
	"fmt"
	"myapp/helper"
	"myapp/sandbox"

	"github.com/akamensky/argparse"
)

func GetListCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command) {
	diffCommand = parser.NewCommand("ls", "List all sandboxes")
	return diffCommand
}

func ExecuteListCommand() error {
	allSandboxes, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}

	if isPiped := helper.IsOutputPiped(); isPiped {
		for _, sandbox := range allSandboxes {
			fmt.Println(sandbox.SandboxId)
		}
		return nil
	} else {
		fmt.Println("All sandboxes:")
		for _, sandbox := range allSandboxes {
			fmt.Printf("Sandbox id=%s\n", sandbox.SandboxId)
		}
	}

	return nil
}
