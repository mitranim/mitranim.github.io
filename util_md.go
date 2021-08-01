package main

import (
	"io"
	"regexp"
	"strings"
	tt "text/template"

	chtml "github.com/alecthomas/chroma/formatters/html"
	clexers "github.com/alecthomas/chroma/lexers"
	cstyles "github.com/alecthomas/chroma/styles"
	"github.com/mitranim/try"
	"github.com/pkg/errors"
	bf "github.com/russross/blackfriday/v2"
	"github.com/shurcooL/sanitized_anchor_name"
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

	HEADING_TAGS = [...]string{
		1: "h1",
		2: "h2",
		3: "h3",
		4: "h4",
		5: "h5",
		6: "h6",
	}

	DETAILS_START       = `<details class="details fan-typo">`
	DETAILS_END         = `</details>`
	SUMMARY_START       = `<summary>`
	SUMMARY_END         = `</summary>`
	ANGLE_OPEN          = "<"
	ANGLE_OPEN_SLASH    = "</"
	ANGLE_CLOSE         = ">"
	ANCHOR_TAG          = "a"
	EXTERNAL_LINK_ATTRS = ` target="_blank" rel="noopener noreferrer"`
	HREF_START          = ` href="`
	HREF_END            = `"`
	HASH_PREFIX         = `<span class="hash-prefix noprint" aria-hidden="true">#</span>`
	HEADING_PREFIX      = `<span class="heading-prefix" aria-hidden="true"></span>`
	BLOCKQUOTE_START    = `<blockquote class="blockquote">`
	BLOCKQUOTE_END      = `</blockquote>`

	DETAIL_TAG_REG    = regexp.MustCompile(`details"([^"\s]*)"(\S*)?`)
	EXTERNAL_LINK_REG = regexp.MustCompile(`^\w+://`)
	HASH_LINK_REG     = regexp.MustCompile(`^#`)
)

func stringMdToHtml(src string) string {
	return bytesToMutableString(mdToHtml([]byte(src)))
}

func mdToHtml(src []byte) []byte {
	return bf.Run(src, mdOpts()...)
}

func mdTplToHtml(src []byte, val interface{}) []byte {
	return mdToHtml(mdTplToMd(bytesToMutableString(src), val))
}

func mdTplToMd(src string, val interface{}) []byte {
	tpl := makeTpl("")
	tplParseMd(tpl, src)
	return tplToBytes(tpl, val)
}

/*
Note: we create a new renderer for every page because `bf.HTMLRenderer` is
stateful and is not meant to be reused between unrelated texts. In particular,
reusing it between pages causes `bf.AutoHeadingIDs` to suffix heading IDs,
making them unique across multiple pages. We don't want that.
*/
func mdOpts() []bf.Option {
	return []bf.Option{
		bf.WithExtensions(
			bf.Autolink | bf.Strikethrough | bf.FencedCode | bf.HeadingIDs | bf.AutoHeadingIDs | bf.Tables,
		),
		bf.WithRenderer(&MdRen{bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		})}),
	}
}

type MdRen struct{ *bf.HTMLRenderer }

func (self *MdRen) RenderNode(out io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
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
		if tag == "" {
			panic(errors.Errorf("unrecognized heading level: %v", headingLevel))
		}
		if entering {
			ioWriteString(out, ANGLE_OPEN)
			ioWriteString(out, tag)
			if node.HeadingID != "" {
				ioWriteString(out, ` id="`+node.HeadingID+`"`)
			}
			ioWriteString(out, ANGLE_CLOSE)
			ioWriteString(out, HEADING_PREFIX)
		} else {
			if node.HeadingID != "" {
				ioWriteString(out, `<a href="#`+node.HeadingID+`" class="heading-anchor" aria-hidden="true"></a>`)
			}
			ioWriteString(out, ANGLE_OPEN_SLASH)
			ioWriteString(out, tag)
			ioWriteString(out, ANGLE_CLOSE)
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

	Note: somewhat duplicates `exta`.
	*/
	case bf.Link:
		if entering {
			ioWriteString(out, ANGLE_OPEN)
			ioWriteString(out, ANCHOR_TAG)
			ioWriteString(out, HREF_START)
			ioWrite(out, node.LinkData.Destination)
			ioWriteString(out, HREF_END)
			if EXTERNAL_LINK_REG.Match(node.LinkData.Destination) {
				ioWriteString(out, EXTERNAL_LINK_ATTRS)
			}
			ioWriteString(out, ANGLE_CLOSE)
			if HASH_LINK_REG.Match(node.LinkData.Destination) {
				ioWriteString(out, HASH_PREFIX)
			}
		} else {
			if EXTERNAL_LINK_REG.Match(node.LinkData.Destination) {
				ioWriteString(out, string(SvgExternalLink))
			}
			ioWriteString(out, ANGLE_OPEN_SLASH)
			ioWriteString(out, ANCHOR_TAG)
			ioWriteString(out, ANGLE_CLOSE)
		}
		return bf.GoToNext

	/**
	Differences from default:

		* code highlighting

		* supports special directives like rendering <details>
	*/
	case bf.CodeBlock:
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
		if DETAIL_TAG_REG.Match(tag) {
			match := DETAIL_TAG_REG.FindSubmatch(tag)
			title := match[1]
			lang := match[2]

			ioWriteString(out, DETAILS_START)
			ioWriteString(out, SUMMARY_START)
			ioWrite(out, title)
			ioWriteString(out, SUMMARY_END)

			if len(lang) > 0 {
				// As code
				node.CodeBlockData.Info = lang
				self.RenderNode(out, node, entering)
			} else {
				// As regular text
				ioWrite(out, mdToHtml(node.Literal))
			}

			ioWriteString(out, DETAILS_END)
			return bf.SkipChildren
		}

		lexer := clexers.Get(bytesToMutableString(tag))
		if lexer == nil {
			panic(errors.Errorf(`no lexer for %q`, tag))
		}

		iterator, err := lexer.Tokenise(nil, bytesToMutableString(node.Literal))
		try.To(errors.Wrap(err, "tokenizer error"))

		err = CHROMA_FORMATTER.Format(out, CHROMA_STYLE, iterator)
		try.To(errors.Wrap(err, "formatter error"))

		return bf.SkipChildren

	case bf.BlockQuote:
		if entering {
			ioWriteString(out, BLOCKQUOTE_START)
		} else {
			ioWriteString(out, BLOCKQUOTE_END)
		}
		return bf.GoToNext
	}
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
	var buf strings.Builder
	headings := mdHeadings(content)

	levelOffset := MAX_INT
	for _, val := range headings {
		if val.Level < levelOffset {
			levelOffset = val.Level
		}
	}

	for _, val := range headings {
		for i := val.Level - levelOffset; i > 0; i-- {
			buf.WriteString("  ")
		}
		buf.WriteString("* [")
		buf.Write(val.Text)
		buf.WriteString("](#")
		buf.WriteString(val.Id)
		buf.WriteString(")\n")
	}

	return buf.String()
}

func mdHeadings(content []byte) []MdHeading {
	var out []MdHeading

	node := bf.New(mdOpts()...).Parse(content)

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
			if textNode != nil && heading.Id == "" {
				heading.Id = sanitized_anchor_name.Create(bytesToMutableString(textNode.Literal))
			}

			if len(heading.Text) > 0 && heading.Id != "" {
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
