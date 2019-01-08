# Usage
#
#   "make"      -- build or rebuild manually
#   "make w -j" -- build, start server, watch and rebuild
#
# Notes
#
#   "$$" rather than "$" -- prevent the interpolation from happening in make,
#   delaying it until shell execution.
#
#   "${varname#prefix}" or "$${varname#prefix}" -- interpolate the variable
#   while stripping the prefix.
#
#   "make w"        -- requires concurrency and MUST be run with "-j"
#   "fswatch <dir>" -- writes absolute paths to stdout, newline-separated
#   "fswatch -l N"  -- set latency; default is 1 second, too slow
#
# Dependencies
#
#   https://golang.org (then "go get -d" in this directory)
#   https://github.com/sass/sassc
#   https://github.com/tdewolff/minify/cmd/minify
#   http://www.graphicsmagick.org
#   https://github.com/emcrisostomo/fswatch
#
# TODO
#
#   Minify HTML to avoid whitespace gotchas?
#   Restart Nginx on config changes.
#   Watch HTML from cmd.go for faster rebuilds.

ABSTRACT   = .PHONY
FSWATCH    = fswatch -l 0.1 # writes absolute paths to stdout
CLEAR_TERM = printf "\x1bc\x1b[3J"

$(ABSTRACT): all
all: cmd static html styles images

# Requires "-j": "make w -j"
$(ABSTRACT): w
w: all cmd-w static-w html-w styles-w images-w server make-w

cmd: cmd.go
	@go build cmd.go

$(ABSTRACT): cmd-w
cmd-w:
	@$(FSWATCH) cmd.go |    \
	while read;             \
	do                      \
		$(CLEAR_TERM) &&    \
		$(MAKE) cmd html && \
		echo "[cmd] done";  \
	done

$(ABSTRACT): static
static: static/**/*
	@rsync -r static/ public/

$(ABSTRACT): static-w
static-w:
	@$(FSWATCH) static |      \
	while read;               \
	do                        \
		$(CLEAR_TERM) &&      \
		$(MAKE) static &&     \
		echo "[static] done"; \
	done

$(ABSTRACT): html
html: public/%.html

# The "styles" dependency is for asset hashing, for asset links.
public/%.html: cmd styles templates/**/*
	@./cmd

$(ABSTRACT): html-w
html-w:
	@$(FSWATCH) templates | \
	while read;             \
	do                      \
		$(CLEAR_TERM) &&    \
		$(MAKE) html;       \
	done

$(ABSTRACT): styles
styles: public/styles/main.css

public/styles/main.css: styles/*.scss
	@mkdir -p public/styles
	@sassc styles/main.scss | minify --type=css > "${@}"
	@echo "[styles] wrote ${@}"

$(ABSTRACT): styles-w
styles-w:
	@$(FSWATCH) styles | \
	while read;          \
	do                   \
		$(CLEAR_TERM) && \
		$(MAKE) styles;  \
	done

$(ABSTRACT): images
images: images/*
	@mkdir -p public/images
	@# Create a multiline batch file and pipe it to graphicsmagick.
	@(\
		for file in ${?};                                               \
		do                                                              \
			echo "convert" "$${file}" "public/images/$${file#images/}"; \
		done                                                            \
	) | gm batch -

# Note: we truncate `pwd` because fswatch gives us absolute paths.
$(ABSTRACT): images-w
images-w:
	@$(FSWATCH) images |                                                  \
	while read file;                                                      \
	do                                                                    \
		$(CLEAR_TERM) &&                                                  \
		gm convert "$${file}" "public/images/$${file#$$(pwd)/images/}" && \
		echo "[images] wrote public/images/$${file#$$(pwd)/images/}";     \
	done

$(ABSTRACT): server
server:
	@echo "Starting server at http://localhost:52693"
	@nginx -p . -c srv.nginx

# Note: we truncate `pwd` because fswatch gives us absolute paths.
$(ABSTRACT): make-w
make-w:
	@$(FSWATCH) $(MAKEFILE_LIST) |                                     \
	while read file;                                                   \
	do                                                                 \
		$(CLEAR_TERM) &&                                               \
		echo "$${file#$$(pwd)/} has changed, don't forget to restart"; \
	done

$(ABSTRACT): clean
clean:
	@rm -f cmd && rm -rf public/*

$(ABSTRACT): deploy
deploy: clean all
	@                                                                          \
	url=$$(git remote get-url origin)        &&                                \
	source=$$(git symbolic-ref --short head) &&                                \
	target=master                            &&                                \
	if                                                                         \
		[ "$${source}" == "$${target}" ];                                      \
	then                                                                       \
		echo "expected source branch to be distinct from \"$${target}\"" >&2;  \
		exit 1;                                                                \
	else                                                                       \
		cd public                                 &&                           \
		rm -rf .git                               &&                           \
		git init                                  &&                           \
		git remote add origin "$${url}"           &&                           \
		git add -A .                              &&                           \
		git commit -a --allow-empty-message -m '' &&                           \
		git branch -m "$${target}"                &&                           \
		git push -f origin "$${target}"           &&                           \
		rm -rf .git;                                                           \
	fi
