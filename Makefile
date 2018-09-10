BASE = $(shell pwd)

export VERSION ?= dev
export GOCMD ?= go

.PHONY: all
all: kubar

.PHONY: kubar
kubar:
	@echo "GOPATH is ${GOPATH}"
	cd $(BASE)/kubar && $(GOCMD) build -ldflags "-X main.KubarBuildVersion=$(VERSION)"

.PHONY: install
install:
	@echo "GOPATH: ${GOPATH}"
	cd $(BASE)/kubar && $(GOCMD) install -ldflags "-X main.KubarBuildVersion=$(VERSION)"

.PHONY: clean
clean:
	@rm -f kubar/kubar
