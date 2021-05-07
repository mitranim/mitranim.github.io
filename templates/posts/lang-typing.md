For brevity, this document uses "typed" to mean "statically typed" and "untyped" to mean "dynamically typed" (rather than "everything is a void pointer").

## Always Typed

After writing various complicated programs in various languages, I never want to use an untyped language again.

..........

Years of experience, of writing complicated programs in various languages, have convinced me that dynamic typing is overused. I never want to use an untyped language again.

This hasn't always been true. In the past decades, typed languages were unwieldy. Their type systems sucked, If you wanted things done,

## Cast vs Conversion

"Cast" means reinterpreting existing memory as another type, for free. It should only be allowed between types sharing the same memory layout. We can cast concrete types or pointers, like in Go. Casting has a syntax, for example `Type(val)` or `cast(val Type)`.

"Conversion" means doing CPU work to create the new value. Conversion doesn't have a syntax, and involves calling a function. Integer-to-float and backwards might be a conversion, not a cast. Other examples include integer-to-string, string-to-bool, and so on.

If type constraints exist (see below), casts into constrained types are not free; they include a constraint check.

## Type Constraints

Suppose we could define type constraints, like in SQL. These would be checked whenever a value of that type is constructed or altered. "Construction" also includes casts.

Where possible, constraints should be evaluated at compile time.

This feature requires ergonomic errors. When every operation involving a given type can produce an error, it might not be practical to handle each of those errors manually. Unchecked exceptions would solve this, but we want checked errors. We could consider Swift-style or Rust-style syntax for optional and result types. Alternatively, we could consider a macro that implicitly does that for the optional/result expressions it contains.

This could be generalized into more than type constraints. You could define these compile-time "contracts" and attach them to functions or other arbitrary places in the code. This may not require any special syntax. We might be able to teach the compiler to evaluate as much of the program as possible at compile time, finding anything that's guaranteed to fail. This could slow down compilation, and made optional.
