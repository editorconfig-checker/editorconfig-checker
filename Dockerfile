FROM golang:1.12.4-alpine as build

RUN apk add --no-cache git
WORKDIR /ec
COPY . /ec
RUN GO111MODULE=on CGO_ENABLED=0 go build -o bin/ec cmd/editorconfig-checker/main.go

#

FROM alpine:latest

RUN apk add --no-cache git
WORKDIR /check/
COPY --from=build /ec/bin/ec /

CMD ["/ec"]
