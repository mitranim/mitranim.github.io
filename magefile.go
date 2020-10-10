// +build mage

/*
See `readme.md` for dependencies and build commands.
*/

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

const FS_EVENTS = notify.Create | notify.Remove | notify.Write

// We could avoid this by acquiring a random port and reporting it to the
// terminal, but a consistent port is more convenient for developing a website.
const SERVER_PORT = 52693

const PUBLIC_DIR = "public"

type Flags struct {
	DEV bool
}

var FLAGS = Flags{
	DEV: os.Getenv("DEV") == "true" || os.Getenv("DEV") == "",
}

var Default = Build

// Rebuild everything.
func Build() {
	mg.Deps(Static, Styles, Images, Html)
}

// Rebuild, then watch and rebuild on changes.
func Watch() {
	mg.Deps(StaticW, StylesW, ImagesW, HtmlW, Server)
}

// Remove built artifacts.
func Clean() error {
	return os.RemoveAll(PUBLIC_DIR)
}

// Copy files from "./static" to the target directory.
func Static(ctx context.Context) error {
	const DIR = "static"

	err := filepath.Walk(DIR, func(srcPath string, info os.FileInfo, fsErr error) (err error) {
		defer rec(&err)
		must(errors.WithStack(fsErr))

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(DIR, srcPath)
		must(errors.WithStack(err))

		outPath := filepath.Join(PUBLIC_DIR, rel)

		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		must(errors.WithStack(err))

		src, err := os.Open(srcPath)
		must(errors.WithStack(err))
		defer src.Close()

		out, err := os.Create(outPath)
		must(errors.WithStack(err))
		defer out.Close()

		_, err = io.Copy(out, src)
		must(errors.WithStack(err))
		return nil
	})
	return errors.WithStack(err)
}

// Watch static files and rerun the static task on changes.
func StaticW(ctx context.Context) error {
	fsEvents := make(chan notify.EventInfo, 1)
	err := notify.Watch("static/...", fsEvents, FS_EVENTS)
	if err != nil {
		return errors.WithStack(err)
	}
	defer notify.Stop(fsEvents)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fsEvent := <-fsEvents:
			info, err := os.Stat(fsEvent.Path())
			if err != nil {
				log.Println("[static] error:", errors.WithStack(err))
				continue
			}

			if info.IsDir() {
				continue
			}

			log.Println("[static] FS event:", fsEvent)

			err = Static(ctx)
			if err != nil {
				log.Println("[static] error:", errors.WithStack(err))
				continue
			}

			notifyClients(nil)
		}
	}
}

// Build styles; requires `dart-sass`.
func Styles(ctx context.Context) error {
	const CMD = "sass"
	args := []string{"--no-source-map"}
	if !FLAGS.DEV {
		args = append(args, "--style=compressed")
	}
	args = append(args, "styles/main.scss", PUBLIC_DIR+"/styles/main.css")

	cmd := exec.CommandContext(ctx, CMD, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Watch and rebuild styles.
//
// This could use Sass's "--watch" option, but on errors, it outputs to both
// stdout and stderr. We'd have to read from them concurrently, and use brittle
// timing-sensitive synchronization to detect errors. It's simpler and more
// reliable to do our own watching. Fortunately, the command is fast enough for
// our purposes.
func StylesW(ctx context.Context) error {
	fsEvents := make(chan notify.EventInfo, 1)
	err := notify.Watch("styles/...", fsEvents, FS_EVENTS)
	if err != nil {
		return errors.WithStack(err)
	}
	defer notify.Stop(fsEvents)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fsEvent := <-fsEvents:
			log.Println("[styles] FS event:", fsEvent)
			err = Styles(ctx)
			if err != nil {
				log.Println("[styles] error:", errors.WithStack(err))
				continue
			}
			notifyClients(nil)
		}
	}
}

// Resize and optimize images; requires GraphicsMagick.
//
// Uses "filepath.Walk" instead of "filepath.Glob" because the latter can't
// find everything we need in a single call.
func Images(ctx context.Context) (err error) {
	defer rec(&err)
	var batch string

	err = filepath.Walk("images", func(srcPath string, info os.FileInfo, fsErr error) (err error) {
		defer rec(&err)
		must(errors.WithStack(fsErr))

		if info.IsDir() || !isImagePath(srcPath) {
			return nil
		}

		outPath, err := makeImagePath(srcPath)
		must(err)

		batch += "convert " + srcPath + " " + outPath + "\n"
		return nil
	})
	must(errors.WithStack(err))

	if batch == "" {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gm", "batch", "-")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pipeIn, err := cmd.StdinPipe()
	must(errors.WithStack(err))

	err = cmd.Start()
	must(errors.WithStack(err))

	_, err = pipeIn.Write([]byte(batch))
	must(errors.WithStack(err))
	must(errors.WithStack(pipeIn.Close()))

	err = cmd.Wait()
	must(errors.WithStack(err))
	return nil
}

// Watch and rebuild images.
func ImagesW(ctx context.Context) error {
	fsEvents := make(chan notify.EventInfo, 1)
	err := notify.Watch("images/...", fsEvents, FS_EVENTS)
	if err != nil {
		return errors.WithStack(err)
	}
	defer notify.Stop(fsEvents)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fsEvent := <-fsEvents:
			absPath := fsEvent.Path()

			cwd, err := os.Getwd()
			if err != nil {
				log.Println("[images] error:", errors.WithStack(err))
				continue
			}

			srcPath, err := filepath.Rel(cwd, absPath)
			if err != nil {
				log.Println("[images] error:", errors.WithStack(err))
				continue
			}

			if !isImagePath(srcPath) {
				continue
			}

			log.Println("[images] FS event:", fsEvent)

			outPath, err := makeImagePath(srcPath)
			if err != nil {
				log.Println("[images] error:", errors.WithStack(err))
				continue
			}

			cmd := exec.CommandContext(ctx, "gm", "convert", srcPath, outPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Println("[images] error:", errors.WithStack(err))
				continue
			}

			notifyClients(nil)
		}
	}
}

func isImagePath(pt string) bool {
	switch filepath.Ext(pt) {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

// "mkdir" is required for GraphicsMagick, which doesn't create directories.
func makeImagePath(srcPath string) (string, error) {
	rel, err := filepath.Rel("images", srcPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	outPath := filepath.Join(PUBLIC_DIR, "images", rel)
	return outPath, os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
}

var CLIENTS sync.Map // sync.Map<string, *Client>

// Note: Gorilla's websockets support only one concurrent reader and one
// concurrent writer, and require external synchronization.
type Client struct {
	*websocket.Conn
	sync.Mutex
}

// Serve static files and notify websocket clients about file changes.
func Server() error {
	const port = SERVER_PORT
	log.Println("Starting server on", fmt.Sprintf("http://localhost:%v", port))
	return http.ListenAndServe(fmt.Sprintf(":%v", port), http.HandlerFunc(serve))
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
func serve(rew http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/ws" {
		if req.Method != http.MethodGet {
			rew.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		initClientConn(rew, req)
		return
	}

	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
	default:
		http.Error(rew, "", http.StatusMethodNotAllowed)
		return
	}

	reqPath := req.URL.Path
	filePath := filepath.Join(PUBLIC_DIR, reqPath)

	// Ends with slash? Return error 404 for hygiene. Directory links must not end
	// with a slash. It's unnecessary, and GH Pages will do a 301 redirect to a
	// non-slash URL, which is a good feature but adds latency.
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
		candidatePath := filepath.Join(filePath, "index.html")
		if fileExists(candidatePath) {
			http.ServeFile(rew, req, candidatePath)
			return
		}
	}

notFound:
	// Minor issue: sends code 200 instead of 404 if "404.html" is found; not
	// worth fixing for local development.
	http.ServeFile(rew, req, filepath.Join(PUBLIC_DIR, "404.html"))
}

func fileExists(filePath string) bool {
	stat, _ := os.Stat(filePath)
	return stat != nil && !stat.IsDir()
}

func initClientConn(rew http.ResponseWriter, req *http.Request) {
	up := websocket.Upgrader{CheckOrigin: skipOriginCheck}
	conn, err := up.Upgrade(rew, req, nil)
	if err != nil {
		log.Printf("Failed to init connection at %v: %v", req.RemoteAddr, errors.WithStack(err))
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
		log.Printf("Failed to notify socket: %+v", errors.WithStack(err))
	}
}

// Build in "production" mode and deploy. Stop all other tasks before running
// this!
func Deploy() (err error) {
	defer rec(&err)

	FLAGS.DEV = false
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

/*
Runs a command for side effects, connecting its stdout and stderr to the parent
process.
*/
func runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return errors.WithStack(err)
}

/*
Runs a command and returns its stdout. Stderr is connected to the parent
process. TODO: should this also return stderr?
*/
func runCmdOut(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(bytes.TrimSpace(out)), errors.WithStack(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func rec(ptr *error) {
	val := recover()
	if val == nil {
		return
	}

	recErr, ok := val.(error)
	if ok {
		*ptr = recErr
		return
	}

	panic(val)
}
