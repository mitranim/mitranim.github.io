package main

import (
	"strings"
	tt "text/template"
	"time"

	"github.com/gotidy/ptr"
	x "github.com/mitranim/gax"
	"github.com/mitranim/try"
)

func initSite() Site {
	tpl := makeTpl("")

	return Site{
		Tpl: tpl,
		Ipages: Ipages{
			Page{
				Path:  "404.html",
				Title: "Page Not Found",
				Fun:   Page404,
			},
			Page{
				Path:        "index.html",
				Title:       "about:mitranim",
				Description: "About me: bio, works, posts",
				MdTpl:       tryRead(fpj(TEMPLATE_DIR, "index.md")),
				Fun:         PageIndex,
			},
			Page{
				Path:        "works.html",
				Title:       "Works",
				Description: "Software I'm involved in",
				MdTpl:       tryRead(fpj(TEMPLATE_DIR, "works.md")),
				Fun:         PageWorks,
			},
			Page{
				Path:        "posts.html",
				Title:       "Blog Posts",
				Description: "Random notes and thoughts",
				Fun:         PagePosts,
			},
			Page{
				Path:        "demos.html",
				Title:       "Demos",
				Description: "Silly little demos",
				MdTpl:       tryRead(fpj(TEMPLATE_DIR, "demos.md")),
				Fun:         PageDemos,
			},
			Page{
				Path:        "resume.html",
				Title:       "Resume",
				Description: "Nelo Mitranim's resume",
				MdTpl:       tryRead(fpj(TEMPLATE_DIR, "resume.md")),
				GlobalClass: "color-scheme-light",
				Fun:         PageResume,
			},
			Post{
				Page: Page{
					Path:        "posts/spaces-tabs.html",
					Title:       "Always Spaces, Never Tabs",
					Description: "Objective arguments that decided my personal preference",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/spaces-tabs.md")),
				},
				PublishedAt: tryTimePtr("2020-10-23T06:48:15Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/lisp-sexpr-hacks.html",
					Title:       "Hacks around S-expressions in Lisps",
					Description: "How far people are willing to go to get prefix and infix in a Lisp syntax",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/lisp-sexpr-hacks.md")),
				},
				PublishedAt: tryTimePtr("2020-10-21T06:34:24Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/lang-var-minus.html",
					Title:       "Language Design: Gotchas With Variadic Minus",
					Description: "Treating the minus operator as a function can be tricky and dangerous",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/lang-var-minus.md")),
				},
				PublishedAt: tryTimePtr("2020-10-17T07:20:06Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/lang-case-conventions.html",
					Title:       "Language Design: Case Conventions",
					Description: "Objective arguments to solve case conventions and move on",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/lang-case-conventions.md")),
				},
				PublishedAt: tryTimePtr("2020-10-16T15:30:41Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/lang-homoiconic.html",
					Title:       "Language Design: Homoiconicity",
					Description: "Thoughts on homoiconicity, an interesting language quality seen in Lisps",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/lang-homoiconic.md")),
				},
				PublishedAt: tryTimePtr("2020-10-16T12:41:58Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/warframe-headcanon.html",
					Title:       "Warframe Headcanon (Spoilers)",
					Description: "Collection of Warframe headcanon co-authored with friends",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/warframe-headcanon.md")),
				},
				PublishedAt: tryTimePtr("2020-10-10T12:25:32Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/thoughts-on-the-egg.html",
					Title:       "Thoughts on The Egg: a short story by Andy Weir, animated by Kurzgesagt",
					Description: "",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/thoughts-on-the-egg.md")),
				},
				PublishedAt: tryTimePtr("2020-04-30T08:25:16Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/gameplay-conjecture.html",
					Title:       "Gameplay Conjecture",
					Description: "Amount of gameplay ≈ amount of required decisions",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/gameplay-conjecture.md")),
				},
				IsListed: !FLAGS.PROD,
			},
			Post{
				Page: Page{
					Path:        "posts/tips-and-tricks-doom-2016.html",
					Title:       "Tips and Tricks: Doom 2016",
					Description: "General tips, notes on difficulty, enemies, runes, weapons",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/tips-and-tricks-doom-2016.md")),
				},
				PublishedAt: tryTimePtr("2019-04-25T12:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/game-impressions-doom-2016.html",
					Title:       "Game Impressions: Doom 2016",
					Description: "I really like Doom 2016, here's why",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/game-impressions-doom-2016.md")),
				},
				PublishedAt: tryTimePtr("2019-04-25T11:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/astrotips.html",
					Title:       "Announcing Astrotips: Video Guides on Astroneer",
					Description: "A series of video guides, tips and tricks on Astroneer, an amazing space exploration and building game",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/astrotips.md")),
				},
				PublishedAt: tryTimePtr("2019-02-22T11:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/camel-case-abbr.html",
					Title:       "Don't Abbreviate in CamelCase",
					Description: `CamelCase identifiers should avoid abbreviations, e.g. "JsonText" rather than "JSONText"`,
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/camel-case-abbr.md")),
				},
				PublishedAt: tryTimePtr("2019-01-17T07:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/remove-from-go.html",
					Title:       "Things I Would Remove From Go",
					Description: "If less is more, Go could gain by losing weight",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/remove-from-go.md")),
				},
				PublishedAt: tryTimePtr("2019-01-15T01:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/back-from-hiatus-2019.html",
					Title:       "Back from Hiatus (2019)",
					Description: "Back to blogging after three and a half years",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/back-from-hiatus-2019.md")),
				},
				PublishedAt: tryTimePtr("2019-01-15T00:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/cheating-for-performance-pjax.html",
					Title:       "Cheating for Performance with Pjax",
					Description: "Faster page transitions, for free",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/cheating-for-performance-pjax.md")),
				},
				RedirFrom:   []string{"thoughts/cheating-for-performance-pjax.html"},
				PublishedAt: tryTimePtr("2015-07-25T00:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/cheating-for-website-performance.html",
					Title:       "Cheating for Website Performance",
					Description: "Frontend tips for speeding up websites",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/cheating-for-website-performance.md")),
				},
				RedirFrom:   []string{"thoughts/cheating-for-website-performance.html"},
				PublishedAt: tryTimePtr("2015-03-11T00:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/keeping-things-simple.html",
					Title:       "Keeping Things Simple",
					Description: "Musings on simplicity in programming",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/keeping-things-simple.md")),
				},
				RedirFrom:   []string{"thoughts/keeping-things-simple.html"},
				PublishedAt: tryTimePtr("2015-03-10T00:00:00Z"),
				IsListed:    true,
			},
			Post{
				Page: Page{
					Path:        "posts/next-generation-today.html",
					Title:       "Next Generation Today",
					Description: "EcmaScript 2015/2016 workflow with current web frameworks",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/next-generation-today.md")),
				},
				RedirFrom:   []string{"thoughts/next-generation-today.html"},
				PublishedAt: tryTimePtr("2015-05-18T00:00:00Z"),
				IsListed:    false,
			},
			Post{
				Page: Page{
					Path:        "posts/old-posts.html",
					Title:       "Old Posts",
					Description: "some old stuff from around the net",
					MdTpl:       tryRead(fpj(TEMPLATE_DIR, "posts/old-posts.md")),
				},
				RedirFrom:   []string{"thoughts/old-posts.html"},
				PublishedAt: tryTimePtr("2015-01-01T00:00:00Z"),
				IsListed:    true,
			},
		},
	}
}

type Site struct {
	Ipages
	Tpl *tt.Template
}

type Ipages []Ipage

func (self Ipages) Posts() (out []Post) {
	for _, val := range self {
		switch val := val.(type) {
		case Post:
			out = append(out, val)
		}
	}
	return
}

func (self Ipages) ListedPosts() (out []Post) {
	for _, val := range self.Posts() {
		if val.IsListed {
			out = append(out, val)
		}
	}
	return
}

type Ipage interface {
	GetPath() string
	GetTitle() string
	GetDescription() string
	GetType() string
	GetImage() string
	GetGlobalClass() string
	Make(Site) error
}

type Page struct {
	Path        string
	Title       string
	Description string
	MdTpl       []byte
	Type        string
	Image       string
	GlobalClass string
	Fun         func(Site, Page) []byte
}

func (self Page) GetPath() string        { return self.Path }
func (self Page) GetTitle() string       { return self.Title }
func (self Page) GetDescription() string { return self.Description }
func (self Page) GetType() string        { return self.Type }
func (self Page) GetImage() string       { return self.Image }
func (self Page) GetGlobalClass() string { return self.GlobalClass }

func (self Page) Make(site Site) error {
	return writePublic(self.Path, self.Fun(site, self))
}

type Post struct {
	Page
	RedirFrom   []string
	PublishedAt *time.Time
	UpdatedAt   *time.Time
	IsListed    bool
}

func (self Post) ExistsAsFile() bool {
	return self.PublishedAt != nil || !FLAGS.PROD
}

func (self Post) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

func (self Post) UrlFromSiteRoot() string {
	return ensureLeadingSlash(trimExt(self.Path))
}

// Somewhat inefficient but shouldn't be measurable.
func (self Post) TimeString() string {
	var out []string

	if self.PublishedAt != nil {
		out = append(out, `published `+timeFmtHuman(*self.PublishedAt))
		if self.UpdatedAt != nil {
			out = append(out, `updated `+timeFmtHuman(*self.UpdatedAt))
		}
	}

	return strings.Join(out, ", ")
}

func (self Post) Make(site Site) (err error) {
	defer try.Rec(&err)

	try.To(writePublic(self.Path, PagePost(site, self)))

	for _, path := range self.RedirFrom {
		try.To(writePublic(path, Ebui(func(E E) {
			E(`meta`, A{{`http-equiv`, `refresh`}, {`content`, `0;URL='` + self.UrlFromSiteRoot() + `'`}})
		}).Bytes()))
	}

	return
}

func (self Post) FeedItem() FeedItem {
	href := siteBase() + self.UrlFromSiteRoot()

	return FeedItem{
		XmlBase:     href,
		Title:       self.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: self.Page.Description,
		Id:          href,
		Published:   self.PublishedAt,
		Updated:     timeCoalesce(self.PublishedAt, self.UpdatedAt, ptr.Time(time.Now().UTC())),
		Content:     Ebui(func(E E) { FeedPostLayout(E, self) }).String(),
	}
}

func Page404(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`div`, A{aRole(`main`), aId(`main`), aClass("fancy-typography")}, func() {
			E(`h2`, nil, page.GetTitle())
			E(`p`, nil, `Sorry, this page is not found.`)
			E(`p`, nil, func() {
				E(`a`, A{aHref(`/`)}, `Return to homepage.`)
			})
		})
	})
}

func PageIndex(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("fancy-typography")},
			x.Bytes(tryMd(page.MdTpl, page)),
		)
	})
}

func PagePosts(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`div`, A{aRole(`main`), aId(`main`), aClass("flex col-start-stretch gaps-v-4")}, func() {
			E(`div`, A{aClass("fancy-typography gaps-v-2")}, func() {
				E(`h1`, nil, `Blog Posts`)

				posts := site.ListedPosts()

				if len(posts) > 0 {
					for _, post := range posts {
						E(`div`, A{aClass("gaps-v-1")}, func() {
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
			})

			E(`div`, A{aClass("fancy-typography")}, func() {
				E(`h1`, nil, `Feed Links`)
				FeedLinks(E)
			})
		})
	})
}

func PageWorks(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("fancy-typography")},
			x.Bytes(tryMd(page.MdTpl, page)),
		)
	})
}

func PageResume(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("fancy-typography limit-width padding-t-1 padding-b-2")},
			x.Bytes(tryMd(page.MdTpl, page)),
		)
	})
}

func PageDemos(site Site, page Page) []byte {
	return Html(page, func(E E) {
		E(`article`, A{aRole(`main article`), aId(`main`), aClass("fancy-typography")},
			x.Bytes(tryMd(page.MdTpl, page)),
		)
	})
}

func PagePost(site Site, page Post) []byte {
	return Html(page, func(E E) {
		E(`div`, A{aRole(`main`), aId(`main`), aClass("fancy-typography flex-1 flex col-start-stretch gaps-v-2")}, func(b *Bui) {
			E(`article`, A{aRole("article"), aClass(`fancy-typography`)},
				func() {
					// Should be kept in sync with "MdRen.RenderNode" logic for headings
					E(`h1`, nil, x.Bytes(HEADING_PREFIX), page.Title)
					if page.Description != "" {
						E(`p`, A{aRole("doc-subtitle"), aClass("size-large font-italic")}, page.Description)
					}
					if page.TimeString() != "" {
						E(`p`, A{aClass("fg-gray-close size-small")}, page.TimeString())
					}
				},
				x.Bytes(tryMd(page.MdTpl, page)),
			)

			E(`hr`, A{aStyle("margin-top: auto")})

			E(`div`, A{aClass("gaps-v-1")}, func() {
				E(`p`, nil, func(b *Bui) {
					text := b.EscString

					text(`This blog currently doesn't support comments. Feel free to `)
					Exta(E, "https://twitter.com/mitranim", "tweet")
					text(` at me, email to `)
					E(`a`, A{aHref(`mailto:me@mitranim.com?subject=Re: ` + page.Title)}, `me@mitranim.com`)
					text(`, or use the `)
					E(`a`, A{aHref("/#contacts")}, `other contacts.`)
					text(`.`)
				})

				FeedLinks(E)
			})
		})
	})
}
