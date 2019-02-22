# Usage
#
#   "make -j"       -- build or rebuild manually
#   "make w -j"     -- build, start server, watch and rebuild
#   "make deploy"   -- deploy while "w" is not running
#   "make w-deploy" -- deploy while "w" is running
#
# Dependencies
#
#   Global dependencies are listed on the "deps" task.
#   Run "make deps" to install.
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
# Env Vars
#
#   Env vars in Make have a major gotcha. Make uses file timestamps to decide
#   when to rebuild a target or skip the rebuild. Our Makefile is explicitly
#   written to allow this behavior wherever possible. Merely changing an env
#   var and rerunning a rule doesn't guarantee a rebuild. This means that any
#   env var change must be accompanied by a "make clean".
#
#   Using "<RULE>: export <VAR> = <VAL>", a rule can set an environment var for
#   itself and all its dependencies. The "=" can be replaced with various forms
#   such as ":=" or "?=". The latter is a "default" that runs only if the var is
#   not already defined. This allows an external override.
#
#   See docs on variable assignment:
#
#     https://www.gnu.org/software/make/manual/make.html#Reading-Makefiles
#     https://www.gnu.org/software/make/manual/make.html#Variables_002fRecursion
#
# Change Detection
#
#   Make avoids rebuilding targets that are newer than thair sources. This
#   feature is biased towards single-file targets, and doesn't work as expected
#   with multi-file targets. For example, if one of the previously created
#   output files is no longer generated or updated, Make will always rerun the
#   rule because that file will always be older than some of the sources.
#
#   To work properly, multi-file outputs need either a phony target or an empty
#   target. A phony target simply causes the rule to always run. An empty
#   target is an empty file that marks the timestamp of the last rebuild,
#   allowing Make to compare timestamps to decide when to rebuild, just like
#   with single-file outputs. That's what we use.
#
# TODO
#
#   * Minify HTML to avoid whitespace gotchas
#   * Restart Nginx on config changes

ABSTRACT   = .PHONY
FSWATCH    = fswatch -l 0.1 # writes absolute paths to stdout
CLEAR_TERM = printf "\x1bc\x1b[3J"
REFRESH    = curl http://localhost:52694/broadcast

$(ABSTRACT): all
all: cmd static html styles img

# Requires "-j": "make w -j"
$(ABSTRACT): w
w: export DEV ?= true
w: all cmd-w static-w html-w styles-w img-w server notify-w make-w

cmd: cmd.go
	@go build cmd.go

$(ABSTRACT): cmd-w
cmd-w:
	@$(FSWATCH) cmd.go |     \
	while read;              \
	do                       \
		$(CLEAR_TERM) &&     \
		$(MAKE) cmd &&       \
		echo "[cmd] done" && \
		$(MAKE) html &&      \
		$(REFRESH);          \
	done

$(ABSTRACT): static
static: public/timestamps/static

public/timestamps/static: static/* static/*/*
	@rsync -r static/ public/
	@mkdir -p public/timestamps && touch public/timestamps/static

$(ABSTRACT): static-w
static-w:
	@$(FSWATCH) static |        \
	while read;                 \
	do                          \
		$(CLEAR_TERM) &&        \
		$(MAKE) static &&       \
		echo "[static] done" && \
		$(REFRESH);             \
	done

$(ABSTRACT): html
html: public/timestamps/html

# Note: asset dependencies are used for link hashing.
public/timestamps/html: cmd public/styles/main.css templates/* templates/*/*
	@./cmd
	@mkdir -p public/timestamps && touch public/timestamps/html

$(ABSTRACT): html-w
html-w:
	@$(FSWATCH) templates | \
	while read;             \
	do                      \
		$(CLEAR_TERM) &&    \
		$(MAKE) html &&     \
		$(REFRESH);         \
	done

$(ABSTRACT): styles
styles: public/styles/main.css

public/styles/main.css: styles/*.scss
	@mkdir -p public/styles
	@#sassc styles/main.scss | minify --type=css > "${@}"
	@sassc styles/main.scss > "${@}"
	@echo "[styles] wrote ${@}"

$(ABSTRACT): styles-w
styles-w:
	@$(FSWATCH) styles |  \
	while read;           \
	do                    \
		$(CLEAR_TERM) &&  \
		$(MAKE) styles && \
		$(REFRESH);       \
	done

$(ABSTRACT): img
img: public/timestamps/img

# Optimizes raster images using graphicsmagick in batch mode.
#
# Notes on the "images" dependency. The directory's timestamp changes when
# adding or deleting a child file or directory. Image files tend to be copied
# around without changing their timestamps, so this is the only way to detect
# such changes. This is also why the task is called "img" rather than "images".
#
# Note: this breaks on filenames with spaces. This happens for so many reasons
# that the best workaround is to avoid spaces in names.
public/timestamps/img: images images/*
	@(\
		find images -type f -name '*.jpg' -o -name '*.png' |  \
		while read path;                                      \
		do                                                    \
			if ! mkdir -p "public/$$(dirname $${path})";      \
			then exit 1;                                      \
			fi;                                               \
		done                                                  \
	)
	@(                                                        \
		find images -type f -name '*.jpg' -o -name '*.png' |  \
		while read path;                                      \
		do                                                    \
			echo "convert" "$${path}" "public/$${path}";      \
		done                                                  \
	) | gm batch -
	@mkdir -p public/timestamps && touch public/timestamps/img

# Note: we truncate `pwd` because fswatch gives us absolute paths.
$(ABSTRACT): img-w
img-w:
	@$(FSWATCH) images |                            \
	while read file;                                \
	do                                              \
		$(CLEAR_TERM)                            && \
		path=$${file#$$(pwd)/}                   && \
		mkdir -p "public/$$(dirname "$${path}")" && \
		gm convert "$${file}" "public/$${path}"  && \
		echo "[img] wrote public/$${path}"       && \
		$(REFRESH);                                 \
	done

# Unlike most processes, Nginx doesn't terminate on SIGHUP and continues
# running in the background when the terminal tab is closed. To terminate it
# along with other processes, we have to trap SIGHUP and kill nginx manually.
# Trapping the signal requires us to run nginx in the background; otherwise the
# trap handler wouldn't run until nginx terminates, defeating the purpose.
# Also, at least in Bash 3.2, the trap handler must be registered AFTER
# starting the background process.
$(ABSTRACT): server
server:
	@echo "Starting server at http://localhost:52693"
	@nginx -p . -c srv.nginx & trap 'jobs -p | xargs kill' INT HUP && wait

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
		$(CLEAR_TERM) &&                                                      \
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

# Note: "make clean" and "make all" ensure that the targets are rebuilt in
# "production mode" rather than "development mode".
$(ABSTRACT): deploy
deploy:
	@                                                                          \
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
		rm -rf .git;                                                           \
	fi

# When "make w -j" is running, this should be used instead of "make deploy".
$(ABSTRACT): w-deploy
w-deploy:
	@$(MAKE) deploy
	@$(MAKE) clean
	@DEV=true $(MAKE) all -j
