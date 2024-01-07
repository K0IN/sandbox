package main

import (
	"fmt"
	registry_client "myapp/registry"
	container "myapp/simple-container"
	"os"
	"strings"
	"syscall"

	"github.com/akamensky/argparse"
)

// Remove the TryConfig struct and replace it with individual parameters in the executeTry function.
func executeTry(command string, workdir *string, skipDiff *bool) {
	if command == "" {
		command = os.Getenv("SHELL")
		if command == "" {
			command = "/bin/sh"
		}
	}

	sandboxDir, err := os.MkdirTemp("", "sandbox")
	if err != nil {
		panic(fmt.Errorf("cannot create sandbox dir: %w", err))
	}

	overlayFs, err := container.CreateOverlayFs(sandboxDir)
	defer overlayFs.Unmount()
	if err != nil {
		panic(fmt.Errorf("cannot create overlayfs: %w", err))
	}

	err = overlayFs.Mount()
	if err != nil {
		panic(fmt.Errorf("cannot mount overlayfs: %w", err))
	}

	if workdir == nil {
		currentDir, err := os.Getwd()
		if err == nil {
			workdir = &currentDir
		}
	}

	exec := container.ExecConfig{
		Env:            os.Environ(),
		WorkDir:        *workdir,
		Rootfs:         overlayFs.GetRootDir(),
		NameSpaceFlags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	err = container.ExecuteCommand(command, &exec)
	if err != nil {
		panic(err)
	}

	if !*skipDiff {
		// 	overlayFs.ShowDiff()
	}
}

func executeContainer(imageWithTag string) {
	client := registry_client.NewDockerRegistryClient(registry_client.RegistryBaseURL, "", "")
	imageWithTagSplit := strings.Split(imageWithTag, ":")
	if len(imageWithTagSplit) != 2 {
		panic(fmt.Errorf("invalid image with tag: %s", imageWithTag))
	}

	image := imageWithTagSplit[0]
	tag := imageWithTagSplit[1]

	destination, err := os.MkdirTemp("", "rootfs")
	if err != nil {
		panic(fmt.Errorf("cannot create container rootfs dir: %w", err))
	}

	fmt.Printf("starting image: %s, tag: %s, destination: %s\n", image, tag, destination)

	config, err := registry_client.ExtractAndAssembleImage(client, image, tag, destination)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Config:", config)

	cmd := strings.Join(config.Config.Cmd, "")
	err = container.ExecuteCommand(cmd, &container.ExecConfig{
		Env:            config.Config.Env,
		WorkDir:        "/",
		Rootfs:         destination,
		NameSpaceFlags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	//arguments
	// --allow-network
	// --allow-proc
	// --allow-env
	// --rootfs
	// --entrypoint
	// --env
	// --volume
	// --workdir
	// --user
	// --group
	// --hostname
	// --mount OR --overlay
	// --mount-proc
	// --mount-dev
	fmt.Printf("user id %d\n", os.Getuid())
	fmt.Printf("group id %d\n", os.Getgid())
	parser := argparse.NewParser("sandbox", "run a command in a sandbox")

	tryCommand := parser.NewCommand("try", "execute a command inside a sandbox and review the changes")
	workDir := tryCommand.String("c", "workdir", &argparse.Options{Required: false, Help: "workdir"})
	args := parser.StringPositional(&argparse.Options{Required: true, Help: "command to execute"})
	skipDiff := tryCommand.Flag("s", "skip-diff", &argparse.Options{Required: false, Default: false, Help: "skip diff"})

	containerCommand := parser.NewCommand("container", "start a container in the most cracked way possible (please note that this is just chroot with a custom namespace, no overlayfs)")
	imageWithTag := containerCommand.String("i", "image", &argparse.Options{Required: true, Help: "image with tag, example: library/python:latest"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if tryCommand.Happened() {

		executeTry(*args, workDir, skipDiff)
	} else if containerCommand.Happened() {
		executeContainer(*imageWithTag)
	} else {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
}
