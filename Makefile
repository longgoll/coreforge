# ──────────────────────────────────────────────────────────────
# Forge CLI — Build & Release
# ──────────────────────────────────────────────────────────────

APP_NAME    := forge
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS     := -s -w -X 'github.com/longgoll/forge-cli/cmd.Version=$(VERSION)' -X 'github.com/longgoll/forge-cli/cmd.BuildDate=$(BUILD_DATE)'
BUILD_DIR   := dist

# ── Development ──────────────────────────────────────────────

.PHONY: build
build: ## Build for current OS
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME).exe .

.PHONY: run
run: build ## Build and run help
	./$(APP_NAME).exe --help

.PHONY: test
test: ## Run tests
	go test ./... -v

.PHONY: tidy
tidy: ## Tidy dependencies
	go mod tidy

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) $(APP_NAME).exe

# ── Cross-compile (all platforms) ────────────────────────────

.PHONY: build-all
build-all: clean ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(APP_NAME) $(VERSION) for all platforms..."
	@echo ""

	@echo "→ windows/amd64"
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

	@echo "→ windows/arm64"
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .

	@echo "→ linux/amd64"
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

	@echo "→ linux/arm64"
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .

	@echo "→ darwin/amd64 (macOS Intel)"
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

	@echo "→ darwin/arm64 (macOS Apple Silicon)"
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

	@echo ""
	@echo "✓ All builds complete! Check $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

# ── Checksums ────────────────────────────────────────────────

.PHONY: checksums
checksums: ## Generate SHA256 checksums
	cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@cat $(BUILD_DIR)/checksums.txt

# ── Help ─────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
