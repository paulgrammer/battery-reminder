BINARY_NAME := battery-reminder
BUILD_DIR := $(PWD)/build
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
VERSION := 1.0.0
BUILD_TIME := $(shell date +%FT%T%z)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

.PHONY: build clean dev

dev: 
	go run main.go

build: build-darwin-amd64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_PATH)

clean:
	rm -f $(BINARY_PATH)-*