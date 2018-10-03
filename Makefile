SRC_DIR = src
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

bin/ec: $(SOURCES)
	@go build -o bin/ec src/main.go

build: bin/ec

run: build
	@./bin/ec
