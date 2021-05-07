**TLDR**: Programming language syntax is over-specialized and constantly reinvented because of that. We would be better served with something more general that doesn't need to be remade.

One [possible solution](#solution) is included, with examples. Basically a variant of [M-expressions](https://en.wikipedia.org/wiki/M-expression) without any "sugar".

## TOC

{{tableOfContents .InputPath}}

## Problem

## Solution

Convert special syntax to "calls". It really is that simple.

```
fn some_func(a, b, c) {d; e}      | fn(some_func(a b c) d e)
some_func(10, 20, 30)             | some_func(10 20 30)
Some_list{10, 20, 30}             | Some_list(10 20 30)
Some_struct{one: 10, two: 20}     | Some_struct(one 10 two 20)
Some_dict{"one": 10, "two": 20}   | Some_dict("one" 10 "two" 20)
Dict<String, Int>                 | Dict(String Int)
-10 + 20 + 30                     | +(-10 20 30)
return 10                         | ret(10)
let val = 10                      | let(val 10)
let val: Int = 10                 | let(val Int 10)
coll[key0][key1]                  | at(coll key0 key1)
coll[key0][key1]                  | coll.get(key0).get(key1)
coll[key0][key1] = val            | set(at(coll key0 key1) val)
coll[key0][key1] = val            | coll.get(key0).set(key1 val)
```

Most of us won't give up `-10` or `one.two.three`. This can be solved generally:

* Ban operator mixing to avoid precedence.
* Allow unary operators without parens.
* Allow the concept of infix operators, but limit them to just `.`.
* Define the following precedence: infix > postfix `()` > prefix.

Result:

```
@one.two.three() ≡ @(.(one two three)())
```

**Show me the code!** Very well.

**NOT FINAL BY ANY MEANS**:

```
use(uuid)

fn(main()
  let(list Slice(Int)(10 20 30))
  let(head list.get(0))
  list.set(0 +(head 40))

  let(dict Dict(String Int)("one" 10 "two" 20 "three" 30))
  let(val dict.get("one"))
  dict.set("one" +(val 40))

  let(set Set(Int)(10 20 30))
  set.del(10)
  set.add(+(10 40))

  let(person Person(id uuid.new() name ""))
)

type(Person struct(
  id   Uuid
  name String
))

type(Tree(A) variant(
  Slice(Tree(A))
  A
))
```

This isn't new. This resembles M-expressions devised by the author of Lisp, more recently found in the Wolfram Language.
