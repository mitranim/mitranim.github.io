package main

import (
	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
)

type Page struct {
	Path        string
	Title       string
	Description string
	MdTpl       []byte
	Type        string
	Image       string
	GlobalClass string
	MdHtml      []byte // Compiled once and reused, if necessary.
}

func (self Page) GetPath() string        { return self.Path }
func (self Page) GetTitle() string       { return self.Title }
func (self Page) GetDescription() string { return self.Description }
func (self Page) GetType() string        { return self.Type }
func (self Page) GetImage() string       { return self.Image }
func (self Page) GetGlobalClass() string { return self.GlobalClass }

func (self Page) Make(site Site) {
	panic(gg.Errf(`"Make" is not implemented for page %#v`, self))
}

func (self Page) MdOnce(val any) x.Bui {
	if self.MdTpl != nil && self.MdHtml == nil {
		self.MdHtml = self.Md(val, MdOpt{})
	}
	return self.MdHtml
}

func (self Page) Md(val any, opt MdOpt) x.Bui {
	defer gg.Detailf(`unable to parse and render %q as Markdown`, self.Path)
	return MdTplToHtml(self.MdTpl, opt, val)
}

func (self Page) GetLink() string {
	return ensureLeadingSlash(trimExt(self.GetPath()))
}

func initSitePages() []Ipage {
	return []Ipage{
		Page404{Page{
			Path:  `404.html`,
			Title: `Page Not Found`,
		}},
		PageIndex{Page{
			Path:        `index.html`,
			Title:       `about:mitranim`,
			Description: `About me: bio, works, posts`,
			MdTpl:       readTemplate(`index.md`),
		}},
		PageWorks{
			Page: Page{
				Path:        `works.html`,
				Title:       `Works`,
				Description: `Software I'm involved in`,
				MdTpl:       readTemplate(`works.md`),
			},
			Works: initWorks(),
		},
		PagePosts{Page{
			Path:        `posts.html`,
			Title:       `Blog Posts`,
			Description: `Random notes and thoughts`,
		}},
		PageGames{Page{
			Path:        `games.html`,
			Title:       `Game Recommendations`,
			Description: `Collection of games I've played, with impressions and recommendations`,
		}},
		PageDemos{Page{
			Path:        `demos.html`,
			Title:       `Demos`,
			Description: `Silly little demos`,
			MdTpl:       readTemplate(`demos.md`),
		}},
		PageResume{Page{
			Path:        `resume.html`,
			Title:       `Resume`,
			Description: `Nelo Mitranim's resume`,
			MdTpl:       readTemplate(`resume.md`),
		}},
	}
}

type Page404 struct{ Page }

func (self Page404) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMainArticleMd(),
			E(`h2`, nil, self.GetTitle()),
			E(`p`, nil, `Sorry, this page is not found.`),
			E(`p`, nil, E(`a`, AP(`href`, `/`), `Return to homepage.`)),
		),
	))
}

type PageIndex struct{ Page }

func (self PageIndex) GetLink() string { return `/` }

func (self PageIndex) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(), self.MdOnce(self)),
	))
}

type PagePosts struct{ Page }

func (self PagePosts) Make(site Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMain().Add(`class`, `post-previews`),
			E(`h1`, nil, `Blog Posts`),

			func(bui B) {
				src := site.ListedPosts()

				if len(src) > 0 {
					for _, post := range src {
						self.PostPreview(bui, post)
					}
				} else {
					bui.E(`p`, nil, `Oops! It appears there are no public posts yet.`)
				}
			},

			E(`h1`, nil, `Feed Links`),
			FeedLinks(),
		),
	))
}

func (self PagePosts) PostPreview(bui B, src PagePost) {
	bui.E(`div`, AP(`class`, `post-preview`), func() {
		bui.E(`h2`, nil,
			E(`a`, AP(`href`, src.GetLink()), src.Title),
		)
		if src.Description != `` {
			bui.E(`p`, nil, src.Description)
		}
		if src.TimeString() != `` {
			bui.E(`p`, AP(`class`, `fg-gray-near size-small`), src.TimeString())
		}
	})
}

type PageWorks struct {
	Page
	Works []Work
}

func (self PageWorks) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(), self.MdOnce(self)),
	))
}

func (self PageWorks) Table() x.Bui {
	return F(
		E(`table`, AP(`class`, `table sm-hide`),
			E(`thead`, nil,
				E(`th`, nil, `Name`),
				E(`th`, nil, `Desc`),
				E(`th`, nil, `Role`),
				E(`th`, nil, `Tech`),
				E(`th`, nil, `Start`),
				E(`th`, nil, `Status/End`),
			),
			E(`tbody`, nil, func(bui B) {
				for _, work := range self.Works {
					bui.E(`tr`, nil,
						E(`td`, nil, LinkExt(urlParse(work.Link).String(), work.Name)),
						E(`td`, nil, MdToHtmlStr(work.Desc)),
						E(`td`, AP(`class`, `fg-gray-near`), work.Role),
						E(`td`, AP(`class`, `fg-gray-near`), work.Tech),
						E(`td`, AP(`class`, `fg-gray-near`), work.Start),
						E(`td`, AP(`class`, `fg-gray-near wspace-nowrap`), work.StatusEnd),
					)
				}
			}),
		),
	)
}

func (self PageWorks) List() x.Bui {
	return F(
		E(`ul`, AP(`class`, `non-sm-hide`), func(bui B) {
			for _, work := range self.Works {
				bui.E(`li`, nil,
					LinkExt(work.Link, work.Name),
					` `,
					E(`span`, AP(`class`, `fg-gray-near`), `(`, work.Meta, `)`),
					` `,
					MdToHtmlStr(work.Desc),
				)
			}
		}),
	)
}

type PageGames struct{ Page }

func (self PageGames) Make(site Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMain().Add(`class`, `article`),
			self.Head(),
			self.Content(site),
		),
	))
}

func (self PageGames) Head() x.Ren {
	return F(
		E(`h1`, nil, self.Title),
		E(`details`, AP(`class`, `details details-spaced`),
			E(`summary`, AP(`class`, `summary`), `Click for additional notes.`),
			MdToHtmlStr(`

These are my current, modern recommendations. There are many other games I've
greatly enjoyed, which I would not recommend right now, either because they're
too outdated (e.g. Etherlords 2), or because the online community that made
them great no longer exists (e.g. WoW).

Even if you prefer MacOS or Linux for general use, you should use a dedicated
Windows system for games. Many games don't exist on other platforms, or take
years to release a port, usually with compatibility issues and poor
performance. Many games have essential mods only available on Windows. Windows
also allows a much better selection of hardware.

Always, _always_ check [PC Gaming Wiki](https://pcgamingwiki.com) for essential
tweaks and mods for any given game. For many games, it's also worth using mods
from [Nexus Mods](https://nexusmods.com), but beware of spoilers.

`),
		),
	)
}

func (self PageGames) Content(site Site) x.Ren {
	src := site.Games.Listed()

	if gg.IsEmpty(src) {
		return self.PlaceholderEmpty()
	}

	return F(
		NoscriptInteractivity().AttrAdd(`class`, `mar-bot-1`),
		self.TimeSinks(src),
		self.Tags(src),
		self.GameGrid(src),
		Script(`/scripts/games.mjs`),
	)
}

func (PageGames) TimeSinks(src Games) x.Ren {
	vals := src.TimeSinks()
	if gg.IsEmpty(vals) {
		return nil
	}

	return E(
		`tag-likes`,
		AP(`class`, `tag-likes`, `data-role`, `filter`),
		// TODO clicking this should clear the filter.
		E(`span`, AP(`class`, `help`, `aria-label`, `combined by logical "or"`), `Time sinks:`),
		vals,
	)
}

func (PageGames) Tags(src Games) x.Ren {
	vals := src.Tags()
	if gg.IsEmpty(vals) {
		return nil
	}

	return E(
		`tag-likes`,
		AP(`class`, `tag-likes`, `data-role`, `filter`),
		// TODO clicking this should clear the filter.
		E(`span`, AP(`class`, `help`, `aria-label`, `combined by logical "and"`), `Tags:`),
		vals,
	)
}

func (self PageGames) GameGrid(src Games) x.Ren {
	return E(`filter-list`, AP(`class`, `game-grid`),
		self.PlaceholderNothingFound().AttrSet(`hidden`, `true`),
		gg.Map(src, self.GameGridItem),
	)
}

func (PageGames) GameGridItem(src Game) x.Ren {
	return E(`filter-item`, AP(`class`, `game-grid-item`),
		E(`img`, AP(
			`src`, src.Img,
			`class`, `game-grid-item-img`,
		)),
		E(`h3`, nil, src.RenderName()),
		func(bui B) {
			if gg.IsNotZero(src.Desc) {
				bui.Child(MdToHtmlStr(src.Desc))
			}
			if gg.IsNotZero(src.TimeSink) || gg.IsNotEmpty(src.Tags) {
				bui.E(`div`, AP(`class`, `tag-likes`),
					src.TimeSink,
					gg.Sorted(src.Tags),
				)
			}
		},
	)
}

func (self PageGames) PlaceholderEmpty(src ...x.Attr) x.Elem {
	return self.Placeholder(
		`Oops! It appears there are no game recommendations yet.`,
	)
}

func (self PageGames) PlaceholderNothingFound(src ...x.Attr) x.Elem {
	return self.Placeholder(
		`Nothing found. Try changing the filters.`,
	)
}

func (self PageGames) Placeholder(text string) x.Elem {
	return E(
		`p`,
		AP(`is`, `filter-placeholder`, `class`, `filter-placeholder`),
		text,
	)
}

type PageResume struct{ Page }

func (self PageResume) Make(site Site) {
	index := PageByType[PageIndex](site)
	works := PageByType[PageWorks](site)

	PageWrite(self, Html(
		self,
		E(`article`, AttrsMainArticleMd().Add(`class`, `pad-body`),
			self.MdOnce(self),
			index.Md(index, MdOpt{}),
			MdToHtmlStr(`# Works`),
			works.Md(works, MdOpt{HeadingLevelOffset: 1}),
		),
	))
}

type PageDemos struct{ Page }

func (self PageDemos) Make(_ Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(), self.MdOnce(self)),
	))
}

func (self PagePost) Render(_ Site) x.Bui {
	return HtmlCommon(
		self,
		E(`article`, AttrsMainArticleMd(),
			// Should be kept in sync with `MdRen.RenderNode` logic for headings.
			E(`h1`, A(HEADING_PREFIX), self.Title),
			func(bui B) {
				if self.Description != `` {
					bui.E(`p`, AP(`role`, `doc-subtitle`, `class`, `size-large italic`), self.Description)
				}
				if self.TimeString() != `` {
					bui.E(`p`, AP(`class`, `fg-gray-near size-small`), self.TimeString())
				}
			},
			self.MdOnce(self),
		),
		// TODO avoid spamming horizontal padding classes.
		E(`hr`, AP(`class`, `hr mar-ver-1 pad-hor-body`)),
		PostsFooterLess().AttrAdd(`class`, `pad-hor-body`),
		FeedLinks().AttrAdd(`class`, `pad-hor-body`),
	)
}

// nolint:deadcode
func PostsFooterMore(page Ipage) x.Elem {
	return E(`p`, nil,
		`This blog currently doesn't support comments. Feel free to `,
		LinkExt(`https://twitter.com/mitranim`, `tweet`),
		` at me, email to `,
		E(`a`, AP(`href`, mailto(page.GetTitle())), EMAIL),
		`, or use the `,
		E(`a`, AP(`href`, `/#contacts`), `other contacts`),
		`.`,
	)
}

func PostsFooterLess() x.Elem {
	return E(`p`, nil,
		`This blog currently doesn't support comments. Write to me via `,
		E(`a`, AP(`href`, `/#contacts`, `class`, `link-deco`), `contacts`),
		`.`,
	)
}
