package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image"
	"io"
	l "log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	tt "text/template"
	"time"
	"unsafe"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gotidy/ptr"
	x "github.com/mitranim/gax"
	"github.com/mitranim/try"
)

const (
	SERVER_PORT  = 52693
	PUBLIC_DIR   = "public"
	TEMPLATE_DIR = "templates"
	MAX_INT      = int(^uint(0) >> 1)
	EMDASH       = "—"
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
	Ebui         = x.Ebui
)

type Flags struct{ PROD bool }

func fpj(path ...string) string { return filepath.Join(path...) }

func timeParse(input string) time.Time {
	inst, err := time.Parse(time.RFC3339, input)
	try.To(err)
	return inst
}

func timeParsePtr(val string) *time.Time { return ptr.Time(timeParse(val)) }
func timeFmtHuman(date time.Time) string { return date.Format("Jan 02 2006") }

func timeCoalesce(vals ...*time.Time) *time.Time {
	for _, val := range vals {
		if val != nil && !val.IsZero() {
			return val
		}
	}
	return nil
}

func trimLeadingSlash(val string) string   { return strings.TrimPrefix(val, "/") }
func ensureLeadingSlash(val string) string { return ensurePrefix(val, "/") }

func ensurePrefix(val, pre string) string {
	if strings.HasPrefix(val, pre) {
		return val
	}
	return pre + val
}

func trimExt(pt string) string  { return strings.TrimSuffix(pt, filepath.Ext(pt)) }
func baseName(pt string) string { return trimExt(filepath.Base(pt)) }

func writePublic(path string, bytes []byte) {
	path = fpj(PUBLIC_DIR, path)
	try.To(os.MkdirAll(filepath.Dir(path), os.ModePerm))
	try.To(os.WriteFile(path, bytes, os.ModePerm))
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

func readFile(path string) []byte { return try.ByteSlice(os.ReadFile(path)) }

// Inefficient but not measurable in our case.
func trimLines(val string) string {
	return strings.TrimSpace(strings.Join(reLines.Split(val, -1), ""))
}

var reLines = regexp.MustCompile(`\s*(?:\r|\n)\s*`)

func randomHex() string {
	var buf [32]byte
	_ = try.Int(rand.Read(buf[:]))
	return hex.EncodeToString(buf[:])
}

func tplToBytes(temp *tt.Template, val interface{}) []byte {
	var buf bytes.Buffer
	try.To(temp.Execute(&buf, val))
	return buf.Bytes()
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
func runCmd(command string, args ...string) {
	try.To(makeCmd(command, args...).Run())
}

/*
Runs a command and returns its stdout. Stderr is connected to the parent
process.
*/
func runCmdOut(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	return bytesToMutableString(bytes.TrimSpace(try.ByteSlice(cmd.Output())))
}

func walkFiles(dir string, fun func(string)) {
	try.To(filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		try.To(err)
		if ignorePath(path) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fun(path)
		return nil
	}))
}

func ignorePath(path string) bool {
	return strings.EqualFold(filepath.Base(path), ".DS_Store")
}

// "mkdir" is required for GraphicsMagick, which doesn't create directories.
func makeImagePath(srcPath string) string {
	rel := try.String(filepath.Rel("images", srcPath))
	outPath := fpj(PUBLIC_DIR, "images", rel)
	try.To(os.MkdirAll(filepath.Dir(outPath), os.ModePerm))
	return outPath
}

func xmlEncode(input interface{}) []byte {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	try.To(enc.Encode(input))
	return buf.Bytes()
}

func timing(name string) func() {
	start := time.Now()
	log.Printf("[%v] starting", name)

	return func() {
		log.Printf("[%v] done in %v", name, time.Since(start))
	}
}

//nolint:unused,deadcode
func withTiming(str string, fun func()) {
	defer timing(str)()
	fun()
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe.
*/
func bytesToMutableString(val []byte) string  { return *(*string)(unsafe.Pointer(&val)) }
func stringToBytesAlloc(val string) []byte    { return []byte(val) }
func ioWrite(out io.Writer, val []byte)       { try.Int(out.Write(val)) }
func ioWriteString(out io.Writer, val string) { try.Int(io.WriteString(out, val)) }

func parseUrl(val string) Url {
	out, err := url.Parse(val)
	try.To(err)
	return Url(*out)
}

// Unfucks `*url.URL` by making it a non-pointer. TODO move to separate lib.
type Url url.URL

func (self Url) Query() url.Values               { return (*url.URL)(&self).Query() }
func (self Url) String() string                  { return (*url.URL)(&self).String() }
func (self Url) MarshalText() ([]byte, error)    { return (*url.URL)(&self).MarshalBinary() }
func (self *Url) UnmarshalText(val []byte) error { return (*url.URL)(self).UnmarshalBinary(val) }

func strJoin(sep string, vals ...string) (out string) {
	for _, val := range vals {
		if val == "" {
			continue
		}
		if out != "" {
			out += sep
		}
		out += val
	}
	return
}

func chdir(val string) { try.To(os.Chdir(val)) }
