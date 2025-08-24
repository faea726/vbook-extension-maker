package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type FileUtils struct{}

func NewFileUtils() *FileUtils {
	return &FileUtils{}
}

// CopyDirectory copies a directory and all its contents to the destination
func (fu *FileUtils) CopyDirectory(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return NewFileSystemError("stat", src, err)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return NewFileSystemError("create", dst, err)
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return NewFileSystemError("read", src, err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := fu.CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := fu.CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst
func (fu *FileUtils) CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return NewFileSystemError("read", src, err)
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return NewFileSystemError("stat", src, err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return NewFileSystemError("write", dst, err)
	}
	defer dstFile.Close()

	// Copy file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return NewFileSystemError("write", dst, err)
	}

	return nil
}

// ValidateFile checks if a file exists and is accessible
func (fu *FileUtils) ValidateFile(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return NewFileSystemError("stat", path, fmt.Errorf("file does not exist"))
	}
	if err != nil {
		return NewFileSystemError("stat", path, err)
	}
	if info.IsDir() {
		return NewFileSystemError("stat", path, fmt.Errorf("path is a directory, not a file"))
	}
	return nil
}

// ValidateDirectory checks if a directory exists and is accessible
func (fu *FileUtils) ValidateDirectory(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return NewFileSystemError("stat", path, fmt.Errorf("directory does not exist"))
	}
	if err != nil {
		return NewFileSystemError("stat", path, err)
	}
	if !info.IsDir() {
		return NewFileSystemError("stat", path, fmt.Errorf("path is a file, not a directory"))
	}
	return nil
}

// ReadJSONFile reads and parses a JSON file into the target interface
func (fu *FileUtils) ReadJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return NewFileSystemError("read", path, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return NewFileSystemError("parse", path, fmt.Errorf("invalid JSON: %w", err))
	}

	return nil
}

// WriteJSONFile writes an interface as JSON to a file
func (fu *FileUtils) WriteJSONFile(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return NewFileSystemError("marshal", path, err)
	}

	if err := os.WriteFile(path, jsonData, 0o644); err != nil {
		return NewFileSystemError("write", path, err)
	}

	return nil
}

// EnsureDirectory creates a directory if it doesn't exist
func (fu *FileUtils) EnsureDirectory(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return NewFileSystemError("create", path, err)
	}
	return nil
}

// GetExecutablePath returns the path to the current executable
func (fu *FileUtils) GetExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}

// NormalizePath normalizes a file path for the current platform
func (fu *FileUtils) NormalizePath(path string) string {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	// Clean the path (removes . and .. elements)
	return filepath.Clean(absPath)
}

// IsExecutable checks if a file is executable on the current platform
func (fu *FileUtils) IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// On Windows, check file extension
	if runtime.GOOS == "windows" {
		ext := filepath.Ext(path)
		return ext == ".exe" || ext == ".bat" || ext == ".cmd"
	}

	// On Unix-like systems, check execute permission
	return info.Mode()&0o111 != 0
}

// GetTempDir returns a platform-appropriate temporary directory
func (fu *FileUtils) GetTempDir() string {
	return os.TempDir()
}

// JoinPath joins path elements using the platform-appropriate separator
func (fu *FileUtils) JoinPath(elements ...string) string {
	return filepath.Join(elements...)
}

// SplitPath splits a path into directory and filename
func (fu *FileUtils) SplitPath(path string) (dir, file string) {
	return filepath.Split(path)
}

// GetFileSize returns the size of a file in bytes
func (fu *FileUtils) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, NewFileSystemError("stat", path, err)
	}
	return info.Size(), nil
}
