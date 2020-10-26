// +build mage

/*
See `readme.md` for dependencies and build commands.
*/

package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

type Flags struct {
	PROD bool
}

var FLAGS = Flags{
	PROD: os.Getenv("PROD") == "true",
}

// Default command for `mage` without args.
var Default = Build

// Rebuild everything.
func Build() {
	mg.Deps(Static, Styles, Images, Templates)
}

// Rebuild, then watch and rebuild on changes.
func Watch() {
	mg.Deps(StaticW, StylesW, ImagesW, TemplatesW, Server)
}

// Remove built artifacts.
func Clean() error {
	return os.RemoveAll(PUBLIC_DIR)
}

// Copy files from "./static" to the target directory.
func Static() error {
	const DIR = "static"
	return walkFiles(DIR, func(path string) error {
		rel, err := filepath.Rel(DIR, path)
		if err != nil {
			return errors.WithStack(err)
		}
		return copyFile(path, fpj(PUBLIC_DIR, rel))
	})
}

// Watch static files and rerun the static task on changes.
func StaticW() error {
	return watch(fpj("static", "..."), notify.All, func(event notify.EventInfo) {
		logger.Println("[static] FS event:", event)
		err := Static()
		if err != nil {
			logger.Println("[static] error:", err)
			return
		}
		notifyClients(nil)
	})
}

// Build styles; requires `dart-sass`.
func Styles() error {
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
func StylesW() error {
	return watch(fpj("styles", "..."), notify.All, func(event notify.EventInfo) {
		logger.Println("[styles] FS event:", event)
		err := Styles()
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
func Images() error {
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
func ImagesW() error {
	return watch(fpj("images", "..."), notify.Create|notify.Write, func(event notify.EventInfo) {
		logger.Println("[images] FS event:", event)
		err := convertImage(event.Path())
		if err != nil {
			logger.Println("[images] error:", err)
			return
		}
		notifyClients(nil)
	})
}

func convertImage(path string) (err error) {
	defer rec(&err)

	cwd, err := os.Getwd()
	must(errors.WithStack(err))

	rel, err := filepath.Rel(cwd, path)
	must(errors.WithStack(err))

	outPath, err := makeImagePath(rel)
	must(err)

	return runCmd("gm", "convert", rel, outPath)
}

// Serve static files and notify websocket clients about file changes.
func Server() error {
	const port = SERVER_PORT
	logger.Println("Starting server on", fmt.Sprintf("http://localhost:%v", port))
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

	serveFile(rew, req)
}

/*
Serves static files, resolving URL/HTML in a fashion similar to the default
Nginx config, Github Pages, and Netlify.

Note: this has a race condition between checking for a file's existence and
actually serving it. In a production-grade version, this condition should
probably be addressed. Serving a file is not an atomic operation; the file may
be deleted or changed midway. This development server doesn't need to handle
this problem.
*/
func serveFile(rew http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
	default:
		http.Error(rew, "", http.StatusMethodNotAllowed)
		return
	}

	reqPath := req.URL.Path
	filePath := fpj(PUBLIC_DIR, reqPath)

	/**
	Ends with slash? Return error 404 for hygiene. Directory links must not end
	with a slash. It's unnecessary, and GH Pages will do a 301 redirect to a
	non-slash URL, which is a good feature but adds latency.
	*/
	if len(reqPath) > 1 && reqPath[len(reqPath)-1] == '/' {
		goto notFound
	}

	if fileExists(filePath) {
		http.ServeFile(rew, req, filePath)
		return
	}

	// Has extension? Don't bother looking for +".html" or +"/index.html".
	if path.Ext(reqPath) != "" {
		goto notFound
	}

	// Try +".html".
	{
		candidatePath := filePath + ".html"
		if fileExists(candidatePath) {
			http.ServeFile(rew, req, candidatePath)
			return
		}
	}

	// Try +"/index.html".
	{
		candidatePath := fpj(filePath, "index.html")
		if fileExists(candidatePath) {
			http.ServeFile(rew, req, candidatePath)
			return
		}
	}

notFound:
	// Minor issue: sends code 200 instead of 404 if "404.html" is found; not
	// worth fixing for local development.
	http.ServeFile(rew, req, fpj(PUBLIC_DIR, "404.html"))
}

var CLIENTS sync.Map // sync.Map<string, *Client>

// Note: Gorilla's websockets support only one concurrent reader and one
// concurrent writer, and require external synchronization.
type Client struct {
	*websocket.Conn
	sync.Mutex
}

func initClientConn(rew http.ResponseWriter, req *http.Request) {
	up := websocket.Upgrader{CheckOrigin: skipOriginCheck}
	conn, err := up.Upgrade(rew, req, nil)
	if err != nil {
		logger.Printf("Failed to init connection at %v: %v", req.RemoteAddr, errors.WithStack(err))
		return
	}

	key := req.RemoteAddr
	CLIENTS.Store(key, &Client{Conn: conn})
	defer CLIENTS.Delete(key)

	// Flush and ignore client messages, if any
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func skipOriginCheck(*http.Request) bool { return true }

func notifyClients(msg []byte) {
	CLIENTS.Range(func(_, val interface{}) bool {
		go notifyClient(val.(*Client), msg)
		return true
	})
}

func notifyClient(client *Client, msg []byte) {
	client.Lock()
	defer client.Unlock()

	err := client.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		logger.Printf("Failed to notify socket: %+v", errors.WithStack(err))
	}
}

// Build in "production" mode and deploy. Stop all other tasks before running
// this!
func Deploy() (err error) {
	defer rec(&err)

	FLAGS.PROD = true
	mg.SerialDeps(Clean, Build)

	originUrl, err := runCmdOut("git", "remote", "get-url", "origin")
	must(err)

	sourceBranch, err := runCmdOut("git", "symbolic-ref", "--short", "head")
	must(err)

	const targetBranch = "master"

	if sourceBranch == targetBranch {
		return errors.Errorf("expected source branch %q to be distinct from target branch %q",
			sourceBranch, targetBranch)
	}

	cwd, err := os.Getwd()
	must(errors.WithStack(err))

	must(os.Chdir(PUBLIC_DIR))
	must(os.RemoveAll(".git"))
	must(runCmd("git", "init"))
	must(runCmd("git", "remote", "add", "origin", originUrl))
	must(runCmd("git", "add", "-A", "."))
	must(runCmd("git", "commit", "-a", "--allow-empty-message", "-m", ""))
	must(runCmd("git", "branch", "-m", targetBranch))
	must(runCmd("git", "push", "-f", "origin", targetBranch))
	must(os.RemoveAll(".git"))
	must(os.Chdir(cwd))

	return nil
}
