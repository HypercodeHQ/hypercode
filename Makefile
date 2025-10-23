.DEFAULT_GOAL := help

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-air: ## Install Air for hot-reloading during development
	@echo "installing Air..."
	@go install github.com/air-verse/air@latest
	@echo "Air installed successfully!"
	@echo "Run 'make dev' to start the server with hot-reloading"

build-css: ## Build and minify CSS from Tailwind
	@mkdir -p public/assets
	bunx @tailwindcss/cli -i ./views/styles/main.css -o ./public/assets/styles.css

build: build-css ## Build production binary to bin/hypercode
	@echo "building hypercode binary with embedded assets..."
	@mkdir -p bin
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/hypercode ./cmd/server
	@echo "binary built successfully at bin/hypercode"

dev: ## Start development server with hot-reloading
	@command -v air >/dev/null 2>&1 || { echo "Air is not installed. Run 'make install-air' first."; exit 1; }
	air

clean: ## Clean build artifacts and temporary files
	go clean
	rm -rf bin/
	rm -rf tmp/
	rm -f hypercode.db
	rm -f build-errors.log
