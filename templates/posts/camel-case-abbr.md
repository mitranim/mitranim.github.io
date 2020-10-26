**Edit 2020-10-21**: see the newer post [Language Design: Case Conventions](/posts/lang-case-conventions).

Programming has the concept of an "identifier". Identifiers are used for keywords, variable names, etc. Most languages restrict identifiers to Latin letters, digits, and an underscore.

An identifier may consist of several words without spaces. The commonly used [case styles](https://en.wikipedia.org/wiki/Letter_case) can distinguish individual words:

```
oneTwo       -- lower camel case
OneTwo       -- title camel case
one_two      -- lower snake case
ONE_TWO      -- upper snake case
one-two      -- lower kebab case
ONE-TWO      -- upper kebab case
```

All have at least two desirable properties: "separability" and "consistency". Words must be separable; consistency ensures this rule is followed without exceptions.

A problem peculiar to `TitleCamelCase` is how to treat abbreviations. Behold this monstrosity from JavaScript's DOM API:

```
XMLHttpRequest
```

It's inconsistent: `XML` is spelled in capitals, while `HTTP` is spelled in title case, like a word. What gives? There were three ways to spell it out:

1. Ignore abbreviations: `XmlHttpRequest`
2. Let them combine into one: `XMLHTTPRequest`
3. Use inconsistent casing: `XMLHttpRequest`

We see that (2) breaks separability while (3) breaks consistency. The general conclusion is that insisting on abbreviations leads to weird names, and is not compatible with the desirable properties of case styles.

The only generally consistent approach is to ignore abbreviations, i.e. treat them as words:

```
XmlHttpRequest
```

As a bonus, non-abbreviated `TitleCamelCase` is easier to automatically parse, convert to other cases, and reverse. Example:

```
XmlHttpRequest -> xml_http_request -> XmlHttpRequest
```

For automatic tools, parsing inconsistent abbreviation is not impossible; for example, my Sublime Text [plugin](https://github.com/mitranim/sublime-caser) for converting between cases can handle this. But there's just no good reason to make it harder.

Finally, having just _one_ choice means less thinking, which is good.

That's all.
