**TLDR**: variadic `-`, as seen in Lisps, has gotchas; it may be allowed syntactically, but not as a variadic function.

`-` tends to be overloaded with two different operations: negation and subtraction. Negation is always unary. Subtraction can be variadic. Unary subtraction is an identity function that returns the first argument unchanged _without negating it_.

```
ƒ negate(a)         = 0 - a

ƒ subtract(a)       = a
ƒ subtract(a b)     = a - b
ƒ subtract(a b c)   = (a - b) - c
ƒ subtract(a b c d) = ((a - b) - c) - d
```

In math and many programming languages, there's no ambiguity because `-` is either unary prefix (negation) or binary infix (subtraction):

```
-A       |    Negation.
B - C    |    Subtraction.
```

But in Lisps, `-` is always prefix, always variadic, and when called with a single argument, it always negates it.

The following examples use [Racket](https://racket-lang.org). Let's dynamically pass N arguments to `-`:

```scm
#lang racket/base

(define (subtract . args) (apply - args))

(println (subtract 11 33 55))
(println (subtract 11 33))
(println (subtract 11))
```

```scm
-77
-22
-11 ; Performed negation, not subtraction!
```

The last call performed _negation_ on its only argument.

Correct variadic subtraction:

```scm
#lang racket/base

(define (flip fun) (lambda (a b) (fun b a)))
(define (foldl1 fun seq) (foldl fun (car seq) (cdr seq)))
(define (subtract . args) (foldl1 (flip -) args))

(println (subtract 11 33 55))
(println (subtract 11 33))
(println (subtract 11))
```

```
-77
-22
11
```

Now, `11` was correctly returned as-is.

Worth comparing to Haskell, which also generalizes operators into functions, but handles `-` differently. In Haskell, the function `-` is always binary subtraction:

```hs
main = do
  print (foldl1 (-) [11, 33, 55])
  print (foldl1 (-) [11, 33])
  print (foldl1 (-) [11])
```

```
-77
-22
11
```

Haskell doesn't allow to overload functions on parameter count. You can't define `-` as both unary and binary. So they special-cased unary `-` in the _syntax_, converting it to `negate`:

```hs
main = do
  print (-11)
  print (negate 11)
```

```hs
-11
-11
```

Lisp and Haskell create this problem for themselves by treating `-` as a function while overloading it with _two_ different functions. Most languages don't have this problem because they don't have `-` as a function. Languages with operator overloading tend to differentiate between negation and subtraction. For example, Rust has `ops::Neg` and `ops::Sub`. Literal `-` is converted into calls to one of those. When passing it to a higher-order function, you either pass `ops::Neg::neg`, or `ops::Sub::sub`, avoiding the problem completely.

