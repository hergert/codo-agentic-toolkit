# Codo Agentic Toolkit - Development Commands

# Default: show help
default:
    @just --list

# === CLI Commands (delegate to cli/justfile) ===

# Build codo binary with embedded pack
build:
    cd cli && just build

# Install codo to ~/.local/bin
install:
    cd cli && just install

# Build for development (no install)
dev:
    cd cli && just dev

# === Testing & Validation ===

# Run tests
test:
    cd cli && just test

# Run tests with race detector
test-race:
    cd cli && just test-race

# Run go vet
vet:
    cd cli && just vet

# Run all CI checks (matches GitHub Actions)
ci: vet test-race
    @echo "✅ All CI checks passed!"

# Quick validation before push
validate:
    cd cli && just validate

# === Build & Release ===

# GoReleaser snapshot build (for testing releases)
snapshot:
    cd cli && just snapshot

# Clean build artifacts
clean:
    @rm -rf dist
    @rm -rf cli/.embedded_pack

# === Development Tools ===

# Run build and watch for changes (requires watchexec)
watch:
    cd cli && watchexec -e go -r -- just build

# Run tests and watch for changes
watch-test:
    cd cli && watchexec -e go -r -- just test

# === Code Quality ===

# Format code
fmt:
    cd cli && just fmt

# Run linter
lint:
    cd cli && just lint

# Full check: format, lint, and test
check:
    cd cli && just check

# === Dependencies ===

# Install dependencies
deps:
    cd cli && go mod download
    cd cli && go mod tidy

# Update dependencies
update-deps:
    cd cli && go get -u ./...
    cd cli && go mod tidy