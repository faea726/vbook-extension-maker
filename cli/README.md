# Vbook CLI Tool

A cross-platform command-line tool for testing, building, and installing Vbook extensions. This CLI tool replicates the core functionality of the Vbook Extension Maker VSCode extension as a standalone application.

## Features

- **Test Extensions**: Test Vbook extension scripts against a running Vbook app
- **Build Extensions**: Create distributable ZIP packages from extension projects
- **Install Extensions**: Install extensions directly to a Vbook app
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Standalone**: No dependencies on VSCode or other tools

## Installation

### Download Pre-built Binaries

Download the appropriate binary for your platform from the releases page:

- **Windows**: `vbook-cli-windows-amd64.exe` or `vbook-cli-windows-arm64.exe`
- **macOS**: `vbook-cli-darwin-amd64` or `vbook-cli-darwin-arm64`
- **Linux**: `vbook-cli-linux-amd64` or `vbook-cli-linux-arm64`

### Build from Source

Requirements:

- Go 1.21 or later

```bash
git clone <repository-url>
cd cli
go build -o vbook-cli .
```

Or use the build scripts:

```bash
# PowerShell (Windows)
.\build.ps1 -Version "1.0.0"

# Make (Unix-like systems)
make build-all
```

## Usage

### Test Extension Scripts

Test a Vbook extension script against a running Vbook app:

```bash
vbook-cli test src/detail.js --app-url http://192.168.1.100:8080
vbook-cli test src/search.js --app-url http://192.168.1.100:8080 --params "query,page"
```

**Options:**

- `--app-url`: URL of the running Vbook app (required)
- `--params`: Test parameters (comma-separated, optional)

### Build Extensions

Build a Vbook extension into a distributable ZIP file:

```bash
# Build current directory
vbook-cli build

# Build specific project
vbook-cli build ./my-extension
vbook-cli build /path/to/my-extension
```

The command creates a `plugin.zip` file in the project root containing:

- `plugin.json` (metadata)
- `icon.png` (extension icon)
- All files from the `src/` directory

### Install Extensions

Install a Vbook extension directly to a running Vbook app:

```bash
# Install from current directory
vbook-cli install --app-url http://192.168.1.100:8080

# Install specific project
vbook-cli install ./my-extension --app-url http://192.168.1.100:8080
vbook-cli install /path/to/my-extension --app-url http://192.168.1.100:8080
```

**Options:**

- `--app-url`: URL of the running Vbook app (required)

### Global Options

- `--verbose, -v`: Enable verbose output with timestamps
- `--help, -h`: Show help information
- `--version`: Show version information

## Project Structure

A valid Vbook extension project must have the following structure:

```
my-extension/
├── plugin.json          # Extension metadata (required)
├── icon.png             # Extension icon (required)
└── src/                 # Source code directory (required)
    ├── main.js          # JavaScript files
    ├── utils.js
    └── ...
```

### plugin.json Format

```json
{
  "metadata": {
    "name": "My Extension",
    "author": "Author Name",
    "version": 1,
    "source": "my-extension",
    "description": "Extension description"
  },
  "script": {
    "detail": "detail.js",
    "toc": "toc.js",
    "chap": "chap.js"
  }
}
```

## Examples

### Complete Workflow

1. **Create or navigate to your extension project:**

   ```bash
   cd my-vbook-extension
   ```

2. **Test your extension:**

   ```bash
   vbook-cli test src/detail.js --app-url http://192.168.1.100:8080 --verbose
   ```

3. **Build the extension:**

   ```bash
   vbook-cli build --verbose
   ```

4. **Install for testing:**
   ```bash
   vbook-cli install --app-url http://192.168.1.100:8080 --verbose
   ```

### URL Formats

The CLI tool accepts various URL formats for the `--app-url` parameter:

- `http://192.168.1.100:8080` (complete URL)
- `https://192.168.1.100:8080` (HTTPS)
- `http://192.168.1.100` (defaults to port 8080)
- `192.168.1.100` (defaults to http:// and port 8080)

## Error Handling

The CLI tool provides clear, actionable error messages:

- **Missing files**: Specific information about which required files are missing
- **Network errors**: Clear messages about connection issues with suggestions
- **Validation errors**: Detailed information about project structure problems
- **Verbose mode**: Additional technical details when using `--verbose`

## Development

### Running Tests

```bash
go test -v ./...
```

### Building for Multiple Platforms

```bash
# Using PowerShell script
.\build.ps1 -Version "1.0.0"

# Using Makefile
make build-all

# Manual cross-compilation
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vbook-cli-linux-amd64 .
```

### Project Structure

```
cli/
├── cmd/                 # CLI command implementations
│   ├── root.go         # Root command and global flags
│   ├── test.go         # Test command
│   ├── build.go        # Build command
│   └── install.go      # Install command
├── internal/           # Internal packages
│   ├── builder/        # Extension building logic
│   ├── project/        # Project validation
│   ├── utils/          # Utility functions
│   └── vbook/          # Vbook communication
├── build.ps1           # Windows build script
├── Makefile           # Unix build script
├── go.mod             # Go module definition
└── main.go            # Application entry point
```

## Troubleshooting

### Common Issues

1. **"plugin.json not found"**

   - Ensure you're in a valid Vbook extension directory
   - Check that `plugin.json` exists in the project root

2. **"Failed to connect to Vbook app"**

   - Verify the Vbook app is running
   - Check the URL format and network connectivity
   - Ensure the port is correct (default: 8080)

3. **"Missing required files"**

   - Ensure `plugin.json`, `icon.png`, and `src/` directory exist
   - Check that the `src/` directory contains JavaScript files

4. **Network interface errors**
   - The CLI tool automatically detects network interfaces
   - If issues persist, check your network configuration

### Getting Help

```bash
# General help
vbook-cli --help

# Command-specific help
vbook-cli test --help
vbook-cli build --help
vbook-cli install --help

# Version information
vbook-cli --version
```

## License

This project is part of the Vbook Extension Maker ecosystem.
