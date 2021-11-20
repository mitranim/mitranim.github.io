package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	tt "text/template"
	"time"
	"unsafe"

	_ "image/jpeg"
	_ "image/png"

	x "github.com/mitranim/gax"
	"github.com/mitranim/gt"
	"github.com/mitranim/try"
)

const (
	SERVER_PORT  = 52693
	PUBLIC_DIR   = "public"
	TEMPLATE_DIR = "templates"
	EMDASH       = "—"
	EMAIL        = `me@mitranim.com`
)

var (
	FLAGS   = Flags{PROD: os.Getenv("PROD") == "true"}
	E       = x.E
	F       = x.F
	A       = x.A
	AP      = x.AP
	TARBLAN = x.Attr{`target`, `_blank`}
	RELNO   = x.Attr{`rel`, `noopener noreferrer`}
	ABLAN   = A(TARBLAN, RELNO)
	MAILTO  = gt.ParseNullUrl(`mailto:` + EMAIL)
)

type (
	Bui  = x.Bui
	B    = *Bui
	Time = gt.NullTime
)

type Flags struct{ PROD bool }

func fpj(path ...string) string { return filepath.Join(path...) }

func timeNow() Time { return gt.NullTimeNow().UTC() }

func timeParse(src string) (val Time) {
	try.To(val.Parse(src))
	return
}

func timeFmtHuman(date Time) string { return date.Format("Jan 02 2006") }

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
	return strings.TrimSpace(strings.Join(reLines.Split(val, -1), ``))
}

var reLines = regexp.MustCompile(`\s*(?:\r|\n)\s*`)

func randomHex() string {
	var buf [32]byte
	_ = try.Int(rand.Read(buf[:]))
	return hex.EncodeToString(buf[:])
}

func tplToBytes(temp *tt.Template, val interface{}) (buf x.NonEscWri) {
	try.To(temp.Execute(&buf, val))
	return buf
}

func xmlEncode(val interface{}) (buf x.NonEscWri) {
	buf = append(buf, xml.Header...)
	enc := xml.NewEncoder(&buf)
	enc.Indent(``, "  ")
	try.To(enc.Encode(val))
	return buf
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
func bytesString(val []byte) string { return *(*string)(unsafe.Pointer(&val)) }

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
		if val == `` {
			continue
		}
		if out != `` {
			out += sep
		}
		out += val
	}
	return
}

func idToHash(val string) string {
	if val == `` {
		return ``
	}
	return `#` + val
}

func walkDirs(path string, fun func(string, fs.DirEntry)) {
	if fun == nil {
		return
	}

	try.To(fs.WalkDir(os.DirFS(try.String(os.Getwd())), path, func(path string, ent fs.DirEntry, err error) error {
		try.To(err)
		if ent.IsDir() {
			fun(path, ent)
		}
		return nil
	}))
}
