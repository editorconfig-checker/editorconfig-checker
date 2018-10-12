SRC_DIR = src
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

install-deps:
	go get -u gopkg.in/editorconfig/editorconfig-core-go.v1

setup: install-deps build

bin/ec: $(SOURCES)
	@go build -o bin/ec src/main.go

build: bin/ec

run: build
	@./bin/ec
