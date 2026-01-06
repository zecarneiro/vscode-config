SHELL := /bin/bash

ROOT := $(CURDIR)
BINARY_DIR := $(ROOT)/bin
PKG_DIR := $(ROOT)/pkg
VENDOR_DIR := $(ROOT)/vendor/golangutils
BINARY := $(BINARY_DIR)/vscode-config-linux
BINARY_WIN := $(BINARY_DIR)/vscode-config-win.exe
SRC := $(ROOT)/src
GO := go

.PHONY: all build clean check-deps process-go-dependencies list-go-dependencies help

all: help

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build                    Build windows and linux binaries"
	@echo "  clean                    Remove build outputs and downloads"
	@echo "  check-deps               Verify required tools (go)"
	@echo "  process-go-dependencies  Processo all go dependencies for this project"
	@echo "  list-go-dependencies     List all go dependencies for this project"

	@echo

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { echo "[ERROR] Please install golang!"; exit 1; }

process-go-dependencies: check-deps
	@echo ">>> Process for $(SRC)"
	@cd $(SRC) && $(GO) mod tidy
	@echo ">>> Process for $(VENDOR_DIR)"
	@cd $(VENDOR_DIR) && $(GO) mod tidy
	@cd $(ROOT)

list-go-dependencies: check-deps
	@echo ">>> List for $(SRC)"
	@cd $(SRC) && $(GO) list -m all
	@echo
	@echo ">>> List for $(VENDOR_DIR)"
	@cd $(VENDOR_DIR) && $(GO) list -m all
	@cd $(ROOT)

build: check-deps
	@mkdir -p $(BINARY_DIR)
	@echo "[INFO] Build WINDOWS app..."
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_WIN) $(SRC)
	@echo "[INFO] Build LINUX app..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY) $(SRC)

clean:
	@rm -rf $(BINARY_DIR)
	@rm -rf $(PKG_DIR)
	
