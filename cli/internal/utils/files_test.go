package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileUtils_CopyFile(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	testContent := "Hello, World!"
	os.WriteFile(srcFile, []byte(testContent), 0o644)

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{"valid copy", srcFile, filepath.Join(tempDir, "dest.txt"), false},
		{"non-existent source", filepath.Join(tempDir, "non-existent.txt"), filepath.Join(tempDir, "dest2.txt"), true},
		{"invalid destination", srcFile, filepath.Join(tempDir, "nonexistent", "dest.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.CopyFile(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was copied
				if _, err := os.Stat(tt.dst); os.IsNotExist(err) {
					t.Errorf("Destination file was not created: %s", tt.dst)
					return
				}

				// Verify content matches
				dstContent, err := os.ReadFile(tt.dst)
				if err != nil {
					t.Errorf("Failed to read destination file: %v", err)
					return
				}

				if string(dstContent) != testContent {
					t.Errorf("File content mismatch: got %s, want %s", string(dstContent), testContent)
				}
			}
		})
	}
}

func TestFileUtils_CopyDirectory(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dir-copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "source")
	createTestDirectoryStructure(t, srcDir)

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{"valid copy", srcDir, filepath.Join(tempDir, "dest"), false},
		{"non-existent source", filepath.Join(tempDir, "non-existent"), filepath.Join(tempDir, "dest2"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.CopyDirectory(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory structure was copied
				verifyDirectoryStructure(t, tt.dst)
			}
		})
	}
}

func TestFileUtils_ValidateFile(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "validate-file-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0o644)

	// Create test directory
	testDir := filepath.Join(tempDir, "testdir")
	os.MkdirAll(testDir, 0o755)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid file", testFile, false},
		{"directory instead of file", testDir, true},
		{"non-existent file", filepath.Join(tempDir, "non-existent.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.ValidateFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileUtils_ValidateDirectory(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "validate-dir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0o644)

	// Create test directory
	testDir := filepath.Join(tempDir, "testdir")
	os.MkdirAll(testDir, 0o755)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid directory", testDir, false},
		{"file instead of directory", testFile, true},
		{"non-existent directory", filepath.Join(tempDir, "non-existent"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.ValidateDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileUtils_ReadJSONFile(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "json-read-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create valid JSON file
	validJSON := map[string]interface{}{
		"name":    "test",
		"version": 1,
		"enabled": true,
	}
	validJSONFile := filepath.Join(tempDir, "valid.json")
	data, _ := json.Marshal(validJSON)
	os.WriteFile(validJSONFile, data, 0o644)

	// Create invalid JSON file
	invalidJSONFile := filepath.Join(tempDir, "invalid.json")
	os.WriteFile(invalidJSONFile, []byte(`{"name": "test", "invalid": }`), 0o644)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid JSON", validJSONFile, false},
		{"invalid JSON", invalidJSONFile, true},
		{"non-existent file", filepath.Join(tempDir, "non-existent.json"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := fu.ReadJSONFile(tt.path, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify JSON was parsed correctly
				if result["name"] != "test" {
					t.Errorf("JSON parsing failed: expected name=test, got %v", result["name"])
				}
			}
		})
	}
}

func TestFileUtils_WriteJSONFile(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "json-write-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testData := map[string]interface{}{
		"name":    "test",
		"version": 1,
		"enabled": true,
	}

	tests := []struct {
		name    string
		path    string
		data    interface{}
		wantErr bool
	}{
		{"valid data", filepath.Join(tempDir, "output.json"), testData, false},
		{"invalid path", filepath.Join(tempDir, "nonexistent", "output.json"), testData, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.WriteJSONFile(tt.path, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created and contains correct data
				var result map[string]interface{}
				if err := fu.ReadJSONFile(tt.path, &result); err != nil {
					t.Errorf("Failed to read written JSON file: %v", err)
					return
				}

				if result["name"] != "test" {
					t.Errorf("Written JSON incorrect: expected name=test, got %v", result["name"])
				}
			}
		})
	}
}

func TestFileUtils_EnsureDirectory(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ensure-dir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"new directory", filepath.Join(tempDir, "newdir"), false},
		{"nested directory", filepath.Join(tempDir, "nested", "deep", "dir"), false},
		{"existing directory", tempDir, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fu.EnsureDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory exists
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("Directory was not created: %s", tt.path)
				}
			}
		})
	}
}

func TestFileUtils_NormalizePath(t *testing.T) {
	fu := NewFileUtils()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"relative path", "test/path", filepath.Join(getCurrentDir(t), "test", "path")},
		{"path with dots", "test/../path", filepath.Join(getCurrentDir(t), "path")},
		{"already absolute", filepath.Join(getCurrentDir(t), "test"), filepath.Join(getCurrentDir(t), "test")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fu.NormalizePath(tt.path)
			if result != tt.expected {
				t.Errorf("NormalizePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileUtils_IsExecutable(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "executable-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	regularFile := filepath.Join(tempDir, "regular.txt")
	os.WriteFile(regularFile, []byte("test"), 0o644)

	var executableFile string
	if runtime.GOOS == "windows" {
		executableFile = filepath.Join(tempDir, "test.exe")
		os.WriteFile(executableFile, []byte("test"), 0o644)
	} else {
		executableFile = filepath.Join(tempDir, "executable")
		os.WriteFile(executableFile, []byte("#!/bin/bash\necho test"), 0o755)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"regular file", regularFile, false},
		{"executable file", executableFile, true},
		{"non-existent file", filepath.Join(tempDir, "non-existent"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fu.IsExecutable(tt.path)
			if result != tt.expected {
				t.Errorf("IsExecutable() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileUtils_GetFileSize(t *testing.T) {
	fu := NewFileUtils()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filesize-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file with known content
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"
	os.WriteFile(testFile, []byte(testContent), 0o644)

	tests := []struct {
		name         string
		path         string
		expectedSize int64
		wantErr      bool
	}{
		{"valid file", testFile, int64(len(testContent)), false},
		{"non-existent file", filepath.Join(tempDir, "non-existent.txt"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := fu.GetFileSize(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && size != tt.expectedSize {
				t.Errorf("GetFileSize() = %v, want %v", size, tt.expectedSize)
			}
		})
	}
}

func TestFileUtils_JoinPath(t *testing.T) {
	fu := NewFileUtils()

	tests := []struct {
		name     string
		elements []string
		expected string
	}{
		{"simple join", []string{"path", "to", "file"}, filepath.Join("path", "to", "file")},
		{"single element", []string{"file"}, "file"},
		{"empty elements", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fu.JoinPath(tt.elements...)
			if result != tt.expected {
				t.Errorf("JoinPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileUtils_SplitPath(t *testing.T) {
	fu := NewFileUtils()

	tests := []struct {
		name         string
		path         string
		expectedDir  string
		expectedFile string
	}{
		{"simple path", filepath.Join("path", "to", "file.txt"), filepath.Join("path", "to") + string(filepath.Separator), "file.txt"},
		{"root file", filepath.Join(string(filepath.Separator), "file.txt"), string(filepath.Separator), "file.txt"},
		{"file only", "file.txt", "", "file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, file := fu.SplitPath(tt.path)
			if dir != tt.expectedDir || file != tt.expectedFile {
				t.Errorf("SplitPath() = (%v, %v), want (%v, %v)", dir, file, tt.expectedDir, tt.expectedFile)
			}
		})
	}
}

// Helper functions for testing

func createTestDirectoryStructure(_ *testing.T, baseDir string) {
	os.MkdirAll(baseDir, 0o755)

	// Create files in root
	os.WriteFile(filepath.Join(baseDir, "file1.txt"), []byte("content1"), 0o644)
	os.WriteFile(filepath.Join(baseDir, "file2.txt"), []byte("content2"), 0o644)

	// Create subdirectory with files
	subDir := filepath.Join(baseDir, "subdir")
	os.MkdirAll(subDir, 0o755)
	os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("content3"), 0o644)

	// Create nested subdirectory
	nestedDir := filepath.Join(subDir, "nested")
	os.MkdirAll(nestedDir, 0o755)
	os.WriteFile(filepath.Join(nestedDir, "file4.txt"), []byte("content4"), 0o644)
}

func verifyDirectoryStructure(t *testing.T, baseDir string) {
	expectedFiles := []string{
		"file1.txt",
		"file2.txt",
		filepath.Join("subdir", "file3.txt"),
		filepath.Join("subdir", "nested", "file4.txt"),
	}

	for _, expectedFile := range expectedFiles {
		fullPath := filepath.Join(baseDir, expectedFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not found: %s", fullPath)
		}
	}
}

func getCurrentDir(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	return dir
}
