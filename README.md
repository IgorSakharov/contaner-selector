# Container Selector

A fast and intuitive CLI tool for selecting and connecting to Docker containers. Built with Go, it provides a fuzzy-finder interface for container selection and seamless command execution.

## Features

- 🔍 **Interactive Fuzzy Finder**: Quickly search and select containers by name
- 🚀 **Auto-Selection**: Skip the interactive menu with pattern matching
- 📦 **Direct Command Execution**: Run commands directly without entering the container
- 🎯 **Smart TTY Detection**: Automatically handles interactive vs non-interactive commands
- ⚡ **Fast**: Native Go binary with minimal overhead
- 🔧 **Flexible**: Multiple ways to specify containers and commands

## Installation

### Prerequisites

- Go 1.21 or higher
- Docker installed and running
- Access to Docker daemon (user must be in docker group or have sudo privileges)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/container-selector.git
cd container-selector

# Build the binary
go build -o container-selector

# Optional: Install to PATH
sudo mv container-selector /usr/local/bin/
```

### Install Dependencies

The project uses Go modules for dependency management. Dependencies will be automatically downloaded during build:

```bash
go mod download
```

## Usage

### Basic Usage

```bash
# Interactive mode - select container and enter command
./container-selector

# Run bash in selected container (no command prompt)
./container-selector --no-prompt
```

### Auto-Selection with Filter

Skip the interactive menu by providing a filter pattern:

```bash
# Auto-select container by partial name
./container-selector -f redis

# Auto-select container by image name
./container-selector -f nginx:latest

# Filter is case-insensitive and matches both name and image
./container-selector -f API
```

### Direct Command Execution

Execute commands directly without entering the container:

```bash
# Run a simple command
./container-selector -f myapp -c "ls -la"

# Run command with pipes and redirects
./container-selector -f myapp -c "ps aux | grep node"

# Check container environment
./container-selector -f myapp -c "env | sort"
```

### Command Line Options

```
Flags:
  -c, --command string   Command to run in the container
  -f, --filter string    Auto-select container matching this pattern (skips fuzzy finder)
  -h, --help            Help for container-selector
      --no-prompt       Skip command prompt and use default bash
```

## Examples

### Interactive Workflow

1. **Select and explore a container**:
   ```bash
   ./container-selector
   # Use fuzzy finder to select container
   # Enter command: bash
   ```

2. **Quick shell access**:
   ```bash
   ./container-selector --no-prompt
   # Select container and immediately get bash shell
   ```

### Automated Workflows

1. **Check logs in specific container**:
   ```bash
   ./container-selector -f web-server -c "tail -n 50 /var/log/app.log"
   ```

2. **Run database queries**:
   ```bash
   ./container-selector -f postgres -c "psql -U user -d database -c 'SELECT COUNT(*) FROM users;'"
   ```

3. **Debug application state**:
   ```bash
   ./container-selector -f myapp -c "ps aux"
   ./container-selector -f myapp -c "netstat -tulpn"
   ./container-selector -f myapp -c "df -h"
   ```

## Filter Behavior

The filter (`-f`) option provides smart container matching:

- **Case-insensitive**: Matches `Redis`, `redis`, or `REDIS`
- **Partial matching**: `red` matches `redis-cache`, `redis-session`, etc.
- **Searches both name and image**: Filter applies to container names and image names
- **Unique match required**: If multiple containers match, shows all matches and asks for a more specific filter

### Examples:

```bash
# Multiple matches - shows error
$ ./container-selector -f app
Multiple containers match filter 'app':
  - app-frontend (node:16-alpine)
  - app-backend (python:3.9)
  - app-database (postgres:13)
Error: please use a more specific filter

# Unique match - proceeds
$ ./container-selector -f app-front -c "npm list"
Auto-selected container: app-frontend
```

## Advanced Usage

### Shell Aliases

Add these to your `.bashrc` or `.zshrc` for quick access:

```bash
# Quick container shell
alias dsh='container-selector --no-prompt'

# Container exec with filter
alias dex='container-selector -f'

# Common containers
alias redis-cli='container-selector -f redis -c "redis-cli"'
alias mysql-cli='container-selector -f mysql -c "mysql -u root -p"'
```

### Integration with Scripts

The tool works well in scripts for automated tasks:

```bash
#!/bin/bash
# Backup script example

# Get database backup
container-selector -f postgres -c "pg_dump -U user dbname" > backup.sql

# Check web server status
if container-selector -f nginx -c "nginx -t" > /dev/null 2>&1; then
    echo "Nginx configuration is valid"
else
    echo "Nginx configuration error!"
fi
```

## Architecture

### Components

1. **Docker SDK Integration**: Uses official Docker Go SDK for container listing
2. **Fuzzy Finder**: Leverages `go-fuzzyfinder` for interactive selection
3. **Command Execution**: Shells out to `docker exec` for proper TTY handling
4. **CLI Framework**: Built with Cobra for robust command-line parsing

### Why Go?

- Single static binary - no runtime dependencies
- Fast startup time
- Cross-platform compatibility
- Excellent Docker SDK support
- Strong concurrency for future features

## Development

### Project Structure

```
container-selector/
├── main.go           # Entry point
├── cmd/
│   └── root.go      # Command implementation
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
└── README.md        # This file
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o container-selector-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o container-selector-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o container-selector.exe
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Common Issues

1. **"Cannot connect to Docker daemon"**
   - Ensure Docker is running
   - Check if your user is in the docker group: `groups $USER`
   - Try with sudo if necessary

2. **"No running containers found"**
   - Verify containers are running: `docker ps`
   - Check Docker daemon status: `systemctl status docker`

3. **"TTY error" when using -c flag**
   - This is handled automatically now, but ensure you're using the latest version
   - For interactive commands, connect without -c flag

### Debug Mode

Set environment variable for verbose output:
```bash
DEBUG=1 ./container-selector
```

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Docker SDK for Go](https://github.com/docker/docker)
- [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder)
- [Cobra CLI framework](https://github.com/spf13/cobra)

## Roadmap

- [ ] Add container status information in selection menu
- [ ] Support for docker-compose project filtering
- [ ] Configuration file support for common filters
- [ ] Shell completion scripts
- [ ] Container log viewing mode
- [ ] Multi-container command execution