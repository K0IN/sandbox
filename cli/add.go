package cli

import (
	"fmt"
	"myapp/sandbox"
	"path/filepath"

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
	sb, err := sandbox.LoadSandbox(*statusCommandArgs.sandboxId)
	if err != nil {
		return fmt.Errorf("failed to load sandbox: %w", err)
	}

	overlay := sb.GetOverlay()
	// first we get all the changed files
	changedFiles, err := overlay.GetChangedFiles()
	if err != nil {
		return err
	}

	for _, file := range changedFiles {
		if matched, err := filepath.Match(*statusCommandArgs.fileSelector, file); err == nil && matched {
			if *statusCommandArgs.remove {
				if err := overlay.UnstageFile(file); err != nil {
					return err
				}
				fmt.Printf("Unstaged %s\n", file)
			} else {
				if err := overlay.StageFile(file); err != nil {
					return err
				}
				fmt.Printf("Staged %s\n", file)
			}
		}
	}
	return nil
}
