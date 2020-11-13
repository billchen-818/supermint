.PHONY: build install docker-image

SUPERMINT_VERSION?=latest
BUILD_TAGS?='supermint'
VERSION := $(shell git describe --always)
BUILD_FLAGS = -ldflags "-X github.com/vbhp/supermintION)

.DEFAULT_GOAL := build

build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -tags $(BUILD_TAGS) -o ../../build/tm-signer-harness main.go

install:
	CGO_ENABLED=0 go install $(BUILD_FLAGS) -tags $(BUILD_TAGS) .

docker-image:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -tags $(BUILD_TAGS) -o tm-signer-harness main.go
	docker build \
		--build-arg SUPERMINT_VERSION=$(SUPERMINT_VERSION) \
		-t supermint/tm-signer-harness:$(SUPERMINT_VERSION) .
	rm -rf tm-signer-harness
