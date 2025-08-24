package utils

import (
	"testing"
)

func TestNetworkUtils_ValidateURL(t *testing.T) {
	nu := NewNetworkUtils()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http URL", "http://192.168.1.100:8080", false},
		{"valid https URL", "https://192.168.1.100:8080", false},
		{"valid URL without port", "http://192.168.1.100", false},
		{"empty URL", "", true},
		{"URL without scheme", "192.168.1.100:8080", true},
		{"invalid scheme", "ftp://192.168.1.100:8080", true},
		{"URL without host", "http://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nu.ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkUtils_ParseVbookURL(t *testing.T) {
	nu := NewNetworkUtils()

	tests := []struct {
		name     string
		url      string
		wantHost string
		wantPort int
		wantErr  bool
	}{
		{"complete URL", "http://192.168.1.100:8080", "192.168.1.100", 8080, false},
		{"URL without port", "http://192.168.1.100", "192.168.1.100", 8080, false},
		{"IP only", "192.168.1.100", "192.168.1.100", 8080, false},
		{"HTTPS URL", "https://192.168.1.100:9000", "192.168.1.100", 9000, false},
		{"invalid URL", "invalid-url", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := nu.ParseVbookURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVbookURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Host != tt.wantHost {
					t.Errorf("ParseVbookURL() host = %v, want %v", result.Host, tt.wantHost)
				}
				if result.Port != tt.wantPort {
					t.Errorf("ParseVbookURL() port = %v, want %v", result.Port, tt.wantPort)
				}
			}
		})
	}
}

func TestNetworkUtils_normalizeVbookURL(t *testing.T) {
	nu := NewNetworkUtils()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"complete http URL", "http://192.168.1.100:8080", "http://192.168.1.100:8080"},
		{"complete https URL", "https://192.168.1.100:8080", "https://192.168.1.100:8080"},
		{"http without port", "http://192.168.1.100", "http://192.168.1.100:8080"},
		{"https without port", "https://192.168.1.100", "https://192.168.1.100:8080"},
		{"IP only", "192.168.1.100", "http://192.168.1.100:8080"},
		{"with whitespace", "  192.168.1.100  ", "http://192.168.1.100:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nu.normalizeVbookURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVbookURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNetworkUtils_GetInterfacePrefix(t *testing.T) {
	nu := NewNetworkUtils()

	tests := []struct {
		name     string
		vbookURL *VbookURL
		expected string
	}{
		{"standard IP", &VbookURL{Host: "192.168.1.100", Port: 8080}, "192.168."},
		{"different subnet", &VbookURL{Host: "10.0.1.100", Port: 8080}, "10.0."},
		{"single part IP", &VbookURL{Host: "localhost", Port: 8080}, "192.168."}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nu.GetInterfacePrefix(tt.vbookURL)
			if result != tt.expected {
				t.Errorf("GetInterfacePrefix() = %v, want %v", result, tt.expected)
			}
		})
	}
}
