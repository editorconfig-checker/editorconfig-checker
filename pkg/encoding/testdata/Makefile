
PHONY+=all
all: nuke download clean sum

DIRS = chardet encoding-test-files mimetype sdl text uchardet utf-8-test

chardet:
	git clone --depth 1 https://github.com/baulk/chardet chardet.tmp
	mkdir -p chardet
	cp -prv chardet.tmp/testdata/* chardet

encoding-test-files:
	git clone --depth 1 https://github.com/stain/encoding-test-files encoding-test-files.tmp
	mkdir -p encoding-test-files
	cp -prv encoding-test-files.tmp/*.txt encoding-test-files

mimetype:
	git clone https://github.com/gabriel-vasile/mimetype mimetype.tmp
	cd mimetype.tmp && \
		git checkout 2f8f7d81db79138423084a7fe1577efba28dcd29
	mkdir -p mimetype
	rm -f  mimetype.tmp/testdata/json*.txt
	cp -prv mimetype.tmp/testdata/*.txt mimetype.tmp/testdata/*.html mimetype

sdl:
	mkdir -p sdl.tmp
	wget -O sdl.tmp/utf8-sdl.txt https://raw.githubusercontent.com/libsdl-org/SDL/HEAD/test/utf8.txt
	mkdir -p sdl
	cp -prv sdl.tmp/*.txt sdl

text:
	git clone --depth 1 https://github.com/golang/text text.tmp
	mkdir -p text
	cp -prv text.tmp/encoding/testdata/* text

uchardet:
	git clone https://gitlab.freedesktop.org/uchardet/uchardet.git uchardet.tmp
	find uchardet.tmp -type f -name "*.c" -delete
	mkdir -p uchardet
	rm -f uchardet.tmp/CMakeLists.txt
	cp -prv uchardet.tmp/test/* uchardet

utf-8-test:
	mkdir -p utf-8-test
	cd utf-8-test && \
		wget -O UTF-8-demo.html https://www.w3.org/2001/06/utf-8-test/UTF-8-demo.html

PHONY+=download
download: $(DIRS)

PHONY+=clean
clean:
	# find -maxdepth 1 -type d -name "*.tmp" | xargs rm -fr

PHONY+=sum
sum:
	-find $(DIRS) -type f | xargs sha256sum >sha256sums.txt

PHONY+=nuke
nuke:
	find -maxdepth 1 -type d -name "*.tmp" | xargs rm -fr
	rm -fr $(DIRS)

.DEFAULT: all

.PHONY: $(PHONY)
