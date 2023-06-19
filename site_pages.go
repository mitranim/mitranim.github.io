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
		self.MdHtml = self.Md(val, nil)
	}
	return self.MdHtml
}

func (self Page) Md(val any, opt *MdOpt) x.Bui {
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
		E(`div`, AttrsMain().Set(`class`, `post-previews`),
			E(`h1`, nil, `Blog Posts`),

			func(bui B) {
				src := site.ListedPosts()

				if len(src) > 0 {
					for _, post := range src {
						bui.E(`div`, AP(`class`, `post-preview`), func() {
							bui.E(`h2`, nil,
								E(`a`, AP(`href`, post.GetLink()), post.Title),
							)
							if post.Description != `` {
								bui.E(`p`, nil, post.Description)
							}
							if post.TimeString() != `` {
								bui.E(`p`, AP(`class`, `fg-gray-near size-small`), post.TimeString())
							}
						})
					}
				} else {
					bui.E(`p`, nil, `Oops! It appears there are no public posts yet.`)
				}
			},

			E(`h1`, nil, `Feed Links`),
			FeedLinks,
		),
	))
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

type PageResume struct{ Page }

func (self PageResume) Make(site Site) {
	index := PageByType[PageIndex](site)
	works := PageByType[PageWorks](site)

	PageWrite(self, Html(
		self,
		// Top padding is a replacement for the missing header.
		E(`article`, AttrsMainArticleMd().Add(`class`, `pad-top-1`),
			self.MdOnce(self),
			index.Md(index, nil),
			MdToHtmlStr(`# Works`),
			works.Md(works, &MdOpt{HeadingLevelOffset: 1}),
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
			E(`h1`, nil, HEADING_PREFIX, self.Title),
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
		E(`hr`, AP(`class`, `hr mar-ver-1`)),
		PostsFooterLess,
		FeedLinks,
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

var PostsFooterLess = E(`p`, nil,
	`This blog currently doesn't support comments. Write to me via `,
	E(`a`, AP(`href`, `/#contacts`), `contacts`),
	`.`,
)
