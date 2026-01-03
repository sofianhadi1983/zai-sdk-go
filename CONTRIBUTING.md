# Contributing to Z.ai SDK for Go

Thank you for your interest in contributing to the Z.ai SDK for Go!

## Development Process

### 1. Fork and Clone

Fork the repository and clone it locally:

```bash
git clone https://github.com/YOUR_USERNAME/zai-sdk-go.git
cd zai-sdk-go
```

### 2. Set Up Development Environment

```bash
# Install dependencies
make deps

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### 3. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 4. Make Your Changes

- Write clean, idiomatic Go code
- Follow the project's architecture patterns (see CLAUDE.md)
- Add tests for new functionality
- Update documentation as needed

### 5. Test Your Changes

```bash
# Run tests
make test

# Run linters
make lint

# Check code coverage
make test-cover
```

### 6. Commit Your Changes

Follow conventional commit messages:

```
feat: add new feature
fix: fix bug
docs: update documentation
test: add tests
refactor: refactor code
```

### 7. Submit a Pull Request

- Push your changes to your fork
- Create a pull request against the `main` branch
- Describe your changes clearly
- Link any related issues

## Code Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` and `goimports` for formatting
- Maximum line length: 120 characters
- Use meaningful variable and function names

### Testing

- Write table-driven tests
- Use `t.Parallel()` where appropriate
- Aim for >80% code coverage
- Include both unit and integration tests

### Documentation

- Document all exported types, functions, and packages
- Use GoDoc format
- Include code examples where helpful

## Architecture Guidelines

See [CLAUDE.md](../z-ai-sdk-python/CLAUDE.md) for detailed architecture guidelines including:

- Clean Architecture principles
- Interface-driven development
- Error handling patterns
- Testing strategies
- Observability with OpenTelemetry

## Pull Request Checklist

- [ ] Tests pass (`make test`)
- [ ] Linters pass (`make lint`)
- [ ] Code is formatted (`make format`)
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (if applicable)
- [ ] Commit messages follow conventions

## Getting Help

- Check existing issues and pull requests
- Ask questions in your PR
- Contact: user_feedback@z.ai

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
