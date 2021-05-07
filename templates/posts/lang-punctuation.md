**TLDR**: Programming languages should be designed to minimize punctuation such as `, ;`.

## Punctuation

Unnecessary punctuation costs a lot of typing and rerunning a program after fixing missing commas.

A well-designed syntax doesn't need commas for separating expressions, or semicolons after separating statements. Many modern languages don't use semicolons (Go, Swift). In most modern languages, comma is an atavism: the language already knows where expressions start and end, and omitting a comma makes it complain about the missing `,`, rather than fusing expressions together.

SQL really needs those commas and semicolons because of its sentence-oriented design.

Part of the reason is prefix keywords. Commas create "dead space" in the syntax, reserving space for future extensions like new keywords.

```
func str_append(inout a String, b String) { a += b }
```

Another part of the reason is infix operators. Compare:

```
some_func(10 + 20, 30)
some_func(10 + 20 30)
```

The latter can be unambiguous for a compiler, but visually ambiguous for humans.

By preferring the "call" form we can also avoid punctuation, without any visual ambiguity:

```
func(str_append(a Ref(String) b String) +=(a b))

some_func(+(10 20) 30)
```

**SQL example**. SQL has this:

```
select
  10,
  20 two,
  30 as three
```

A trailing comma at the end would be a syntax error. This gets in the way of development. Sometimes I test a small piece of code by frequently rerunning it. This can involve commenting and uncommenting lines. Because of commas, I can't just comment/uncomment a line or two, I have to edit another adjacent line.

To start dropping commas, you'd have to redesign the entire SQL syntax. Consider this possibility:

```
select(
  row(
    one   10
    two   20
    three 30
  )
)
```

## Observations

Preference towards or against punctuation is influenced by how you edit code. If you tend to move and delete by "word", punctuation tends to get in the way. Vim users may move and delete from/until specific characters, and punctuation can be useful there, while lack of punctuation can be a bother.
