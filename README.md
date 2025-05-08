# Docker PSA - Human Readable Container Listings

A Docker CLI plugin that provides a more human-readable format for container listings, enhancing the output of `docker ps`.

## Features

- 🟢 Colorful status indicators for running, exited, and created containers
- 💚 Health status indicators (green heart for healthy, orange heart for unhealthy)
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
- Build the plugin
- Create the Docker CLI plugins directory if it doesn't exist
- Copy the binary to the Docker CLI plugins directory (`~/.docker/cli-plugins/`)
- Set appropriate permissions

3. Verify the plugin is installed:
```bash
docker psa --help
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

## License

MIT