// +build mage

package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	ht "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/chroma"
	chtml "github.com/alecthomas/chroma/formatters/html"
	clexers "github.com/alecthomas/chroma/lexers"
	cstyles "github.com/alecthomas/chroma/styles"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
	bf "github.com/russross/blackfriday/v2"
)

const TEMPLATE_DIR = "templates"
const FILE_MODE = 0600
const DIR_MODE = 0700
const HUMAN_TIME_FORMAT = "Jan 02, 2006"

var SITE_FILE = filepath.Join(TEMPLATE_DIR, "site.toml")

// Rebuilds HTML.
func Html() error {
	mg.Deps(Styles)
	t0 := time.Now()
	err := buildSite()
	if err != nil {
		return err
	}
	t1 := time.Now()
	log.Printf("[html] built in %v", t1.Sub(t0))
	return nil
}

/*
Watches templates and rebuilds HTML.

If rebuilds become too slow from too many files, this could reinitialize and
re-render only the changed templates. Keeping it simple for now.
*/
func HtmlW(ctx context.Context) error {
	fsEvents := make(chan notify.EventInfo, 1)
	err := notify.Watch(TEMPLATE_DIR+"/...", fsEvents, FS_EVENTS)
	if err != nil {
		return err
	}
	defer notify.Stop(fsEvents)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case fsEvent := <-fsEvents:
			log.Println("[html] FS event:", fsEvent)

			err := Html()
			if err != nil {
				log.Println("[html] error:", err)
				continue
			}

			notifyClients()
		}
	}
}

var SITE struct {
	Pages []Page
	Posts []Post
}

var FEED_AUTHOR = &FeedAuthor{Name: "Nelo Mitranim", Email: "me@mitranim.com"}

var SITE_BASE = func() string {
	if FLAGS.DEV {
		return "http://localhost:" + SERVER_PORT
	}
	return "https://mitranim.com"
}()

var SITE_FEED = Feed{
	Title:   "Software, Tech, Philosophy, Games",
	XmlBase: SITE_BASE,
	AltLink: &FeedLink{
		Rel:  "alternate",
		Type: "text/html",
		Href: SITE_BASE + "/posts",
	},
	SelfLink: &FeedLink{
		Rel:  "self",
		Type: "application/atom+xml",
		Href: SITE_BASE + "/feed.xml",
	},
	Author:      FEED_AUTHOR,
	Created:     time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
	Updated:     time.Now(),
	Id:          SITE_BASE + "/posts",
	Description: `Random thoughts about technology`,
	Items:       nil,
}

var TEMPLATES *ht.Template

var TEMPLATE_FUNCS = ht.FuncMap{
	"asHtml":              asHtml,
	"asAttr":              asAttr,
	"toMarkdown":          toMarkdown,
	"externalAnchor":      externalAnchor,
	"current":             currentAttr,
	"now":                 func() string { return formatDateForHumans(time.Now().UTC()) },
	"formatDateForHumans": formatDateForHumans,
	"years":               years,
	"listedPosts":         listedPosts,
	"include":             includeTemplate,
	"includeWith":         includeTemplateWith,
	"joinPath":            path.Join,
	"linkWithHash":        linkWithHash,
	"raw":                 func(text string) ht.HTML { return ht.HTML(text) },
	"headingPrefix":       func() ht.HTML { return HEADING_PREFIX_HTML },
	"FLAGS":               func() Flags { return FLAGS },
}

var ASSET_HASHES = map[string]string{}

var CHROMA_FORMATTER = chtml.New()

/*
// Light
cstyles.Colorful
cstyles.Tango
cstyles.VisualStudio
cstyles.Xcode
cstyles.Pygments

// Dark
cstyles.Dracula
cstyles.Fruity
cstyles.Native
cstyles.Monokai
*/
var CHROMA_STYLE = cstyles.Monokai

/*
Note: we create a new renderer for every page because `bf.HTMLRenderer` is
stateful and is not meant to be reused between unrelated texts. In particular,
reusing it between pages causes `bf.AutoHeadingIDs` to suffix heading IDs,
making them unique across multiple pages. We don't want that.
*/
func markdownOpts() []bf.Option {
	return []bf.Option{
		bf.WithExtensions(
			bf.Autolink | bf.Strikethrough | bf.FencedCode | bf.HeadingIDs | bf.AutoHeadingIDs,
		),
		bf.WithRenderer(&MarkdownRenderer{bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		})}),
	}
}

type Page struct {
	Path        string
	Title       string
	Description string
	Type        string
	Image       string
	ForceLight  bool
}

type Post struct {
	Page
	InputPath    string
	ObsoletePath string
	HtmlBody     []byte
	Created      time.Time
	Updated      time.Time
	Public       flagBool
	Listed       flagBool
}

func (self Post) UrlFromSiteRoot() string {
	return "/" + strings.TrimSuffix(self.Path, filepath.Ext(self.Path))
}

type flagBool bool

func (self *flagBool) UnmarshalText(input []byte) error {
	switch string(input) {
	case "true":
		*self = true
		return nil

	case "false":
		*self = false
		return nil

	// Somehow arrives unquoted, just like true and false
	case "dev":
		*self = flagBool(FLAGS.DEV)
		return nil

	default:
		return errors.Errorf(`unrecognized flagBool value: %v; must be true, false, or dev`, input)
	}
}

func buildSite() error {
	err := initSite()
	if err != nil {
		return err
	}

	for _, page := range SITE.Pages {
		err := buildPage(page)
		if err != nil {
			return err
		}
	}

	feed := SITE_FEED
	for _, post := range SITE.Posts {
		if !post.Public {
			continue
		}
		feed, err = buildPost(post, feed)
		if err != nil {
			return err
		}
	}

	buf, err := xmlEncode(feed.AtomFeed())
	if err != nil {
		return err
	}
	err = writePublic("feed.xml", buf.Bytes())
	if err != nil {
		return err
	}

	buf, err = xmlEncode(feed.RssFeed())
	if err != nil {
		return err
	}
	err = writePublic("feed_rss.xml", buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func buildPage(page Page) error {
	temp, err := findTemplate(TEMPLATES, page.Path)
	if err != nil {
		return err
	}

	output, err := renderTemplate(temp, page)
	if err != nil {
		return err
	}

	err = writePublic(page.Path, output)
	if err != nil {
		return err
	}

	return nil
}

func buildPost(post Post, feed Feed) (Feed, error) {
	bodyTemp, err := findTemplate(TEMPLATES, "post-body.html")
	if err != nil {
		return feed, err
	}

	// Used for the page and the feed entry, enclosed in different layouts.
	post.HtmlBody, err = renderTemplate(bodyTemp, post)
	if err != nil {
		return feed, err
	}

	layoutTemp, err := findTemplate(TEMPLATES, "post-layout.html")
	if err != nil {
		return feed, err
	}

	content, err := renderTemplate(layoutTemp, post)
	if err != nil {
		return feed, err
	}

	feedPostLayoutTemp, err := findTemplate(TEMPLATES, "feed-post-layout.html")
	if err != nil {
		return feed, err
	}

	feedPostContent, err := renderTemplate(feedPostLayoutTemp, post)
	if err != nil {
		return feed, err
	}

	err = writePublic(post.Path, content)
	if err != nil {
		return feed, err
	}

	if post.ObsoletePath != "" {
		meta := fmt.Sprintf(`<meta http-equiv="refresh" content="0;URL='%v'" />`, post.UrlFromSiteRoot())
		err = writePublic(post.ObsoletePath, []byte(meta))
		if err != nil {
			return feed, err
		}
	}

	if !post.Listed {
		return feed, nil
	}

	href := SITE_BASE + post.UrlFromSiteRoot()
	feed.Items = append(feed.Items, FeedItem{
		XmlBase:     href,
		Title:       post.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: post.Page.Description,
		Id:          href,
		Created:     post.Created,                                          // TODO get from git?
		Updated:     anyTime(post.Created, post.Updated, time.Now().UTC()), // TODO get from git?
		Content:     string(feedPostContent),
	})

	return feed, nil
}

func initSite() error {
	_, err := toml.DecodeFile(SITE_FILE, &SITE)
	if err != nil {
		return err
	}

	temp := ht.New("")
	temp.Funcs(TEMPLATE_FUNCS)

	/**
	The following code is similar to `temp.ParseGlob()`, but:
		* accepts empty matches
		* rejects duplicates
		* preprocesses .md templates to preserve raw code blocks
	*/

	matches, err := globs(
		filepath.Join(TEMPLATE_DIR, "*.html"),
		filepath.Join(TEMPLATE_DIR, "*.md"),
		filepath.Join(TEMPLATE_DIR, "**/*.html"),
		filepath.Join(TEMPLATE_DIR, "**/*.md"),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, fsPath := range matches {
		virtPath := strings.TrimPrefix(filepath.ToSlash(fsPath), TEMPLATE_DIR+"/")

		if temp.Lookup(virtPath) != nil {
			return errors.Errorf("duplicate template %q", virtPath)
		}

		bytes, err := ioutil.ReadFile(fsPath)
		if err != nil {
			return errors.WithStack(err)
		}

		content := string(bytes)
		if filepath.Ext(fsPath) == ".md" {
			/**
			Modify the template to preserve content between ``` as-is. We
			need it raw for Markdown and code highlighting.
			*/
			content = codeBlockReg.ReplaceAllStringFunc(content, codeBlockToRaw)
		}

		_, err = temp.New(virtPath).Parse(content)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	TEMPLATES = temp
	return nil
}

var codeBlockReg = regexp.MustCompile("(?:^|\\n)```\\S*\\r?\\n((?:[^`]|`[^`]|``[^`])*)```")

func codeBlockToRaw(input string) string {
	return "{{raw (print `" + strings.Replace(input, "`", "` \"`\" `", -1) + "`)}}"
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

func xmlEncode(input interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")
	err := enc.Encode(input)
	return buf, errors.WithStack(err)
}

func findTemplate(root *ht.Template, name string) (*ht.Template, error) {
	temp := root.Lookup(name)
	if temp != nil {
		return temp, nil
	}

	var names []string
	for _, temp := range root.Templates() {
		if temp.Name() != "" {
			names = append(names, temp.Name())
		}
	}
	return nil, errors.Errorf("template %q not found; known templates: %v", name, names)
}

func renderTemplate(temp *ht.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := temp.Execute(&buf, data)
	return buf.Bytes(), errors.WithStack(err)
}

func includeTemplate(name string) (ht.HTML, error) {
	return includeTemplateWith(name, nil)
}

func includeTemplateWith(name string, data interface{}) (ht.HTML, error) {
	temp, err := findTemplate(TEMPLATES, name)
	if err != nil {
		return "", err
	}

	bytes, err := renderTemplate(temp, data)
	if err != nil {
		return "", err
	}

	return ht.HTML(bytes), nil
}

func writePublic(path string, bytes []byte) error {
	path = filepath.Join(PUBLIC_DIR, path)

	err := os.MkdirAll(filepath.Dir(path), DIR_MODE)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(path, bytes, FILE_MODE)
	if err != nil {
		return errors.WithStack(err)
	}

	// log.Printf("Wrote %v\n", path)
	return nil
}

var featherIconExternalLink = strings.TrimSpace(`
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display: inline-block; width: 1.5ex; height: 1.5ex; margin-left: 0.3ch;" aria-hidden="true">
	<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
	<polyline points="15 3 21 3 21 9" />
	<line x1="10" y1="14" x2="21" y2="3" />
</svg>
`)

var featherIconExternalLinkBytes = []byte(featherIconExternalLink)

// Note: somewhat duplicated in `MarkdownRenderer.RenderNode`.
func externalAnchor(href string, text string) ht.HTML {
	return ht.HTML(fmt.Sprintf(`<a href="%v" target="_blank" rel="noopener noreferrer" class="decorate-link">%v%v</a>`,
		href, text, featherIconExternalLink))
}

func currentAttr(href string, data interface{}) ht.HTMLAttr {
	var path string
	switch data := data.(type) {
	case Page:
		path = data.Path
	case Post:
		path = data.Path
	}
	if href == path {
		return `aria-current="page"`
	}
	return ""
}

func formatDateForHumans(date time.Time) string {
	return date.Format(HUMAN_TIME_FORMAT)
}

func years() string {
	start := 2014
	now := time.Now().UTC().Year()
	if now > start {
		return fmt.Sprintf("%v—%v", start, now)
	}
	return fmt.Sprint(start)
}

func asHtml(input interface{}) ht.HTML {
	return ht.HTML(toString(input))
}

func asAttr(input interface{}) ht.HTMLAttr {
	return ht.HTMLAttr(toString(input))
}

func toString(input interface{}) string {
	switch input := input.(type) {
	case []byte:
		return string(input)
	case string:
		return input
	case ht.HTML:
		return string(input)
	case ht.HTMLAttr:
		return string(input)
	default:
		panic(errors.Errorf("unrecognized input: %v", input))
	}
}

func toMarkdown(input interface{}) ht.HTML {
	return ht.HTML(bf.Run(toBytes(input), markdownOpts()...))
}

func toBytes(input interface{}) []byte {
	switch input := input.(type) {
	case []byte:
		return input
	case string:
		return []byte(input)
	case ht.HTML:
		return []byte(input)
	case ht.HTMLAttr:
		return []byte(input)
	default:
		panic(errors.Errorf("unrecognized input: %v", input))
	}
}

func listedPosts() (out []Post) {
	for _, post := range SITE.Posts {
		if post.Public && post.Listed {
			out = append(out, post)
		}
	}
	return
}

func linkWithHash(assetPath string) (string, error) {
	out := ASSET_HASHES[assetPath]

	if out == "" {
		path := filepath.Join(PUBLIC_DIR, assetPath)
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return "", errors.WithStack(err)
		}

		hash := crc32.ChecksumIEEE(bytes)
		if hash == 0 {
			out = assetPath
		} else {
			out = fmt.Sprintf("%v?%v", assetPath, hash)
		}
		ASSET_HASHES[assetPath] = out
	}

	return out, nil
}

var (
	detailTagReg = regexp.MustCompile(`details"([^"\s]*)"(\S*)?`)

	HEADING_TAGS = map[int][]byte{
		1: []byte("h1"),
		2: []byte("h2"),
		3: []byte("h3"),
		4: []byte("h4"),
		5: []byte("h5"),
		6: []byte("h6"),
	}

	DETAILS_START       = []byte(`<details class="details fancy-typography">`)
	DETAILS_END         = []byte(`</details>`)
	SUMMARY_START       = []byte(`<summary>`)
	SUMMARY_END         = []byte(`</summary>`)
	ANGLE_OPEN          = []byte("<")
	ANGLE_OPEN_SLASH    = []byte("</")
	ANGLE_CLOSE         = []byte(">")
	ANCHOR_TAG          = []byte("a")
	EXTERNAL_LINK_ATTRS = []byte(` target="_blank" rel="noopener noreferrer"`)
	HREF_START          = []byte(` href="`)
	HREF_END            = []byte(`"`)
	SPACE               = []byte(` `)
	HASH_PREFIX         = []byte(`<span class="hash-prefix noprint" aria-hidden="true">#</span>`)
	HEADING_PREFIX      = []byte(`<span class="heading-prefix" aria-hidden="true"></span>`)
	BLOCKQUOTE_START    = []byte(`<blockquote class="blockquote">`)
	BLOCKQUOTE_END      = []byte(`</blockquote>`)
)

var (
	HEADING_PREFIX_HTML = ht.HTML(HEADING_PREFIX)
)

var (
	externalLinkReg = regexp.MustCompile(`^\w+://`)
	hashLinkReg     = regexp.MustCompile(`^#`)
)

type MarkdownRenderer struct{ *bf.HTMLRenderer }

func (self *MarkdownRenderer) RenderNode(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	default:
		return self.HTMLRenderer.RenderNode(out, node, entering)

	/**
	Differences from default:

		* fancy prefix indicating heading level, hidden from screen readers

		* ID anchor suffix, hidden from screen readers
	*/
	case bf.Heading:
		headingLevel := self.HTMLRenderer.HTMLRendererParameters.HeadingLevelOffset + node.Level
		tag := HEADING_TAGS[headingLevel]
		if tag == nil {
			panic(errors.Errorf("unrecognized heading level: %v", headingLevel))
		}
		if entering {
			out.Write(ANGLE_OPEN)
			out.Write(tag)
			if node.HeadingID != "" {
				out.Write([]byte(` id="` + node.HeadingID + `"`))
			}
			out.Write(ANGLE_CLOSE)
			out.Write(HEADING_PREFIX)
		} else {
			if node.HeadingID != "" {
				out.Write([]byte(`<a href="#` + node.HeadingID + `" class="heading-anchor" aria-hidden="true"></a>`))
			}
			out.Write(ANGLE_OPEN_SLASH)
			out.Write(tag)
			out.Write(ANGLE_CLOSE)
		}
		return bf.GoToNext

	/**
	Differences from default:

		* external links get attributes like `target="_blank"` and an external
		  link icon

		* intra-page hash links, like `href="#blah"`, are prefixed with a hash
		  symbol hidden from screen readers

	"External href" is defined as "starts with a protocol".

	Note: currently doesn't support some flags and extensions.

	Note: somewhat duplicates `externalAnchor`.
	*/
	case bf.Link:
		if entering {
			out.Write(ANGLE_OPEN)
			out.Write(ANCHOR_TAG)
			out.Write(HREF_START)
			out.Write(node.LinkData.Destination)
			out.Write(HREF_END)
			if externalLinkReg.Match(node.LinkData.Destination) {
				out.Write(EXTERNAL_LINK_ATTRS)
			}
			out.Write(ANGLE_CLOSE)
			if hashLinkReg.Match(node.LinkData.Destination) {
				out.Write(HASH_PREFIX)
			}
		} else {
			if externalLinkReg.Match(node.LinkData.Destination) {
				out.Write(featherIconExternalLinkBytes)
			}
			out.Write(ANGLE_OPEN_SLASH)
			out.Write(ANCHOR_TAG)
			out.Write(ANGLE_CLOSE)
		}
		return bf.GoToNext

	/**
	Differences from default:

		* code highlighting

		* supports special directives like rendering <details>
	*/
	case bf.CodeBlock:
		tag := string(node.CodeBlockData.Info)

		/**
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
				out.Write(bf.Run(node.Literal, markdownOpts()...))
			}

			out.Write(DETAILS_END)
			return bf.SkipChildren
		}

		text := string(node.Literal)
		lexer := findLexer(tag, text)
		iterator, err := lexer.Tokenise(nil, text)
		if err != nil {
			panic(errors.Wrap(err, "tokenizer error"))
		}

		err = CHROMA_FORMATTER.Format(out, CHROMA_STYLE, iterator)
		if err != nil {
			panic(errors.Wrap(err, "formatter error"))
		}

		return bf.SkipChildren

	case bf.BlockQuote:
		if entering {
			out.Write(BLOCKQUOTE_START)
		} else {
			out.Write(BLOCKQUOTE_END)
		}
		return bf.GoToNext
	}
}

/*
TODO: instantiating some lexers is EXTREMELY SLOW (tens of milliseconds). This
takes at least as much time as the rest of the build. The worst offender is JS.
HTML also auto-detects and includes JS.
*/
func findLexer(tag string, text string) (out chroma.Lexer) {
	if len(tag) > 0 {
		out = clexers.Get(tag)
	} else {
		out = clexers.Analyse(text)
	}
	if out == nil {
		out = clexers.Fallback
	}
	return out
}

/*
This and other feed-related types are copied from
https://github.com/gorilla/feeds with minor modifications.
*/
type Feed struct {
	Title       string
	XmlBase     string
	AltLink     *FeedLink
	SelfLink    *FeedLink
	Description string
	Author      *FeedAuthor
	Created     time.Time
	Updated     time.Time
	Id          string
	Subtitle    string
	Items       []FeedItem
	Copyright   string
	Image       *FeedImage
}

type FeedLink struct {
	Rel    string
	Type   string
	Href   string
	Length string
}

type FeedAuthor struct {
	Name  string
	Email string
}

func (self FeedAuthor) RssAuthor() string {
	if len(self.Email) > 0 {
		str := self.Email
		if len(self.Name) > 0 {
			str += " (" + self.Name + ")"
		}
		return str
	}
	return self.Name
}

type FeedImage struct {
	Url    string
	Title  string
	Link   string
	Width  int
	Height int
}

type FeedItem struct {
	XmlBase     string
	Title       string
	Link        *FeedLink
	Source      *FeedLink
	Author      *FeedAuthor
	Description string // used as description in rss, summary in atom
	Id          string // used as guid in rss, id in atom
	Created     time.Time
	Updated     time.Time
	Enclosure   *FeedEnclosure
	Content     string
}

type FeedEnclosure struct {
	Url    string
	Length string
	Type   string
}

func (self Feed) AtomFeed() AtomFeed {
	feed := AtomFeed{
		Xmlns:    "https://www.w3.org/2005/Atom",
		XmlBase:  self.XmlBase,
		Title:    self.Title + " | Atom | mitranim",
		Subtitle: self.Description,
		Updated:  AtomTime(self.Updated),
		Rights:   self.Copyright,
	}

	if self.AltLink != nil {
		feed.Id = self.AltLink.Href
		feed.Links = append(feed.Links, AtomLink{
			Rel:  self.AltLink.Rel,
			Type: self.AltLink.Type,
			Href: self.AltLink.Href,
		})
	}

	if self.SelfLink != nil {
		feed.Links = append(feed.Links, AtomLink{
			Rel:  self.SelfLink.Rel,
			Type: self.SelfLink.Type,
			Href: self.SelfLink.Href,
		})
	}

	if self.Author != nil {
		feed.Author = &AtomAuthor{
			AtomPerson: AtomPerson{
				Name:  self.Author.Name,
				Email: self.Author.Email,
			},
		}
	}

	for _, item := range self.Items {
		var name string
		var email string
		if item.Author != nil {
			name = item.Author.Name
			email = item.Author.Email
		}

		entry := AtomEntry{
			XmlBase:   item.XmlBase,
			Title:     item.Title,
			Id:        item.Id,
			Published: AtomTime(item.Created),
			Updated:   AtomTime(anyTime(item.Updated, item.Created, time.Now().UTC())),
			Summary:   &AtomSummary{Type: "html", Content: item.Description},
		}

		var linkRel string
		if item.Link != nil {
			link := AtomLink{
				Href: item.Link.Href,
				Rel:  item.Link.Rel,
				Type: item.Link.Type,
			}
			if link.Rel == "" {
				link.Rel = "alternate"
			}
			linkRel = link.Rel
			entry.Links = append(entry.Links, link)
		}

		if item.Enclosure != nil && linkRel != "enclosure" {
			entry.Links = append(entry.Links, AtomLink{
				Href:   item.Enclosure.Url,
				Rel:    "enclosure",
				Type:   item.Enclosure.Type,
				Length: item.Enclosure.Length,
			})
		}

		// If there's content, assume it's html
		if len(item.Content) > 0 {
			entry.Content = &AtomContent{Type: "html", Content: item.Content}
		}

		if len(name) > 0 || len(email) > 0 {
			entry.Author = &AtomAuthor{AtomPerson: AtomPerson{Name: name, Email: email}}
		}

		feed.Entries = append(feed.Entries, entry)
	}

	return feed
}

func (self Feed) RssFeed() RssFeed {
	var author string
	if self.Author != nil {
		author = self.Author.RssAuthor()
	}

	var image *RssImage
	if self.Image != nil {
		image = &RssImage{
			Url:    self.Image.Url,
			Title:  self.Image.Title,
			Link:   self.Image.Link,
			Width:  self.Image.Width,
			Height: self.Image.Height,
		}
	}

	feed := RssFeed{
		Version:          "2.0",
		ContentNamespace: "http://purl.org/rss/1.0/modules/content/",
		XmlBase:          self.XmlBase,
		Channel: &RssChannel{
			Title:          self.Title + " | RSS | mitranim",
			Description:    self.Description,
			ManagingEditor: author,
			PubDate:        RssTime(anyTime(self.Created, time.Now().UTC())),
			LastBuildDate:  RssTime(anyTime(self.Updated, self.Created, time.Now().UTC())),
			Copyright:      self.Copyright,
			Image:          image,
		},
	}

	if self.AltLink != nil {
		feed.Channel.AltLink = self.AltLink.Href
	}

	for _, item := range self.Items {
		rssItem := RssItem{
			XmlBase:     item.XmlBase,
			Title:       item.Title,
			Description: item.Description,
			Guid:        item.Id,
			PubDate:     RssTime(item.Created),
		}

		if item.Link != nil {
			rssItem.Link = item.Link.Href
		}

		if len(item.Content) > 0 {
			rssItem.Content = &RssContent{Content: item.Content}
		}

		if item.Source != nil {
			rssItem.Source = item.Source.Href
		}

		if item.Enclosure != nil && item.Enclosure.Type != "" && item.Enclosure.Length != "" {
			rssItem.Enclosure = &RssEnclosure{
				Url:    item.Enclosure.Url,
				Type:   item.Enclosure.Type,
				Length: item.Enclosure.Length,
			}
		}

		if item.Author != nil {
			rssItem.Author = item.Author.RssAuthor()
		}

		feed.Channel.Items = append(feed.Channel.Items, rssItem)
	}

	return feed
}

type AtomFeed struct {
	XMLName     xml.Name    `xml:"feed"`
	Xmlns       string      `xml:"xmlns,attr"`
	XmlBase     string      `xml:"xml:base,attr,omitempty"`
	Title       string      `xml:"title"`   // required
	Id          string      `xml:"id"`      // required
	Updated     AtomTime    `xml:"updated"` // required
	Category    string      `xml:"category,omitempty"`
	Icon        string      `xml:"icon,omitempty"`
	Logo        string      `xml:"logo,omitempty"`
	Rights      string      `xml:"rights,omitempty"` // copyright used
	Subtitle    string      `xml:"subtitle,omitempty"`
	Author      *AtomAuthor `xml:"author,omitempty"`
	Contributor *AtomContributor
	Links       []AtomLink
	Entries     []AtomEntry
}

// Multiple links with different rel can coexist
// Atom 1.0 <link rel="enclosure" type="audio/mpeg" title="MP3" href="https://www.example.org/myaudiofile.mp3" length="1234" />
type AtomLink struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Href    string   `xml:"href,attr"`
	Length  string   `xml:"length,attr,omitempty"`
}

type AtomAuthor struct {
	XMLName xml.Name `xml:"author"`
	AtomPerson
}

type AtomContributor struct {
	XMLName xml.Name `xml:"contributor"`
	AtomPerson
}

type AtomPerson struct {
	Name  string `xml:"name,omitempty"`
	Uri   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

type AtomEntry struct {
	XMLName     xml.Name     `xml:"entry"`
	Xmlns       string       `xml:"xmlns,attr,omitempty"`
	XmlBase     string       `xml:"xml:base,attr,omitempty"`
	Title       string       `xml:"title"` // required
	Id          string       `xml:"id"`    // required
	Category    string       `xml:"category,omitempty"`
	Content     *AtomContent ``
	Rights      string       `xml:"rights,omitempty"`
	Source      string       `xml:"source,omitempty"`
	Published   AtomTime     `xml:"published,omitempty"`
	Updated     AtomTime     `xml:"updated"` // required
	Contributor *AtomContributor
	Links       []AtomLink   // required if no content
	Summary     *AtomSummary // required if content has src or is base64
	Author      *AtomAuthor  // required if feed lacks an author
}

type AtomContent struct {
	XMLName xml.Name `xml:"content"`
	Content string   `xml:",cdata"`
	Type    string   `xml:"type,attr"`
}

type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",cdata"`
	Type    string   `xml:"type,attr"`
}

type AtomTime time.Time

func (self AtomTime) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if time.Time(self).IsZero() {
		return nil
	}
	enc.EncodeToken(start)
	enc.EncodeToken(xml.CharData(time.Time(self).Format(time.RFC3339)))
	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

type RssFeed struct {
	XMLName          xml.Name `xml:"rss"`
	XmlBase          string   `xml:"xml:base,attr,omitempty"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:content,attr"`
	Channel          *RssChannel
}

type RssChannel struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	AltLink        string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        RssTime  `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  RssTime  `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Ttl            int      `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []RssItem
}

type RssImage struct {
	XMLName xml.Name `xml:"image"`
	Url     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

type RssItem struct {
	XMLName     xml.Name      `xml:"item"`
	XmlBase     string        `xml:"xml:base,attr,omitempty"`
	Title       string        `xml:"title"`       // required
	Link        string        `xml:"link"`        // required
	Description string        `xml:"description"` // required
	Content     *RssContent   ``
	Author      string        `xml:"author,omitempty"`
	Category    string        `xml:"category,omitempty"`
	Comments    string        `xml:"comments,omitempty"`
	Enclosure   *RssEnclosure ``
	Guid        string        `xml:"guid,omitempty"`    // Id used
	PubDate     RssTime       `xml:"pubDate,omitempty"` // created or updated
	Source      string        `xml:"source,omitempty"`
}

type RssContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

// RSS 2.0 <enclosure url="https://example.com/file.mp3" length="123456789" type="audio/mpeg" />
type RssEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type RssTime time.Time

func (self RssTime) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if time.Time(self).IsZero() {
		return nil
	}
	enc.EncodeToken(start)
	enc.EncodeToken(xml.CharData(time.Time(self).Format(time.RFC1123Z)))
	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func anyTime(vals ...time.Time) time.Time {
	for _, val := range vals {
		if !val.IsZero() {
			return val
		}
	}
	return time.Time{}
}
