# This Dockerfile is built using goreleaser, which provides the ready built binary in the context.
# With that we only need to copy the binary into the expected location.
# Doing it this way with goreleaser allows us to not build the binary again, but reuse the binaries built already.

FROM alpine:latest

ARG TARGETPLATFORM

RUN apk add --no-cache git \
    && git config --global --add safe.directory /check
WORKDIR /check/

COPY $TARGETPLATFORM/editorconfig-checker /usr/bin/

CMD ["/usr/bin/editorconfig-checker"]
