BINARY_NAME := cmt
VERSION := $(shell node -p "require('./package.json').version")
BUILD_DIR := bin

TARGETS := darwin-arm64 darwin-amd64 linux-arm64 linux-amd64
LDFLAGS := -w -s -X main.AppVersion=$(VERSION) -X main.AppName=$(BINARY_NAME)

build: $(TARGETS)

$(TARGETS):
	GOOS=$(word 1,$(subst -, ,$@)) GOARCH=$(word 2,$(subst -, ,$@)) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$@ -ldflags "$(LDFLAGS)"

clean:
	rm -rf $(BUILD_DIR)

local:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags "$(LDFLAGS)"
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

version:
	@echo $(VERSION)

demo:
	@vhs demo.tape

.PHONY: build clean dev version demo $(TARGETS)
