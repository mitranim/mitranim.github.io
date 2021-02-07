package main

import (
	"bytes"
	"encoding/xml"
	ht "html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	g "github.com/mitranim/gtg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

const (
	SERVER_PORT           = 52693
	PUBLIC_DIR            = "public"
	FS_MODE_FILE          = 0666
	TEMPLATE_DIR          = "templates"
	HUMAN_TIME_FORMAT     = "Jan 02 2006"
	ESC                   = "\x1b"
	TERM_CLEAR_SOFT       = ESC + "c"
	TERM_CLEAR_SCROLLBACK = ESC + "[3J"
	TERM_CLEAR_HARD       = TERM_CLEAR_SOFT + TERM_CLEAR_SCROLLBACK
)

var (
	info = log.New(os.Stderr, "", 0)
	verb = log.New(os.Stderr, "", 0)
	fpj  = filepath.Join
)

type Flags struct {
	PROD bool
	VERB bool
}

var FLAGS = Flags{
	PROD: os.Getenv("PROD") == "true",
	VERB: os.Getenv("VERB") == "" || os.Getenv("VERB") == "true",
}

func init() {
	if !FLAGS.VERB {
		verb.SetOutput(ioutil.Discard)
	}
}

// "mkdir" is required for GraphicsMagick, which doesn't create directories.
func makeImagePath(srcPath string) (string, error) {
	rel, err := filepath.Rel("images", srcPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	outPath := fpj(PUBLIC_DIR, "images", rel)

	err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return outPath, nil
}

func fileExists(filePath string) bool {
	stat, _ := os.Stat(filePath)
	return stat != nil && !stat.IsDir()
}

func ignorePath(path string) bool {
	return strings.EqualFold(filepath.Base(path), ".DS_Store")
}

func makeCmd(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

/*
Runs a command for side effects, connecting its stdout and stderr to the parent
process.
*/
func runCmd(command string, args ...string) error {
	cmd := makeCmd(command, args...)
	return errors.WithStack(cmd.Run())
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

func walkFiles(dir string, fun func(string) error) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if ignorePath(path) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		return fun(path)
	})
}

func watch(pattern string, types notify.Event, fun func(event notify.EventInfo)) error {
	events := make(chan notify.EventInfo, 1)
	err := notify.Watch(pattern, events, types)
	if err != nil {
		return errors.WithStack(err)
	}

	for event := range events {
		if ignorePath(event.Path()) {
			continue
		}
		fun(event)
	}
	return nil
}

func copyFile(srcPath string, outPath string) (err error) {
	defer rec(&err)

	err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	must(errors.WithStack(err))

	src, err := os.Open(srcPath)
	must(errors.WithStack(err))
	defer src.Close()

	out, err := os.Create(outPath)
	must(errors.WithStack(err))
	defer out.Close()

	_, err = io.Copy(out, src)
	return errors.WithStack(err)
}

func withTiming(name string, fun func() error) error {
	defer timing(name)()
	return fun()
}

func taskTiming(fun g.TaskFunc) func() {
	if FLAGS.VERB {
		// return g.Timing(task)
		return timing(fun.ShortName())
	}
	return noop
}

func timing(name string) func() {
	if !FLAGS.VERB {
		return noop
	}

	t0 := time.Now()
	verb.Printf("[%v] starting\n", name)

	return func() {
		t1 := time.Now()
		verb.Printf("[%v] done in %v\n", name, t1.Sub(t0))
	}
}

func noop() {}

func globs(patterns ...string) ([]string, error) {
	var out []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return out, errors.WithStack(err)
		}
		out = append(out, matches...)
	}
	return out, nil
}

func xmlEncode(input interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	err := enc.Encode(input)
	return buf.Bytes(), errors.WithStack(err)
}

func renderTemplate(temp *ht.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := temp.Execute(&buf, data)
	return buf.Bytes(), errors.WithStack(err)
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
	sync.Mutex
	*websocket.Conn
}

func initClientConn(rew http.ResponseWriter, req *http.Request) {
	up := websocket.Upgrader{CheckOrigin: skipOriginCheck}
	conn, err := up.Upgrade(rew, req, nil)
	if err != nil {
		info.Printf("failed to init connection at %v: %v", req.RemoteAddr, errors.WithStack(err))
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
		info.Printf("failed to notify socket: %+v", errors.WithStack(err))
	}
}

func skipOriginCheck(*http.Request) bool { return true }

func clearTerminal() {
	os.Stdout.Write([]byte(TERM_CLEAR_HARD))
}

func onFsEvent(task g.Task, event notify.EventInfo) {
	clearTerminal()
	// verb.Printf("[%v] FS event: %v", task.TaskFunc().ShortName(), event)
}
