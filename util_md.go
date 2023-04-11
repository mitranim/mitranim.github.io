package main

import (
	"io"
	"regexp"
	"strings"
	tt "text/template"

	chtml "github.com/alecthomas/chroma/formatters/html"
	clexers "github.com/alecthomas/chroma/lexers"
	cstyles "github.com/alecthomas/chroma/styles"
	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
	"github.com/mitranim/gt"
	bf "github.com/russross/blackfriday/v2"
	san "github.com/shurcooL/sanitized_anchor_name"
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

	MD_INDENT               = `  `
	HEADING_TAGS            = [...]string{1: `h1`, 2: `h2`, 3: `h3`, 4: `h4`, 5: `h5`, 6: `h6`}
	HEADING_PREFIX          = F(E(`span`, AP(`class`, `heading-prefix`, `aria-hidden`, `true`)))
	RE_DETAIL_TAG_PREFIX    = regexp.MustCompile(`^details\b`)
	RE_DETAIL_TAG           = regexp.MustCompile(`^details(?:\s*[|]\s*(\w*)\s*[|]\s*(.*))?`)
	RE_PROTOCOL             = regexp.MustCompile(`^\w+://`)
	DETAIL_SUMMARY_FALLBACK = gg.ToBytes(`Click for details`)
)

func stringMdToHtml(src string, opt *MdOpt) string {
	return gg.ToString(mdToHtml(gg.ToBytes(src), opt))
}

func mdToHtml(src []byte, opt *MdOpt) []byte {
	return bf.Run(src, mdOpts(opt)...)
}

func mdTplToHtml(src []byte, opt *MdOpt, val any) []byte {
	return mdToHtml(mdTplExec(gg.ToString(src), val), opt)
}

func mdTplExec(src string, val any) []byte {
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

Note: somewhat duplicates `Exta`.
*/
func (self *MdRen) RenderLink(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	href := gg.ToString(node.LinkData.Destination)

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
	if RE_DETAIL_TAG_PREFIX.Match(tag) {
		return self.RenderCodeBlockDetails(out, node, entering)
	}
	return self.RenderCodeBlockHighlighted(out, node, entering)
}

/**
Special magic for code blocks like these:

```details | lang | summary
(some text)
```

This gets wrapped in a <details> element with the given <summary>. The lang tag
is optional; if present, the block is processed as code, otherwise as regular
text.
*/
func (self *MdRen) RenderCodeBlockDetails(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	tag := node.CodeBlockData.Info
	match := RE_DETAIL_TAG.FindSubmatch(tag)
	if match == nil {
		panic(gg.Errf(`invalid code block tag %q`, tag))
	}

	lang, summary := match[1], match[2]
	if len(summary) == 0 {
		summary = DETAIL_SUMMARY_FALLBACK
	}

	var bui Bui
	bui.E(
		`details`,
		AP(`class`, `details fan-typo`),
		E(`summary`, nil, Bui(mdToHtml(summary, nil))),
		func() {
			if len(lang) > 0 {
				// As code.
				node.CodeBlockData.Info = lang
				self.RenderCodeBlockHighlighted((*x.NonEscWri)(&bui), node, entering)
			} else {
				// As regular text.
				bui.NonEscBytes(mdToHtml(node.Literal, nil))
			}
		},
	)

	ioWrite(out, bui)
	return bf.SkipChildren
}

func (self *MdRen) RenderCodeBlockHighlighted(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	tag := node.CodeBlockData.Info
	lexer := clexers.Get(gg.ToString(tag))
	if lexer == nil {
		panic(gg.Errf(`no lexer for %q`, tag))
	}

	iterator, err := lexer.Tokenise(nil, gg.ToString(node.Literal))
	gg.Try(gg.Wrapf(err, `tokenizer error`))

	err = CHROMA_FORMATTER.Format(out, CHROMA_STYLE, iterator)
	gg.Try(gg.Wrapf(err, `formatter error`))

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
		panic(gg.Errf(`unrecognized heading level: %v`, headingLevel))
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
		ioWrite(out, `<blockquote class="blockquote">`)
	} else {
		ioWrite(out, `</blockquote>`)
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
	offset := gg.MinPrimBy(headings, MdHeading.GetLevel)

	for ind := range headings {
		tar := &headings[ind]
		tar.Level = gg.MaxPrim(0, tar.Level-offset)
	}

	var buf gg.Buf
	for _, val := range headings {
		buf = val.AppendTo(buf)
		buf.AppendNewline()
	}
	return buf.String()
}

func mdHeadings(content []byte) (out []MdHeading) {
	node := bf.New(mdOpts(nil)...).Parse(content)

	node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Document {
			return bf.GoToNext
		}

		if node.Type == bf.Heading {
			heading := MdHeading{
				Text:  node.Literal,
				Id:    node.HeadingData.HeadingID,
				Level: node.HeadingData.Level,
			}

			textNode := bfNodeFind(node, bf.Text)
			if textNode != nil && len(heading.Text) == 0 {
				heading.Text = textNode.Literal
			}
			if textNode != nil && heading.Id == `` {
				heading.Id = san.Create(gg.ToString(textNode.Literal))
			}
			if heading.IsValid() {
				out = append(out, heading)
			}
		}

		return bf.SkipChildren
	})

	return
}

type MdHeading struct {
	Text  []byte
	Id    string
	Level int
}

func (self MdHeading) IsValid() bool {
	return len(self.Text) > 0 && self.Id != ``
}

func (self MdHeading) GetLevel() int { return self.Level }

func (self MdHeading) AppendTo(buf gg.Buf) gg.Buf {
	buf = self.AppendIndentTo(buf)
	buf = self.AppendContentTo(buf)
	return buf
}

func (self MdHeading) AppendIndentTo(buf gg.Buf) gg.Buf {
	switch self.Level {
	case 0:
	case 1:
		buf.AppendString(MD_INDENT)
	default:
		/**
		We should be able to simply use `self.Level` here. However, there's a quirk
		either in Markdown itself or in the Markdown library we use that forces
		this workaround.
		*/
		buf.AppendStringN(MD_INDENT, 1+(self.Level*2))
	}
	return buf
}

func (self MdHeading) AppendContentTo(buf gg.Buf) gg.Buf {
	buf.AppendString(`* [`)
	buf.AppendBytes(self.Text)
	buf.AppendString(`](#`)
	buf.AppendString(self.Id)
	buf.AppendString(`)`)
	return buf
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
		id := `id` + gt.RandomUuid().String()
		funs[id] = func() string { return val }
		return `{{` + id + `}}`
	})

	gg.Try1(tpl.Funcs(funs).Parse(text))
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

func mdLink(text, link string) string {
	if link == `` {
		return text
	}
	if text == `` {
		return link
	}
	return `[` + text + `](` + link + `)`
}
