package main

type Work struct {
	Name string
	Link string
	Role string
	Tech string
	Desc string
	Lifecycle
}

func (self Work) Meta(bui B) {
	SemicolonList{
		F(self.Role),
		F(self.Tech),
		F(self.Range()),
		F(self.StatusLink),
	}.Render(bui)
}

type Lifecycle struct {
	Start  string
	End    string
	Status string
	Link   string
}

func (self Lifecycle) Range() string {
	if self.Start != `` && self.End != `` {
		return self.Start + EMDASH + self.End
	}
	if self.Start != `` && self.End == `` {
		return self.Start + `+`
	}
	return ``
}

func (self Lifecycle) StatusLink(bui B) {
	if self.Link != `` {
		bui.Child(LinkExt(self.Link, self.Status))
		return
	}

	if self.Status != `` {
		bui.Child(self.Status)
		return
	}
}

func (self Lifecycle) StatusEnd(bui B) {
	if !buiChild(bui, self.StatusLink) {
		bui.Child(self.End)
	}
}

func initWorks() []Work {
	return []Work{
		{
			Name:      `gg`,
			Link:      `https://github.com/mitranim/gg`,
			Desc:      `Essential tools missing from the Go standard library.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2022`},
		},
		{
			Name:      `js`,
			Link:      `https://github.com/mitranim/js`,
			Desc:      `Kinda "JS standard library" that doesn't suck. Also a tiny framework for JS apps.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2022`},
		},
		{
			Name: `ur`,
			Link: `https://github.com/mitranim/ur`,
			Desc: `Superior URL and query implementation for JS. Similar to built-in URL but actually usable.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2022`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name: `test`,
			Link: `https://github.com/mitranim/test`,
			Desc: `Superior testing and benchmarking library for JS. Runs in all environments. High benchmark accuracy.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name:      `oas`,
			Link:      `https://github.com/mitranim/oas`,
			Desc:      `OpenAPI specs for your Go server, generated at server runtime using reflection.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `gr`,
			Link:      `https://github.com/mitranim/gr`,
			Desc:      `Short for "Go Request-Response". Shortcuts for making HTTP requests and reading HTTP responses in Go.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `ded`,
			Link:      `https://github.com/mitranim/ded`,
			Desc:      `Experimental tool for deduplicating concurrent background operations in Go, with limited-time caching.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `gt`,
			Link:      `https://github.com/mitranim/gt`,
			Desc:      `Short for "Go Types". Important data types missing from the Go standard library.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `cc`,
			Link:      `https://github.com/mitranim/cc`,
			Desc:      `Tiny Go tool for running multiple functions concurrently and collecting their results into an error slice.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `cmd`,
			Link:      `https://github.com/mitranim/cmd`,
			Desc:      "Missing feature of the Go standard library: ability to define subcommands while using `flag`.",
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `sublime-scss`,
			Link:      `https://github.com/mitranim/sublime-scss`,
			Desc:      `Redesigned CSS and SCSS syntaxes for Sublime Text. Built on open-ended principles. Designed for forward compatibility.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `ProstoPoi SSG`,
			Link:      `https://github.com/mitranim/pp`,
			Desc:      `Poi community website. Runs since 2014. Now converted from Django (Python) to static generation (JS), open sourced.`,
			Role:      `implementer`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `gax`,
			Link:      `https://github.com/mitranim/gax`,
			Desc:      `Simple system for writing HTML as Go code. Use normal Go conditionals, loops and functions. Benefit from typing and code analysis. Better performance than templating.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name: `imperouter`,
			Link: `https://github.com/mitranim/imperouter`,
			Desc: `Simple imperative router for hybrid SSR+SPA.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name: `jol`,
			Link: `https://github.com/mitranim/jol`,
			Desc: `JS Collection Classes. Tiny extensions on JS built-in classes, with nice features such as easy-to-use typed collections, dictionary with structured keys, and more.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name:      `web-starter`,
			Link:      `https://github.com/mitranim/web-starter`,
			Desc:      `Starter templates for minimal web apps, from simplest to complex. Sucks less than X. Work in progress.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name: `afr`,
			Link: `https://github.com/mitranim/afr`,
			Desc: `Flexible library for serving files, with optional client integration for CSS injection and page reload.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name: `espo`,
			Link: `https://github.com/mitranim/espo`,
			Desc: `Observables via proxies. Particularly suited for UI programming. Supports automatic, implicit sub/resub and resource deinit.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name: `prax`,
			Link: `https://github.com/mitranim/prax`,
			Desc: `Simple rendering library for hybrid SSR+SPA. Superior replacement for rendering frameworks like React.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2021`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name:      `jtg`,
			Link:      `https://github.com/mitranim/jtg`,
			Desc:      `"JS Task Group". Simple JS-based replacement for Make, Gulp, etc.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `try`,
			Link:      `https://github.com/mitranim/try`,
			Desc:      "Shorter error handling in Go. Supports two styles: explicit `try` and exceptions.",
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `srv`,
			Link:      `https://github.com/mitranim/srv`,
			Desc:      `Extremely simple Go tool that serves files out of a given folder, using a file resolution algorithm similar to Github Pages, Netlify, or the default Nginx config.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `gtg`,
			Link:      `https://github.com/mitranim/gtg`,
			Desc:      `Go task group / task graph. Good for CLI task orchestration. Replaces Make and Mage.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `gord`,
			Link:      `https://github.com/mitranim/gord`,
			Desc:      `Simple ordered sets for Go.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `goh`,
			Link:      `https://github.com/mitranim/goh`,
			Desc:      `Go HTTP handlers. Utility types that represent a not-yet-sent HTTP response as a value (status, header, body) with NO added abstractions or interfaces.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `frac`,
			Link:      `https://github.com/mitranim/frac`,
			Desc:      `Missing feature of Go stdlib: integers ↔︎ fractional numeric strings, without rounding errors or bignums. Arbitrary fraction precision and radix.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `rout`,
			Link:      `https://github.com/mitranim/rout`,
			Desc:      `Imperative router for Go HTTP servers. Procedural control flow with declarative syntax. Doesn't need middleware.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `sublime-sql`,
			Link:      `https://github.com/mitranim/sublime-sql`,
			Desc:      `Sublime Text syntax definitions for SQL, rebuilt with better semantics. Currently only Postgres dialect.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `sublime-lisp`,
			Link:      `https://github.com/mitranim/sublime-lisp`,
			Desc:      `Lisp support for Sublime Text. Supports multiple dialects. Immature work in progress.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `emptty`,
			Link:      `https://github.com/mitranim/emptty`,
			Desc:      `Clears the terminal, for real.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `sublime-forth`,
			Link:      `https://github.com/mitranim/sublime-forth`,
			Desc:      `Sublime Text syntax for the Forth programming language.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2021`},
		},
		{
			Name:      `sublime-rebol`,
			Link:      `https://github.com/mitranim/sublime-rebol`,
			Desc:      `Immature syntax for Rebol/Red in Sublime Text.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `jsonfmt`,
			Link:      `https://github.com/mitranim/jsonfmt`,
			Desc:      `Flexible JSON formatter. Supports comments, single-line until width limit, fixes punctuation. Library and optional CLI.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name: `jel`,
			Link: `https://github.com/mitranim/jel`,
			Desc: `"JSON Expession Language". Expresses a whitelisted subset of SQL with simple JSON structures.`,
			Role: `author`,
			Tech: `Go`,
			Lifecycle: Lifecycle{
				Start:  `2020`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/sqlb`,
			},
		},
		{
			Name:      `sqlb`,
			Link:      `https://github.com/mitranim/sqlb`,
			Desc:      `SQL query builder in Go. Highly flexible and efficient.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `sqlp`,
			Link:      `https://github.com/mitranim/sqlp`,
			Desc:      `Parser for rewriting foreign code embedded in SQL.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `gos`,
			Link:      `https://github.com/mitranim/gos`,
			Desc:      `Tool for mapping between Go structs and plain SQL. Not an ORM.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `rd`,
			Link:      `https://github.com/mitranim/rd`,
			Desc:      `Tool for decoding HTTP requests into Go structs. Transparently supports multiple formats: JSON, URL-encoded, multipart.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `untext`,
			Link:      `https://github.com/mitranim/untext`,
			Desc:      `Missing feature of the Go standard library: unmarshal arbitrary string into arbitrary value.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `rf`,
			Link:      `https://github.com/mitranim/rf`,
			Desc:      "Important utilities missing from the `reflect` package in the Go standard library.",
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `sublime-fmt`,
			Link:      `https://github.com/mitranim/sublime-fmt`,
			Desc:      "Sublime Text generic formatter plugin; formats arbitrary code by calling arbitrary executables, such as `gofmt`.",
			Role:      `author`,
			Tech:      `Python`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `Core Spirit`,
			Link:      `https://corespirit.com`,
			Desc:      `Current employer. Platform for practitioners of spiritual arts. Combines articles, services, and more.`,
			Role:      `tech lead`,
			Tech:      `Postgres, Go, JS`,
			Lifecycle: Lifecycle{Start: `2020`},
		},
		{
			Name:      `eth`,
			Link:      `https://github.com/purelabio/eth`,
			Desc:      "Client library for interacting with Ethereum from Go. Superior alternative to the \"official\" client libraries provided with `go-ethereum`.",
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2018`, End: `2018`},
		},
		{
			Name:      `gow`,
			Link:      `https://github.com/mitranim/gow`,
			Desc:      "Missing watch mode for Go commands. Watch Go files and execute a command like `go run` or `go test`.",
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `repr`,
			Link:      `https://github.com/mitranim/repr`,
			Desc:      `Pretty-print Go data structures as valid Go code.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `sublime-goasm`,
			Link:      `https://github.com/mitranim/sublime-goasm`,
			Desc:      `Syntax definition for Go assembly for Sublime Text.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `sublime-caser`,
			Link:      `https://github.com/mitranim/sublime-caser`,
			Desc:      `Sublime Text plugin for converting typographic cases.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `sublime-gox`,
			Link:      `https://github.com/mitranim/sublime-gox`,
			Desc:      `[Merged into ST] Sublime Text syntax for Go.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `sb`,
			Link:      `https://github.com/mitranim/sb`,
			Desc:      `Minimal CSS foundation.`,
			Role:      `author`,
			Tech:      `CSS`,
			Lifecycle: Lifecycle{Start: `2018`},
		},
		{
			Name:      `papyre`,
			Link:      `https://github.com/mitranim/papyre`,
			Desc:      `Build tool for static websites. Bring your own rendering engine. Works well with React and Netlify CMS.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2018`, End: `2018`},
		},
		{
			Name:      `sublime-themes`,
			Link:      `https://github.com/mitranim/sublime-themes`,
			Desc:      `Custom color schemes for Sublime Text.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2017`},
		},
		{
			Name:      `wordplay`,
			Link:      `https://github.com/mitranim/wordplay`,
			Desc:      `the PUNS 🔥`,
			Role:      `colab`,
			Tech:      `🧠`,
			Lifecycle: Lifecycle{Start: `2017`},
		},
		{
			Name:      `epdf`,
			Link:      `https://github.com/mitranim/epdf`,
			Desc:      `Render any URL to PDF using Electron.`,
			Role:      `author`,
			Tech:      `JS, Node`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2017`},
		},
		{
			Name:      `posterus`,
			Link:      `https://github.com/mitranim/posterus`,
			Desc:      `Asynchronous primitives. Superior replacement for JS promises. Synchronous by default. Supports true cancelation. Supports fibers.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2017`, Status: `paused`},
		},
		{
			Name:      `Bolala`,
			Link:      `https://bolala.ru`,
			Desc:      `E-commerce platform. Had some interesting tech but didn't launch. Permanently down.`,
			Role:      `frontend`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2018`},
		},
		{
			Name:      `clj-ws-client`,
			Link:      `https://github.com/mitranim/clj-ws-client`,
			Desc:      `WebSocket client (not server) written in pure Clojure with no dependencies. Less bad than most alternatives.`,
			Role:      `author`,
			Tech:      `Clojure`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2017`},
		},
		{
			Name:      `clojure-datomic-starter`,
			Link:      `https://github.com/mitranim/clojure-datomic-starter`,
			Desc:      `Quickstart/template for a Clojure/Ring webserver with Datomic.`,
			Role:      `author`,
			Tech:      `Clojure`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2017`},
		},
		{
			Name:      `clojure-auth0-starter`,
			Link:      `https://github.com/mitranim/clojure-auth0-starter`,
			Desc:      `Quickstart/template for a Clojure/Ring webserver with Auth0.`,
			Role:      `author`,
			Tech:      `Clojure`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2017`},
		},
		{
			Name:      `clojure-forge`,
			Link:      `https://github.com/mitranim/clojure-forge`,
			Desc:      `Development tool for Clojure. Especially useful for Ring servers. Watches files, reloads code, restarts system, displays system errors on a webpage.`,
			Role:      `author`,
			Tech:      `Clojure`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2018`},
		},
		{
			Name:      `sublime-clojure`,
			Link:      `https://github.com/mitranim/sublime-clojure`,
			Desc:      `[Merged into ST] Sublime Text syntax for Clojure.`,
			Role:      `author`,
			Tech:      `Sublime`,
			Lifecycle: Lifecycle{Start: `2017`},
		},
		{
			Name:      `Shanzhai City`,
			Link:      `https://shanzhaicity.com`,
			Desc:      `Various webapps and websites for Shanzhai City, a joint US-Chinese startup aiming at making charity effective.`,
			Role:      `tech lead`,
			Tech:      `Go, JS, Clojure`,
			Lifecycle: Lifecycle{Start: `2017`, End: `2018`},
		},
		{
			Name:      `Render.js`,
			Link:      `https://renderjs.io`,
			Desc:      `Experimental service for prerendering JS SPA into HTML. An order of magnitude faster than the alternatives.`,
			Role:      `member`,
			Tech:      `JS, Node`,
			Lifecycle: Lifecycle{Start: `2016`, End: `2017`},
		},
		{
			Name: `fpx`,
			Link: `https://github.com/mitranim/fpx`,
			Desc: `Functional programming utils and type assertions for JS.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2016`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name:      `emerge`,
			Link:      `https://github.com/mitranim/emerge`,
			Desc:      `Utils for using plain JS objects as immutable data structures with extremely memory-efficient updates. Heavily inspired by clojure.core. Much lighter and simpler than the popular alternatives.`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2015`, Status: `paused`},
		},
		{
			Name:      `chat`,
			Link:      `https://github.com/mitranim/chat`,
			Desc:      `Realtime chat demo made with Firebase and React.`,
			Role:      `author`,
			Tech:      `JS, Firebase, React`,
			Lifecycle: Lifecycle{Start: `2015`},
		},
		{
			Name: `statil`,
			Link: `https://github.com/mitranim/statil`,
			Desc: `Simple templating utility that uses JS for embedded scripting. Superseded by Prax.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2015`,
				End:    `2021`,
				Status: `replaced`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name: `alder`,
			Link: `https://github.com/mitranim/alder`,
			Desc: `Experimental rendering library inspired by React and Reagent. Represents view components with plain functions and DOM with plain JavaScript data structures. Superseded by Prax.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2015`,
				End:    `2015`,
				Status: `replaced`,
				Link:   `https://github.com/mitranim/prax`,
			},
		},
		{
			Name:      `ToBox`,
			Link:      `https://tobox.com`,
			Desc:      `Stylish, visual platform for creating online shops. Part of the web frontend team. Permanently down.`,
			Role:      `member`,
			Tech:      `JS, React`,
			Lifecycle: Lifecycle{Start: `2015`, End: `2016`},
		},
		{
			Name: `atril`,
			Link: `https://mitranim.com/atril/`,
			Desc: `Experimental rendering library inspired by React and Angular. Documented with flashy demos, but unused. Superseded by Prax.`,
			Role: `author`,
			Tech: `TypeScript`,
			Lifecycle: Lifecycle{
				Start:  `2015`,
				End:    `2015`,
				Status: `replaced`,
				Link:   `https://github.com/mitranim/prax`,
			},
		},
		{
			Name:      `GorodDeti`,
			Link:      `https://mitranim.com/kindergarten`,
			Desc:      `First version of a website about the kindergarten ran by a friend's friend. See [https://goroddeti.ru](https://goroddeti.ru) for the current version (not mine).`,
			Role:      `implementer`,
			Tech:      `JS, SCSS`,
			Lifecycle: Lifecycle{Start: `2015`, End: `2015`},
		},
		{
			Name:      `simple-pjax`,
			Link:      `https://github.com/mitranim/simple-pjax`,
			Desc:      `Drop-in JS library that dramatically speeds up page transitions on server-rendered websites. See the explanatory [blog post](/posts/cheating-for-performance-pjax/).`,
			Role:      `author`,
			Tech:      `JS`,
			Lifecycle: Lifecycle{Start: `2015`},
		},
		{
			Name: `xhttp`,
			Link: `https://github.com/mitranim/xhttp`,
			Desc: `Shortcuts for the native JS fetch/Request/Response API. Provides a fluent builder-style API for request building and response reading.`,
			Role: `author`,
			Tech: `JS`,
			Lifecycle: Lifecycle{
				Start:  `2014`,
				Status: `subsumed`,
				Link:   `https://github.com/mitranim/js`,
			},
		},
		{
			Name:      `stylific`,
			Link:      `https://github.com/mitranim/stylific`,
			Desc:      `CSS library/framework. Similar to [Bootstrap](https://getbootstrap.com), built on different principles. Accompanied by [stylific-lite](https://mitranim.com/stylific-lite/), succeeded by [sb](https://github.com/mitranim/sb).`,
			Role:      `author`,
			Tech:      `SCSS`,
			Lifecycle: Lifecycle{Start: `2015`, Status: `paused`},
		},
		{
			Name:      `codex`,
			Link:      `https://github.com/mitranim/codex`,
			Desc:      `Generator of random synthetic words or names.`,
			Role:      `author`,
			Tech:      `Go`,
			Lifecycle: Lifecycle{Start: `2015`, Status: `paused`},
		},
		{
			Name:      `ProstoPoi`,
			Link:      `https://prostopoi.ru`,
			Desc:      `Poi community website. We have our own video lessons, go check us out!`,
			Role:      `implementer`,
			Tech:      `Python, Django, React`,
			Lifecycle: Lifecycle{Start: `2014`},
		},
		{
			Name:      `jisp`,
			Link:      `https://mitranim.com/jisp/`,
			Desc:      `Lisp-style language that compiles to JavaScript.`,
			Role:      `author`,
			Tech:      `JS, Jisp`,
			Lifecycle: Lifecycle{Start: `2014`, Status: `update upcoming`},
		},
	}
}
