package main

import (
	"fmt"
	"myapp/cli"
	sandbox "myapp/sandbox"
	"os"

	"github.com/akamensky/argparse"
)

func executeSandbox(hostname *string) int {
	if os.Getuid() != 0 || os.Getgid() != 0 {
		panic("started in namespace mode but not as root")
	}
	// we are inside the new namespace
	shell := sandbox.GetPrimaryShell()
	println("Starting sandbox with shell", shell)
	sandboxResult := sandbox.CreateSandboxInsideNamespace(shell, *hostname)
	return sandboxResult
}

func main() {
	parser := argparse.NewParser("sandbox", "run a command in a sandbox")
	tryParser, tryArguments := cli.GetTryCommandParser(parser)
	diffParser := cli.GetDiffCommandParser(parser)
	statusParser := cli.GetStatusCommandParser(parser)
	listParser := cli.GetListCommandParser(parser)
	// skipDiff := tryCommand.Flag("s", "skip-diff", &argparse.Options{Required: false, Default: false, Help: "skip diff"})
	// workDir := tryCommand.String("c", "workdir", &argparse.Options{Required: false, Help: "workdir"})
	// args := parser.StringPositional(&argparse.Options{Required: true, Help: "command to execute"})
	// commit := parser.NewCommand("commit", "commit the previously added files to disk")
	// rm := parser.NewCommand("rm", "forcefully remove all sandbox files")

	// do NOT use this command directly, use try instead, else you will run the container code on your host!
	execute := parser.NewCommand("sandbox-entry", argparse.DisableDescription)
	sandboxHostName := execute.String("", "hostname", &argparse.Options{Required: true, Help: "hostname"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
	} else if execute.Happened() {
		sandboxResult := executeSandbox(sandboxHostName)
		os.Exit(sandboxResult)
	} else if diffParser.Happened() {
		if err := cli.ExecuteDiffCommand(); err != nil {
			panic(err)
		}
	} else if statusParser.Happened() {
		if err := cli.ExecuteStatusCommand(); err != nil {
			panic(err)
		}
	} else if listParser.Happened() {
		if err := cli.ExecuteListCommand(); err != nil {
			panic(err)
		}
	} else if tryParser.Happened() {
		executeResult := cli.ExecuteTryCommand(tryArguments)
		os.Exit(executeResult)
	}
}
