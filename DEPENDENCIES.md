# Dependencies

This document lists all dependencies used in the Z.ai SDK for Go.

## Core Dependencies

### Production Dependencies

- **github.com/golang-jwt/jwt/v5** (v5.3.0)
  - Purpose: JWT token generation and validation for API authentication
  - License: MIT
  - Used in: `internal/auth/jwt.go`

### Development Dependencies

- **github.com/stretchr/testify** (v1.11.1)
  - Purpose: Testing utilities and assertions
  - License: MIT
  - Used in: All test files (*_test.go)

- **go.uber.org/mock** (v0.6.0)
  - Purpose: Mock generation for testing
  - License: Apache-2.0
  - Used in: `test/mocks/`

## Tool Dependencies

Tools are tracked in `internal/tools.go` with build tag `//go:build tools`.

- **mockgen** (from go.uber.org/mock)
  - Purpose: Generate mock implementations of interfaces
  - Install: `go install go.uber.org/mock/mockgen@latest`
  - Usage: `mockgen -source=interface.go -destination=mocks/mock_interface.go`

## Future Dependencies

The following dependencies will be added as implementation progresses:

### Observability (Phase 2+)
- **go.opentelemetry.io/otel** - OpenTelemetry SDK for tracing and metrics
- **go.opentelemetry.io/otel/trace** - Distributed tracing
- **go.opentelemetry.io/otel/metric** - Metrics collection

### Utilities (as needed)
- **golang.org/x/sync/errgroup** - Goroutine synchronization (may use standard library)

## Dependency Management

### Adding Dependencies

```bash
# Add a production dependency
go get github.com/example/package

# Add a development dependency
go get -d github.com/example/package

# Update go.mod and go.sum
go mod tidy
```

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update specific dependency
go get -u github.com/example/package
```

### Vendoring (Optional)

If you need to vendor dependencies for reproducible builds:

```bash
go mod vendor
```

## Security

Dependencies are regularly checked for vulnerabilities using:

```bash
# Check for known vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## Version Constraints

- Go Version: 1.25.5 (minimum 1.23)
- All dependencies use semantic versioning
- Use `go list -m all` to see full dependency tree
