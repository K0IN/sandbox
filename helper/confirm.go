package helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Confirm(s string) bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", s)
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if len(res) < 2 {
		return false
	}
	return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
}
