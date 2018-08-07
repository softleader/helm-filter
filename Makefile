HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-filter/
HAS_GLIDE := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_PLUGIN_DIR)
	cp filter $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o filter -ldflags $(LDFLAGS)

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o filter -ldflags $(LDFLAGS)
	tar -zcvf $(DIST)/helm-filter-linux-$(VERSION).tgz filter README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o filter -ldflags $(LDFLAGS)
	tar -zcvf $(DIST)/helm-filter-macos-$(VERSION).tgz filter README.md LICENSE plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o filter.exe -ldflags $(LDFLAGS)
	tar -zcvf $(DIST)/helm-filter-windows-$(VERSION).tgz filter.exe README.md LICENSE plugin.yaml

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
