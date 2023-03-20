package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	tt "text/template"
	"time"

	_ "image/jpeg"
	_ "image/png"

	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
	"github.com/mitranim/gt"
)

const (
	SERVER_PORT  = 52693
	PUBLIC_DIR   = `public`
	TEMPLATE_DIR = `templates`
	EMDASH       = `—`
	EMAIL        = `me@mitranim.com`
)

var (
	FLAGS   = Flags{PROD: os.Getenv(`PROD`) == `true`}
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
	Url  = gt.NullUrl
)

type Flags struct{ PROD bool }

func fpj(path ...string) string { return filepath.Join(path...) }

func timeNow() Time { return gt.NullTimeNow().UTC() }

func timeParse(src string) (val Time) {
	gg.Try(val.Parse(src))
	return
}

func timeFmtHuman(date Time) string { return date.Format(`2006-Jan-02`) }

func trimLeadingSlash(val string) string   { return strings.TrimPrefix(val, `/`) }
func ensureLeadingSlash(val string) string { return ensurePrefix(val, `/`) }

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
	gg.Try(os.MkdirAll(filepath.Dir(path), os.ModePerm))
	gg.Try(os.WriteFile(path, bytes, os.ModePerm))
}

func yearsElapsed() string {
	start := 2014
	now := time.Now().UTC().Year()
	if now > start {
		return fmt.Sprintf(`%v—%v`, start, now)
	}
	return fmt.Sprint(start)
}

func imgConfig(path string) image.Config {
	file := gg.Try1(os.Open(path))
	defer file.Close()

	conf, _ := gg.Try2(image.DecodeConfig(file))
	return conf
}

// Inefficient but not measurable in our case.
func trimLines(val string) string {
	return strings.TrimSpace(strings.Join(reLines.Split(val, -1), ``))
}

func isDir(path string) bool {
	stat, _ := os.Stat(path)
	return stat != nil && stat.IsDir()
}

var reLines = regexp.MustCompile(`\s*(?:\r|\n)\s*`)

func tplToBytes(temp *tt.Template, val any) (buf x.NonEscWri) {
	gg.Try(temp.Execute(&buf, val))
	return buf
}

func xmlEncode(val any) (buf x.NonEscWri) {
	buf = append(buf, xml.Header...)
	enc := xml.NewEncoder(&buf)
	enc.Indent(``, `  `)
	gg.Try(enc.Encode(val))
	return buf
}

func timing(name string) func() {
	start := time.Now()
	log.Printf(`[%v] starting`, name)

	return func() {
		log.Printf(`[%v] done in %v`, name, time.Since(start))
	}
}

//nolint:unused,deadcode
func withTiming(str string, fun func()) {
	defer timing(str)()
	fun()
}

func ioWrite[Out io.Writer, Src gg.Text](out Out, src Src) {
	gg.Try1(out.Write(gg.ToBytes(src)))
}

func urlParse(src string) (out Url) {
	gg.Try(out.Parse(src))
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

	gg.Try(fs.WalkDir(os.DirFS(gg.Try1(os.Getwd())), path, func(path string, ent fs.DirEntry, err error) error {
		gg.Try(err)
		if ent.IsDir() {
			fun(path, ent)
		}
		return nil
	}))
}

func buiChild(bui B, val any) bool {
	size := len(*bui)
	bui.Child(val)
	return len(*bui) > size
}

type SemicolonList []Bui

func (self SemicolonList) Render(bui B) {
	has := false

	for _, val := range self {
		if len(val) == 0 {
			continue
		}

		if has {
			bui.T(`; `)
		}

		has = true
		val.Render(bui)
	}
}

func readTemplate(path string) []byte {
	return gg.ReadFile[[]byte](fpj(TEMPLATE_DIR, path))
}
