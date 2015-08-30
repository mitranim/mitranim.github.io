Lately I've been trying to figure out how to write shorter programs. Or, more
generally, how to design simple solutions.

I often hear that "less is more", that you should
[KISS](http://en.wikipedia.org/wiki/KISS_principle) and follow
[YAGNI](http://en.wikipedia.org/wiki/YAGNI), yada yada. A small program is easy
to understand and cover with tests. A simple API is pleasant to use. But that's
still abstract. What's a practical recipe for keeping things small? We might
define two attack vectors:

* Reducing the scope of the problem.
* Seeking general case solutions to special case problems.

## Scope Reduction

This approach is as simple as it gets. Saying no to a _problem_ spares you from
having to implement a solution.

Sometimes you need to draw a line and say that this feature shouldn't be in the
library, the user should write a bit of glue code instead. Or that this extra
concept is not worth the code savings it produces.

For programs with one well-defined function, this is known as the Unix
philosophy and is straightforward to follow. But it's also useful for programs
with a potentially unbounded scope, like a data modeling library or a [language
compiler](http://golang.org). A surprising number of ideas turns out to be dead
weight after a while.

Curiously, this takes willpower, or _restraint_, which seems to be an unpopular
feature with developers. Adding moving parts is interesting. Being lazy is not
enough; you have to apply mental _effort_ to refuse additions and keep things
simple.

## General Case Solutions

Programs with an unbounded scope accumulate complexity as a result of tackling
new problems, usually in response to feedback. Feedback tends to focus on
specific use cases. Addressing them individually leads to accumulating special
case solutions, even for problems that could be addressed with a general case
feature, if this class of problems could be foreseen in advance.

Feature feedback also indicates that the application scope _perceived_ by users
outranges its design scope. Including a new feature or addressing a new use case
would expand its implementation scope, which should be defined by the design
scope, not the other way around. Which means agreeing to expand a program should
begin by exploring and expanding its design scope, as if the system was being
designed anew.

Therefore the default reaction to a feature request should be figuring out what
class of problems it represents, and either refusing it entirely, or addressing
the entire _class_ instead.

## Conclusion

Every person is different, but for me, both things boil down to restraint. It's
tempting to add new moving parts. It's tempting to address a special case
instead of figuring out a wider class of problems and a solution that covers
them all. You need to stop yourself, take a step back, and remember that taking
the time to find the _right_ problem to solve will spare you from
[throwing](https://github.com/Mitranim/datacore/commit/2ce33186c0a45024c632ea8f5a113e6780cfb398)
solutions away.
