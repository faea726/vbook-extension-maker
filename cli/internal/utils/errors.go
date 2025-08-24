package utils

import "fmt"

// ValidationError represents validation-related errors
type ValidationError struct {
	Field   string
	Message string
	Value   string
}

func (e ValidationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("validation failed for %s: %s (value: %s)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message, value string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// NetworkError represents network-related errors
type NetworkError struct {
	Operation string
	URL       string
	Cause     error
}

func (e NetworkError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Cause)
	}
	return fmt.Sprintf("network error during %s: %v", e.Operation, e.Cause)
}

func (e NetworkError) Unwrap() error {
	return e.Cause
}

func NewNetworkError(operation, url string, cause error) *NetworkError {
	return &NetworkError{
		Operation: operation,
		URL:       url,
		Cause:     cause,
	}
}

// FileSystemError represents file system-related errors
type FileSystemError struct {
	Operation string
	Path      string
	Cause     error
}

func (e FileSystemError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("file system error during %s on %s: %v", e.Operation, e.Path, e.Cause)
	}
	return fmt.Sprintf("file system error during %s: %v", e.Operation, e.Cause)
}

func (e FileSystemError) Unwrap() error {
	return e.Cause
}

func NewFileSystemError(operation, path string, cause error) *FileSystemError {
	return &FileSystemError{
		Operation: operation,
		Path:      path,
		Cause:     cause,
	}
}

// ProjectError represents project-related errors
type ProjectError struct {
	ProjectPath string
	Message     string
	Cause       error
}

func (e ProjectError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("project error in %s: %s (%v)", e.ProjectPath, e.Message, e.Cause)
	}
	return fmt.Sprintf("project error in %s: %s", e.ProjectPath, e.Message)
}

func (e ProjectError) Unwrap() error {
	return e.Cause
}

func NewProjectError(projectPath, message string, cause error) *ProjectError {
	return &ProjectError{
		ProjectPath: projectPath,
		Message:     message,
		Cause:       cause,
	}
}

// VbookError represents Vbook app communication errors
type VbookError struct {
	Operation string
	AppURL    string
	Message   string
	Cause     error
}

func (e VbookError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Vbook error during %s to %s: %s (%v)", e.Operation, e.AppURL, e.Message, e.Cause)
	}
	return fmt.Sprintf("Vbook error during %s to %s: %s", e.Operation, e.AppURL, e.Message)
}

func (e VbookError) Unwrap() error {
	return e.Cause
}

func NewVbookError(operation, appURL, message string, cause error) *VbookError {
	return &VbookError{
		Operation: operation,
		AppURL:    appURL,
		Message:   message,
		Cause:     cause,
	}
}

// GetUserFriendlyError converts technical errors to user-friendly messages
func GetUserFriendlyError(err error) string {
	switch e := err.(type) {
	case *ValidationError:
		return fmt.Sprintf("Invalid %s: %s", e.Field, e.Message)
	case *NetworkError:
		if e.URL != "" {
			return fmt.Sprintf("Failed to connect to %s. Please check the URL and ensure the Vbook app is running.", e.URL)
		}
		return "Network connection failed. Please check your internet connection."
	case *FileSystemError:
		switch e.Operation {
		case "read":
			return fmt.Sprintf("Cannot read file: %s. Please check if the file exists and you have permission to read it.", e.Path)
		case "write":
			return fmt.Sprintf("Cannot write file: %s. Please check if you have permission to write to this location.", e.Path)
		case "create":
			return fmt.Sprintf("Cannot create directory: %s. Please check if you have permission to create directories here.", e.Path)
		default:
			return fmt.Sprintf("File system error: %s", e.Error())
		}
	case *ProjectError:
		return fmt.Sprintf("Project error: %s", e.Message)
	case *VbookError:
		switch e.Operation {
		case "test":
			return fmt.Sprintf("Failed to test script with Vbook app at %s: %s", e.AppURL, e.Message)
		case "install":
			return fmt.Sprintf("Failed to install extension to Vbook app at %s: %s", e.AppURL, e.Message)
		default:
			return fmt.Sprintf("Vbook communication error: %s", e.Message)
		}
	default:
		return err.Error()
	}
}
