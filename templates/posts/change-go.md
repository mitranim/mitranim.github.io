Changes that I would like to see in Go.

These changes seem unrealistic, but a person can dream, right?

Note that I would _not_ accept any of these changes as a tradeoff against compilation speed. To me, fast compilation is a critically important language feature. Let's not trade it off.

# Table of Contents

* Generics
    * Revise the standard library
        * Replace reflection with generics where possible
        * Utility functions, like the "math" package, should become generic over any type convertible to their current input
* Statically dispatched interfaces (e.g. like Rust's generic traits)
    * Revise the standard library
        * Replace virtual interfaces with static where possible
* Unify strings and bytes into one type
    * The unified type is mutable
    * The unified type is equatable by `==`; other slice types are not
    * Since this type is mutable, we can no longer define string constants
    * Revise the standard library
        * Unify the `strings` and `bytes` packages, remove redundancies
        * Unify the regexp API
        * Etc.
    * All string-using code must revise assumptions about mutability and copying
* General approach to nullability that doesn't involve pointers and plays well with JSON, SQL, etc. Ideally, it shouldn't require data definitions to "know" which fields are nullable?
* Unify competing APIs for efficiently appending bytes; namely, `io.Writer` and other writer interfaces must become as efficient as appending to a byte slice

# Convention Changes

## "NewX" → "MakeX"

Bonus idea. This is a not a language change, and it seems unrealistic due to the sheer amount of code written under the "New" convention.

It would be better if we called constructors "MakeSomething" rather than "NewSomething", for a practical reason: "make" is a verb, which allows modifications like "MustMakeSomething". Trying to include "must" into a "New"-style constructor is painful.
