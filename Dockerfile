# This Dockerfile is built using goreleaser, which provides the ready built binary in the context.
# With that we only need to copy the binary into the expected location.
# Doing it this way with goreleaser allows us to not build the binary again, but reuse the binaries built already.

FROM alpine:latest

RUN apk add --no-cache git \
    && git config --global --add safe.directory /check
WORKDIR /check/

COPY editorconfig-checker /usr/bin/

# Rationale: We observed repackagers use COPY --from editorconfig-checker /usr/bin/ec /usr/bin/editorconfig-checker
# (found at https://github.com/super-linter/super-linter/blob/7b76efbd69ef471b83d5273d4b5d8b3cbd8e5e3f/Dockerfile#L335C34-L335C45)
RUN ln /usr/bin/editorconfig-checker /usr/bin/ec

CMD ["/usr/bin/editorconfig-checker"]
