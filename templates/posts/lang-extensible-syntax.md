## Why

The idea of "extensible syntax" might be foreign to some programmers, so I have to explain it first. In short, it's the ability to write code that looks like the language was designed for that particular program.

Closely related to [bottom-up programming](https://www.paulgraham.com/progbot.html): making smaller building blocks, combining them into larger blocks, then into larger, and so on. A stone castle made of blocks, rather than a sand castle made of assembly.

Bottom-up programming and shaping language syntax for the given problem are both essential foundations of software. It's why we don't usually write assembly, why languages and DSLs exist, why standard libraries exist, why frameworks exist. Segregating common functionality into reusable functions is widely taught and should be a no-brainer to everyone.

Some programmers will take whichever building blocks you gave them, and write the rest of the program in the lowest level "assembly" they know. Extensible syntax isn't for them. It's for those who write bottom-up, make libraries and frameworks, and want programs they can actually maintain. Read on!

## Extensibility

A well-designed syntax would be infinitely extensible for both language developers and language users. Every feature costs design space. Poorly-designed syntax features create conflicts, making it harder to add more features.

For a syntax to be _infinitely_ extensible, it must also be _fixed_, in the sense of having few features that can express anything. Instead of adding new syntax at the base level, we add new "forms" expressed in the base syntax. This automatically has a _cost_, or rather a tradeoff: code appears very uniform.

We have a working example: the Lisp family, or more specifically, S-expressions. It's not the only syntax that meets our requirements. I'll propose something else below.

```
(anything
  (anything anything)
  (anything (anything (10) "20")
    anything
  )
)
```

With only a few elements (symbols, numbers, strings, and lists), this can express everything. Special forms, such as conditionals, define their own _sub-syntax_ expressed _in terms_ of the base syntax:

```
(if <predicate> <then_branch> <else_branch>)

(case <predicate0> <branch0> <predicate1> <branch1> <else_branch>)
```

Worth noting that Lisps don't use _just_ S-expressions. They also have hacks for prefix and infix operators:

```lisp
; Prefix operators embedded in number literals.
-10
+10

; Prefix operators evaluated at read time.
'sym ; (quote sym)
@sym ; (deref sym)

; Namespacing embedded in symbols.
one:two:three
one/two.three
```

Lisps didn't _actually_ stick with pure S-expressions. From the start, they had prefix operators such as `'`.

 Common Lisp has package namespacing embedded in symbols: `one:two:three`. Clojure has _two_ namespacing syntaxes embedded in symbols: `one/two.three`. Those are basically hacks for infix operators. They have the special prefix operators `'` `,`. Clojure also has the prefix operator `@`. They're evaluated at "read time", but they're still prefix operators. They also have `+10 -10` where `+ -` is embedded in the number literal. Another prefix operator hack.

Belated realization that Lisps didn't _fully_ do what they claimed to do, with their syntax. Namely, it's not pure S-expressions. They've _always_ had prefix operators, and later on, hacks for infix operators. They _did_ avoid precedence problems. Precedence doesn't arise when your calls are `(a b c)` and you have only prefix operators. By embedding infix inside symbols, they dodged the problem.

Precedence problems arise in Mox is because it has `a(b c)` calls _and_ real infix operators. Between the trifecta of: prefix operators; postfix calls; infix operators; **any two** of them create precedence problems.

So why `-10` must be a prefix operator rather than part of the number literal? First, numbers support more than just `+ -`; there are other unary operations, like bitwise negation. It's an open set, and we don't want to special case them in the parser. Second, it's required for a typed language that allows to overload operators and value literals. For example, you should be able to define a big number type and declare a value `-10`, where `10` invokes the constructor and `-` invokes negation.

## Layering

We must define several layers of the syntax.

* Base level: data literals, identifiers, delimiters.
* Language level: semantics.

## Misc

Being able to elide everything not essential to the problem you're solving. This requires either type inference or dynamic typing, among other things.
