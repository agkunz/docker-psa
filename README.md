# Docker PSA - Human Readable Container Listings

A Docker CLI plugin that provides a more human-readable format for container listings, enhancing the out### Download Pre-built Binaries

You ## CI/CD Pipeline

The project uses GitHub Actions for:

- **Testing** - Automated tests on multiple Go versions
- **Code Quality** - Format checking, vetting, and linting
- **Multi-Platform Builds** - Cross-compilation for different OS/architectures
- **Semantic Versioning** - Automated version management
- **GitHub Releases** - Automatic releases with build artifacts

### Release Process

- **Automatic Versioning** - Version numbers determined by conventional commit messages
- **Multi-Platform Builds** - Binaries built for Linux, macOS, and Windows
- **Release Notes** - Changelog automatically generated from commits
- **GitHub Releases** - Release assets uploaded automaticallye-built binaries from the GitHub releases page:

- **GitHub**: [Releases](https://github.com/username/docker-psa/releases) *(Update with your GitHub URL)*

Available binaries:
- `docker-psa-linux-amd64` - Linux x86_64
- `docker-psa-darwin-amd64` - macOS Intel
- `docker-psa-darwin-arm64` - macOS Apple Silicon
- `docker-psa-windows-amd64.exe` - Windows x86_64er ps`.

## Features

- 🟢 Colorful status indicators for running, exited, created, and restarting containers
- 💚 Health status indicators (green heart for healthy, orange heart for unhealthy)
- 🔥 Critical indicators for restarting containers (fire emoji with bold red text)
- ⏱️ Human-readable time formatting (e.g., "2 hours ago", "Yesterday")
- 📋 Clean tabular output for better readability
- 🔄 Shows all containers by default (like `docker ps -a`)
- 🔍 Filter containers by name or image using regex patterns
- 📊 Multiple verbosity levels for detailed information

## Installation

### Building from Source

1. Clone this repository
```bash
git clone git@github.com:agkunz/docker-psa.git
cd docker-psa
```

2. Build and install the plugin using the Makefile
```bash
make install
```

This will:
- Build the plugin to the `build/` directory
- Create the Docker CLI plugins directory if it doesn't exist
- Copy the binary to the Docker CLI plugins directory (`~/.docker/cli-plugins/`)
- Set appropriate permissions

3. Verify the plugin is installed:
```bash
docker psa --help
```

### Development Setup

For contributors, set up the development environment:

```bash
# Install development tools and git hooks
make dev-setup

# Run format checks and linting
make format-check

# Format code
make format
```

**Note:** The development setup automatically creates a Python virtual environment (`.venv`) for pre-commit hooks. This directory is excluded from git and will be created locally on each developer's machine.

### Build Options

```bash
# Basic build (outputs to build/ directory)
make build

# Cross-compile for all platforms
make build-all

# Individual platform builds
make build-linux    # Linux AMD64
make build-darwin   # macOS (both Intel and Apple Silicon)
make build-windows  # Windows AMD64

# Clean build artifacts
make clean
```

## Usage

### Basic Command

Simply run:

```bash
docker psa
```

This will show all containers in a human-readable format with status, name, image, and ports (if any).

### Filtering Containers

You can filter containers by name or image using a regex pattern:

```bash
docker psa <regex-pattern>
```

Examples:
```bash
docker psa nginx        # Show containers with "nginx" in name or image
docker psa mongo        # Show containers related to MongoDB
docker psa "web|api"    # Show containers with "web" OR "api" in name/image
docker psa "^front"     # Show containers with names starting with "front"
```

### Verbosity Levels

The plugin supports different verbosity levels for detailed information:

```bash
docker psa -v           # Verbose mode with additional details
docker psa -vv          # Very verbose mode with maximum information
```

#### Verbose Mode (-v)
Shows container ID, age, creation time, image, command, and ports.

#### Very Verbose Mode (-vv)
Shows all the information from verbose mode plus:
- Network details with IP addresses (each network on its own line)
- Volume mounts (each volume mapping on its own line)

### Combining Options

You can combine filtering with verbosity levels:

```bash
docker psa -v mongo     # Show MongoDB containers with verbose details
docker psa mongo -v     # Same as above (flags can be before or after the pattern)
docker psa -vv nginx    # Show nginx containers with maximum details
```

### Health Status Indicators

The plugin uses different indicators for container health status:
- 💚 (Green heart) - Running container with healthy status
- 🧡 (Orange heart) - Running container with unhealthy status
- 🟢 (Green circle) - Running container without health check
- 🔴 (Red circle) - Exited/stopped container
- 🟡 (Yellow circle) - Created but not started container
- 🔥 (Fire) - Restarting container (critical state)

## Examples

### Basic Output
```
💚 Up 2 hours | my-webapp | nginx:latest | Ports: 80:80
🟢 Up 3 days | my-database | postgres:13 | Ports: 5432
🔴 Exited (0) 5 hours ago | my-batch-job | python:3.9
```

### Verbose Output (-v)
```
────────────────────────────────────────────────────────────────────────────────
🐋 3a7c21f85d32 | my-webapp
   💚 Up 2 hours (healthy) | Age: 2 hours | Created: 2 hours ago
   Image: nginx:latest | Command: nginx -g daemon off;
   Ports: 80:80/tcp
────────────────────────────────────────────────────────────────────────────────
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature-name`
3. Make your changes following our coding standards
4. Use conventional commit messages (e.g., `feat: add new feature`)
5. Run tests: `go test ./...`
6. Submit a pull request

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning:

- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `feat!:` or `fix!:` - Breaking changes (major version bump)
- `docs:`, `style:`, `refactor:`, `test:`, `chore:` - Other changes

## Releases

Releases are automated using semantic-release. When commits are pushed to the main branch:

1. **Automatic Versioning** - Version numbers are determined by commit messages
2. **Multi-Platform Builds** - Binaries are built for Linux, macOS, and Windows
3. **Release Notes** - Changelog is automatically generated
4. **GitHub/GitLab Releases** - Release assets are uploaded automatically

### Download Pre-built Binaries

You can download pre-built binaries from the releases page on either platform:

- **GitLab**: [Releases](../../releases)
- **GitHub**: [Releases](https://github.com/username/docker-psa/releases) *(Update with your GitHub URL)*

Available binaries:
- `docker-psa-linux-amd64` - Linux x86_64
- `docker-psa-darwin-amd64` - macOS Intel
- `docker-psa-darwin-arm64` - macOS Apple Silicon
- `docker-psa-windows-amd64.exe` - Windows x86_64

## GitHub Workflow

This project uses GitHub for hosting and CI/CD:

### Quick GitHub Setup

```bash
# Setup GitHub remote (if not already configured)
make setup-github

# Push changes to GitHub
make push
```

## CI/CD Pipeline

The project uses both GitLab CI/CD and GitHub Actions for:

- **Testing** - Automated tests on multiple Go versions
- **Code Quality** - Format checking, vetting, and linting
- **Multi-Platform Builds** - Cross-compilation for different OS/architectures
- **Semantic Versioning** - Automated version management
- **Dual Releases** - Automatic releases on both GitLab and GitHub

### Platform-Specific Features

- **GitLab CI**: Uses `.releaserc.json` and GitLab's integrated CI/CD
- **GitHub Actions**: Uses `.releaserc.github.json` and GitHub's workflow system
- **Synchronized**: Both platforms get the same version numbers and release notes

## License

MIT
