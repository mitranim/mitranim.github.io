**TLDR**: Identifiers in programming languages should use only `snake_case`, `Title_snake_case`, `UPPER_SNAKE_CASE`, ignore abbreviations, and be limited to ASCII alphanumerics with `_`.

This post will also touch on the structure of identifiers.

There was an earlier, more specialized post: [Don't Abbreviate In Camel-Case](/posts/camel-case-abbr). This one is more general.

## Lower case

_Objective_ arguments in favor of `snake_case` over `camelCase`:

* Works in case-insensitive systems such as SQL.
* Can remap `_` to type without Shift. (I did.)
* Unambiguously convertible when numbers are involved.
* Refactoring requires only 1 renaming instead of 2.

Conversion:

```
one_123_two <-> one 123 two

one123two   <-> one123two
one123two   <-> one123 two
one123two   <-> one 123 two

one123Two   <-> one123 two
one123Two   <-> one 123 two
```

## Title case

_Objective_ arguments in favor of `Title_snake_case` over `TitleCamelCase`:

* Can properly support abbreviations, for example `XML_HTTP_request`. No schizophrenia such as `XMLHttpRequest`.
* Can remap `_` to type without Shift. Titled identifiers require only one Shift press. (I did.)
* Consistent with `snake_case` in a language that uses it for lowercase identifiers.
* Unambiguously convertible when numbers are involved.

Conversion:

```
One_123_two <-> one 123 two

One123two   <-> one123two
One123two   <-> one123 two
One123two   <-> one 123 two

One123Two   <-> one123 two
One123Two   <-> one 123 two
```

## Abbreviations

_Objective_ arguments in favor of avoiding abbreviations, for example `Json_encoder` over `JSON_encoder`, or `JsonEncoder` over `JSONEncoder`:

* Fewer rules.
* Less thinking.
* Simpler code.
* No schizophrenia such as `XMLHttpRequest`.

Example from work.

At some point I had contact with a code base involving generating Go code from Swagger. The generator had a variety of special cases for `id`, `xml`, and some other abbreviations. A field named `xml_setting_id` would become `XMLSettingID`. However, if you used an abbreviation _unknown_ to the generator, for example XSD (XML Schema Definition), `xsd_setting_id` would become `XsdSettingID`.

The goal was noble: be consistent with the Go standard library, which stupidly uses abbreviations, for example `MarshalXML`. But unlike the standard library, you couldn't just remember "abbreviations are uppercase", your brain needed the database of the _exact_ abbreviations special-cased in that generator. So don't. Don't use abbreviations in identifiers, and don't special-case them in code generators or parsers.

## Characters

_Objective_ arguments in favor of restricting identifiers to ASCII alphanumerics with `_`:

* Interoperable between all languages.
* Works in all encodings.
* Works in all Latin keyboard layouts.

Example from work.

At some point, we at Purelab were using Clojure and Datomic to build apps. Clojure symbols (Lisp equivalent of identifiers) use `kebab-case` and may contain operator characters such as `-?`. Booleans are expected to end with a question: `hidden?` instead of `is_hidden`.

Datomic has its own idiosyncrasy: column names are global and include the entity type. So, instead of this:

```sql
create table persons (is_email_verified bool not null default false);
```

...you use this:

```clj
{:db/ident     :person/email-verified?
 :db/valueType :db.type/boolean}
```

For simplicity, let's suppose we use Postgres, and have a JS client. You have to either break the SQL and JS conventions by quoting the field:

```sql
create table persons ("email-verified?" bool not null default false);
```

```js
person['email-verified?']
```

...or break the Clojure convention by using the interoperable format:

```clj
:is_email_verified
```

## Footnote on Lisp symbols

Lisps allow identifiers like `email-verified?` because they don't distinguish identifiers and operators. They just have "symbols". This has various problems.

* People define custom operators, creating inscrutable code. Popular in Haskell. What the hell is `>>=`? With `bind`, you can at least start _guessing_ the purpose, or pronounce it, or google it, what a feat!
* Leads to hacks like embedding `: / .` in symbols to implement namespacing (Common Lisp, Clojure). This requires re-parsing the symbol, something the AST should have done for you. Clojure symbols are classes with "namespace" and "name" parts, indicating that they were combined prematurely. The AST should split alphanumerics and operators from the start.

## Conclusion

When making a language, follow the conventions listed at the top. Let's solve this forever and move on.
