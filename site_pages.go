package main

import (
	"github.com/mitranim/gax"
	x "github.com/mitranim/gax"
)

type Page404 struct{ Page }

func (self Page404) Make(_ Site) {
	self.Write(Html(self, func(E E) {
		Navbar(E, self)
		E(`div`, A{aRole(`main`), aId(`main`), aClass("wid-lim fan-typo")}, func() {
			E(`h2`, nil, self.GetTitle())
			E(`p`, nil, `Sorry, this page is not found.`)
			E(`p`, nil, func() {
				E(`a`, A{aHref(`/`)}, `Return to homepage.`)
			})
		})
		Footer(E, self)
	}))
}

type PageIndex struct{ Page }

func (self PageIndex) Make(_ Site) {
	self.Write(Html(self, func(E E) {
		Navbar(E, self)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(self.MdOnce(self)),
		)
		Footer(E, self)
	}))
}

type PagePosts struct{ Page }

func (self PagePosts) Make(site Site) {
	self.Write(Html(self, func(E E) {
		Navbar(E, self)

		E(`div`, A{aRole(`main`), aId(`main`), aClass("wid-lim fan-typo")}, func() {
			E(`h1`, nil, `Blog Posts`)

			posts := site.ListedPosts()

			if len(posts) > 0 {
				for _, post := range posts {
					E(`div`, A{aClass("mar-top-2 gap-ver-1")}, func() {
						E(`h2`, nil, func() {
							E(`a`, A{aHref(post.UrlFromSiteRoot())}, post.Title)
						})
						if post.Description != "" {
							E(`p`, nil, post.Description)
						}
						if post.TimeString() != "" {
							E(`p`, A{aClass("fg-gray-close size-small")}, post.TimeString())
						}
					})
				}
			} else {
				E(`p`, nil, `Oops! It appears there are no public posts yet.`)
			}

			E(`h1`, nil, `Feed Links`)
			FeedLinks(E)
		})

		Footer(E, self)
	}))
}

type PageWorks struct {
	Page
	Works []Work
}

func (self PageWorks) Make(_ Site) {
	self.Write(Html(self, func(E E) {
		Navbar(E, self)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(self.MdOnce(self)),
		)
		Footer(E, self)
	}))
}

func (self PageWorks) Table() string {
	return Ebui(func(E E) {
		E(`table`, A{aClass(`sm-hide`)}, func() {
			E(`thead`, nil, func() {
				E(`th`, nil, `Name`)
				E(`th`, nil, `Desc`)
				E(`th`, nil, `Role`)
				E(`th`, nil, `Tech`)
				E(`th`, nil, `Start`)
				E(`th`, nil, `End`)
			})
			E(`tbody`, nil, func() {
				for _, work := range self.Works {
					E(`tr`, nil, func() {
						E(`td`, nil, func() { Exta(E, parseUrl(work.Link).String(), work.Name) })
						E(`td`, nil, gax.String(stringMdToHtml(work.Desc)))
						E(`td`, aFade, work.Role)
						E(`td`, aFade, work.Tech)
						E(`td`, aFade, work.Start)
						E(`td`, aFade, work.End)
					})
				}
			})
		})
	}).String()
}

func (self PageWorks) List() string {
	return Ebui(func(E E) {
		E(`ul`, A{aClass(`non-sm-hide`)}, func() {
			for _, work := range self.Works {
				E(`li`, nil, func(b *Bui) {
					T := b.NonEscString

					Exta(E, work.Link, work.Name)
					T(` `)
					E(`span`, aFade, `(`, work.Meta(), `)`)
					T(` `)
					T(stringMdToHtml(work.Desc))
				})
			}
		})
	}).String()
}

type PageResume struct{ Page }

func (self PageResume) Make(_ Site) {
	self.Write(Html(self, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo pad-top-1 pad-bot-2")},
			x.Bytes(self.MdOnce(self)),
		)
	}))
}

type PageDemos struct{ Page }

func (self PageDemos) Make(_ Site) {
	self.Write(Html(self, func(E E) {
		Navbar(E, self)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(self.MdOnce(self)),
		)
		Footer(E, self)
	}))
}

func (self PagePost) Render(_ Site) []byte {
	return Html(self, func(E E) {
		Navbar(E, self)
		E(`div`, A{aRole(`main`), aId(`main`), aClass("wid-lim fan-typo flex-1 flex col-sta-str gap-ver-2")}, func(b *Bui) {
			E(`article`, A{aRole("article"), aClass(`fan-typo`)},
				func() {
					// Should be kept in sync with "MdRen.RenderNode" logic for headings
					E(`h1`, nil, x.Bytes(HEADING_PREFIX), self.Title)
					if self.Description != "" {
						E(`p`, A{aRole("doc-subtitle"), aClass("size-large italic")}, self.Description)
					}
					if self.TimeString() != "" {
						E(`p`, A{aClass("fg-gray-close size-small")}, self.TimeString())
					}
				},
				x.Bytes(self.MdOnce(self)),
			)

			E(`hr`, A{aStyle("margin-top: auto")})

			E(`div`, A{aClass("gap-ver-1")}, func() {
				E(`p`, nil, func(b *Bui) {
					T := b.EscString

					T(`This blog currently doesn't support comments. Feel free to `)
					Exta(E, "https://twitter.com/mitranim", "tweet")
					T(` at me, email to `)
					E(`a`, A{aHref(`mailto:me@mitranim.com?subject=Re: ` + self.Title)}, `me@mitranim.com`)
					T(`, or use the `)
					E(`a`, A{aHref("/#contacts")}, `other contacts.`)
					T(`.`)
				})

				FeedLinks(E)
			})
		})
		Footer(E, self)
	})
}
