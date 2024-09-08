FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build

RUN apk add --no-cache git
WORKDIR /ec
COPY . /ec
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-X main.version=$(cat VERSION | tr -d '\n')" -o bin/editorconfig-checker ./cmd/editorconfig-checker/main.go

#

FROM alpine:latest

RUN apk add --no-cache git
RUN git config --global --add safe.directory /check
WORKDIR /check/
COPY --from=build /ec/bin/editorconfig-checker /usr/bin

CMD ["editorconfig-checker"]
