
In all languages, all functions have exactly one input (`$0`) and one output (`$1`). Usually, the input must be a heterogeneous tuple with built-in deconstruction syntax. Sometimes, the output must also be a heterogeneous tuple (Go, Odin).

The restriction makes no sense. It should be possible to define functions with an arbitrary input type, addressed in a canonical built-in way (say, `$0`). The essense of the input is often communicated through function name or input type. Languages shouldn't also force us to name the only argument. Programmers often fail at it anyway.
