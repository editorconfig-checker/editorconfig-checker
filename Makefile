ifeq ($(OS),Windows_NT)
	STDERR=con
else
	STDERR=/dev/stderr
endif

ifeq ($(wildcard "C:/Program Files/Git/usr/bin/*"),)
export PATH:="C:/Program Files/Git/usr/bin:$(PATH)"
endif

EXE=bin/ec$(EXEEXT)
SRC_DIR := $(shell dirname "$(realpath "$(firstword $(MAKEFILE_LIST))")")
SOURCES = $(shell find "$(SRC_DIR)" -type f -name "*.go")
CURRENT_VERSION = $(shell cat VERSION | tr -d "\n")

prefix = /usr/local
bindir = /bin
mandir = /share/man

all: build ## Build bin/ec

clean: ## Clean bin/ directory
	rm -f ./bin/*

define _build
go build -ldflags "-X main.version=$(CURRENT_VERSION)" -o $1 ./cmd/editorconfig-checker/main.go
endef

$(EXE): $(SOURCES) VERSION
	$(call _build,$(EXE))

build: $(EXE) ## Build bin/ec

install: build ## Build and install executable in PATH
	install -D $(EXE) $(DESTDIR)$(prefix)$(bindir)/editorconfig-checker
	install -D docs/editorconfig-checker.1 $(DESTDIR)$(prefix)$(mandir)/man1/editorconfig-checker.1

uninstall: ## Remove executable from PATH
	rm -f $(DESTDIR)$(prefix)$(bindir)/editorconfig-checker
	rm -f $(DESTDIR)$(prefix)$(mandir)/man1/editorconfig-checker.1

test: ## Run test suite
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go vet ./...
	@test -z $(shell gofmt -s -l . | tee $(STDERR)) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

bench: ## Run benchmark
	go test -bench=. ./**/*/

run: build ## Build and run bin/ec
	@./bin/ec --exclude "\\.git" --exclude "\\.exe$$"

run-verbose: build ## Build and run bin/ec --verbose
	@./bin/ec --verbose --exclude "\\.git" --exclude "\\.exe$$"

nix-build: ## Build for nix
	nix-build -E 'with import <nixpkgs> { }; callPackage ./default.nix {}'

nix-install: ## Install on nix
	nix-env -i -f install.nix

nix-update-dependencies: ## Update nix dependencies
	nix-shell -p vgo2nix --run vgo2nix

PHONY += help
help: ## Display available commands
	@awk -F ':.*##[ \t]*' '/^[^#: \t]+:.*##/ {printf "\033[36m%-23s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

PHONY += dumpvars
dumpvars: ## Dump variables
	@echo CURRENT_VERSION=$(CURRENT_VERSION)
	@echo EXE=$(EXE)
	@echo EXEEXT=$(EXEEXT)
	@echo STDERR=$(STDERR)
	@echo SRC_DIR=$(SRC_DIR)
	@echo SOURCES=$(SOURCES)

.PHONY: $(PHONY) clean build install uninstall test bench run run-verbose release _is_main_branch _git_branch_is_up_to_date current_version _do_release _tag_version _build-all-binaries _compress-all-binaries _release_dockerfile _build_dockerfile _push_dockerfile nix-build nix-install nix-update-dependencies
