{{mdToToc .MdTpl}}

## Myths

> "Go doesn't have exceptions".

Go has panics. Panics are exceptions.

> "Errors as values are simpler than exceptions".

Decent argument that doesn't apply to Go. Go already has both.

> "All errors are in function signatures".

The stdlib has many documented panics. New releases frequently add more. Panics are not in function signatures.

> "Panics are reserved for unrecoverable errors".

Untrue in Go. Panics are recoverable and actionable. For example, HTTP servers respond with 500 and error details instead of crashing.

> "Explicit errors lead to more reliable code."

Decent argument that doesn't apply to Go. Go has panics. Reliable code _must_ handle panics. Code that assumes "no panics" will have leaks, data corruption, and other unexpected states.

> "Panics are expensive".

Cheap enough for most code. Just minimize them in bottlenecks.

## Observations

* "Just panics" is simpler than "error values and panics". 1 < 2.
* "Just panics" is more reliable than "error values and panics".
* Requires some un-doctrination.
* Performance is fine.
* Avoids mishandling of `err` variables.
* Exceptions and stacktraces are orthogonal. You want both.

Combination of `defer` `panic` `recover` allows terse and flexible exception handling.

Brevity:

```golang
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

```golang
func outer() {
  defer try.Detail(`failed to do X`)
  someFunc()
  anotherFunc()
  moreFunc()
}
```

## Performance

Exceptions are slightly more expensive than checks and returns. In CPU-heavy hotspot code, the overhead can be noticeable. In IO-heavy control code, the overhead is often not measurable, and has the biggest benefit, since most IO operations are failable.

TLDR: OK for most code, avoid in bottlenecks.

## Compatibility

Libraries should use error values rather than panics.

## Stacktraces

Stacktraces must be mentioned because they're essential.

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
  return something()
}

func moreFunc() (err error) {
  defer try.WithMessage(&err, `moreFunc`)
  return something()
}
```

🔔 Alarm bells should be ringing in your head. This emulates a stacktrace, doing manually what other languages have automated decades ago.
