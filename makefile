MAKEFLAGS := --silent
PAR := $(MAKE) -j 128
TAR := public
STATIC := static
BIN := ./bin
CMD := $(BIN)/cmd
SASS := sass --no-source-map -I . styles/main.scss:$(TAR)/styles/main.css
WATCH := watchexec -r -p -c -d=0 -n
WATCH_CMD := $(WATCH) --no-ignore -w=$(CMD)

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
	CP = cp -r "$(1)"/* "$(2)"
endif

.PHONY: watch
watch: clean
	$(PAR) styles-w cmd-w afr srv pages-w images-w static-w

.PHONY: build
build: clean
	$(PAR) styles pages images static

.PHONY: styles-w
styles-w:
	$(SASS) --style=$(SASS_STYLE) --watch

.PHONY: styles
styles:
	$(SASS) --style=$(SASS_STYLE)

.PHONY: afr
afr:
	deno run -A --unstable --no-check https://deno.land/x/afr@0.5.1/afr.ts --port 52692
# 	afr -v -p 52692

.PHONY: cmd-w
cmd-w: $(CMD)
	$(WATCH) -e=go,mod -- $(MAKE) $(CMD)

$(CMD): *.go go.mod
	go build -o $(CMD)

.PHONY: srv
srv: $(CMD)
	$(CMD) srv

.PHONY: pages-w
pages-w: pages
	$(WATCH_CMD) -w=templates -- $(CMD) pages

.PHONY: pages
pages: $(CMD)
	$(CMD) pages

.PHONY: images-w
images-w: images
	$(WATCH) -w=images -- $(CMD) images

.PHONY: images
images: $(CMD)
	$(CMD) images

.PHONY: static-w
static-w: static
	$(WATCH) -w=static -- $(MAKE) static

.PHONY: static
static:
	$(call MKDIR,$(TAR))
	$(call CP,$(STATIC),$(TAR))

.PHONY: lint
lint:
	golangci-lint run

.PHONY: deploy
deploy: export PROD=true
deploy: $(CMD) build
	$(CMD) deploy

.PHONY: clean
clean:
	$(call RM,$(TAR))

.PHONY: deps
deps:
ifeq ($(OS), Windows_NT)
	scoop install sass go watchexec deno
else
	brew install -q sass/sass/sass go watchexec deno
endif
# 	go install github.com/mitranim/afr@latest
