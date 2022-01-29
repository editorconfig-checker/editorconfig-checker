SRC_DIR = $(PWD)
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')
BINARIES = $(wildcard bin/*)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_BRANCH_UP_TO_DATE = $(shell git remote show origin | tail -n1 | sed 's/.*(\(.*\))/\1/')
CURRENT_VERSION = $(shell cat VERSION | tr -d '\n')
COMPILE_COMMAND = go build -ldflags "-X main.version=$(CURRENT_VERSION)" -o bin/ec ./cmd/editorconfig-checker/main.go

prefix = /usr/local
bindir = /bin
mandir = /share/man

all: build

clean:
	rm -f ./bin/*

bin/ec: $(SOURCES) VERSION
	$(COMPILE_COMMAND)

build: bin/ec

install: build
	install -D bin/ec $(DESTDIR)$(prefix)$(bindir)/editorconfig-checker
	install -D docs/editorconfig-checker.1 $(DESTDIR)$(prefix)$(mandir)/man1/editorconfig-checker.1

uninstall:
	rm -f $(DESTDIR)$(prefix)$(bindir)/editorconfig-checker
	rm -f $(DESTDIR)$(prefix)$(mandir)/man1/editorconfig-checker.1

test:
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go vet ./...
	@test -z $(shell gofmt -s -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

bench:
	go test -bench=. ./**/*/

run: build
	@./bin/ec

run-verbose: build
	@./bin/ec --verbose

release: _is_master_branch _git_branch_is_up_to_date _do_release
	@echo Release done. Go to Github and create a release.

_is_master_branch:
ifneq ($(GIT_BRANCH),master)
	@echo You are not on the master branch.
	@echo Please check out the master and try to release again
	@false
endif

_git_branch_is_up_to_date:
ifneq ($(GIT_BRANCH_UP_TO_DATE),up to date)
	@echo Your master branch is not up to date.
	@echo Please push your changes or pull changes from the remote.
	@false
endif

current_version:
	@echo the current version is: $(CURRENT_VERSION)

_do_release: _checkout clean _tag_version build run _build-all-binaries _compress-all-binaries
	git checkout master
	git merge --no-ff release
	git push origin master && git push origin master --tags

_checkout:
	git branch -D release; git checkout -b release

_tag_version: current_version
	@read -p "Enter version to release: " version && \
	echo $${version} > VERSION && \
	sed -i "s/VERSION=".*"/VERSION=\"$${version}\"/" ./README.md && \
	sed -i "s/\"Version\": \".*\",/\"Version\": \"$${version}\",/" .ecrc && \
	sed -i "s/\"Version\": \".*\",/\"Version\": \"$${version}\",/" testfiles/generated-config.json && \
	sed -i "s/${CURRENT_VERSION}/$${version}/" ./pkg/config/config_test.go && \
	git add . && git commit -m "chore(release): $${version}" && git tag "$${version}"

_build-all-binaries:
	# doesn't work on my machine and not in travis, see: https://github.com/golang/go/wiki/GoArm
	# GOOS=android GOARCH=arm  $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-android-arm
	# GOOS=darwin  GOARCH=arm $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm
	GOOS=darwin  GOARCH=arm64 $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm64
	CGO_ENABLED=0 GOOS=darwin    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin    GOARCH=arm64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm64
	CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-dragonfly-amd64
	CGO_ENABLED=0 GOOS=freebsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-386
	CGO_ENABLED=0 GOOS=freebsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-amd64
	CGO_ENABLED=0 GOOS=freebsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-arm
	CGO_ENABLED=0 GOOS=linux     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-386
	CGO_ENABLED=0 GOOS=linux     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-amd64
	CGO_ENABLED=0 GOOS=linux     GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-arm
	CGO_ENABLED=0 GOOS=linux     GOARCH=arm64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-arm64
	CGO_ENABLED=0 GOOS=linux     GOARCH=ppc64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-ppc64
	CGO_ENABLED=0 GOOS=linux     GOARCH=ppc64le  $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-ppc64le
	CGO_ENABLED=0 GOOS=linux     GOARCH=mips     $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips
	CGO_ENABLED=0 GOOS=linux     GOARCH=mipsle   $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mipsle
	CGO_ENABLED=0 GOOS=linux     GOARCH=mips64   $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips64
	CGO_ENABLED=0 GOOS=linux     GOARCH=mips64le $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips64le
	CGO_ENABLED=0 GOOS=netbsd    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-386
	CGO_ENABLED=0 GOOS=netbsd    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-amd64
	CGO_ENABLED=0 GOOS=netbsd    GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-arm
	CGO_ENABLED=0 GOOS=openbsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-386
	CGO_ENABLED=0 GOOS=openbsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-amd64
	CGO_ENABLED=0 GOOS=openbsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-arm
	CGO_ENABLED=0 GOOS=plan9     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-plan9-386
	CGO_ENABLED=0 GOOS=plan9     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-plan9-amd64
	CGO_ENABLED=0 GOOS=solaris   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-solaris-amd64
	CGO_ENABLED=0 GOOS=windows   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-windows-386.exe
	CGO_ENABLED=0 GOOS=windows   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-windows-amd64.exe

_compress-all-binaries:
	for f in $(BINARIES); do      \
		tar czf $$f.tar.gz $$f;    \
	done
	@rm $(BINARIES)

_release_dockerfile: _build_dockerfile _push_dockerfile

_build_dockerfile:
	docker build -t mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/') .

_push_dockerfile:
	docker push mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/')

nix-build:
	nix-build -E 'with import <nixpkgs> { }; callPackage ./default.nix {}'

nix-install:
	nix-env -i -f install.nix

nix-update-dependencies:
	nix-shell -p vgo2nix --run vgo2nix

.PHONY: clean build install uninstall test bench run run-verbose release _is_master_branch _git_branch_is_up_to_date current_version _do_release _tag_version _build-all-binaries _compress-all-binaries _release_dockerfile _build_dockerfile _push_dockerfile nix-build nix-install nix-update-dependencies
