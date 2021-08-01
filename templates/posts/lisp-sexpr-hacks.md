**TLDR**: nobody wants to write pure S-expressions, and Lisps are full of hacks around them.

Disclaimer: Lisps have decades of history and many dialects with a variety of hacks. The following is just what I happened to come across. There might be more.

Examples on this page use [Racket](https://racket-lang.org).

## Definitions

[S-expressions](https://en.wikipedia.org/wiki/S-expression) is a syntax for binary trees. The base notation has only atoms, pairs, and nil:

```
symbol       |    atom
"string"     |    atom
10           |    atom
(10 . 20)    |    pair
()           |    nil
```

The "abbreviated" notation omits `.` from pairs that end with another pair or nil, combining them into lists:

```
(10)            ->    (10 . ())
(10 20)         ->    (10 . (20 . ()))
(10 20 30)      ->    (10 . (20 . (30 . ())))
(10 20 . 30)    ->    (10 . (20 . 30))
```

When talking about S-expressions as code, we usually mean the abbreviated notation, as in Lisps. Writing code in the base notation is out of the question, but pairs will come back to haunt us later. Example Lisp code:

```rkt
(define add (lambda (a b) (+ a b)))

(define some_var (add 10 20))
```

## Why

* Extremely simple.
* Can express any computation.
* Infinitely extensible.

We can express new concepts by adding meaning to symbols such as `lambda`, `if`, and so on. Each such "form" will have its internal "syntax", usually extremely simple, but we don't have to change the base notation. The cost of adding and learning new features is lower compared to other syntaxes. This also makes it easy to give _users_ the ability to extend it, via AST-based macros.

Sidenote. Personally I like the S-expression syntax, but advocate against dynamic typing and [homoiconity](/posts/lang-homoiconic) as seen in Lisps. We could and should use S-expressions for statically typed languages.

## Hacks

### Number Literals

S-expressions require unary negation to be written like this:

```
(- num)
(- 10)
```

But `-10` was too hard to give up, so they built `+-` _into number literals_. The language's parser supports `+10` `-10` where the operator is part of the number's syntax. Note that `+ 10` `- 10` (with a space) don't work that way. Of course, this limited special case works _only_ for literal numbers, not variables, and doesn't extend to other unary operators such as bitwise negation.

### Prefix Operators

Despite claiming the opposite, Lisps have always had many prefix operators, not just `-10`.

Lisps have a concept of "quoting" code. Because the code notation _happens_ to be a data notation, the quoted code can be evaluated as data. This also serves as the language's AST, used internally.

```rkt
; Evaluate as code, result is `30`
(add 10 20)

; Evaluate as data, result is `(add 10 20)`
(quote (add 10 20))
```

Writing `(quote)` and others was too much, so they added prefix shortcuts.

```rkt
'(add 10 20)       ->    (quote (add 10 20))
`(add 10 20)       ->    (quasiquote (add 10 20))
`(add 10 ,expr)    ->    (quasiquote (add 10 (unquote expr)))
`(add ,@exprs)     ->    (quasiquote (add (unquote-splicing exprs)))
```

In general, all Lisp prefix operators are aliases for "expanded" forms. They're converted after or during parsing text into AST. Parsing text and converting prefix operators is combined into a step called "reading", which returns a canonical AST.

[Clojure's reader](https://clojure.org/reference/reader) has more prefix operators, such as `@A` → `(deref A)`, and a somewhat-generalized `#`.

Upside: because this is done once at "read time", no other code has to deal with prefix operators. Downside: standard library and user code either can't define new prefix operators, or must use an API different from functions and macros.

### Curly Infix

People have written large documents and reference implementations suggesting `{}` for infix. See [SRFI 105](https://srfi.schemers.org/srfi-105/srfi-105.html). Code inside `{}` would be implicitly and unambiguously converted to the canonical form by the reader.

```
{10 + 20 + 30}      ->    (+ 10 20 30)
{{10 + 20} * 30}    ->    (* (+ 10 20) 30)
```

Veiled in-joke or serious request? Can't tell...

It can be observed that this proposal has grouping, but no precedence. Grouping is both necessary and sufficient. Precedence is not necessary and not sufficient. Programming languages have lots of operators that don't exist in math, and their precedence is inconsistent between languages. Precedence errors are so insidious that some languages, like Pony, ban most forms of operator mixing and enforce grouping. This proposal, while ludicrous in the context of Lisp, has at least one good idea at its core.

### Racket Infix Hack

Racket has a special infix hack.

Remember the unabbreviated `(a . b)` syntax for pairs? Racket folks have found unused "dead space" in the syntax they could exploit. In addition to binary `(a . b)` which makes a _pair_, it supports ternary `(a . b . c)` which makes a _reordered list_. They use _one_ infix operator to enable _other_ infix operators or functions in a "general" way.

```
(10 . + . 20)               ->    (+ 10 20)
((10 . + . 20) . * . 30)    ->    (* (+ 10 20) 30)
```

It's often said that forbidden fruit is desired more strongly. Evidence suggests that when Lisp bereaves its users of infix, they develop a strong desire for more, _more_ infix! (We herd you like infix, so we put more infix in your infix...)

### Namespacing in Symbols

Most languages have some form of namespacing. Some mix several forms.

```
one.two.three
one->two->three
one:two:three
one::two::three
one/two.three
```

Since inception, Lisps have allowed special characters inside symbols, and avoided infix operators. It naturally followed that Lisp package systems implement namespacing inside symbols. Common Lisp and Racket use `:`, Clojure uses `/` and `.`.

```
package:identifier
namespace/identifier
value.method
```

Still a hack, because _useful_ applications of these symbols involve sub-parsing them. Conceptually, these are separate identifiers combined by an infix operator. The parser (or "reader") should have parsed them for you, storing the pieces in the AST. That's what Clojure does: its symbols are classes with separate "namespace" and "name" parts. This indicates that they were combined prematurely.

Sidenote. One simple alternative is to extend "reader macros" by supporting infix `:`, converting `one:two:three` to canonical `:(one two three)`. Lisps already special-case `.` in a similar way; `:` would have a higher precedence. As long as there's no other infix, this should parse unambiguously. Alternatively, we could ditch the pair syntax and use `.` for namespacing. Improper pairs could be printed as `(cons a b)`.

The major downside of the solution above, aside from added complexity, is that it's non-extensible, as adding more infix would create parsing ambiguities, which we can't resolve because we can't afford `()` for grouping. I would appreciate a simple and flexible approach that doesn't seem hacky.

## Conclusion

If Lisp people didn't stick with pure S-expressions, nobody will. Languages designed for practical use must include common prefix and infix shortcuts. To me, everything above seems hacky or complicated. Elegant approaches are topics for other posts.
