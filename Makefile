BIN = askii

ARCH ?= $(shell go env GOOS)-$(shell go env GOARCH)

OUTPUT_DIR ?= _oputput

# Go environment
platform = $(subst -, ,$(ARCH))
GOOS = $(word 1, $(platform))
GOARCH = $(word 2, $(platform))
GOPROXY ?= "https://proxy.golang.org,direct"

.PHONY: all
all:
	@$(MAKE) build

build: _output/bin/$(GOOS)/$(GOARCH)/$(BIN)

_output/bin/$(GOOS)/$(GOARCH)/$(BIN): build-dirs
	@echo "Building $(BIN) for $(ARCH)..."
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	BIN=$(BIN) \
	OUTPUT_DIR=$$(pwd)/_output/bin/$(GOOS)/$(GOARCH) \
	go build -o _output/bin/$(GOOS)/$(GOARCH)/$(BIN) main.go
	@echo "Build complete: $(OUTPUT_DIR)/bin/$(GOOS)/$(GOARCH)/$(BIN)"

build-dirs:
	@mkdir -p _output/bin/$(GOOS)/$(GOARCH)

.PHONY: clean
clean:
	rm -rf _output

.PHONY: install
install: build
	cp _output/bin/$(GOOS)/$(GOARCH)/$(BIN) /usr/local/bin/$(BIN)
	@echo "$(BIN) installed to /usr/local/bin/$(BIN)"

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/$(BIN)
	@echo "$(BIN) uninstalled from /usr/local/bin/$(BIN)"

