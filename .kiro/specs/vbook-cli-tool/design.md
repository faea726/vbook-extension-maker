# Design Document

## Overview

The Vbook CLI tool is a cross-platform command-line application written in Go that replicates the core functionality of the existing VSCode extension. The tool is completely self-contained within the `cli/` directory and provides three main commands: test, build, and install, allowing developers to manage existing Vbook extensions without requiring VSCode. Project creation is handled by the VSCode extension. The CLI tool operates independently with its own Go module, dependencies, and build process. The design follows Go best practices and the KISS principle, using standard library packages where possible and maintaining a clean separation of concerns.

## Architecture

The application follows a layered architecture with clear separation between CLI interface, business logic, and external dependencies:

```
┌─────────────────────────────────────┐
│           CLI Layer (cmd/)          │
│     Command parsing & validation    │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│        Business Logic Layer        │
│         (internal/)                │
│  ┌─────────────┬─────────────────┐  │
│  │   Project   │     Vbook       │  │
│  │ Operations  │ Communication   │  │
│  └─────────────┼─────────────────┘  │
│  ┌─────────────▼─────────────────┐  │
│  │   Builder   │    Utilities    │  │
│  │ Operations  │                 │  │
│  └─────────────┴─────────────────┘  │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│      External Dependencies         │
│  File System │ Network │ Archives  │
└─────────────────────────────────────┘
```

## Components and Interfaces

### CLI Layer (cmd/)

The CLI layer uses the Cobra framework for command structure and argument parsing:

```go
// Root command structure
type RootCmd struct {
    Version string
    Verbose bool
}

// Command interfaces
type Command interface {
    Execute(args []string) error
    Validate(args []string) error
}
```

**Commands:**

- `test <script-path> --app-url <vbook-app-url> [--params <param1,param2>]` - Script testing command
- `build [project-path]` - Extension building command (defaults to current directory)
- `install <project-path> --app-url <vbook-app-url>` - Extension installation command
- Global flags: `--help/-h`, `--version/-v`, `--verbose`

### Business Logic Layer (internal/)

#### Project Operations (internal/project/)

```go
type ProjectValidator interface {
    ValidateProject(projectPath string) error
    ValidateForTesting(projectPath string) error    // Checks plugin.json exists
    ValidateForBuilding(projectPath string) error   // Checks plugin.json, icon.png, src/ exist
    ValidateForInstall(projectPath string) error    // Checks plugin.json, icon.png exist
    ResolveProjectPath(path string) (string, error) // Handles relative/absolute paths, defaults to current dir
    CheckRequiredFiles(projectPath string, files []string) error
}
```

#### Vbook Communication (internal/vbook/)

```go
type VbookTester interface {
    TestScript(scriptPath, appURL string, params []string) (*TestResult, error)
    StartLocalServer(port int, projectPath string) (*http.Server, error)
    ValidateAppURL(url string) error
}

type VbookInstaller interface {
    InstallExtension(projectPath, appURL string) error
    PreparePluginData(projectPath string) (*PluginData, error)
    ValidateAppURL(url string) error
}

type VbookClient interface {
    SendTestRequest(appURL string, data *TestRequest) (*TestResponse, error)
    SendInstallRequest(appURL string, data *PluginData) error
}
```

#### Extension Builder (internal/builder/)

```go
type ExtensionBuilder interface {
    BuildExtension(projectPath string) (string, error)
    ValidateBuildFiles(projectPath string) error
    CreateArchive(projectPath, outputPath string) error
}

type Archiver interface {
    CreateZip(sourceDir, outputPath string, includePatterns []string) error
    AddFileToZip(writer *zip.Writer, filePath, zipPath string) error
}
```

#### Utilities (internal/utils/)

```go
type NetworkUtils interface {
    GetLocalIP(interfacePrefix string) (string, error)
    ValidateURL(url string) error
    ParseVbookURL(url string) (*VbookURL, error)
}

type FileUtils interface {
    CopyDirectory(src, dst string) error
    ValidateFile(path string) error
    ReadJSONFile(path string, target interface{}) error
}

type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
}
```

## Data Models

### Core Data Structures

```go
// Plugin configuration from plugin.json
type PluginConfig struct {
    Metadata struct {
        Source      string `json:"source"`
        Name        string `json:"name"`
        Author      string `json:"author"`
        Version     string `json:"version"`
        Description string `json:"description"`
    } `json:"metadata"`
    Script struct {
        Language string `json:"language"`
        Main     string `json:"main"`
    } `json:"script"`
}

// Test request data for Vbook communication
type TestRequest struct {
    IP       string      `json:"ip"`
    Root     string      `json:"root"`
    Language string      `json:"language"`
    Script   string      `json:"script"`
    Input    [][]string  `json:"input"`
}

// Test response from Vbook app
type TestResponse struct {
    Status  string                 `json:"status"`
    Output  string                 `json:"output"`
    Error   string                 `json:"error"`
    Details map[string]interface{} `json:"details"`
}

// Plugin data for installation
type PluginData struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Author      string            `json:"author"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Icon        string            `json:"icon"` // base64 encoded
    Language    string            `json:"language"`
    Main        string            `json:"main"`
    Enabled     bool              `json:"enabled"`
    Debug       bool              `json:"debug"`
    Data        string            `json:"data"` // JSON string of script files
}

// Vbook URL components
type VbookURL struct {
    Host string
    Port int
    Base string
}
```

### Error Types

```go
type ValidationError struct {
    Field   string
    Message string
    Missing []string // Specific missing files for detailed error messages
}

type NetworkError struct {
    Operation string
    URL       string
    Cause     error
}

type FileSystemError struct {
    Operation string
    Path      string
    Cause     error
}

// Specific error types for clear user messaging
type ProjectValidationError struct {
    ProjectPath string
    MissingFiles []string
    Message string
}

type URLValidationError struct {
    URL string
    Reason string
}
```

## Project Structure

The entire CLI tool is self-contained within the `cli/` directory and operates independently of the VSCode extension:

```
cli/
├── cmd/
│   ├── root.go          # Root command and global flags
│   ├── test.go          # Test script command
│   ├── build.go         # Build extension command
│   └── install.go       # Install extension command
├── internal/
│   ├── project/
│   │   └── validator.go # Project validation logic
│   ├── vbook/
│   │   ├── tester.go    # Script testing implementation
│   │   ├── installer.go # Extension installation logic
│   │   └── client.go    # HTTP/TCP client utilities
│   ├── builder/
│   │   └── archiver.go  # ZIP creation and compression
│   └── utils/
│       ├── network.go   # Network utility functions
│       ├── files.go     # File operation utilities
│       └── logger.go    # Structured logging
├── go.mod               # Independent Go module
├── go.sum
└── main.go             # Application entry point
```

**Design Principles for CLI Isolation:**

- The CLI tool is a completely separate Go module with its own dependencies
- No shared code or dependencies with the VSCode extension
- All functionality is implemented within the cli/ directory structure
- The CLI can be built, distributed, and run independently
- Project operations work on any valid Vbook extension project directory

## Network Communication

### Script Testing Protocol

The testing functionality replicates the original VSCode extension's network protocol:

1. **Local HTTP Server**: Starts on `appPort - 10` to serve project files
2. **TCP Connection**: Establishes connection to Vbook app on specified port
3. **HTTP-over-TCP**: Sends HTTP-formatted request with test data
4. **Response Handling**: Parses HTTP response and displays results

```go
// HTTP request format sent over TCP
GET /test HTTP/1.1
Host: {vbook-host}
Connection: close
data: {base64-encoded-test-data}

```

### Installation Protocol

Extension installation uses standard HTTP POST requests:

```go
// HTTP request to Vbook app
GET /install HTTP/1.1
Host: {vbook-host}
data: {base64-encoded-plugin-data}
```

## File Operations

### Archive Creation

ZIP archives are created with maximum compression using the standard `archive/zip` package:

```go
func CreateZip(sourceDir, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    zipWriter := zip.NewWriter(file)
    defer zipWriter.Close()

    // Add files with compression
    return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        // Add file to zip with proper structure
    })
}
```

## Error Handling

The application uses Go's idiomatic error handling with wrapped errors for context and provides specific, actionable error messages as required:

```go
func (v *ProjectValidator) ValidateForBuilding(projectPath string) error {
    var missing []string

    if !fileExists(filepath.Join(projectPath, "plugin.json")) {
        missing = append(missing, "plugin.json")
    }
    if !fileExists(filepath.Join(projectPath, "icon.png")) {
        missing = append(missing, "icon.png")
    }
    if !dirExists(filepath.Join(projectPath, "src")) {
        missing = append(missing, "src directory")
    }

    if len(missing) > 0 {
        return &ProjectValidationError{
            ProjectPath: projectPath,
            MissingFiles: missing,
            Message: fmt.Sprintf("Missing required files for building: %s", strings.Join(missing, ", ")),
        }
    }

    return nil
}
```

**Error Message Design Principles:**

- **Specific**: Clearly identify what files are missing or what went wrong
- **Actionable**: Provide guidance on how to fix the issue
- **Context-aware**: Include relevant paths and operation details
- **User-friendly**: Avoid technical jargon in user-facing messages

```go
func (e ProjectValidationError) Error() string {
    return fmt.Sprintf("validation failed for project at %s: %s", e.ProjectPath, e.Message)
}

func (e URLValidationError) Error() string {
    return fmt.Sprintf("invalid Vbook app URL '%s': %s", e.URL, e.Reason)
}
```

## Testing Strategy

### Unit Testing

- **Business Logic**: Test all interfaces with mock implementations
- **Validation**: Table-driven tests for input validation
- **File Operations**: Use temporary directories for file system tests
- **Network Operations**: Mock HTTP clients and servers

```go
func TestProjectGenerator_CreateProject(t *testing.T) {
    tests := []struct {
        name        string
        projectName string
        targetDir   string
        wantErr     bool
    }{
        {"valid project", "test-project", "/tmp", false},
        {"empty name", "", "/tmp", true},
        {"invalid characters", "test/project", "/tmp", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Testing

- **End-to-End**: Test complete workflows with real file system
- **Network Integration**: Test against mock Vbook app server
- **Cross-Platform**: Automated testing on Windows, macOS, and Linux

### Performance Testing

- **Large Projects**: Test with projects containing many files
- **Network Latency**: Test with simulated network delays
- **Memory Usage**: Profile memory usage during archive creation

## Cross-Platform Considerations

### File System

- Use `filepath` package for all path operations
- Handle different path separators automatically
- Support both absolute and relative paths

### Network Interfaces

- Detect local IP addresses across different platforms
- Handle platform-specific network interface naming
- Support both IPv4 and IPv6 where applicable

### Build and Distribution

- Cross-compile for Windows, macOS, and Linux
- Provide platform-specific installation instructions
- Use GitHub Actions for automated builds and releases

## Configuration Management

### Command-Line Interface Design

**Global Flags:**

```bash
--verbose, -v    Enable verbose logging
--help, -h       Show comprehensive help information with examples
--version, -v    Show version information
```

**Command Structure:**

```bash
# Test command - requires script path and app URL
vbook-cli test <script-path> --app-url <vbook-app-url> [--params <param1,param2>]

# Build command - project path defaults to current directory
vbook-cli build [project-path]

# Install command - requires project path and app URL
vbook-cli install <project-path> --app-url <vbook-app-url>

# Help system
vbook-cli --help                    # Show general help
vbook-cli <command> --help          # Show command-specific help with examples
vbook-cli --version                 # Show version information
```

**Help System Implementation:**

- Each command includes usage examples in help output
- Invalid commands suggest using `--help` for guidance
- Command-specific help shows required vs optional parameters
- Help text includes common error scenarios and solutions

### Environment Variables

```bash
VBOOK_APP_URL     # Default Vbook app URL
VBOOK_LOG_LEVEL   # Logging level (debug, info, error)
VBOOK_CONFIG_DIR  # Configuration directory
```

### Configuration File

Optional YAML configuration file for persistent settings:

```yaml
app_url: "http://192.168.1.100:3000"
log_level: "info"
default_params: ["param1", "param2"]
```

**Path Resolution Logic:**

- When no project path is provided, commands default to current working directory
- Relative paths are resolved relative to current working directory
- Absolute paths are used as-is
- All paths are validated to ensure they point to valid Vbook extension projects (contain plugin.json)

This design provides a robust, maintainable, and cross-platform CLI tool that faithfully replicates the VSCode extension functionality while following Go best practices and the KISS principle.
