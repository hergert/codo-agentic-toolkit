# Match CI validation steps from .github/workflows/ci.yml

# Default: run all CI checks
default: ci

# Build the project
build:
    cd cli && go build ./...

# Run go vet
vet:
    cd cli && go vet ./...

# Run tests
test:
    cd cli && go test ./...

# Run tests with race detector
race:
    cd cli && go test -race ./...

# Run all CI checks (matches GitHub Actions)
ci: build vet race
    @echo "✅ All CI checks passed!"

# Quick validation before push (alias for ci)
validate: ci

# GoReleaser snapshot build (requires goreleaser installed)
snapshot:
    #!/usr/bin/env sh
    if ! command -v goreleaser >/dev/null 2>&1; then
        echo "❌ goreleaser not installed. Install with: go install github.com/goreleaser/goreleaser/v2@latest"
        exit 1
    fi
    cd cli && goreleaser release --clean --skip=publish --snapshot

# Clean build artifacts
clean:
    rm -rf cli/dist/
    cd cli && go clean ./...
    cd cli && go clean -testcache

# Run build and watch for changes (requires watchexec)
watch:
    cd cli && watchexec -e go -r -- go build ./...

# Run tests and watch for changes
watch-test:
    cd cli && watchexec -e go -r -- go test ./...

# Install dependencies
deps:
    cd cli && go mod download
    cd cli && go mod tidy

# Update dependencies
update-deps:
    cd cli && go get -u ./...
    cd cli && go mod tidy

# Format code
fmt:
    cd cli && go fmt ./...
    
# Run linter (requires golangci-lint)
lint:
    #!/usr/bin/env sh
    if ! command -v golangci-lint >/dev/null 2>&1; then
        echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        echo "   Running basic go vet instead..."
        cd cli && go vet ./...
    else
        cd cli && golangci-lint run
    fi

# Full check: format, lint, and test
check: fmt lint ci

# Show all available commands
help:
    @just --list