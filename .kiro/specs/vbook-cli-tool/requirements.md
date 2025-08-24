# Requirements Document

## Introduction

This document outlines the requirements for creating a cross-platform CLI tool based on the existing Vbook Extension Maker VSCode extension. The CLI tool will provide test, build, and install functionality as a standalone command-line application written in Go. The tool will enable developers to test, build, and install existing Vbook extensions from any directory without requiring VSCode. Project creation functionality is not included as developers can use the VSCode extension for that purpose.

## Requirements

### Requirement 1

**User Story:** As a developer, I want to test Vbook extension scripts from the command line, so that I can validate my extension functionality without using VSCode.

#### Acceptance Criteria

1. WHEN the user runs `vbook-cli test <script-path> --app-url <vbook-app-url>` THEN the system SHALL test the specified script against the Vbook app
2. WHEN testing a script THEN the system SHALL validate that plugin.json exists in the project root
3. WHEN testing a script THEN the system SHALL establish a TCP connection to the Vbook app using the provided URL
4. WHEN testing a script THEN the system SHALL start a local HTTP server for script serving
5. WHEN testing a script THEN the system SHALL send the script content and metadata to the Vbook app via HTTP request
6. WHEN the test completes THEN the system SHALL display the response from the Vbook app including any output or errors
7. IF the Vbook app URL is invalid or unreachable THEN the system SHALL display an appropriate error message
8. WHEN testing with parameters THEN the system SHALL accept optional `--params` flag to pass test parameters

### Requirement 2

**User Story:** As a developer, I want to build Vbook extensions into distributable packages from the command line, so that I can create plugin.zip files for distribution.

#### Acceptance Criteria

1. WHEN the user runs `vbook-cli build <project-path>` THEN the system SHALL create a plugin.zip file in the project root
2. WHEN building an extension THEN the system SHALL validate that required files exist (plugin.json, icon.png, src directory)
3. WHEN building an extension THEN the system SHALL include plugin.json, icon.png, and all files from the src directory in the zip
4. WHEN building an extension THEN the system SHALL use maximum compression level for the zip file
5. WHEN build is successful THEN the system SHALL display the output file path and size
6. IF required files are missing THEN the system SHALL display specific error messages indicating which files are missing

### Requirement 3

**User Story:** As a developer, I want to install Vbook extensions directly to the Vbook app from the command line, so that I can deploy extensions for testing without manual installation.

#### Acceptance Criteria

1. WHEN the user runs `vbook-cli install <project-path> --app-url <vbook-app-url>` THEN the system SHALL install the extension to the specified Vbook app
2. WHEN installing an extension THEN the system SHALL validate that plugin.json and icon.png exist in the project root
3. WHEN installing an extension THEN the system SHALL read all JavaScript files from the src directory
4. WHEN installing an extension THEN the system SHALL prepare plugin data including metadata, icon as base64, and script content
5. WHEN installing an extension THEN the system SHALL send an HTTP GET request to the Vbook app's /install endpoint with plugin data
6. WHEN installation completes THEN the system SHALL display a success or error message based on the response
7. IF the Vbook app URL is invalid or unreachable THEN the system SHALL display an appropriate error message

### Requirement 4

**User Story:** As a developer, I want the CLI tool to work on all major platforms (Windows, macOS, Linux), so that I can use it regardless of my operating system.

#### Acceptance Criteria

1. WHEN the CLI tool is built THEN it SHALL compile and run on Windows, macOS, and Linux
2. WHEN handling file paths THEN the system SHALL use platform-appropriate path separators
3. WHEN creating network connections THEN the system SHALL work consistently across all platforms
4. WHEN displaying output THEN the system SHALL handle platform-specific terminal capabilities appropriately
5. WHEN the tool is distributed THEN it SHALL be available as standalone executables for each platform

### Requirement 5

**User Story:** As a developer, I want the CLI tool to accept a project path as an argument, so that I can operate on Vbook extension projects from any directory.

#### Acceptance Criteria

1. WHEN any command requires a project path THEN the system SHALL accept both absolute and relative paths
2. WHEN a relative path is provided THEN the system SHALL resolve it relative to the current working directory
3. WHEN an invalid path is provided THEN the system SHALL display a clear error message
4. WHEN no path is provided for commands that require it THEN the system SHALL use the current working directory as default
5. WHEN validating project paths THEN the system SHALL check for the presence of plugin.json to confirm it's a valid Vbook extension project

### Requirement 6

**User Story:** As a developer, I want clear help and usage information for the CLI tool, so that I can understand how to use all available commands and options.

#### Acceptance Criteria

1. WHEN the user runs `vbook-cli --help` or `vbook-cli -h` THEN the system SHALL display comprehensive usage information
2. WHEN the user runs `vbook-cli <command> --help` THEN the system SHALL display help specific to that command
3. WHEN the user runs `vbook-cli --version` or `vbook-cli -v` THEN the system SHALL display the current version number
4. WHEN displaying help THEN the system SHALL include examples for each command
5. WHEN an invalid command is entered THEN the system SHALL display an error message and suggest using --help

### Requirement 7

**User Story:** As a developer, I want the CLI tool to follow the KISS (Keep It Simple, Stupid) principle, so that the codebase is maintainable and the tool is easy to understand and use.

#### Acceptance Criteria

1. WHEN implementing functionality THEN the system SHALL use straightforward, readable code without unnecessary complexity
2. WHEN designing the command interface THEN the system SHALL use intuitive command names and argument patterns
3. WHEN handling errors THEN the system SHALL provide clear, actionable error messages
4. WHEN structuring the code THEN the system SHALL separate concerns into distinct, focused modules
5. WHEN adding features THEN the system SHALL avoid feature creep and maintain focus on core functionality
