package main

import (
	"errors"
	"jnoronhautils"
	"main/libs"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		jnoronhautils.ProcessError(errors.New("Invalid JSON with extensions and settings"))
	}
	libs.Start(os.Args[1])
}
