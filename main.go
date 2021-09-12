package main

import (
	"time"

	"github.com/mitranim/cmd"
)

var commands = cmd.Map{}

func init() {
	time.Local = nil
}

func main() {
	defer cmd.Report()
	commands.Get()()
}
