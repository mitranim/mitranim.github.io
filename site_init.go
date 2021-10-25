package main

func initSite() Site {
	return Site{
		Pages: initSitePages(),
		Posts: initSitePosts(),
	}
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
			MdTpl:       readFile(fpj(TEMPLATE_DIR, `index.md`)),
		}},
		PageWorks{
			Page: Page{
				Path:        `works.html`,
				Title:       `Works`,
				Description: `Software I'm involved in`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `works.md`)),
			},
			Works: initWorks(),
		},
		PagePosts{Page{
			Path:        `posts.html`,
			Title:       `Blog Posts`,
			Description: `Random notes and thoughts`,
		}},
		PageDemos{Page{
			Path:        `demos.html`,
			Title:       `Demos`,
			Description: `Silly little demos`,
			MdTpl:       readFile(fpj(TEMPLATE_DIR, `demos.md`)),
		}},
		PageResume{Page{
			Path:        `resume.html`,
			Title:       `Resume`,
			Description: `Nelo Mitranim's resume`,
			MdTpl:       readFile(fpj(TEMPLATE_DIR, `resume.md`)),
			// GlobalClass: `color-scheme-light`,
		}},
	}
}

func initSitePosts() []PagePost {
	return []PagePost{
		PagePost{
			Page: Page{
				Path:        `posts/goex.html`,
				Title:       `Use Exceptions in Go`,
				Description: `Go secretly favors exceptions. Using them has tremendous benefits.`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/goex.md`)),
			},
			PublishedAt: timeParse(`2021-09-07T16:17:36Z`),
			IsListed:    !FLAGS.PROD,
		},
		PagePost{
			Page: Page{
				Path:        `posts/spaces-tabs.html`,
				Title:       `Always Spaces, Never Tabs`,
				Description: `Objective arguments that decided my personal preference`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/spaces-tabs.md`)),
			},
			PublishedAt: timeParse(`2020-10-23T06:48:15Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lisp-sexpr-hacks.html`,
				Title:       `Hacks around S-expressions in Lisps`,
				Description: `How far people are willing to go to get prefix and infix in a Lisp syntax`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/lisp-sexpr-hacks.md`)),
			},
			PublishedAt: timeParse(`2020-10-21T06:34:24Z`),
			UpdatedAt:   timeParse(`2021-08-20T07:16:38Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-var-minus.html`,
				Title:       `Language Design: Gotchas With Variadic Minus`,
				Description: `Treating the minus operator as a function can be tricky and dangerous`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/lang-var-minus.md`)),
			},
			PublishedAt: timeParse(`2020-10-17T07:20:06Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-case-conventions.html`,
				Title:       `Language Design: Case Conventions`,
				Description: `Objective arguments to solve case conventions and move on`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/lang-case-conventions.md`)),
			},
			PublishedAt: timeParse(`2020-10-16T15:30:41Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-homoiconic.html`,
				Title:       `Language Design: Homoiconicity`,
				Description: `Thoughts on homoiconicity, an interesting language quality seen in Lisps`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/lang-homoiconic.md`)),
			},
			PublishedAt: timeParse(`2020-10-16T12:41:58Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/warframe-headcanon.html`,
				Title:       `Warframe Headcanon (Spoilers)`,
				Description: `Collection of Warframe headcanon co-authored with friends`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/warframe-headcanon.md`)),
			},
			PublishedAt: timeParse(`2020-10-10T12:25:32Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/thoughts-on-the-egg.html`,
				Title:       `Thoughts on The Egg: a short story by Andy Weir, animated by Kurzgesagt`,
				Description: ``,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/thoughts-on-the-egg.md`)),
			},
			PublishedAt: timeParse(`2020-04-30T08:25:16Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/gameplay-conjecture.html`,
				Title:       `Gameplay Conjecture`,
				Description: `Amount of gameplay ≈ amount of required decisions`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/gameplay-conjecture.md`)),
			},
			IsListed: !FLAGS.PROD,
		},
		PagePost{
			Page: Page{
				Path:        `posts/tips-and-tricks-doom-2016.html`,
				Title:       `Tips and Tricks: Doom 2016`,
				Description: `General tips, notes on difficulty, enemies, runes, weapons`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/tips-and-tricks-doom-2016.md`)),
			},
			PublishedAt: timeParse(`2019-04-25T12:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/game-impressions-doom-2016.html`,
				Title:       `Game Impressions: Doom 2016`,
				Description: `I really like Doom 2016, here's why`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/game-impressions-doom-2016.md`)),
			},
			PublishedAt: timeParse(`2019-04-25T11:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/astrotips.html`,
				Title:       `Announcing Astrotips: Video Guides on Astroneer`,
				Description: `A series of video guides, tips and tricks on Astroneer, an amazing space exploration and building game`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/astrotips.md`)),
			},
			PublishedAt: timeParse(`2019-02-22T11:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/camel-case-abbr.html`,
				Title:       `Don't Abbreviate in CamelCase`,
				Description: "CamelCase identifiers should avoid abbreviations, e.g. `JsonText` rather than `JSONText`",
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/camel-case-abbr.md`)),
			},
			PublishedAt: timeParse(`2019-01-17T07:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/remove-from-go.html`,
				Title:       `Things I Would Remove From Go`,
				Description: `If less is more, Go could gain by losing weight`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/remove-from-go.md`)),
			},
			PublishedAt: timeParse(`2019-01-15T01:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/back-from-hiatus-2019.html`,
				Title:       `Back from Hiatus (2019)`,
				Description: `Back to blogging after three and a half years`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/back-from-hiatus-2019.md`)),
			},
			PublishedAt: timeParse(`2019-01-15T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/cheating-for-performance-pjax.html`,
				Title:       `Cheating for Performance with Pjax`,
				Description: `Faster page transitions, for free`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/cheating-for-performance-pjax.md`)),
			},
			RedirFrom:   []string{`thoughts/cheating-for-performance-pjax.html`},
			PublishedAt: timeParse(`2015-07-25T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/cheating-for-website-performance.html`,
				Title:       `Cheating for Website Performance`,
				Description: `Frontend tips for speeding up websites`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/cheating-for-website-performance.md`)),
			},
			RedirFrom:   []string{`thoughts/cheating-for-website-performance.html`},
			PublishedAt: timeParse(`2015-03-11T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/keeping-things-simple.html`,
				Title:       `Keeping Things Simple`,
				Description: `Musings on simplicity in programming`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/keeping-things-simple.md`)),
			},
			RedirFrom:   []string{`thoughts/keeping-things-simple.html`},
			PublishedAt: timeParse(`2015-03-10T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/next-generation-today.html`,
				Title:       `Next Generation Today`,
				Description: `EcmaScript 2015/2016 workflow with current web frameworks`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/next-generation-today.md`)),
			},
			RedirFrom:   []string{`thoughts/next-generation-today.html`},
			PublishedAt: timeParse(`2015-05-18T00:00:00Z`),
			IsListed:    false,
		},
		PagePost{
			Page: Page{
				Path:        `posts/old-posts.html`,
				Title:       `Old Posts`,
				Description: `some old stuff from around the net`,
				MdTpl:       readFile(fpj(TEMPLATE_DIR, `posts/old-posts.md`)),
			},
			RedirFrom:   []string{`thoughts/old-posts.html`},
			PublishedAt: timeParse(`2015-01-01T00:00:00Z`),
			IsListed:    true,
		},
	}
}

func initWorks() []Work {
	return []Work{
		{
			Name:  `sublime-scss`,
			Link:  `https://github.com/mitranim/sublime-scss`,
			Desc:  `Redesigned CSS and SCSS syntaxes for Sublime Text. Built on open-ended principles. Designed for forward compatibility.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2021`,
		},
		{
			Name:  `ProstoPoi SSG`,
			Link:  `https://github.com/mitranim/pp`,
			Desc:  `Poi community website. Runs since 2014. Now converted from Django (Python) to static generation (JS), open sourced.`,
			Role:  `implementer`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `gax`,
			Link:  `https://github.com/mitranim/gax`,
			Desc:  `Simple system for writing HTML as Go code. Use normal Go conditionals, loops and functions. Benefit from typing and code analysis. Better performance than templating.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `imperouter`,
			Link:  `https://github.com/mitranim/imperouter`,
			Desc:  `Simple imperative router for hybrid SSR+SPA.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `jol`,
			Link:  `https://github.com/mitranim/jol`,
			Desc:  `JS Collection Classes. Tiny extensions on JS built-in classes, with nice features such as easy-to-use typed collections, dictionary with structured keys; and more.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `web-starter`,
			Link:  `https://github.com/mitranim/web-starter`,
			Desc:  `Starter templates for minimal web apps, from simplest to complex. Sucks less than X. Work in progress.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `afr`,
			Link:  `https://github.com/mitranim/afr`,
			Desc:  `Flexible library for serving files, with optional client integration for CSS injection and page reload.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `espo`,
			Link:  `https://github.com/mitranim/espo`,
			Desc:  `Observables via proxies. Particularly suited for UI programming. Supports automatic, implicit sub/resub and resource deinit.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `prax`,
			Link:  `https://github.com/mitranim/prax`,
			Desc:  `Simple rendering library for hybrid SSR+SPA. Superior replacement for rendering frameworks like React.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `jtg`,
			Link:  `https://github.com/mitranim/jtg`,
			Desc:  `"JS Task Group". Simple JS-based replacement for Make, Gulp, etc.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2021`,
		},
		{
			Name:  `try`,
			Link:  `https://github.com/mitranim/try`,
			Desc:  "Shorter error handling in Go. Supports two styles: explicit `try` and exceptions.",
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `srv`,
			Link:  `https://github.com/mitranim/srv`,
			Desc:  `Extremely simple Go tool that serves files out of a given folder, using a file resolution algorithm similar to Github Pages, Netlify, or the default Nginx config.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `gtg`,
			Link:  `https://github.com/mitranim/gtg`,
			Desc:  `Go task group / task graph. Good for CLI task orchestration. Replaces Make and Mage.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `gord`,
			Link:  `https://github.com/mitranim/gord`,
			Desc:  `Simple ordered sets for Go.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `goh`,
			Link:  `https://github.com/mitranim/goh`,
			Desc:  `Go HTTP handlers. Utility types that represent a not-yet-sent HTTP response as a value (status, header, body) with NO added abstractions or interfaces.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `frac`,
			Link:  `https://github.com/mitranim/frac`,
			Desc:  `Missing feature of Go stdlib: integers ↔︎ fractional numeric strings, without rounding errors or bignums. Arbitrary fraction precision and radix.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `rout`,
			Link:  `https://github.com/mitranim/rout`,
			Desc:  `Imperative router for Go HTTP servers. Procedural control flow with declarative syntax. Doesn't need middleware.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `sublime-sql`,
			Link:  `https://github.com/mitranim/sublime-sql`,
			Desc:  `Sublime Text syntax definitions for SQL, rebuilt with better semantics. Currently only Postgres dialect.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2021`,
		},
		{
			Name:  `sublime-lisp`,
			Link:  `https://github.com/mitranim/sublime-lisp`,
			Desc:  `Lisp support for Sublime Text. Supports multiple dialects. Immature work in progress.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2021`,
		},
		{
			Name:  `emptty`,
			Link:  `https://github.com/mitranim/emptty`,
			Desc:  `Clears the terminal, for real.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2021`,
		},
		{
			Name:  `sublime-forth`,
			Link:  `https://github.com/mitranim/sublime-forth`,
			Desc:  `Sublime Text syntax for the Forth programming language.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2021`,
		},
		{
			Name:  `sublime-rebol`,
			Link:  `https://github.com/mitranim/sublime-rebol`,
			Desc:  `Immature syntax for Rebol/Red in Sublime Text.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2020`,
		},
		{
			Name:  `jsonfmt`,
			Link:  `https://github.com/mitranim/jsonfmt`,
			Desc:  `Flexible JSON formatter. Supports comments, single-line until width limit, fixes punctuation. Library and optional CLI.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `jel`,
			Link:  `https://github.com/mitranim/jel`,
			Desc:  `"JSON Expession Language". Expresses a whitelisted subset of SQL with simple JSON structures.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `sqlb`,
			Link:  `https://github.com/mitranim/sqlb`,
			Desc:  `SQL query builder in Go. Oriented towards plain SQL. Supports composition.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `sqlp`,
			Link:  `https://github.com/mitranim/sqlp`,
			Desc:  `Parser for rewriting foreign code embedded in SQL.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `gos`,
			Link:  `https://github.com/mitranim/gos`,
			Desc:  `Tool for mapping between Go structs and plain SQL. Not an ORM.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `reqdec`,
			Link:  `https://github.com/mitranim/reqdec`,
			Desc:  `Tool for decoding HTTP requests into Go structs. Transparently supports multiple formats: JSON, URL-encoded, multipart.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `untext`,
			Link:  `https://github.com/mitranim/untext`,
			Desc:  `Missing feature of the Go standard library: unmarshal arbitrary string into arbitrary value.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `refut`,
			Link:  `https://github.com/mitranim/refut`,
			Desc:  "Reflect utils: utilities missing from the `reflect` package in the Go standard library.",
			Role:  `author`,
			Tech:  `Go`,
			Start: `2020`,
		},
		{
			Name:  `sublime-fmt`,
			Link:  `https://github.com/mitranim/sublime-fmt`,
			Desc:  "Sublime Text generic formatter plugin; formats arbitrary code by calling arbitrary executables, such as `gofmt`.",
			Role:  `author`,
			Tech:  `Python`,
			Start: `2020`,
		},
		{
			Name:  `Core Spirit`,
			Link:  `https://corespirit.com`,
			Desc:  `Current employer. Platform for practitioners of spiritual arts. Combines articles, services, and more.`,
			Role:  `tech lead`,
			Tech:  `Postgres, Go, JS`,
			Start: `2020`,
		},
		{
			Name:  `eth`,
			Link:  `https://github.com/purelabio/eth`,
			Desc:  "Client library for interacting with Ethereum from Go. Superior alternative to the `official` client libraries provided with `go-ethereum`.",
			Role:  `author`,
			Tech:  `Go`,
			Start: `2018`,
			End:   `2018`,
		},
		{
			Name:  `gow`,
			Link:  `https://github.com/mitranim/gow`,
			Desc:  "Missing watch mode for Go commands. Watch Go files and execute a command like `go run` or `go test`.",
			Role:  `author`,
			Tech:  `Go`,
			Start: `2018`,
		},
		{
			Name:  `repr`,
			Link:  `https://github.com/mitranim/repr`,
			Desc:  `Pretty-print Go data structures as valid Go code.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2018`,
		},
		{
			Name:  `sublime-goasm`,
			Link:  `https://github.com/mitranim/sublime-goasm`,
			Desc:  `Syntax definition for Go assembly for Sublime Text.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2018`,
		},
		{
			Name:  `sublime-caser`,
			Link:  `https://github.com/mitranim/sublime-caser`,
			Desc:  `Sublime Text plugin for converting typographic cases.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2018`,
		},
		{
			Name:  `sublime-gox`,
			Link:  `https://github.com/mitranim/sublime-gox`,
			Desc:  `[Merged into ST] Sublime Text syntax for Go.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2018`,
		},
		{
			Name:  `style-base`,
			Link:  `https://github.com/mitranim/style-base`,
			Desc:  `Minimal CSS foundation.`,
			Role:  `author`,
			Tech:  `CSS`,
			Start: `2018`,
		},
		{
			Name:  `papyre`,
			Link:  `https://github.com/mitranim/papyre`,
			Desc:  `Build tool for static websites. Bring your own rendering engine. Works well with React and Netlify CMS.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2018`,
			End:   `2018`,
		},
		{
			Name:  `sublime-themes`,
			Link:  `https://github.com/mitranim/sublime-themes`,
			Desc:  `Custom color schemes for Sublime Text.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2017`,
		},
		{
			Name:  `wordplay`,
			Link:  `https://github.com/mitranim/wordplay`,
			Desc:  `the PUNS 🔥`,
			Role:  `colab`,
			Tech:  `🧠`,
			Start: `2017`,
		},
		{
			Name:  `epdf`,
			Link:  `https://github.com/mitranim/epdf`,
			Desc:  `Render any URL to PDF using Electron.`,
			Role:  `author`,
			Tech:  `JS, Node`,
			Start: `2017`,
			End:   `2017`,
		},
		{
			Name:  `posterus`,
			Link:  `https://github.com/mitranim/posterus`,
			Desc:  `Asynchronous primitives. Superior replacement for JS promises. Synchronous by default. Supports true cancelation. Supports fibers.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2017`,
		},
		{
			Name:  `Bolala`,
			Link:  `https://bolala.ru`,
			Desc:  `E-commerce platform. Had some interesting tech but didn't launch. Permanently down.`,
			Role:  `frontend`,
			Tech:  `JS`,
			Start: `2017`,
			End:   `2018`,
		},
		{
			Name:  `clj-ws-client`,
			Link:  `https://github.com/mitranim/clj-ws-client`,
			Desc:  `WebSocket client (not server) written in pure Clojure with no dependencies. Less bad than most alternatives.`,
			Role:  `author`,
			Tech:  `Clojure`,
			Start: `2017`,
		},
		{
			Name:  `clojure-datomic-starter`,
			Link:  `https://github.com/mitranim/clojure-datomic-starter`,
			Desc:  `Quickstart/template for a Clojure/Ring webserver with Datomic.`,
			Role:  `author`,
			Tech:  `Clojure`,
			Start: `2017`,
		},
		{
			Name:  `clojure-auth0-starter`,
			Link:  `https://github.com/mitranim/clojure-auth0-starter`,
			Desc:  `Quickstart/template for a Clojure/Ring webserver with Auth0.`,
			Role:  `author`,
			Tech:  `Clojure`,
			Start: `2017`,
		},
		{
			Name:  `clojure-forge`,
			Link:  `https://github.com/mitranim/clojure-forge`,
			Desc:  `Development tool for Clojure. Especially useful for Ring servers. Watches files, reloads code, restarts system, displays system errors on a webpage.`,
			Role:  `author`,
			Tech:  `Clojure`,
			Start: `2017`,
		},
		{
			Name:  `sublime-clojure`,
			Link:  `https://github.com/mitranim/sublime-clojure`,
			Desc:  `[Merged into ST] Sublime Text syntax for Clojure.`,
			Role:  `author`,
			Tech:  `Sublime`,
			Start: `2017`,
		},
		{
			Name:  `Shanzhai City`,
			Link:  `https://shanzhaicity.com`,
			Desc:  `Various webapps and websites for Shanzhai City, a joint US-Chinese startup aiming at making charity effective.`,
			Role:  `tech lead`,
			Tech:  `Go, JS, Clojure`,
			Start: `2017`,
			End:   `2018`,
		},
		{
			Name:  `Render.js`,
			Link:  `https://renderjs.io`,
			Desc:  `Experimental service for prerendering JS SPA into HTML. An order of magnitude faster than the alternatives.`,
			Role:  `member`,
			Tech:  `JS, Node`,
			Start: `2016`,
			End:   `2017`,
		},
		{
			Name:  `fpx`,
			Link:  `https://github.com/mitranim/fpx`,
			Desc:  `Functional programming utils and type assertions for JS.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2016`,
		},
		{
			Name:  `emerge`,
			Link:  `https://github.com/mitranim/emerge`,
			Desc:  `Utils for using plain JS objects as immutable data structures with extremely memory-efficient updates. Heavily inspired by clojure.core. Much lighter and simpler than the popular alternatives.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2015`,
		},
		{
			Name:  `chat`,
			Link:  `https://github.com/mitranim/chat`,
			Desc:  `Realtime chat demo made with Firebase and React.`,
			Role:  `author`,
			Tech:  `JS, Firebase, React`,
			Start: `2015`,
		},
		{
			Name:  `statil`,
			Link:  `https://github.com/mitranim/statil`,
			Desc:  `Simple templating utility that uses JS for embedded scripting. Superseded by Prax.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2015`,
			End:   `2021`,
		},
		{
			Name:  `alder`,
			Link:  `https://github.com/mitranim/alder`,
			Desc:  `Experimental rendering library inspired by React and Reagent. Represents view components with plain functions and DOM with plain JavaScript data structures. Superseded by Prax.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2015`,
			End:   `2015`,
		},
		{
			Name:  `ToBox`,
			Link:  `https://tobox.com`,
			Desc:  `Stylish, visual platform for creating online shops. Part of the web frontend team. Permanently down.`,
			Role:  `member`,
			Tech:  `JS, React`,
			Start: `2015`,
			End:   `2016`,
		},
		{
			Name:  `atril`,
			Link:  `https://mitranim.com/atril/`,
			Desc:  `Experimental rendering library inspired by React and Angular. Documented with flashy demos, but unused. Superseded by Prax.`,
			Role:  `author`,
			Tech:  `TypeScript`,
			Start: `2015`,
			End:   `2015`,
		},
		{
			Name:  `GorodDeti`,
			Link:  `https://mitranim.com/kindergarten`,
			Desc:  `First version of a website about the kindergarten ran by a friend's friend. See [https://goroddeti.ru](https://goroddeti.ru) for the current version (not mine).`,
			Role:  `implementer`,
			Tech:  `JS, SCSS`,
			Start: `2015`,
			End:   `2015`,
		},
		{
			Name:  `simple-pjax`,
			Link:  `https://github.com/mitranim/simple-pjax`,
			Desc:  `Drop-in JS library that dramatically speeds up page transitions on server-rendered websites. See the explanatory [blog post](/posts/cheating-for-performance-pjax/).`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2015`,
		},
		{
			Name:  `xhttp`,
			Link:  `https://github.com/mitranim/xhttp`,
			Desc:  `Isomorphic library for making HTTP requests in Node.js and browsers. Superseded by abortable fetch.`,
			Role:  `author`,
			Tech:  `JS`,
			Start: `2014`,
			End:   `2021`,
		},
		{
			Name:  `stylific`,
			Link:  `https://github.com/mitranim/stylific`,
			Desc:  `CSS library/framework. Similar to [Bootstrap](https://getbootstrap.com), built on different principles. Accompanied by [stylific-lite](https://mitranim.com/stylific-lite/), succeeded by [style-base](https://github.com/mitranim/style-base).`,
			Role:  `author`,
			Tech:  `SCSS`,
			Start: `2015`,
		},
		{
			Name:  `codex`,
			Link:  `https://github.com/mitranim/codex`,
			Desc:  `Generator of random synthetic words or names.`,
			Role:  `author`,
			Tech:  `Go`,
			Start: `2015`,
		},
		{
			Name:  `ProstoPoi`,
			Link:  `https://prostopoi.ru`,
			Desc:  `Poi community website. We have our own video lessons, go check us out!`,
			Role:  `implementer`,
			Tech:  `Python, Django, React`,
			Start: `2014`,
		},
		{
			Name:  `jisp`,
			Link:  `https://mitranim.com/jisp/`,
			Desc:  `Lisp-style language that compiles to JavaScript. Currently on pause.`,
			Role:  `author`,
			Tech:  `JS, Jisp`,
			Start: `2014`,
			End:   `2015`,
		},
	}
}
