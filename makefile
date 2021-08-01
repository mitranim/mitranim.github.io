MAKEFLAGS := --silent
PAR := $(MAKE) -j 128
TAR := public
STATIC := static
SASS := sass --no-source-map -I . styles/main.scss:$(TAR)/styles/main.css

ifeq ($(PROD), true)
	SASS_STYLE := compressed
else
	SASS_STYLE := expanded
endif

ifeq ($(OS), Windows_NT)
	RM = if exist "$(1)" rmdir /s /q "$(1)"
else
	RM = rm -rf "$(1)"
endif

ifeq ($(OS), Windows_NT)
	MKDIR = if not exist "$(1)" mkdir "$(1)"
else
	MKDIR = mkdir -p "$(1)"
endif

ifeq ($(OS), Windows_NT)
	CP = copy "$(1)"\* "$(2)" >nul
else
	CP = cp "$(1)"/* "$(2)"
endif

.PHONY: watch
watch: clean
	$(PAR) styles-w afr-w cmd-w srv pages-w images-w static-w

.PHONY: build
build: clean
	$(PAR) styles pages images static

.PHONY: styles-w
styles-w:
	$(SASS) --style=$(SASS_STYLE) --watch

.PHONY: styles
styles:
	$(SASS) --style=$(SASS_STYLE)

.PHONY: afr-w
afr-w:
	deno run -A --unstable --no-check https://deno.land/x/afr@0.5.1/afr.ts --port 52692 --verbose true

# May compile twice on startup, should probably fix.
.PHONY: cmd-w
cmd-w:
	watchexec -r -c -d=0 -e=go,mod -n -- $(MAKE) cmd

cmd: *.go go.mod
	go build -o cmd

.PHONY: srv
srv: cmd
	./cmd srv

.PHONY: pages-w
pages-w: cmd
	watchexec -r -d=0 --no-ignore -w=cmd -w=templates -n -- ./cmd pages

.PHONY: pages
pages: cmd
	./cmd pages

.PHONY: images-w
images-w:
	watchexec -r -d=0 -w=images -n -- ./cmd images

.PHONY: images
images: cmd
	./cmd images

.PHONY: static-w
static-w:
	watchexec -r -d=0 -w=static -n -- $(MAKE) static

.PHONY: static
static:
	$(call MKDIR,"$(TAR)")
	$(call CP,"$(STATIC)","$(TAR)")

.PHONY: deploy
deploy: export PROD=true
deploy: cmd build
	./cmd deploy

.PHONY: clean
clean:
	$(call RM,"$(TAR)")

.PHONY: deps
deps:
ifeq ($(OS), Windows_NT)
	scoop install sass go deno watchexec
else
	brew install -q sass/sass/sass go deno watchexec
endif
