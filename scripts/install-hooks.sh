#!/bin/bash
# Install Git hooks for the project
# Usage: ./scripts/install-hooks.sh

set -e

HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Installing Git hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Copy pre-commit hook
cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "✓ Installed pre-commit hook"
echo ""
echo "The following checks will run before each commit:"
echo "  - golangci-lint (linting)"
echo "  - go test (unit tests)"
echo ""
echo "To bypass (not recommended): git commit --no-verify"
