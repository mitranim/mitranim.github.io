package main

import "github.com/mitranim/gt"

func initSiteGames() GameColl {
	return makeSiteGames().SortedByKeys(
		`astroneer`,
		`subnautica`,
		`subnautica_below_zero`,
		`undertale`,
		`ftl`,
		`slay_the_spire`,
		`beneath_oresa`,
		`talos_principle`,
		`steins_gate`,
		`darkest_dungeon`,
		`divinity_original_sin_2`,
		`horizon_zero_dawn`,
		`mass_effect_trilogy`,
		`mass_effect_andromeda`,
		`deux_ex_human_revolution`,
		`deus_ex_mankind_divided`,
		`metal_gear_rising_revengeance`,
		`nier_automata`,
		`prey`,
		`control`,
		`doom_2016`,
		`singularity`,
		`borderlands`,
		`borderlands_tps`,
		`borderlands_2`,
		`tales_from_the_borderlands`,
		`witcher_2`,
		`witcher_3`,
		`jedi_fallen_order`,
		`portal`,
		`portal_2`,
		`kotor`,
		`kotor_2`,
		`nwn`,
		`bastion`,
		`half_life`,
		`half_life_2`,
		`half_life_2_episode_one`,
		`half_life_2_episode_two`,
		`no_mans_sky`,
		`warframe`,
	)
}

// Sketch, wildly incomplete.
func makeSiteGames() (out GameColl) {
	out.AddUniq(Games{
		{
			Id:       `ftl`,
			Name:     `FTL: Faster Than Light`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/212680`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/212680/header.jpg`),
			Desc:     `One of the greatest roguelikes. Highly influential, and highly enjoyable. Start on Easy difficulty. One run ≈ 2-3 hours.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagWindows, TagMac, TagConsole, TagMobile,
				Tag2d,
				TagSolo, TagPve,
				TagStrategy, TagTurnBased, TagSciFi, TagRoguelike, TagSpace,
			),
			IsListed: true,
		},
		{
			Id:       `slay_the_spire`,
			Name:     `Slay the Spire`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/646570`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/646570/header.jpg`),
			Desc:     `Strategic roguelike deck building dungeon crawler. One run ≈ 2-3 hours.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, TagMac, TagMobile,
				Tag2d,
				TagSolo, TagPve,
				TagStrategy, TagDeckBuilder, TagRoguelike, TagTurnBased,
			),
			IsListed: true,
		},
		{
			Id:       `beneath_oresa`,
			Name:     `Beneath Oresa`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/1803400`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/1803400/header.jpg`),
			Desc:     `Strategic roguelike deck building dungeon crawler inspired by _Slay the Spire_. One run ≈ 2-3 hours.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagStrategy, TagDeckBuilder, TagRoguelike, TagTurnBased,
			),
			IsListed: true,
		},
		{
			Id:       `undertale`,
			Name:     `Undertale`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/391540`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/391540/header.jpg`),
			Desc:     `Touching, uplifting story. _Highly_ recommended. Don't let pixel art deter you.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagMobile, TagWindows, TagMac,
				Tag2d,
				TagSolo, TagPve,
				TagRolePlay, TagStory, TagFantasy, TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `astroneer`,
			Name:     `Astroneer`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/361420`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/361420/header.jpg`),
			Desc:     `Planetary exploration, survival, resource gathering, building, terraforming. Watch my [video guides](https://www.youtube.com/playlist?list=PLfygJGWNJ-9WaNWXim4P7lLwZ0ooSWLQ4)!`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagThirdPerson,
				TagSolo, TagCoop, TagPve,
				TagCrafting, TagBuilding, TagExploration, TagSpace, TagSciFi, TagOpenWorld,
				TagGrind, TagPacifist,
			),
			IsListed: true,
		},
		{
			Id:       `horizon_zero_dawn`,
			Name:     `Horizon Zero Dawn`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/1151640`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/1151640/header.jpg`),
			Desc:     `True masterpiece of an RPG, set in a post-apocalyptic world where advanced robots roam the wilds. Excellent game design, story, characters, dialogues, music, graphics, UI, marred only by atrocious mouse handling.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagOpenWorld, TagExploration, TagShooter, TagRolePlay, TagStory, TagSciFi,
				TagCrafting,
			),
			IsListed: true,
		},
		{
			Id:       `subnautica`,
			Name:     `Subnautica`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/264710`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/264710/header.jpg`),
			Desc:     `Exploring, surviving, crafting, and building in an amazingly beautiful alien ocean. Doesn't hand-hold, tell you where to go, or what to do. Highly recommended!`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagMac,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagSciFi, TagExploration, TagCrafting, TagBuilding, TagUnderwater, TagStory, TagPacifist,
			),
			IsListed: true,
		},
		{
			Id:       `subnautica_below_zero`,
			Name:     `Subnautica: Below Zero`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/848450`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/848450/header.jpg`),
			Desc:     `_Subnautica_'s expansion that became its own game. Slightly different environment and toys, but basically more of the same. Smaller than the original.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagMac,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagSciFi, TagExploration, TagCrafting, TagBuilding, TagUnderwater, TagStory, TagPacifist,
			),
			IsListed: true,
		},
		{
			Id:       `divinity_original_sin_2`,
			Name:     `Divinity Original Sin 2`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/435150`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/435150/header.jpg`),
			Desc:     `Marvelous RPG with deep mechanics, amazing writing, characters, dialogues, story, music, and more. Check my [tips & tricks](/posts/divinity-original-sin-2)!`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagIsometric,
				TagSolo, TagCoop, TagPve,
				TagRolePlay, TagStory, TagFantasy,
				TagMagic, TagTurnBased, TagParty, TagCrafting,
			),
			IsListed: true,
		},
		{
			Id:       `talos_principle`,
			Name:     `Talos Principle`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/257510`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/257510/header.jpg`),
			Desc:     `Well-designed puzzle game with amazingly deep and thought-provoking monologues and dialogues. Excellent relaxing music. Take your time, play it slow, and _think_.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagWindows, TagMac, TagConsole, TagMobile,
				Tag3d, TagFirstPerson, TagThirdPerson,
				TagSolo, TagPuzzle, TagSciFi, TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `portal`,
			Name:     `Portal`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/400`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/400/header.jpg`),
			Desc:     `Highly influential and acclaimed puzzle game. Excellent in its minimalism. For added context, I recommend playing _Half-Life 2_ first; _Portal_ takes place in the same world, at the same time.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagSolo,
				TagFirstPerson, TagPuzzle, TagSciFi, TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `portal_2`,
			Name:     `Portal 2`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/620`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/620/header.jpg`),
			Desc:     `Sequel to _Portal_, with more mechanics and a longer campaign. Well-made, well-rounded, even if not particularly revolutionary.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagCoop,
				TagPuzzle, TagSciFi, TagPacifist,
			),
			IsListed: true,
		},
		{
			Id:       `darkest_dungeon`,
			Name:     `Darkest Dungeon`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/262060`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/262060/header.jpg`),
			Desc:     `Turn-based, party-based dungeon crawler inspired by Lovecraftian horror stories, with fairly unique mechanics.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagMac, TagConsole, TagMobile,
				Tag2d,
				TagSolo, TagPve,
				TagStrategy, TagTurnBased, TagParty, TagFantasy,
			),
			IsListed: true,
		},
		{
			Id:       `mass_effect_trilogy`,
			Name:     `Mass Effect Trilogy (Legendary Edition)`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/1328670`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/1328670/header.jpg`),
			Desc:     `One of the greatest sci-fi RPG experiences ever made. Huge influence on other games and on many people's lives. Mechanics and graphics are outdated by modern standards, but the setting, characters, dialogues, story, music and more are timeless.`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagShooter, TagStory, TagRolePlay, TagSciFi, TagParty, TagMagic,
				TagExploration, TagSpace,
			),
			IsListed: true,
		},
		{
			Id:       `mass_effect_andromeda`,
			Name:     `Mass Effect: Andromeda`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/1238000`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/1238000/header.jpg`),
			Desc:     `Worthy successor to previous ME titles. Excellent sci-fi RPG in its own right. Requires unfucking. Read [my recommendations](/posts/andromeda).`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagShooter, TagStory, TagRolePlay, TagSciFi, TagParty, TagMagic,
				TagOpenWorld, TagExploration, TagSpace, TagCrafting,
			),
			IsListed: true,
		},
		{
			Id:       `deux_ex_human_revolution`,
			Name:     `Deus Ex: Human Revolution`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/238010`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/238010/header.jpg`),
			Desc:     `Marvelous RPG set in near future, exploring themes of human cyborg augmentation, secret societies, and more. Worthy successor to the original _Deus Ex_. I wrote a [review](https://blog-blogger.mitranim.com/2012/11/a-thank-you-to-eidos-montreal-for-dehr.html) that's closer to a love letter!`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson, TagThirdPerson,
				TagSolo, TagPve,
				TagShooter, TagSciFi,
				TagRolePlay, TagStory, TagStealth,
				TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `deus_ex_mankind_divided`,
			Name:     `Deus Ex: Mankind Divided`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/337000`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/337000/header.jpg`),
			Desc:     `Excellent sequel to _DEHR_. More of the same, with massively improved graphics and lowered stakes. Great writing, dialogues, characters, music, atmosphere, gameplay. Highly satisfying to a _Deus Ex_ fan.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson, TagThirdPerson,
				TagSolo, TagPve,
				TagShooter, TagSciFi,
				TagRolePlay, TagStory, TagStealth,
				TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `metal_gear_rising_revengeance`,
			Name:     `Metal Gear Rising: Revengeance`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/235460`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/235460/header.jpg`),
			Desc:     `Badass melee slasher with cyborgs, vibroblades, giant robots, metal music, and some surprisingly thoughtful dialogues. Requires a gamepad.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagSciFi, TagStory, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `nier_automata`,
			Name:     `Nier Automata`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/524220`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/524220/header.jpg`),
			Desc:     `Surprisingly thought-provoking and touching. Brought me to tears. _Avoid external spoilers!_ Go for "true ending". Use original Japanese voiceovers and English subtitles. Requires unfucking via external tools, check [PCGW](https://www.pcgamingwiki.com/wiki/Nier_Automata). Requires a gamepad.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagShooter, TagSciFi, TagStory, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `prey`,
			Name:     `Prey`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/480490`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/480490/header.jpg`),
			Desc:     `True successor to System Shock. Stealth gameplay on an abandoned space station infested by truly alien outsiders. Love the environments, gameplay, music. Thoughtful. See [PCGW](https://www.pcgamingwiki.com/wiki/Prey_\(2017\)) for tweaks and tools.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagSciFi, TagSpace, TagCrafting, TagStory, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `control`,
			Name:     `Control`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/870780`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/870780/header.jpg`),
			Desc:     `Shooter / magic hybrid in a setting where all the paranormal sci-fi tropes are real. If you've read "Понедельник начинается в субботу", this may feel familiar.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagShooter, TagMagic, TagSciFi, TagStory, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `witcher_2`,
			Name:     `Witcher 2: Assassins of Kings`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/20920`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/20920/header.jpg`),
			Desc:     `RPG following after the events of the original _Witcher_ books. Excellent writing, characters, dialogues, music. Requires unfucking, check [my recommendations](/posts/witcher).`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagMagic, TagStory, TagRolePlay, TagFantasy, TagCrafting,
			),
			IsListed: true,
		},
		{
			Id:       `witcher_3`,
			Name:     `Witcher 3: Wild Hunt`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/292030`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/292030/header.jpg`),
			Desc:     `Legendary RPG beloved by all. Excellent writing, characters, dialogues, music, combat mechanics, and more, and it is **HUGE**. Requires unfucking, read [my recommendations](/posts/witcher).`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagWindows, TagConsole,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagMagic, TagStory, TagRolePlay, TagFantasy,
				TagOpenWorld, TagExploration, TagCrafting,
			),
			IsListed: true,
		},
		{
			Id:       `warframe`,
			Name:     `Warframe`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/230410`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/230410/header.jpg`),
			Desc:     `Third person shooter / melee hybrid with deep mechanics and a bewildering bevy of toys. Read [stupid headcanon](/posts/warframe-headcanon) authored by me and friends (spoilers!).`,
			TimeSink: TimeSinkExtreme,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagThirdPerson,
				TagSolo, TagCoop, TagPve,
				TagShooter, TagMelee, TagCrafting, TagBuilding, TagGrind, TagSciFi,
			),
			IsListed: true,
		},
		{
			Id:       `no_mans_sky`,
			Name:     `No Man's Sky`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/275850`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/275850/header.jpg`),
			Desc:     `Exploring an open-world galaxy where everything is procedurally generated. Lots of features. Use a [save editor](https://github.com/zencq/NomNom) to minimize grinding.`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson, TagThirdPerson,
				TagSolo, TagCoop, TagPve,
				TagShooter, TagSciFi, TagSpace,
				TagOpenWorld, TagExploration, TagCrafting, TagBuilding, TagGrind,
			),
			IsListed: true,
		},
		{
			Id:       `kotor`,
			Name:     `Star Wars: Knights of the Old Republic`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/32370`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/32370/header.jpg`),
			Desc:     `Classic SW RPG with excellent writing and dialogues. Uses simplified D&D mechanics. Dated but well worth it. Avoid spoilers!`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagMobile,
				TagWindows, TagMac,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagMagic, TagTurnBased, TagSciFi, TagFantasy, TagParty,
				TagRolePlay, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `kotor_2`,
			Name:     `Star Wars: Knights of the Old Republic 2: the Sith Lords`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/208580`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/208580/header.jpg`),
			Desc:     `Thoughtful sequel exploring the fallout of your "heroic deeds" in the previous game. Excellent writing and dialogues. Requires community patch.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagMobile,
				TagWindows, TagMac,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagMagic, TagTurnBased, TagSciFi, TagFantasy, TagParty,
				TagRolePlay, TagStory, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `steins_gate`,
			Name:     `Steins Gate`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/412830`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/412830/header.jpg`),
			Desc:     `One of the greatest time travel stories. Gripping, coherent, based on well-researched science, brimming with otaku references, a true geek party! Also watch the anime: _Steins Gate_, then _Steins Gate 0_. Use essential [community patch](https://sonome.dareno.me/projects/sghd.html).`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, Tag2d,
				TagSolo, TagPuzzle, TagSciFi, TagStory, TagTimeTravel,
				TagPacifist, TagPhilosophical,
			),
			IsListed: true,
		},
		{
			Id:       `half_life`,
			Name:     `Half-Life: Source`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/280`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/280/header.jpg`),
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagStory,
			),
			IsListed: false,
		},
		{
			Id:       `half_life_2`,
			Name:     `Half-Life 2`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/220`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/220/header.jpg`),
			Desc:     `Influential first-person shooter with fluid on-the-fly story narration. Somewhat primitive by modern standards, but still recommended, especially if you plan to play _Portal_.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `half_life_2_episode_one`,
			Name:     `Half-Life 2: Episode One`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/380`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/380/header.jpg`),
			Desc:     `Continues the story of _Half-Life 2_, followed by _Episode Two_.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `half_life_2_episode_two`,
			Name:     `Half-Life 2: Episode Two`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/420`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/420/header.jpg`),
			Desc:     `Continues the story of _Half-Life 2: Episode One_. The story is left unfinished.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `jedi_fallen_order`,
			Name:     `Star Wars: Jedi: Fallen Order`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/1172380`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/1172380/header.jpg`),
			Desc:     `Apprentice Jedi exploring various planets and fighting stormtroopers. Mechanically similar to Dark Souls. Somewhat simplistic but well-delivered. Graphics require unfucking via ReShade.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagThirdPerson,
				TagSolo, TagPve,
				TagMelee, TagMagic, TagSciFi, TagFantasy,
				TagExploration, TagOpenWorld, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `singularity`,
			Name:     `Singularity`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/42670`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/42670/header.jpg`),
			Desc:     `Well-wrought time travel story and a competent first-person shooter. Set on an abandoned Soviet research island. Very atmospheric. Bilingual bonus for Russian speakers.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve,
				TagShooter, TagSciFi, TagStory, TagTimeTravel,
			),
			IsListed: true,
		},
		{
			Id:       `borderlands`,
			Name:     `Borderlands`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/729040`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/729040/header.jpg`),
			Desc:     `Stylish, atmospheric shooter / looter with unique enemy and environment designs, and loot mechanics similar to _Diablo_.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, TagMac,
				Tag3d, TagFirstPerson,
				TagSolo, TagCoop, TagPve,
				TagShooter, TagSciFi, TagExploration, TagOpenWorld, TagStory, TagGrind,
			),
			IsListed: true,
		},
		{
			Id:       `borderlands_tps`,
			Name:     `Borderlands: The Pre-Sequel`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/261640`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/261640/header.jpg`),
			Desc:     `More of the same, with its own unique environments and mechanics. Does a great job bridging the story between _BL1_ and _BL2_.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, TagMac,
				Tag3d, TagFirstPerson,
				TagSolo, TagCoop, TagPve,
				TagShooter, TagSciFi, TagExploration, TagOpenWorld, TagStory, TagGrind,
			),
			IsListed: true,
		},
		{
			Id:       `borderlands_2`,
			Name:     `Borderlands 2`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/49520`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/49520/header.jpg`),
			Desc:     `The ultimate Borderlands game. Culmination of the story, with highly memorable characters and environments, and with lots of excellent DLCs. Should be played after _BL:TPS_.`,
			TimeSink: TimeSinkDanger,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, TagMac,
				Tag3d, TagFirstPerson,
				TagSolo, TagCoop, TagPve,
				TagShooter, TagSciFi, TagExploration, TagOpenWorld, TagStory, TagGrind,
			),
			IsListed: true,
		},
		{
			Id:       `tales_from_the_borderlands`,
			Name:     `Tales from the Borderlands`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/330830`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/330830/header.jpg`),
			Desc:     `Pure story / adventure set on Pandora after _Borderlands 2_, featuring new protagonists. If you loved previous Borderlands games, this is worth a playthrough.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows, TagMac,
				Tag3d, TagThirdPerson,
				TagSolo, TagPuzzle, TagSciFi, TagStory, TagPacifist,
			),
			IsListed: true,
		},
		{
			Id:       `doom_2016`,
			Name:     `Doom (2016)`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/379720`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/379720/header.jpg`),
			Desc:     `Excellent first-person shooter. Competent sequel to the classic Doom games. Doesn't outstay its welcome. Read my [impressions](/posts/game-impressions-doom-2016) and [tips & tricks](/posts/tips-and-tricks-doom-2016)!`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagWindows,
				Tag3d, TagFirstPerson,
				TagSolo, TagPve, TagShooter, TagSciFi,
			),
			IsListed: true,
		},
		{
			Id:       `nwn`,
			Name:     `Neverwinter Nights`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/704450`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/704450/header.jpg`),
			Desc:     `Classic RPG using D&D mechanics. Don't be deterred by boring OG campaign; expansions get better, and there's fun in the mechanics. Use console commands to un-slog.`,
			TimeSink: TimeSinkModerate,
			Tags: Slice(
				TagPc, TagWindows,
				Tag3d, TagThirdPerson, TagIsometric,
				TagSolo, TagCoop, TagPve,
				TagMelee, TagMagic, TagTurnBased, TagFantasy, TagParty, TagRolePlay, TagStory,
			),
			IsListed: true,
		},
		{
			Id:       `bastion`,
			Name:     `Bastion`,
			Link:     gt.ParseNullUrl(`https://store.steampowered.com/app/107100`),
			Img:      gt.ParseNullUrl(`https://cdn.cloudflare.steamstatic.com/steam/apps/107100/header.jpg`),
			Desc:     `Charming top-down adventure game. Touching story, music, narration. Originally for mobile devices, simplistic but worth a playthrough regardless.`,
			TimeSink: TimeSinkSafe,
			Tags: Slice(
				TagPc, TagConsole, TagMobile, TagWindows, TagMac,
				Tag2d, TagThirdPerson, TagIsometric,
				TagSolo, TagPve, TagFantasy, TagStory,
			),
			IsListed: true,
		},
		// {
		// 	Id:       `star_craft_2`,
		// 	Name:     `StarCraft 2`,
		// 	Link:     gt.ParseNullUrl(`https://starcraft2.com`),
		// 	Img:      gt.ParseNullUrl(``), // Needs an image.
		// 	Desc:     ``,
		// 	TimeSink: TimeSinkModerate,
		// 	IsListed: false,
		// },
	}...)
	return
}
