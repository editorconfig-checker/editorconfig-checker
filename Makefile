SRC_DIR = $(PWD)
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')
BINARIES = $(wildcard bin/*)
COMPILE_COMMAND = go build -o bin/ec ./cmd/editorconfig-checker/main.go
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_BRANCH_UP_TO_DATE = $(shell git remote show origin | tail -n1 | sed 's/.*(\(.*\))/\1/')
CURRENT_VERSION = $(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/')

clean:
	rm ./bin/*

bin/ec: $(SOURCES)
	$(COMPILE_COMMAND)

build: bin/ec

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

release: _is_master_branch _git_branch_is_up_to_date current_version _tag_version _do_release _release_dockerfile
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

_do_release: clean test build run _build-all-binaries _compress-all-binaries

_tag_version:
	@read -p "Enter version to release: " version && \
	sed -i "s/const version string = \".*\"/const version string = \"$${version}\"/" ./cmd/editorconfig-checker/main.go && \
	git add . && git commit -m "chore(release): $${version}" && git tag "$${version}" && \
	git push origin master && git push origin master --tags

_build-all-binaries:
	# doesn't work on my machine and not in travis, see: https://github.com/golang/go/wiki/GoArm
	# GOOS=android GOARCH=arm  $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-android-arm
	# GOOS=darwin  GOARCH=arm $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm
	# GOOS=darwin  GOARCH=arm64 $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm64
	CGO_ENABLED=0 GOOS=darwin    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-386
	CGO_ENABLED=0 GOOS=darwin    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-amd64
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

docker-build:
	docker  build -t mstruebing/editorconfig-checker:latest .

docker-run:
	docker run --rm --volume=$$PWD:/check mstruebing/editorconfig-checker:1.1.3

_release_dockerfile: _build_dockerfile _push_dockerfile

_build_dockerfile:
	docker build -t mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/') .

_push_dockerfile:
	docker push mstruebing/editorconfig-checker:$(shell grep 'const version' cmd/editorconfig-checker/main.go | sed 's/.*"\(.*\)"/\1/')


