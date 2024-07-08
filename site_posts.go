package main

import (
	"strings"

	"github.com/mitranim/gg"
)

// TODO: add `.WrittenAt`, which often doesn't match `.PublishedAt`.
type PagePost struct {
	Page
	RedirFrom   []string
	PublishedAt Time
	UpdatedAt   Time
	IsListed    bool
}

func (self PagePost) ExistsAsFile() bool {
	return !self.PublishedAt.IsZero() || !FLAGS.PROD
}

func (self PagePost) ExistsInFeeds() bool {
	return self.ExistsAsFile() && bool(self.IsListed)
}

// Somewhat inefficient but shouldn't be measurable.
func (self PagePost) TimeString() string {
	var out []string

	if !self.PublishedAt.IsZero() {
		out = append(out, `published `+timeFmtHuman(self.PublishedAt))
		if !self.UpdatedAt.IsZero() {
			out = append(out, `updated `+timeFmtHuman(self.UpdatedAt))
		}
	}

	return strings.Join(out, `, `)
}

func (self PagePost) Make(site Site) {
	PageWrite(self, self.Render(site))

	for _, path := range self.RedirFrom {
		writePublic(path, F(
			E(`meta`, AP(`http-equiv`, `refresh`, `content`, `0;URL='`+self.GetLink()+`'`)),
		))
	}
}

func (self PagePost) MakeMd() []byte {
	if self.MdHtml == nil {
		self.MdHtml = self.Md(self, MdOpt{})
	}
	return self.MdHtml
}

func (self PagePost) FeedItem() FeedItem {
	href := siteBaseUrl.Get().WithPath(self.GetLink()).String()

	return FeedItem{
		XmlBase:     href,
		Title:       self.Page.Title,
		Link:        &FeedLink{Href: href},
		Author:      FEED_AUTHOR,
		Description: self.Page.Description,
		Id:          href,
		Published:   self.PublishedAt.MaybeTime(),
		Updated:     self.GetUpdatedAt().MaybeTime(),
		Content:     FeedPost(self).String(),
	}
}

func (self PagePost) GetIsListed() bool { return self.IsListed }

func (self PagePost) GetUpdatedAt() Time {
	return gg.Or(self.UpdatedAt, self.PublishedAt, timeNow())
}

func initSitePosts() []PagePost {
	return []PagePost{
		PagePost{
			Page: Page{
				Path:        `posts/hades.html`,
				Title:       `Hades: tweak recommendations`,
				Description: `Suggestions for how to play Hades, an excellent single-player roguelike game.`,
				MdTpl:       readTemplate(`posts/hades.md`),
			},
			PublishedAt: timeParse(`2023-08-25T15:42:29Z`),
			UpdatedAt:   timeParse(`2024-02-16T13:22:47Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/speed.html`,
				Title:       `Using speedhacks in single player games`,
				Description: `Explanation and instructions on speedhacking, a surprisingly handy tool in gaming.`,
				MdTpl:       readTemplate(`posts/speed.md`),
			},
			PublishedAt: timeParse(`2023-08-25T14:00:44Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/witcher.html`,
				Title:       `Witcher franchise: how to enjoy`,
				Description: `Essential tips and tricks for Witcher games. Spoiler-free!`,
				MdTpl:       readTemplate(`posts/witcher.md`),
			},
			PublishedAt: timeParse(`2023-03-20T23:40:42Z`),
			UpdatedAt:   timeParse(`2023-08-25T14:00:44Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/divinity-original-sin-2.html`,
				Title:       `Divinity Original Sin 2: how to play and enjoy`,
				Description: `Mod recommendations and gameplay suggestions. Spoiler-free!`,
				MdTpl:       readTemplate(`posts/divinity-original-sin-2.md`),
			},
			PublishedAt: timeParse(`2023-03-17T12:01:03Z`),
			UpdatedAt:   timeParse(`2023-08-25T14:00:44Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:  `posts/steins-gate.html`,
				Title: `[Draft] Impressions: Steins Gate series (games and anime).`,
				MdTpl: readTemplate(`posts/steins-gate.md`),
			},
		},
		PagePost{
			Page: Page{
				Path:        `posts/anime-impressions-parasyte.html`,
				Title:       `Anime impressions: Parasyte`,
				Description: `Thoughts and analysis on this surprisingly deep anime. Spoilers!`,
				MdTpl:       readTemplate(`posts/anime-impressions-parasyte.md`),
			},
			PublishedAt: timeParse(`2022-03-08T07:02:11Z`),
			UpdatedAt:   timeParse(`2022-09-05T11:40:59Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/anime-impressions-evangelion.html`,
				Title:       `Anime impressions: Evangelion`,
				Description: `How to watch: Neon Genesis Evangelion, End of Evangelion.`,
				MdTpl:       readTemplate(`posts/anime-impressions-evangelion.md`),
			},
			PublishedAt: timeParse(`2022-03-08T06:31:41Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/anime.html`,
				Title:       `Anime impressions and recommendations`,
				Description: `Periodically-updated gist. Check later for more.`,
				MdTpl:       readTemplate(`posts/anime.md`),
			},
			PublishedAt: timeParse(`2022-03-08T05:48:55Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/andromeda.html`,
				Title:       `Game impressions: Mass Effect Andromeda`,
				Description: `Enjoyed, highly recommended.`,
				MdTpl:       readTemplate(`posts/andromeda.md`),
			},
			PublishedAt: timeParse(`2022-01-23T07:43:31Z`),
			UpdatedAt:   timeParse(`2022-06-19T11:03:04Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/goex.html`,
				Title:       `Shorten your Go code by using exceptions`,
				Description: `Go secretly favors exceptions. Use them.`,
				MdTpl:       readTemplate(`posts/goex.md`),
			},
			PublishedAt: timeParse(`2021-11-20T11:47:36Z`),
			UpdatedAt:   timeParse(`2023-10-31T11:55:26Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/spaces-tabs.html`,
				Title:       `Always spaces, never tabs`,
				Description: `Objective arguments that decided my personal preference.`,
				MdTpl:       readTemplate(`posts/spaces-tabs.md`),
			},
			PublishedAt: timeParse(`2020-10-23T06:48:15Z`),
			UpdatedAt:   timeParse(`2024-02-16T13:23:19Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lisp-sexpr-hacks.html`,
				Title:       `Hacks around S-expressions in Lisps`,
				Description: `How far people are willing to go to get prefix and infix in a Lisp syntax.`,
				MdTpl:       readTemplate(`posts/lisp-sexpr-hacks.md`),
			},
			PublishedAt: timeParse(`2020-10-21T06:34:24Z`),
			UpdatedAt:   timeParse(`2021-08-20T07:16:38Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-var-minus.html`,
				Title:       `Language design: gotchas with variadic minus`,
				Description: `Treating the minus operator as a function can be tricky and dangerous.`,
				MdTpl:       readTemplate(`posts/lang-var-minus.md`),
			},
			PublishedAt: timeParse(`2020-10-17T07:20:06Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-case-conventions.html`,
				Title:       `Language design: case conventions`,
				Description: `Objective arguments to solve case conventions and move on.`,
				MdTpl:       readTemplate(`posts/lang-case-conventions.md`),
			},
			PublishedAt: timeParse(`2020-10-16T15:30:41Z`),
			UpdatedAt:   timeParse(`2023-03-17T11:58:53Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/lang-homoiconic.html`,
				Title:       `Language design: homoiconicity`,
				Description: `Thoughts on homoiconicity, an interesting language quality seen in Lisps.`,
				MdTpl:       readTemplate(`posts/lang-homoiconic.md`),
			},
			PublishedAt: timeParse(`2020-10-16T12:41:58Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/warframe-headcanon.html`,
				Title:       `Warframe headcanon (spoilers)`,
				Description: `Collection of Warframe headcanon co-authored with friends.`,
				MdTpl:       readTemplate(`posts/warframe-headcanon.md`),
			},
			PublishedAt: timeParse(`2020-10-10T12:25:32Z`),
			UpdatedAt:   timeParse(`2023-04-11T15:43:24Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:  `posts/thoughts-on-the-egg.html`,
				Title: `Thoughts on The Egg: a short story by Andy Weir, animated by Kurzgesagt`,
				MdTpl: readTemplate(`posts/thoughts-on-the-egg.md`),
			},
			PublishedAt: timeParse(`2020-04-30T08:25:16Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/gameplay-conjecture.html`,
				Title:       `Gameplay conjecture`,
				Description: `Amount of gameplay ≈ amount of required decisions.`,
				MdTpl:       readTemplate(`posts/gameplay-conjecture.md`),
			},
			IsListed: !FLAGS.PROD,
		},
		PagePost{
			Page: Page{
				Path:        `posts/tips-and-tricks-doom-2016.html`,
				Title:       `Tips and tricks: Doom 2016`,
				Description: `General tips, notes on difficulty, enemies, runes, weapons.`,
				MdTpl:       readTemplate(`posts/tips-and-tricks-doom-2016.md`),
			},
			PublishedAt: timeParse(`2019-04-25T12:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/game-impressions-doom-2016.html`,
				Title:       `Game impressions: Doom 2016`,
				Description: `I really like Doom 2016, here's why.`,
				MdTpl:       readTemplate(`posts/game-impressions-doom-2016.md`),
			},
			PublishedAt: timeParse(`2019-04-25T11:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/astrotips.html`,
				Title:       `Announcing Astrotips: video guides on Astroneer`,
				Description: `A series of video guides, tips and tricks on Astroneer, an amazing space exploration and building game.`,
				MdTpl:       readTemplate(`posts/astrotips.md`),
			},
			PublishedAt: timeParse(`2019-02-22T11:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/camel-case-abbr.html`,
				Title:       `Don't abbreviate in camelCase`,
				Description: "CamelCase identifiers should avoid abbreviations, e.g. `JsonText` rather than `JSONText`.",
				MdTpl:       readTemplate(`posts/camel-case-abbr.md`),
			},
			PublishedAt: timeParse(`2019-01-17T07:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/remove-from-go.html`,
				Title:       `Things I would remove from Go`,
				Description: `If less is more, Go could gain by losing weight.`,
				MdTpl:       readTemplate(`posts/remove-from-go.md`),
			},
			PublishedAt: timeParse(`2019-01-15T01:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/back-from-hiatus-2019.html`,
				Title:       `Back from hiatus (2019)`,
				Description: `Back to blogging after three and a half years.`,
				MdTpl:       readTemplate(`posts/back-from-hiatus-2019.md`),
			},
			PublishedAt: timeParse(`2019-01-15T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/cheating-for-performance-pjax.html`,
				Title:       `Cheating for performance with pjax`,
				Description: `Faster page transitions, for free.`,
				MdTpl:       readTemplate(`posts/cheating-for-performance-pjax.md`),
			},
			RedirFrom:   []string{`thoughts/cheating-for-performance-pjax.html`},
			PublishedAt: timeParse(`2015-07-25T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/cheating-for-website-performance.html`,
				Title:       `Cheating for website performance`,
				Description: `Frontend tips for speeding up websites.`,
				MdTpl:       readTemplate(`posts/cheating-for-website-performance.md`),
			},
			RedirFrom:   []string{`thoughts/cheating-for-website-performance.html`},
			PublishedAt: timeParse(`2015-03-11T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/keeping-things-simple.html`,
				Title:       `Keeping things simple`,
				Description: `Musings on simplicity in programming.`,
				MdTpl:       readTemplate(`posts/keeping-things-simple.md`),
			},
			RedirFrom:   []string{`thoughts/keeping-things-simple.html`},
			PublishedAt: timeParse(`2015-03-10T00:00:00Z`),
			IsListed:    true,
		},
		PagePost{
			Page: Page{
				Path:        `posts/next-generation-today.html`,
				Title:       `Next generation today`,
				Description: `EcmaScript 2015/2016 workflow with current web frameworks.`,
				MdTpl:       readTemplate(`posts/next-generation-today.md`),
			},
			RedirFrom:   []string{`thoughts/next-generation-today.html`},
			PublishedAt: timeParse(`2015-05-18T00:00:00Z`),
			IsListed:    false,
		},
		PagePost{
			Page: Page{
				Path:        `posts/old-posts.html`,
				Title:       `Old posts`,
				Description: `Some old stuff from around the net.`,
				MdTpl:       readTemplate(`posts/old-posts.md`),
			},
			RedirFrom:   []string{`thoughts/old-posts.html`},
			PublishedAt: timeParse(`2015-01-01T00:00:00Z`),
			IsListed:    true,
		},
	}
}
