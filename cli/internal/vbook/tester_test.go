package vbook

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestVbookTester_prepareInputParams(t *testing.T) {
	tester := NewVbookTester()

	tests := []struct {
		name     string
		params   []string
		expected [][]string
	}{
		{
			name:     "empty params",
			params:   []string{},
			expected: [][]string{{}},
		},
		{
			name:     "single param",
			params:   []string{"param1"},
			expected: [][]string{{"param1"}},
		},
		{
			name:     "multiple params",
			params:   []string{"param1", "param2"},
			expected: [][]string{{"param1"}, {"param2"}},
		},
		{
			name:     "comma-separated param",
			params:   []string{"param1,param2,param3"},
			expected: [][]string{{"param1", "param2", "param3"}},
		},
		{
			name:     "mixed params",
			params:   []string{"single", "comma,separated,values", "another"},
			expected: [][]string{{"single"}, {"comma", "separated", "values"}, {"another"}},
		},
		{
			name:     "params with whitespace",
			params:   []string{" param1 ", " comma , separated , values "},
			expected: [][]string{{"param1"}, {"comma", "separated", "values"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tester.prepareInputParams(tt.params)
			if !equalStringSlices(result, tt.expected) {
				t.Errorf("prepareInputParams() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVbookTester_TestScript(t *testing.T) {
	tester := NewVbookTester()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tester-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		setupProject  func(string) string                   // Returns script path
		setupServer   func() (net.Listener, func(), string) // Returns listener, cleanup, and URL
		params        []string
		wantErr       bool
		errorContains string
		validate      func(*testing.T, *TestResponse)
	}{
		{
			name: "successful test",
			setupProject: func(dir string) string {
				return createValidTesterProject(t, dir)
			},
			setupServer: func() (net.Listener, func(), string) {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					t.Fatalf("Failed to create listener: %v", err)
				}

				go func() {
					for {
						conn, err := listener.Accept()
						if err != nil {
							return // Listener closed
						}

						go func(c net.Conn) {
							defer c.Close()

							// Read request
							buffer := make([]byte, 4096)
							c.Read(buffer)

							// Send successful response
							response := TestResponse{
								Status: "success",
								Output: "Script executed successfully",
							}
							responseJSON, _ := json.Marshal(response)

							httpResponse := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s",
								len(responseJSON), string(responseJSON))
							c.Write([]byte(httpResponse))
						}(conn)
					}
				}()

				addr := listener.Addr().(*net.TCPAddr)
				url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
				return listener, func() { listener.Close() }, url
			},
			params:  []string{"test", "param"},
			wantErr: false,
			validate: func(t *testing.T, resp *TestResponse) {
				if resp.Status != "success" {
					t.Errorf("Expected status 'success', got %s", resp.Status)
				}
				if resp.Output != "Script executed successfully" {
					t.Errorf("Expected output 'Script executed successfully', got %s", resp.Output)
				}
			},
		},
		{
			name: "invalid project structure",
			setupProject: func(dir string) string {
				// Create invalid project (missing plugin.json)
				os.MkdirAll(dir, 0o755)
				scriptPath := filepath.Join(dir, "test.js")
				os.WriteFile(scriptPath, []byte("console.log('test');"), 0o644)
				return scriptPath
			},
			setupServer: func() (net.Listener, func(), string) {
				return nil, func() {}, "http://127.0.0.1:8080"
			},
			wantErr:       true,
			errorContains: "invalid workspace",
		},
		{
			name: "invalid app URL",
			setupProject: func(dir string) string {
				return createValidTesterProject(t, dir)
			},
			setupServer: func() (net.Listener, func(), string) {
				return nil, func() {}, "invalid-url"
			},
			wantErr:       true,
			errorContains: "invalid app URL",
		},
		{
			name: "connection failed",
			setupProject: func(dir string) string {
				return createValidTesterProject(t, dir)
			},
			setupServer: func() (net.Listener, func(), string) {
				return nil, func() {}, "http://127.0.0.1:9999" // Non-existent port
			},
			wantErr:       true,
			errorContains: "test request failed",
		},
		{
			name: "missing script file",
			setupProject: func(dir string) string {
				createValidTesterProject(t, dir)
				return filepath.Join(dir, "non-existent.js") // Return non-existent script path
			},
			setupServer: func() (net.Listener, func(), string) {
				return nil, func() {}, "http://127.0.0.1:8080"
			},
			wantErr:       true,
			errorContains: "failed to read script file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := filepath.Join(tempDir, tt.name)
			scriptPath := tt.setupProject(projectDir)

			var appURL string
			var cleanup func()

			if tt.name == "connection failed" || tt.name == "invalid app URL" || tt.name == "invalid project structure" || tt.name == "missing script file" {
				_, cleanup, appURL = tt.setupServer()
			} else {
				_, cleanupFunc, url := tt.setupServer()
				cleanup = cleanupFunc
				appURL = url

				// Give server time to start
				time.Sleep(100 * time.Millisecond)
			}
			defer cleanup()

			resp, err := tester.TestScript(scriptPath, appURL, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("TestScript() error = %v, should contain %v", err, tt.errorContains)
				}
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

func TestVbookTester_StartLocalServer(t *testing.T) {
	tester := NewVbookTester()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "server-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test project
	createValidTesterProject(t, tempDir)

	tests := []struct {
		name        string
		port        int
		projectPath string
		wantErr     bool
	}{
		{
			name:        "valid server start",
			port:        0, // Use any available port
			projectPath: tempDir,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Find an available port
			listener, err := net.Listen("tcp", ":0")
			if err != nil {
				t.Fatalf("Failed to find available port: %v", err)
			}
			port := listener.Addr().(*net.TCPAddr).Port
			listener.Close()

			server, err := tester.StartLocalServer(port, tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartLocalServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				defer server.Stop()

				// Give server time to start
				time.Sleep(200 * time.Millisecond)

				// Verify server is running by checking if port is in use
				conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				if err != nil {
					t.Errorf("Server is not listening on port %d: %v", port, err)
				} else {
					conn.Close()
				}
			}
		})
	}
}

// Helper functions for testing

func createValidTesterProject(_ *testing.T, projectDir string) string {
	os.MkdirAll(projectDir, 0o755)

	// Create plugin.json
	pluginJson := `{
		"metadata": {
			"name": "Test Extension",
			"author": "Test Author",
			"version": 1,
			"source": "test-source",
			"description": "Test Description"
		},
		"script": {
			"detail": "detail.js",
			"toc": "toc.js",
			"chap": "chap.js"
		}
	}`
	os.WriteFile(filepath.Join(projectDir, "plugin.json"), []byte(pluginJson), 0o644)

	// Create icon.png
	os.WriteFile(filepath.Join(projectDir, "icon.png"), []byte("fake png data"), 0o644)

	// Create src directory with JavaScript files
	srcDir := filepath.Join(projectDir, "src")
	os.MkdirAll(srcDir, 0o755)

	scriptContent := `console.log('Hello, World!');`
	scriptPath := filepath.Join(srcDir, "main.js")
	os.WriteFile(scriptPath, []byte(scriptContent), 0o644)

	os.WriteFile(filepath.Join(srcDir, "utils.js"), []byte("function helper() { return 'help'; }"), 0o644)

	return scriptPath
}

func equalStringSlices(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}
