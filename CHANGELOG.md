# Changelog

All notable changes to the Z.ai Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-01-03

### Added
- Initial implementation of Z.ai Go SDK
- Support for all 15 Z.ai API services:
  - Chat Completions (with streaming, function calling, multimodal support)
  - Embeddings (with batch processing)
  - Image Generation
  - Files (upload, download, management)
  - Audio (transcription)
  - Videos (async generation)
  - Assistant (conversational AI)
  - Batch (batch processing)
  - Web Search (AI-powered search)
  - Moderations (content safety)
  - Tools (function calling)
  - Agents (agent invocation)
  - Voice (voice cloning)
  - OCR (handwriting recognition)
  - File Parser (document parsing)
  - Web Reader (web content extraction)
- Comprehensive GoDoc documentation for all exported types and functions
- 15 example applications demonstrating API usage
- Migration guide from Python SDK
- Builder pattern for fluent API usage
- Automatic JWT token generation and caching
- Retry logic with exponential backoff
- Context support for cancellation and timeout
- Type-safe request/response handling
- Streaming support via channels
- Custom HTTP client support
- Configurable logging
- Support for both international and Chinese (Zhipu) endpoints

### Documentation
- Comprehensive README with usage examples for all APIs
- Examples README with instructions for running all examples
- Migration guide for Python SDK users
- GoDoc comments with examples for all exported functions
- Best practices and common patterns documentation
- RELEASING.md with complete release process guide
- CHANGELOG.md with semantic versioning strategy

### CI/CD
- GitHub Actions workflows for testing, security, and releases
- Multi-platform testing (Ubuntu, macOS, Windows) across Go 1.21-1.23
- Automated security scanning (govulncheck, gosec, OSV Scanner, Trivy)
- Performance regression tracking with benchmarks
- Dependabot for automated dependency updates
- Codecov integration for coverage reporting

## Versioning Strategy

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR version** (X.0.0): Incompatible API changes
- **MINOR version** (0.X.0): New functionality in a backwards-compatible manner
- **PATCH version** (0.0.X): Backwards-compatible bug fixes

### Version Compatibility

The Go SDK version numbers are independent from the Z.ai API version and the Python SDK version. However, we aim to maintain feature parity with the official Python SDK where possible.

### Deprecation Policy

- Deprecated features will be marked with `// Deprecated:` comments in code
- Deprecated features will remain functional for at least one minor version
- Migration instructions will be provided in deprecation notices
- Deprecated features will be listed in CHANGELOG with alternatives

## Release Notes Template

For future releases, use this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features and capabilities

### Changed
- Changes to existing functionality

### Deprecated
- Features marked for removal in future versions

### Removed
- Features removed in this version

### Fixed
- Bug fixes

### Security
- Security improvements and vulnerability fixes
```

## Upgrade Guide

### From 0.x to 1.0 (Future)

When upgrading to version 1.0 (future release):
- Review all deprecated features in 0.x releases
- Update code to use recommended alternatives
- Test thoroughly with the new version
- Check for breaking changes in the release notes

## Links

- [Go SDK Repository](https://github.com/sofianhadi1983/zai-sdk-go)
- [Python SDK Repository](https://github.com/zhipuai/zhipuai-sdk-python)
- [Official Z.ai API Documentation](https://open.bigmodel.cn/dev/api)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

---

**Note**: This CHANGELOG will be updated with each release. See [releases](https://github.com/sofianhadi1983/zai-sdk-go/releases) for published versions.
