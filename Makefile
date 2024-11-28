# COLORS
TARGET_MAX_CHAR_NUM := 10
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# The binary to build (just the basename).
PWD := $(shell pwd)
NOW := $(shell date +%s)
APPNAME := k8s-ldap-auth
BIN ?= $(APPNAME)

ORG ?= registry.aegir.bouchaud.org
PKG := hopopops/$(APPNAME)
PLATFORM ?= "linux/arm/v7,linux/arm64/v8,linux/amd64"
GO ?= go
SED ?= sed
GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f)
GOVERSION := $(shell $(GO) version | $(SED) -r 's/go version go(.+)\s.+/\1/')
PACKAGES ?= $(shell $(GO) list ./...)

# This version-strategy uses git tags to set the version string
GIT_TAG := $(shell git describe --tags --always --dirty || echo unsupported)
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo unsupported)
BUILDTIME := $(shell date -u +"%FT%TZ%:z")
TAG ?= $(GIT_TAG)
VERSION ?= $(GIT_TAG)

.PHONY: fmt fmt-check vet test test-coverage cover install hooks docker tag push help clean dev
default: help

## Format go source code
fmt:
	$(GOFMT) -w $(GOFILES)

## Check if source code is formatted correctly
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

## Check source code for common errors
vet:
	$(GO) vet ${PACKAGES}

## Execute unit tests
test:
	$(GO) test ${PACKAGES}

## Execute unit tests & compute coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ${PACKAGES}

## Compute coverage
cover: test-coverage
	$(GO) tool cover -html=coverage.out

## Tidy dependencies
tidy:
	$(GO) mod tidy

## Install dependencies used for development
install: hooks tidy
	$(GO) mod download

## Install git hooks for post-checkout & pre-commit
hooks:
	@cp -f ./scripts/post-checkout .git/hooks/
	@cp -f ./scripts/pre-commit .git/hooks/
	@chmod +x .git/hooks/post-checkout
	@chmod +x .git/hooks/pre-commit

## Build the docker images
docker:
	@docker buildx build \
		--push \
		--build-arg COMMITHASH="$(GIT_COMMIT)" \
		--build-arg BUILDTIME="$(BUILDTIME)" \
		--build-arg VERSION="$(VERSION)" \
		--build-arg PKG="$(PKG)" \
		--build-arg APPNAME="$(APPNAME)" \
		--platform $(PLATFORM) \
		--tag $(ORG)/$(BIN):$(TAG) \
		--tag $(ORG)/$(BIN):latest \
		.

## Clean artifacts
clean:
	rm -f $(BIN) $(BIN)-dev $(BIN)-packed

$(APPNAME):
	$(GO) build \
		-trimpath \
		-buildmode=pie \
		-mod=readonly \
		-modcacherw \
		-o $(BIN) \
		-ldflags "\
				-X $(PKG)/version.APPNAME=$(APPNAME) \
				-X $(PKG)/version.VERSION=$(VERSION) \
				-X $(PKG)/version.GOVERSION=$(GOVERSION) \
				-X $(PKG)/version.BUILDTIME=$(BUILDTIME) \
				-X $(PKG)/version.COMMITHASH=$(GIT_COMMIT) \
				-s -w"

$(APPNAME)-dev:
	$(GO) build \
		-o $(BIN)-dev -ldflags "\
				-X $(PKG)/version.APPNAME=$(APPNAME) \
				-X $(PKG)/version.VERSION=$(VERSION) \
				-X $(PKG)/version.GOVERSION=$(GOVERSION) \
				-X $(PKG)/version.BUILDTIME=$(BUILDTIME) \
				-X $(PKG)/version.COMMITHASH=$(GIT_COMMIT)"

## Dev build outside of docker, not stripped
dev: $(APPNAME)-dev

$(APPNAME)-packed: $(APPNAME)
	upx --best $(APPNAME) -o $(APPNAME)-packed

## Release build outside of docker, stripped and packed
release: $(APPNAME)-packed

## Print this help message
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
