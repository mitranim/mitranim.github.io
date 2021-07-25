package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	err := runMain(os.Args[1])
	if err != nil {
		fmt.Printf("%T: %+v\n", err, err)
		os.Exit(1)
	}
}

func runMain(cmd string) error {
	switch cmd {
	case "srv":
		return cmdSrv()
	case "pages":
		return cmdPages()
	case "images":
		return cmdImages()
	case "deploy":
		return cmdDeploy()
	default:
		return errors.Errorf(`unknown cmd %q`, cmd)
	}
}
