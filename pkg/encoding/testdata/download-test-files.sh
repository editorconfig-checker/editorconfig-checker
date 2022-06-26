#!/usr/bin/env bash

rm -fr *.html *.txt mimetype text encoding-test-files chardet

git clone --depth 1 https://github.com/gabriel-vasile/mimetype
git clone --depth 1 https://github.com/golang/text/
git clone --depth 1 https://github.com/stain/encoding-test-files
git clone --depth 1 https://github.com/baulk/chardet

wget -O utf8-sdl.txt https://github.com/libsdl-org/SDL/blob/HEAD/test/utf8.txt

cp --force --backup=simple --suffix=-mimetype.html mimetype/testdata/*.html .

cp --force --backup=simple --suffix=-mimetype.txt mimetype/testdata/*.txt .

cp --force --backup=simple --suffix=-text.txt text/encoding/testdata/*.txt .

cp --force --backup=simple --suffix=-encoding-test-files.txt encoding-test-files/*.txt .

cp --force --backup=simple --suffix=-chardet.html chardet/testdata/*.html .

rm -f json.float.txt json.int.txt json.string.txt not.srt.2.txt not.srt.txt

sha256sum *.html *.txt | grep -v sha256sums\.txt >sha256sums.txt
