FROM golang:1.16-alpine as build

RUN apk add --no-cache git
WORKDIR /ec
COPY . /ec
RUN GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-X main.version=$(cat VERSION | tr -d '\n')" -o bin/ec ./cmd/editorconfig-checker/main.go

#

FROM alpine:latest

RUN apk add --no-cache git
WORKDIR /check/
COPY --from=build /ec/bin/ec /usr/bin

CMD ["ec"]
