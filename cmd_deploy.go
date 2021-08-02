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

	cwd := try.String(os.Getwd())
	try.To(os.Chdir(PUBLIC_DIR))
	defer os.Chdir(cwd)

	try.To(os.RemoveAll(".git"))
	runCmd("git", "init", "-q", "-b", targetBranch)
	runCmd("git", "remote", "add", "origin", originUrl)
	runCmd("git", "add", "-A", ".")
	runCmd("git", "commit", "-q", "-a", "--allow-empty-message", "-m", "")
	runCmd("git", "branch", "-m", targetBranch)
	runCmd("git", "push", "-f", "origin", targetBranch)
	try.To(os.RemoveAll(".git"))
}
