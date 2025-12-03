# ====== Config ======
GO        ?= go
BINARY    ?= wakyma-backend
MAIN_PKG  ?= ./...
BUILD_DIR ?= bin
PORT      ?= 8080

# Detect OS (Windows uses cmd.exe)
ifeq ($(OS),Windows_NT)
  SET_ENV = set PORT=$(PORT) &&
else
  SET_ENV = PORT=$(PORT)
endif

# ====== Default target ======
.PHONY: all
all: test build

# ====== Build & Run ======
.PHONY: build
build:
	@echo ">> Building binary..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY) .

.PHONY: run
run:
	@echo ">> Running app on port $(PORT)..."
	$(SET_ENV) $(GO) run .

.PHONY: dev
dev:
	@echo ">> Running in dev mode (auto-reload requires 'air' installed)..."
	$(SET_ENV) air

# ====== Code Quality ======
.PHONY: fmt
fmt:
	@echo ">> Formatting code..."
	$(GO) fmt $(MAIN_PKG)

.PHONY: tidy
tidy:
	@echo ">> Tidying modules..."
	$(GO) mod tidy

.PHONY: vet
vet:
	@echo ">> Running go vet..."
	$(GO) vet $(MAIN_PKG)

# Requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: lint
lint:
	@echo ">> Running golangci-lint..."
	golangci-lint run ./...

# ====== Tests ======
.PHONY: test
test:
	@echo ">> Running tests..."
	$(GO) test ./...

.PHONY: test-cover
test-cover:
	@echo ">> Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

# ====== Clean ======
.PHONY: clean
clean:
	@echo ">> Cleaning build artifacts..."
	$(GO) clean
	rm -rf $(BUILD_DIR) coverage.out

# ====== Help ======
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make run         - Run the app (PORT=$(PORT) by default)"
	@echo "  make dev         - Run with air (auto-reload, optional)"
	@echo "  make fmt         - Format Go code"
	@echo "  make tidy        - go mod tidy"
	@echo "  make vet         - go vet"
	@echo "  make lint        - Run golangci-lint (if installed)"
	@echo "  make test        - Run tests"
	@echo "  make test-cover  - Run tests with coverage"
	@echo "  make clean       - Clean artifacts"
