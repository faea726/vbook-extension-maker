package vbook

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestVbookInstaller_PreparePluginData(t *testing.T) {
	installer := NewVbookInstaller()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		setupFunc func(string)
		wantErr   bool
		validate  func(*testing.T, *PluginData)
	}{
		{
			name: "valid project",
			setupFunc: func(dir string) {
				createValidInstallerProject(t, dir)
			},
			wantErr: false,
			validate: func(t *testing.T, data *PluginData) {
				if data.Name != "Test Extension" {
					t.Errorf("Expected name 'Test Extension', got %s", data.Name)
				}
				if data.Author != "Test Author" {
					t.Errorf("Expected author 'Test Author', got %s", data.Author)
				}
				if data.Version != "1" {
					t.Errorf("Expected version '1', got %s", data.Version)
				}
				if data.Source == "" {
					t.Error("Expected source to be set")
				}
				if !data.Enabled {
					t.Error("Expected plugin to be enabled")
				}
				if !data.Debug {
					t.Error("Expected plugin to be in debug mode")
				}
				if data.Icon == "" {
					t.Error("Expected icon data to be present")
				}
				if data.Data == "" {
					t.Error("Expected script data to be present")
				}

				// Verify script data is valid JSON
				var scriptData map[string]string
				if err := json.Unmarshal([]byte(data.Data), &scriptData); err != nil {
					t.Errorf("Script data is not valid JSON: %v", err)
				}

				// Verify script files are included
				if _, exists := scriptData["main.js"]; !exists {
					t.Error("Expected main.js to be in script data")
				}
				if _, exists := scriptData["utils.js"]; !exists {
					t.Error("Expected utils.js to be in script data")
				}
			},
		},
		{
			name: "missing plugin.json",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "main.js"), []byte("console.log('test');"), 0o644)
			},
			wantErr: true,
		},
		{
			name: "missing icon.png",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				createValidPluginJsonForInstaller(t, filepath.Join(dir, "plugin.json"))
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "main.js"), []byte("console.log('test');"), 0o644)
			},
			wantErr: true,
		},
		{
			name: "missing src directory",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				createValidPluginJsonForInstaller(t, filepath.Join(dir, "plugin.json"))
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
			},
			wantErr: true,
		},
		{
			name: "no JavaScript files in src",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				createValidPluginJsonForInstaller(t, filepath.Join(dir, "plugin.json"))
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "readme.txt"), []byte("not a js file"), 0o644)
			},
			wantErr: true,
		},
		{
			name: "invalid plugin.json",
			setupFunc: func(dir string) {
				os.MkdirAll(dir, 0o755)
				os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"invalid": json}`), 0o644)
				os.WriteFile(filepath.Join(dir, "icon.png"), []byte("fake png"), 0o644)
				os.MkdirAll(filepath.Join(dir, "src"), 0o755)
				os.WriteFile(filepath.Join(dir, "src", "main.js"), []byte("console.log('test');"), 0o644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := filepath.Join(tempDir, tt.name)
			tt.setupFunc(projectDir)

			data, err := installer.PreparePluginData(projectDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("PreparePluginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, data)
			}
		})
	}
}

func TestVbookInstaller_readPluginConfig(t *testing.T) {
	installer := NewVbookInstaller()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create valid plugin.json
	validConfigDir := filepath.Join(tempDir, "valid")
	os.MkdirAll(validConfigDir, 0o755)
	createValidPluginJsonForInstaller(t, filepath.Join(validConfigDir, "plugin.json"))

	// Create invalid plugin.json
	invalidConfigDir := filepath.Join(tempDir, "invalid")
	os.MkdirAll(invalidConfigDir, 0o755)
	os.WriteFile(filepath.Join(invalidConfigDir, "plugin.json"), []byte(`{"invalid": json}`), 0o644)

	tests := []struct {
		name        string
		projectPath string
		wantErr     bool
		validate    func(*testing.T, *PluginConfig)
	}{
		{
			name:        "valid config",
			projectPath: validConfigDir,
			wantErr:     false,
			validate: func(t *testing.T, config *PluginConfig) {
				if config.Metadata.Name != "Test Extension" {
					t.Errorf("Expected name 'Test Extension', got %s", config.Metadata.Name)
				}
				if config.Metadata.Author != "Test Author" {
					t.Errorf("Expected author 'Test Author', got %s", config.Metadata.Author)
				}
				if config.Metadata.Version != 1 {
					t.Errorf("Expected version 1, got %d", config.Metadata.Version)
				}
			},
		},
		{
			name:        "invalid config",
			projectPath: invalidConfigDir,
			wantErr:     true,
		},
		{
			name:        "missing config",
			projectPath: filepath.Join(tempDir, "missing"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := installer.readPluginConfig(tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPluginConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

func TestVbookInstaller_readAndEncodeIcon(t *testing.T) {
	installer := NewVbookInstaller()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "icon-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test icon file
	iconData := []byte("fake png data")
	validIconDir := filepath.Join(tempDir, "valid")
	os.MkdirAll(validIconDir, 0o755)
	os.WriteFile(filepath.Join(validIconDir, "icon.png"), iconData, 0o644)

	tests := []struct {
		name        string
		projectPath string
		wantErr     bool
		validate    func(*testing.T, string)
	}{
		{
			name:        "valid icon",
			projectPath: validIconDir,
			wantErr:     false,
			validate: func(t *testing.T, iconBase64 string) {
				expectedPrefix := "data:image/*;base64,"
				if len(iconBase64) <= len(expectedPrefix) {
					t.Error("Icon base64 string is too short")
					return
				}

				if iconBase64[:len(expectedPrefix)] != expectedPrefix {
					t.Errorf("Expected prefix %s, got %s", expectedPrefix, iconBase64[:len(expectedPrefix)])
				}

				// Decode and verify content
				encodedData := iconBase64[len(expectedPrefix):]
				decodedData, err := base64.StdEncoding.DecodeString(encodedData)
				if err != nil {
					t.Errorf("Failed to decode base64: %v", err)
					return
				}

				if string(decodedData) != string(iconData) {
					t.Error("Decoded icon data doesn't match original")
				}
			},
		},
		{
			name:        "missing icon",
			projectPath: filepath.Join(tempDir, "missing"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iconBase64, err := installer.readAndEncodeIcon(tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndEncodeIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, iconBase64)
			}
		})
	}
}

func TestVbookInstaller_readScriptFiles(t *testing.T) {
	installer := NewVbookInstaller()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "script-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project with JavaScript files
	validScriptDir := filepath.Join(tempDir, "valid")
	os.MkdirAll(filepath.Join(validScriptDir, "src"), 0o755)
	os.WriteFile(filepath.Join(validScriptDir, "src", "main.js"), []byte("console.log('main');"), 0o644)
	os.WriteFile(filepath.Join(validScriptDir, "src", "utils.js"), []byte("function helper() {}"), 0o644)
	os.WriteFile(filepath.Join(validScriptDir, "src", "readme.txt"), []byte("not a js file"), 0o644) // Should be ignored

	// Create project with no JavaScript files
	noJsDir := filepath.Join(tempDir, "no-js")
	os.MkdirAll(filepath.Join(noJsDir, "src"), 0o755)
	os.WriteFile(filepath.Join(noJsDir, "src", "readme.txt"), []byte("not a js file"), 0o644)

	tests := []struct {
		name        string
		projectPath string
		wantErr     bool
		validate    func(*testing.T, map[string]string)
	}{
		{
			name:        "valid scripts",
			projectPath: validScriptDir,
			wantErr:     false,
			validate: func(t *testing.T, scripts map[string]string) {
				if len(scripts) != 2 {
					t.Errorf("Expected 2 script files, got %d", len(scripts))
				}

				if scripts["main.js"] != "console.log('main');" {
					t.Errorf("Unexpected content for main.js: %s", scripts["main.js"])
				}

				if scripts["utils.js"] != "function helper() {}" {
					t.Errorf("Unexpected content for utils.js: %s", scripts["utils.js"])
				}

				// Verify non-JS files are not included
				if _, exists := scripts["readme.txt"]; exists {
					t.Error("Non-JS file should not be included")
				}
			},
		},
		{
			name:        "no JavaScript files",
			projectPath: noJsDir,
			wantErr:     true,
		},
		{
			name:        "missing src directory",
			projectPath: filepath.Join(tempDir, "missing"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scripts, err := installer.readScriptFiles(tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readScriptFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, scripts)
			}
		})
	}
}

// Helper functions for testing

func createValidInstallerProject(t *testing.T, projectDir string) {
	os.MkdirAll(projectDir, 0o755)

	// Create plugin.json
	createValidPluginJsonForInstaller(t, filepath.Join(projectDir, "plugin.json"))

	// Create icon.png
	os.WriteFile(filepath.Join(projectDir, "icon.png"), []byte("fake png data"), 0o644)

	// Create src directory with JavaScript files
	srcDir := filepath.Join(projectDir, "src")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "main.js"), []byte("console.log('Hello, World!');"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "utils.js"), []byte("function helper() { return 'help'; }"), 0o644)
}

func createValidPluginJsonForInstaller(t *testing.T, path string) {
	config := PluginConfig{}
	config.Metadata.Name = "Test Extension"
	config.Metadata.Author = "Test Author"
	config.Metadata.Version = 1
	config.Metadata.Source = "test-source"
	config.Metadata.Description = "Test Description"
	config.Script.Detail = "detail.js"
	config.Script.Toc = "toc.js"
	config.Script.Chap = "chap.js"

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal plugin.json: %v", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("Failed to write plugin.json: %v", err)
	}
}
