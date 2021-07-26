package main

import (
	"strings"

	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

/*
Resize and optimize images; requires GraphicsMagick.

Doesn't use "filepath.Glob" because the latter can't find everything we need in
a single call.
*/
func cmdImages() (err error) {
	defer try.Rec(&err)
	defer timing("images")()

	var batch string

	try.To(walkFiles("images", func(srcPath string) error {
		outPath := try.String(makeImagePath(srcPath))
		batch += "convert " + srcPath + " " + outPath + "\n"
		return nil
	}))

	if batch == "" {
		return
	}

	cmd := makeCmd("gm", "batch", "-")
	cmd.Stdin = strings.NewReader(batch)
	return errors.WithStack(cmd.Run())
}
