# Kitty-Beads Makefile
# Beads issue tracker with Spec Kitty dashboard UI

.PHONY: all build run clean deps test help

# Default target
all: build

# Build the server binary
build:
	@echo "Building kitty-beads server..."
	cd backend && go build -o ../bin/kitty-beads ./cmd/server

# Run the server (builds first if needed)
run: build
	@echo "Starting kitty-beads server..."
	./bin/kitty-beads -port 8080

# Run in development mode (backend only, serves pre-built frontend)
dev:
	@echo "Starting server on http://localhost:8080"
	cd backend && go run ./cmd/server -port 8080

# Run with frontend hot reload (backend + Vite dev server)
dev-hot:
	@echo "Starting backend on :8080 and frontend on :5173 (with hot reload)"
	@echo "Open http://localhost:5173 for hot reload"
	@trap 'kill 0' EXIT; \
		(cd backend && go run ./cmd/server -port 8080) & \
		(cd frontend && pnpm dev --host)

# Build frontend (run after making frontend changes)
build-frontend:
	@echo "Building frontend..."
	cd frontend && pnpm build
	@echo "Copying dist to server for embedding..."
	rm -rf backend/cmd/server/frontend-dist
	cp -r frontend/dist backend/cmd/server/frontend-dist
	@echo "Done! Refresh http://localhost:8080"

# Install dependencies
deps:
	cd backend && go mod download

# Tidy dependencies
tidy:
	cd backend && go mod tidy

# Run tests
test:
	cd backend && go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf backend/cmd/server/kitty-beads

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
	@echo "  make dev       - Run backend only (serves pre-built frontend)"
	@echo "  make dev-hot   - Run with frontend hot reload (use localhost:5173)"
	@echo "  make deps      - Download dependencies"
	@echo "  make tidy      - Tidy Go modules"
	@echo "  make test      - Run tests"
	@echo "  make build-frontend - Rebuild frontend after changes"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help"
	@echo ""
	@echo "Server runs on http://localhost:8080"
