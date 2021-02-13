package main

import (
	"bytes"
	"encoding/xml"
	ht "html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mitranim/emptty"
	g "github.com/mitranim/gtg"
	"github.com/mitranim/try"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
)

const (
	SERVER_PORT       = 52693
	PUBLIC_DIR        = "public"
	FS_MODE_FILE      = 0666
	TEMPLATE_DIR      = "templates"
	HUMAN_TIME_FORMAT = "Jan 02 2006"
)

var (
	logger = log.New(os.Stderr, "", 0)
)

type Flags struct {
	PROD bool
}

var FLAGS = Flags{
	PROD: os.Getenv("PROD") == "true",
}

func fpj(path ...string) string { return filepath.Join(path...) }

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
	defer try.Rec(&err)

	try.To(os.MkdirAll(filepath.Dir(outPath), os.ModePerm))

	src, err := os.Open(srcPath)
	try.To(err)
	defer src.Close()

	out, err := os.Create(outPath)
	try.To(err)
	defer out.Close()

	_ = try.Int64(io.Copy(out, src))
	return nil
}

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
		logger.Printf("failed to init connection at %v: %v", req.RemoteAddr, errors.WithStack(err))
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
		logger.Printf("failed to notify socket: %+v", errors.WithStack(err))
	}
}

func skipOriginCheck(*http.Request) bool { return true }

func onFsEvent(_ g.Task, _ notify.EventInfo) {
	emptty.ClearHard()
}
