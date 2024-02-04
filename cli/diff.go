package cli

import (
	"fmt"
	"myapp/sandbox"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type DiffCommandArguments struct {
	sandboxId *string
	diffAll   *bool
}

func GetDiffCommandParser(parser *argparse.Parser) (diffCommand *argparse.Command, statusCommandArgs DiffCommandArguments) {
	diffCommand = parser.NewCommand("diff", "Show all the changes in a sandbox as a diff")
	return diffCommand, DiffCommandArguments{
		sandboxId: diffCommand.StringPositional(&argparse.Options{Required: true, Help: "the sandbox to commit"}),
		diffAll:   diffCommand.Flag("a", "all", &argparse.Options{Required: false, Default: false, Help: "show all changes, staged and unstaged"}),
	}
}

func ExecuteDiffCommand(statusCommandArgs DiffCommandArguments) error {
	sb, err := sandbox.LoadSandbox(*statusCommandArgs.sandboxId)
	if err != nil {
		return fmt.Errorf("failed to load sandbox: %w", err)
	}

	overlayFs := sb.GetOverlay()
	diff, err := overlayFs.GetChangedFiles()
	if err != nil {
		return err
	}

	// normally skip staged files

	for _, file := range diff {
		if overlayFs.IsStaged(file) && !*statusCommandArgs.diffAll {
			continue
		}

		overlayRealPath := overlayFs.GetFsPathForFile(file)
		// now we check on our system
		realPath := path.Join("/", file)

		if info, err := os.Stat(overlayRealPath); err != nil || info.IsDir() {
			continue
		}

		realFileContent := ""
		if info, err := os.Stat(realPath); os.IsNotExist(err) || info.IsDir() {
			realFileContent = ""
		} else {
			// read the file
			dat, err := os.ReadFile(realPath)
			if err != nil {
				return err
			}
			realFileContent = string(dat)
		}

		overlayContent, err := os.ReadFile(overlayRealPath)
		if err != nil {
			return err
		}

		// now we read the file and compare

		edits := myers.ComputeEdits(span.URIFromPath(realPath), realFileContent, string(overlayContent))
		fmt.Println(gotextdiff.ToUnified(realPath, overlayRealPath, realFileContent, edits))
	}

	return nil
}
