MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
GOBIN_PATH := $(PROJECT_PATH)/build/bin
PROJECTSEARCHPATH := $(PROJECT_PATH)/build/ext/bin
SEARCH_PATH := $(GOBIN_PATH):$(PROJECTSEARCHPATH):$(PATH)

export GOBIN := ${GOBIN_PATH}
export PATH := ${SEARCH_PATH}

all: ci

.PHONY: all ci lint shellcheck golangci-lint test setup clean

## ci: run all CI steps
ci: setup lint test

# lint: run all linting checks
lint: shellcheck golangci-lint

## shellscript: check all scripts with shellcheck
shellcheck: scripts/*
	shellcheck $<

## golangci-lint: run golang linter
golangci-lint:
	golangci-lint run ./...

## check: run checks
check: check-docs

## check-doc: check if the documentation is up to date
check-docs:
	gomarkdoc --check

## doc: update the documentation
doc:
	gomarkdoc

## test: run all tests
test:
	gotestsum ./...

## setup: download all external dependencies for a build on linux
setup:
	./scripts/setup_toolchain
	touch build/go.mod

## clean: delete all build artefacts
clean:
	rm -rf build
