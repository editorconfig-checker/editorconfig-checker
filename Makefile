SRC_DIR = $(PWD)
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

install-deps:
	go get -u gopkg.in/editorconfig/editorconfig-core-go.v1

setup: install-deps build

bin/ec: $(SOURCES)
	@go build -o bin/ec ./cmd/editorconfig-checker/main.go

build: bin/ec

test:
	@go test ./validators -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=github.com/editorconfig-checker/editorconfig-checker/validators/...
	@go test ./pkg/... -v -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=github.com/editorconfig-checker/editorconfig-checker/pkg/...
	@go tool vet .
	@test -z $(shell gofmt -s -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

bench:
	go test -bench=. ./**/*/

run: build
	@./bin/ec
