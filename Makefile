# What to build.
BUILD_TARGET := cmd/posto_ipirangad/*.go
# The built binary names (just the basename).
BIN := posto_ipirangad

# Tooling images
BUILD_IMAGE ?= golang:1.17.3-bullseye
GOLANGCI_LINT_IMAGE ?= golangci/golangci-lint:v1.42.1-alpine

# Where to push the docker image.
REGISTRY ?= 726369664624.dkr.ecr.ap-northeast-1.amazonaws.com/posto_ipiranga-api
IMAGE_NAME := posto_ipiranga-api

# This version-strategy uses a manual value to set the version string
# Versions should be generated sourcing from Conventional Commits using standard-version by CI
VERSION ?= x

###
### These variables should not need tweaking.
###
# Platforms to build when the user runs all-cross-build
ALL_PLATFORMS := linux/amd64 linux/arm linux/arm64 linux/ppc64le linux/s390x darwin/amd64 windows/amd64

# $SRC_DIRS/... will be used for `go test`
SRC_DIRS := . # directories which hold app source (not vendored)

# Used internally.  Users should pass GOOS and/or GOARCH.
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

IMAGE := $(REGISTRY)/$(IMAGE_NAME)
TAG := $(VERSION)

# If you want to cross build all binaries, see the 'all-cross-build' rule.
.PHONY: all
all: build

# For the following OS/ARCH expansions, we transform OS/ARCH into OS_ARCH
# because make pattern rules don't match with embedded '/' characters.

build-%:
	@$(MAKE) build \
	    --no-print-directory \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

ci-build-%:
	@$(MAKE) ci-build \
	    --no-print-directory \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

##@ Build

.PHONY: build
build: bin/$(OS)_$(ARCH)/$(BIN) ## Build the application in a containerized environment, ${GOOS}_${GOARCH} can be provided after a dash, e.g. build-linux_amd64

# Directories that we need created to build/test.
BUILD_DIRS := bin/$(OS)_$(ARCH) \
              .go/bin/$(OS)_$(ARCH) \
              .go/cache

# The following structure defeats Go's (intentional) behavior to always touch
# result files, even if they have not changed.  This will still run `go` but
# will not trigger further work if nothing has actually changed.
OUTBIN = bin/$(OS)_$(ARCH)/$(BIN)
$(OUTBIN): .go/$(OUTBIN).stamp
	@true

# This will build the binary under ./.go and update the real binary iff needed.
.PHONY: .go/$(OUTBIN).stamp
.go/$(OUTBIN).stamp: $(BUILD_DIRS)
	@echo "making $(OUTBIN)"
	@docker run \
	    -i \
	    --rm \
	    -u $$(id -u):$$(id -g) \
	    -v $$(pwd):/src \
	    -w /src \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH) \
	    -v $$(pwd)/.go/cache:/.cache \
	    --env HTTP_PROXY=$(HTTP_PROXY) \
	    --env HTTPS_PROXY=$(HTTPS_PROXY) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c " \
	        ARCH=$(ARCH) \
	        OS=$(OS) \
	        BINARY_OUTPUT_ROOT=/go/bin \
			BUILD_TARGET=$(BUILD_TARGET) \
			BINARY_NAME=$(BIN) \
	        VERSION=$(VERSION) \
	        ./scripts/build.sh \
	    "
	@if ! cmp -s .go/$(OUTBIN) $(OUTBIN); then \
	    mv .go/$(OUTBIN) $(OUTBIN); \
	    date >$@; \
	fi

ci-build: ## Build the application in the current environment (without container), ${GOOS}_${GOARCH} can be provided after a dash, e.g. ci-build-linux_amd64
	@echo "making $(OUTBIN)"
	@sh -c " \
		ARCH=$(ARCH) \
		OS=$(OS) \
		BINARY_OUTPUT_ROOT=bin \
		BUILD_TARGET=$(BUILD_TARGET) \
		BINARY_NAME=$(BIN) \
		VERSION=$(VERSION) \
		./scripts/build.sh \
	"

##@ Cross-build

all-cross-build: $(addprefix build-, $(subst /,_, $(ALL_PLATFORMS))) ## Cross-build the application for multiple platforms in a containerized environment

ci-all-cross-build: $(addprefix ci-build-, $(subst /,_, $(ALL_PLATFORMS))) ## Cross-build the application for multiple platforms in the current environment (without container)

##@ Lint
.PHONY: lint
lint: $(BUILD_DIRS) ## Lint the application with golangci-lint
	@docker run \
		-i \
		--rm \
		-u $$(id -u):$$(id -g) \
		-v $$(pwd):/src \
		-w /src \
		-v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin \
		-v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH) \
		-v $$(pwd)/.go/cache:/.cache \
		--env HTTP_PROXY=$(HTTP_PROXY) \
		--env HTTPS_PROXY=$(HTTPS_PROXY) \
		$(GOLANGCI_LINT_IMAGE) golangci-lint run

##@ Test

.PHONY: test
test: $(BUILD_DIRS) ## Test the application in a containerized environment
	@docker run \
	    -i \
	    --rm \
	    -u $$(id -u):$$(id -g) \
	    -v $$(pwd):/src \
	    -w /src \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH) \
	    -v $$(pwd)/.go/cache:/.cache \
	    --env HTTP_PROXY=$(HTTP_PROXY) \
	    --env HTTPS_PROXY=$(HTTPS_PROXY) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c " \
	        ARCH=$(ARCH) \
	        OS=$(OS) \
	        VERSION=$(VERSION) \
	        ./scripts/test.sh $(SRC_DIRS) \
	    "

ci-test: ## Test the application in the current environment (without container)
	@echo "making $(OUTBIN)"
	@sh -c " \
		ARCH=$(ARCH) \
		OS=$(OS) \
		VERSION=$(VERSION) \
		./scripts/test.sh $(SRC_DIRS) \
	"

##@ Container

container: Dockerfile ## Build the container
	@docker build \
		-t $(IMAGE):$(TAG) \
		--build-arg VERSION=$(VERSION) \
		.
	@echo "container: $(IMAGE):$(TAG)"

push-container: ## Push the container to the registry
	@docker push $(IMAGE):$(TAG)
	@echo "pushed: $(IMAGE):$(TAG)"

$(BUILD_DIRS):
	@mkdir -p $@

##@ Cleanup

.PHONY: clean
clean: ## Clean the built products
	rm -rf .go bin

##@ Misc.

version: ## Print the application version
	@echo $(VERSION)

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nCashflow Statement Excel Renderer Makefile\n\n\033[1mUsage\033[0m:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
