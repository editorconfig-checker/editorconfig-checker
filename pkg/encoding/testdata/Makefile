
PHONY+=default
default: help ;

DIRS = licenses chardet encoding-test-files mimetype sdl text uchardet utf-8-test wpt

licenses:
	mkdir -p $@

chardet:
	git clone --sparse --depth 1 https://github.com/baulk/chardet $@.tmp
	cd $@.tmp && \
		git sparse-checkout set testdata
	mkdir -p $@
	cp -prv $@.tmp/testdata/* $@
	cp -pv $@.tmp/LICENSE licenses/chardet-LICENSE

encoding-test-files:
	git clone --sparse --depth 1 https://github.com/stain/encoding-test-files $@.tmp
	mkdir -p $@
	cp -prv $@.tmp/*.txt $@
	cp -pv $@.tmp/LICENSE licenses/encoding-test-files-LICENSE

mimetype:
	git clone --sparse https://github.com/gabriel-vasile/mimetype $@.tmp
	cd $@.tmp && \
		git sparse-checkout set testdata && \
		git checkout 2f8f7d81db79138423084a7fe1577efba28dcd29
	mkdir -p $@
	rm -f  $@.tmp/testdata/json*.txt
	cp -prv $@.tmp/testdata/*.txt $@.tmp/testdata/*.html $@
	cp -pv $@.tmp/LICENSE licenses/mimetype-LICENSE

sdl:
	git clone --sparse --depth 1 https://github.com/libsdl-org/SDL $@.tmp
	cd $@.tmp && \
		git sparse-checkout set test
	mkdir -p $@
	cp -prv $@.tmp/test/utf8.txt $@/utf8-sdl.txt
	cp -pv $@.tmp/LICENSE.txt licenses/sdl-LICENSE.txt

text:
	git clone --sparse --depth 1 https://github.com/golang/text $@.tmp
	cd $@.tmp && \
		git sparse-checkout set encoding/testdata
	mkdir -p $@
	cp -prv $@.tmp/encoding/testdata/* $@
	cp -pv $@.tmp/LICENSE licenses/text-LICENSE

uchardet:
	git clone --sparse --depth 1 https://gitlab.freedesktop.org/uchardet/uchardet.git $@.tmp
	cd $@.tmp && \
		git sparse-checkout set test
	find $@.tmp -type f -name "*.c" -delete
	find $@.tmp -name 'CMakeLists.txt' -delete
	mkdir -p $@
	cp -prv $@.tmp/test/* $@
	cp -prv $@.tmp/COPYING licenses/uchardet-COPYING

utf-8-test:
	mkdir -p $@.tmp
	wget -O $@.tmp/UTF-8-demo.html https://www.w3.org/2001/06/utf-8-test/UTF-8-demo.html
	mkdir -p $@
	cp -prv $@.tmp/* $@
	-wget -O licenses/utf-8-test-LICENSE.html https://www.w3.org/copyright/document-license-2023/

wpt:
	mkdir -p $@.tmp
	git clone --sparse --depth 1 https://github.com/web-platform-tests/wpt $@.tmp
	cd $@.tmp && \
		git sparse-checkout set encoding/
	mkdir -p $@
	cd $@.tmp && \
		cp -prv \
		encoding/resources/ \
		encoding/legacy-mb-japanese/ \
		encoding/legacy-mb-korean/ \
		encoding/legacy-mb-schinese/ \
		encoding/legacy-mb-tchinese/ \
		../$@
	find $@ -type f -not -name '*.*ml' -delete
	find $@ -type d -empty -delete
	cp -prv $@.tmp/LICENSE.md licenses/wpt-LICENSE.md

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
