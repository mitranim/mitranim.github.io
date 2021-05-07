(Self-reminder: the top of the post should say that it's constantly updated and include a link to the file history in git. Recently added entries should be marked with a colored "new".)

I've spent an ungodly amount of time playing Warframe. It's one of my favorite games of all time. Instead of talking about what's great in it, I'd rather talk about what's not — so we can make it even better.

The target audience for this post are Digital Extremes developers and game designers. If anyone reading this has a contact at DE, please forward it to them. Forward it to other players too, to increase the outreach.

## Table of Contents

{{tableOfContents .InputPath}}

## Bugs

... so many to list. Where to start? Might limit this post to _mechanics_ since everyone can agree on bugs.

## Controls

### Held Controls

Warframe is currently extremely bad at handling "held" controls such as sprint, aim, shoot/charge, Void Mode, crouch, and more. This must be ~~unfucked~~ rectified.

Currently most of these actions only activate if you press the key at the "right" time. I want them to activate whenever possible if the key is held, no matter _when_ you started holding it.

One example is Transference into Void Mode. Currently after pressing 5 to Transfer out, you have to wait some _indeterminate amount of time_ which is _host latency dependent_ before you're allowed to activate Void Mode and use Void Dash. If you press it too soon, Void Mode won't activate, and trying to use Dash will perform a limp jump. What _needs_ to happen is that you press Transfer and hold crouch — once the operator is out, it enters Void Mode as soon as possible, and pressing jump is guaranteed to perform a Dash. The same applies to activating Void Mode at other times, but Transference is the most egregious one.

Another example is shooting and charging. With full-auto guns, it's not so bad. Most of them will shoot whenever you're holding the key regardless of the previous actions. There's a few full-auto guns that don't behave this way. But the worst offenders are "charge" triggers such as bows, Opticor variants, etc.. What needs to happen is simple: whenever you're holding the fire key and it's possible for the shot to charge, it must charge. No exceptions, no excuses. What _currently_ happens is if you started holding the key during a recovery animation, a roll, a knockdown, or simply at the wrong time after the shot, it won't start charging until you release and press again, at some indeterminate "right" time. Some actions such as rolling also break the currently held charge, and it fails to charge again. This makes controlling those guns inconsistent and annoying.

Aiming has a similar problem. Currently jumping while aiming kicks you out of aiming. This is somewhat understandable because aim gliding exists. But a better behavior would be to automatically activate aim gliding at the peak of the jump if the aim key was held before the jump. This behavior should apply if you get knocked _up_ for any other reason, such as through a Butcher or Terra Trencher ground slam. If you fall _down_ from a ledge while holding aim, aim gliding should activate immediately; IIRC it already works this way.

Sprinting tends to work most of the time, but a few actions kick you out of sprint and it fails to reactivate until you re-press the key. Just like other held actions, it needs a better implementation that always works.

### Action Lockouts

Warframe lacks the concept of "UI actions" that can be performed anytime independently of everything else. A variety of controls that I consider UI actions, or non-warframe actions, are currently treated like warframe abilities, and can't be performed during a recovery animation, a melee animation, and in a variety of other circumstances. For me, this is very inconvenient. Examples of such actions:

* Toggling the map.
* Triggering a weapon switch or selecting "held" gear such as scanners. I want it to be reflected in the UI instantly and queued up after the current animation, knockdown recovery, etc.; there should _never_ be a case where the weapon switch command is simply ignored. Currently it's very inconsistent: if you trigger a weapon switch and immediately shoot, sometimes it will finish the switch before shooting, but sometimes shooting interrupts the switch and uses the previous weapon. This must never happen.
* Cycling ability variants such as Wisp reservoirs, Ivara arrows, etc.. Currently when placing reservoirs, you have to wait for the placement animation to finish before you can cycle to the next reservoir. If you want to place them as fast as possible, this means repeatedly holding and releasing the ability key, but releasing it too soon may accidentally place the same reservoir again. It's very inconvenient and unnecessary. Cycling ability variants must be made a UI action doable at any time. (Or better yet, replaced with instant activation, see the Ability Cycling section.)
* Transference. Currently Transference is treated as a warframe ability that can't be performed during recovery animations. This often means mashing the button while waiting for some animation to finish. This is made worse by the _delay_ on Transference-out, combined with _no delay_ on Transference-in, which sometimes causes you to instantly Transfer back in after transfering out, wasting time and forcing another attempt. This is made worse by the fact that the delays are _latency dependent_ which varies between hosts so you can't even get used to a consistent timing. To make Transference less annoying, it should be made usable during any warframe action, and any artificial delays on Transference-out and -in should be removed. In addition, Transference should be made client-side rather than host-based and latency-dependent, but that's an entirely separate change. Even converting it to a non-warframe action and removing what seems like client-side activation delays would be a big improvement.

### Ability Cycling

Some warframes have abilities which cycle between variants: Wisp reservoirs, Ivara arrows, Grendel feasts, Titania tributes. Currently you activate by holding and cycle by tapping. There's an option to invert it: activate by tapping and cycle by holding. In both cases, you have to alternate tapping and holding, which is fiddly and error-prone, and you have to watch which variant is currently selected. The input method delays your actions.

I want to _instantly_ activate any variant, with _no cycling_, ideally without adding more hotkeys. I will suggest one _possible_ option, but I'm not arguing specifically for this input method. It's just an example.

Let's say we number ability variants (not just abilities) from 1 to 4, and activate them by pressing two ability keys at once, except for the variant whose number matches the number of its ability.

Confusing? Example on Ivara:

  * Tap 1 to shoot Cloaking Arrow.
  * Hold 1 and tap 2 to shoot Dashwire Arrow.
  * Hold 1 and tap 3 to shoot Noise Arrow.
  * Hold 1 and tap 4 to shoot Sleep Arrow.

On Titania:

  * Hold 2 and tap 1 for Thorns Tribute.
  * Tap 2 for Dust Tribute.
  * Hold 2 and tap 3 for Full Moon Tribute.
  * Hold 2 and tap 4 for Entangle Tribute.

This can be instantaneous, like bullet jumping by pressing crouch + jump.

This convention probably wouldn't work for gamepads. For them, we could keep the cycling, or try to come up with a more flexible convention. The important part is being able to activate specific variants instantly.

### Always-Sprint

The "toggle sprint" option never worked properly. A variety of abilities toggle sprinting off. Frankly, this is insane. The game should simply remember your last mode (sprinting or not sprinting), _persist it between missions_, and _never change it by itself_. Shapeshifting abilities such as Wukong's Cloudwalker or Inaros' Sandstorm should have no interaction with sprinting because they have their own movement speeds. Shapeshifting abilities that _are_ affected by sprinting should not mess with it; you said it must be on, so it must be on. Ivara's Prowl should suppress "automatic" sprinting for its duration without toggling it off. Tapping the sprint key while Prowling should always force the "on" behavior, breaking into a sprint. Any other effect that breaks your Prowl, such as bullet jumping or getting nullified, should restore sprinting to the state it had before activating Prowl.

It should be noted that before the archwing rework brought to us by Empyrean, archwing controls on outdoor maps were practically incompatible with always-on sprinting, but thanks to the rework, always-on sprint is now perfectly viable for AW.

I should mention that in operator mode, many actions don't work while sprinting. Even some really basic stuff like activating objects or opening containers. As a result, you're discouraged from sprinting in operator despite having the option. This also needs to be fixed to make always-sprint usable.

### Melee Heavy Attacks

I want a dedicated hotkey that always performs a melee heavy attack no matter which weapon is selected. In fact, this hotkey _exists but does not actually do that_. Instead, it acts identically to the combined "alt fire / heavy attack" hotkey, requiring you to perform a quick melee or lock into melee before using a heavy. This needs to be fixed. The lack of a functioning heavy attack hotkey is the main reason why I don't run heavy attacks builds; constantly minding the selected weapon and switching with quick melee is fiddlesome.

### Gear UI

As a keyboard and mouse user, I loath the gear wheel and want a gear _grid_ that shows _all gear at once_ without having to wiggle your mouse like drawing a spiral in MS Paint. I understand that Warframe supports dynamically switching between KBM and gamepad, and there's value in having a consistent gear UI that works for both. But I find the wheel so abysmally suboptimal for KBM. Can we get an all-at-once grid as an option?

Obviously the screen size is limited, so with a very large amount of equipped gear, you'd have to scroll down to see the rest of it. This just means scaling down the icons and making it denser. Have it fit, let's say, 64 gear icons on the screen. I have close to 30 gear items equipped and would add maybe a handful more.

Of course, the same UI should apply to emotes. Whether the user prefers the grid or the wheel, it should be consistently used for both gear and emotes.

I could also see the grid UI being preferred by gamepad users if navigated by a cursor that _consistently_ spawns in the middle of the grid, rather than some random place. Taking the cursor a known distance in a known direction could be at least as fast as scrolling through the wheel.

## Game Mechanics

### Nullifiers

Instead of dispelling abilities and objects, they should temporarily suppress their _effects_ within the bubble. For example, entering a bubble with Rhino's Iron Skin should suppress the effects of Iron Skin and any other active casts, and prevent ability casts, but after leaving the bubble, Iron Skin and any other buffs should still be there. Buff duration should still tick down, and channeled abilities should still drain energy. Contact of a null bubble with ability-created objects such as Wisp's reservoirs, Limbo's Cataclysm, Frost's Snow Globe, should temporarily disable their effects, possibly indicated by visual transparency, but once the null bubble is gone, the objects' effects should be restored. The link between the objects and the warframe, such as the duration of Limbo's Cataclysm or the restriction on the maximum number of reservoirs and Snow Globes, should persist throughout the null effect.

Null bubble preventing Transference (either in or out) never made any in-universe sense to me, and I don't see it serving a useful mechanical role. To me, operator abilities feel somewhat like "cheating", acting outside the scope of "normal" warframe tech, and being limited by a nullifier effect aimed at _warframe_ abilities makes no sense. I would prefer this restriction removed.

If I recall correctly, some Comba units emit a pulse that disables ability _activation_ without dispelling active abilities. This seems like a step in the right direction. The nullifier alteration proposed above is more complete, as it also suppresses active abilities and the effects of conjured objects.

### Movement Stoppers

* Full-body ability animations that stop movement.
* Summoning AWG.
* Shooting charged AWG at full charge.
* Melee combos that force movement.
* ... and more.

I strongly dislike forced stops, forced self-stun, forced self-stagger, and don't understand why this exists in a game about agile space ninjas.

Let's start with abilities that stop movement. I don't care if the animator had a cool idea for a theatric bow or some other cool gesture. I simply want my warframe to move while casting abilities. In other words, I want them to be upper-body animations, not full-body animations. I don't see forced stops serving any useful mechanical role. Sure, you eventually learn to bullet jump and aim glide through ability casts, which adds a tiny bit of skill to the gameplay, but that's a workaround for a clumsy mechanic.

Some of the more egregious offenders: Mesa and Mirage with their double buff refresh; Valkyr with _extremely_ long cast times on Warcry and Hysteria; Equinox with the _extremely_ long Maim cast time. On Equinox I run Natural Talent just to make it usable; perhaps without the movement restriction, it would be an option rather than a requirement. Perhaps the worst offender is Loki's Radial Disarm which both _stops_ your movement and forces a _step forward_; it's annoying when moving and annoying when standing.

I shouldn't have to explain the problems with the AWG self-stuns. I'm talking about both the stun on deployment, and the stun on shooting a fully charged Velocitus, Corvas, or alt-mode Larkspur. This is highly subjective, but does anyone enjoy frequently losing control of their character for 2-3 seconds? This also results in a very low fire rate and APM compared to normal guns, making them unviable for fast-paced missions. I simply don't understand why anyone would want this.

### Hard Landings

Warframe should take a cue from other 3rd person games and make hard landings _cosmetic_. It's about agile space ninjas who can roll, double jump, bullet jump, aim glide, wall latch, wall run, slide, and possibly more. Tenno are the most agile characters out of any game I've played. Personally, I find that forced hard landings mess up the vibe. Keep the animation there — but make it optional, like in so many other games (DMC, Nier Automata, there's plenty of examples). This means performing the recovery animation only if the player is not holding a movement key, and being able to cancel out of the animation by performing any action.

"Duh, just get used to it". I did. I'm very used to crouching or aim gliding to avoid a hard landing, and this does add a little bit of skill to the movement. But I'd still prefer it to be cosmetic, especially with things like Terra Trenchers (hard landing from 3m? really?).

To keep Kavat's Grace useful, its effect could be replaced with something like "+X% bullet jump velocity after landing at ≥ Y m/s".

### Gunblade Combos

A lot of combos mix shooting and meleeing and/or force movement. This is extremely impractical. I would prefer if aim combos were all-shoot and upper-body-only (compatible with arbitrary movement, or no movement), and non-aim combos were all-melee.

Gunblade shots can't be performed in midair. Heavy attack shots root you to the ground for the entire duration of the lo-o-o-ong-ass recovery animation. Part of that animation lets you run but not jump. This is extremely annoying. Gunblade shots should be converted to upper-body animations, doable while running _like they used to be_, and doable in midair.

Shots at close range often ragdoll the enemy away _instead_ of shooting it. This is extremely annoying. I'd like to have a word with whoever came up with this.

Gunblades somewhat encourage heavy attack builds, but if you have a gun equipped, there's no way to immediately perform a melee heavy attack. You have to either tap melee, or lock into melee, both of which take time, often too much time. The quick melee may ragdoll your target away which can be extremely inconvenient. If you were planning to heal with that heavy attack using Life Strike, the quick melee may kill the target, preventing you from healing off of it. Locking into melee is impractical when alternating gun and gunblade shots. This must be addressed by fixing the dedicated heavy attack hotkey which already exists but doesn't do what it should.

### Despawn Dispel

"Despawning" happens when you end up outside the allowed space, such as by falling into chasms. I would greatly appreciate if despawning didn't dispel warframe buffs. I understand that it's supposed to make chasms feel dangerous, forcing you to play carefully. But the current implementation disproportionately punishes some buff-dependent frames (Mesmer Skin Revenant, Assimilate Nyx, Chroma, etc.), while barely giving a slap on a wrist to others like Hildryn, Inaros, etc.

The mechanic also punishes the innocents: too many maps have bullshit despawn areas that look traversable. Some tiles in the Europa tileset have inviting open spaces that despawn you if you dare explore. In some places, simply _bullet jumping up_ can get you despawned. Fixing the maps might be impractical since it's hard to track down all such places, which makes for a solid argument in favor of removing the dispel.

In my opinion, stopping your movement, teleporting you, and forcing a recovery animation is punishment enough.

All that said, maps need to be fixed too. Some despawn zones are utter bullshit.

### Energy Leeches

Energy drain as a threat has its spice. It can be fun to play around magnetic door fields, EM clouds, Plains of Eidolon water, and so on. All of these are stationary, telegraphed, avoidable.

Energy leeches, however, drain your second most important resource by merely existing. They don't need line of sight, often in a different room. There's no visual or audio indicator, no telegraph, nothing. You have to constantly watch your energy UI to spot their presence, and that's not fun.

In general, I don't like the idea of enemies having an unavoidable effect on you through their mere _presence_, without any abilities or a LoS requirement. This leaves no place for interplay such as dodging or CC.

What needs to change: energy leeches should require line of sight and have a clear audio-visual effect, for example a fuzzy beam between them and the victim. The LoS requirement allows some interplay such as dashing into cover. We could go further and make it an ability they have to channel; the downside is having to implement special animations for a large number of enemies, lowering the uptime of the threat, and possibly stopping its movement, which all together might make it laughable.

To all veterans puffing at their monitor: "Just run ahead and kill all the things!", keep in mind that the game isn't balanced solely for fully-modded weapons and ultra-AoE abilities. At lower MR, especially in solo play, people don't instantly obliterate everything they see. Which is a good thing.

### Buff Refresh

By default, you should be able to refresh any buff before it expires. Coming from other games where this is standard, Warframe's approach is surprising and inconvenient. Some abilities allow this and some don't, which makes absolutely no sense.

Buff refresh is made doubly inconvenient by the [movement-stopping](#movement-stoppers) animations, which also need to be fixed.

### Energy Recovery while Channeling

Currently, energy recovery while channeling is extremely inconsistent. It disables energy recovery from _some_ sources, but not from others, on a case-by-case basis.

Allowed effects: Rage and Hunter Adrenaline, energy orbs, Arcane Energize, Orokin death orbs, etc.

Disabled effects: Energizing Dash, Energy Vampire, Thurible, Octavia's passive, Arcane Brave, consumables, etc.

Strange case-by-case rules:

* Desecrate and Peace & Provoke have no drain over time and allow all forms of energy recovery.
* Effigy and Mend & Maim have the standard channeling restrictions but make an exception to allow Energy Vampire.
* Artemis Bow has no drain over time but _still_ impedes energy recovery. It makes an exception to allow Energizing Dash.
* Spectral Scream has the standard channeling restrictions but makes an exception to allow Rift Plane.
* ...probably more.

This doesn't even work for power balance. Equinox's Maim and Titania's Razorwing are very strong, but Oberon's Renewal is very weak. Nyx's Absorb + Assimilate negates a lot of damage, but Revenant's Mesmer Skin negates infinite damage without impeding energy recovery or movement.

Renewal is probably the worst offender. You must frequently recast it to re-apply it to new targets or targets that have lost the buff, and then channel it, lest they lose that puny little heal over time. Given the loss of energy recovery from various effects such as Zenurik, EV, consumables etc., the energy cost is disproportionate to its small effect. This ability just screams "outdated" and needs to be converted from channeling to a timed buff.

Non-channeling frames can benefit from Energizing Dash and consumables, while channeling frames can't, and must resort to Arcane Energize and Primed Flow. When DE reworked Ember, they made statements such as "converting her 4 from channeled to non-channeled will help with energy recovery".

What's up with these special rules? Apparently, frames are not created equal: some are born with channeled abilities and punished for it. Then if someone at DE nags the rest about their special annoyance, they come up with special exceptions or warframe reworks. How about fixing this for everyone? Just drop the energy recovery restrictions. Arcane Energize already enables ability spam, but it locks one of your arcane slots and needlessly raises the entry bar on some frames.

## Bosses

### Profit-Taker Orb

I consider Profit-Taker an exemplary masterclass in bad boss design. It holds a ground-shattering record on packing the most annoyances into one encounter. Unfortunately I don't have a clear idea how to improve this fight, so just gonna rant to get this off my chest.

* Reinforcements are a constant stream. If you feel threatened and want to clear the field before focusing on the boss, tough luck, this tactic is not viable. It seems to have a lower bound on how much infantry is present, and _instantly_ replenishes them.
* Reinforcements are simultaneously an annoyance and something you're meant to ignore while focusing the boss. Both are backwards. Enemies should be something you engage actively. (Compare Vomvalysts which are a _resource_.)
* When _anyone_ comes close to the boss, _everyone_ is punished by bullshit energy-draining magnetic waves. Even when you're careful, some boss shots and infantry attacks knock you _toward_ the boss, forcing this. Note that the optimal anti-PT-shield weapons (Fulmin, Catchmoon, Redeemer, any melee) are close-range.
* Terra Trenchers stagger you with every hit. Their charge knocks you up, forcing a hard landing from a 3m height. They barely deal any damage. In a game about agile space ninjas, this enemy feels extremely insulting rather than threatening. (See also: secret lab amalgam MOAs on Jupiter.)
* Multiple boss abilities knock you down and are often impractical to avoid.
* Some MOAs knock you down with a grapple + kick. With so much happening at once, they're often impractical to avoid. This just punishes you for... existing?
* The fight forces the use of AWG which stuns you on deployment. Some stun you on every charged shot.
* AWG can't be deployed in the air, which is where you're _encouraged_ to be with all the ground-side threats.
* Boss shields appear instantly with no warning, denying the action you were about to perform. You don't know when to expect them, you take a position, prep a shot, it's instantly denied and you're knocked down for good measure. This happens over and over.
* On one of the arenas, you're likely to be knocked into water, losing buffs such as Vex Armor. Some of the pylons spawn near water, forcing you to go there.
* Enemy Fluctus hits push you. This can push you towards the boss, triggering magnetic bullshit, or push you away _while reviving someone_, which is normally supposed to make you CC-immune. This _even affects Void Mode_, making it unreasonably annoying to revive your trusty pet. This happens at the time when you're already threatened by someone going down, doubling down on the annoyance.
* ... and more

Given the above, I'm surprised the fight doesn't also feature melee nullifiers rushing at you to dispel that pesky Vex Armor!

"Just get used to it." I did. I have a Chroma loadout that makes Profit-Taker decently comfortable. It consists of anti-annoyance features: knockdown and stagger recovery, Kavat's Grace against bullshit hard landings from 3m, Arcane Nullifier against magnetic bullshit, and more.

If you personally don't find knockdown, stagger, self-stuns, hard landings, energy drain, instant reinforcements annoying, you might not understand why I'm complaining. It's subjective. It doesn't justify the current design. Someone who's not you being annoyed already makes it bad.

But wait, there's more! In addition to all the above, the fight is short and you're encouraged to run it multiple times, making it very repetitive with little room to change things up. And there's no real benefit from running it with a team.

### Exploiter Orb

Fewer annoyances than Profit-Taker, but still extremely bad.

* ... rant here

## Chat

Add the missing text editing shortcuts such as Ctrl+arrows, Ctrl+Backspace/Delete, Ctrl+A, Ctrl+X.

Preserve group chat history.

Add hotkeys for closing tabs.

Fix the bug where Tab and Shift+Tab are sometimes reversed.

Make all items linkable regardless if you have them.

Consider supporting text selection.

## Dojo

Allow us to clip decorations through anything.

Allow to bypass room destruction timers. No amount of timers will protect against a determined vandal.

Enable region chat. Bored veterans have two things to do: chat in region and decorate dojos. But you're not allowed both at once. That's Chaotic Stupid.

LR: Add _all_ decorations, especially Orokin, not a cherry-picked subset.

Leaving the polychrome should automatically activate the preview. Also, leaving one of the color menus should bring you back to the top-level polychrome menu.

## Laggy Mouse Cursor

On displays below 144 Hz, Warframe's software mouse cursor feels horribly laggy. I've probably played tens if not hundreds of cross-platform games which didn't have this issue. There's no excuse to keep it that way.

## Miscellaneous

* Remove the midair Transference-in limit of 1. It serves no discernible purpose and can be inconvenient on outdoor maps.
* Allow effects such as energy and ammo consumables to apply to transferred-out frames.
* Shorten the operator recovery animation from hard landings.
* Make weapon and gear choice stick through Transference, for real this time. Currently it's very inconsistent and depends on whether you're host or client.
* Remove door lockdowns.
* Provide a way to revive dead sentinels and pets in long missions, including arbitrations. This could be done via a craftable gear item.
* Replace the slider for global Ordis volume with a toggle for Ordis' Orbiter quips. It's useful to have Ordis on for quests. When new quests involving Ordis are added in the future, the current setting will degrade the quest experience for anyone who keeps it disabled because of the annoying Orbiter quips.
* Add hotkeys for toggling music and sound.
* Allow to bind hotkeys with modifiers, such as Ctrl+M for music.
* Add a hotkey for taking a screenshot without UI, replacing the global toggle we currently have.
* Remove the per-hit energy drain from Prowl or throttle it to something reasonable like no more than 10 per second.
* Troop transports (dropships) should be louder and/or heard from farther away. They should constantly emit noise instead of flying silently.
* Fix bullshit offhost staggers from Kuva Guardians and Guardsmen unaffected by Pain Threshold and stagger resistance.
* Remove bullshit energy drain from Guardsman block.
* Allow Deconstructor on non-Helios.
