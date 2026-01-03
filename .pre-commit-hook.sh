#!/bin/bash
# Pre-commit hook for Go code quality checks
# To install: cp .pre-commit-hook.sh .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

set -e

echo "🔍 Running pre-commit checks..."

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "⚠️  golangci-lint not found. Install with:"
    echo "   brew install golangci-lint"
    echo "   OR"
    echo "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "✅ No Go files staged, skipping checks"
    exit 0
fi

echo "📝 Checking $(echo "$STAGED_GO_FILES" | wc -l | tr -d ' ') staged Go files..."

# Run gofmt
echo "🔧 Running gofmt..."
UNFORMATTED=$(gofmt -l $STAGED_GO_FILES)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ The following files are not formatted:"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: make format"
    exit 1
fi

# Run goimports (if available)
if command -v goimports &> /dev/null; then
    echo "🔧 Running goimports..."
    UNFORMATTED_IMPORTS=$(goimports -l $STAGED_GO_FILES)
    if [ -n "$UNFORMATTED_IMPORTS" ]; then
        echo "❌ The following files have incorrect imports:"
        echo "$UNFORMATTED_IMPORTS"
        echo ""
        echo "Run: make format"
        exit 1
    fi
fi

# Run go vet
echo "🔍 Running go vet..."
go vet ./... || {
    echo "❌ go vet failed"
    exit 1
}

# Run golangci-lint on staged files
echo "🔍 Running golangci-lint..."
golangci-lint run --new-from-rev=HEAD~ --config=.golangci.yml || {
    echo "❌ golangci-lint found issues"
    echo ""
    echo "Run: make lint"
    exit 1
}

# Run tests
echo "🧪 Running tests..."
go test -short -race ./... || {
    echo "❌ Tests failed"
    exit 1
}

echo "✅ All pre-commit checks passed!"
exit 0
