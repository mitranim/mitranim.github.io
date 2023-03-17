package main

import (
	"net/url"

	x "github.com/mitranim/gax"
)

type Page404 struct{ Page }

func (self Page404) Make(_ Site) {
	pageWrite(self, Html(
		self,
		Header(self),
		E(`div`, AP(`role`, `main`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			E(`h2`, nil, self.GetTitle()),
			E(`p`, nil, `Sorry, this page is not found.`),
			E(`p`, nil, E(`a`, AP(`href`, `/`), `Return to homepage.`)),
		),
		Footer(self),
	))
}

type PageIndex struct{ Page }

func (self PageIndex) GetLink() string { return `/` }

func (self PageIndex) Make(_ Site) {
	pageWrite(self, Html(
		self,
		Header(self),
		E(`article`, AP(`role`, `main article`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			Bui(self.MdOnce(self)),
		),
		Footer(self),
	))
}

type PagePosts struct{ Page }

func (self PagePosts) Make(site Site) {
	pageWrite(self, Html(
		self,
		Header(self),

		E(`div`, AP(`role`, `main`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			E(`h1`, nil, `Blog Posts`),

			func(bui B) {
				posts := site.ListedPosts()

				if len(posts) > 0 {
					for _, post := range posts {
						bui.E(`div`, AP(`class`, `mar-top-2 gap-ver-1`), func() {
							bui.E(`h2`, nil,
								E(`a`, AP(`href`, post.GetLink()), post.Title),
							)
							if post.Description != `` {
								bui.E(`p`, nil, post.Description)
							}
							if post.TimeString() != `` {
								bui.E(`p`, AP(`class`, `fg-gray-close size-small`), post.TimeString())
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
		Footer(self),
	))
}

type PageWorks struct {
	Page
	Works []Work
}

func (self PageWorks) Make(_ Site) {
	pageWrite(self, Html(
		self,
		Header(self),
		E(`article`, AP(`role`, `main article`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			Bui(self.MdOnce(self)),
		),
		Footer(self),
	))
}

func (self PageWorks) Table() Bui {
	return F(
		E(`table`, AP(`class`, `sm-hide`),
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
						E(`td`, nil, Exta(urlParse(work.Link).String(), work.Name)),
						E(`td`, nil, x.Str(stringMdToHtml(work.Desc, nil))),
						E(`td`, AP(`class`, `fg-gray-close`), work.Role),
						E(`td`, AP(`class`, `fg-gray-close`), work.Tech),
						E(`td`, AP(`class`, `fg-gray-close`), work.Start),
						E(`td`, AP(`class`, `fg-gray-close wspace-nowrap`), work.StatusEnd),
					)
				}
			}),
		),
	)
}

func (self PageWorks) List() Bui {
	return F(
		E(`ul`, AP(`class`, `non-sm-hide`), func(bui B) {
			for _, work := range self.Works {
				bui.E(`li`, nil,
					Exta(work.Link, work.Name),
					` `,
					E(`span`, AP(`class`, `fg-gray-close`), `(`, work.Meta, `)`),
					` `,
					x.Str(stringMdToHtml(work.Desc, nil)),
				)
			}
		}),
	)
}

type PageResume struct{ Page }

func (self PageResume) Make(site Site) {
	index := PageByType[PageIndex](site)
	works := PageByType[PageWorks](site)

	pageWrite(self, Html(
		self,
		E(`article`, AP(`role`, `main article`, `id`, ID_MAIN, `class`, `wid-lim fan-typo pad-top-1 pad-bot-2`),
			Bui(self.MdOnce(self)),
			Bui(index.Md(index, nil)),
			x.Str(stringMdToHtml(`# Works`, nil)),
			Bui(works.Md(works, &MdOpt{HeadingLevelOffset: 1})),
		),
	))
}

type PageDemos struct{ Page }

func (self PageDemos) Make(_ Site) {
	pageWrite(self, Html(
		self,
		Header(self),
		E(`article`, AP(`role`, `main article`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			Bui(self.MdOnce(self)),
		),
		Footer(self),
	))
}

func (self PagePost) Render(_ Site) Bui {
	return Html(
		self,
		Header(self),
		E(`div`, AP(`role`, `main`, `id`, ID_MAIN, `class`, `wid-lim fan-typo`),
			E(`article`, AP(`role`, `article`, `class`, `fan-typo`),
				// Should be kept in sync with `MdRen.RenderNode` logic for headings.
				E(`h1`, nil, HEADING_PREFIX, self.Title),
				func(bui B) {
					if self.Description != `` {
						bui.E(`p`, AP(`role`, `doc-subtitle`, `class`, `size-large italic`), self.Description)
					}
					if self.TimeString() != `` {
						bui.E(`p`, AP(`class`, `fg-gray-close size-small`), self.TimeString())
					}
				},
				Bui(self.MdOnce(self)),
			),

			E(`hr`, nil),
			PostsFooterLess,

			FeedLinks,
		),
		Footer(self),
	)
}

// nolint:deadcode
func PostsFooterMore(page Ipage) x.Elem {
	return E(`p`, nil,
		`This blog currently doesn't support comments. Feel free to `,
		Exta(`https://twitter.com/mitranim`, `tweet`),
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

func mailto(subj string) string {
	return MAILTO.WithQuery(emailSubj(subj)).String()
}

func emailSubj(val string) url.Values {
	if val != `` {
		return url.Values{`subject`: {`Re: ` + val}}
	}
	return nil
}
