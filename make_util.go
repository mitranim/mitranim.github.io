package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image"
	l "log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	tt "text/template"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gotidy/ptr"
	x "github.com/mitranim/gax"
	"github.com/mitranim/try"
	"github.com/pkg/errors"
)

const (
	SERVER_PORT  = 52693
	PUBLIC_DIR   = "public"
	TEMPLATE_DIR = "templates"
	MAX_INT      = int(^uint(0) >> 1)
)

type (
	A    = x.A
	E    = x.E
	Bui  = x.Bui
	Attr = x.Attr
)

var (
	log          = l.New(os.Stderr, "", 0)
	FLAGS        = Flags{PROD: os.Getenv("PROD") == "true"}
	TARGET_BLANK = Attr{`target`, `_blank`}
	REL_NOP      = Attr{`rel`, `noopener noreferrer`}
)

type Flags struct{ PROD bool }

func fpj(path ...string) string { return filepath.Join(path...) }

func tryTime(val time.Time, err error) time.Time {
	try.To(err)
	return val
}

func timeParse(input string) (time.Time, error) {
	inst, err := time.Parse(time.RFC3339, input)
	return inst, errors.WithStack(err)
}

func tryTimePtr(val string) *time.Time {
	return ptr.Time(tryTime(timeParse(val)))
}

func timeFmtHuman(date time.Time) string {
	return date.Format("Jan 02 2006")
}

func timeCoalesce(vals ...*time.Time) *time.Time {
	for _, val := range vals {
		if val != nil && !val.IsZero() {
			return val
		}
	}
	return nil
}

func trimLeadingSlash(val string) string {
	return strings.TrimPrefix(val, "/")
}

func ensureLeadingSlash(val string) string {
	return ensurePrefix(val, "/")
}

func ensurePrefix(val, pre string) string {
	if strings.HasPrefix(val, pre) {
		return val
	}
	return pre + val
}

func trimExt(pt string) string {
	return strings.TrimSuffix(pt, filepath.Ext(pt))
}

func baseName(pt string) string {
	return trimExt(filepath.Base(pt))
}

func writePublic(path string, bytes []byte) (err error) {
	defer try.Rec(&err)

	path = fpj(PUBLIC_DIR, path)
	try.To(os.MkdirAll(filepath.Dir(path), os.ModePerm))
	try.To(os.WriteFile(path, bytes, os.ModePerm))

	return
}

func yearsElapsed() string {
	start := 2014
	now := time.Now().UTC().Year()
	if now > start {
		return fmt.Sprintf("%v—%v", start, now)
	}
	return fmt.Sprint(start)
}

func imgConfig(path string) image.Config {
	file, err := os.Open(path)
	try.To(err)
	defer file.Close()

	conf, _, err := image.DecodeConfig(file)
	try.To(err)
	return conf
}

func tryRead(path string) []byte { return try.ByteSlice(os.ReadFile(path)) }

// Inefficient but not measurable in our case.
func trimLines(val string) string {
	return strings.TrimSpace(strings.Join(reLines.Split(val, -1), ""))
}

var reLines = regexp.MustCompile(`\s*(?:\r|\n)\s*`)

func Ebui(fun func(E E)) Bui {
	var bui Bui
	fun(bui.E)
	return bui
}

func randomHex() string {
	const pre = `id`
	var buf [32]byte
	_ = try.Int(rand.Read(buf[:]))
	return hex.EncodeToString(buf[:])
}

func tplToBytes(temp *tt.Template, val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := temp.Execute(&buf, val)
	return buf.Bytes(), errors.WithStack(err)
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
process.
*/
func runCmdOut(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(bytes.TrimSpace(out)), errors.WithStack(err)
}

func walkFiles(dir string, fun func(string) error) (err error) {
	defer try.Rec(&err)

	try.To(filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		try.To(err)
		if ignorePath(path) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		try.To(fun(path))
		return nil
	}))

	return
}

func ignorePath(path string) bool {
	return strings.EqualFold(filepath.Base(path), ".DS_Store")
}

// "mkdir" is required for GraphicsMagick, which doesn't create directories.
func makeImagePath(srcPath string) (_ string, err error) {
	defer try.Rec(&err)
	rel := try.String(filepath.Rel("images", srcPath))
	outPath := fpj(PUBLIC_DIR, "images", rel)
	try.To(os.MkdirAll(filepath.Dir(outPath), os.ModePerm))
	return outPath, nil
}

func xmlEncode(input interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	err := enc.Encode(input)
	return buf.Bytes(), errors.WithStack(err)
}

func timing(name string) func() {
	start := time.Now()
	log.Printf("[%v] starting", name)

	return func() {
		end := time.Now()
		log.Printf("[%v] done in %v", name, end.Sub(start))
	}
}
