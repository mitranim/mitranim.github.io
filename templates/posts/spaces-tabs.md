**TLDR**: always spaces, never tabs; 2 spaces rather than 4.

## Arguments

_Objective_ arguments in favor of spaces over tabs:

* Spaces are both necessary and sufficient. Tabs are not necessary and not sufficient (imagine separating words with tabs). Adding them adds complexity.
* Tabs break when copy-pasting. For example, some terminals and websites render tabs as spaces, usually as 8.
* Distinguishing tabs from spaces requires special editor support. (Rendering special whitespace symbols.)
* Mixing tabs with spaces causes indentation to break in different editors. People _will_ mix them. Plenty of languages and editors don't have autoformatters. It will always stay this way.
* Many environments, such as browsers, use Tab for navigation and don't support insertion of the tab character.

_Objective_ arguments in favor of tabs:

* Configurable visual indentation level.
* Fewer characters. Sometimes fewer keystrokes. (Note: in decent code editors, using spaces takes just as few keystrokes.)

_Objective_ arguments against tabs:

* When copy-pasting code, tabs often get converted to spaces, causing breakage and busywork.
* When copy-pasting code with tabs, some REPLs interpret tabs as completion requests rather than indentation.

_Objective_ arguments in favor of 2 spaces over 4 spaces:

* Fewer characters and keystrokes.
* Easier to type in non-specialized editors, such as chat input boxes, which don't have indentation shortcuts.
* Highly-nested code fits better on the screen. Relevant for markup such as XML. Relevant for side-by-side file viewing and diffing.

_Objective_ arguments in favor of 2 spaces over 1 space:

* Can distinguish from line wrapping. In some editors, when line wrapping is enabled, the secondary lines are intended by 1. With 2-space indentation, you can tell them apart.

## Bias

Your preference is influenced by your display pixel density, resolution, OS, font family, font size, eyesight, and habits. Someone with a very large but low-DPI display is likely to prefer 4 spaces. Someone who writes code on a small display, in an IDE that uses 20% of the screen area for the actual code, is likely to prefer 2 spaces.

If you don't have a strong preference, 2 spaces seems like a better default, based on the arguments above.

<!--
## Variable Indentation

Some people use variable indentation. See the post [Use Fixed-Size Indentation](/posts/indent-fixed) on that.
-->
