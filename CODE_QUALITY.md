# Code Quality Guidelines

This document outlines the code quality standards and tools used in the Z.ai SDK for Go.

## Quick Start

### Install Development Tools

```bash
make install-tools
```

This installs:
- `golangci-lint` - Comprehensive linter
- `goimports` - Import formatting tool
- `mockgen` - Mock generation tool

### Run All Checks

```bash
make check
```

This runs:
1. Code formatting (`gofmt`, `goimports`)
2. Static analysis (`go vet`)
3. Linting (`golangci-lint`)
4. Unit tests

## Linting

### Configuration

Linting is configured in `.golangci.yml` with the following enabled linters:

#### Core Linters
- **errcheck** - Check for unchecked errors
- **gosimple** - Simplify code
- **govet** - Reports suspicious constructs
- **ineffassign** - Detect ineffectual assignments
- **staticcheck** - Advanced static analysis
- **unused** - Check for unused code

#### Additional Linters
- **gofmt** - Check code formatting
- **goimports** - Check import formatting
- **gosec** - Security analysis
- **misspell** - Spell checking
- **revive** - Fast, configurable linter
- **goconst** - Find repeated strings
- **gocyclo** - Cyclomatic complexity
- **dupl** - Code duplication detection
- **gocritic** - Comprehensive diagnostics

### Run Linting

```bash
# Run all linters
make lint

# Run linters and auto-fix issues
make lint-fix
```

### Linting Rules

- **Cyclomatic Complexity:** Maximum 15 per function
- **Code Duplication:** Minimum 3 occurrences to report
- **String Constants:** Minimum length 3, minimum occurrences 3
- **Error Checking:** All errors must be checked (except explicitly ignored)

### Suppressing Linters

Use `//nolint` directives sparingly and always with explanation:

```go
//nolint:gosec // G104: This error is intentionally ignored because...
func example() {
    // code
}
```

## Code Formatting

### Standards

- **Indentation:** Tabs (width 4)
- **Line Length:** 120 characters (soft limit)
- **Imports:** Grouped and sorted by goimports
- **Local Imports:** `github.com/z-ai/zai-sdk-go` packages listed last

### Run Formatting

```bash
make format
```

This runs:
1. `gofmt` - Standard Go formatting
2. `goimports` - Import organization

### EditorConfig

The project includes `.editorconfig` for consistent editor settings:
- Tabs for Go files
- Spaces for YAML, JSON, Markdown
- UTF-8 encoding
- LF line endings
- Trim trailing whitespace

## Static Analysis

### Go Vet

```bash
make vet
```

Checks for:
- Suspicious constructs
- Shadowed variables
- Printf format string errors
- Struct tag errors
- Unreachable code

## Pre-commit Hooks

### Install

```bash
make pre-commit-install
```

Or manually:

```bash
cp .pre-commit-hook.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### What It Checks

The pre-commit hook runs:
1. ✅ Code formatting (`gofmt`, `goimports`)
2. ✅ Static analysis (`go vet`)
3. ✅ Linting (`golangci-lint`)
4. ✅ Unit tests (short mode)

### Skip Pre-commit Hook

If needed (not recommended):

```bash
git commit --no-verify
```

## Testing Requirements

### Coverage

- **Minimum Coverage:** 80% for core packages (`pkg/`, `internal/`)
- **Test Files:** All packages must have `*_test.go` files
- **Table-Driven Tests:** Preferred for unit tests
- **Parallel Tests:** Use `t.Parallel()` where applicable

### Run Tests

```bash
# Unit tests
make test

# With coverage report
make test-cover

# Integration tests
make test-integration
```

## CI/CD Integration

All quality checks run in CI/CD pipelines:

1. **Formatting Check:** Code must be formatted
2. **Linting:** All linters must pass
3. **Tests:** All tests must pass with >80% coverage
4. **Security Scan:** `gosec` must pass

## Best Practices

### Error Handling

```go
// ✅ Good - Always check errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ❌ Bad - Unchecked error
result, _ := doSomething()
```

### Naming Conventions

```go
// ✅ Good - Clear, concise names
func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

// ❌ Bad - Stuttering
func (c *Client) ClientCreateChatCompletion(ctx context.Context, req *ChatCompletionRequest)
```

### Documentation

```go
// ✅ Good - Clear GoDoc comment
// NewClient creates a new Z.ai API client with the provided options.
// It returns an error if the API key is not provided.
func NewClient(opts ...ClientOption) (*Client, error)

// ❌ Bad - No documentation
func NewClient(opts ...ClientOption) (*Client, error)
```

### Context Usage

```go
// ✅ Good - Always accept context
func (c *Client) Create(ctx context.Context, req *Request) (*Response, error)

// ❌ Bad - No context
func (c *Client) Create(req *Request) (*Response, error)
```

### Interface Design

```go
// ✅ Good - Small, focused interfaces
type Reader interface {
    Read(ctx context.Context, id string) (*Data, error)
}

// ❌ Bad - Large interface
type Everything interface {
    Read(...) error
    Write(...) error
    Delete(...) error
    Update(...) error
    // ... many more methods
}
```

## Troubleshooting

### Linter Failures

```bash
# View detailed linter output
golangci-lint run --config=.golangci.yml -v ./...

# Run specific linter
golangci-lint run --disable-all --enable=errcheck ./...
```

### False Positives

If a linter reports a false positive:

1. Verify it's truly a false positive
2. Add `//nolint` with explanation
3. Consider opening an issue with the linter

### Slow Linting

```bash
# Run only fast linters
golangci-lint run --fast ./...

# Enable linter cache
golangci-lint cache clean
```

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
