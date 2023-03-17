**TLDR**: Programming language syntax tends to be over-specialized. Sticking to only one delimiter such as `()` helps avoid this.

For simplicity, examples in this post use `()`, but the choice is arbitrary. It can be `{}`, `[]` or other, if you prefer.

## TOC

{{mdToToc .MdTpl}}

## Delimiters

Some modern languages use _all_ of these delimiters, while overloading every single one.

```
() [] {} <>
```

Example from Rust:

```
() func def / parameters    |    fn some_func()
() func call / arguments    |    some_func()
() grouping                 |    one && (two || three)
() tuple type               |    (i16, i32, i64)
() tuple literal            |    (10, 20, 30)
[] indexing                 |    value[key]
[] array type               |    [i32; 11]
[] slice type               |    [i32]
[] array literal            |    [10, 20, 30]
{} block                    |    {statement; statement; statement}
{} struct literal           |    SomeStruct{field: "value"}
<> generics                 |    TypeA<TypeB<TypeC>>
<> operators                |    < > << >>
```

Little-known fact: you could do everything with one pair:

```
()
```

We can support all these and more, by simply using the `call()` form:

```
() func def / parameters    |    fn(some_func())
() func call / arguments    |    some_func()
() grouping                 |    and(one or(two three))
() tuple type               |    Tuple(Int16 Int32 Int64)
() tuple literal            |    SomeTuple(10 20 30)
() indexing                 |    at(value key) or value.get(key)
() array type               |    Array(11 Int32)
() slice type               |    Slice(Int32)
() array literal            |    Some_array(10 20 30)
() block                    |    do(statement statement statement)
() struct literal           |    SomeStruct(field "value")
() generics                 |    TypeA(TypeB(TypeC))
```

## Issues

Multiple delimiters would make sense if they weren't overloaded. If each pair had exactly _one_ fixed role, you could tell on a glance what it does. For example:

```
fn add(a, b) {ret a + b}
add(10, 20)

{
  var list = [10, 20, 30]
  var value = list.get(0) // No indexing syntax.
}
```

But this already can't express dicts without overloading. The smallest amount of overloading would be to use `[:]` for dicts, like Swift. And we're still ignoring generics.

```
var dict = ["one": 10, "two": 20]
```

Next thing you know, people want set literals. If lists and dicts deserve their own literals, why not sets? In a typed language, we could overload list literals by relying on type inference:

```
let set: Set = [10, 20, 30]
```

But if our language is dynamic, or people want full inference, set literals require their own syntax. Or, we just give up and use constructors.

```
let set = #[10, 20, 30]
let set = new Set(10, 20, 30)
```

At this point, the language designer has difficulties cramming new requirements into the original syntax because it was _over-specialized_.

## Consistency

When first learning JS, I was confused by the differences between `{}` for blocks and `{}` for objects.

```
{
  var one = 10;
  var two = 20;
}
{
  one: 10,
  two: 20
}
```

Both of these define a "scope" with "variables" within `{}`. But a block uses `=;`, while an object uses `:,`. What?! This was a headscratcher {{Emoji "🤯" ""}}. They could have used the same syntax!

```
{
  one = 10
  two = 20
}
var dict = {
  one = 10
  two = 20
}
```

Furthermore, a dict literal could be just a block, executed sequentially, and you would get the _scope object_, where variables become properties:

```
var dict = {
  var one = 10
  var two = one + 20
}
```

Why not? Dynamic languages can do this easily (Python does). What's weird is that _it wasn't done this way_.

C reduces the inconsistency by using `=` for struct fields, but the mandatory `.` and `,` spoil the fun:

```c
struct {int x; int y;} point = {.x = 10, .y = 20};
```

Python _almost_ allows consistent declaration syntax, but the mandatory `,` spoils the fun:

```py
one = 10
two = 20
val = dict(
  one = 10,
  two = 20,
)
```

Learned programmers handle this easily, but it was an obstacle for me, and probably many other novices. There are similar consistency issues in `[]`, and so on.

On the other hand, a simpler, punctuation-free syntax presents fewer problems for learning.

```lisp
(do
  (var one 10)
  (var two 20)
)
(var dict (:
  one 10
  two 20
))
```

## Precedents



## Parameters

Many modern languages enclose function parameters in parens. This doesn't need any changes.

```
fn some_func(a, b, c) {}
fn(some_func(a b c))
```

## Calling

Many modern languages use parens for calls. This doesn't need any changes.

```
some_func(a, b, c)
some_func(a b c)
(some_func a b c)
```

## Precedence

Many modern languages allow to control precedence by enclosing expressions in parens.

```
10 * (20 + 30)
```

This can be disambiguated from other parts of the language by punctuation. Go allows the following, by implicitly inserting `;` after `fmt.Println()`:

```
func main() {
  fmt.Println()
  ("one" + Printer("two")).Println()
}

type Printer string

func (self Printer) Println() { fmt.Println(self) }
```

But the following wouldn't compile:

```
fmt.Println()("one" + Printer("two")).Println()
```

If our language doesn't rely on newlines or punctuation, we can't afford bare `()` for grouping because these two examples would be indistinguishable to the compiler.

One simple solution is to just "call" operators, which automatically groups the expressions:

```
*(10 +(20 30))
(* 10 (+ 20 30))
```

## Indexing

Many modern languages use `[]` for indexing lists and dicts. This is used for both reading and writing.

```
val = some_list[10]
val = some_list[expression_as_index()]
some_list[index] = val

val = some_dict["literal_key"]
val = some_dict[expression_as_key()]
some_dict[key] = val
```

We could instead treat "calls" of lists and dicts as indexing. This would be unambiguous in any language where functions can't be collections, and collections can't be functions:

```
val = some_list(10)
val = some_list(expression_as_index())
some_list(index) = val

val = some_dict("literal_key")
val = some_dict(expression_as_key())
some_dict(key) = val
```

Some languages require distinct forms for indexing and calling. For example, in JS functions can also be treated as dicts. In Python you can define `__call__` and `__getitem__` on the same object.

We could define a pseudo-function understood by the compiler:

```
val = at(some_list 10)
val = at(some_list expression_as_index())
at(some_list index) = val

val = at(some_dict "literal_key")
val = at(some_dict expression_as_key())
at(some_dict key) = val
```

Or support this special syntax:

```
val = some_list.(10)
val = some_list.(expression_as_index())
some_list.(index) = val

val = some_dict.("literal_key")
val = some_dict.(expression_as_key())
some_dict.(key) = val
```

Finally, we could skip syntax and use methods. They could be predefined on built-in collections. They could compile to the exact same assembly instructions:

```
val = some_list.get(10)
val = some_list.get(expression_as_index())
some_list.set(index val)

val = some_dict.get("literal_key")
val = some_dict.get(expression_as_key())
some_dict.set(key val)
```

## Blocks

Just about every language since the 60s uses blocks and lexical scoping. Many modern languages denote blocks with `{}`. (Some use `begin end` or indentation, but we'll focus on `{}` for brevity.)

Blocks are often required in function definitions and various control flow constructs.

```
fn add(a, b) { return a + b }
if a { b } else { c }
loop { break }
```

Standalone `{}` creates a lexical sub-scope. In some languages, `{}` can be an expression.

```
{
  let val: Int = 10;
  {
    let val: String = "20"; // Unused variable.
  }
  val += {
    let _ = 40; // Unused.
    50
  };
}
```

We could just as easily use `()`. In function definitions and control flow constructs, it just works:

```
fn add(a, b) (return a + b)
fn(add(a b) +(a b))
```

For standalone blocks, we can use a pseudo-function:

```
val = do(
  _ = 10     // Unused.
  _ = 20     // Unused.
  30         // Actual result.
)
```

## Generics

Generics have multiple parts.

Defining and instantiating a generic type.

```
struct Vec<T> {
  ptr: *T,
  len: usize,
  cap: usize,
}

let vec: Vec<String> = ...;
```

Defining and instantiating a generic function.

```
fn add<T: Add>(a: T, b: T) -> T::Output { a + b }
add::<i64>(10, 20)
```

This could be converted to `()` without ambiguity with the other parts of the language.

```
type(
  Vec(T)
  struct(
    ptr *T
    len Usize
    cap Usize
  )
)

let(vec Vec(String) ...)
```

```
fn(
  add(a T b T) Ret(T.Output) where(T Add)
  a + b
)

let(num Int64 add(10 20))

cast(Int64 add(10 20))
```

## Constructors

---

`()` can be used for constructors. Precedent: Swift, Python.

```swift
struct SomeType {
  let one: Int
  let two: String
}

let someVal = SomeType(one: 10, two: "20")
```

```py
class SomeType:
  def __init__(self, one, two):
    self.one = one
    self.two = two

some_val = SomeType(one = 10, two = "20")
```

## Delimiter choice

While you could technically use any pair of characters, the most sensible choice on ASCII layouts is one of: `() [] {} <>`.

_Objective_ argument that excludes `<>`: these characters are used for math operators, while the other delimiters are not.

_Objective_ argument in favor of `()` over the others: easier to write by hand.

Unless better arguments are presented, `()` wins.
