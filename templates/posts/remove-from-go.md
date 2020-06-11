The [Go programming language](https://golang.org) espouses ["less is more"](https://commandcenter.blogspot.com/2012/06/less-is-exponentially-more.html). It prefers fewer features and "one way of doing things". However, it still has some fat to lose! This article highlights what I consider unnecessary, and suggests the path to gradual deprecation and removal.

Goes without saying: **this is an opinion piece**. If we disagree, that's cool!

This is just what I consider _relatively easy to remove_. I have other complaints about Go, mostly related to its deep fundamentals that would be very hard or impossible to change. They're not mentioned in this piece.

We're not allowed to break existing code under Go1. However, it seems plausible to migrate most existing code in advance, preparing it for the hypothetical Go2 that removes the deprecated features, alongside other breaking changes it's expected to make. The following migration strategy seems realistic:

* Go1 adds two minor syntactic features (see below)
* A tool like `go fix` converts existing code to the "new" style, avoiding the "deprecated" features
* Both "old" and "new" code continues to run under Go1
* Go2 drops the unnecessary features

## TOC

* [Language](#language)
    * [Remove `:=` in favor of `var`](#prefer-var)
    * [Remove parenthesized lists](#remove-paren-lists)
    * [Maybe remove `iota`](#remove-iota)
    * [Remove `new` in favor of `&`](#remove-new)
    * [Remove dot-import](#remove-dot-import)
    * [Remove if-assignment](#remove-if-assignment)
    * [Remove short float syntax](#remove-short-float-syntax)
* [Tools](#tools)
    * [Gofmt: align adjacent assignments](#gofmt-declarations)
* [Misc](#misc)

## Language Changes {#language}

### Remove `:=` in favor of `var` {#prefer-var}

#### Arguments

1\. Having two equivalent assignment forms is redundant.

2\. `:=` can't justify itself with brevity. Compared to `var`, it requires one or two fewer keystrokes to type, but involves `Shift` and an awkward movement between `:` and `=`. Subjectively, I find `var` easier and faster to type.

3\. Code sometimes needs to be converted between `:=`, `var` and `const`. For example, you have a string that's initially produced by `fmt.Sprintf`, but as you edit the code, it becomes a `const`. Or vice versa. I find these conversions fiddly and awkward. Converting between `var` and `const` is noticeably easier.

Moving a declaration between local and global scopes also involves converting between `:=` and `var`. This should be unnecessary.

4\. Some idiomatic code already prefers `var`. For example, it's commonly used for zero values:

```go
var buf bytes.Buffer
buf.WriteString("hello world!")
_ = buf.Bytes()
```

5\. As shown above, `var` allows to specify the type. Type inference is nice, but sometimes you _have_ to spell it out:

```go
num := 10

num := float64(10)

var num float64 = 10

var num = float64(10)
```

Without `:=`, you'd have less choice, which is good.

6\. `var` also allows the blank identifier:

```go
var _ = 123 // compiles
_ := 123    // doesn't compile
```

7\. `var` is also better for code highlighting. While writing a Go [syntax definition](https://github.com/mitranim/sublime-gox) for Sublime Text, I found that it's impossible to correctly scope the following:

```go
one,
    two := someExpression
```

Scoping the variable names as declarations with `:=` requires multiline lookahead or backtracking, neither of which is supported in the modern Sublime Text syntax engine.

With `var`, this can be properly scoped without multiline lookahead or backtracking:

```go
var one,
    two = someExpression
```

#### Migration

Completely embracing `var` requires an addition to the language. Various forms of `if`, `for`, `select`, and `switch` currently support `:=` but not `var`:

```go
// compiles ok
select {
    case err := <-errChan:
    case msg := <-msgChan:
}

// doesn't compile
select {
    case var err = <-errChan:
    case var msg = <-msgChan:
}
```

For Go1, adding the missing `var` support would be a safe, backwards-compatible change.

See the related [gofmt change](#gofmt-declarations).

### Remove parenthesized lists from `var`, `const`, `type`, `import` {#remove-paren-lists}

Let's start with arguments in favor of the feature.

Currently, parenthesized lists have exactly _one_ non-aesthetic reason to exist: `const (...)` enables the use of `iota`, acting as its scope.

`import` is traditionally listed, so the keyword doesn't repeat:

```go
import (
    "bytes"
    "encoding"
    "encoding/base64"
)
```

```go
import "bytes"
import "encoding"
import "encoding/base64"
```

That's a weak-ass justification for an entire language feature, made even weaker by [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) which edits your `imports` automatically.

Now, arguments against the feature.

Code should be convenient to type and edit. I think having options hinders that. Every time you write adjacent vars, some of your neurons are wasted on choosing between:

```go
var one = _
var two = _
```

and:

```go
var (
    one = _
    two = _
)
```

Worse, it occasionally leads to menial conversions between the two. That's a waste of brainpower and typing. Let's say you have a single var:

```go
const ichi = 10
```

Now you're adding another:

```go
const ichi = 10
const ni = 20
```

You might be compelled to convert to the list style:

```go
const (
    echi = 10
    ni   = 20
)
```

We've now wasted some brainpower and typing. Without lists, this would not have happened.

For consistency, the `go.mod` syntax should also remove lists.

### Maybe remove `iota` due to removing lists {#remove-iota}

`iota` requires parenthesized `const (...)` for scoping. Removing lists also leads to removing `iota`.

While I tend to avoid `iota`, I don't have a strong argument against it. If keeping `iota` in the language is important, then instead of removing lists entirely, we could just consider them non-idiomatic _unless_ `iota` is used.

### Remove `new` in favor of `&` {#remove-new}

`new` was relevant when `&` was allowed only on "storage locations" such as variables and inner fields. Now that `&` is allowed on [composite literals](https://golang.org/ref/spec#Composite_literals), `new` is close to obsolete.

`new` is limited to a zero value, while `&` allows content:

```go
client := new(http.Client)
client.Timeout = time.Minute

client = &http.Client{Timeout: time.Minute}
```

Currently, `&` doesn't work with non-composite literals:

```go
// doesn't compile
_ = &"hello world!"
```

Before `new` can be removed, `&` needs to be extended to support primitive literals. That would make it strictly more powerful than `new`.

Allowing `&` on primitives would also make it easier to print Go data structures as code. Currently, pretty-printing libraries have to resort to ugly workarounds to support those types.

Note that most code can already be converted to `&`. Code like `new(string)` or `new(int)` should be extremely rare in the wild.

For Go1, extending `&` to primitive literals would be a safe, backwards-compatible change.

### Remove dot-import: `import . "some-package"` {#remove-dot-import}

Dot-import splurges all exported definitions from another package into the current scope:

```go
import . "fmt"

func main() {
    Println("hello world!")
}
```

Having read a considerable amount of code in multiple languages with this import style, I'm convinced that it's always a bad idea. Subjectively, it makes the code harder to understand and harder to track down the definitions. Objectively, it makes the code more fragile against changes.

### Remove if-assignment and derivatives: `if _ := _ ; _ {}` {#remove-if-assignment}

Subjectively, I find this form annoying to type and annoying to read. Objectively, it's a choice, and this post is predicated on "choice is bad". This wastes everyone's brainpower; anyone reading the code has to be aware of both syntactic forms.

Instead of two options:

```go
if ok := _; ok { _ }

ok := _
if ok { _ }
```

Let's leave just _one_ option:

```go
var ok = _
if ok { _ }
```

If subscoping the variable is vital, just use a block. This also allows you to subscope more than one variable.

```go
{
    var ok = _
    if ok { _ }
}
```

### Remove short float syntax {#remove-short-float-syntax}

<span class="fg-faded">(This entry was added at 2020-06-11.)</span>

In Go, the following forms are equivalent:

```go
var _ = 0.123
var _ = .123
```

The short form works only for numbers below `0` and is not essential. The long form is essential and more general. Subjectively, I find the short form slightly harder to read; my brain starts thinking about typos and other syntactic forms involving dots. Objectively, it creates an unnecessary choice. Let's leave just one option: the "long" form.

## Tool Changes {#tools}

### Gofmt: align adjacent non-listed `var`, `const`, `type`, `import` {#gofmt-declarations}

Currently, `gofmt` aligns adjacent assignments only in parenthesized lists:

```go
const (
    ichi = 10
    ni   = 20
    san  = 30
)

const ichi = 10
const ni = 20
const san = 30
```

After [removing parenthesized lists](#remove-paren-lists), we probably want `gofmt` to align adjacent non-parenthesized assignments:

```go
const ichi = 10
const ni   = 20
const san  = 30
```

## Misc {#misc}

While writing this post, I tried to argue that complex numbers should be moved from built-ins to the standard library, but ended unconvinced.

Arguments _pro_:

* removing built-ins simplifies the language
* can implement additional math functions as methods
* can implement encoding and decoding methods for various formats

Arguments _contra_:

* breaks code
* additional functions can be provided as a package, mirroring `math`
* support for encoding and decoding can be added to the corresponding packages: `strconv`, `fmt`, `encoding/json`, `encoding/xml`, etc.

In the end, I'm not convinced that it's worthwhile.

---

Have any thoughts? Let me know!
