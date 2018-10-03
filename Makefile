build:
	go build -o bin/ec src/main.go

run: build
	./bin/ec
