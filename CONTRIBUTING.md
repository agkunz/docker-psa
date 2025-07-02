# Contributing to Docker PSA

Thank you for your interest in contributing to Docker PSA! This guide provides guidelines for contributing to the project.

## Development Setup

### Prerequisites
- Go 1.23 or later
- Docker (for testing the plugin)
- Git

### Clone and Setup
```bash
git clone git@github.com:agkunz/docker-psa.git
cd docker-psa
go mod download
make dev-setup  # Install Go tools and set up git hooks
```

### GitHub Setup (Optional)
```bash
# Setup GitHub remote if not already configured
make setup-github
```

### Build and Test
```bash
make build        # Build to build/ directory
make install      # Install the plugin locally
docker psa        # Test the plugin

# Cross-platform builds
make build-all    # Build for all platforms

# Individual platform builds
make build-linux    # Linux AMD64
make build-darwin   # macOS (both Intel and Apple Silicon)
make build-windows  # Windows AMD64

# Clean up
make clean        # Remove build artifacts
```

### Development Environment Notes

The development setup automatically creates a Python virtual environment (`.venv`) for pre-commit hooks. This directory is excluded from git and will be created locally on each developer's machine.

## Development Workflow

1. **Fork and Branch**
   - Fork the repository
   - Create a feature branch: `git checkout -b feat/your-feature-name`

2. **Make Changes**
   - Write code following Go best practices
   - Add tests for new functionality
   - Ensure all tests pass: `go test ./...`
   - Format code: `go fmt ./...` (automatically done by git hooks)

3. **Commit**
   - Use conventional commit messages (enforced by git hooks)
   - Git hooks will automatically:
     - Format your Go code
     - Run `go vet`
     - Tidy go.mod/go.sum
     - Run quick tests
     - Validate commit message format
   - Make atomic commits (one logical change per commit)

4. **Submit Pull Request**
   - Push to your fork: `git push origin feature-branch`
   - Create a pull request on GitHub with a clear description
   - Link any related issues
   ```

## Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/) for our commit messages. This enables automated versioning and changelog generation.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit

### Examples

```bash
# Feature addition
feat: add watch mode for continuous container monitoring

# Bug fix
fix: handle containers with no exposed ports gracefully

# Breaking change
feat!: redesign output format with improved readability

# With scope
feat(formatting): add color-coded status indicators

# With body and footer
fix: resolve memory leak in watch mode

Previously, the watch mode would accumulate container data
without cleaning up old entries, causing memory usage to
grow over time.

Closes #123
```

## Pull Request Process

1. **Fork and Branch**
   - Fork the repository
   - Create a feature branch: `git checkout -b feat/your-feature-name`

2. **Make Changes**
   - Write code following Go best practices
   - Add tests for new functionality
   - Ensure all tests pass: `go test ./...`
   - Format code: `go fmt ./...`

3. **Commit**
   - Use conventional commit messages
   - Make atomic commits (one logical change per commit)

5. **Submit Pull Request**
   - Push to your fork: `git push origin feature-branch`
   - Create a pull request on GitHub with a clear description
   - Link any related issues

## Release Process

Releases are automated using semantic-release:

- **Patch Release** (0.1.0 → 0.1.1): `fix:`, `docs:`, `perf:`, `refactor:`
- **Minor Release** (0.1.0 → 0.2.0): `feat:`
- **Major Release** (0.1.0 → 1.0.0): `feat!:`, `fix!:`, or any commit with `BREAKING CHANGE:` in the footer

### Automated Release Pipeline

When commits are pushed to the main branch, the CI/CD pipeline will:

1. **Automatic Versioning** - Version numbers determined by conventional commit messages
2. **Multi-Platform Builds** - Binaries built for Linux, macOS, and Windows
3. **Release Notes** - Changelog automatically generated from commits
4. **GitHub Releases** - Release assets uploaded automatically

The pipeline uses both GitLab CI/CD and GitHub Actions for:
- **Testing** - Automated tests on multiple Go versions
- **Code Quality** - Format checking, vetting, and linting
- **Multi-Platform Builds** - Cross-compilation for different OS/architectures
- **Semantic Versioning** - Automated version management
- **Dual Releases** - Automatic releases on both GitLab and GitHub

### Platform-Specific Features

- **GitLab CI**: Uses `.releaserc.json` and GitLab's integrated CI/CD
- **GitHub Actions**: Uses `.releaserc.github.json` and GitHub's workflow system
- **Synchronized**: Both platforms get the same version numbers and release notes

### GitHub Workflow Commands

```bash
# Setup GitHub remote (if not already configured)
make setup-github

# Push changes to GitHub
make push
```

## Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and reasonably sized
- Write tests for new functionality

## Testing

- Run tests with `go test ./...`
- Test the plugin manually with `make build && make install && docker psa`
- Ensure the plugin works with various Docker container states

## Code Quality Tools

This project uses **Go-native tools** to maintain code quality:

### Git Hooks (Pure Go/Shell)
Pre-commit hooks automatically run on every commit to:
- **Format Go code** - Runs `go fmt` and `goimports`, stages formatted files
- **Run go vet** - Catches common Go issues
- **Tidy modules** - Ensures go.mod/go.sum are clean
- **Run quick tests** - Tests affected packages
- **Validate commit messages** - Go-based conventional commit validation

### Go Development Tools
```bash
# Install recommended Go tools
make install-tools

# Tools installed:
# - goimports: Import organization
# - staticcheck: Advanced static analysis
# - gotestsum: Better test output
# - golangci-lint: Comprehensive linting
```

### Manual Code Quality Checks
```bash
# Run all format checks (same as CI)
make format-check

# Format code with Go tools
make format

# Run advanced linting
make lint

# Run comprehensive tests
make test

# One-time setup (tools + hooks)
make dev-setup
```

### Bypassing Hooks (Emergency Only)
```bash
# Skip hooks (not recommended)
git commit --no-verify -m "emergency fix"
```

## Questions?

Feel free to open an issue for any questions about contributing!
