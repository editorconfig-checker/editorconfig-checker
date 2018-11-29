SRC_DIR = $(PWD)
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

install-deps:
	go get -u gopkg.in/editorconfig/editorconfig-core-go.v1

setup: install-deps build

bin/ec: $(SOURCES)
	@go build -o bin/ec ./cmd/editorconfig-checker/main.go

build: bin/ec

test:
	@go test -p=1 -cover -v ./...
	@go tool vet .
	@test -z $(shell gofmt -s -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

bench:
	go test -bench=. ./cmd/editorconfig-checker/

run: build
	@./bin/ec
