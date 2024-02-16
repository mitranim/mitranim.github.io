### Misc

* Conversions cost time. Having only _one_ way is better, no matter which.

* When typing multiple lines of single-line comments, you have to retype the comment character for each line. If the editor auto-inserts them, then you have to constantly delete them instead.

* Having a prefix on each line adds visual noise.

* Migrating text between comments and other places, such as Markdown documentation, is objectively easier without a prefix on each line.

* Single-line works better without syntax highlighting, or when colorblind.

* Single-line works better when code is viewed partially rather than fully, for example in text search.

* Multi-line comments cause arguments whether the body should be indented. Some prefer it indented, some unindented.

* Multi-line comments allow some code to become single line and then multi line without loss of information. Untrue for comments that use ASCII art, Markdown formatting, etc.

* Editing large single-line comments requires more micromanagement or special editor support for automatic insertion and deletion of comment prefix.

* Multi-line comments are easier to edit without special editor support, particularly when joining and breaking lines at the right width.

* Multi-line comments lend themselves to syntax hacks such as using a comment as a tag for syntax highlighting, for example `/* sql */` preceding an SQL string.

* Multi-line comments are more compatible with indentation shortcuts in editors: preserve indent when creating the next line, indent-unindent, bullets.

* Multi-line comments are more compatible with embedding code or pseudo-code, such as special annotations parsed by external tools.

* Multi-line comments, if nestable, allow to disable and enable chunks of code without affecting intermediary lines, and thus without affecting source control history for those lines.

### Inherent tradeoffs

For any large comment with multiple line breaks: single-line is easier to read without special support, harder to write without special support; multi-line is easier to write without special support, harder to read without special support.

### Observations

Popular, long-lived languages eventually acquire both comment types. Committing to just single-line or just multi-line comments always proves impossible in the long run. For example, Python officially has only single-line comments, but in practice, multi-line strings are used for large comments, mostly docstrings. Conversely, some people prefer having a comment prefix on every line, and when given only a multi-line comment format, they invent a prefix convention that emulates single-line comments inside of multi-line comments.

If a language supports both single-line and multi-line formats, it should use exactly one comment delimiter, to make it easier to convert between the formats. In such cases, the comment "fence" should be resizable (`|`, `||`, `|||` and so on), to avoid conflicts with the same character inside the comment.
