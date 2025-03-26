
PHONY+=default
default: help ;

DIRS = licenses chardet encoding-test-files mimetype sdl text uchardet utf-8-test wpt

licenses:
	mkdir -p licenses

chardet:
	git clone --depth 1 https://github.com/baulk/chardet chardet.tmp
	mkdir -p chardet
	cp -prv chardet.tmp/testdata/* chardet
	cp -pv chardet.tmp/LICENSE licenses/chardet-LICENSE

encoding-test-files:
	git clone --depth 1 https://github.com/stain/encoding-test-files encoding-test-files.tmp
	mkdir -p encoding-test-files
	cp -prv encoding-test-files.tmp/*.txt encoding-test-files
	cp -pv encoding-test-files.tmp/LICENSE licenses/encoding-test-files-LICENSE

mimetype:
	git clone https://github.com/gabriel-vasile/mimetype mimetype.tmp
	cd mimetype.tmp && \
		git checkout 2f8f7d81db79138423084a7fe1577efba28dcd29
	mkdir -p mimetype
	rm -f  mimetype.tmp/testdata/json*.txt
	cp -prv mimetype.tmp/testdata/*.txt mimetype.tmp/testdata/*.html mimetype
	cp -pv mimetype.tmp/LICENSE licenses/mimetype-LICENSE

sdl:
	git clone https://github.com/libsdl-org/SDL sdl.tmp
	mkdir -p sdl
	cp -prv sdl.tmp/test/utf8.txt sdl/utf8-sdl.txt
	cp -pv sdl.tmp/LICENSE.txt licenses/sdl-LICENSE.txt

text:
	git clone --depth 1 https://github.com/golang/text text.tmp
	mkdir -p text
	cp -prv text.tmp/encoding/testdata/* text
	cp -pv text.tmp/LICENSE licenses/text-LICENSE

uchardet:
	git clone --depth 1 https://gitlab.freedesktop.org/uchardet/uchardet.git uchardet.tmp
	find uchardet.tmp -type f -name "*.c" -delete
	mkdir -p uchardet
	rm -f uchardet.tmp/test/CMakeLists.txt
	cp -prv uchardet.tmp/test/* uchardet
	cp -prv uchardet.tmp/COPYING licenses/uchardet-COPYING

utf-8-test:
	mkdir -p utf-8-test.tmp
	wget -O utf-8-test.tmp/UTF-8-demo.html https://www.w3.org/2001/06/utf-8-test/UTF-8-demo.html
	mkdir -p utf-8-test
	cp -prv utf-8-test.tmp/* utf-8-test
	-wget -O licenses/utf-8-test-LICENSE.html https://www.w3.org/copyright/document-license-2023/

wpt:
	mkdir -p wpt.tmp
	git clone --depth 1 https://github.com/web-platform-tests/wpt wpt.tmp
	mkdir -p wpt
	cp -prv wpt.tmp/encoding/resources/ wpt
	cp -prv wpt.tmp/encoding/legacy-mb-japanese/ wpt
	cp -prv wpt.tmp/encoding/legacy-mb-korean/ wpt
	cp -prv wpt.tmp/encoding/legacy-mb-schinese/ wpt
	cp -prv wpt.tmp/encoding/legacy-mb-tchinese/ wpt
	find wpt -type f -not -name '*.*ml' -delete
	find wpt -type d -empty -delete
	cp -prv wpt.tmp/LICENSE.md licenses/wpt-LICENSE.md

PHONY+=all
all: nuke download clean sum ## Run `make nuke download clean sum`

PHONY+=nuke
nuke: ## Nuke all the downloaded files
	find -maxdepth 1 -type d -name "*.tmp" | xargs rm -fr
	rm -fr $(DIRS)

PHONY+=download
download: $(DIRS) ## Download all the files

PHONY+=clean
clean: ## Clean only the .tmp directories
	find -maxdepth 1 -type d -name "*.tmp" | xargs rm -fr

PHONY+=sum
sum: ## Generate the shasum256.txt file
	find $(DIRS) -type f | xargs sha256sum >sha256sums.txt

fmt=\033[36m%-10s\033[0m %s\n
PHONY+=help
help: ## Display a help screen
	@printf '$(fmt)' 'Target' 'Description'
	@printf '$(fmt)' '------' '------------------------------------------'
	@awk -F ':.*##[ \t]*' '/^[^#: \t]+:.*##/ {printf "$(fmt)", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: $(PHONY)
