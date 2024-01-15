package main

import (
	"fmt"
	sandbox "myapp/sandbox"
	"os"
	"os/exec"
	"syscall"

	"github.com/akamensky/argparse"
)

func createSandboxInsideNamespace(entryCommand string) {
	sandboxDir, rootFs, _, err := sandbox.CreateSandboxDirectories()
	if err != nil {
		panic(err)
	}

	println("Created sandbox directories", sandboxDir)

	if err := sandbox.MountDevices(rootFs); err != nil {
		panic(err)
	}

	if err := sandbox.MountProc(rootFs); err != nil {
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
		Chroot: rootFs,
	}
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	secret := "sandbox-secret"
	if os.Args[1] == secret {
		if os.Getuid() != 0 || os.Getgid() != 0 {
			panic("started in namespace mode but not as root")
		}
		// we are inside the new namespace
		shell := sandbox.GetPrimaryShell()
		println("Starting sandbox with shell", shell)
		createSandboxInsideNamespace(shell)
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
	//
	// commit := parser.NewCommand("commit", "commit the previously added files to disk")
	// rm := parser.NewCommand("rm", "forcefully remove all sandbox files")

	if tryCommand.Happened() {
		sandbox.ForkSelfIntoNewNamespace(os.Args) // this will call us again with an argument
	}

}
