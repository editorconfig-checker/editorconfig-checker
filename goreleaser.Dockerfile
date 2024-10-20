FROM alpine:latest

RUN apk add --no-cache git \
    && git config --global --add safe.directory /check
WORKDIR /check/

COPY editorconfig-checker /usr/bin/

# reasonale: I've see repackagers use COPY --from editorconfig-checker /usr/bin/ec /usr/bin/editorconfig-checker
# (found at https://github.com/super-linter/super-linter/blob/7b76efbd69ef471b83d5273d4b5d8b3cbd8e5e3f/Dockerfile#L335C34-L335C45)
RUN ln /usr/bin/editorconfig-checker /usr/bin/ec

CMD ["/editorconfig-checker"]
