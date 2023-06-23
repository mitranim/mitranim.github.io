package main

import (
	"github.com/mitranim/gg"
)

type GameColl struct{ gg.Coll[string, Game] }

func (self GameColl) SortedByKeys(src ...string) GameColl {
	self.ValidateOrd(src...)
	self.Coll = gg.CollOf[string, Game](gg.Map(src, self.GetReq)...)
	return self
}

func (self GameColl) ValidateOrd(src ...string) {
	keys := gg.Reject(self.Pks(), gg.SetOf(src...).Has)

	if gg.IsNotEmpty(keys) {
		panic(gg.Errf(`entity keys missing from ordering: %q`, keys))
	}
}

func (self GameColl) Listed() Games {
	return gg.Filter(self.Slice, Game.GetIsListed)
}

// TODO consider moving to `gg.Coll`.
func (self GameColl) Pks() []string {
	return gg.Map(self.Slice, gg.ValidPk[string, Game])
}

type Games []Game

func (self Games) TimeSinks() []TimeSink {
	return gg.Sorted(gg.Compact(gg.MapUniq(self, Game.GetTimeSink)))
}

func (self Games) Tags() []Tag {
	return gg.Sorted(gg.Compact(gg.MapFlatUniq(self, Game.GetTags)))
}

type Game struct {
	Id       string
	Name     string
	Link     Url
	Img      Url
	Desc     string
	TimeSink TimeSink
	Tags     []Tag
	IsListed bool
}

// Implement `gg.Pker`.
func (self Game) Pk() string {
	if gg.IsZero(self.Id) && gg.IsNotZero(self.Name) {
		panic(gg.Errf(`missing id in %q`, self.Name))
	}
	return self.Id
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
			`is`, `time-sink`,
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
	TagTimeTravel  = Tag{`time_travel`, TagOrd.Next()}

	TagPacifist      = Tag{`pacifist`, TagOrd.Next()}
	TagPhilosophical = Tag{`philosophical`, TagOrd.Next()}
)

// TODO consider additional explanations for some tags.
// TODO "exotic" tag.
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
		AP(`is`, `a-tag`, `type`, `button`, `class`, `tag`),
		self.String(),
	)
}
