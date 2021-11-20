package main

import (
	"io"
	"math"
	"regexp"
	"strings"
	tt "text/template"

	chtml "github.com/alecthomas/chroma/formatters/html"
	clexers "github.com/alecthomas/chroma/lexers"
	cstyles "github.com/alecthomas/chroma/styles"
	x "github.com/mitranim/gax"
	"github.com/mitranim/gt"
	"github.com/mitranim/try"
	e "github.com/pkg/errors"
	bf "github.com/russross/blackfriday/v2"
	"github.com/shurcooL/sanitized_anchor_name"
)

type (
	MdOpt = bf.HTMLRendererParameters
)

var (
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

	HEADING_TAGS            = [...]string{1: "h1", 2: "h2", 3: "h3", 4: "h4", 5: "h5", 6: "h6"}
	HEADING_PREFIX          = F(E(`span`, AP(`class`, `heading-prefix`, `aria-hidden`, `true`)))
	RE_DETAIL_TAG           = regexp.MustCompile(`details"([^"\r\n]*)"(\S*)?`)
	RE_PROTOCOL             = regexp.MustCompile(`^\w+://`)
	DETAIL_SUMMARY_FALLBACK = stringToBytesAlloc(`Click for details`)
)

func stringMdToHtml(src string, opt *MdOpt) string {
	return bytesString(mdToHtml(stringToBytesAlloc(src), opt))
}

func mdToHtml(src []byte, opt *MdOpt) []byte {
	return bf.Run(src, mdOpts(opt)...)
}

func mdTplToHtml(src []byte, opt *MdOpt, val interface{}) []byte {
	return mdToHtml(mdTplExec(bytesString(src), val), opt)
}

func mdTplExec(src string, val interface{}) []byte {
	tpl := makeTpl(``)
	tplParseMd(tpl, src)
	return tplToBytes(tpl, val)
}

/*
Note: we create a new renderer for every page because `bf.HTMLRenderer` is
stateful and is not meant to be reused between unrelated texts. In particular,
reusing it between pages causes `bf.AutoHeadingIDs` to suffix heading IDs,
making them unique across multiple pages. We don't want that.
*/
func mdOpts(opt *MdOpt) []bf.Option {
	if opt == nil {
		opt = &MdOpt{}
	}
	if opt.Flags == 0 {
		opt.Flags = bf.CommonHTMLFlags
	}

	return []bf.Option{
		bf.WithExtensions(
			bf.Autolink | bf.Strikethrough | bf.FencedCode | bf.HeadingIDs | bf.AutoHeadingIDs | bf.Tables,
		),
		bf.WithRenderer(&MdRen{bf.NewHTMLRenderer(*opt)}),
	}
}

type MdRen struct{ *bf.HTMLRenderer }

func (self *MdRen) RenderNode(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Heading:
		return self.RenderHeading(out, node, entering)
	case bf.Link:
		return self.RenderLink(out, node, entering)
	case bf.CodeBlock:
		return self.RenderCodeBlock(out, node, entering)
	case bf.BlockQuote:
		return self.RenderBlockQuote(out, node, entering)
	default:
		return self.HTMLRenderer.RenderNode(out, node, entering)
	}
}

/**
Differences from default:

	* external links get attributes like `target="_blank"` and an external
	  link icon

	* intra-page hash links, like `href="#blah"`, are prefixed with a hash
	  symbol hidden from screen readers

"External href" is defined as "starts with a protocol".

Note: currently doesn't support some flags and extensions.

Note: somewhat duplicates `exta`.
*/
func (self *MdRen) RenderLink(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	href := bytesString(node.LinkData.Destination)

	if entering {
		attrs := AP(`href`, href)
		if isLinkExternal(href) {
			attrs = attrs.A(ABLAN...)
		}

		var bui Bui
		bui.Begin(`a`, attrs)
		if isLinkHash(href) {
			bui.E(`span`, AP(`class`, `hash-prefix noprint`, `aria-hidden`, `true`), `#`)
		}
		ioWrite(out, bui)
	} else {
		var bui Bui
		if isLinkExternal(href) {
			bui.F(SvgExternalLink)
		}
		bui.End(`a`)
		ioWrite(out, bui)
	}

	return bf.GoToNext
}

func isLinkExternal(val string) bool  { return hasLinkProtocol(val) }
func hasLinkProtocol(val string) bool { return RE_PROTOCOL.MatchString(val) }
func isLinkHash(val string) bool      { return strings.HasPrefix(val, `#`) }

/**
Differences from default:

	* code highlighting

	* supports special directives like rendering <details>
*/
func (self *MdRen) RenderCodeBlock(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	tag := node.CodeBlockData.Info

	if len(tag) == 0 {
		return self.HTMLRenderer.RenderNode(out, node, entering)
	}

	/**
	Special magic for code blocks like these:

	```details"title"lang
	(some text)
	```

	This gets wrapped in a <details> element, with the string in the middle
	as <summary>. The lang tag is optional; if present, the block is
	processed as code, otherwise as regular text.
	*/
	match := RE_DETAIL_TAG.FindSubmatch(tag)
	if match != nil {
		title := match[1]
		if len(title) == 0 {
			title = DETAIL_SUMMARY_FALLBACK
		}

		lang := match[2]
		var bui Bui

		bui.E(
			`details`,
			AP(`class`, `details fan-typo`),
			E(`summary`, nil, Bui(mdToHtml(title, nil))),
			func() {
				if len(lang) > 0 {
					// As code.
					node.CodeBlockData.Info = lang
					self.RenderNode((*x.NonEscWri)(&bui), node, entering)
				} else {
					// As regular text.
					bui.NonEscBytes(mdToHtml(node.Literal, nil))
				}
			},
		)

		ioWrite(out, bui)
		return bf.SkipChildren
	}

	lexer := clexers.Get(bytesString(tag))
	if lexer == nil {
		panic(e.Errorf(`no lexer for %q`, tag))
	}

	iterator, err := lexer.Tokenise(nil, bytesString(node.Literal))
	try.To(e.Wrap(err, "tokenizer error"))

	err = CHROMA_FORMATTER.Format(out, CHROMA_STYLE, iterator)
	try.To(e.Wrap(err, "formatter error"))

	return bf.SkipChildren
}

/**
Differences from default:

	* Fancy prefix indicating heading level, hidden from screen readers;
	  speaking it aloud is redundant because screen readers will indicate the
	  heading level anyway.

	* ID anchor suffix, hidden from screen readers; hearing it all the time
	  quickly gets tiring.
*/
func (self *MdRen) RenderHeading(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	headingLevel := self.HTMLRenderer.HTMLRendererParameters.HeadingLevelOffset + node.Level
	tag := HEADING_TAGS[headingLevel]
	if tag == `` {
		panic(e.Errorf("unrecognized heading level: %v", headingLevel))
	}

	if entering {
		var bui Bui
		bui.Begin(tag, A(aId(node.HeadingID)))
		bui.F(HEADING_PREFIX)
		ioWrite(out, bui)
	} else {
		var bui Bui
		if node.HeadingID != `` {
			bui.E(`a`, AP(
				`href`, `#`+node.HeadingID,
				`class`, `heading-anchor`,
				`aria-hidden`, `true`,
			))
		}
		bui.End(tag)
		ioWrite(out, bui)
	}

	return bf.GoToNext
}

func (self *MdRen) RenderBlockQuote(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	if entering {
		ioWriteString(out, `<blockquote class="blockquote">`)
	} else {
		ioWriteString(out, `</blockquote>`)
	}
	return bf.GoToNext
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
func mdToToc(content []byte) string {
	headings := mdHeadings(content)
	levelOffset := math.MaxInt
	for _, val := range headings {
		if val.Level < levelOffset {
			levelOffset = val.Level
		}
	}

	var buf gt.Raw
	for _, val := range headings {
		for i := val.Level - levelOffset; i > 0; i-- {
			buf = append(buf, `  `...)
		}
		buf = append(buf, `* [`...)
		buf = append(buf, val.Text...)
		buf = append(buf, `](#`...)
		buf = append(buf, val.Id...)
		buf = append(buf, `)`...)
		buf = append(buf, "\n"...)
	}
	return buf.String()
}

func mdHeadings(content []byte) []MdHeading {
	var out []MdHeading

	node := bf.New(mdOpts(nil)...).Parse(content)

	node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Document {
			return bf.GoToNext
		}

		if node.Type == bf.Heading {
			heading := MdHeading{
				Level: node.HeadingData.Level,
				Text:  node.Literal,
				Id:    node.HeadingData.HeadingID,
			}

			textNode := bfNodeFind(node, bf.Text)
			if textNode != nil && len(heading.Text) == 0 {
				heading.Text = textNode.Literal
			}
			if textNode != nil && heading.Id == `` {
				heading.Id = sanitized_anchor_name.Create(bytesString(textNode.Literal))
			}

			if len(heading.Text) > 0 && heading.Id != `` {
				out = append(out, heading)
			}
		}

		return bf.SkipChildren
	})

	return out
}

type MdHeading struct {
	Level int
	Text  []byte
	Id    string
}

func bfNodeFind(node *bf.Node, nodeType bf.NodeType) (out *bf.Node) {
	node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == nodeType {
			out = node
			return bf.SkipChildren
		}
		return bf.GoToNext
	})
	return
}

/*
Modifies the template to preserve content between ``` as-is. We need it raw for
Markdown and code highlighting.
*/
func tplParseMd(tpl *tt.Template, cont string) {
	funs := tt.FuncMap{}

	text := replaceCodeBlocks(cont, func(val string) string {
		id := `id` + randomHex()
		funs[id] = func() string { return val }
		return "{{" + id + "}}"
	})

	_, err := tpl.Funcs(funs).Parse(text)
	try.To(err)
}

/*
Known limitations:

	* Doesn't support single-backticks.
	* Doesn't support indented code blocks.
	* Doesn't support fences other than ```.
	* Doesn't support escapes.

Markdown technically allows fences from "``" to any amount of "`". Proper
matching requires backreferences, which are unsupported by Go regexps. In this
codebase, only "```" should be used.
*/
var replaceCodeBlocks = reFencedCodeBlock.ReplaceAllStringFunc

var reFencedCodeBlock = regexp.MustCompile(
	"```((?:[^`]|`[^`]|``[^`])*)```",
)
