APP=octonotify
VERSION=$(shell git describe)
GLIDE=$(GOPATH)/bin/glide
ICNSIFY=$(GOPATH)/bin/icnsify
GOBINDATA=$(GOPATH)/bin/go-bindata


all: build/$(APP)

icons/bindata.go: $(GOBINDATA) $(wildcard icons/*.png)
	$(GOBINDATA) -ignore='.*(svg|go)$$' -o $@ -pkg icons -prefix icons icons

build/$(APP): $(wildcard */*.go) vendor icons/bindata.go 
	go build -o $@ -ldflags='-X main.version=$(VERSION)' $(wildcard cmd/$(APP)/*.go)

build/$(APP).app: build/$(APP) icons/octonotify.icns
	mkdir -p $@/Contents/MacOS
	cp build/$(APP) $@/Contents/MacOS/
	mkdir -p $@/Contents/Resources
	cp icons/octonotify.icns $@/Contents/Resources
	go run -ldflags='-X main.version=$(VERSION)' scripts/genplist.go

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
	rm -rf $(APP) vendor

.PHONY: clean
