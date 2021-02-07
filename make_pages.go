package main

import (
	"encoding/xml"
	"fmt"
	"hash/crc32"
	ht "html/template"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/chroma"
	chtml "github.com/alecthomas/chroma/formatters/html"
	clexers "github.com/alecthomas/chroma/lexers"
	cstyles "github.com/alecthomas/chroma/styles"
	g "github.com/mitranim/gtg"
	"github.com/pkg/errors"
	"github.com/rjeczalik/notify"
	bf "github.com/russross/blackfriday/v2"
	"github.com/shurcooL/sanitized_anchor_name"
)

// Rebuild pages (HTML and feed XML).
func Pages(task g.Task) error {
	g.MustWait(task, Styles)
	defer taskTiming(Pages)()
	return buildSite()
}

/*
Watch and rebuild pages.

If rebuilds become too slow because of too many files, this could reinitialize
and re-render only the changed pages rather than everything. Keeping it simple
for now.
*/
func PagesW(task g.Task) error {
	g.MustWait(task, g.Opt(Pages))

	return watch(fpj(TEMPLATE_DIR, "..."), notify.All, func(event notify.EventInfo) {
		onFsEvent(task, event)
		err := Pages(task)
		if err != nil {
			info.Println("[pages] error:", err)
			return
		}
		notifyClients(nil)
	})
}

var (
	SITE_FILE = fpj(TEMPLATE_DIR, "site.toml")

	// See `templates/site.toml`.
	SITE struct {
		Pages []Page
		Posts []Post
	}

	FEED_AUTHOR = &FeedAuthor{Name: "Nelo Mitranim", Email: "me@mitranim.com"}

	TEMPLATES *ht.Template
)

var TEMPLATE_FUNCS = ht.FuncMap{
	"FLAGS":               func() Flags { return FLAGS },
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
	"pathWithoutExt":      pathWithoutExt,
	"baseName":            baseName,
	"imgBoxWithLink":      imgBoxWithLink,
	"imgBox":              imgBox,
	"ariaHidden":          ariaHidden,
	"emoji":               emoji,
	"formatFloat":         formatFloat,
	"formatFloatPercent":  formatFloatPercent,
	"tableOfContents":     tableOfContents,
}

func siteBase() string {
	if FLAGS.PROD {
		return "https://mitranim.com"
	}
	return fmt.Sprintf("http://localhost:%v", SERVER_PORT)
}

func siteFeed() Feed {
	base := siteBase()

	return Feed{
		Title:   "Software, Tech, Philosophy, Games",
		XmlBase: base,
		AltLink: &FeedLink{
			Rel:  "alternate",
			Type: "text/html",
			Href: base + "/posts",
		},
		SelfLink: &FeedLink{
			Rel:  "self",
			Type: "application/atom+xml",
			Href: base + "/feed.xml",
		},
		Author:      FEED_AUTHOR,
		Published:   timePtr(time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)),
		Updated:     timePtr(time.Now()),
		Id:          base + "/posts",
		Description: `Random thoughts about technology`,
		Items:       nil,
	}
}

var (
	// Concurrency-unsafe like many other globals, but should only be called from
	// templating functions which are run sequentially.
	ASSET_HASHES = map[string]string{}

	CHROMA_FORMATTER = chtml.New()

	/**
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
	CHROMA_STYLE = cstyles.Monokai
)

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
	ColorScheme string
}

type Post struct {
	Page
	InputPath    string
	ObsoletePath string
	HtmlBody     []byte
	PublishedAt  *time.Time
	UpdatedAt    *time.Time
	IsListed     fudgedBool
}

func (self Post) ExistsAsFile() bool {
	return self.PublishedAt != nil || !FLAGS.PROD
}

func (self Post) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

func (self Post) UrlFromSiteRoot() string {
	return "/" + pathWithoutExt(self.Path)
}

// Somewhat inefficient but shouldn't be measurable.
func (self Post) TimeString() string {
	var out []string

	if self.PublishedAt != nil {
		out = append(out, `published `+formatDateForHumans(*self.PublishedAt))
		if self.UpdatedAt != nil {
			out = append(out, `updated `+formatDateForHumans(*self.UpdatedAt))
		}
	}

	return strings.Join(out, ", ")
}

type fudgedBool bool

func (self *fudgedBool) UnmarshalText(input []byte) error {
	switch string(input) {
	case "true":
		*self = true
		return nil

	case "false":
		*self = false
		return nil

	// Somehow arrives unquoted, just like "true" and "false". 🤨
	case "dev":
		*self = fudgedBool(!FLAGS.PROD)
		return nil

	default:
		return errors.Errorf(`unrecognized fudgedBool value: %q; must be "true", "false", or "dev"`, input)
	}
}

func buildSite() (err error) {
	defer rec(&err)

	must(initSite())
	for _, page := range SITE.Pages {
		must(buildPage(page))
	}

	feed := siteFeed()
	for i := range SITE.Posts {
		post := &SITE.Posts[i]
		must(maybeBuildPost(post))
		must(maybeAppendPost(&feed.Items, *post))
	}

	content, err := xmlEncode(feed.AtomFeed())
	must(err)

	must(writePublic("feed.xml", content))

	content, err = xmlEncode(feed.RssFeed())
	must(err)

	must(writePublic("feed_rss.xml", content))
	return nil
}

func buildPage(page Page) (err error) {
	defer rec(&err)

	tpl, err := findTemplate(TEMPLATES, page.Path)
	must(err)

	output, err := renderTemplate(tpl, page)
	must(err)

	return writePublic(page.Path, output)
}

func maybeBuildPost(post *Post) (err error) {
	defer rec(&err)

	if !post.ExistsAsFile() {
		return nil
	}

	bodyTemp, err := findTemplate(TEMPLATES, "post-body.html")
	must(err)

	// Used for the page and the feed entry, enclosed in different layouts.
	post.HtmlBody, err = renderTemplate(bodyTemp, post)
	must(err)

	layoutTemp, err := findTemplate(TEMPLATES, "post-layout.html")
	must(err)

	content, err := renderTemplate(layoutTemp, post)
	must(err)

	must(writePublic(post.Path, content))

	if post.ObsoletePath != "" {
		meta := fmt.Sprintf(`<meta http-equiv="refresh" content="0;URL='%v'" />`, post.UrlFromSiteRoot())
		must(writePublic(post.ObsoletePath, []byte(meta)))
	}

	return nil
}

func postToFeedItem(post Post) (_ FeedItem, err error) {
	defer rec(&err)

	feedPostLayoutTemp, err := findTemplate(TEMPLATES, "feed-post-layout.html")
	must(err)

	feedPostContent, err := renderTemplate(feedPostLayoutTemp, post)
	must(err)

	href := siteBase() + post.UrlFromSiteRoot()

	return FeedItem{
		XmlBase:     href,
		Title:       post.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: post.Page.Description,
		Id:          href,
		Published:   post.PublishedAt,                                                     // TODO get from git?
		Updated:     anyTime(post.PublishedAt, post.UpdatedAt, timePtr(time.Now().UTC())), // TODO get from git?
		Content:     string(feedPostContent),
	}, nil
}

func maybeAppendPost(feedItems *[]FeedItem, post Post) error {
	if !post.ExistsInFeeds() {
		return nil
	}

	item, err := postToFeedItem(post)
	if err != nil {
		return err
	}

	*feedItems = append(*feedItems, item)
	return nil
}

func initSite() (err error) {
	defer rec(&err)

	zero(&SITE)
	_, err = toml.DecodeFile(SITE_FILE, &SITE)
	must(errors.WithStack(err))

	tpl := ht.New("")
	tpl.Funcs(TEMPLATE_FUNCS)

	/**
	The following code is similar to `tpl.ParseGlob()`, but:
		* accepts empty matches
		* rejects duplicates
		* preprocesses .md templates to preserve raw code blocks
	*/

	matches, err := globs(
		fpj(TEMPLATE_DIR, "*.html"),
		fpj(TEMPLATE_DIR, "*.md"),
		fpj(TEMPLATE_DIR, "**/*.html"),
		fpj(TEMPLATE_DIR, "**/*.md"),
	)
	must(errors.WithStack(err))

	for _, fsPath := range matches {
		virtPath := strings.TrimPrefix(filepath.ToSlash(fsPath), TEMPLATE_DIR+"/")

		if tpl.Lookup(virtPath) != nil {
			return errors.Errorf("duplicate template %q", virtPath)
		}

		bytes, err := ioutil.ReadFile(fsPath)
		must(errors.WithStack(err))

		content := string(bytes)
		if filepath.Ext(fsPath) == ".md" {
			/**
			Modify the template to preserve content between ``` as-is. We
			need it raw for Markdown and code highlighting.

			TODO extend this to single grave quotes, and use a decent parser instead
			of a regexp. It might be viable to use Blackfriday to parse a Markdown
			file that includes Go templating directives.
			*/
			content = codeBlockReg.ReplaceAllStringFunc(content, codeBlockToRaw)
		}

		_, err = tpl.New(virtPath).Parse(content)
		must(errors.WithStack(err))
	}

	TEMPLATES = tpl
	return nil
}

var codeBlockReg = regexp.MustCompile("(?:^|\\n)```\\S*\\r?\\n((?:[^`]|`[^`]|``[^`])*)```")

/*
TODO: check if the standard library has a better quoting function. Might be able
to use `strconv.Quote`.
*/
func codeBlockToRaw(input string) string {
	return fmt.Sprintf(
		"{{raw (print `%v`)}}",
		strings.Replace(input, "`", "` \"`\" `", -1),
	)
}

func findTemplate(root *ht.Template, templateName string) (*ht.Template, error) {
	tpl := root.Lookup(templateName)
	if tpl != nil {
		return tpl, nil
	}

	var names []string
	for _, tpl := range root.Templates() {
		if tpl.Name() != "" {
			names = append(names, tpl.Name())
		}
	}
	return nil, errors.Errorf("template %q not found; known templates: %v", templateName, names)
}

func includeTemplate(templateName string) (ht.HTML, error) {
	return includeTemplateWith(templateName, nil)
}

func includeTemplateWith(templateName string, data interface{}) (ht.HTML, error) {
	tpl, err := findTemplate(TEMPLATES, templateName)
	if err != nil {
		return "", err
	}

	bytes, err := renderTemplate(tpl, data)
	if err != nil {
		return "", err
	}

	return ht.HTML(bytes), nil
}

/*
Scans the given Markdown template and generates a TOC from the headings.

Note: the Markdown library we're using has its own TOC feature, but it's
unusable for our purposes. Fortunately, it exposes the parser and AST, allowing
us to extract the heading data.

Note: we currently render markdown content as a Go template, which includes
parsing it for the TOC, then we render it as markdown, which involves parsing it
again. An ideal implementation would parse only once.

We could technically find the template by name and call `.Tree.Root.String()`
instead of reading from disk, but reading from disk is simpler and doesn't
depend on an obscure API.
*/
func tableOfContents(templateName string) (ht.HTML, error) {
	pt := fpj(TEMPLATE_DIR, templateName)
	content, err := ioutil.ReadFile(pt)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return ht.HTML(tableOfContentsFromMarkdown(content)), nil
}

func tableOfContentsFromMarkdown(content []byte) []byte {
	var out []byte

	headings := markdownHeadings(content)

	levelOffset := int(^uint(0) >> 1) // max `int`
	for _, heading := range headings {
		if heading.Level < levelOffset {
			levelOffset = heading.Level
		}
	}

	const indent = "  "
	for _, heading := range headings {
		for i := heading.Level - levelOffset; i > 0; i-- {
			out = append(out, indent...)
		}
		out = append(out, "* ["...)
		out = append(out, heading.Text...)
		out = append(out, "](#"...)
		out = append(out, heading.Id...)
		out = append(out, ")\n"...)
	}

	return out
}

type MarkdownHeading struct {
	Level int
	Text  []byte
	Id    string
}

func markdownHeadings(content []byte) []MarkdownHeading {
	var out []MarkdownHeading

	node := bf.New(markdownOpts()...).Parse(content)

	node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Document {
			return bf.GoToNext
		}

		if node.Type == bf.Heading {
			heading := MarkdownHeading{
				Level: node.HeadingData.Level,
				Text:  node.Literal,
				Id:    node.HeadingData.HeadingID,
			}

			textNode := bfNodeFind(node, bf.Text)
			if textNode != nil && len(heading.Text) == 0 {
				heading.Text = textNode.Literal
			}
			if textNode != nil && heading.Id == "" {
				heading.Id = sanitized_anchor_name.Create(string(textNode.Literal))
			}

			if len(heading.Text) > 0 && heading.Id != "" {
				out = append(out, heading)
			}
		}

		return bf.SkipChildren
	})

	return out
}

func bfNodeFind(node *bf.Node, nodeType bf.NodeType) *bf.Node {
	var out *bf.Node

	node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == nodeType {
			out = node
			return bf.SkipChildren
		}
		return bf.GoToNext
	})

	return out
}

func writePublic(path string, bytes []byte) error {
	path = fpj(PUBLIC_DIR, path)

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(path, bytes, FS_MODE_FILE)
	return errors.WithStack(err)
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
		if post.ExistsInFeeds() {
			out = append(out, post)
		}
	}
	return
}

func linkWithHash(assetPath string) (string, error) {
	out := ASSET_HASHES[assetPath]

	if out == "" {
		path := fpj(PUBLIC_DIR, assetPath)
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

		* Fancy prefix indicating heading level, hidden from screen readers;
		  speaking it aloud is redundant because screen readers will indicate the
		  heading level anyway.

		* ID anchor suffix, hidden from screen readers; hearing it all the time
		  quickly gets tiring.
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
	Published   *time.Time
	Updated     *time.Time
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
	Width  int64
	Height int64
}

type FeedItem struct {
	XmlBase     string
	Title       string
	Link        *FeedLink
	Source      *FeedLink
	Author      *FeedAuthor
	Description string // used as description in rss, summary in atom
	Id          string // used as guid in rss, id in atom
	Published   *time.Time
	Updated     *time.Time
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
		Updated:  (*AtomTime)(self.Updated),
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
			Published: (*AtomTime)(item.Published),
			Updated:   (*AtomTime)(anyTime(item.Updated, item.Published, timePtr(time.Now().UTC()))),
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
			PubDate:        (*RssTime)(anyTime(self.Published, timePtr(time.Now().UTC()))),
			LastBuildDate:  (*RssTime)(anyTime(self.Updated, self.Published, timePtr(time.Now().UTC()))),
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
			PubDate:     (*RssTime)(item.Published),
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
	Updated     *AtomTime   `xml:"updated"` // required
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
	Published   *AtomTime    `xml:"published,omitempty"`
	Updated     *AtomTime    `xml:"updated"` // required
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
	PubDate        *RssTime `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  *RssTime `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Ttl            int64    `xml:"ttl,omitempty"`
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
	Width   int64    `xml:"width,omitempty"`
	Height  int64    `xml:"height,omitempty"`
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
	PubDate     *RssTime      `xml:"pubDate,omitempty"` // created or updated
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
	enc.EncodeToken(start)
	enc.EncodeToken(xml.CharData(time.Time(self).Format(time.RFC1123Z)))
	enc.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func anyTime(vals ...*time.Time) *time.Time {
	for _, val := range vals {
		if val != nil && !val.IsZero() {
			return val
		}
	}
	return nil
}

func stripLeadingSlash(str string) string {
	return strings.TrimPrefix(str, "/")
}

func pathWithoutExt(pt string) string {
	return strings.TrimSuffix(pt, filepath.Ext(pt))
}

func baseName(pt string) string {
	return pathWithoutExt(filepath.Base(pt))
}

func imgConfig(pt string) (image.Config, error) {
	file, err := os.Open(pt)
	if err != nil {
		return image.Config{}, errors.WithStack(err)
	}
	defer file.Close()

	conf, _, err := image.DecodeConfig(file)
	return conf, errors.WithStack(err)
}

/*
Renders an image box. Scans the image file on disk to determine its dimentions.
Includes the height/width proportion into the template, which allows to ensure
fixed image dimensions and therefore prevent layout reflow on image load.
*/
func imgBoxWithLink(src string, caption string, href string) (ht.HTML, error) {
	// Takes tens of microseconds on my system, good enough for now.
	conf, err := imgConfig(stripLeadingSlash(src))
	if err != nil {
		return "", err
	}

	input := struct {
		Src     string
		Href    string // TODO: does this need to be `ht.HTMLAttr`?
		Caption string
		Width   int
		Height  int
	}{
		Src:     src,
		Href:    href,
		Caption: caption,
		Width:   conf.Width,
		Height:  conf.Height,
	}

	return includeTemplateWith("img-box.html", input)
}

func imgBox(src string, caption string) (ht.HTML, error) {
	return imgBoxWithLink(src, caption, "")
}

func ariaHidden(str string) ht.HTML {
	return ht.HTML(fmt.Sprintf(`<span aria-hidden="true">%v</span>`, str))
}

/*
Causes the MacOS VoiceOver to read the label followed by the word "image".
`role="img"` prevents it from reading the original name of the emoji, but forces
it to say "image". TODO: find a way to make it read only the label.
*/
func emoji(emoji string, label string) ht.HTML {
	if emoji == "" {
		return ""
	}
	if label == "" {
		return ariaHidden(emoji)
	}
	return ht.HTML(fmt.Sprintf(`<span aria-label="%v" role="img">%v</span>`, label, emoji))
}

func formatFloat(value float64, prec int) string {
	return strconv.FormatFloat(value, 'f', prec, 64)
}

func formatFloatPercent(value float64, prec int) string {
	return strconv.FormatFloat(value*100, 'f', prec, 64) + "%"
}

func timePtr(inst time.Time) *time.Time { return &inst }

func zero(val interface{}) {
	rval := reflect.ValueOf(val).Elem()
	rval.Set(reflect.New(rval.Type()).Elem())
}
