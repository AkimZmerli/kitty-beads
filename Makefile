# Kitty-Beads Makefile
# Beads issue tracker with Spec Kitty dashboard UI

.PHONY: all build run clean deps test help

# Default target
all: build

# Build the server binary
build:
	@echo "Building kitty-beads server..."
	cd src && go build -o ../bin/kitty-beads ./cmd/server

# Run the server (builds first if needed)
run: build
	@echo "Starting kitty-beads server..."
	./bin/kitty-beads -port 8080

# Run in development mode with auto-reload
dev:
	@echo "Starting development server..."
	cd src && go run ./cmd/server -port 8080

# Install dependencies
deps:
	cd src && go mod download

# Tidy dependencies
tidy:
	cd src && go mod tidy

# Run tests
test:
	cd src && go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf src/cmd/server/kitty-beads

# Initialize beads in current directory (for testing)
init-beads:
	@if [ ! -d ".beads" ]; then \
		echo "Initializing .beads directory..."; \
		mkdir -p .beads; \
		echo "Run 'bd init' after installing beads CLI to fully initialize"; \
	else \
		echo ".beads already exists"; \
	fi

# Clean up upstream repos (optional, saves space)
clean-upstream:
	rm -rf beads-upstream spec-kitty-upstream

# Show help
help:
	@echo "Kitty-Beads - Beads + Spec Kitty UI"
	@echo ""
	@echo "Usage:"
	@echo "  make build     - Build the server binary"
	@echo "  make run       - Build and run the server"
	@echo "  make dev       - Run in development mode"
	@echo "  make deps      - Download dependencies"
	@echo "  make tidy      - Tidy Go modules"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help"
	@echo ""
	@echo "Server runs on http://localhost:8080 by default"
