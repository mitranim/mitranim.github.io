package main

import (
	"log"
	"time"

	"github.com/mitranim/cmd"
)

var commands = cmd.Map{}

func init() {
	time.Local = nil
	log.SetFlags(0)
}

func main() {
	defer cmd.Report()
	commands.Get()()
}
