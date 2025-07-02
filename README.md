# Docker PSA - Human Readable Container Listings

A Docker CLI plugin that provides a more human-readable format for container listings, enhancing the standard `docker ps` output.

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

### Option 1: Download Pre-built Binary (Recommended)

1. Download the appropriate binary for your platform from [GitHub Releases](https://github.com/agkunz/docker-psa/releases):
   - `docker-psa-linux-amd64` - Linux x86_64
   - `docker-psa-darwin-amd64` - macOS Intel
   - `docker-psa-darwin-arm64` - macOS Apple Silicon
   - `docker-psa-windows-amd64.exe` - Windows x86_64

2. Create the Docker CLI plugins directory:
```bash
mkdir -p ~/.docker/cli-plugins
```

3. Move the binary to the plugins directory and make it executable:
```bash
# Linux/macOS
mv docker-psa-* ~/.docker/cli-plugins/docker-psa
chmod +x ~/.docker/cli-plugins/docker-psa

# Windows (PowerShell)
Move-Item docker-psa-windows-amd64.exe $env:USERPROFILE\.docker\cli-plugins\docker-psa.exe
```

4. Verify the plugin is installed:
```bash
docker psa --help
```

### Option 2: Build from Source

1. Clone and build:
```bash
git clone git@github.com:agkunz/docker-psa.git
cd docker-psa
make install
```

2. Verify installation:
```bash
docker psa --help
```

For detailed development setup instructions, see [CONTRIBUTING.md](CONTRIBUTING.md).

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

Interested in contributing? Check out [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, build instructions, and contribution guidelines.

## License

MIT
