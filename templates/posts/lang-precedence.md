**TLDR**: Programming languages should avoid per-operator precedence and enforce grouping.

## Precedent (examples from existing languages)

Precedent: Lisps.

```
(+ 10 (* 20 -30))
```

Precedent: Pony. See https://tutorial.ponylang.io/expressions/ops.html#precedence.

```
10 + (20 * -30)
```

## Precedence

These are probably familiar to you:

```
+ - * /
```

Traditionally, math is written as tersely as possible, with single-letter variables, precedence rules over explicit grouping, and other shortcuts. `*`, or rather `✕` or `·`, is usually omitted.

People learn this in school and think this should translate 1-1 into code. Only seems natural, right?

Then, programming languages add more operators, like 3-5 times as many, for which precedence rules are not taught in school, and which vary between languages to boot. I can't believe people vouch _for a way to make errors_ when we could solve this forever and move on.

Programming languages should ban operator mixing and enforce grouping. There should be _no_ precedence rules.

## Operators and Precedence

**The bad**. I don't want to _ever_ read code like this:

```
10 + 20 * 30 - 40 / 50
@one.two.three()
```

**The good**. This is unambiguous:

```
-(*(+(10 20) 30) /(40 50))
@(one.two.three())
```

**Rules**.

Every operator, regardless of arity, has the canonical "call" form:

```
+(10)
+(10 20 30)
-(10)
-(10 20 30)
#(val)
@(ptr)
.(one two)
.(one two three)
```

Every unary operator can be used without parens, which is implicitly treated like the canonical "call" form:

```
+10                    ->     +(10)
-10                    ->     -(10)
#val                   ->     #(val)
@ref                   ->     @(ref)
```

Infix is implicitly treated like the canonical "call" form. Either variadic or left-associative binary. Some operators are inherently variadic.

```
10 + 20 + 30           ->     +(10 20 30)       |  +(+(10 20) 30)
one.two.three          ->     .(one two three)
```

Operators can be mixed with the postfix "call" form, which has the lowest precedence:

```
one.two.three()        ->     .(one two three)()
@ptr()                 ->     @(ptr)()
one + two + three()    ->     +(one two three)()
```

Operator mixing is _banned_ and must be resolved using the canonical "call" form:

```
invalid                ->     valid
                       ->
-10 + 20 + 30          ->     -(10) + 20 + 30
-10 + 20 + 30          ->     +(-10 20 30)
^ -10                  ->     ^(-10)
- ^10                  ->     -(^10)
- @ptr                 ->     -(@ptr)
@one.two.three         ->     @(one).two.three
@one.two.three         ->     @(one.two.three)
```
