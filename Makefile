#!/usr/bin/make -fv

.ONESHELL:
SHELLFLAGS := -u nounset -ec

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

.PHONY: build build-env serve watch run

build: build-env
	echo $(GO) build $(GOFLAGS)
	for os in linux; do # darwin
	  GOOS=$${os} GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(APPNAME)-$${os}-$(GOARCH)
	done
	ln -svf "./ghostship-$${os}-$(GOARCH)" "./ghostship"

clean:
	rm -vf "$(APPNAME)"-*

tidy:
	$(GRC) go get -u ./...  # Upgrade all packages
	$(GRC) go mod tidy

build-env:
	$(GRC) go mod init $(APPNAME)
	$(GRC) go mod download

run: build
	bash -c "time ./$(APPNAME)"

test:
	$(GRC) go test ./... -v

watch:
	watcher

update: build-env tidy
	@:

version:
	./$(APPNAME) version

