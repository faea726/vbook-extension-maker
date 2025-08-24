# Implementation Plan

- [x] 1. Analyze original VSCode extension logic and context

  - Study the existing VSCode extension code to understand current implementation
  - Document the network protocols, data structures, and communication patterns
  - Identify key functionality that needs to be replicated in CLI tool
  - Map out the test, build, and install workflows from the VSCode extension
  - _Requirements: 1.1, 2.1, 3.1, 7.1_

- [x] 2. Clean up and analyze existing CLI code

  - Remove old/outdated files and code from cli/ directory
  - Analyze existing code structure and identify reusable components
  - Document what needs to be kept, improved, or completely rewritten
  - Clean up directory structure to match new design requirements
  - _Requirements: 7.4, 7.1_

- [x] 3. Set up CLI project structure and core interfaces

  - Initialize or update Go module in cli/ directory with go.mod
  - Create proper directory structure (cmd/, internal/project/, internal/vbook/, internal/builder/, internal/utils/)
  - Define core interfaces for ProjectValidator, VbookTester, VbookInstaller, ExtensionBuilder
  - Create main.go entry point with basic Cobra root command setup
  - _Requirements: 7.4, 4.1_

- [x] 4. Implement project validation and path resolution

  - Create ProjectValidator interface implementation with validation methods
  - Implement ValidateForTesting, ValidateForBuilding, ValidateForInstall methods
  - Add ResolveProjectPath function to handle relative/absolute paths and default to current directory
  - Write unit tests for path resolution and validation logic
  - _Requirements: 5.1, 5.2, 5.4, 5.5_

- [x] 5. Create data models and JSON parsing

  - Implement PluginConfig struct for plugin.json parsing
  - Create TestRequest, TestResponse, and PluginData structs
  - Add JSON marshaling/unmarshaling methods with proper error handling
  - Write unit tests for data model validation and JSON operations
  - _Requirements: 1.2, 2.2, 3.2_

- [x] 6. Implement file utilities and operations

  - Create FileUtils interface with file validation and JSON reading methods
  - Implement file existence checks and directory validation
  - Add utility functions for reading JavaScript files from src/ directory
  - Write unit tests for file operations using temporary directories
  - _Requirements: 2.2, 3.3, 5.5_

- [x] 7. Build ZIP archive creation functionality

  - Implement ExtensionBuilder interface with BuildExtension method
  - Create Archiver interface for ZIP file creation with maximum compression
  - Add logic to include plugin.json, icon.png, and all src/ files in ZIP
  - Write unit tests for archive creation with various file structures
  - _Requirements: 2.1, 2.3, 2.4_

- [x] 8. Implement network utilities and URL validation

  - Create NetworkUtils interface with URL validation methods
  - Implement ValidateURL and ParseVbookURL functions
  - Add GetLocalIP method for detecting local network interface
  - Write unit tests for network utility functions
  - _Requirements: 1.7, 3.7, 4.3_

- [x] 9. Create Vbook communication client

  - Implement VbookClient interface for HTTP/TCP communication
  - Add SendTestRequest method for script testing protocol
  - Implement SendInstallRequest method for extension installation
  - Write unit tests with mock HTTP servers
  - _Requirements: 1.3, 1.5, 3.5_

- [x] 10. Build script testing functionality

  - Implement VbookTester interface with TestScript method
  - Create StartLocalServer method to serve project files on appPort-10
  - Add TCP connection logic for HTTP-over-TCP communication
  - Implement test parameter handling for --params flag
  - Write integration tests for complete testing workflow
  - _Requirements: 1.1, 1.4, 1.6, 1.8_

- [x] 11. Implement extension installation functionality

  - Create VbookInstaller interface implementation
  - Add PreparePluginData method to read project files and encode icon as base64
  - Implement InstallExtension method with HTTP GET to /install endpoint
  - Write unit tests for plugin data preparation and installation logic
  - _Requirements: 3.1, 3.4, 3.6_

- [x] 12. Create CLI command implementations

  - Implement test command with script-path and app-url arguments
  - Create build command with optional project-path (defaults to current directory)
  - Add install command with project-path and app-url arguments
  - Implement global flags (--help, --version, --verbose)
  - _Requirements: 1.1, 2.1, 3.1, 6.1, 6.2, 6.3_

- [x] 13. Implement comprehensive error handling

  - Create custom error types (ValidationError, NetworkError, FileSystemError)
  - Add ProjectValidationError and URLValidationError with specific messaging
  - Implement error wrapping with context for all operations
  - Write tests for error scenarios and message formatting
  - _Requirements: 1.7, 2.6, 3.7, 6.5, 7.3_

- [x] 14. Add help system and command validation

  - Implement comprehensive help text with examples for each command
  - Add command-specific help with --help flag support
  - Create input validation for all command arguments and flags
  - Write tests for help system and validation logic
  - _Requirements: 6.1, 6.2, 6.4, 6.5_

- [x] 15. Implement logging and verbose output

  - Create Logger interface with Info, Error, and Debug methods
  - Add structured logging with different levels based on --verbose flag
  - Implement progress indicators for long-running operations
  - Write tests for logging functionality
  - _Requirements: 7.1, 7.3_

- [x] 16. Finalize CLI tool with build and distribution setup

  - Set up cross-compilation for multiple platforms (Windows, macOS, Linux)
  - Create build scripts for automated releases
  - Ensure single binary distribution
  - Write integration tests for complete CLI workflows
  - _Requirements: 4.1, 4.5_
