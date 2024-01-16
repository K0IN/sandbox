package main

import (
	"fmt"
	sandbox "myapp/sandbox"
	"os"
	"os/exec"
	"syscall"

	"github.com/akamensky/argparse"
)

func createSandboxInsideNamespace(entryCommand string) int {
	sandboxPaths, err := sandbox.CreateSandboxDirectories()
	if err != nil {
		panic(err)
	}

	println("Created sandbox directories", sandboxPaths.SandboxDir)

	if err := sandbox.MountDevices(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}

	if err := sandbox.MountProc(sandboxPaths.RootFsBasePath); err != nil {
		panic(err)
	}

	_ = sandbox.SetSandboxHostname()

	// current dir
	currentWorkingDir := "/"
	if workDir, err := os.Getwd(); err == nil {
		currentWorkingDir = workDir
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd \"%s\" && %s", currentWorkingDir, entryCommand))
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: sandboxPaths.RootFsBasePath,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// panic(err)
	}

	return cmd.ProcessState.ExitCode()
}

func main() {
	secret := "sandbox-secret"
	if len(os.Args) > 1 && os.Args[1] == secret {
		if os.Getuid() != 0 || os.Getgid() != 0 {
			panic("started in namespace mode but not as root")
		}
		// we are inside the new namespace
		shell := sandbox.GetPrimaryShell()
		println("Starting sandbox with shell", shell)
		sandboxResult := createSandboxInsideNamespace(shell)
		os.Exit(sandboxResult)
		return
	}

	parser := argparse.NewParser("sandbox", "run a command in a sandbox")
	tryCommand := parser.NewCommand("try", "execute a command inside a sandbox and review the changes")
	// --disable-network
	// --disable-env
	// --env key=value

	// skipDiff := tryCommand.Flag("s", "skip-diff", &argparse.Options{Required: false, Default: false, Help: "skip diff"})
	// workDir := tryCommand.String("c", "workdir", &argparse.Options{Required: false, Help: "workdir"})
	// args := parser.StringPositional(&argparse.Options{Required: true, Help: "command to execute"})

	// skipDiff := tryCommand.Flag("s", "skip-diff", &argparse.Options{Required: false, Default: false, Help: "skip diff"})

	// diffCommand := parser.NewCommand("diff", "show a diff of all files")
	// addCommand := parser.NewCommand("add", "add files from the sandbox to be committed")
	// --restore staged

	// status := parser.NewCommand("status", "show the status of the sandbox")

	// commit := parser.NewCommand("commit", "commit the previously added files to disk")

	// rm := parser.NewCommand("rm", "forcefully remove all sandbox files")

	// diff := parser.NewCommand("diff", "show a diff of all files")

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if tryCommand.Happened() {
		sandbox.ForkSelfIntoNewNamespace(sandbox.SandboxConfig{AllowNetwork: false, AllowEnv: true, Arguments: os.Args[1:]}) // this will call us again with an argument
	}
}
