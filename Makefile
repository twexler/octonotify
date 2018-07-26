APP_NAME=octonotify
SHELL=bash
VERSION = $(shell bash scripts/git_version.sh)
GLIDE = $(GOPATH)/bin/glide
ICNSIFY = $(GOPATH)/bin/icnsify
GOBINDATA = $(GOPATH)/bin/go-bindata
APP_EXEC = build/$(APP_NAME)-$(VERSION)-$(shell uname | tr A-Z a-z)-$(shell uname -m)

ifeq ("$(APPVEYOR)","True")
APP_EXEC := build/$(APP_NAME)-$(VERSION)-windows-x64.exe
endif

APP := $(APP_EXEC)

ifeq ("$(TRAVIS_OS_NAME)","osx")
APP := build/$(APP_NAME).app.zip
endif

all: $(APP) 

icons/bindata.go: $(GOBINDATA) $(wildcard icons/*.png)
	$(GOBINDATA) -ignore='.*(svg|go)$$' -o $@ -pkg icons -prefix icons icons

build:
	mkdir $@

$(APP_EXEC): build $(wildcard */*.go) vendor icons/bindata.go
	go build -o $@ -ldflags='-X main.version=$(VERSION)' $(wildcard cmd/$(APP_NAME)/*.go)

build/$(APP_NAME).app: $(APP_EXEC) icons/octonotify.icns
	mkdir -p $@/Contents/MacOS
	cp $(APP_EXEC) $@/Contents/MacOS/
	mkdir -p $@/Contents/Resources
	cp icons/octonotify.icns $@/Contents/Resources
	go run -ldflags='-X main.version=$(VERSION)' scripts/genplist.go

build/$(APP_NAME).app.zip: build/$(APP_NAME).app
	pushd build/; zip -r $(shell basename $@) $(shell basename $<)

icons/octonotify.icns: $(ICNSIFY) icons/octonotify-small.png
	$(ICNSIFY) -i icons/octonotify-small.png -o $@

$(ICNSIFY):
	go get github.com/JackMordaunt/icns/cmd/icnsify

$(GOBINDATA):
	go get github.com/jteeuwen/go-bindata/...

$(GLIDE):
	go get github.com/Masterminds/glide

vendor: $(GLIDE)
	$(GLIDE) install

clean:
	rm -rf vendor build

.PHONY: clean
