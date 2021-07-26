package main

import (
	"io"
	"regexp"
	"strings"
	tt "text/template"

	"github.com/alecthomas/chroma"
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

	HEADING_TAGS = [...][]byte{
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

	DETAIL_TAG_REG    = regexp.MustCompile(`details"([^"\s]*)"(\S*)?`)
	EXTERNAL_LINK_REG = regexp.MustCompile(`^\w+://`)
	HASH_LINK_REG     = regexp.MustCompile(`^#`)
)

func mdToHtml(src []byte) []byte {
	return bf.Run(src, mdOpts()...)
}

func tryMd(src []byte, val interface{}) []byte {
	return mdToHtml(try.ByteSlice(mdTplToMd(string(src), val)))
}

func mdTplToMd(src string, val interface{}) (_ []byte, err error) {
	defer try.Rec(&err)
	tpl := makeTpl("")
	try.To(tplParseMd(tpl, src))
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
			bf.Autolink | bf.Strikethrough | bf.FencedCode | bf.HeadingIDs | bf.AutoHeadingIDs,
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

	Note: somewhat duplicates `exta`.
	*/
	case bf.Link:
		if entering {
			out.Write(ANGLE_OPEN)
			out.Write(ANCHOR_TAG)
			out.Write(HREF_START)
			out.Write(node.LinkData.Destination)
			out.Write(HREF_END)
			if EXTERNAL_LINK_REG.Match(node.LinkData.Destination) {
				out.Write(EXTERNAL_LINK_ATTRS)
			}
			out.Write(ANGLE_CLOSE)
			if HASH_LINK_REG.Match(node.LinkData.Destination) {
				out.Write(HASH_PREFIX)
			}
		} else {
			if EXTERNAL_LINK_REG.Match(node.LinkData.Destination) {
				out.Write(SvgExternalLink)
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
		if DETAIL_TAG_REG.MatchString(tag) {
			match := DETAIL_TAG_REG.FindStringSubmatch(tag)
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
				out.Write(bf.Run(node.Literal, mdOpts()...))
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
func tplParseMd(tpl *tt.Template, cont string) error {
	funs := tt.FuncMap{}

	text := replaceCodeBlocks(cont, func(val string) string {
		id := `id` + randomHex()
		funs[id] = func() string { return val }
		return "{{" + id + "}}"
	})

	_, err := tpl.Funcs(funs).Parse(text)
	return errors.WithStack(err)
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
