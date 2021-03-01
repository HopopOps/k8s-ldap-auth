# COLORS
TARGET_MAX_CHAR_NUM := 10
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# The binary to build (just the basename).
PWD := $(shell pwd)
NOW := $(shell date +%s)
BIN := ldap-auth

ORG ?= registry.aegir.bouchaud.org
NAMESPACE := legion/kubernetes
PKG := bouchaud.org/${NAMESPACE}/${BIN}
GO ?= go
GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f)
GOVERSION := 1.13
PACKAGES ?= $(shell $(GO) list ./...)

# This version-strategy uses git tags to set the version string
GIT_TAG := $(shell git describe --tags --always --dirty || echo unsupported)
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo unsupported)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")
BUILDTIME := $(shell date -u +"%FT%TZ%:z")
ARCH := $(shell uname -m)
TAG ?= $(GIT_TAG)

.PHONY: fmt fmt-check vet test test-coverage cover install hooks docker-amd64 docker-arm64 docker tag-amd64 tag-arm64 tag push-amd64 push-arm64 push manifest manifest-push dev all help
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

## Install dependencies used for development
install: hooks
	$(GO) mod download

## Install git hooks for post-checkout & pre-commit
hooks:
	@cp -f ./scripts/post-checkout .git/hooks
	@cp -f ./scripts/pre-commit .git/hooks
	@chmod +x .git/hooks/post-checkout
	@chmod +x .git/hooks/pre-commit

# base image are removed between each build for both amd64 and arm64:
# https://github.com/moby/moby/issues/36552#issuecomment-538061565 -
# it's necessary for multiarch build if we want to be able to select the
# correct arch for base image since both arch use the same tag

## Build the amd64 docker image
docker-amd64:
	@docker rmi --force --no-prune $(shell grep FROM Dockerfile | sed -E 's/from ([^ ]+).*/\1/mig')
	@docker build \
		--pull \
		--platform "linux/amd64" \
		--build-arg COMMITHASH="$(GIT_COMMIT)" \
		--build-arg BUILDTIME="$(BUILDTIME)" \
		--build-arg VERSION="$(GIT_TAG)" \
		--build-arg PKG="$(PKG)" \
		--build-arg APPNAME="$(BIN)" \
		--build-arg GOVERSION="$(GOVERSION)" \
		--tag "$(ORG)/$(BIN):amd64-latest" \
		.

## Build the arm64 docker image
docker-arm64:
	@docker rmi --force --no-prune $(shell grep FROM Dockerfile | sed -E 's/from ([^ ]+).*/\1/mig')
ifeq ($(ARCH),x86_64)
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
else
	@echo "current arch is $(ARCH), expect the unexpected"
endif
	@docker build \
		--pull \
		--platform "linux/arm64" \
		--build-arg COMMITHASH="$(GIT_COMMIT)" \
		--build-arg BUILDTIME="$(BUILDTIME)" \
		--build-arg VERSION="$(GIT_TAG)" \
		--build-arg PKG="$(PKG)" \
		--build-arg APPNAME="$(BIN)" \
		--build-arg GOVERSION="$(GOVERSION)" \
		--tag "$(ORG)/$(BIN):arm64-latest" \
		.

## Build the docker images
docker: docker-amd64 docker-arm64

## Tag the latest docker image to current version, arm64
tag-arm64: docker-arm64
	@docker tag "$(ORG)/$(BIN):arm64-latest" "$(ORG)/$(BIN):arm64-$(TAG)"

## Tag the latest docker image to current version, amd64
tag-amd64: docker-amd64
	@docker tag "$(ORG)/$(BIN):amd64-latest" "$(ORG)/$(BIN):amd64-$(TAG)"

## Tag both amd64 and arm64 images
tag: tag-arm64 tag-amd64

## Push the latest and current version tags to registry, arm64
push-arm64: tag-arm64
	@docker push "$(ORG)/$(BIN):arm64-latest"
	@docker push "$(ORG)/$(BIN):arm64-$(TAG)"

## Push the latest and current version tags to registry, amd64
push-amd64: tag-amd64
	@docker push "$(ORG)/$(BIN):amd64-latest"
	@docker push "$(ORG)/$(BIN):amd64-$(TAG)"

## Push both amd64 and arm64 images
push: push-arm64 push-amd64

## Create the manifests for the current and latest version
manifest: push
	@docker manifest create "$(ORG)/$(BIN):latest" --amend "$(ORG)/$(BIN):arm64-latest" --amend "$(ORG)/$(BIN):amd64-latest"
	@docker manifest create "$(ORG)/$(BIN):$(TAG)" --amend "$(ORG)/$(BIN):arm64-$(TAG)" --amend "$(ORG)/$(BIN):amd64-$(TAG)"

## Push the manifests to the registry
manifest-push: manifest
	@docker manifest push "$(ORG)/$(BIN):latest" --purge
	@docker manifest push "$(ORG)/$(BIN):$(TAG)" --purge

## Dev build outside of docker, not stripped
dev:
	$(GO) build \
		-o $(BIN) -ldflags "\
				-X $(PKG)/version.APPNAME=$(BIN) \
				-X $(PKG)/version.VERSION=$(GIT_TAG) \
				-X $(PKG)/version.GOVERSION=$(GOVERSION) \
				-X $(PKG)/version.BUILDTIME=$(BUILDTIME) \
				-X $(PKG)/version.COMMITHASH=$(GIT_COMMIT)"

## Same as manifest-push
all: manifest-push

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
