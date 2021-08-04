package main

import (
	x "github.com/mitranim/gax"
)

type Page404 struct{ Page }

func (self Page404) Make(_ Site) {
	pageWrite(self, Html(
		self,
		Header(self),
		E(`div`, AP(`role`, `main`, `id`, `main`, `class`, `wid-lim fan-typo`),
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
		E(`article`, AP(`role`, `main article`, `id`, `main`, `class`, `wid-lim fan-typo`),
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

		E(`div`, AP(`role`, `main`, `id`, `main`, `class`, `wid-lim fan-typo`),
			E(`h1`, nil, `Blog Posts`),

			func(b B) {
				posts := site.ListedPosts()

				if len(posts) > 0 {
					for _, post := range posts {
						b.E(`div`, AP(`class`, "mar-top-2 gap-ver-1"), func() {
							b.E(`h2`, nil,
								E(`a`, AP(`href`, post.GetLink()), post.Title),
							)
							if post.Description != "" {
								b.E(`p`, nil, post.Description)
							}
							if post.TimeString() != "" {
								b.E(`p`, AP(`class`, "fg-gray-close size-small"), post.TimeString())
							}
						})
					}
				} else {
					b.E(`p`, nil, `Oops! It appears there are no public posts yet.`)
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
		E(`article`, AP(`role`, `main article`, `id`, `main`, `class`, `wid-lim fan-typo`),
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
				E(`th`, nil, `End`),
			),
			E(`tbody`, nil, func(b B) {
				for _, work := range self.Works {
					b.E(`tr`, nil,
						E(`td`, nil, Exta(parseUrl(work.Link).String(), work.Name)),
						E(`td`, nil, x.Str(stringMdToHtml(work.Desc, nil))),
						E(`td`, aFade, work.Role),
						E(`td`, aFade, work.Tech),
						E(`td`, aFade, work.Start),
						E(`td`, aFade, work.End),
					)
				}
			}),
		),
	)
}

func (self PageWorks) List() Bui {
	return F(
		E(`ul`, AP(`class`, `non-sm-hide`), func(b B) {
			for _, work := range self.Works {
				b.E(`li`, nil,
					Exta(work.Link, work.Name),
					` `,
					E(`span`, aFade, `(`, work.Meta(), `)`),
					` `,
					x.Str(stringMdToHtml(work.Desc, nil)),
				)
			}
		}),
	)
}

type PageResume struct{ Page }

func (self PageResume) Make(site Site) {
	index := site.PageByType(PageIndex{}).(PageIndex)
	works := site.PageByType(PageWorks{}).(PageWorks)

	pageWrite(self, Html(
		self,
		E(`article`, AP(`role`, `main article`, `id`, `main`, `class`, `wid-lim fan-typo pad-top-1 pad-bot-2`),
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
		E(`article`, AP(`role`, `main article`, `id`, `main`, `class`, `wid-lim fan-typo`),
			Bui(self.MdOnce(self)),
		),
		Footer(self),
	))
}

func (self PagePost) Render(_ Site) Bui {
	return Html(
		self,
		Header(self),
		E(`div`, AP(`role`, `main`, `id`, `main`, `class`, `wid-lim fan-typo`),
			E(`article`, AP(`role`, `article`, `class`, `fan-typo`),
				// Should be kept in sync with "MdRen.RenderNode" logic for headings
				E(`h1`, nil, HEADING_PREFIX, self.Title),
				func(b B) {
					if self.Description != "" {
						b.E(`p`, AP(`role`, "doc-subtitle", `class`, "size-large italic"), self.Description)
					}
					if self.TimeString() != "" {
						b.E(`p`, AP(`class`, "fg-gray-close size-small"), self.TimeString())
					}
				},
				Bui(self.MdOnce(self)),
			),

			E(`hr`, nil),

			E(`p`, nil,
				`This blog currently doesn't support comments. Feel free to `,
				Exta("https://twitter.com/mitranim", "tweet"),
				` at me, email to `,
				E(`a`, AP(`href`, `mailto:me@mitranim.com?subject=Re: `+self.Title), `me@mitranim.com`),
				`, or use the `,
				E(`a`, AP(`href`, "/#contacts"), `other contacts`),
				`.`,
			),
			FeedLinks,
		),
		Footer(self),
	)
}
