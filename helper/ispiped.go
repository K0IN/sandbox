package helper

import "os"

func IsOutputPiped() (bool, error) {
	outputStat, err := os.Stdout.Stat()
	if err != nil {
		return false, err
	}
	isPiped := outputStat.Mode()&os.ModeCharDevice == 0
	return isPiped, nil
}
