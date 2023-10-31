This post is informed by many years of Go, and months of Go with exceptions. **I am well aware** of many arguments for error values. Some of them are addressed below.

Reddit discussion: https://www.reddit.com/r/golang/comments/r2h31i/shorten_your_go_code_by_using_exceptions/

**Update 2023-10-23.** The original version of this post referred to https://github.com/mitranim/try. The updated post refers to https://github.com/mitranim/gg, which subsumes the previous library and offers more features.

{{mdToToc .MdTpl}}

## Myths to debunk

> "Go doesn't have exceptions".

Go has panics, which are exceptions.

> "Errors-as-values is simpler than exceptions".

Decent argument that doesn't apply to Go. Go already has both. We don't get to choose to use just one.

> "All errors are in function signatures".

The stdlib has many documented panics. New releases frequently add more. Panics are not in function signatures.

> "Panics are reserved for unrecoverable errors".

Untrue in Go. Panics are recoverable and actionable. For example, HTTP servers respond with 500 and error details instead of crashing.

> "Explicit errors lead to more reliable code."

Decent argument that doesn't apply to Go. Go has panics. Reliable code _must_ handle panics in addition to error values. Code that assumes "no panics" or "panics always crash the process" will have leaks, data corruption, and other unexpected states.

> "Panics are expensive".

Panics are cheap. Stack traces have a minor cost.

## Observations

* "Just panics" is objectively simpler than "error values and panics". 1 is objectively less than 2.
* "Just panics" is more reliable than "error values and panics". You only need to handle 1, not 2.
* Requires some un-doctrination, after years of trying to believe in error values.
* Performance is nearly the same.
* Avoids mishandling of `err` variables.
* Exceptions and stack traces are orthogonal. You want both.

Combination of `defer` `panic` `recover` allows terse and flexible exception handling.

Brevity:

```golang
import "github.com/mitranim/gg"

func outer() {
  defer gg.Detail(`failed to do X`)
  someFunc()
  anotherFunc()
  moreFunc()
}
```

Same without panics:

```golang
func outer() (err error) {
  defer ErrWrapf(&err, `failed to do X`)

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

// Suboptimal implementation, only for example purposes.
func ErrWrapf(out *error, pat string, msg ...any) {
  if out != nil && *out != nil {
    *out = fmt.Errorf(fmt.Sprintf(pat, msg...)+`: %w`, *out)
  }
}
```

## Performance

In modern Go (1.17 and higher), there is barely any difference. Defer/panic/recover is usable even in CPU-heavy hotspot code.

Generating stack traces has a far larger cost. The examples in this post use `github.com/mitranim/gg` which automatically adds stack traces. If you're using stack traces with error values, that cost is already dominant, compared to the cost of defer/panic/recover.

## Stack traces

Stack traces are essential to debugging, with or without exceptions.

* Exceptions and stack traces are orthogonal.
* Exceptions don't require stack traces.
* You _always_ want stack traces for debugging.
  * Many languages elide them for performance, but you still want them.
  * Don't show stack traces to your users. They should be printed only in debug logging.
* Lack of stack traces causes developers to _manually emulate stack traces_.

Some real Go code, written by experienced developers, has errors annotated with function names, like this:

```golang
func someFunc() error {
  err := anotherFunc()
  if err != nil {
    return fmt.Errorf(`someFunc: anotherFunc: %w`, err)
  }

  err = moreFunc()
  if err != nil {
    return fmt.Errorf(`someFunc: moreFunc: %w`, err)
  }

  return nil
}
```

You can simplify this with `defer`, as shown above:

```golang
func someFunc() (err error) {
  defer ErrWrapf(&err, `someFunc`)

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

func anotherFunc() (err error) {
  defer ErrWrapf(&err, `anotherFunc`)
  return someErroringOperation()
}

func moreFunc() (err error) {
  defer ErrWrapf(&err, `moreFunc`)
  return anotherErroringOperation()
}

// Suboptimal implementation, only for example purposes.
func ErrWrapf(out *error, pat string, msg ...any) {
  if out != nil && *out != nil {
    *out = fmt.Errorf(fmt.Sprintf(pat, msg...)+`: %w`, *out)
  }
}
```

🔔 Alarm bells should be ringing in your head. This emulates a stack trace, doing manually what other languages have automated decades ago.

So stop doing that. Automate your stack traces, and shorten your code:

```golang
import "github.com/mitranim/gg"

func someFunc() {
  defer gg.Detail(`failed to do X`)
  anotherFunc()
  moreFunc()
}

func anotherFunc() {
  gg.Try(someErroringOperation())
}

func moreFunc() {
  gg.Try(anothrErroringOperation())
}
```
