package vbook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"vbook-cli/internal/utils"
)

type VbookClient struct {
	timeout time.Duration
}

func NewVbookClient() *VbookClient {
	return &VbookClient{
		timeout: 30 * time.Second,
	}
}

// SendTestRequest sends a test request to the Vbook app via TCP connection
func (vc *VbookClient) SendTestRequest(appURL string, data *TestRequest) (*TestResponse, error) {
	// Standardize the URL first
	netUtils := utils.NewNetworkUtils()
	normalizedURL := netUtils.NormalizeVbookURL(appURL)
	if normalizedURL == "" {
		return nil, fmt.Errorf("invalid app URL format: %s", appURL)
	}
	appURL = normalizedURL

	// Parse URL to get host and port
	var host string
	var port string

	if strings.HasPrefix(appURL, "http://") {
		appURL = strings.TrimPrefix(appURL, "http://")
	} else if strings.HasPrefix(appURL, "https://") {
		appURL = strings.TrimPrefix(appURL, "https://")
	}

	parts := strings.Split(appURL, ":")
	if len(parts) == 2 {
		host = parts[0]
		port = parts[1]
	} else {
		host = parts[0]
		port = "8080"
	}

	// Establish TCP connection
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), vc.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Vbook app: %w", err)
	}
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(vc.timeout))

	// Prepare request data
	requestData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create HTTP request over TCP with base64 encoded data
	base64Data := base64.StdEncoding.EncodeToString(requestData)
	request := fmt.Sprintf("GET /test HTTP/1.1\r\nHost: %s\r\nConnection: close\r\ndata: %s\r\n\r\n",
		host, base64Data)

	// Send request
	_, err = conn.Write([]byte(request))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	response, err := vc.readHTTPResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return response, nil
}

// SendInstallRequest sends an installation request to the Vbook app
func (vc *VbookClient) SendInstallRequest(appURL string, data *PluginData) error {
	// Standardize the URL first
	netUtils := utils.NewNetworkUtils()
	normalizedURL := netUtils.NormalizeVbookURL(appURL)
	if normalizedURL == "" {
		return fmt.Errorf("invalid app URL format: %s", appURL)
	}
	appURL = normalizedURL

	// Prepare request data
	requestData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin data: %w", err)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: vc.timeout,
	}

	// Create request
	req, err := http.NewRequest("GET", appURL+"/install", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add data header with base64 encoding
	base64Data := base64.StdEncoding.EncodeToString(requestData)
	req.Header.Set("data", base64Data)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send install request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("install request failed with status: %d", resp.StatusCode)
	}

	// Read response body to check for actual success
	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to read install response: %w", err)
	}

	responseStr := string(body[:n])

	// Parse the JSON response to check status
	var installResponse map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &installResponse); err != nil {
		// If JSON parsing fails, check for error keywords in plain text
		if strings.Contains(strings.ToLower(responseStr), "error") ||
			strings.Contains(strings.ToLower(responseStr), "failed") {
			return fmt.Errorf("installation failed: %s", responseStr)
		}
		return nil // Assume success if no error keywords found
	}

	// Check the status field
	if status, ok := installResponse["status"]; ok {
		if statusCode, ok := status.(float64); ok {
			if statusCode != 0 {
				return fmt.Errorf("installation failed with status code: %.0f", statusCode)
			}
		}
	}

	return nil
}

// readHTTPResponse reads and parses HTTP response from TCP connection
func (vc *VbookClient) readHTTPResponse(conn net.Conn) (*TestResponse, error) {
	// Read all data from connection
	var chunks []byte
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			chunks = append(chunks, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	if len(chunks) == 0 {
		return &TestResponse{
			Status: "error",
			Error:  "No response received",
		}, nil
	}

	responseStr := string(chunks)

	// Split headers and body using double CRLF
	parts := strings.Split(responseStr, "\r\n\r\n")
	if len(parts) < 2 {
		// Try with single LF
		parts = strings.Split(responseStr, "\n\n")
	}

	if len(parts) < 2 {
		return &TestResponse{
			Status: "error",
			Error:  "Invalid HTTP response format",
			Output: responseStr,
		}, nil
	}

	bodyStr := strings.Join(parts[1:], "\n\n")

	// Clean up the body string - remove any null bytes or invalid UTF-8 sequences
	bodyStr = strings.ReplaceAll(bodyStr, "\x00", "")

	// Check if the body contains valid UTF-8
	if !isValidUTF8(bodyStr) {
		// If not valid UTF-8, try to clean it up
		bodyStr = cleanInvalidUTF8(bodyStr)
	}

	// Try to parse as JSON
	var response TestResponse
	var allFields map[string]any

	// First try to parse as a complete response object
	if err := json.Unmarshal([]byte(bodyStr), &response); err != nil {
		// If that fails, try to parse as a generic map to capture all fields
		if err := json.Unmarshal([]byte(bodyStr), &allFields); err != nil {
			// If JSON parsing still fails, check if it's a nested JSON string
			if cleanedBody := tryParseNestedJSON(bodyStr); cleanedBody != "" {
				if err := json.Unmarshal([]byte(cleanedBody), &allFields); err == nil {
					// Successfully parsed nested JSON
					response.Details = allFields
					response.Status = "success"

					if output, ok := allFields["output"]; ok {
						response.Output = fmt.Sprintf("%v", output)
					}
					if errorMsg, ok := allFields["error"]; ok {
						response.Error = fmt.Sprintf("%v", errorMsg)
					}
					if status, ok := allFields["status"]; ok {
						response.Status = fmt.Sprintf("%v", status)
					}

					return &response, nil
				}
			}

			// If all JSON parsing fails, treat as plain text response
			return &TestResponse{
				Status: "success",
				Output: bodyStr,
			}, nil
		}

		// Convert map to TestResponse and capture all fields
		response.Details = allFields
		response.Status = "success"

		if output, ok := allFields["output"]; ok {
			response.Output = fmt.Sprintf("%v", output)
		}
		if errorMsg, ok := allFields["error"]; ok {
			response.Error = fmt.Sprintf("%v", errorMsg)
		}
		if status, ok := allFields["status"]; ok {
			response.Status = fmt.Sprintf("%v", status)
		}
	} else {
		// If direct parsing worked, also try to get all fields for Details
		json.Unmarshal([]byte(bodyStr), &allFields)
		response.Details = allFields
	}

	return &response, nil
}

// isValidUTF8 checks if a string contains valid UTF-8
func isValidUTF8(s string) bool {
	return len(s) == len([]rune(s))
}

// cleanInvalidUTF8 removes invalid UTF-8 sequences from a string
func cleanInvalidUTF8(s string) string {
	// Convert to runes and back to string to remove invalid sequences
	runes := []rune(s)
	return string(runes)
}

// tryParseNestedJSON attempts to extract and parse nested JSON strings
func tryParseNestedJSON(s string) string {
	// Look for JSON-like patterns in the string
	s = strings.TrimSpace(s)

	// Try to find JSON object boundaries
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return s
	}

	// Look for escaped JSON strings
	if strings.Contains(s, "\\\"") {
		// Try to unescape
		unescaped := strings.ReplaceAll(s, "\\\"", "\"")
		unescaped = strings.ReplaceAll(unescaped, "\\\\", "\\")
		if strings.HasPrefix(unescaped, "{") && strings.HasSuffix(unescaped, "}") {
			return unescaped
		}
	}

	return ""
}
