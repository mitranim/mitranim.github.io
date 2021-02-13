/*
See `readme.md` for dependencies and build commands.
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	g "github.com/mitranim/gtg"
	"github.com/mitranim/srv"
	"github.com/mitranim/try"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

func main() {
	g.MustRunCmd(Build, Watch, Deploy)
}

// Rebuild everything.
func Build(task g.Task) error {
	defer g.TaskTiming(Build)()
	return g.Wait(task, g.Ser(Clean, g.Par(Static, Styles, Images, Pages)))
}

// Rebuild, then watch and rebuild on changes.
func Watch(task g.Task) error {
	return g.Wait(task, g.Ser(Clean, g.Par(StaticW, StylesW, ImagesW, PagesW, Server)))
}

// Remove built artifacts.
func Clean(task g.Task) error {
	defer g.TaskTiming(Clean)()
	return os.RemoveAll(PUBLIC_DIR)
}

// Copy files from "./static" to the target directory.
func Static(task g.Task) error {
	defer g.TaskTiming(Static)()

	const DIR = "static"

	return walkFiles(DIR, func(path string) error {
		rel := try.String(filepath.Rel(DIR, path))
		return copyFile(path, fpj(PUBLIC_DIR, rel))
	})
}

// Watch static files and rerun the static task on changes.
func StaticW(task g.Task) error {
	g.MustWait(task, g.Opt(Static))

	return watch(fpj("static", "..."), notify.All, func(event notify.EventInfo) {
		onFsEvent(task, event)
		err := Static(task)
		if err != nil {
			logger.Println("[static] error:", err)
			return
		}
		notifyClients(nil)
	})
}

// Build styles; requires `dart-sass`.
func Styles(task g.Task) error {
	defer g.TaskTiming(Styles)()

	var style string
	if FLAGS.PROD {
		style = "compressed"
	} else {
		style = "expanded"
	}

	return runCmd("sass",
		"--no-source-map",
		"--style",
		style,
		fpj("styles", "main.scss"),
		fpj(PUBLIC_DIR, "styles", "main.css"))
}

/*
Watch and rebuild styles.

The reason we don't use Sass's "--watch" option is because we run
`notifyClients` on successful rebuilds, which would be relatively hard to detect
from the subcommand. We can't just assume "any output to stdout = success",
because on errors it outputs to BOTH stdout and stderr. It's simpler and more
reliable to do our own watching. Fortunately, the command is fast enough for our
purposes.
*/
func StylesW(task g.Task) error {
	g.MustWait(task, g.Opt(Styles))

	return watch(fpj("styles", "..."), notify.All, func(event notify.EventInfo) {
		onFsEvent(task, event)
		err := Styles(task)
		if err != nil {
			logger.Println("[styles] error:", err)
			return
		}
		notifyClients(nil)
	})
}

/*
Resize and optimize images; requires GraphicsMagick.

Doesn't use "filepath.Glob" because the latter can't find everything we need in
a single call.
*/
func Images(task g.Task) error {
	defer g.TaskTiming(Images)()

	var batch string

	err := walkFiles("images", func(srcPath string) error {
		outPath, err := makeImagePath(srcPath)
		if err != nil {
			return err
		}
		batch += "convert " + srcPath + " " + outPath + "\n"
		return nil
	})

	if err != nil {
		return err
	}

	if batch == "" {
		return nil
	}

	cmd := makeCmd("gm", "batch", "-")
	cmd.Stdin = strings.NewReader(batch)
	return errors.WithStack(cmd.Run())
}

// Watch and rebuild images.
func ImagesW(task g.Task) error {
	g.MustWait(task, g.Opt(Images))

	return watch(fpj("images", "..."), notify.Create|notify.Write, func(event notify.EventInfo) {
		onFsEvent(task, event)
		err := convertImage(event.Path())
		if err != nil {
			logger.Println("[images] error:", err)
			return
		}
		notifyClients(nil)
	})
}

func convertImage(path string) (err error) {
	defer try.Rec(&err)

	cwd := try.String(os.Getwd())
	rel := try.String(filepath.Rel(cwd, path))
	outPath := try.String(makeImagePath(rel))

	return runCmd("gm", "convert", rel, outPath)
}

// Serve static files and notify websocket clients about file changes.
func Server(_ g.Task) error {
	const port = SERVER_PORT
	logger.Println("[Server] starting on", fmt.Sprintf("http://localhost:%v", port))
	return http.ListenAndServe(fmt.Sprintf(":%v", port), http.HandlerFunc(serve))
}

func serve(rew http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/ws" {
		if req.Method != http.MethodGet {
			rew.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		initClientConn(rew, req)
		return
	}

	srv.FileServer(PUBLIC_DIR).ServeHTTP(rew, req)
}

// Build in "production" mode and deploy. Stop all other tasks before running
// this!
func Deploy(task g.Task) error {
	FLAGS.PROD = true
	g.MustWait(task, Clean)
	g.MustWait(task, Build)

	defer g.TaskTiming(Deploy)()

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
