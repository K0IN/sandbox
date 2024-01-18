package cli

import (
	"github.com/akamensky/argparse"
)

func GetDiffCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command) {
	diffCommand = parser.NewCommand("diff", "Compare two files line by line")
	return diffCommand
}

func ExecuteDiffCommand() error {
	return nil
}
