package main

import (
	"path"

	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
)

const (
	ID_MAIN = `main`
	ID_TOP  = `top`
)

func HtmlCommon[A Ipage](page A, chi ...any) x.Bui {
	return Html(page, Header(page), chi, Footer(page))
}

func Html[A Ipage](page A, chi ...any) x.Bui {
	return F(
		x.Str(x.Doctype),
		E(`html`, AP(`class`, page.GetGlobalClass()),
			E(`head`, nil, HtmlHead(page)),
			E(`body`, AP(`id`, ID_TOP),
				// SkipToContent,
				chi,
			),
		),
	)
}

func HtmlHead[A Ipage](page A) x.Bui {
	return F(
		E(`meta`, AP(`charset`, `utf-8`)),
		E(`meta`, AP(`http-equiv`, `X-UA-Compatible`, `content`, `IE=edge,chrome=1`)),
		E(`meta`, AP(`name`, `viewport`, `content`, `width=device-width, minimum-scale=1, maximum-scale=2, initial-scale=1, user-scalable=yes`)),
		E(`link`, AP(`rel`, `icon`, `href`, `data:;base64,=`)),
		E(`link`, AP(`rel`, `stylesheet`, `type`, `text/css`, `href`, `/styles/main.css`)),

		func(bui B) {
			if page.GetTitle() != `` {
				bui.E(`title`, nil, page.GetTitle())
			} else {
				bui.E(`title`, nil, `about:mitranim`)
			}

			bui.E(
				`link`,
				AP(`rel`, `alternate`, `type`, `application/atom+xml`, `title`, `Mitranim's Posts (Atom)`, `href`, `/feed.xml`),
			)
			bui.E(
				`link`,
				AP(`rel`, `alternate`, `type`, `application/rss+xml`, `title`, `Mitranim's Posts (RSS)`, `href`, `/feed_rss.xml`),
			)
			bui.E(`meta`, AP(`name`, `author`, `content`, `Nelo Mitranim`))
			if page.GetTitle() != `` {
				bui.E(`meta`, AP(`property`, `og:title`, `content`, page.GetTitle()))
			}
			if page.GetDescription() != `` {
				bui.E(`meta`, AP(`name`, `description`, `content`, page.GetDescription()))
			}
			if page.GetImage() != `` {
				bui.E(`meta`, AP(`property`, `og:image`, `content`, path.Join(`/images`, page.GetImage())))
			}
			if page.GetType() != `` {
				bui.E(`meta`, AP(`property`, `og:type`, `content`, page.GetType()))
			}
			bui.E(`meta`, AP(`property`, `og:site_name`, `content`, `about:mitranim`))

			if !FLAGS.PROD {
				bui.E(`script`, AP(`type`, `module`, `src`, `http://localhost:52693/afr/client.mjs`))
			}
		},
	)
}

func Header[A Ipage](page A) x.Elem {
	const link = `header-link --busy`

	return E(`header`, AP(`class`, `header`),
		E(`nav`, AP(`class`, `flex row-sta-str`),
			E(`a`, AP(`class`, link).A(Cur(page, `/`)...), `home`),
			E(`a`, AP(`class`, link).A(Cur(page, `/works`)...), `works`),
			E(`a`, AP(`class`, link).A(Cur(page, `/posts`)...), `posts`),
			E(`a`, AP(`class`, link).A(Cur(page, `/games`)...), `games`),
		),

		E(`span`, AP(`class`, `flex-1`)),

		E(`span`, AP(`class`, `fg-blue flex row-cen-cen pad-body sm-hide`),
			`Updated: `+today(),
		),
	)
}

func MainContent(chi ...any) x.Elem { return E(`div`, AttrsMain(), chi...) }

func AttrsMain() x.Attrs {
	return AP(`role`, `main`, `id`, ID_MAIN, `class`, `main`)
}

func AttrsMainArticleMd() x.Attrs { return AttrsMain().Add(`class`, `article`) }

func Footer[A Ipage](page A) x.Elem {
	return E(`footer`, AP(`class`, `footer pad-bot-body`),
		E(`span`, AP(`class`, `flex-1 flex row-sta-cen gap-hor-0x5 pad-lef-body`),
			E(`span`, AP(`class`, `text-lef`), yearsElapsed()),
			LinkExt(`https://github.com/mitranim/mitranim.github.io`, `code`).AttrAdd(`class`, `wspace-nowrap`),
		),

		E(`span`, AP(`class`, `flex-1 text-cen`), func(bui B) {
			if page.GetLink() != `/` {
				bui.E(`a`, AP(`href`, `/#contacts`, `class`, `link-deco`), `touch me`)
			}
		}),

		E(`span`, AP(`class`, `flex-1 flex row-end-cen`),
			E(
				`a`,
				AP(
					`href`, idToHash(ID_TOP),
					`class`, `fill-gray-fg-near pad-body theme-plain-bg-gray`,
					`onclick`, `event.preventDefault(); window.scrollTo(0, 0)`,
				),
				SvgArrowUp,
			),
		),
	)
}

func FeedLinks() x.Elem {
	return E(`p`, AP(`class`, `feed-links`),
		E(`span`, nil, `Subscribe using one of:`),
		LinkExt(`/feed.xml`, `Atom`),
		LinkExt(`/feed_rss.xml`, `RSS`),
		LinkExt(`https://feedly.com/i/subscription/feed/https://mitranim.com/feed.xml`, `Feedly`),
	)
}

func FeedPost(page PagePost) x.Elem {
	return E(`article`, AP(`role`, `main article`, `class`, `typography`),
		FeedPostDesc(page),
		page.MdOnce(page),
	)
}

/*
Unsure if we want this. The Atom or RSS `FeedItem` already contains the
desription as a separate field. The RSS readers I use tend to show it in feed
item previews, but not when viewing the article. However, other readers may
include the description into the article, making this redundant.
*/
func FeedPostDesc(page PagePost) any {
	if true || page.Description == `` {
		return nil
	}
	return E(`p`, AP(`role`, `doc-subtitle`, `class`, `size-large italic`),
		page.Description,
	)
}

/**
Must be revised. This should not be accidentally read by voiceover utils, and on
click, it must skip to the content without changing the URL or polluting the
browser history. Previously, it seemed to work with the MacOS VoiceOver, but
right now on MacOS BigSur, it's flaky.
*/
// nolint:deadcode
var SkipToContent = E(
	`a`,
	AP(
		`href`, idToHash(ID_MAIN),
		`class`, `skip-to-content`,
		`onclick`, `event.preventDefault(); if (document.getElementById('main')) {document.getElementById('main').scrollIntoView()}`,
	),
	`Skip to content`,
)

func LinkExt(href, text string) x.Elem {
	if href == `` {
		panic(gg.Errf(`unexpected empty link`))
	}

	return E(
		`a`,
		AP(`href`, href, `class`, `link-deco`).A(ABLAN...),
		gg.Or(text, href),
		// We would prefer to use CSS `::after` with SVG as `background-image`, but
		// it doesn't seem to be able to inherit `currentColor`. Inline SVG avoids
		// that issue.
		SvgExternalLink,
	)
}

type ImgMeta struct {
	Src     string
	Href    string
	Caption string
	Width   int
	Height  int
}

/**
Note: using <figure> and <figcaption> would cause the MacOS VoiceOver to read
the caption twice when reading the content sequentially.
*/
func ImgBox(meta ImgMeta) x.Elem {
	if meta.Caption == `` {
		meta.Caption = baseName(meta.Src)
	}

	linkAttrs := AP(`class`, `img-box-link`)

	inner := E(`img`, AP(
		`src`, meta.Src,
		`alt`, meta.Caption,
		`class`, `img-box-img`,
		`style`, gg.Str(`aspect-ratio: `, meta.Width, `/`, meta.Height),
	))

	return E(`div`, AP(`class`, `img-box`), func(bui B) {
		if meta.Href != `` {
			bui.E(`a`, linkAttrs.AP(`href`, meta.Href).A(ABLAN...), inner)
		} else {
			bui.E(`div`, linkAttrs, inner)
		}
		bui.E(`span`, AP(`class`, `img-box-caption`, `aria-hidden`, `true`), meta.Caption)
	})
}

var partials = map[string]string{
	`svg-book`:          string(SvgBook),
	`svg-external-link`: string(SvgExternalLink),
	`svg-skype`:         string(SvgSkype),
	`svg-mobile-alt`:    string(SvgMobileAlt),
	`svg-paper-plane`:   string(SvgPaperPlane),
	`svg-github`:        string(SvgGithub),
	`svg-youtube`:       string(SvgYoutube),
	`svg-twitter`:       string(SvgTwitter),
	`svg-linkedin`:      string(SvgLinkedin),
	`svg-facebook`:      string(SvgFacebook),
	`svg-discord`:       string(SvgDiscord),
	`svg-arrow-up`:      string(SvgArrowUp),
	`svg-rss`:           string(SvgRss),
	`svg-rss-square`:    string(SvgRssSquare),
	`svg-print`:         string(SvgPrint),
}

/*
Source: https://feathericons.com

TODO consider using an SVG sprite. In particular, we repeat `SvgExternalLink` so
many times on some pages, that a sprite may reduce total size.
*/
var SvgBook x.Str = `<svg class="svg-icon stroke-fg" width="1em" height="1em" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-book-open"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path></svg>`
var SvgExternalLink = x.Str(trimLines(`
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="display: inline-block; width: 1.5ex; height: 1.5ex; margin-left: 0.3ch;" aria-hidden="true">
	<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
	<polyline points="15 3 21 3 21 9" />
	<line x1="10" y1="14" x2="21" y2="3" />
</svg>
`))

/* Source: https://github.com/FortAwesome/Font-Awesome */
var SvgSkype x.Str = `<svg class="svg-icon fill-skype" viewBox="0 0 448 512" width="1em" height="1em"><path d="M424.7 299.8c2.9-14 4.7-28.9 4.7-43.8 0-113.5-91.9-205.3-205.3-205.3-14.9 0-29.7 1.7-43.8 4.7C161.3 40.7 137.7 32 112 32 50.2 32 0 82.2 0 144c0 25.7 8.7 49.3 23.3 68.2-2.9 14-4.7 28.9-4.7 43.8 0 113.5 91.9 205.3 205.3 205.3 14.9 0 29.7-1.7 43.8-4.7 19 14.6 42.6 23.3 68.2 23.3 61.8 0 112-50.2 112-112 .1-25.6-8.6-49.2-23.2-68.1zm-194.6 91.5c-65.6 0-120.5-29.2-120.5-65 0-16 9-30.6 29.5-30.6 31.2 0 34.1 44.9 88.1 44.9 25.7 0 42.3-11.4 42.3-26.3 0-18.7-16-21.6-42-28-62.5-15.4-117.8-22-117.8-87.2 0-59.2 58.6-81.1 109.1-81.1 55.1 0 110.8 21.9 110.8 55.4 0 16.9-11.4 31.8-30.3 31.8-28.3 0-29.2-33.5-75-33.5-25.7 0-42 7-42 22.5 0 19.8 20.8 21.8 69.1 33 41.4 9.3 90.7 26.8 90.7 77.6 0 59.1-57.1 86.5-112 86.5z"/></svg>`
var SvgMobileAlt x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 320 512" width="1em" height="1em"><path d="M272 0H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h224c26.5 0 48-21.5 48-48V48c0-26.5-21.5-48-48-48zM160 480c-17.7 0-32-14.3-32-32s14.3-32 32-32 32 14.3 32 32-14.3 32-32 32zm112-108c0 6.6-5.4 12-12 12H60c-6.6 0-12-5.4-12-12V60c0-6.6 5.4-12 12-12h200c6.6 0 12 5.4 12 12v312z"/></svg>`
var SvgPaperPlane x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 512 512" width="1em" height="1em"><path d="M440 6.5L24 246.4c-34.4 19.9-31.1 70.8 5.7 85.9L144 379.6V464c0 46.4 59.2 65.5 86.6 28.6l43.8-59.1 111.9 46.2c5.9 2.4 12.1 3.6 18.3 3.6 8.2 0 16.3-2.1 23.6-6.2 12.8-7.2 21.6-20 23.9-34.5l59.4-387.2c6.1-40.1-36.9-68.8-71.5-48.9zM192 464v-64.6l36.6 15.1L192 464zm212.6-28.7l-153.8-63.5L391 169.5c10.7-15.5-9.5-33.5-23.7-21.2L155.8 332.6 48 288 464 48l-59.4 387.3z"/></svg>`
var SvgGithub x.Str = `<svg class="svg-icon fill-github" viewBox="0 0 496 512" width="1em" height="1em"><path d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3.3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5.3-6.2 2.3zm44.2-1.7c-2.9.7-4.9 2.6-4.6 4.9.3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3.7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3.3 2.9 2.3 3.9 1.6 1 3.6.7 4.3-.7.7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3.7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3.7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"/></svg>`
var SvgYoutube x.Str = `<svg class="svg-icon fill-youtube" viewBox="0 0 576 512" width="1em" height="1em"><path d="M549.655 124.083c-6.281-23.65-24.787-42.276-48.284-48.597C458.781 64 288 64 288 64S117.22 64 74.629 75.486c-23.497 6.322-42.003 24.947-48.284 48.597-11.412 42.867-11.412 132.305-11.412 132.305s0 89.438 11.412 132.305c6.281 23.65 24.787 41.5 48.284 47.821C117.22 448 288 448 288 448s170.78 0 213.371-11.486c23.497-6.321 42.003-24.171 48.284-47.821 11.412-42.867 11.412-132.305 11.412-132.305s0-89.438-11.412-132.305zm-317.51 213.508V175.185l142.739 81.205-142.739 81.201z"/></svg>`
var SvgTwitter x.Str = `<svg class="svg-icon fill-twitter" viewBox="0 0 512 512" width="1em" height="1em"><path d="M459.37 151.716c.325 4.548.325 9.097.325 13.645 0 138.72-105.583 298.558-298.558 298.558-59.452 0-114.68-17.219-161.137-47.106 8.447.974 16.568 1.299 25.34 1.299 49.055 0 94.213-16.568 130.274-44.832-46.132-.975-84.792-31.188-98.112-72.772 6.498.974 12.995 1.624 19.818 1.624 9.421 0 18.843-1.3 27.614-3.573-48.081-9.747-84.143-51.98-84.143-102.985v-1.299c13.969 7.797 30.214 12.67 47.431 13.319-28.264-18.843-46.781-51.005-46.781-87.391 0-19.492 5.197-37.36 14.294-52.954 51.655 63.675 129.3 105.258 216.365 109.807-1.624-7.797-2.599-15.918-2.599-24.04 0-57.828 46.782-104.934 104.934-104.934 30.213 0 57.502 12.67 76.67 33.137 23.715-4.548 46.456-13.32 66.599-25.34-7.798 24.366-24.366 44.833-46.132 57.827 21.117-2.273 41.584-8.122 60.426-16.243-14.292 20.791-32.161 39.308-52.628 54.253z"/></svg>`
var SvgLinkedin x.Str = `<svg class="svg-icon fill-linkedin" viewBox="0 0 448 512" width="1em" height="1em"><path d="M416 32H31.9C14.3 32 0 46.5 0 64.3v383.4C0 465.5 14.3 480 31.9 480H416c17.6 0 32-14.5 32-32.3V64.3c0-17.8-14.4-32.3-32-32.3zM135.4 416H69V202.2h66.5V416zm-33.2-243c-21.3 0-38.5-17.3-38.5-38.5S80.9 96 102.2 96c21.2 0 38.5 17.3 38.5 38.5 0 21.3-17.2 38.5-38.5 38.5zm282.1 243h-66.4V312c0-24.8-.5-56.7-34.5-56.7-34.6 0-39.9 27-39.9 54.9V416h-66.4V202.2h63.7v29.2h.9c8.9-16.8 30.6-34.5 62.9-34.5 67.2 0 79.7 44.3 79.7 101.9V416z"/></svg>`
var SvgFacebook x.Str = `<svg class="svg-icon fill-facebook" viewBox="0 0 448 512" width="1em" height="1em"><path d="M448 56.7v398.5c0 13.7-11.1 24.7-24.7 24.7H309.1V306.5h58.2l8.7-67.6h-67v-43.2c0-19.6 5.4-32.9 33.5-32.9h35.8v-60.5c-6.2-.8-27.4-2.7-52.2-2.7-51.6 0-87 31.5-87 89.4v49.9h-58.4v67.6h58.4V480H24.7C11.1 480 0 468.9 0 455.3V56.7C0 43.1 11.1 32 24.7 32h398.5c13.7 0 24.8 11.1 24.8 24.7z"/></svg>`
var SvgDiscord x.Str = `<svg class="svg-icon fill-discord" viewBox="0 0 448 512" width="1em" height="1em"><path d="M297.216 243.2c0 15.616-11.52 28.416-26.112 28.416-14.336 0-26.112-12.8-26.112-28.416s11.52-28.416 26.112-28.416c14.592 0 26.112 12.8 26.112 28.416zm-119.552-28.416c-14.592 0-26.112 12.8-26.112 28.416s11.776 28.416 26.112 28.416c14.592 0 26.112-12.8 26.112-28.416.256-15.616-11.52-28.416-26.112-28.416zM448 52.736V512c-64.494-56.994-43.868-38.128-118.784-107.776l13.568 47.36H52.48C23.552 451.584 0 428.032 0 398.848V52.736C0 23.552 23.552 0 52.48 0h343.04C424.448 0 448 23.552 448 52.736zm-72.96 242.688c0-82.432-36.864-149.248-36.864-149.248-36.864-27.648-71.936-26.88-71.936-26.88l-3.584 4.096c43.52 13.312 63.744 32.512 63.744 32.512-60.811-33.329-132.244-33.335-191.232-7.424-9.472 4.352-15.104 7.424-15.104 7.424s21.248-20.224 67.328-33.536l-2.56-3.072s-35.072-.768-71.936 26.88c0 0-36.864 66.816-36.864 149.248 0 0 21.504 37.12 78.08 38.912 0 0 9.472-11.52 17.152-21.248-32.512-9.728-44.8-30.208-44.8-30.208 3.766 2.636 9.976 6.053 10.496 6.4 43.21 24.198 104.588 32.126 159.744 8.96 8.96-3.328 18.944-8.192 29.44-15.104 0 0-12.8 20.992-46.336 30.464 7.68 9.728 16.896 20.736 16.896 20.736 56.576-1.792 78.336-38.912 78.336-38.912z"/></svg>`
var SvgArrowUp x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 448 512" width="1em" height="1em"><path d="M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V456c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z"/></svg>`
var SvgRss x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 448 512" width="1em" height="1em"><path d="M128.081 415.959c0 35.369-28.672 64.041-64.041 64.041S0 451.328 0 415.959s28.672-64.041 64.041-64.041 64.04 28.673 64.04 64.041zm175.66 47.25c-8.354-154.6-132.185-278.587-286.95-286.95C7.656 175.765 0 183.105 0 192.253v48.069c0 8.415 6.49 15.472 14.887 16.018 111.832 7.284 201.473 96.702 208.772 208.772.547 8.397 7.604 14.887 16.018 14.887h48.069c9.149.001 16.489-7.655 15.995-16.79zm144.249.288C439.596 229.677 251.465 40.445 16.503 32.01 7.473 31.686 0 38.981 0 48.016v48.068c0 8.625 6.835 15.645 15.453 15.999 191.179 7.839 344.627 161.316 352.465 352.465.353 8.618 7.373 15.453 15.999 15.453h48.068c9.034-.001 16.329-7.474 16.005-16.504z"/></svg>`
var SvgRssSquare x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 448 512" width="1em" height="1em"><path d="M400 32H48C21.49 32 0 53.49 0 80v352c0 26.51 21.49 48 48 48h352c26.51 0 48-21.49 48-48V80c0-26.51-21.49-48-48-48zM112 416c-26.51 0-48-21.49-48-48s21.49-48 48-48 48 21.49 48 48-21.49 48-48 48zm157.533 0h-34.335c-6.011 0-11.051-4.636-11.442-10.634-5.214-80.05-69.243-143.92-149.123-149.123-5.997-.39-10.633-5.431-10.633-11.441v-34.335c0-6.535 5.468-11.777 11.994-11.425 110.546 5.974 198.997 94.536 204.964 204.964.352 6.526-4.89 11.994-11.425 11.994zm103.027 0h-34.334c-6.161 0-11.175-4.882-11.427-11.038-5.598-136.535-115.204-246.161-251.76-251.76C68.882 152.949 64 147.935 64 141.774V107.44c0-6.454 5.338-11.664 11.787-11.432 167.83 6.025 302.21 141.191 308.205 308.205.232 6.449-4.978 11.787-11.432 11.787z"/></svg>`
var SvgPrint x.Str = `<svg class="svg-icon fill-fg" viewBox="0 0 512 512" width="1em" height="1em"><path d="M464 192h-16V81.941a24 24 0 0 0-7.029-16.97L383.029 7.029A24 24 0 0 0 366.059 0H88C74.745 0 64 10.745 64 24v168H48c-26.51 0-48 21.49-48 48v132c0 6.627 5.373 12 12 12h52v104c0 13.255 10.745 24 24 24h336c13.255 0 24-10.745 24-24V384h52c6.627 0 12-5.373 12-12V240c0-26.51-21.49-48-48-48zm-80 256H128v-96h256v96zM128 224V64h192v40c0 13.2 10.8 24 24 24h40v96H128zm304 72c-13.254 0-24-10.746-24-24s10.746-24 24-24 24 10.746 24 24-10.746 24-24 24z"/></svg>`

// Short fur "current link".
func Cur(page Ipage, href string) x.Attrs {
	if page.GetLink() == href {
		return AP(`href`, href, `aria-current`, `page`)
	}
	return AP(`href`, href)
}

func aId(val string) (_ x.Attr) {
	if val != `` {
		return x.Attr{`id`, val}
	}
	return
}

func NoscriptInteractivity() x.Elem {
	return E(`noscript`, AP(`class`, `block fg-blue`),
		`This page has interactive features that require JavaScript. Consider enabling JS.`,
	)
}

func Script(src string) x.Elem {
	return E(`script`, AP(`src`, src, `type`, `module`))
}

func AttrsHidden(ok bool) x.Attrs {
	if ok {
		return AP(`hidden`, `true`)
	}
	return nil
}
