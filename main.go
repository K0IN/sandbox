package main

import (
	"fmt"
	"myapp/cli"
	"myapp/helper"
	sandbox "myapp/sandbox"
	"os"

	"github.com/akamensky/argparse"
)

// this function is called when we are inside the namespace and sets up the mounts and executes the command.
func executeSandbox(hostname, hostPath, command *string) int {
	// a fail save to prevent running the sandbox code on the host machine
	if os.Getuid() != 0 || os.Getgid() != 0 {
		panic("started in namespace mode but not as root")
	}
	sandboxResult := sandbox.CreateSandboxInsideNamespace(*command, *hostname, *hostPath)
	return sandboxResult
}

func main() {
	parser := argparse.NewParser("sandbox", "run a command in a sandbox")
	tryParser, tryArguments := cli.GetTryCommandParser(parser)
	diffParser, diffArguments := cli.GetDiffCommandParser(parser)
	statusParser, statusArguments := cli.GetStatusCommandParser(parser)
	listParser := cli.GetListCommandParser(parser)
	removeParser, removeArguments := cli.GetRemoveCommandParser(parser)
	addParser, addArguments := cli.GetAddCommandParser(parser)

	execute := parser.NewCommand("sandbox-entry", argparse.DisableDescription)
	sandboxHostName := execute.String("", "hostname", &argparse.Options{Required: true, Help: "hostname"})
	sandboxHostPath := execute.String("", "sandboxdir", &argparse.Options{Required: true, Help: "the directory on the host to mount the sandbox to"})
	sandboxEntryCommand := execute.String("", "command", &argparse.Options{Required: true, Default: helper.GetPrimaryShell(), Help: "the command to execute inside the sandbox"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
	} else if execute.Happened() {
		sandboxResult := executeSandbox(sandboxHostName, sandboxHostPath, sandboxEntryCommand)
		os.Exit(sandboxResult)
	} else if diffParser.Happened() {
		if err := cli.ExecuteDiffCommand(diffArguments); err != nil {
			panic(err)
		}
	} else if statusParser.Happened() {
		if err := cli.ExecuteStatusCommand(statusArguments); err != nil {
			panic(err)
		}
	} else if listParser.Happened() {
		if err := cli.ExecuteListCommand(); err != nil {
			panic(err)
		}
	} else if removeParser.Happened() {
		if err := cli.ExecuteRemoveCommand(removeArguments); err != nil {
			panic(err)
		}
	} else if tryParser.Happened() {
		executeResult := cli.ExecuteTryCommand(tryArguments)
		os.Exit(executeResult)
	} else if addParser.Happened() {
		if err := cli.ExecuteAddCommand(addArguments); err != nil {
			panic(err)
		}
	} else {
		fmt.Print(parser.Usage(nil))
	}
}
