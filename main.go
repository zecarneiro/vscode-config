package main

import (
	"errors"
	"jnoronha_golangutils"
	"main/libs"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		jnoronha_golangutils.ProcessError(errors.New("Invalid JSON with extensions and settings"))
	}
	libs.Start(os.Args[1])
}
