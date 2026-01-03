# Release Process

This document describes the release process for the Z.ai Go SDK.

## Version Numbers

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (X.0.0): Incompatible API changes
- **MINOR** version (0.X.0): New functionality, backwards-compatible
- **PATCH** version (0.0.X): Backwards-compatible bug fixes

**Pre-release versions:**
- Alpha: `v0.1.0-alpha.1`
- Beta: `v0.1.0-beta.1`
- Release Candidate: `v0.1.0-rc.1`

## Pre-Release Checklist

Before creating a release, ensure all items are complete:

### Code Quality
- [ ] All tests passing on main branch
  - `make test` passes locally
  - GitHub Actions tests pass
  - Coverage meets minimum threshold (80%)
- [ ] No security vulnerabilities detected
  - `govulncheck ./...` passes
  - Dependabot alerts resolved
- [ ] Code review completed for all changes
- [ ] All examples tested manually
  - Test with real API credentials
  - Verify all 15 examples work
- [ ] Linters pass without warnings
  - `make lint` passes
  - `go vet ./...` passes

### Documentation
- [ ] CHANGELOG.md updated
  - Version number and date added
  - All changes categorized (Added, Changed, Fixed, etc.)
  - Breaking changes highlighted
  - Migration guide updated if needed
- [ ] README.md reviewed and updated
  - Installation instructions current
  - Examples work with new version
  - Links are valid
- [ ] GoDoc comments reviewed
  - New APIs documented
  - Examples added for new features
- [ ] Migration guide updated (if breaking changes)
- [ ] Version number updated in:
  - CHANGELOG.md
  - README.md examples (if applicable)
  - Any version constants in code

### Testing
- [ ] Unit tests cover new functionality
- [ ] Integration tests pass (if applicable)
- [ ] Benchmarks reviewed (no regressions >10%)
- [ ] Examples verified to work
- [ ] Tested on multiple platforms:
  - Linux (GitHub Actions)
  - macOS (GitHub Actions)
  - Windows (GitHub Actions)

## Release Steps

### 1. Create Release Branch (Optional for Major/Minor)

For major or minor releases, create a release branch:

```bash
git checkout main
git pull origin main
git checkout -b release/v0.1.0
```

For patch releases, you can release directly from main.

### 2. Update Version Information

Update the CHANGELOG.md:

```bash
# Replace [Unreleased] with [0.1.0] - YYYY-MM-DD
# Add a new [Unreleased] section
vim CHANGELOG.md
git add CHANGELOG.md
git commit -m "chore: prepare release v0.1.0"
```

### 3. Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.1.0 -m "Release v0.1.0"

# Push tag to trigger release workflow
git push origin v0.1.0
```

**Important:** The tag must follow the format `vX.Y.Z` exactly (e.g., `v0.1.0`, `v1.2.3`).

### 4. Monitor Release Workflow

1. Go to GitHub Actions: https://github.com/sofianhadi1983/zai-sdk-go/actions
2. Watch the "Release" workflow
3. Verify all jobs pass:
   - Validate Release
   - Create Release

### 5. Verify Release

Check that the release was created properly:

1. **GitHub Release:**
   - Visit: https://github.com/sofianhadi1983/zai-sdk-go/releases
   - Verify release notes from CHANGELOG
   - Check release assets (source archives, checksums)

2. **pkg.go.dev:**
   - Visit: https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go@v0.1.0
   - Verify documentation appears (may take a few minutes)
   - Check that all packages are listed

3. **Go Module Proxy:**
   ```bash
   # Verify the version is available
   go list -m -versions github.com/sofianhadi1983/zai-sdk-go
   ```

### 6. Test Installation

Test that users can install the new version:

```bash
# In a temporary directory
mkdir test-install && cd test-install
go mod init test
go get github.com/sofianhadi1983/zai-sdk-go@v0.1.0

# Verify version
go list -m github.com/sofianhadi1983/zai-sdk-go
```

## Post-Release Tasks

### 1. Announcements

- [ ] Update project README badges (if needed)
- [ ] Create announcement (if major/minor release):
  - GitHub Discussions
  - Social media (if applicable)
  - Documentation website (if applicable)

### 2. Monitor for Issues

- [ ] Watch for new issues on GitHub
- [ ] Monitor CI/CD for any failures
- [ ] Check pkg.go.dev for documentation issues

### 3. Update Main Branch

If you created a release branch, merge it back or update main:

```bash
git checkout main
git pull origin main
# If there were any release-specific commits, cherry-pick them
```

### 4. Prepare for Next Release

Create a new "Unreleased" section in CHANGELOG.md:

```markdown
## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security
```

## Hotfix Process

For critical bugs that need immediate release:

### 1. Create Hotfix Branch

```bash
git checkout v0.1.0  # checkout the tag
git checkout -b hotfix/v0.1.1
```

### 2. Make the Fix

```bash
# Make minimal changes to fix the bug
git add .
git commit -m "fix: critical bug description"
```

### 3. Update CHANGELOG

```bash
# Add [0.1.1] section with the fix
vim CHANGELOG.md
git add CHANGELOG.md
git commit -m "chore: prepare hotfix v0.1.1"
```

### 4. Tag and Release

```bash
git tag -a v0.1.1 -m "Hotfix v0.1.1"
git push origin v0.1.1
```

### 5. Backport to Main

```bash
git checkout main
git cherry-pick <hotfix-commit-sha>
git push origin main
```

## Rollback Process

If a release has critical issues:

### 1. Mark Release as Pre-release

1. Go to GitHub Releases
2. Edit the release
3. Check "This is a pre-release"
4. Add warning to release notes

### 2. Create Hotfix

Follow the hotfix process above to create a fixed version.

### 3. Communicate

- Create GitHub issue explaining the problem
- Update release notes with warning
- Announce on relevant channels

## Troubleshooting

### Release Workflow Failed

1. Check the workflow logs in GitHub Actions
2. Common issues:
   - Tests failing: Fix tests and re-tag
   - Invalid tag format: Delete tag and recreate with correct format
   - Permission issues: Check repository settings

### pkg.go.dev Not Updating

1. Wait 10-15 minutes (it's cached)
2. Force refresh: https://pkg.go.dev/github.com/sofianhadi1983/zai-sdk-go@v0.1.0?tab=overview
3. If still not working:
   ```bash
   # Verify the module is accessible
   GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/sofianhadi1983/zai-sdk-go@v0.1.0
   ```

### Tag Already Exists

If you need to recreate a tag:

```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag
git push origin :refs/tags/v0.1.0

# Recreate tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

**Warning:** Only do this if the release hasn't been widely distributed yet.

## Release Cadence

- **Patch releases:** As needed for bug fixes
- **Minor releases:** Monthly or when significant features are ready
- **Major releases:** When breaking changes are necessary

## Security Releases

For security vulnerabilities:

1. **Do not** create a public issue
2. Follow the hotfix process
3. Coordinate with security@anthropic.com if needed
4. Add security advisory after release
5. Update CHANGELOG with Security section

## Checklist Summary

```markdown
## Pre-Release
- [ ] Tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Examples tested
- [ ] Security scans pass

## Release
- [ ] Create tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Verify GitHub release created
- [ ] Verify pkg.go.dev updated
- [ ] Test installation

## Post-Release
- [ ] Monitor for issues
- [ ] Create announcements
- [ ] Update main branch
```

## Additional Resources

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Publishing Go Modules](https://go.dev/blog/publishing-go-modules)
