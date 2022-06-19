package main

import (
	"log"
	"time"

	"github.com/mitranim/cmd"
	"github.com/mitranim/gg"
)

var commands = cmd.Map{}

func init() {
	time.Local = nil
	log.SetFlags(0)
	gg.TraceBaseDir = gg.Cwd()
}

func main() {
	defer cmd.Report()
	commands.Get()()
}
