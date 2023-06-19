MAKEFLAGS := --silent
PAR := $(MAKE) -j 128
TAR := public
CMD := ./bin/cmd
MOD := modules
SASS := sass --no-source-map -I $(MOD) styles/main.scss:$(TAR)/styles/main.css
GO_FLAGS := -tags=$(tags) -mod=mod
WATCH := watchexec -c -r -d=0 -n
W_CMD := --no-ignore -w=$(CMD)
W_GO := -e=go,mod

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
watch: clean cmd
	$(PAR) cmd_w srv pages_w styles_w cp_w

.PHONY: build
build: clean_tar
	$(PAR) styles pages cp

.PHONY: cmd_w
cmd_w: cmd
	$(WATCH) $(W_GO) -p -- $(MAKE) cmd

.PHONY: cmd
cmd: $(CMD)

$(CMD): *.go go.mod
	go build $(GO_FLAGS) -o $(CMD)

.PHONY: srv
srv: cmd
	$(CMD) srv

.PHONY: pages_w
pages_w:
	$(WATCH) $(W_CMD) -w=templates -- $(CMD) pages

.PHONY: pages
pages: cmd
	$(CMD) pages

.PHONY: games_steam_w
games_steam_w:
	$(WATCH) $(W_CMD) -w=misc/steam_apps.json -- $(CMD) games_steam

.PHONY: games_steam
games_steam: cmd
	$(CMD) games_steam

.PHONY: styles_w
styles_w:
	$(SASS) --style=$(SASS_STYLE) --watch

.PHONY: styles
styles:
	$(SASS) --style=$(SASS_STYLE)

.PHONY: test_w
test_w:
	$(WATCH) $(W_GO) -- $(MAKE) test

# The pattern `*_test.go` is needed here due to a bug/gotcha in Go's test
# runner.
.PHONY: test
test:
	go test -count=1 -failfast -short -run="$(run)" *_test.go

.PHONY: cp_w
cp_w:
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

.PHONY: clean_tar
clean_tar:
	$(call RM,$(TAR))

.PHONY: clean
clean: clean_tar
	$(call RM,$(CMD))

.PHONY: mod
mod:
	git submodule update --init --recursive --quiet

# Usage: `mod set commit=<hash>`.
.PHONY: mod_set
mod_set:
	cd $(MOD)/sb && git checkout $(commit)

.PHONY: deps
deps:
ifeq ($(OS), Windows_NT)
	scoop install sass go watchexec
else
	brew install -q sass/sass/sass go watchexec
endif
	$(MAKE) mod
