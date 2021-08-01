package main

import (
	x "github.com/mitranim/gax"
)

func Page404(site Site, page Page) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)
		E(`div`, A{aRole(`main`), aId(`main`), aClass("wid-lim fan-typo")}, func() {
			E(`h2`, nil, page.GetTitle())
			E(`p`, nil, `Sorry, this page is not found.`)
			E(`p`, nil, func() {
				E(`a`, A{aHref(`/`)}, `Return to homepage.`)
			})
		})
		Footer(E, page)
	})
}

func PageIndex(site Site, page Page) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(page.MakeMd()),
		)
		Footer(E, page)
	})
}

func PagePosts(site Site, page Page) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)

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

		Footer(E, page)
	})
}

func PageWorks(site Site, page Page) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(page.MakeMd()),
		)
		Footer(E, page)
	})
}

func PageResume(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo pad-top-1 pad-bot-2")},
			x.Bytes(page.MakeMd()),
		)
	})
}

func PageDemos(site Site, page Page) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("wid-lim fan-typo")},
			x.Bytes(page.MakeMd()),
		)
		Footer(E, page)
	})
}

func PagePost(site Site, page Post) []byte {
	return Html(page, func(E E) {
		Navbar(E, page)
		E(`div`, A{aRole(`main`), aId(`main`), aClass("wid-lim fan-typo flex-1 flex col-sta-str gap-ver-2")}, func(b *Bui) {
			E(`article`, A{aRole("article"), aClass(`fan-typo`)},
				func() {
					// Should be kept in sync with "MdRen.RenderNode" logic for headings
					E(`h1`, nil, x.Bytes(HEADING_PREFIX), page.Title)
					if page.Description != "" {
						E(`p`, A{aRole("doc-subtitle"), aClass("size-large italic")}, page.Description)
					}
					if page.TimeString() != "" {
						E(`p`, A{aClass("fg-gray-close size-small")}, page.TimeString())
					}
				},
				x.Bytes(page.MakeMd()),
			)

			E(`hr`, A{aStyle("margin-top: auto")})

			E(`div`, A{aClass("gap-ver-1")}, func() {
				E(`p`, nil, func(b *Bui) {
					T := b.EscString

					T(`This blog currently doesn't support comments. Feel free to `)
					Exta(E, "https://twitter.com/mitranim", "tweet")
					T(` at me, email to `)
					E(`a`, A{aHref(`mailto:me@mitranim.com?subject=Re: ` + page.Title)}, `me@mitranim.com`)
					T(`, or use the `)
					E(`a`, A{aHref("/#contacts")}, `other contacts.`)
					T(`.`)
				})

				FeedLinks(E)
			})
		})
		Footer(E, page)
	})
}