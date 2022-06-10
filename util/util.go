package util

import (
	"fmt"
	"os"
)

func ErrAndExit(err error) {
	fmt.Println("error: " + err.Error())
	os.Exit(1)
}

func PrintAndExit(f string, args ...any) {
	fmt.Println(fmt.Sprintf(f, args...))
	os.Exit(1)
}
