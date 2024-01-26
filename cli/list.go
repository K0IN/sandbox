package cli

import (
	"github.com/akamensky/argparse"
)

func GetListCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command) {
	diffCommand = parser.NewCommand("ls", "List all sandboxes")
	return diffCommand
}

func ExecuteListCommand() error {
	return nil // todo
}
