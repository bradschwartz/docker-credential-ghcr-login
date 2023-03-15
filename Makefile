# VERSION := $(shell cat VERSION)

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY_NAME=docker-credential-ghcr-login
LOCAL_BINARY=bin/local/$(BINARY_NAME)

.PHONY: build
build: $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o $(LOCAL_BINARY)
