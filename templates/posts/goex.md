{{mdToToc .MdTpl}}

This post is informed by years of Go, months of Go with exceptions, and years of other languages.

## Myths to Debunk

> "Go doesn't have exceptions".

Go has panics. Panics are exceptions.

> "Errors-as-values is simpler than exceptions".

Decent argument that doesn't apply to Go. Go already has both.

> "All errors are in function signatures".

The stdlib has many documented panics. New releases frequently add more. Panics are not in function signatures.

> "Panics are reserved for unrecoverable errors".

Untrue in Go. Panics are recoverable and actionable. For example, HTTP servers respond with 500 and error details instead of crashing.

> "Explicit errors lead to more reliable code."

Decent argument that doesn't apply to Go. Go has panics. Reliable code _must_ handle panics in addition to error values. Code that assumes "no panics" will have leaks, data corruption, and other unexpected states.

> "Panics are expensive".

Actually cheap enough.

## Observations

* "Just panics" is objectively simpler than "error values and panics". 1 is objectively less than 2.
* "Just panics" is more reliable than "error values and panics". You only need to handle 1, not 2.
* Requires some un-doctrination, after years of trying to believe in error values.
* Performance is nearly the same.
* Avoids mishandling of `err` variables.
* Exceptions and stacktraces are orthogonal. You want both.

Combination of `defer` `panic` `recover` allows terse and flexible exception handling.

Brevity:

```golang
import "github.com/mitranim/try"

func outer() {
  defer try.Detail(`failed to do X`)
  someFunc()
  anotherFunc()
  moreFunc()
}
```

Same without panics:

```golang
import "github.com/mitranim/try"

func outer() (err error) {
  defer try.WithMessage(&err, `failed to do X`)

  err = someFunc()
  if err != nil {
    return
  }

  err = anotherFunc()
  if err != nil {
    return
  }

  err = moreFunc()
  if err != nil {
    return
  }

  return
}
```

## Performance

In modern Go (1.17 and higher), there is barely any difference. Defer/panic/recover is usable even in CPU-heavy hotspot code.

## Stacktraces

Stacktraces are essential to debugging, with or without exceptions.

* Exceptions and stacktraces are orthogonal.
* Exceptions don't require stacktraces.
* You _always_ want stacktraces.
  * Many languages elide them for performance, but you still want them.
* Lack of stacktraces causes developers to _manually emulate stacktraces_.

Some real Go code, written by experienced developers, has errors annotated with function names, like this:

```golang
func someFunc() error {
  err := anotherFunc()
  if err != nil {
    return fmt.Errorf(`someFunc: %w`, err)
  }

  err = moreFunc()
  if err != nil {
    return fmt.Errorf(`someFunc: %w`, err)
  }

  return nil
}
```

You can simplify this with `defer`:

```golang
import "github.com/mitranim/try"

func someFunc() (err error) {
  defer try.WithMessage(&err, `someFunc`)

  err = anotherFunc()
  if err != nil {
    return err
  }

  err = moreFunc()
  if err != nil {
    return err
  }

  return nil
}

func anotherFunc() (err error) {
  defer try.WithMessage(&err, `anotherFunc`)
  return someErroringOperation()
}

func moreFunc() (err error) {
  defer try.WithMessage(&err, `moreFunc`)
  return anotherErroringOperation()
}
```

🔔 Alarm bells should be ringing in your head. This emulates a stacktrace, doing manually what other languages have automated decades ago.

So stop doing that. Automate your stacktraces, and shorten your code:

```golang
import "github.com/mitranim/try"

func someFunc() {
  defer try.Detail(`failed to do X`)
  anotherFunc()
  moreFunc()
}

func anotherFunc() {
  try.To(someErroringOperation())
}

func moreFunc() {
  try.To(anothrErroringOperation())
}
```
