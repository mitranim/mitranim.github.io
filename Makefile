# Usage
#
#   "make -j"   -- build or rebuild manually
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
#   Global dependencies are listed on the "deps" task.
#   Run "make deps" to install.
#
# TODO
#
#   * Minify HTML to avoid whitespace gotchas
#   * Restart Nginx on config changes
#   * Kill Nginx on SIGTERM

ABSTRACT   = .PHONY
FSWATCH    = fswatch -l 0.1 # writes absolute paths to stdout
CLEAR_TERM = printf "\x1bc\x1b[3J"
REFRESH    = curl http://localhost:52694/broadcast

$(ABSTRACT): all
all: cmd static html styles images

# Requires "-j": "make w -j"
$(ABSTRACT): w
w: all cmd-w static-w html-w styles-w images-w server notify-w make-w

cmd: cmd.go
	@go build cmd.go

$(ABSTRACT): cmd-w
cmd-w:
	@$(FSWATCH) cmd.go |     \
	while read;              \
	do                       \
		$(CLEAR_TERM);       \
		$(MAKE) cmd &&       \
		echo "[cmd] done" && \
		$(MAKE) html &&      \
		$(REFRESH);          \
	done

$(ABSTRACT): static
static: static/**/*
	@rsync -r static/ public/

$(ABSTRACT): static-w
static-w:
	@$(FSWATCH) static |        \
	while read;                 \
	do                          \
		$(CLEAR_TERM);          \
		$(MAKE) static &&       \
		echo "[static] done" && \
		$(REFRESH);             \
	done

$(ABSTRACT): html
html: public/%.html

# The "styles" dependency is for asset hashing for asset links.
public/%.html: cmd styles templates/**/*
	@./cmd

$(ABSTRACT): html-w
html-w:
	@$(FSWATCH) templates | \
	while read;             \
	do                      \
		$(CLEAR_TERM);      \
		$(MAKE) html &&     \
		$(REFRESH);         \
	done

$(ABSTRACT): styles
styles: public/styles/main.css

public/styles/main.css: styles/*.scss
	@mkdir -p public/styles
	@sassc styles/main.scss | minify --type=css > "${@}"
	@#sassc styles/main.scss > "${@}"
	@echo "[styles] wrote ${@}"

$(ABSTRACT): styles-w
styles-w:
	@$(FSWATCH) styles |  \
	while read;           \
	do                    \
		$(CLEAR_TERM);    \
		$(MAKE) styles && \
		$(REFRESH);       \
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
		$(CLEAR_TERM);                                                    \
		gm convert "$${file}" "public/images/$${file#$$(pwd)/images/}" && \
		echo "[images] wrote public/images/$${file#$$(pwd)/images/}" &&   \
		$(REFRESH);                                                       \
	done

$(ABSTRACT): server
server:
	@echo "Starting server at http://localhost:52693"
	@nginx -p . -c srv.nginx

$(ABSTRACT): notify-w
notify-w: notify
	@./notify

notify: notify.go
	@go build notify.go

# Note: we truncate `pwd` because fswatch gives us absolute paths.
$(ABSTRACT): make-w
make-w:
	@$(FSWATCH) $(MAKEFILE_LIST) |                                            \
	while read file;                                                          \
	do                                                                        \
		$(CLEAR_TERM);                                                        \
		echo "[make] $${file#$$(pwd)/} has changed, don't forget to restart"; \
	done

# Currently MacOS only. Requires Homebrew: https://brew.sh.
#   https://golang.org
#   https://github.com/sass/sassc
#   http://www.graphicsmagick.org
#   https://github.com/emcrisostomo/fswatch
#   https://github.com/tdewolff/minify/tree/master/cmd/minify
$(ABSTRACT): deps
deps:
	@brew install go sassc graphicsmagick fswatch
	@(cd && go get github.com/tdewolff/minify/cmd/minify)

# Doesn't remove binaries
$(ABSTRACT): clean
clean:
	@rm -rf public/*

$(ABSTRACT): deploy
deploy:
	@                                                                          \
	export PRODUCTION=true                   &&                                \
	$(MAKE) clean                            &&                                \
	$(MAKE) all -j                           &&                                \
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
		rm -rf .git                               &&                           \
		cd ..                                     &&                           \
		export PRODUCTION=''                      &&                           \
		$(MAKE) all -j;                                                        \
	fi
