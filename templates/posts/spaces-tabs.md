**TLDR**: always spaces, never tabs; 2 spaces rather than 4.

## Arguments

_Objective_ arguments in favor of spaces over tabs:

* Tabs break when copy-pasting. For example, terminals render tabs as spaces, usually as 8 (at least mine).
* Mixing tabs with spaces causes indentation to break in different editors. People _will_ mix them. Plenty of languages and editors don't have autoformatters. It will always stay this way.
* Distinguishing tabs from spaces requires special editor support. (Code editors can render special harder-to-see symbols for whitespace.)
* Spaces are both necessary and sufficient. Adding tabs is adding complexity.

_Objective_ arguments in favor of tabs:

* Configurable visual indentation level.
* Fewer characters and keystrokes. Not relevant in specialized editors; Tab and Backspace insert and delete multiple spaces.

_Objective_ arguments in favor of 2 spaces over 4 spaces:

* Fewer characters and keystrokes.
* Easier to type in non-specialized editors, such as chat input boxes, which don't have indentation shortcuts.
* Highly-nested code fits better on the screen. Relevant for markup such as XML.

_Objective_ arguments in favor of 2 spaces over 1 space:

* Easier to distinguish from line wrapping. In some editors, when line wrapping is enabled, the secondary lines are intended by 1. With 2-space indentation, you can tell them apart.

## Relativity

Your preference is influenced by your display pixel density, resolution, OS, font family, font size, eyesight, and habits. Someone with a very large but low-DPI display is likely to prefer 4 spaces. Someone who writes code on a small display, in an IDE that uses 20% of the screen area for the actual code, is likely to prefer 2 spaces.

If you don't have a strong preference, 2 spaces seems like a better default, based on the arguments above.

<!--
## Variable Indentation

Some people use variable indentation. See the post [Use Fixed-Size Indentation](/posts/indent-fixed) on that.
-->
