APP=octonotify
VERSION=$(shell git describe)
GLIDE=$(GOPATH)/bin/glide

all: $(APP)

$(APP): vendor $(wildcard */*.go)
	go build -o $@ -ldflags='-X main.version=$(VERSION)' $(wildcard cmd/*.go)

$(GLIDE):
	go get github.com/Masterminds/glide

vendor: $(GLIDE)
	$(GLIDE) install

clean:
	rm -rf $(APP) vendor

.PHONY: clean