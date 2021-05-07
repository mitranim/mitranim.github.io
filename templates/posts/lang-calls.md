Note: also see `lang-call-style.md`.

## NOTE TO SELF

Might convert this to a more general post about syntax structure, particularly about how many languages try to be like English sentences with a bit of math spliced in, with various arguments pro and contra.

Observation: language structure is all about _operators_ called with _operands_. Regardless of syntax, the structure of every program is all about this. It's extremely important to distinguish operators from operands. And the syntax must help with this. The role of keywords in most syntaxes is to act as such "operators". Keywords tend to be special-cased in parsers and syntax highlighting engines because there's no other means of distinguishing them from operands. This results from the misguided goal of making it look like English.

Once we settled on the delimiter-based calling convention like `fun(a b c)`, it was over. The use of keywords should have stopped right there. They should be simply defined as pseudo-functions or macros, see Lisp.

Special characters get a pass if they always act as operators, and never as operands. The important part is easy visual differentiation.

Uppercasing keywords in SQL is an attempt to differentiate operators from operands.

Punctuation is a consequence of sentence-oriented syntax.

## Single form for calls

Across languages, we have a variety of calling conventions:

```
prefix outside delimiters     func(arg0 arg1)
prefix inside delimiters      (func arg0 arg1)
prefix variadic               func arg0 arg1
prefix unary                  -arg
infix binary                  arg0 + arg1
postfix unary                 arg?
postfix variadic (RPN)        arg0 arg1 func
```

These days, most languages use prefix outside delimiters, unary prefix, and binary infix. Some also use unary postfix (Swift, Rust). Haskell doesn't use delimiters for arguments, and allows prefix and infix for both functions and operators. The Lisp family claims to use only prefix inside delimiters; in reality, they also use prefix unary and infix binary operators.

Technically, the delimited form is usable for everything. Nullary, unary, binary, polyadic calls:

```
+()
-(10)
-(10 20)
-(10 20 30)
.(one two three)
```

In practice, most of us won't give up `-10` and `one.two.three`, necessitating unary prefix and infix operators. Lisps have them too, as hacks:

```
-10
@ref
'form
(10 . 20)
one:two:three
one/two.three
```
