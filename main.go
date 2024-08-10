package main

import (
	"fmt"
	"myapp/cli"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("You must run this program as root!")
		fmt.Printf("Current effective uid: %d, effective gid: %d\n", os.Geteuid(), os.Getegid())
		fmt.Println("Try running it as suid root or sudo.")
		fmt.Printf("sudo chown root:root %s && sudo chmod u+s %s\n", os.Args[0], os.Args[0])
		return
	}

	parser := argparse.NewParser("sandbox", "run a command in a sandbox")
	tryParser, tryArguments := cli.GetTryCommandParser(parser)
	diffParser, diffArguments := cli.GetDiffCommandParser(parser)
	statusParser, statusArguments := cli.GetStatusCommandParser(parser)
	listParser := cli.GetListCommandParser(parser)
	removeParser, removeArguments := cli.GetPruneCommandParser(parser)
	addParser, addArguments := cli.GetAddCommandParser(parser)
	commitParser, confirmArguments := cli.GetCommitCommandParser(parser)
	pruneParser, pruneArguments := cli.GetPruneCommandParser(parser)

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
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
		if err := cli.ExecutePruneCommand(removeArguments); err != nil {
			panic(err)
		}
	} else if tryParser.Happened() {
		executeResult, err := cli.ExecuteTryCommand(tryArguments)
		if err != nil {
			panic(err)
		}
		os.Exit(executeResult)
	} else if addParser.Happened() {
		if err := cli.ExecuteAddCommand(addArguments); err != nil {
			panic(err)
		}
	} else if commitParser.Happened() {
		if err := cli.ExecuteCommitCommand(confirmArguments); err != nil {
			panic(err)
		}
	} else if pruneParser.Happened() {
		if err := cli.ExecutePruneCommand(pruneArguments); err != nil {
			panic(err)
		}
	} else {
		fmt.Print(parser.Usage(nil))
	}
}
