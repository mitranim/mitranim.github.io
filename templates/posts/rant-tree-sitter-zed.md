## Preface

I write code and prose in [Sublime Text](https://sublimetext.com). Over the years, I ended up writing many syntaxes for it, some of which (reworks of Clojure and Go) got officially adopted. This required understanding ST's syntax engine, which is simultaneously simple and special.

Sublime's syntax engine uses a pushdown automaton (PDA), which pushes and pops "contexts". Every context matches some text, gives it a "scope", and optionally modifies the stack by pushing another context (or several) or popping when done. Matching is regex-based, using Sublime's own optimized regex engine. You can assign different scopes to different captures in a single match, and contexts can re-match an earlier match via `\1`. This engine can handle any context-free grammar.

A syntax is a YAML file which defines some metadata (how to detect files), some variables (for reuse), and a bunch of contexts. Writing one is very fluid: the editor detects changes and reloads the syntax on the fly. Furthermore, syntaxes can _inherit_ from each other, with additions and overrides. This allows easy reuse in syntax "families" such as SQL, and the user can easily customize built-in syntaxes by inheriting and adding their own features. I use this for embedding, such as SQL-in-Go.

Sublime also uses syntaxes for symbol indexing. When you open a directory, it runs the syntax engine on every file and collects the symbols. This is fast enough to index large repositories in seconds. Sublime doesn't even bother to cache symbol indexes to disk. Result? Global symbol search and goto any symbol from anywhere, without an LSP, and goto works across languages (important for me).

Strangely, no other major editor implements this.

(Edit: it's been pointed out to me that Sublime _does_ cache symbol index to disk. It appears that something on my system or in my configuration prevents it from being _used_. This simply strengthens the point above, that syntax-based symbol indexing is very practical, since even without caching, it works great for me.)

## Rant

Been playing with the [Zed editor](https://zed.dev). It's very nice, very snappy, behind Sublime in some ways, ahead in some other ways. Like every editor which isn't Sublime, it's fatally defective (for me) due to the [lack](https://github.com/zed-industries/zed/issues/13307) of syntax-based symbol indexing and goto. It relies on LSPs, which don't exist for many languages and don't support cross-language goto. Like many recent editors, it uses [Tree-Sitter](https://tree-sitter.github.io) for syntaxes.

Since I'm working on a new language and have a Sublime syntax for it, I dabbled into making a Tree-Sitter syntax. And what a horror it turned out to be. Install an NPM lib that runs in Node to execute JavaScript to generate C then shell out to a C compiler... What?? You have to run its special CLI to initialize a project and it splurges out binding templates for six languages, which is still not enough, because every language which wants to use this stuff needs its own bindings for _every syntax_ so you get `N*M` integration complexity... What??? When making changes to a syntax, I have to rerun two separate CLI commands just to build this stuff, and another one to test it... What??!!! How _on earth_ did they manage to get this so wrong?!

Sublime's way is so obvious: write _one_ implementation of an automaton that interprets a syntax file described with data; make it a library and write bindings for all languages just _once_; keep syntax files portable, and watch and reload them without having to shell out to a C compiler!

Oh, and after _all that_, T-S _still_ **can't handle heredocs** (variable-length delimiters) which exist in shell, SQL, Markdown, and my own languages, so you **_have to write custom C code_** for what's done with `\1` in Sublime syntax.

I don't understand how they managed to push this garbage approach so hard without people revolting. Rant over.

## PS

I'm not attacking people; I'm attacking bad ideas.

Symbol indexing is an offtopic here. T-S can obviously be used for it. But as a syntax author, I never want to deal with it.

We could still fix this. All we need is one open source implementation of Sublime's engine (convince Sublime HQ to release theirs?), and then other editors could switch to its much less painful approach (from syntax author perspective).

The PDA approach is obviously not perfect. It requires the author to _think_ like a PDA, which is easy enough if you routinely program in assembly or Forth, but otherwise requires training.

However, we have an amazing high level solution: [SBNF](https://github.com/BenjaminSchaaf/sbnf), which produces Sublime syntax files from sources written in a specialized variant of the Backus-Naur form. Those [sources](https://github.com/BenjaminSchaaf/sbnf/blob/master/sbnf/sbnf.sbnf) are just as high level as Tree-Sitter, but can also handle heredocs without custom C code. Currently this requires syntax authors to run SBNF and distribute YAML. But SBNF, or something very similar, could be integrated into editors. Syntax authors would just write SBNF files, and editors would interpret them. No JS→C toolchain needed.
