package helper

import (
	"os"
	"strconv"
)

func GetOriginalUserId() (uint32, uint32) {
	if os.Getgid() != 0 && os.Getuid() != 0 {
		return uint32(os.Getuid()), uint32(os.Getgid())
	}

	uid := os.Getenv("SUDO_UID")
	gid := os.Getenv("SUDO_GID")

	// todo check for doas as well

	// convert string to int atoi
	uidInt, _ := strconv.Atoi(uid)
	gidInt, _ := strconv.Atoi(gid)

	return uint32(uidInt), uint32(gidInt)
}
