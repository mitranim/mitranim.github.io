
* Conversions cost time. Having only _one_ way is better, no matter which.

* When typing multiple lines of single-line comments, you have to retype the comment character for each line. If the editor auto-inserts them, then you have to constantly delete them instead.

* Having a prefix on each line adds visual noise.

* Migrating text between comments and other places, such as Markdown documentation, is objectively easier without a prefix on each line.

* Single-line works better without syntax highlighting, or when colorblind.

* Multiline comments cause arguments whether the body should be indented. Some prefer it indented, some unindented.

* Multiline comments allow some code to become single line and then multi line without loss of information. Untrue for comments that use ASCII art, Markdown formatting, etc.

* Editing large single-line comments requires more micromanagement or special editor support for automatic insertion and deletion of comment prefix.

* Multi-line comments are easier to edit without special editor support, particularly when joining and breaking lines at the right width.

* Multi-line comments lend themselves to syntax hacks such as using a comment as a tag for syntax highlighting, for example `/* sql */` preceding an SQL string.

* Multi-line comments are more compatible with indentation shortcuts in editors: preserve indent when creating the next line, indent-unindent, bullets.

* Multi-line comments are more compatible with embedding code or pseudo-code, such as special annotations parsed by external tools.

Inherent tradeoffs:

* For any large comment with multiple line breaks: single-line is easier to read without special support, harder to write without special support; multi-line is easier to write without special support, harder to read without special support.
