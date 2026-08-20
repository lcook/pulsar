# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) Lewis Cook <lcook@FreeBSD.org>
.PHONY: build container container-push update test lint clean
.DELETE_ON_ERROR:

.DEFAULT_GOAL = build

VER = 0.2.0
PROGS = bot relay

GH_ACCOUNT = lcook
GH_PROJECT = pulsar

GIT_HASH = $(shell git rev-parse --short HEAD)
GIT_BRANCH = $(shell git symbolic-ref HEAD 2>/dev/null | sed 's,refs/heads/,,')
GIT_DIRTY = $(shell git status --porcelain)

ifeq ($(strip $(GIT_HASH)),)
GIT_HASH := unknown
IMAGE_TAG = latest
else
ifneq ($(strip $(GIT_DIRTY)),)
GIT_HASH := $(GIT_HASH)-dirty
endif
IMAGE_TAG = $(GIT_HASH)
endif

ifeq ($(strip $(GIT_BRANCH)),)
VER := $(VER)-$(GIT_HASH)
else
VER := $(GIT_BRANCH)/$(VER)-$(GIT_HASH)
endif

UNAME_S = $(shell uname -s)

ifeq ($(UNAME_S),FreeBSD)
PODMAN_ARGS = --network=host
endif
OCI_REPO ?= localhost
OCI_TAG = $(OCI_REPO)/$(GH_PROJECT):$(IMAGE_TAG)
ifneq ($(OCI_REPO),localhost)
OCI_TAG = $(OCI_REPO)/$(GH_ACCOUNT)/$(GH_PROJECT):$(IMAGE_TAG)
endif

GO_MODULE = github.com/$(GH_ACCOUNT)/$(GH_PROJECT)
GO_FLAGS = -v -ldflags "-s -w -X $(GO_MODULE)/internal/version.Build=$(VER)"

build: $(PROGS)

$(PROGS):
	@echo "|> Building $@@$(VER)"
	go build $(GO_FLAGS) -o $@ cmd/$(GH_PROJECT)-$@/$@.go

container:
	@for prog in $(PROGS); do \
		echo "|> Building $$prog@$(VER) container image"; \
		if [ "$(OCI_REPO)" != "localhost" ]; then \
			TAG_NAME=$$(echo "$(OCI_TAG)" | sed "s/$(GH_PROJECT)/$(GH_PROJECT)\/$$prog/"); \
		else \
			TAG_NAME=$$(echo "$(OCI_TAG)" | sed "s/$(GH_PROJECT)/$(GH_PROJECT)-$$prog/"); \
		fi; \
		podman build $(PODMAN_ARGS) --file container/$(UNAME_S)-$$prog --tag $$TAG_NAME .; \
	done

container-push:
	@for prog in $(PROGS); do \
		echo "|> Pushing $$prog@$(VER) to $(OCI_REPO)"; \
		if [ "$(OCI_REPO)" != "localhost" ]; then \
			TAG_NAME=$$(echo "$(OCI_TAG)" | sed "s/$(GH_PROJECT)/$(GH_PROJECT)\/$$prog/"); \
		else \
			TAG_NAME=$$(echo "$(OCI_TAG)" | sed "s/$(GH_PROJECT)/$(GH_PROJECT)-$$prog/"); \
		fi; \
		podman push $$TAG_NAME; \
	done

update:
	@echo "|> Updating and tidying up Go dependencies"
	go get -u -v ./...
	go mod tidy -v
	go mod verify

test:
	@echo "|> Running Go unit tests"
	go test -v -race -cover ./...

lint:
	@echo "|> Running linter on Go files"
	golangci-lint run

clean:
	@echo "|> Cleaning up project root directory"
	go clean
	rm -f $(PROGS)