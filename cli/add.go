package cli

import (
	"github.com/akamensky/argparse"
)

type AddCommandArguments struct {
	sandboxId    *string
	remove       *bool
	fileSelector *string
}

func GetAddCommandParser(parser *argparse.Parser) (addCommand *argparse.Command, statusCommandArgs AddCommandArguments) {
	addCommand = parser.NewCommand("add", "Add or Remove a file from staging to be committed")
	return addCommand, AddCommandArguments{
		sandboxId:    addCommand.StringPositional(&argparse.Options{Required: true, Help: "The sandbox to add or remove a file from staging"}),
		remove:       addCommand.Flag("r", "remove", &argparse.Options{Required: false, Default: false, Help: "Remove the file from staging"}),
		fileSelector: addCommand.StringPositional(&argparse.Options{Required: true, Help: "The file to add or remove from staging, globs are supported"}),
	}
}

func ExecuteAddCommand(statusCommandArgs AddCommandArguments) error {
	return nil // todo
}
