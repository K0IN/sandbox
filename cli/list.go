package cli

import (
	"myapp/sandbox"

	"github.com/akamensky/argparse"
)

func GetListCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command) {
	diffCommand = parser.NewCommand("ls", "List all sandboxes")
	return diffCommand
}

func ExecuteListCommand() error {
	sandboxes, err := sandbox.ListSandboxes()
	if err != nil {
		return err
	}

	for _, sandbox := range sandboxes {
		println(sandbox)
	}
	return nil // todo
}
