{{mdToToc .MdTpl}}

## Definitions

_Speedhack_: one of:

* Making the game as a whole faster. (Focus of this article.)
* Performing some actions much faster or sooner than normally allowed. (Out of scope for this article.)

_Slowhack_: inverse of speedhack. Typically makes the game as a whole run slower.

I recommend this **for single-player games only**. Do not cheat in multi-player games. Don't be a nuisance to your partners/opponents.

## Why

Many games are generally sluggish, or have sluggish segments such as unskippable cutscenes, forced walking sequences, and more. Whenever you feel bored and wish you could fast-forward, use speedhacks.

Some games, for some players, can be so challenging that the player is unable to progress without cheats. Slowhacking can make this easier in a way that doesn't circumvent any mechanics, merely compensating for reaction time. Slowhacking can be useful for learning and training, before upgrading to "proper" speed.

## How

There are probably many different ways. At the time of writing, I use and recommend _Cheat Engine_. Official site: https://cheatengine.org. CE "attaches" to another app, such as a game, to perform arbitrary hacks on it. CE has many features. Speedhacking is just one of them. I suggest reading CE docs/guides.

_**Kill Cheat Engine**_ before launching any multi-player game, otherwise you _**might get banned**_ from it. Many multi-player games have their own "cheat detection" which can produce false positives.

For convenience, setup global hotkeys in CE, which can be used without alt-tabbing. At the minimum, I suggest the following:

* Attach to current foreground process (example key: `numpad .`).
* Toggle speedhack (example key: `numpad 0`).
* Various speed increments, such as:
  * Speed: 0.5 (example key: `numpad 1`).
  * Speed: 1.5 (example key: `numpad 2`).
  * Speed: 2 (example key: `numpad 3`).
  * Speed: 4 (example key: `numpad 4`).
  * Speed: 6 (example key: `numpad 5`).
  * Speed increase (example key: `numpad +`).
  * Speed decrease (example key: `numpad -`).

After launching both CE and the target game (or another app you want to hack), hit the key to "attach to current process", use the appropriate speed key, and enjoy.

## Pausing

Speed can be zero! Cheat Engine lets you assign a hotkey to pause the selected process. Useful for pausing games which don't have their own pause built in. This also eliminates the process's CPU and GPU usage, allowing to keep temperature and fan noise down when alt-tabbing for prolonged periods.

## Gamepads

This section is out of scope for speedhacking, and might eventually be expanded into its own post.

When using a gamepad, it can be inconvenient to reach for the keyboard to use CE hotkeys, or other global hotkeys. This is fixable by emulating keyboard keys on the gamepad, for example via _Steam Input_ or _DS4Windows_.

Compared to KBM, gamepads have very few keys. You can't spare them for additional global hotkeys. However, you can find _combinations_ which are normally unused. For example, pressing or tilting the right analog stick and simultaneously pressing one of the face buttons. Or similar on the left side. Such combinations never occur in normal gameplay because a thumb can't be in two places at once. Some gamepads have additional less-useful buttons such as "mute", which can be recycled as modifiers. In Steam Input, pressing any button, or tilting an analog stick, or performing another action of your choice, can overlay a different configuration (called "_action set layer_") over existing keys, allowing you to configure a large number of additional actions, some of which can be used for CE and speedhacking.
