package cli

import (
	"github.com/akamensky/argparse"
)

func GetStatusCommandParser(parser *argparse.Parser) (statusCommand *argparse.Command) {
	statusCommand = parser.NewCommand("status", "Show the status of the sandbox")
	return statusCommand
}

func ExecuteStatusCommand() error {
	return nil
}
