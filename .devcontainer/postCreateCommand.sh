#!/bin/bash

# Post-create command for Cursor Admin API Exporter devcontainer
# Minimal, robust setup script

set -e

echo "🚀 Setting up Cursor Admin API Exporter development environment..."

# Verify Go installation
echo "📍 Verifying Go installation..."
if go version; then
    echo "✅ Go is available"
else
    echo "❌ Go not found, this may cause issues"
    exit 1
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
if go mod download; then
    echo "✅ Dependencies downloaded successfully"
else
    echo "⚠️ Failed to download dependencies"
fi

# Install essential Go tools
echo "🔧 Installing essential Go tools..."

echo "📦 Installing goimports..."
if go install golang.org/x/tools/cmd/goimports@latest; then
    echo "✅ goimports installed"
else
    echo "⚠️ Failed to install goimports"
fi

echo "📦 Installing golangci-lint..."
if curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.3; then
    echo "✅ golangci-lint installed"
else
    echo "⚠️ Failed to install golangci-lint"
fi

# Make scripts executable
echo "🔑 Making scripts executable..."
if [ -d "scripts" ]; then
    chmod +x scripts/*.sh
    echo "✅ Scripts made executable"
else
    echo "⚠️ Scripts directory not found"
fi

# Test basic functionality
echo "🧪 Testing basic setup..."
if go fmt ./...; then
    echo "✅ Go formatting works"
else
    echo "⚠️ Go formatting failed"
fi

echo "🎉 Basic development environment setup complete!"
echo ""
echo "Available commands:"
echo "  make help    - Show all available make targets"
echo "  make build   - Build the application"
echo "  make test    - Run tests"
echo "  make fmt     - Format code"
echo "  make dev     - Development mode (requires air)"
echo ""
echo "🎉 Setup completed successfully!"