SRC_DIR = $(PWD)
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')
BINARIES = $(wildcard bin/*)
COMPILE_COMMAND = go build -o bin/ec ./cmd/editorconfig-checker/main.go

install-deps:
	go get -u gopkg.in/editorconfig/editorconfig-core-go.v1

setup: install-deps build

clean:
	rm ./bin/*

bin/ec: $(SOURCES)
	$(COMPILE_COMMAND)

build: bin/ec

test:
	@go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool vet .
	@test -z $(shell gofmt -s -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

bench:
	go test -bench=. ./**/*/

run: build
	@./bin/ec

run-verbose: build
	@./bin/ec --verbose

build-all-binaries: build run test clean
	# doesn't work on my machine and not in travis, see: https://github.com/golang/go/wiki/GoArm
	# GOOS=android GOARCH=arm  $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-android-arm
	# GOOS=darwin  GOARCH=arm $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm
	# GOOS=darwin  GOARCH=arm64 $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-arm64
	GOOS=darwin    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-386
	GOOS=darwin    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-darwin-amd64
	GOOS=dragonfly GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-dragonfly-amd64
	GOOS=freebsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-386
	GOOS=freebsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-amd64
	GOOS=freebsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-freebsd-arm
	GOOS=linux     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-386
	GOOS=linux     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-amd64
	GOOS=linux     GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-arm
	GOOS=linux     GOARCH=arm64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-arm64
	GOOS=linux     GOARCH=ppc64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-ppc64
	GOOS=linux     GOARCH=ppc64le  $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-ppc64le
	GOOS=linux     GOARCH=mips     $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips
	GOOS=linux     GOARCH=mipsle   $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mipsle
	GOOS=linux     GOARCH=mips64   $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips64
	GOOS=linux     GOARCH=mips64le $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-linux-mips64le
	GOOS=netbsd    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-386
	GOOS=netbsd    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-amd64
	GOOS=netbsd    GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-netbsd-arm
	GOOS=openbsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-386
	GOOS=openbsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-amd64
	GOOS=openbsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-openbsd-arm
	GOOS=plan9     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-plan9-386
	GOOS=plan9     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-plan9-amd64
	GOOS=solaris   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-solaris-amd64
	GOOS=windows   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-windows-386
	GOOS=windows GOARCH=amd64 $(COMPILE_COMMAND) && mv ./bin/ec ./bin/ec-windows-amd64

compress-all-binaries: build-all-binaries
	for f in $(BINARIES); do      \
		tar czf $$f.tar.gz $$f;    \
	done
	@rm $(BINARIES)
