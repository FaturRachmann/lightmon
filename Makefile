BINARY=lightmon
VERSION=2.0.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION) -X main.goos=$(GOOS) -X main.goarch=$(GOARCH)"

.PHONY: all build clean run tidy install config help

all: build

## Show this help
help:
	@echo "lightmon build targets:"
	@echo "  make build        - Build for current OS"
	@echo "  make build-all    - Build for Linux amd64 + macOS arm64 + macOS amd64"
	@echo "  make install      - Build and install to /usr/local/bin"
	@echo "  make config       - Generate default config file"
	@echo "  make run          - Run without installing"
	@echo "  make tidy         - Download dependencies"
	@echo "  make clean        - Remove binaries"
	@echo "  make help         - Show this help"

## Download dependencies
tidy:
	go mod tidy

## Build for current OS
build: tidy
	go build $(LDFLAGS) -o $(BINARY) .

## Build for Linux (amd64)
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-linux-amd64 .

## Build for macOS (arm64 / Apple Silicon)
build-mac-arm:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY)-darwin-arm64 .

## Build for macOS (amd64)
build-mac-intel:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-darwin-amd64 .

## Build all platforms
build-all: build-linux build-mac-arm build-mac-intel

## Generate default config file
config:
	@mkdir -p ~/.lightmon
	@cp config.example.yaml ~/.lightmon/config.yaml
	@echo "Default config generated at ~/.lightmon/config.yaml"

## Run without building
run:
	go run . --interval 1s

## Install to /usr/local/bin
install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed to /usr/local/bin/$(BINARY)"

## Remove binaries
clean:
	rm -f $(BINARY) $(BINARY)-*
