package main

import (
	"fmt"
	"os"

	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

func main() {
	err := try.Unpanic(runMain)
	if err != nil {
		fmt.Printf("%T: %+v\n", err, err)
		os.Exit(1)
	}
}

func runMain() {
	cmd := os.Args[1]

	switch cmd {
	case "srv":
		cmdSrv()
	case "pages":
		cmdPages()
	case "images":
		cmdImages()
	case "deploy":
		cmdDeploy()
	default:
		panic(errors.Errorf(`unknown cmd %q`, cmd))
	}
}
