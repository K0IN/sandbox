package main

import (
	"fmt"
	registry_client "myapp/registry"
	container "myapp/simple-container"
	"os"
	"syscall"
)

type TryConfig struct {
	// AllowNetwork bool
	// AllowEnv     bool
	Command  string
	Workdir  *string
	SkipDiff *bool
}

func executeTry(config TryConfig) {
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

	if config.Workdir == nil {
		currentDir, err := os.Getwd()
		if err == nil {
			config.Workdir = &currentDir
		}
	}

	exec := container.ExecConfig{
		Env:            os.Environ(),
		WorkDir:        *config.Workdir,
		Rootfs:         overlayFs.GetRootDir(),
		NameSpaceFlags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	err = container.ExecuteCommand(config.Command, &exec)
	if err != nil {
		panic(err)
	}

	// if !*config.SkipDiff {
	// 	overlayFs.ShowDiff()
	// }
}

func main() {

	client := registry_client.NewDockerRegistryClient(registry_client.RegistryBaseURL, "", "")

	image := "library/python"
	tag := "alpine"

	err := registry_client.ExtractAndAssembleImage(client, image, tag, "rootfs")
	if err != nil {
		panic(err)
	}

	err = container.ExecuteCommand("/bin/sh", &container.ExecConfig{
		Env:            os.Environ(),
		WorkDir:        "/",
		Rootfs:         "rootfs",
		NameSpaceFlags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	})

	if err != nil {
		panic(err)
	}

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

	// parser := argparse.NewParser("sandbox", "run a command in a sandbox")
	// tryCommand := parser.NewCommand("try", "execute a command inside a sandbox and review the changes")
	// // allowNet := tryCommand.Flag("n", "allow-net", &argparse.Options{Required: false, Default: true, Help: "allow network"})
	// // tryCommand.Flag("p", "allow-proc", &argparse.Options{Required: false, Default: true, Help: "allow proc"})
	// // allowEnv := tryCommand.Flag("e", "allow-env", &argparse.Options{Required: false, Default: true, Help: "allow env"})
	// workDir := tryCommand.String("c", "workdir", &argparse.Options{Required: false, Help: "workdir"})
	// args := parser.StringPositional(&argparse.Options{Required: true, Help: "command to execute"})
	// skipDiff := tryCommand.Flag("s", "skip-diff", &argparse.Options{Required: false, Default: false, Help: "skip diff"})
	//
	// // c2 := parser.NewCommand("container", "start a ocid container in the most cracked way possible")
	//
	// err := parser.Parse(os.Args)
	// if err != nil {
	// 	fmt.Print(parser.Usage(err))
	// }
	//
	// if tryCommand.Happened() {
	// 	config := TryConfig{
	// 		// AllowNetwork: *allowNet,
	// 		// AllowEnv:     *allowEnv,
	// 		Command:  *args,
	// 		Workdir:  workDir,
	// 		SkipDiff: skipDiff,
	// 	}
	// 	executeTry(config)
	// }
	//
	// // Finally print the collected string
	// // fmt.Println(*s)
	// //
	// //
	// //
	// // argsWithProg := os.Args
	// // argsWithoutProg := os.Args[1:]
	// // if len(argsWithoutProg) < 1 {
	// // 	fmt.Printf("Usage: %s <command>\n", argsWithProg[0])
	// // 	fmt.Printf("Example: %s ls\n", argsWithProg[0])
	// // 	os.Exit(1)
	// // }
	// //
	// // sandboxDir, err := os.MkdirTemp("", "sandbox")
	// // if err != nil {
	// // 	panic(fmt.Errorf("cannot create sandbox dir: %w", err))
	// // }
	// //
	// // fmt.Println("Created sandbox:", sandboxDir)
	// //
	// // fs, _ := overlayfs.CreateOverlayFs(sandboxDir) // todo -> use fuse instead of overlayfs
	// // err = fs.Mount()
	// // if err != nil {
	// // 	panic(fmt.Errorf("cannot mount overlayfs: %w", err))
	// // }
	// //
	// // defer func() {
	// // 	fs.Unmount()
	// // 	_ = os.RemoveAll(sandboxDir)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	//
	// // execute the command
	// err = ExecuteCommand(argsWithoutProg, fs.GetRootDir())
	// if err != nil {
	// 	panic(err)
	// }

	// show diff
	// fs.ShowDiff()
}
