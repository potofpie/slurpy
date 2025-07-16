.PHONY: build clean test example cli run-example

# Build CLI
build:
	go build -o slurpy ./cli

# Build and run CLI
cli: build
	./slurpy

# Run basic example to generate test data
example:
	go run ./examples/basic

# Run example then CLI
run-example: example
	@echo "Requests have been logged. Starting CLI..."
	@echo "Use ↑/↓ or j/k to navigate, tab to switch panels, q to quit"
	./slurpy

# Clean built artifacts and logs
clean:
	rm -f slurpy
	rm -rf ~/.config/slurpy/logs/*

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy

# Show help
help:
	@echo "Slurpy - HTTP Request Logger & Debugger"
	@echo ""
	@echo "Commands:"
	@echo "  make build       - Build the CLI"
	@echo "  make cli         - Build and run CLI"
	@echo "  make example     - Run basic example to generate test data"
	@echo "  make run-example - Run example then start CLI"
	@echo "  make clean       - Clean built artifacts and logs"
	@echo "  make test        - Run tests"
	@echo "  make deps        - Install dependencies"
	@echo "  make help        - Show this help" 