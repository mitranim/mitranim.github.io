package main

import (
	"strings"

	"github.com/mitranim/try"
)

/*
Resize and optimize images; requires GraphicsMagick.

Doesn't use "filepath.Glob" because the latter can't find everything we need in
a single call.
*/
func cmdImages() {
	defer timing("images")()

	var batch string

	walkFiles("images", func(srcPath string) {
		outPath := makeImagePath(srcPath)
		batch += "convert " + srcPath + " " + outPath + "\n"
	})

	if batch == "" {
		return
	}

	cmd := makeCmd("gm", "batch", "-")
	cmd.Stdin = strings.NewReader(batch)
	try.To(cmd.Run())
}
