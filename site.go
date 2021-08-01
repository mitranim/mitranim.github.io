package main

func initSite() (out Site) {
	out = append(out, initPages()...)
	out = append(out, initPosts()...)
	return
}

func initPages() Site {
	return Site{
		Page{
			Path:  "404.html",
			Title: "Page Not Found",
			Fun:   Page404,
		},
		Page{
			Path:        "index.html",
			Title:       "about:mitranim",
			Description: "About me: bio, works, posts",
			MdTpl:       readFile(fpj(TEMPLATE_DIR, "index.md")),
			Fun:         PageIndex,
		},
		Page{
			Path:        "works.html",
			Title:       "Works",
			Description: "Software I'm involved in",
			MdTpl:       readFile(fpj(TEMPLATE_DIR, "works.md")),
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
			MdTpl:       readFile(fpj(TEMPLATE_DIR, "demos.md")),
			Fun:         PageDemos,
		},
		Page{
			Path:        "resume.html",
			Title:       "Resume",
			Description: "Nelo Mitranim's resume",
			MdTpl:       readFile(fpj(TEMPLATE_DIR, "resume.md")),
			GlobalClass: "color-scheme-light",
			Fun:         PageResume,
		},
	}
}

func initPosts() Site {
	return Site{
		Post{
			Page: Page{
				Path:        "posts/spaces-tabs.html",
				Title:       "Always Spaces, Never Tabs",
				Description: "Objective arguments that decided my personal preference",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/spaces-tabs.md")),
			},
			PublishedAt: timeParsePtr("2020-10-23T06:48:15Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/lisp-sexpr-hacks.html",
				Title:       "Hacks around S-expressions in Lisps",
				Description: "How far people are willing to go to get prefix and infix in a Lisp syntax",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/lisp-sexpr-hacks.md")),
			},
			PublishedAt: timeParsePtr("2020-10-21T06:34:24Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/lang-var-minus.html",
				Title:       "Language Design: Gotchas With Variadic Minus",
				Description: "Treating the minus operator as a function can be tricky and dangerous",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/lang-var-minus.md")),
			},
			PublishedAt: timeParsePtr("2020-10-17T07:20:06Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/lang-case-conventions.html",
				Title:       "Language Design: Case Conventions",
				Description: "Objective arguments to solve case conventions and move on",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/lang-case-conventions.md")),
			},
			PublishedAt: timeParsePtr("2020-10-16T15:30:41Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/lang-homoiconic.html",
				Title:       "Language Design: Homoiconicity",
				Description: "Thoughts on homoiconicity, an interesting language quality seen in Lisps",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/lang-homoiconic.md")),
			},
			PublishedAt: timeParsePtr("2020-10-16T12:41:58Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/warframe-headcanon.html",
				Title:       "Warframe Headcanon (Spoilers)",
				Description: "Collection of Warframe headcanon co-authored with friends",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/warframe-headcanon.md")),
			},
			PublishedAt: timeParsePtr("2020-10-10T12:25:32Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/thoughts-on-the-egg.html",
				Title:       "Thoughts on The Egg: a short story by Andy Weir, animated by Kurzgesagt",
				Description: "",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/thoughts-on-the-egg.md")),
			},
			PublishedAt: timeParsePtr("2020-04-30T08:25:16Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/gameplay-conjecture.html",
				Title:       "Gameplay Conjecture",
				Description: "Amount of gameplay ≈ amount of required decisions",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/gameplay-conjecture.md")),
			},
			IsListed: !FLAGS.PROD,
		},
		Post{
			Page: Page{
				Path:        "posts/tips-and-tricks-doom-2016.html",
				Title:       "Tips and Tricks: Doom 2016",
				Description: "General tips, notes on difficulty, enemies, runes, weapons",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/tips-and-tricks-doom-2016.md")),
			},
			PublishedAt: timeParsePtr("2019-04-25T12:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/game-impressions-doom-2016.html",
				Title:       "Game Impressions: Doom 2016",
				Description: "I really like Doom 2016, here's why",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/game-impressions-doom-2016.md")),
			},
			PublishedAt: timeParsePtr("2019-04-25T11:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/astrotips.html",
				Title:       "Announcing Astrotips: Video Guides on Astroneer",
				Description: "A series of video guides, tips and tricks on Astroneer, an amazing space exploration and building game",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/astrotips.md")),
			},
			PublishedAt: timeParsePtr("2019-02-22T11:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/camel-case-abbr.html",
				Title:       "Don't Abbreviate in CamelCase",
				Description: `CamelCase identifiers should avoid abbreviations, e.g. "JsonText" rather than "JSONText"`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/camel-case-abbr.md")),
			},
			PublishedAt: timeParsePtr("2019-01-17T07:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/remove-from-go.html",
				Title:       "Things I Would Remove From Go",
				Description: "If less is more, Go could gain by losing weight",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/remove-from-go.md")),
			},
			PublishedAt: timeParsePtr("2019-01-15T01:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/back-from-hiatus-2019.html",
				Title:       "Back from Hiatus (2019)",
				Description: "Back to blogging after three and a half years",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/back-from-hiatus-2019.md")),
			},
			PublishedAt: timeParsePtr("2019-01-15T00:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/cheating-for-performance-pjax.html",
				Title:       "Cheating for Performance with Pjax",
				Description: "Faster page transitions, for free",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/cheating-for-performance-pjax.md")),
			},
			RedirFrom:   []string{"thoughts/cheating-for-performance-pjax.html"},
			PublishedAt: timeParsePtr("2015-07-25T00:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/cheating-for-website-performance.html",
				Title:       "Cheating for Website Performance",
				Description: "Frontend tips for speeding up websites",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/cheating-for-website-performance.md")),
			},
			RedirFrom:   []string{"thoughts/cheating-for-website-performance.html"},
			PublishedAt: timeParsePtr("2015-03-11T00:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/keeping-things-simple.html",
				Title:       "Keeping Things Simple",
				Description: "Musings on simplicity in programming",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/keeping-things-simple.md")),
			},
			RedirFrom:   []string{"thoughts/keeping-things-simple.html"},
			PublishedAt: timeParsePtr("2015-03-10T00:00:00Z"),
			IsListed:    true,
		},
		Post{
			Page: Page{
				Path:        "posts/next-generation-today.html",
				Title:       "Next Generation Today",
				Description: "EcmaScript 2015/2016 workflow with current web frameworks",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/next-generation-today.md")),
			},
			RedirFrom:   []string{"thoughts/next-generation-today.html"},
			PublishedAt: timeParsePtr("2015-05-18T00:00:00Z"),
			IsListed:    false,
		},
		Post{
			Page: Page{
				Path:        "posts/old-posts.html",
				Title:       "Old Posts",
				Description: "some old stuff from around the net",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, "posts/old-posts.md")),
			},
			RedirFrom:   []string{"thoughts/old-posts.html"},
			PublishedAt: timeParsePtr("2015-01-01T00:00:00Z"),
			IsListed:    true,
		},
	}
}
