GO=GO111MODULE=on go
GOBUILD=$(GO) build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: vendor build

build:
	$(GOBUILD) ./cli/moov-io_xsd2go

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
