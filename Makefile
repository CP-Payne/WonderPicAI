# Makefile

# Variables
BINARY_NAME=app  # Or your project name
BINARY_PATH=tmp/$(BINARY_NAME)
GO_MAIN_PACKAGE=cmd/app/main.go

TAILWIND_INPUT=./tailwind/input.css
TAILWIND_OUTPUT=./static/css/style.css

# Phony targets (targets that don't represent files)
.PHONY: help run dev build clean install-tools css-build css-watch templ-generate templ-watch tidy all

# Default target (executed when you just run `make`)
all: build

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ------------------------------------------------------------------------------
# Development Tasks
# ------------------------------------------------------------------------------

run: build-dev ## Build for development and run the application.
	@echo "Running $(BINARY_PATH) (dev build)..."
	@$(BINARY_PATH)

dev: ## Run the application with live reloading using air (if installed).
	@echo "Starting application in development mode with air..."
	@echo "Ensure 'air' is installed (go install github.com/cosmtrek/air@latest) and .air.toml is configured."
	@go run -tags dev cmd/app/main.go

watch: ## Concurrently watch for CSS and Templ changes, and run the app with air.
	@echo "Starting watchers for CSS, Templ (temple includes hot reload for all .templ and .go files)..."
	@echo "Consider using a tool like 'overmind' or 'foreman' for managing multiple processes,"
	@echo "or run 'make css-watch', 'make templ-watch', and 'make dev' in separate terminals."
	@# This is a conceptual target. For true concurrent watching, you'd typically
	@# use a process manager like 'overmind', 'goreman', or run them in separate terminals.
	@# Example using a simple backgrounding (might not be ideal for all shells/OS):
	@# (make css-watch & make templ-watch & make dev)
	@echo "For a better experience, run 'make css-watch', 'make templ-watch', and 'make dev' in separate terminals."


# ------------------------------------------------------------------------------
# Build Tasks
# ------------------------------------------------------------------------------

build: css-build templ-generate ## Build the application for production.
	@echo "Building application for production..."
	@go build -o $(BINARY_PATH) $(GO_MAIN_PACKAGE)
	@echo "Production build complete: $(BINARY_PATH)"

build-dev: css-build templ-generate ## Build the application for development (includes dev tag).
	@echo "Building application for development..."
	@go build -tags dev -o $(BINARY_PATH) $(GO_MAIN_PACKAGE)
	@echo "Development build complete: $(BINARY_PATH)"


# ------------------------------------------------------------------------------
# Asset Generation Tasks
# ------------------------------------------------------------------------------

css-build: ## Build Tailwind CSS for production (minified).
	@echo "Building Tailwind CSS..."
	@(cd ./web && npx @tailwindcss/cli -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --minify)

css-watch: ## Watch Tailwind CSS files for changes and rebuild.
	@echo "Watching Tailwind CSS for changes... (Press Ctrl+C to stop)"
	@(cd ./web && npx @tailwindcss/cli -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --watch)

templ-generate: ## Generate Go code from Templ files.
	@echo "Generating Go code from Templ files..."
	@templ generate

templ-watch: ## Watch Templ files for changes and regenerate Go code.
	@echo "Watching Templ files for changes... (Press Ctrl+C to stop)"
	@# The --proxy flag is useful if you have a separate live-reloader for Go (like air)
	@# Adjust the port if your Go app runs on a different one.
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run -tags dev cmd/app/main.go" --open-browser=false  
	@# For ensuring tailindcss has finished generating before reloading browser
	@#templ generate --watch --proxy="http://localhost:8080" --cmd="cd web && npx @tailwindcss/cli -i ./tailwind/input.css -o ./static/css/style.css && cd .. && go run -tags dev cmd/app/main.go" --open-browser=false
	@# If not using a Go live reloader with proxy, just use:
	@# templ generate --watch


# ------------------------------------------------------------------------------
# Dependency Management & Installation
# ------------------------------------------------------------------------------

install-tools: ## Install necessary Go and Node development tools.
	@echo "Installing Go tools (templ CLI)..."
	@go install github.com/a-h/templ/cmd/templ@latest
	@echo "Installing Node.js dependencies (Tailwind CSS)..."
	@(cd ./web && npm install -D tailwindcss @tailwindcss/cli)
	@echo "Tools installation complete."
	@echo "Run 'make tidy' to fetch Go module dependencies."

tidy: ## Tidy Go module files and fetch dependencies.
	@echo "Tidying Go modules..."
	@go mod tidy
	@go mod download
	@echo "Go modules tidied."

# ------------------------------------------------------------------------------
# Cleanup Tasks
# ------------------------------------------------------------------------------

clean: ## Remove build artifacts and generated files.
	@echo "Cleaning up..."
	@rm -f $(BINARY_PATH)
	@rm -f $(TAILWIND_OUTPUT)
	@# Add other cleanup commands if necessary (e.g., removing templ_gen.go files if you don't commit them)
	@# rm -f ./web/template/*_templ.go
	@echo "Cleanup complete."
