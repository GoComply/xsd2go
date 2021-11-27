GO=GO111MODULE=on go
GOBUILD=$(GO) build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

build:
	$(GOBUILD) ./cli/gocomply_xsd2go

.PHONY: pkger vendor
pkger:
ifeq ("$(wildcard $(GOBIN)/pkger)","")
	go get -u -v github.com/markbates/pkger/cmd/pkger
endif

ci-update-bundled-deps: pkger
	$(GOBIN)/pkger -o pkg/template
	go fmt ./pkg/template

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
