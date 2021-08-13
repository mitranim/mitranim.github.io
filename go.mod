module main

go 1.16

require (
	github.com/alecthomas/chroma v0.9.3
	github.com/gotidy/ptr v1.3.0
	github.com/mitranim/gax v0.2.0
	github.com/mitranim/srv v0.0.0-20210207104346-0df64d1a7dff
	github.com/mitranim/try v0.1.2
	github.com/pkg/errors v0.9.1
	github.com/rjeczalik/notify v0.9.2
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0
)

// Remove when upstream releases v0.9.3
replace github.com/alecthomas/chroma => github.com/mitranim/chroma v0.9.3-0.20210731070231-be6172b784ed
