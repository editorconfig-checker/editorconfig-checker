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
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_BRANCH_UP_TO_DATE = $(shell git remote show origin | tail -n1 | sed "s/.*(\(.*\))/\1/")
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

release: _is_main_branch _git_branch_is_up_to_date _do_release ## Create release
	@echo Release done. Go to Github and create a release.

_is_main_branch:
ifneq ($(GIT_BRANCH),main)
	@echo You are not on the main branch.
	@echo Please check out the main and try to release again
	@false
endif

_git_branch_is_up_to_date:
ifneq ($(GIT_BRANCH_UP_TO_DATE),up to date)
	@echo Your main branch is not up to date.
	@echo Please push your changes or pull changes from the remote.
	@false
endif

current_version: ## Display current version
	@echo the current version is: $(CURRENT_VERSION)

_do_release: _checkout clean _tag_version build run _build-all-binaries _compress-all-binaries
	git checkout main
	git merge --no-ff release
	git push origin main && git push origin main --tags

_checkout:
	git branch -D release; git checkout -b release

_tag_version: current_version
	@read -p "Enter version to release: " version && \
	[[ $$version == v* ]]  && \ # check that there is leading v in the version
	echo $${version} > VERSION && \
	sed -i "s/VERSION=".*"/VERSION=\"$${version}\"/" ./README.md && \
	sed -i "s/\"Version\": \".*\",/\"Version\": \"$${version}\",/" .ecrc && \
	sed -i "s/\"Version\": \".*\",/\"Version\": \"$${version}\",/" testfiles/generated-config.json && \
	sed -i "s/${CURRENT_VERSION}/$${version}/" ./pkg/config/config_test.go && \
	git add . && git commit -m "chore(release): $${version}" && git tag "$${version}"

define _build_target
CGO_ENABLED=0 GOOS=$(subst /,,$(dir $1)) GOARCH=$(notdir $1) $(call _build,bin/ec-$(subst /,,$(dir $1))-$(notdir $1)$2)
endef

_build-all-binaries:
	@# To generate list of supported targets, run:
	@#  go tool dist list
	$(call _build_target,aix/ppc64)
	@# gcc.exe: error: unrecognized command-line option '-rdynamic'
	@# $(call _build_target,android/386)
	@# gcc.exe: error: unrecognized command-line option '-rdynamic'
	@# $(call _build_target,android/amd64)
	@# gcc.exe: error: unrecognized command-line option '-marm'; did you mean '-mabm'?
	@# gcc.exe: error: unrecognized command-line option '-rdynamic'
	@# $(call _build_target,android/arm)
	$(call _build_target,android/arm64)
	$(call _build_target,darwin/amd64)
	$(call _build_target,darwin/arm64)
	$(call _build_target,dragonfly/amd64)
	$(call _build_target,freebsd/386)
	$(call _build_target,freebsd/amd64)
	$(call _build_target,freebsd/arm)
	$(call _build_target,freebsd/arm64)
	$(call _build_target,illumos/amd64)
	@# gcc.exe: error: unrecognized command-line option '-framework'
	@# $(call _build_target,ios/amd64)
	@# gcc.exe: error: unrecognized command-line option '-framework'
	@# $(call _build_target,ios/arm64)
	$(call _build_target,js/wasm)
	$(call _build_target,linux/386)
	$(call _build_target,linux/amd64)
	$(call _build_target,linux/arm)
	$(call _build_target,linux/arm64)
	$(call _build_target,linux/mips)
	$(call _build_target,linux/mips64)
	$(call _build_target,linux/mips64le)
	$(call _build_target,linux/mipsle)
	$(call _build_target,linux/ppc64)
	$(call _build_target,linux/ppc64le)
	$(call _build_target,linux/riscv64)
	$(call _build_target,linux/s390x)
	$(call _build_target,netbsd/386)
	$(call _build_target,netbsd/amd64)
	$(call _build_target,netbsd/arm)
	$(call _build_target,netbsd/arm64)
	$(call _build_target,openbsd/386)
	$(call _build_target,openbsd/amd64)
	$(call _build_target,openbsd/arm)
	$(call _build_target,openbsd/arm64)
	$(call _build_target,openbsd/mips64)
	$(call _build_target,plan9/386)
	$(call _build_target,plan9/amd64)
	$(call _build_target,plan9/arm)
	$(call _build_target,solaris/amd64)
	$(call _build_target,windows/386,.exe)
	$(call _build_target,windows/amd64,.exe)
	$(call _build_target,windows/arm,.exe)
	$(call _build_target,windows/arm64,.exe)

_compress-all-binaries:
	for f in bin/*; do      \
		tar czf $$f.tar.gz $$f;    \
		rm -f $$f;    \
	done

_release_dockerfile: _build_dockerfile _push_dockerfile

_build_dockerfile:
	docker build -t mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/') .

_push_dockerfile:
	docker push mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/')

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
	@echo GIT_BRANCH=$(GIT_BRANCH)
	@echo GIT_BRANCH_UP_TO_DATE=$(GIT_BRANCH_UP_TO_DATE)
	@echo STDERR=$(STDERR)
	@echo SRC_DIR=$(SRC_DIR)
	@echo SOURCES=$(SOURCES)

.PHONY: $(PHONY) clean build install uninstall test bench run run-verbose release _is_main_branch _git_branch_is_up_to_date current_version _do_release _tag_version _build-all-binaries _compress-all-binaries _release_dockerfile _build_dockerfile _push_dockerfile nix-build nix-install nix-update-dependencies
