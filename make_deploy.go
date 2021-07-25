package main

import (
	"os"

	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

// Stop all other tasks before running this!
func cmdDeploy() error {
	defer timing(`deploy`)()

	originUrl := try.String(runCmdOut("git", "remote", "get-url", "origin"))
	sourceBranch := try.String(runCmdOut("git", "symbolic-ref", "--short", "head"))
	const targetBranch = "master"

	if sourceBranch == targetBranch {
		return errors.Errorf("expected source branch %q to be distinct from target branch %q",
			sourceBranch, targetBranch)
	}

	try.To(os.Chdir(PUBLIC_DIR))
	try.To(os.RemoveAll(".git"))
	try.To(runCmd("git", "init"))
	try.To(runCmd("git", "remote", "add", "origin", originUrl))
	try.To(runCmd("git", "add", "-A", "."))
	try.To(runCmd("git", "commit", "-a", "--allow-empty-message", "-m", ""))
	try.To(runCmd("git", "branch", "-m", targetBranch))
	try.To(runCmd("git", "push", "-f", "origin", targetBranch))
	try.To(os.RemoveAll(".git"))
	try.To(os.Chdir(try.String(os.Getwd())))

	return nil
}
