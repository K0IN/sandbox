package helper

import "os"

func IsOutputPiped() bool {
	outputStat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	isPiped := outputStat.Mode()&os.ModeCharDevice == 0
	return isPiped
}
