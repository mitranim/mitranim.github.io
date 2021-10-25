MAKEFLAGS := --silent
PAR := $(MAKE) -j 128
TAR := public
BIN := ./bin
CMD := $(BIN)/cmd
SASS := sass --no-source-map -I submodules styles/main.scss:$(TAR)/styles/main.css
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
	$(PAR) styles_w cmd_w srv pages_w cp_w

.PHONY: build
build: clean
	$(PAR) styles pages cp

.PHONY: styles_w
styles_w:
	$(SASS) --style=$(SASS_STYLE) --watch

.PHONY: styles
styles:
	$(SASS) --style=$(SASS_STYLE)

.PHONY: cmd_w
cmd_w: $(CMD)
	$(WATCH) -e=go,mod -- $(MAKE) $(CMD)

$(CMD): *.go go.mod
	go build -o $(CMD)

.PHONY: srv
srv: $(CMD)
	$(CMD) srv

.PHONY: pages_w
pages_w: pages
	$(WATCH_CMD) -w=templates -- $(CMD) pages

.PHONY: pages
pages: $(CMD)
	$(CMD) pages

.PHONY: cp_w
cp_w: cp
	$(WATCH) -w=static -w=images -- $(MAKE) cp

.PHONY: cp
cp:
	$(call MKDIR,$(TAR))
	$(call MKDIR,$(TAR)/images)
	$(call CP,static,$(TAR))
	$(call CP,images,$(TAR)/images)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	$(call RM,$(TAR))

sub:
	git submodule update --init --recursive --quiet

.PHONY: deps
deps:
ifeq ($(OS), Windows_NT)
	scoop install sass go watchexec
else
	brew install -q sass/sass/sass go watchexec
endif
	$(MAKE) sub
