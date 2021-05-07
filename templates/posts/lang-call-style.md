Note: also see `lang-calls.md`.

## Neutral Observations

The `(inside_parens)` calling convention places restrictions on delimiter-free operators. In Lisps, they're all special-cased:

```
-10
'one
```

If you want arbitrary unary operators to be called without delimiters, in this convention they can't _also_ be called with delimiters, because of the following ambiguity:

```
@one
(@ one two three)
(@one two three)
```

You can't tell if everything is passed to `@`, or if it's applied to `one` and the result is called with `two three`.

With the `outside_parens()` calling convention, arbitrary operators can be made callable with and without parens, unambiguously:

```
@one
@(one two three)
```

## _Objective_ advantages of `(inside_parens)` over `outside_parens()`.

* Both humans and compilers can easily tell where the expression begins and ends.

* The compiler doesn't need to group expressions with the following `()`.

* If your language doesn't have paren-free prefix or infix operators, it doesn't need any precedence or invisible expression grouping at all.

## _Objective_ advantages of `outside_parens()` over `(inside_parens)`.

* Works better with autocompletion in editors.

* In a typed language with generics, where constructors will often look like this: `Slice(String)("one" "two")`, this placement allows to elide the type without changing the "data" part of the constructor.

If your language is typed, constructors look like this:

```
List("one" "two" "three")
(List "one" "two" three)
```

This can be unwieldy for pattern matching and deconstruction, and for deeply nested structure literals:

```
let(Tuple(String Err)(val err) some_func())
(let ((Tuple String Err) val err) (some_func))

Slice(Slice(Slice(Int)))(
  Slice(Slice(Int))(Slice(Int)(10) Slice(Int)(20))
  Slice(Slice(Int))(Slice(Int)(30) Slice(Int)(40))
)

((Slice (Slice (Slice Int)))
  ((Slice (Slice Int)) ((Slice Int) 10) ((Slice Int) 20))
  ((Slice (Slice Int)) ((Slice Int) 30) ((Slice Int) 40))
)
```

The `outside_parens()` style allows you to elide the types without affecting the data part:

```
let((val err) some_func())

Slice(Slice(Slice(Int)))(
  ((10) (20))
  ((30) (40))
)
```

But doing this elision in the `(inside_parens)` style would turn values into functions. It would be very ambiguous and confusing. We probably can't do that.

Alternatively, you could come up with ludicrous hacks like this:

```
(let (_ val err) (some_func))

((Slice (Slice (Slice Int)))
  (_ (_ 10) (_ 20))
  (_ (_ 30) (_ 40))
)
```
