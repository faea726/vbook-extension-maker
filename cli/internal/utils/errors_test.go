package utils

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		message  string
		value    string
		expected string
	}{
		{
			name:     "error with empty value",
			field:    "project_name",
			message:  "cannot be empty",
			value:    "",
			expected: "validation failed for project_name: cannot be empty",
		},
		{
			name:     "error without value",
			field:    "url",
			message:  "invalid format",
			value:    "",
			expected: "validation failed for url: invalid format",
		},
		{
			name:     "error with non-empty value",
			field:    "port",
			message:  "must be between 1 and 65535",
			value:    "99999",
			expected: "validation failed for port: must be between 1 and 65535 (value: 99999)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.message, tt.value)
			if err.Error() != tt.expected {
				t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), tt.expected)
			}

			// Test field access
			if err.Field != tt.field {
				t.Errorf("ValidationError.Field = %v, want %v", err.Field, tt.field)
			}
			if err.Message != tt.message {
				t.Errorf("ValidationError.Message = %v, want %v", err.Message, tt.message)
			}
			if err.Value != tt.value {
				t.Errorf("ValidationError.Value = %v, want %v", err.Value, tt.value)
			}
		})
	}
}

func TestNetworkError(t *testing.T) {
	cause := errors.New("connection refused")

	tests := []struct {
		name      string
		operation string
		url       string
		cause     error
		expected  string
	}{
		{
			name:      "error with URL",
			operation: "connect",
			url:       "http://192.168.1.100:8080",
			cause:     cause,
			expected:  "network error during connect to http://192.168.1.100:8080: connection refused",
		},
		{
			name:      "error without URL",
			operation: "resolve",
			url:       "",
			cause:     cause,
			expected:  "network error during resolve: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNetworkError(tt.operation, tt.url, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("NetworkError.Error() = %v, want %v", err.Error(), tt.expected)
			}

			// Test unwrapping
			if errors.Unwrap(err) != tt.cause {
				t.Errorf("NetworkError.Unwrap() = %v, want %v", errors.Unwrap(err), tt.cause)
			}

			// Test field access
			if err.Operation != tt.operation {
				t.Errorf("NetworkError.Operation = %v, want %v", err.Operation, tt.operation)
			}
			if err.URL != tt.url {
				t.Errorf("NetworkError.URL = %v, want %v", err.URL, tt.url)
			}
		})
	}
}

func TestFileSystemError(t *testing.T) {
	cause := errors.New("permission denied")

	tests := []struct {
		name      string
		operation string
		path      string
		cause     error
		expected  string
	}{
		{
			name:      "error with path",
			operation: "read",
			path:      "/path/to/file.txt",
			cause:     cause,
			expected:  "file system error during read on /path/to/file.txt: permission denied",
		},
		{
			name:      "error without path",
			operation: "write",
			path:      "",
			cause:     cause,
			expected:  "file system error during write: permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFileSystemError(tt.operation, tt.path, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("FileSystemError.Error() = %v, want %v", err.Error(), tt.expected)
			}

			// Test unwrapping
			if errors.Unwrap(err) != tt.cause {
				t.Errorf("FileSystemError.Unwrap() = %v, want %v", errors.Unwrap(err), tt.cause)
			}

			// Test field access
			if err.Operation != tt.operation {
				t.Errorf("FileSystemError.Operation = %v, want %v", err.Operation, tt.operation)
			}
			if err.Path != tt.path {
				t.Errorf("FileSystemError.Path = %v, want %v", err.Path, tt.path)
			}
		})
	}
}

func TestProjectError(t *testing.T) {
	cause := errors.New("invalid structure")

	tests := []struct {
		name        string
		projectPath string
		message     string
		cause       error
		expected    string
	}{
		{
			name:        "error with cause",
			projectPath: "/path/to/project",
			message:     "validation failed",
			cause:       cause,
			expected:    "project error in /path/to/project: validation failed (invalid structure)",
		},
		{
			name:        "error without cause",
			projectPath: "/path/to/project",
			message:     "missing files",
			cause:       nil,
			expected:    "project error in /path/to/project: missing files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewProjectError(tt.projectPath, tt.message, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("ProjectError.Error() = %v, want %v", err.Error(), tt.expected)
			}

			// Test unwrapping
			if errors.Unwrap(err) != tt.cause {
				t.Errorf("ProjectError.Unwrap() = %v, want %v", errors.Unwrap(err), tt.cause)
			}

			// Test field access
			if err.ProjectPath != tt.projectPath {
				t.Errorf("ProjectError.ProjectPath = %v, want %v", err.ProjectPath, tt.projectPath)
			}
			if err.Message != tt.message {
				t.Errorf("ProjectError.Message = %v, want %v", err.Message, tt.message)
			}
		})
	}
}

func TestVbookError(t *testing.T) {
	cause := errors.New("connection timeout")

	tests := []struct {
		name      string
		operation string
		appURL    string
		message   string
		cause     error
		expected  string
	}{
		{
			name:      "error with cause",
			operation: "test",
			appURL:    "http://192.168.1.100:8080",
			message:   "request failed",
			cause:     cause,
			expected:  "Vbook error during test to http://192.168.1.100:8080: request failed (connection timeout)",
		},
		{
			name:      "error without cause",
			operation: "install",
			appURL:    "http://192.168.1.100:8080",
			message:   "invalid response",
			cause:     nil,
			expected:  "Vbook error during install to http://192.168.1.100:8080: invalid response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewVbookError(tt.operation, tt.appURL, tt.message, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("VbookError.Error() = %v, want %v", err.Error(), tt.expected)
			}

			// Test unwrapping
			if errors.Unwrap(err) != tt.cause {
				t.Errorf("VbookError.Unwrap() = %v, want %v", errors.Unwrap(err), tt.cause)
			}

			// Test field access
			if err.Operation != tt.operation {
				t.Errorf("VbookError.Operation = %v, want %v", err.Operation, tt.operation)
			}
			if err.AppURL != tt.appURL {
				t.Errorf("VbookError.AppURL = %v, want %v", err.AppURL, tt.appURL)
			}
			if err.Message != tt.message {
				t.Errorf("VbookError.Message = %v, want %v", err.Message, tt.message)
			}
		})
	}
}

func TestGetUserFriendlyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ValidationError",
			err:      NewValidationError("project_name", "cannot be empty", ""),
			expected: "Invalid project_name: cannot be empty",
		},
		{
			name:     "NetworkError with URL",
			err:      NewNetworkError("connect", "http://192.168.1.100:8080", errors.New("connection refused")),
			expected: "Failed to connect to http://192.168.1.100:8080. Please check the URL and ensure the Vbook app is running.",
		},
		{
			name:     "NetworkError without URL",
			err:      NewNetworkError("resolve", "", errors.New("no such host")),
			expected: "Network connection failed. Please check your internet connection.",
		},
		{
			name:     "FileSystemError read",
			err:      NewFileSystemError("read", "/path/to/file.txt", errors.New("permission denied")),
			expected: "Cannot read file: /path/to/file.txt. Please check if the file exists and you have permission to read it.",
		},
		{
			name:     "FileSystemError write",
			err:      NewFileSystemError("write", "/path/to/file.txt", errors.New("permission denied")),
			expected: "Cannot write file: /path/to/file.txt. Please check if you have permission to write to this location.",
		},
		{
			name:     "FileSystemError create",
			err:      NewFileSystemError("create", "/path/to/dir", errors.New("permission denied")),
			expected: "Cannot create directory: /path/to/dir. Please check if you have permission to create directories here.",
		},
		{
			name:     "FileSystemError other operation",
			err:      NewFileSystemError("delete", "/path/to/file.txt", errors.New("permission denied")),
			expected: "File system error: file system error during delete on /path/to/file.txt: permission denied",
		},
		{
			name:     "ProjectError",
			err:      NewProjectError("/path/to/project", "invalid structure", nil),
			expected: "Project error: invalid structure",
		},
		{
			name:     "VbookError test",
			err:      NewVbookError("test", "http://192.168.1.100:8080", "script failed", nil),
			expected: "Failed to test script with Vbook app at http://192.168.1.100:8080: script failed",
		},
		{
			name:     "VbookError install",
			err:      NewVbookError("install", "http://192.168.1.100:8080", "installation failed", nil),
			expected: "Failed to install extension to Vbook app at http://192.168.1.100:8080: installation failed",
		},
		{
			name:     "VbookError other operation",
			err:      NewVbookError("unknown", "http://192.168.1.100:8080", "unknown error", nil),
			expected: "Vbook communication error: unknown error",
		},
		{
			name:     "Generic error",
			err:      errors.New("generic error message"),
			expected: "generic error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserFriendlyError(tt.err)
			if result != tt.expected {
				t.Errorf("GetUserFriendlyError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorChaining(t *testing.T) {
	// Test that errors can be properly chained and unwrapped
	originalErr := errors.New("original error")
	networkErr := NewNetworkError("connect", "http://example.com", originalErr)

	// Test that we can unwrap to get the original error
	if !errors.Is(networkErr, originalErr) {
		t.Error("NetworkError should wrap the original error")
	}

	// Test with errors.As
	var netErr *NetworkError
	if !errors.As(networkErr, &netErr) {
		t.Error("Should be able to extract NetworkError with errors.As")
	}

	if netErr.Operation != "connect" {
		t.Errorf("NetworkError.Operation = %v, want connect", netErr.Operation)
	}
}

func TestErrorTypeAssertions(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectType  string
		expectField string
	}{
		{
			name:        "ValidationError type assertion",
			err:         NewValidationError("field", "message", "value"),
			expectType:  "ValidationError",
			expectField: "field",
		},
		{
			name:        "NetworkError type assertion",
			err:         NewNetworkError("operation", "url", nil),
			expectType:  "NetworkError",
			expectField: "operation",
		},
		{
			name:        "FileSystemError type assertion",
			err:         NewFileSystemError("operation", "path", nil),
			expectType:  "FileSystemError",
			expectField: "operation",
		},
		{
			name:        "ProjectError type assertion",
			err:         NewProjectError("path", "message", nil),
			expectType:  "ProjectError",
			expectField: "path",
		},
		{
			name:        "VbookError type assertion",
			err:         NewVbookError("operation", "url", "message", nil),
			expectType:  "VbookError",
			expectField: "operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch e := tt.err.(type) {
			case *ValidationError:
				if tt.expectType != "ValidationError" {
					t.Errorf("Expected %s, got ValidationError", tt.expectType)
				}
				if e.Field != tt.expectField {
					t.Errorf("Expected field %s, got %s", tt.expectField, e.Field)
				}
			case *NetworkError:
				if tt.expectType != "NetworkError" {
					t.Errorf("Expected %s, got NetworkError", tt.expectType)
				}
				if e.Operation != tt.expectField {
					t.Errorf("Expected operation %s, got %s", tt.expectField, e.Operation)
				}
			case *FileSystemError:
				if tt.expectType != "FileSystemError" {
					t.Errorf("Expected %s, got FileSystemError", tt.expectType)
				}
				if e.Operation != tt.expectField {
					t.Errorf("Expected operation %s, got %s", tt.expectField, e.Operation)
				}
			case *ProjectError:
				if tt.expectType != "ProjectError" {
					t.Errorf("Expected %s, got ProjectError", tt.expectType)
				}
				if e.ProjectPath != tt.expectField {
					t.Errorf("Expected path %s, got %s", tt.expectField, e.ProjectPath)
				}
			case *VbookError:
				if tt.expectType != "VbookError" {
					t.Errorf("Expected %s, got VbookError", tt.expectType)
				}
				if e.Operation != tt.expectField {
					t.Errorf("Expected operation %s, got %s", tt.expectField, e.Operation)
				}
			default:
				t.Errorf("Unexpected error type: %T", tt.err)
			}
		})
	}
}
