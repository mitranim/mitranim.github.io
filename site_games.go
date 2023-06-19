package main

import (
	x "github.com/mitranim/gax"
	"github.com/mitranim/gg"
)

type Games []Game

func (self Games) Listed() Games { return gg.Filter(self, Game.GetIsListed) }

func (self Games) HasListed() bool { return gg.Some(self, Game.GetIsListed) }

func (self Games) TimeSinks() []TimeSink {
	return gg.Sorted(gg.Compact(gg.MapUniq(self, Game.GetTimeSink)))
}

func (self Games) Tags() []Tag {
	return gg.Sorted(gg.Compact(gg.MapFlatUniq(self, Game.GetTags)))
}

type Game struct {
	Name     string
	Link     Url
	Img      Url
	Desc     string
	TimeSink TimeSink
	Tags     []Tag
	IsListed bool
}

func (self Game) GetIsListed() bool { return self.IsListed }

func (self Game) RenderName() any {
	link := self.Link.String()
	if link == `` {
		return self.Name
	}
	return LinkExt(link, self.Name)
}

func (self Game) GetTags() []Tag        { return self.Tags }
func (self Game) GetTimeSink() TimeSink { return self.TimeSink }

var TimeSinkOrd Ord

var (
	TimeSinkSafe     = TimeSink{`safe`, TimeSinkOrd.Next()}
	TimeSinkModerate = TimeSink{`moderate`, TimeSinkOrd.Next()}
	TimeSinkDanger   = TimeSink{`danger`, TimeSinkOrd.Next()}
	TimeSinkExtreme  = TimeSink{`extreme`, TimeSinkOrd.Next()}
	TimeSinkSuicide  = TimeSink{`suicide`, TimeSinkOrd.Next()}
)

type TimeSink struct {
	Name string
	Ord  uint64
}

func (self TimeSink) String() string         { return self.Name }
func (self TimeSink) Less(val TimeSink) bool { return self.Ord < val.Ord }

// TODO more colors.
func (self TimeSink) Mode() string {
	switch self {
	case TimeSinkSafe, TimeSinkModerate:
		return `--safe`
	case TimeSinkDanger, TimeSinkExtreme, TimeSinkSuicide:
		return `--danger`
	default:
		panic(errUnrecognized(self))
	}
}

// Implement `gax.Ren`.
func (self TimeSink) Render(bui B) {
	bui.E(
		`button`,
		AP(
			`is`, `btn-time-sink`,
			`type`, `button`,
			`class`, gg.SpacedOpt(`time-sink`, self.Mode()),
		),
		self.String(),
	)
}

var TagOrd Ord

var (
	TagPc      = Tag{`pc`, TagOrd.Next()}
	TagConsole = Tag{`console`, TagOrd.Next()}
	TagMobile  = Tag{`mobile`, TagOrd.Next()}

	TagWindows = Tag{`windows`, TagOrd.Next()}
	TagMac     = Tag{`mac`, TagOrd.Next()}
	TagAndroid = Tag{`android`, TagOrd.Next()}

	Tag2d = Tag{`2d`, TagOrd.Next()}
	Tag3d = Tag{`3d`, TagOrd.Next()}

	TagFirstPerson = Tag{`first_person`, TagOrd.Next()}
	TagThirdPerson = Tag{`third_person`, TagOrd.Next()}
	TagIsometric   = Tag{`isometric`, TagOrd.Next()}
	TagTopdown     = Tag{`topdown`, TagOrd.Next()}
	TagSideway     = Tag{`sideway`, TagOrd.Next()}

	TagSolo = Tag{`solo`, TagOrd.Next()}
	TagCoop = Tag{`coop`, TagOrd.Next()}
	TagPvp  = Tag{`pvp`, TagOrd.Next()}
	TagPve  = Tag{`pve`, TagOrd.Next()}

	TagStrategy  = Tag{`strategy`, TagOrd.Next()}
	TagPuzzle    = Tag{`puzzle`, TagOrd.Next()}
	TagShooter   = Tag{`shooter`, TagOrd.Next()}
	TagMelee     = Tag{`melee`, TagOrd.Next()}
	TagMagic     = Tag{`magic`, TagOrd.Next()}
	TagTurnBased = Tag{`turn_based`, TagOrd.Next()}
	TagSciFi     = Tag{`sci_fi`, TagOrd.Next()}
	TagFantasy   = Tag{`fantasy`, TagOrd.Next()}
	TagParty     = Tag{`party`, TagOrd.Next()}

	TagExploration = Tag{`exploration`, TagOrd.Next()}
	TagCrafting    = Tag{`crafting`, TagOrd.Next()}
	TagBuilding    = Tag{`building`, TagOrd.Next()}
	TagSpace       = Tag{`space`, TagOrd.Next()}
	TagOpenWorld   = Tag{`open_world`, TagOrd.Next()}
	TagUnderwater  = Tag{`underwater`, TagOrd.Next()}
	TagRolePlay    = Tag{`role_play`, TagOrd.Next()}
	TagStory       = Tag{`story`, TagOrd.Next()}
	TagDeckBuilder = Tag{`deck_builder`, TagOrd.Next()}
	TagRoguelike   = Tag{`roguelike`, TagOrd.Next()}
	TagGrind       = Tag{`grind`, TagOrd.Next()}
	TagStealth     = Tag{`stealth`, TagOrd.Next()}

	TagPacifist      = Tag{`pacifist`, TagOrd.Next()}
	TagPhilosophical = Tag{`philosophical`, TagOrd.Next()}
)

type Tag TimeSink

func (self Tag) String() string    { return TimeSink(self).String() }
func (self Tag) Less(val Tag) bool { return TimeSink(self).Less(TimeSink(val)) }

// Implement `gax.Ren`.
// Placeholder. Must be interactive.
func (self Tag) Render(bui B) {
	if self.Name == `` {
		return
	}
	bui.E(
		`button`,
		AP(`is`, `btn-tag`, `type`, `button`, `class`, `tag`),
		self.String(),
	)
}

type PageGames struct{ Page }

func (self PageGames) Make(site Site) {
	PageWrite(self, HtmlCommon(
		self,
		E(`div`, AttrsMain().Set(`class`, `article`),
			self.Head(),
			self.Content(site),
		),
	))
}

func (self PageGames) Head() x.Ren {
	return F(
		E(`h1`, nil, `Game Recommendations (work in progress)`),
		MdToHtmlStr(`

This list is carefully selected. There are many other games I've greatly
enjoyed, which I would not recommend right now, either because they're too
outdated (e.g. Etherlords 2), or because the online community that made them
great no longer exists (e.g. WoW).

Even if you prefer MacOS or Linux for general use, you should use a dedicated
Windows system for games. Many games don't exist on other platforms, or take
years to release a port, usually with compatibility issues and poor
performance. Many games have essential mods only available on Windows. Windows
also allows a much better selection of hardware.

Always, _always_ check [PC Gaming Wiki](https://pcgamingwiki.com) for essential
tweaks and mods for any given game.

`),
	)
}

func (self PageGames) Content(site Site) x.Ren {
	src := site.Games.Listed()
	inner := self.Grid(src)

	if gg.IsEmpty(src) {
		return inner
	}

	return F(
		NoscriptInteractivity().AttrAdd(`class`, `mar-bot-1`),
		self.TimeSinks(src),
		self.Tags(src),
		inner,
		Script(`scripts/games.mjs`),
	)
}

func (self PageGames) TimeSinks(src Games) x.Ren {
	vals := src.TimeSinks()
	if gg.IsEmpty(vals) {
		return nil
	}

	return E(
		`tag-likes`,
		AP(`class`, `tag-likes`, `data-role`, `filter`),
		// TODO clicking this should clear the filter.
		E(`span`, AP(`class`, `help`, `aria-label`, `combined by logical "or"`), `Time sinks:`),
		vals,
	)
}

func (self PageGames) Tags(src Games) x.Ren {
	vals := src.Tags()
	if gg.IsEmpty(vals) {
		return nil
	}

	return E(
		`tag-likes`,
		AP(`class`, `tag-likes`, `data-role`, `filter`),
		// TODO clicking this should clear the filter.
		E(`span`, AP(`class`, `help`, `aria-label`, `combined by logical "and"`), `Tags:`),
		vals,
	)
}

func (self PageGames) Grid(src Games) x.Ren {
	return E(`filter-list`, AP(`class`, `game-grid`),
		self.Placeholder(AttrsHidden(gg.IsNotEmpty(src))...),
		func(bui B) {
			for _, game := range src {
				bui.E(`filter-item`, AP(`class`, `game-grid-item`),
					func(bui B) {
						bui.E(`img`, AP(`src`, game.Img.String()))
						bui.E(`h3`, nil, game.RenderName())
						if gg.IsNotZero(game.Desc) {
							bui.E(`div`, nil, MdToHtmlStr(game.Desc))
						}
						if gg.IsNotZero(game.TimeSink) || gg.IsNotEmpty(game.Tags) {
							bui.E(`div`, AP(`class`, `tag-likes`),
								game.TimeSink,
								game.Tags,
							)
						}
					},
				)
			}
		},
	)
}

func (self PageGames) Placeholder(src ...x.Attr) x.Ren {
	return E(
		`p`,
		AP(`is`, `filter-placeholder`, `class`, `filter-placeholder`).A(src...),
		`Oops! It appears there are no game recommendations yet.`,
	)
}
