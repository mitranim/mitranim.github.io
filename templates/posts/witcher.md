{{mdToToc .MdTpl}}

## Order

The in-universe events happen in this order, and should be read and played in the same order.

* The [book](#book). The original compilation, _without_ "Season of Storms" or any later works.
* [Witcher 1](#witcher-1).
* [Witcher 2](#witcher-2).
* [Witcher 3](#witcher-3).

## Book

I can highly recommend the Russian version of the book. The original is in Polish, which should be comparable. I haven't tried the English version, or "Season of Storms" or any later works, and can't vouch for their quality.

Look for an EPUB version and use a decent reader. Avoid unusable formats like PDF.

## Witcher 1

It's been a long time since I played it, so my ability to provide advice is limited to the following:

* According to rumors, a remake is in the works. Consider waiting until it's released (and patched, and modded).
* Check the game's [page on PC Gaming Wiki](https://www.pcgamingwiki.com/wiki/Witcher) and follow its recommendations.
  * The article provides help with various issues, and recommends various modifications, including graphical improvements.
* Use [speedhacks](#speedhacks) to minimize slog.

## Witcher 2

* Check the game's [page on PC Gaming Wiki](https://www.pcgamingwiki.com/wiki/Witcher_2) and follow its recommendations.
* Prepare for 2 playthroughs. A _huge_ chunk of the game's content is split between mutually exclusive paths: Roche/Kaedwen and Iorveth/Scoia'tael. As a result, the game requires 2 runs to fully complete.
  * If you intend to setup a "perfect" save for importing into Witcher 3, it should be on the Roche/Kaedwen path. In Witcher 3, Roche is present, while Iorveth was cut due to production constraints.
* Essential tweaks:
  * Unfuck config files, as recommended in the PCGW article linked above.
  * [AI Upscale Textures](https://www.nexusmods.com/witcher2/mods/889). Requires also editing config files.
  * [Zero Weight](https://www.nexusmods.com/witcher2/mods/265).
  * Use [speedhacks](#speedhacks) to minimize slog.
* Optional tweaks:
  * Cheat gold (via Cheat Engine), when ferrying loot to vendors gets old.
  * Cheat XP on subsequent playthroughs as a form of "new game plus".
  * Cheat talent points if you feel the default system is too restrictive.

## Witcher 3

Check the game's [page on PC Gaming Wiki](https://www.pcgamingwiki.com/wiki/Witcher_3) and follow its recommendations. Also read forth. The rest of this article is all about Witcher 3.

### Version

At the time of writing, Witcher 3 has a "classic edition" (_CE_, version 1.32) and a "next gen" edition (_NGE_, version 4+). Steam installs NGE by default, but also allows to choose CE. I've played only NGE, and have no opinion on CE. The [mods](#mods) listed below are all compatible with NGE. Some of them may not exist for CE.

### Graphics

* In the game's own custom launcher, switch from DirectX 12 to DirectX 11.
  * In the versions that I've played (NGE 4.0 and slightly above), DX12 has significantly worse performance, with no observable benefit.
* Disable all forms of blur.
* If your display has a decent DPI (120 and higher), disable anti-aliasing.
  * In this game, many forms of AA make the picture much blurrier, while other forms of AA make no observable change to the picture, with no in-between.
* If your GPU has enough memory, crank up texture quality.
* The other settings are mostly expendable in favor of FPS. Tweak them to ensure a comfortable framerate.

### Launch Flags

After [switching](#graphics) from DX12 to DX11, disable the launcher by adding `--launcher-skip` to the game's launch flags.

### Config Tweaks

Check the [PCGW article](https://www.pcgamingwiki.com/wiki/Witcher_3) to find the location of the `user.settings` file for your system and game version. Note that when running DX12, you would be using the file `dx12user.settings`, but you should really [switch](#graphics) to DX11 for better performance. When editing this file, you can simply append entries. The game reorders the content automatically.

By default, the game has extremely obnoxious videos on the loading screen. Add the following to disable them.

```
[LoadingScreen/Debug]
DisableVideos=true
```

Add the following to enable the console. Afterwards, it can be summoned by the `~` key, unless you change the key by editing config files. See [console](#console) tips below.

```
[General]
DBGConsoleOn=true
```

### Mods

This section may seem daunting, but trust me, modding the game is **well worth the effort**!

Most mods can be installed and updated via [Vortex](https://www.nexusmods.com/about/vortex/), the mod manager from [Nexus](https://www.nexusmods.com). I've been using Vortex without issues. However, the authors of Witcher 3 mods tend to recommend against Vortex, and in favor of [the Witcher 3 Mod Manager](https://github.com/stefan3372/The-WItcher-3-Mod-manager), which I haven't tried.

Some mods require manual install. Some require additional tweaks in `input.settings`. Some require invoking the Menu Filelist Updater which is listed below. Check each mod's description for setup instructions.

* [Menu Filelist Updater](https://www.nexusmods.com/witcher3/mods/7171)
  * Some mods have an ingame menu which can be enabled by editing a specific config file. This small program updates the config for you.
  * [Source repository](https://github.com/Aelto/tw3-menufilelist-updater)
* [Brothers in Arms](https://www.nexusmods.com/witcher3/mods/7329)
  * Community patch that fixes lots of bugs and restores significant amounts of missing content.
* [Skip Movies](https://www.nexusmods.com/witcher3/mods/358)
* [Instant Tooltips](https://www.nexusmods.com/witcher3/mods/2019)
  * Needs higher load priority than Smooth GUI. They both work.
* [Smooth GUI](https://www.nexusmods.com/witcher3/mods/7730)
* Increase or remove inventory weight limit. Pick one of these mods:
  * [No Inventory Weight Limit](https://www.nexusmods.com/witcher3/mods/7159)
  * [9000 Weight Saddlebags](https://www.nexusmods.com/witcher3/mods/2948)
* [Item Levels Normalized](https://www.nexusmods.com/witcher3/mods/3095)
* [Remove Item Level Requirements](https://www.nexusmods.com/witcher3/mods/1542)
* [Indestructible Items](https://www.nexusmods.com/witcher3/mods/342)
* [Instant Witcher Senses](https://www.nexusmods.com/witcher3/mods/2428)
* [Sprint in Witcher Senses](https://www.nexusmods.com/witcher3/mods/7407)
* [Jump while in Wicher Senses](https://www.nexusmods.com/witcher3/mods/3665)
* [No Witcher Sense Zoom FX plus Toggle and Range](https://www.nexusmods.com/witcher3/mods/351)
  * Combine two files: `NoWitcherSenseFX` and `WitcherSenseDoubleRange`.
  * `NoWitcherSenseFX` is particularly important because in addition to making witcher sense more usable, it fixes an annoying NGE bug where witcher sense SFX gradually breaks music.
* [AutoLoot All-in-One](https://www.nexusmods.com/witcher3/mods/7198)
  * Allows to instantly loot containers with 1 keypress, disable theft mechanic, make hotkey to loot everything around you. Configurable.
  * Edit `input.settings` to configure a hotkey. I use `IK_Mouse5=(Action=AutoLootRadius)`. The hotkey must be placed in multiple sections.
  * When configuring this mod through ingame menus, when configuring notifications, set "Use Action Log Notification" to "false" because this breaks the ingame action log completely.
* [Stack Your Items](https://www.nexusmods.com/witcher3/mods/7175)
* [Disable Fall Damage](https://www.nexusmods.com/witcher3/mods/7219)
* [Next Gen Movement Input Lag Fix](https://www.nexusmods.com/witcher3/mods/7586)
* [BetterMovement](https://www.nexusmods.com/witcher3/mods/7591)
  * Allows to sprint x2 faster and swim x5 faster, configurable.
* [Alternate Horse Sprint](https://www.nexusmods.com/witcher3/mods/5510)
  * Allows separate keys for canter and gallop.
  * Requires editing `input.settings`. Beware: in settings, terms "canter" and "gallop" are reversed!
* [Galloping In Cities](https://www.nexusmods.com/witcher3/mods/385)
* [Sensible Map Borders](https://www.nexusmods.com/witcher3/mods/7393)
* [Any skill in Mutation Slots](https://www.nexusmods.com/witcher3/mods/7333)
* [Fast Travel from Anywhere](https://www.nexusmods.com/witcher3/mods/324)
  * To avoid breaking quest triggers, avoid fast traveling during scripted sequences.
* [Loot Bags Glow without Witcher Senses](https://www.nexusmods.com/witcher3/mods/3050)
* [Improved Fist Fights](https://www.nexusmods.com/witcher3/mods/3703)
  * Makes fist fights more plausible and less boring. Makes Geralt about as strong as opponents are. Usually on high difficulties Geralt is x10 weaker.
* [All Quest Objectives On Map](https://www.nexusmods.com/witcher3/mods/943)
  * Needs higher load priority than Smooth GUI. They both work.
* [Missing Gwent Cards Tracker](https://www.nexusmods.com/witcher3/mods/7179)
  * Makes the book "Miraculous Guide to Gwent" _actually_ a miraculous guide to finding the missing cards.
* [The Two Gwent Stores](https://www.nexusmods.com/witcher3/mods/7616)
  * Allows to quickly get all Gwent cards that aren't acquired through quests. Can be considered a cheat.
  * I prefer to simply add Gwent cards via a [console command](#console-gwent) and skip this mod.
* [Oils Potions Bombs Tab by Default](https://www.nexusmods.com/witcher3/mods/7509)
* [Smart Alcohol Refill Selection](https://www.nexusmods.com/witcher3/mods/4619)
* [Slimmer Griffin Armor](https://www.nexusmods.com/witcher3/mods/5631)
  * Requires manual install. Well worth it. The Grandmaster version of this gear has been my favorite ever since installing the mod.
* [Less Junk](https://www.nexusmods.com/witcher3/mods/7239)
  * Questionable. Doesn't stop you from finding too much junk.
* [Cut-Throat Razor](https://www.nexusmods.com/witcher3/mods/7858)
  * Alternatively, shave via [console](#console-beard).
* [Omelet Quest](https://www.nexusmods.com/witcher3/mods/7213)
  * Restores a rather fun piece of content that was unfairly cut from the game.
* [Unlimited Enchanting](https://www.nexusmods.com/witcher3/mods/7831)
  * Significantly improves the HoS runesmith, allowing arbitrary enchants on arbitrary gear pieces.
* [No Automatic Trophy Switch](https://www.nexusmods.com/witcher3/mods/8022)

Potentially interesting mods that I haven't tried:

* [No Levels](https://www.nexusmods.com/witcher3/mods/3605)
* [Essential Weapon Rework](https://www.nexusmods.com/witcher3/mods/5104)
* [Gwent Redux](https://www.nexusmods.com/witcher3/mods/4287)
  * I think Gwent is already well-balanced for PvE and does not need a rework.
* [Friendly HUD](https://www.nexusmods.com/witcher3/mods/7290)
  * May cause weird issues unrelated to the HUD.
* [Uniform Horse Armour Stats](https://www.nexusmods.com/witcher3/mods/7930)
* [Improved Horse Controls](https://www.nexusmods.com/witcher3/mods/7229)
* [Sort Everything](https://www.nexusmods.com/witcher3/mods/1710)
  * Has conflicts with Smooth GUI, requires tweaking mod load order.
* [Fix Ciri Invulnerability](https://www.nexusmods.com/witcher3/mods/7919)
  * Requires manual install.

### Misc Tips

When starting a new game, either import a Witcher 2 save (if you've made the right choices), or choose "_do_ simulate Witcher 2 save". If you choose the latter, at some point in Witcher 3 you will be questioned to determine various Witcher 2 choices that have an effect in Witcher 3.

Disable tutorial popups ASAP. You'll learn the game just fine without them, and I personally found them extremely obnoxious. They also break the tourney horse race in B&W.

Ignore damage numbers in skill tooltips. Damage of Signs scales with enemy health (current health for Aard, maximum health for other Signs). Tooltips always lie. Experiment with everything.

You get less XP for lower-level enemies and quests. My suggestion: don't think about it. You'll end up at roughly the same maximum level regardless of completion order. The developers did this to prevent overleveling, which would make combat too easy and less fun.

As soon as is practical, raise combat and Gwent difficulty to max. Makes things more interesting.

### Combat Tips

Use the Fleet Footed ability, and learn when to use small dodge (default <key>Alt</key>) and when to use dodge roll (default <key>Space</key>). Both are extremely useful.

Loading screen tips tell you that light armor is best for stamina regen. This is misleading. Medium armor with Griffin School Techiques is far better for stamina regen.

Fighting multiple opponents may feel very different from fighting one. They don't coordinate attacks and may attack from offscreen. This forces you to dodge or parry more often, making combat significantly slower. Learn to enjoy this. Prolonging combat is a positive rather than a negative, because it makes transitions between ambient and combat music less grating on the ears. Many combat tracks in Witcher 3 can be annoying at the start, especially if heard frequently, but eventually pick up and become interesting. Long combat makes the music better.

Initially, Signs are underwhelming. Their effects are weak and stamina regen is slow. However, Signs become extremely strong with high Sign intensity and stamina regen, obtainable later in the game via slottable greater mutagens, HoS enchants, B&W mutations, Grandmaster Griffin gear. A magic-oriented build has been my favorite in multiple New Game+ playthroughs.

### Console

See [config tweaks](#config-tweaks) for enabling the console. Also see this outdated but still useful reference: https://commands.gg/witcher3. Some particularly useful commands are highlighted below.

#### Console: Gwent

The following command adds almost every Gwent card, in its maximum total obtainable amount:

    addgwintcards

Limitations:

* Does not add cards which are present in the base North deck. Can be added by running `additem` with the appropriate item codes, which can be found on the Witcher wiki. See the [Gwent](https://witcher.fandom.com/wiki/Gwent) article.
* Increasing the total amount of cards increases lag in some parts of the Gwent UI. I suggest using this command 0 times in the first playthrough and 1 or 2 times in New Game+.

#### Console: Beard

Shave:

    setbeard(0)

Maximum beard:

    setbeard(1)

## Speedhacks

Witcher 1 and Witcher 2 can feel like a slog when backtracking or otherwise running around, with no ingame way to speed that up. Speedhacks are highly recommended. See my [post on speedhacking](/posts/speed).
