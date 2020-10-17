**TLDR**: always spaces, never tabs; 2 spaces rather than 4.

_Objective_ arguments in favor of spaces over tabs:

* Tabs break when copy-pasting. For example, terminals render tabs as spaces, usually as 8.
* Mixing tabs with spaces causes indentation to break in different editors. People _will_ mix them. Plenty of languages and editors don't have autoformatters. It will always stay this way.
* Distinguishing tabs from spaces requires special editor support.
* Having both is objectively more complicated.

_Objective_ arguments in favor of tabs:

* Configurable visual indentation level.
* Fewer characters and keystrokes. Not relevant in specialized editors; Tab and Backspace insert and delete multiple spaces.

_Objective_ arguments in favor of 2 spaces over 4 spaces:

* Fewer characters and keystrokes.
* Easier to type in non-specialized editors, such as chat input boxes, which don't have indentation shortcuts.
* Highly-nested code fits better on the screen. Relevant for markup such as XML.

_Relativity_: your preference is influenced by your display pixel density, resolution, OS, font family, font size, eyesight, and habits. Someone with a very large but low-DPI display is likely to prefer 4 spaces. Someone who writes code on a small display, in an IDE that uses 20% of the screen area for the actual code, is likely to prefer 2 spaces.

If you don't have a strong preference, go with 2 spaces.

## Sidenote

Some people use variable indentation. See the post [Use Fixed-Size Indentation](/posts/indent-fixed) on that.
