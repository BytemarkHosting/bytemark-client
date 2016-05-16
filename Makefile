SHELL:=/bin/bash
PKGBASE := github.com/BytemarkHosting/bytemark-client
CHOCOBASE := ports/chocolatey/package
ALL_PACKAGES := $(PKGBASE)/lib $(PKGBASE)/cmds/util $(PKGBASE)/cmds $(PKGBASE)/cmd/bytemark
ALL_SOURCE := lib/*.go mocks/*.go util/*/*.go cmd/**/*.go
TAR_FILES := bytemark doc/bytemark.1
ZIP_FILES := bytemark.exe doc/bytemark.pdf

BUILD_NUMBER ?= 0

LAUNCHER_APP:=ports/mac/launcher.app
RGREP=grep -rn --color=always --exclude=.* --exclude-dir=Godeps --exclude=Makefile

.PHONY: test update-dependencies
.PHONY: bytemark-client.nupkg
.PHONY: find-uk0 find-bugs-todos find-exits
.PHONY: gensrc

all: bytemark 

bytemark-client.zip: $(ZIP_FILES)
	zip $@ $^

bytemark-client.tar.gz: $(TAR_FILES)
	tar czf $@ $^

bytemark-client.deb: $(ALL_SOURCE) cmd/bytemark/debian/*
	cd cmd/bytemark && fakeroot debian/rules binary

bytemark-client.changes: cmd/bytemark/debian/control cmd/bytemark/debian/changes
	cd cmd/bytemark && dpkg-genchanges -b > ../../bytemark-client.changes

bytemark-client.nupkg: VERSION
	cd ports/chocolatey && make VERSION=$(<VERSION)
	mv $(CHOCOBASE)/bytemark.nupkg bytemark-client.nupkg

%.pdf: %.ps
	ps2pdf $< $@

doc/bytemark-client.ps: doc/bytemark.1
	groff -mandoc -T ps $< > $@

bytemark.exe: bytemark
	mv bytemark bytemark.exe

bytemark: $(ALL_SOURCE) gensrc
	GO15VENDOREXPERIMENT=1 go build -o bytemark $(PKGBASE)/cmd/bytemark

install-jessie-golang:
ifeq ($(build_distribution),"jessie")
	echo "deb http://mirror.bytemark.co.uk/debian jessie-backports main" | sudo tee -a /etc/apt/sources.list
	sudo apt-get update
	sudo apt-get install -y --force-yes -t jessie-backports golang-go
endif


install-build-deps: install-jessie-golang
	sudo apt-get install git golang-go-$(GOOS)-$(GOARCH)

install-dpkg-build-deps: install-build-deps
	sudo apt-get install -y --force-yes dh-golang 

#install-rpm-build-deps:

install-windows-build-deps: install-build-deps
	sudo apt-get install -y --force-yes curl wget unzip mono-runtime libmono-system-core4.0-cil libmono-system-componentmodel-dataannotations4.0-cil libmono-windowsbase4.0-cil libmono-system-xml-linq4.0-cil ghostscript

# make changelog opens vim to update the changelog
# then generates a new version.go file.
changelog:
	gen/changelog.sh
	make gensrc

clean:
	rm -rf Bytemark.app rm $(LAUNCHER_APP)
	rm -f bytemark bytemark.exe
	rm -f bytemark-client.zip bytemark-client.tar
	rm -f main.coverage lib.coverage
	rm -f main.coverage.html lib.coverage.html

gensrc:
	BUILD_NUMBER=$(BUILD_NUMBER) go generate ./...

install: bytemark doc/bytemark.1
	cp bytemark /usr/bin/bytemark
	cp doc/bytemark.1 /usr/share/man/man1

coverage: lib.coverage.html main.coverage.html
ifeq (Darwin, $(shell uname -s))
	open lib.coverage.html
	open main.coverage.html
	open cmds.coverage.html
else
	xdg-open lib.coverage.html
	xdg-open main.coverage.html
	xdg-open cmds.coverage.html
endif

main.coverage: cmd/bytemark/*.go
	go test -coverprofile=$@ $(PKGBASE)/cmd/bytemark

util.coverage: cmd/bytemark/util/*.go
	go test -coverprofile=$@ $(PKGBASE)/cmd/bytemark/util

%.coverage.html: %.coverage
	go tool cover -html=$< -o $@

%.coverage: % %/*
	go test -coverprofile=$@ $(PKGBASE)/$<

#docs: doc/*.md
#	for file in doc/*.md; do \
#	    pandoc --from markdown --to html $$file --output $${file%.*}.html; \
#	done

test: gensrc
ifdef $(VERBOSE)
	GO15VENDOREXPERIMENT=1 go test -v $(ALL_PACKAGES)
else 
	GO15VENDOREXPERIMENT=1 go test $(ALL_PACKAGES)
endif

find-uk0: 
	$(RGREP) --exclude=bytemark "uk0" .

find-bugs-todos:
	$(RGREP) -P "// BUG(.*):" . || echo ""
	$(RGREP) -P "// TODO(.*):" .

find-exits:
	$(RGREP) --exclude=exit.go --exclude=main.go -P "panic\(|os.Exit" .
