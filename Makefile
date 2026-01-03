# Linux Packet Visualizer Makefile

.PHONY: all dev generate frontend build install clean help

# Default target
all: dev

# Development: generate contract and start dev server
dev: generate frontend

# Generate JSON contract from Go code
generate:
	@echo "📦 Generating contract JSON..."
	@go run ./cmd/contract -o frontend/public/data/egress_path.json
	@echo "✅ Contract written to frontend/public/data/egress_path.json"

# Start frontend development server
frontend:
	@echo "🚀 Starting frontend dev server..."
	@cd frontend && npm run dev

# Build production frontend
build: generate
	@echo "🏗️  Building production frontend..."
	@cd frontend && npm run build
	@echo "✅ Production build complete in frontend/dist/"

# Install all dependencies
install: install-go install-frontend

install-go:
	@echo "📥 Checking Go modules..."
	@go mod download
	@go mod verify
	@echo "✅ Go dependencies installed"

install-frontend:
	@echo "📥 Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "✅ Frontend dependencies installed"

# Clean generated files
clean:
	@echo "🧹 Cleaning generated files..."
	@rm -f frontend/public/data/egress_path.json
	@rm -rf frontend/dist
	@rm -rf frontend/node_modules/.vite
	@echo "✅ Clean complete"

# Run Go tests
test:
	@echo "🧪 Running Go tests..."
	@go test -v ./...

# Build Go binary
build-cli:
	@echo "🔧 Building CLI binary..."
	@go build -o bin/contract ./cmd/contract
	@echo "✅ Binary built: bin/contract"

# Lint Go code
lint:
	@echo "🔍 Linting Go code..."
	@go vet ./...
	@echo "✅ Lint complete"

# Format Go code
fmt:
	@echo "✨ Formatting Go code..."
	@go fmt ./...
	@echo "✅ Format complete"

# Preview production build
preview: build
	@echo "👀 Previewing production build..."
	@cd frontend && npm run preview

# Generate and validate contract
validate: generate
	@echo "🔍 Validating contract JSON..."
	@cat frontend/public/data/egress_path.json | jq '.version' 
	@cat frontend/public/data/egress_path.json | jq '.paths | length'
	@echo "✅ Contract is valid"

# Show help
help:
	@echo "Linux Packet Visualizer - Make Targets"
	@echo ""
	@echo "Development:"
	@echo "  make dev       - Generate contract and start dev server"
	@echo "  make generate  - Generate JSON contract only"
	@echo "  make frontend  - Start frontend dev server only"
	@echo ""
	@echo "Building:"
	@echo "  make build     - Build production frontend"
	@echo "  make build-cli - Build Go CLI binary"
	@echo "  make preview   - Preview production build"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install   - Install all dependencies"
	@echo ""
	@echo "Quality:"
	@echo "  make test      - Run Go tests"
	@echo "  make lint      - Lint Go code"
	@echo "  make fmt       - Format Go code"
	@echo "  make validate  - Validate generated contract"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean     - Clean generated files"
