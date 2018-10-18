package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"html/template"
	"io"
	"io/ioutil"
	l "log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/pkg/errors"
	bf "github.com/russross/blackfriday/v2"
)

const (
	OUT       = "public"
	FILE_MODE = 0600
	DIR_MODE  = 0700
)

var (
	log = l.New(os.Stderr, "", 0)

	TEMPS = template.New("")

	SITE = []interface{}{
		Page{
			Path:        "index.html",
			Title:       "about:mitranim",
			Description: "About me: details, works, posts",
		},

		Page{
			Path:  "404.html",
			Title: "Page Not Found",
		},

		Page{
			Path:        "works.html",
			Title:       "Works",
			Description: "Software I'm involved in",
		},

		Page{
			Path:        "posts.html",
			Title:       "Blog Posts",
			Description: "Random notes and thoughts",
		},

		Feed{
			Path:        "feed.xml",
			Title:       "Nelo Mitranim's Blog",
			Description: "Occasional notes, mostly about programming",
		},

		Page{
			Path:        "demos.html",
			Title:       "Demos",
			Description: "Silly little demos",
		},

		Page{
			Path:        "resume.html",
			Title:       "Resume",
			Description: "Nelo Mitranim's Resume",
		},

		Post{
			Page: Page{
				Path:        "posts/cheating-for-performance-pjax.html",
				Title:       "Cheating for Performance with Pjax",
				Description: "Faster page transitions, for free",
			},
			Md:     "cheating-for-performance-pjax.md",
			Date:   time.Date(2015, 7, 25, 0, 0, 0, 0, time.UTC),
			Listed: true,
		},

		Post{
			Page: Page{
				Path:        "posts/cheating-for-website-performance.html",
				Title:       "Cheating for Website Performance",
				Description: "Frontend tips for speeding up websites",
			},
			Md:     "cheating-for-website-performance.md",
			Date:   time.Date(2015, 3, 11, 0, 0, 0, 0, time.UTC),
			Listed: true,
		},

		Post{
			Page: Page{
				Path:  "posts/keeping-things-simple.html",
				Title: "Keeping Things Simple",
			},
			Md:     "keeping-things-simple.md",
			Date:   time.Date(2015, 3, 10, 0, 0, 0, 0, time.UTC),
			Listed: true,
		},

		Post{
			Page: Page{
				Path:        "posts/next-generation-today.html",
				Title:       "Next Generation Today",
				Description: "EcmaScript 2015/2016 workflow with current web frameworks",
			},
			Md:     "next-generation-today.md",
			Date:   time.Date(2015, 5, 18, 0, 0, 0, 0, time.UTC),
			Listed: false,
		},

		Post{
			Page: Page{
				Path:        "posts/old-posts.html",
				Title:       "Old Posts",
				Description: "some old stuff from around the net",
			},
			Md:     "old-posts.md",
			Date:   time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			Listed: true,
		},
	}

	FUNCS = map[string]interface{}{
		"external":       external,
		"current":        current,
		"now":            func() string { return isoDate(time.Now().UTC()) },
		"isoDate":        isoDate,
		"years":          years,
		"html":           func(val string) template.HTML { return template.HTML(val) },
		"attr":           func(val string) template.HTMLAttr { return template.HTMLAttr(val) },
		"md":             md,
		"getListedPosts": getListedPosts,
		"render":         renderByName,
		"join":           path.Join,
		"svg":            svg,
		"withHash":       withHash,
		"ngTemplate":     func() string { return NG_TEMPLATE },
	}

	WITH_HASHES = map[string]string{}

	MD_OPTS = []bf.Option{
		bf.WithExtensions(
			bf.Autolink | bf.Strikethrough | bf.FencedCode | bf.HeadingIDs,
		),
		bf.WithRenderer(&MdRenderer{*bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		})}),
	}

	CHROMA_FORMATTER = html.New()

	// CHROMA_STYLE     = styles.Colorful
	CHROMA_STYLE = styles.Pygments
	// CHROMA_STYLE = styles.Tango
	// CHROMA_STYLE = styles.VisualStudio
	// CHROMA_STYLE = styles.Xcode

	stack = errors.WithStack
)

type Page struct {
	Path        string
	Title       string
	Description string
	Type        string
	Image       string
}

type Post struct {
	Page
	Md     string
	Date   time.Time
	Listed bool
}

func (self Post) Slug() string {
	return strings.TrimSuffix(filepath.Base(self.Path), filepath.Ext(self.Path))
}

type Feed struct {
	Path        string
	Title       string
	Description string
}

func main() {
	t0 := time.Now()
	err := build()
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
	t1 := time.Now()
	log.Printf("Built in %v\n", t1.Sub(t0))
}

func build() error {
	err := initTemplates()
	if err != nil {
		return err
	}

	err = renderSite()
	if err != nil {
		return err
	}

	return nil
}

func initTemplates() error {
	TEMPS.Funcs(FUNCS)

	for _, pattern := range []string{
		"templates/*.md",
		"templates/**/*.md",
		"templates/*.html",
	} {
		_, err := TEMPS.ParseGlob(pattern)
		if err != nil {
			return stack(err)
		}
	}

	return nil
}

func renderSite() error {
	for _, entry := range SITE {
		switch entry := entry.(type) {
		case Page:
			path := entry.Path
			temp, err := findTemplate(path)
			if err != nil {
				return err
			}
			err = renderTo(temp, path, entry)
			if err != nil {
				return err
			}

		case Post:
			path := entry.Path
			temp, err := findTemplate("post.html")
			if err != nil {
				return err
			}
			err = renderTo(temp, path, entry)
			if err != nil {
				return err
			}

			// Redirect for old post URL
			meta := fmt.Sprintf(
				`<meta http-equiv="refresh" content="0;URL='https://mitranim.com/posts/%v'" />`,
				entry.Slug(),
			)
			err = writeTo(filepath.Join("thoughts", entry.Slug()), []byte(meta))
			if err != nil {
				return err
			}

		case Feed:
			path := entry.Path
			temp, err := findTemplate(path)
			if err != nil {
				return err
			}
			err = renderTo(temp, path, entry)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findTemplate(name string) (*template.Template, error) {
	temp := TEMPS.Lookup(name)
	if temp != nil {
		return temp, nil
	}

	var names []string
	for _, temp := range TEMPS.Templates() {
		if temp.Name() != "" {
			names = append(names, temp.Name())
		}
	}

	return nil, errors.Errorf("Template %q not found. Known templates: %v", name, names)
}

func render(temp *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := temp.Execute(&buf, data)
	if err != nil {
		return nil, stack(err)
	}
	return buf.Bytes(), nil
}

func renderByName(name string, data interface{}) (template.HTML, error) {
	temp, err := findTemplate(name)
	if err != nil {
		return "", err
	}

	bytes, err := render(temp, data)
	if err != nil {
		return "", err
	}

	return template.HTML(bytes), nil
}

func renderTo(temp *template.Template, path string, data interface{}) error {
	path = filepath.Join(OUT, path)

	err := os.MkdirAll(filepath.Dir(path), DIR_MODE)
	if err != nil {
		return stack(err)
	}

	out, err := os.Create(path)
	if err != nil {
		return stack(err)
	}
	defer out.Close()

	err = temp.Execute(out, data)
	if err != nil {
		return err
	}

	// log.Printf("Wrote %v\n", path)
	return nil
}

func writeTo(path string, bytes []byte) error {
	path = filepath.Join(OUT, path)

	err := os.MkdirAll(filepath.Dir(path), DIR_MODE)
	if err != nil {
		return stack(err)
	}

	err = ioutil.WriteFile(path, bytes, FILE_MODE)
	if err != nil {
		return stack(err)
	}

	// log.Printf("Wrote %v\n", path)
	return nil
}

func external() template.HTMLAttr {
	return template.HTMLAttr(`target="_blank" rel="noopener noreferrer"`)
}

func current(href string, data interface{}) template.HTMLAttr {
	var path string
	switch data := data.(type) {
	case Page:
		path = data.Path
	case Post:
		path = data.Path
	}
	if href == path {
		return "aria-current"
	}
	return ""
}

func isoDate(date time.Time) string {
	// time.Parse uses these magic numbers instead of conventional placeholders
	return date.Format("2006-01-02")
}

func years() string {
	start := 2014
	now := time.Now().UTC().Year()
	if now > start {
		return fmt.Sprintf("%v—%v", start, now)
	}
	return fmt.Sprint(start)
}

func md(val interface{}) template.HTML {
	var input []byte
	switch val := val.(type) {
	case []byte:
		input = val
	case string:
		input = []byte(val)
	case template.HTML:
		input = []byte(val)
	}
	return template.HTML(bf.Run(input, MD_OPTS...))
}

func getListedPosts() (out []Post) {
	for _, val := range SITE {
		switch val := val.(type) {
		case Post:
			if val.Listed {
				out = append(out, val)
			}
		}
	}
	return
}

func svg(name string) (template.HTML, error) {
	return renderByName("svg-"+name, nil)
}

func withHash(assetPath string) (string, error) {
	out := WITH_HASHES[assetPath]

	if out == "" {
		path := filepath.Join(OUT, assetPath)
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return "", stack(err)
		}

		hash := crc32.ChecksumIEEE(bytes)
		if hash == 0 {
			out = assetPath
		} else {
			out = fmt.Sprintf("%v?%v", assetPath, hash)
		}
		WITH_HASHES[assetPath] = out
	}

	return out, nil
}

var (
	detailTagReg  = regexp.MustCompile(`details"([^"\s]*)"(\S*)?`)
	DETAILS_START = []byte(`<details class="details fancy-typography">`)
	DETAILS_END   = []byte(`</details>`)
	SUMMARY_START = []byte(`<summary>`)
	SUMMARY_END   = []byte(`</summary>`)
)

type MdRenderer struct{ bf.HTMLRenderer }

func (self *MdRenderer) RenderNode(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	default:
		return self.HTMLRenderer.RenderNode(out, node, entering)

	case bf.CodeBlock:
		tag := string(node.CodeBlockData.Info)

		/*
			Special magic for code blocks like these:

			```details"title"lang
			(some text)
			```

			This gets wrapped in a <details> element, with the string in the middle
			as <summary>. The lang tag is optional; if present, the block is
			processed as code, otherwise as regular text.
		*/
		if detailTagReg.MatchString(tag) {
			match := detailTagReg.FindStringSubmatch(tag)
			title := match[1]
			lang := match[2]

			out.Write(DETAILS_START)
			out.Write(SUMMARY_START)
			out.Write([]byte(title))
			out.Write(SUMMARY_END)

			if lang != "" {
				// As code
				node.CodeBlockData.Info = []byte(lang)
				self.RenderNode(out, node, entering)
			} else {
				// As regular text
				out.Write(bf.Run(node.Literal, MD_OPTS...))
			}

			out.Write(DETAILS_END)
			return bf.SkipChildren
		}

		text := string(node.Literal)
		lexer := findLexer(tag, text)
		iterator, err := lexer.Tokenise(nil, string(text))
		if err != nil {
			log.Printf("Tokenizer error: %v", err)
			return self.HTMLRenderer.RenderNode(out, node, entering)
		}

		err = CHROMA_FORMATTER.Format(out, CHROMA_STYLE, iterator)
		if err != nil {
			log.Printf("Formatter error: %v", err)
			return self.HTMLRenderer.RenderNode(out, node, entering)
		}

		return bf.SkipChildren
	}
}

// TODO: instantiating some lexers is EXTREMELY SLOW (tens of milliseconds).
// This takes an order of magnitude more CPU time than the the rest of the
// build. The worst offender is JS. HTML also auto-detects and includes JS.
func findLexer(tag string, text string) (out chroma.Lexer) {
	if len(tag) > 0 {
		out = lexers.Get(tag)
	} else {
		out = lexers.Analyse(text)
	}
	if out == nil {
		out = lexers.Fallback
	}
	return out
}

// Must be interpolated raw
const NG_TEMPLATE = `
<div layout="gaps-1-v">
  <!-- Left column: source words -->
  <div flex="1" class="gaps-1-v">
    <h3 theme="text-primary" layout="row-between">
      <span>Source Words</span>
      <span id="indicator"></span>
    </h3>
    <form ng-submit="self.add()" layout="gaps-1-v"
          sf-tooltip="{{self.error}}" sf-trigger="{{!!self.error}}">
      <input flex="11" tabindex="1" ng-model="self.word">
      <button flex="1" theme="primary" tabindex="1">Add</button>
    </form>
    <div ng-repeat="word in self.words" layout="row-between gaps-1-v">
      <span flex="11" layout="cross-center" class="padding-1" style="margin-1-r">{{word}}</span>
      <button flex="1" ng-click="self.remove(word)">✕</button>
    </div>
  </div>

  <!-- Right column: generated results -->
  <div flex="1" class="gaps-1-v">
    <h3 theme="text-accent">Generated Words</h3>
    <form ng-submit="self.generate()" layout>
      <button flex="1" theme="accent" tabindex="1">Generate</button>
    </form>
    <div ng-repeat="word in self.results" layout="row-between">
      <button flex="1" ng-click="self.pick(word)">←</button>
      <span flex="11" layout="cross-center" class="padding-1" style="margin-1-l">{{word}}</span>
    </div>
    <div ng-if="self.depleted" layout="cross-center">
      <span theme="text-warn" class="padding-1">(depleted)</span>
    </div>
  </div>
</div>
`
