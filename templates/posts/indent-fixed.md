Follow-up to the post [Always Spaces, Never Tabs](/posts/spaces-tabs) on indentation.

Indentation is not always fixed. Sometimes people use variable indentation:

```scm
(define (some-func a b c d)
        (+ (- a b)
           (- c d)))
```

```sql
select *
from table0
     cross join table1
where table0.one = 10 and
      table1.two = 20
```

My wild guess is this style was useful before decent code editors. It might be easier to consistently follow if your editor doesn't have support for fixed-size indentation. Modern editors preserve the current indentation level on Enter, have hotkeys for indenting and unindenting by a fixed amount, and automatically increase or decrease indentation based on language context, such as between `{}`. They usually don't have built-in support for the style above.

The main _objective_ problems with this style are:

* Makes it much harder to edit and restructure code. Causes you to reindent more code when editing unrelated lines, compared to the fixed style.
* Makes it harder to tell the nesting level by indentation.

## Tooling

Lisp programmers will point at editor plugins, such as Lisp indentation in Emacs, or more recent algorithms such as Parinfer. (Sidenote: Parinfer handles more than indentation, but _all_ of the problems it handles are self-inflicted just like indentation.)

Innocent bystanders will rightfully raise eyebrows {{Emoji "🤨" ""}}. This indicates bad design and bad conventions.

## Fixed

I recommend fixed-size indentation. The Lisp `let` is particularly bad for indentation, so let's avoid it.

```scm
(define (some-func . args)
  (define one 10)
  (define two 20)
  (three one two)
)
```

```sql
select
  *
from
  table0
  cross join table1
where
  table0.one = 10 and
  table1.two = 20
```
