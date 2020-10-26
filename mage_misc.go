// +build mage

package main

import (
	"bytes"
	"encoding/xml"
	ht "html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	fpj    = filepath.Join
)

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

func logTime(prefix string, fun func() error) error {
	t0 := time.Now()

	err := fun()
	if err != nil {
		return err
	}

	t1 := time.Now()
	logger.Println(prefix, t1.Sub(t0))
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
