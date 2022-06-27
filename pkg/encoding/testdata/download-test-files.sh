#!/usr/bin/env bash

rm -fr *.html *.txt mimetype text encoding-test-files chardet

git clone --depth 1 https://github.com/gabriel-vasile/mimetype
git clone --depth 1 https://github.com/golang/text/
git clone --depth 1 https://github.com/stain/encoding-test-files
git clone --depth 1 https://github.com/baulk/chardet

wget -O utf8-sdl.txt https://raw.githubusercontent.com/libsdl-org/SDL/HEAD/test/utf8.txt

cp -v --force --backup=simple --suffix=-mimetype.html mimetype/testdata/*.html .

cp -v --force --backup=simple --suffix=-mimetype.txt mimetype/testdata/*.txt .

cp -v --force --backup=simple --suffix=-text.txt text/encoding/testdata/*.txt .

cp -v --force --backup=simple --suffix=-encoding-test-files.txt encoding-test-files/*.txt .

cp -v --force --backup=simple --suffix=-chardet.html chardet/testdata/*.html .

rm -f json*.txt not*.txt

sha256sum *.html *.txt | grep -v sha256sums\.txt >sha256sums.txt
