package main

/*
Templating is used for preprocessing Markdown: inserting SVG and some special
tidbits of markup, such as TOC, external links, and more.

This uses `text/template`, rather than `html/template`, because it's free of
weird surprises. We do our own escaping via Gax.
*/

import (
	"fmt"
	tt "text/template"

	"github.com/pkg/errors"
)

var TPL_FUNS = tt.FuncMap{
	"include":    include,
	"imgBox":     imgBox,
	"imgBoxLink": imgBoxLink,
	"emoji":      emoji,
	"exta":       exta,
	"mdToToc":    mdToToc,
	"today":      today,
}

func makeTpl(name string) *tt.Template {
	return tt.New(name).Funcs(TPL_FUNS)
}

func include(key string) string {
	out, ok := partials[key]
	if !ok {
		panic(errors.Errorf(`unknown include %q`, key))
	}
	return out
}

func exta(href string, text string) fmt.Stringer {
	return Exta(href, text)
}

func imgBox(src string, caption string) fmt.Stringer {
	return imgBoxLink(src, caption, "")
}

/*
Renders an image box. Scans the image file on disk to determine its dimentions.
Includes the height/width proportion into the template, which allows to ensure
fixed image dimensions and therefore prevent layout reflow on image load.
*/
func imgBoxLink(src string, caption string, href string) fmt.Stringer {
	// Takes tens of microseconds on my system, good enough for now.
	conf := imgConfig(trimLeadingSlash(src))

	return ImgBox(ImgMeta{
		Src:     src,
		Href:    href,
		Caption: caption,
		Width:   conf.Width,
		Height:  conf.Height,
	})
}

func emoji(emoji, label string) fmt.Stringer {
	if emoji == "" {
		return nil
	}

	if label == "" {
		return E(`span`, AP(`aria-hidden`, `true`), emoji)
	}

	return E(`span`, AP(`aria-label`, label, `role`, `img`), emoji)
}

func today() string {
	return timeFmtHuman(timeNow())
}
