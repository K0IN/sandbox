package sandbox

import (
	"os"
	"os/exec"
)

type NamespaceMode struct {
	AllowNetwork bool
}

func ForkSelfIntoNewNamespace(arguments []string) {
	// todo use cmd.SysProcAttr if unshare is not available
	cmd := exec.Command("unshare", "--mount", "--user", "--map-root-user", "--pid", "--fork", "--uts", os.Args[0], "mode=namespace") //, arguments[1:]...)
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	os.Exit(cmd.ProcessState.ExitCode())
}

// https://blog.scottlowe.org/2013/09/04/introducing-linux-network-namespaces/ and https://lwn.net/Articles/580893/
func AddNetwork(name string, ip string) error {
	// https://github.com/vishvananda/netlink

	// ip link add veth0 type veth peer name veth1
	// ip link set veth1 netns <pid>
	// ip netns exec <pid> ip addr add
	// ip netns exec <pid> ip link set veth1 up
	// ip netns exec <pid> ip route add default via <ip>
	return nil
}

func RemoveNetwork(name string) error {
	return nil
}
