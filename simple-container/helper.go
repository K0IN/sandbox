package container

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func lookupUser() (*user.User, error) {
	// Check if SUDO_USER or DOAS_USER is set. These are typically set by sudo and doas respectively.
	for _, envVar := range []string{"SUDO_USER", "DOAS_USER"} {
		if username, found := os.LookupEnv(envVar); found {
			return user.Lookup(username)
		}
	}

	// If none of the environment variables were set, and we're running as root,
	// then we might be running under setuid. Get the real UID and retrieve user info.
	if syscall.Geteuid() == 0 {
		realUID := syscall.Getuid()
		return user.LookupId(fmt.Sprint(realUID))
	}

	// Fallback to current effective user.
	currentUID := fmt.Sprint(syscall.Geteuid())
	return user.LookupId(currentUID)
}

func GetOriginalUser() (uid int, gid int, err error) {
	user, err := lookupUser()
	if err != nil {
		return 0, 0, err
	}

	uid, err = strconv.Atoi(user.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err = strconv.Atoi(user.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func GetPrimaryShell() string {
	shell := "/bin/sh"
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		shell = "/bin/sh"
	}
	return shell
}
