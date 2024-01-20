package cli

import (
	"fmt"
	"myapp/sandbox"
	"path"
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
	sandbox, err := sandbox.LoadSandboxById(*statusCommandArgs.sandboxId)
	if err != nil {
		return err
	}

	status, err := sandbox.GetStatus()
	if err != nil {
		return err
	}

	// get paths relative to the sandbox

	selectedFiles := []string{}
	for _, file := range status.ChangedFiles {
		relativeDir, _ := filepath.Rel(sandbox.SandboxDir, file)
		if match, err := path.Match(*statusCommandArgs.fileSelector, relativeDir); err == nil && match {
			selectedFiles = append(selectedFiles, file)
		}
	}

	if *statusCommandArgs.remove {
		for _, file := range selectedFiles {
			fmt.Printf("Removing file %s from staging\n", file)
			err := sandbox.RemoveStagedFile(file)
			if err != nil {
				return fmt.Errorf("Error removing file %s from staging: %s", file, err)
			}
		}
	} else {
		for _, file := range selectedFiles {
			fmt.Printf("Adding file %s to staging\n", file)
			err := sandbox.AddStagedFile(file)
			if err != nil {
				return fmt.Errorf("Error adding file %s to staging: %s", file, err)
			}
		}
	}

	return nil
}
