#!/usr/bin/make -fv

.ONESHELL:
SHELL      := /bin/bash
SHELLFLAGS := -u nounset -ec


MAKEFILE  := $(realpath $(lastword $(MAKEFILE_LIST)))
MAKE      := make
MAKEFLAGS += --no-print-directory
MAKEFLAGS += --warn-undefined-variables

APPNAME   := ghostship
MODNAME   := github.com/shalomb/$(APPNAME)
VERSION   := $(shell git describe --tags --long --always)
GO				:= $(shell command -v go)
GOVERSION := $(shell go version | awk '{ print $$3 }')
GOOS      := $(shell go version | awk '{ split($$4, a, "/"); print a[1]  }' )
GOARCH    := $(shell go version | awk '{ split($$4, a, "/"); print a[2]  }' )
GITBRANCH := $(shell git branch --show-current)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
BUILDHOST := $(shell hostname -f)

THIS_MAKEFILE := $(realpath $(lastword $(MAKEFILE_LIST)))
THIS_DIR      := $(shell dirname $(THIS_MAKEFILE))
THIS_PROJECT  := $(APPNAME)

GOLDFLAGS += -X main.AppName=$(APPNAME)
GOLDFLAGS += -X main.Branch=$(GITBRANCH)
GOLDFLAGS += -X main.BuildHost=$(BUILDHOST)
GOLDFLAGS += -X main.BuildTime=$(BUILDTIME)
GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.GoVersion=$(GOVERSION)
GOLDFLAGS += -X main.GoOS=$(GOOS)
GOLDFLAGS += -X main.GoArch=$(GOARCH)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

GRC := $(shell command -v grc)

UNAME_S := $(shell uname -s)
TARGET :=
ifeq ($(UNAME_S),Linux)
	TARGET := "$(APPNAME)-linux"
endif
ifeq ($(UNAME_S),Darwin)
	TARGET := "$(APPNAME)"
endif

# https://dustinrue.com/2021/08/parameters-in-a-makefile/
# setup arguments
RUN_ARGS          := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# ...and turn them into do-nothing targets
$(eval $(RUN_ARGS):;@:)

.PHONY: audit build build-env serve watch run test
.SILENT: clean test

# https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects

audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose .

build:
	@:
	for os in linux; do # darwin
	  GOOS=$${os} GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(APPNAME)-$${os}-$(GOARCH)
	done
	ln -sf "./ghostship-$${os}-$(GOARCH)" "./ghostship"

build-env:
	test -f go.mod || $(GRC) go mod init $(APPNAME)
	$(GRC) go mod download

clean:
	rm -f "$(APPNAME)"-* coverage.*

deploy: run
	install -v -m775 ./"$(APPNAME)" ~/.local/bin/"$(APPNAME)"
	hash -r

lint:
	golangci-lint run --enable-all

run: build
	@:
	source <(./"$(APPNAME)" init bash);
	./"$(APPNAME)" version
	time ./"$(APPNAME)" prompt --terminal-width $${COLUMNS:-80} --status 0 --pipestatus 0 --cmd-duration 1 --prompt-character "T" >/dev/null

test:
	$(GRC) go test -race -buildvcs ./... $(RUN_ARGS) # allow 'make test -- -v'

test-benchmark:
	$(GRC) go test -bench=. -benchmem ./... $(RUN_ARGS)

test-coverage:
	$(GRC) go test -coverprofile coverage.out ./... $(RUN_ARGS)
	$(GRC) go tool cover -func coverage.out
	# $(GRC) go tool cover -html coverage.out -o coverage.html

tidy: build-env
	$(GRC) go get -u ./...  # Upgrade all packages
	$(GRC) go mod tidy

watch:
	watcher

update: build-env tidy
	@:

version:
	./$(APPNAME) version
