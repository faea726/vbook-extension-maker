package vbook

import (
	"encoding/json"
	"testing"
)

func TestPluginConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  PluginConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: PluginConfig{
				Metadata: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     int    `json:"version"`
					Source      string `json:"source"`
					Regexp      string `json:"regexp"`
					Description string `json:"description"`
					Locale      string `json:"locale"`
					Tag         string `json:"tag"`
					Type        string `json:"type"`
				}{
					Name:    "Test Plugin",
					Author:  "Test Author",
					Version: 1,
					Source:  "test-source",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: PluginConfig{
				Metadata: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     int    `json:"version"`
					Source      string `json:"source"`
					Regexp      string `json:"regexp"`
					Description string `json:"description"`
					Locale      string `json:"locale"`
					Tag         string `json:"tag"`
					Type        string `json:"type"`
				}{
					Author:  "Test Author",
					Version: 1,
					Source:  "test-source",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid version",
			config: PluginConfig{
				Metadata: struct {
					Name        string `json:"name"`
					Author      string `json:"author"`
					Version     int    `json:"version"`
					Source      string `json:"source"`
					Regexp      string `json:"regexp"`
					Description string `json:"description"`
					Locale      string `json:"locale"`
					Tag         string `json:"tag"`
					Type        string `json:"type"`
				}{
					Name:    "Test Plugin",
					Author:  "Test Author",
					Version: 0,
					Source:  "test-source",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PluginConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTestRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request TestRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: TestRequest{
				IP:       "http://192.168.1.100:8070",
				Root:     "test-extension/src",
				Language: "javascript",
				Script:   "console.log('test');",
				Input:    [][]string{{"param1", "param2"}},
			},
			wantErr: false,
		},
		{
			name: "missing IP",
			request: TestRequest{
				Root:     "test-extension/src",
				Language: "javascript",
				Script:   "console.log('test');",
			},
			wantErr: true,
		},
		{
			name: "missing script",
			request: TestRequest{
				IP:       "http://192.168.1.100:8070",
				Root:     "test-extension/src",
				Language: "javascript",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginData_JSONMarshaling(t *testing.T) {
	pluginData := PluginData{
		ID:          "test-plugin",
		Name:        "Test Plugin",
		Author:      "Test Author",
		Version:     "1.0.0",
		Description: "A test plugin",
		Icon:        "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
		Source:      "test-source",
		Enabled:     true,
		Debug:       true,
		Data:        `{"main.js": "console.log('Hello World');"}`,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(pluginData)
	if err != nil {
		t.Fatalf("Failed to marshal PluginData: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled PluginData
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal PluginData: %v", err)
	}

	// Verify data integrity
	if unmarshaled.ID != pluginData.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, pluginData.ID)
	}
	if unmarshaled.Name != pluginData.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, pluginData.Name)
	}
	if unmarshaled.Enabled != pluginData.Enabled {
		t.Errorf("Enabled mismatch: got %v, want %v", unmarshaled.Enabled, pluginData.Enabled)
	}
}
