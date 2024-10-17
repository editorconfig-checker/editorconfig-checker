FROM scratch
ENTRYPOINT ["/editorconfig-checker"]
COPY editorconfig-checker /
# we currently increase the image size by providing the binary in the old place as well,
# but I fail to find a way to provide a real binary as /usr/bin/ec
# reasonale: I've see repackagers use COPY --from editorconfig-checker /usr/bin/ec /usr/bin/editorconfig-checker
# (found at https://github.com/super-linter/super-linter/blob/7b76efbd69ef471b83d5273d4b5d8b3cbd8e5e3f/Dockerfile#L335C34-L335C45)
COPY editorconfig-checker /usr/bin/ec
