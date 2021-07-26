package main

import (
	"os"

	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

// Stop all other tasks before running this!
func cmdDeploy() {
	defer timing(`deploy`)()

	originUrl := runCmdOut("git", "remote", "get-url", "origin")
	sourceBranch := runCmdOut("git", "symbolic-ref", "--short", "head")
	const targetBranch = "master"

	if sourceBranch == targetBranch {
		panic(errors.Errorf("expected source branch %q to be distinct from target branch %q",
			sourceBranch, targetBranch))
	}

	try.To(os.Chdir(PUBLIC_DIR))
	try.To(os.RemoveAll(".git"))
	runCmd("git", "init")
	runCmd("git", "remote", "add", "origin", originUrl)
	runCmd("git", "add", "-A", ".")
	runCmd("git", "commit", "-a", "--allow-empty-message", "-m", "")
	runCmd("git", "branch", "-m", targetBranch)
	runCmd("git", "push", "-f", "origin", targetBranch)
	try.To(os.RemoveAll(".git"))
	try.To(os.Chdir(try.String(os.Getwd())))
}
