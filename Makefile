# Usage
#
#   "make"      -- build or rebuild manually
#   "make w -j" -- build, start server, watch and rebuild
#
# Notes
#
#   "$$" escapes shell variables.
#   "$${varname#prefix}" interpolates the variable without the prefix.
#   "fswatch -l N" means latency = N (default 1 second, too slow)
#   "fswatch -o" means empty output
#   "xargs -n1 -I{}" means invoke without arguments
#   "make w" requires concurrency and MUST be run with "-j"
#
# Dependencies
#
#   https://golang.org (then "go get -d" in this directory)
#   https://github.com/sass/sassc
#   https://github.com/tdewolff/minify/cmd/minify
#   http://www.graphicsmagick.org
#   https://github.com/emcrisostomo/fswatch
#
# Optional dependencies
#
#   https://github.com/Mitranim/gorun
#
# TODO
#
#   Report when rebuilding in watch mode; old errors may confuse
#   Minify HTML to avoid whitespace gotchas?

# This writes absolute paths, space-separated, to stdout
FSWATCH_LINE = fswatch -l 0.1 -0
FSWATCH_MUTE = fswatch -l 0.1 -o
INVOKE = xargs -n1 -I{}
ABSTRACT = .PHONY

$(ABSTRACT): all
all: cmd static html styles images

# Requires "make w -j"
$(ABSTRACT): w
w: all cmd-w static-w html-w styles-w images-w server make-w

cmd: cmd.go
	@go build cmd.go

$(ABSTRACT): cmd-w
cmd-w:
	@$(FSWATCH_MUTE) cmd.go | $(INVOKE) $(MAKE) cmd html

$(ABSTRACT): static
static: static/**/*
	@rsync -r static/ public/

$(ABSTRACT): static-w
static-w:
	@$(FSWATCH_MUTE) static | $(INVOKE) $(MAKE) static

$(ABSTRACT): html
html: public/%.html

# The styles are for asset hashing
public/%.html: cmd styles templates/**/*
	@./cmd

$(ABSTRACT): html-w
html-w:
	@$(FSWATCH_MUTE) templates | $(INVOKE) $(MAKE) html

$(ABSTRACT): styles
styles: public/styles/main.css

public/styles/main.css: styles/*.scss
	@mkdir -p public/styles
	@sassc styles/main.scss | minify --type=css > $@
	@echo "[styles] Wrote $@"

$(ABSTRACT): styles-w
styles-w:
	@$(FSWATCH_MUTE) styles | $(INVOKE) $(MAKE) styles

$(ABSTRACT): images
images: images/*
	@mkdir -p public/images
	@# Create a multiline batch file and pipe it to graphicsmagick.
	@(for file in $?; do\
		echo convert $$file public/images/$${file#images/};\
	done) | gm batch -

# Note: fswatch gives us absolute paths
$(ABSTRACT): images-w
images-w:
	@$(FSWATCH_LINE) images | while read -d "" file; do\
		gm convert $$file public/images/$${file#$$(pwd)/images/};\
	done

$(ABSTRACT): server
server:
	@if command -v gorun > /dev/null; then\
		gorun -w server.go;\
	else\
		go run server.go;\
	fi

# Note: fswatch gives us absolute paths
$(ABSTRACT): make-w
make-w:
	@$(FSWATCH_LINE) $(MAKEFILE_LIST) | while read -d "" file; do\
		echo \"$${file#$$(pwd)/}\" has changed. Don\'t forget to restart.;\
	done

$(ABSTRACT): clean
clean:
	@rm -f cmd
	@rm -rf public/*

$(ABSTRACT): deploy
deploy: clean all
	@(\
		cd public &&\
		rm -rf .git &&\
		git init &&\
		git remote add origin https://github.com/Mitranim/mitranim.github.io.git &&\
		git add -A . &&\
		git commit -a -m gh &&\
		git push -f origin master\
	)
