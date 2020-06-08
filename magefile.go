// +build mage

/*
See `readme.md` for dependencies and build commands.
*/

package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

const FS_EVENTS = notify.Create | notify.Remove | notify.Write
const SERVER_PORT = "52693"
const PUBLIC_DIR = "public"

var FLAGS = struct{ DEV bool }{DEV: os.Getenv("DEV") == "true" || os.Getenv("DEV") == ""}

var Default = Build

// Builds everything.
func Build() {
	mg.Deps(Static, Styles, Images, Html)
}

// Rebuilds, then watches and rebuilds on changes.
func Watch() {
	mg.Deps(StaticW, StylesW, ImagesW, HtmlW, Server)
}

// Removes built artifacts.
func Clean() error {
	return os.RemoveAll(PUBLIC_DIR)
}

// Copies files from ./static to the target directory.
func Static(ctx context.Context) error {
	const DIR = "static"
	return filepath.Walk(DIR, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(DIR, srcPath)
		if err != nil {
			return errors.WithStack(err)
		}

		outPath := filepath.Join(PUBLIC_DIR, rel)

		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return errors.WithStack(err)
		}

		src, err := os.Open(srcPath)
		if err != nil {
			return errors.WithStack(err)
		}
		defer src.Close()

		out, err := os.Create(outPath)
		if err != nil {
			return errors.WithStack(err)
		}
		defer out.Close()

		_, err = io.Copy(out, src)
		return errors.WithStack(err)
	})
}

// Watches static files and reruns the static task on changes.
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

			notifyClients()
		}
	}
}

// Builds styles. Requires dart-sass.
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

// Watches and rebuilds styles.
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
			notifyClients()
		}
	}
}

// Resizes and optimizes images. Requires GraphicsMagick.
//
// Uses "filepath.Walk" instead of "filepath.Glob" because the latter can't
// find everything we need in a single call.
func Images(ctx context.Context) error {
	var batch string

	err := filepath.Walk("images", func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() || !isImage(srcPath) {
			return nil
		}

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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gm", "batch", "-")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pipeIn, err := cmd.StdinPipe()
	if err != nil {
		return errors.WithStack(err)
	}

	err = cmd.Start()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = pipeIn.Write([]byte(batch))
	if err != nil {
		return errors.WithStack(err)
	}
	pipeIn.Close()

	return cmd.Wait()
}

// Watches and rebuilds images.
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

			if !isImage(srcPath) {
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

			notifyClients()
		}
	}
}

func isImage(pt string) bool {
	switch filepath.Ext(pt) {
	case ".jpg", ".png":
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

// Serves static files and notifies about file changes over websockets.
func Server() error {
	log.Println("Starting server on", "http://localhost:"+SERVER_PORT)
	return http.ListenAndServe(":"+SERVER_PORT, http.HandlerFunc(serve))
}

// Serves static files, resolving URL/HTML similarly to the default Nginx
// config, Github Pages, or Netlify.
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

	// Ends with slash -> error 404 for hygiene. Directory links must not end
	// with a slash. It's unnecessary, and GH Pages will do a 301 redirect,
	// introducing an additional delay.
	if len(reqPath) > 1 && reqPath[len(reqPath)-1] == '/' {
		goto notFound
	}

	{
		stat, _ := os.Stat(filePath)
		if fileExists(stat) {
			http.ServeFile(rew, req, filePath)
			return
		}
	}

	// Has extension -> don't bother looking for +.html or +/index.html
	if path.Ext(reqPath) != "" {
		goto notFound
		return
	}

	// Try +.html
	{
		candidatePath := filePath + ".html"
		candidateStat, _ := os.Stat(candidatePath)
		if fileExists(candidateStat) {
			http.ServeFile(rew, req, candidatePath)
			return
		}
	}

	// Try +/index.html
	{
		candidatePath := filepath.Join(filePath, "index.html")
		stat, _ := os.Stat(candidatePath)
		if fileExists(stat) {
			http.ServeFile(rew, req, candidatePath)
			return
		}
	}

notFound:
	// Sends code 200 if 404.html is found; not worth fixing for local development
	http.ServeFile(rew, req, filepath.Join(PUBLIC_DIR, "404.html"))
}

func fileExists(stat os.FileInfo) bool {
	return stat != nil && !stat.IsDir()
}

func initClientConn(rew http.ResponseWriter, req *http.Request) {
	up := websocket.Upgrader{CheckOrigin: skipOriginCheck}
	conn, err := up.Upgrade(rew, req, nil)
	if err != nil {
		log.Printf("failed to init connection at %v: %v", req.RemoteAddr, errors.WithStack(err))
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

func notifyClients() {
	CLIENTS.Range(func(_, val interface{}) bool {
		go notifyClient(val.(*Client))
		return true
	})
}

func notifyClient(client *Client) {
	client.Lock()
	defer client.Unlock()

	err := client.WriteMessage(websocket.TextMessage, nil)
	if err != nil {
		log.Printf("failed to notify socket: %+v", errors.WithStack(err))
	}
}

// Builds in "production" mode and deploys. Must not run concurrently with any
// other tasks.
func Deploy() error {
	FLAGS.DEV = false
	mg.SerialDeps(Clean, Build)

	originUrl, err := shellCmd("git", "remote", "get-url", "origin")
	if err != nil {
		return err
	}
	sourceBranch, err := shellCmd("git", "symbolic-ref", "--short", "head")
	if err != nil {
		return err
	}
	const targetBranch = "master"

	if sourceBranch == targetBranch {
		return errors.Errorf("expected source branch %q to be distinct from target branch %q",
			sourceBranch, targetBranch)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}

	return unpanic(func() {
		must(os.Chdir(PUBLIC_DIR))
		must(os.RemoveAll(".git"))
		mustShellCmd("git", "init")
		mustShellCmd("git", "remote", "add", "origin", originUrl)
		mustShellCmd("git", "add", "-A", ".")
		mustShellCmd("git", "commit", "-a", "--allow-empty-message", "-m", "")
		mustShellCmd("git", "branch", "-m", targetBranch)
		mustShellCmd("git", "push", "-f", "origin", targetBranch)
		must(os.RemoveAll(".git"))
		must(os.Chdir(cwd))
	})
}

func unpanic(fun func()) (err error) {
	defer func() {
		if err == nil {
			err, _ = recover().(error)
		}
	}()
	fun()
	return
}

func mustShellCmd(command string, args ...string) string {
	out, err := shellCmd(command, args...)
	must(err)
	return out
}

func shellCmd(command string, args ...string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), errors.WithStack(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
