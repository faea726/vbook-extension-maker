package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectValidator_ValidateProject(t *testing.T) {
	validator := NewProjectValidator()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vbook-validator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid project structure
	validProjectDir := filepath.Join(tempDir, "valid-project")
	createValidProject(t, validProjectDir)

	// Create an invalid project (missing files)
	invalidProjectDir := filepath.Join(tempDir, "invalid-project")
	os.MkdirAll(invalidProjectDir, 0o755)

	tests := []struct {
		name        string
		projectPath string
		wantErr     bool
	}{
		{"valid project", validProjectDir, false},
		{"invalid project", invalidProjectDir, true},
		{"non-existent path", filepath.Join(tempDir, "non-existent"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateProject(tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectValidator_CheckRequiredFiles(t *testing.T) {
	validator := NewProjectValidator()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vbook-files-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		setupFunc     func(string)
		wantErr       bool
		errorContains string
	}{
		{
			name: "all files present",
			setupFunc: func(dir string) {
				createValidProject(t, dir)
			},
			wantErr: false,
		},
		{
			name: "missing plugin.json",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "main.js"), []byte("console.log('test');"), 0o644)
			},
			wantErr:       true,
			errorContains: "plugin.json",
		},
		{
			name: "missing icon.png",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				createValidPluginJson(t, filepath.Join(dir, "plugin.json"))
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "main.js"), []byte("console.log('test');"), 0o644)
			},
			wantErr:       true,
			errorContains: "icon.png",
		},
		{
			name: "missing src directory",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				createValidPluginJson(t, filepath.Join(dir, "plugin.json"))
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
			},
			wantErr:       true,
			errorContains: "src",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tempDir, tt.name)
			tt.setupFunc(testDir)

			err := validator.ValidateProject(testDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProject() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errorContains != "" && err != nil {
				if !contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateProject() error = %v, should contain %v", err, tt.errorContains)
				}
			}
		})
	}
}

// Helper functions for testing

func createValidProject(t *testing.T, projectDir string) {
	os.MkdirAll(projectDir, 0o755)

	// Create plugin.json
	createValidPluginJson(t, filepath.Join(projectDir, "plugin.json"))

	// Create icon.png
	os.WriteFile(filepath.Join(projectDir, "icon.png"), []byte("fake png data"), 0o644)

	// Create src directory with a JavaScript file
	srcDir := filepath.Join(projectDir, "src")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "main.js"), []byte("console.log('Hello, World!');"), 0o644)
}

func createValidPluginJson(t *testing.T, path string) {
	config := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":        "Test Extension",
			"author":      "Test Author",
			"version":     1,
			"source":      "test-source",
			"description": "Test Description",
		},
		"script": map[string]interface{}{
			"detail": "detail.js",
			"toc":    "toc.js",
			"chap":   "chap.js",
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal plugin.json: %v", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("Failed to write plugin.json: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
