package vbook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vbook-cli/internal/project"
	"vbook-cli/internal/utils"
)

type VbookTester struct {
	client    *VbookClient
	validator *project.ProjectValidator
	netUtils  *utils.NetworkUtils
}

func NewVbookTester() *VbookTester {
	return &VbookTester{
		client:    NewVbookClient(),
		validator: project.NewProjectValidator(),
		netUtils:  utils.NewNetworkUtils(),
	}
}

// TestScript tests a Vbook extension script against the Vbook app
func (vt *VbookTester) TestScript(scriptPath, appURL string, params []string) (*TestResponse, error) {
	// Find project root by looking for plugin.json
	projectRoot, err := vt.validator.CheckPluginJsonExists(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace: %w", err)
	}

	// Validate project structure
	if err := vt.validator.ValidateProject(projectRoot); err != nil {
		return nil, fmt.Errorf("project validation failed: %w", err)
	}

	// Parse and validate Vbook URL
	vbookURL, err := vt.netUtils.ParseVbookURL(appURL)
	if err != nil {
		return nil, fmt.Errorf("invalid app URL: %w", err)
	}

	// Calculate server port (app port - 10)
	serverPort := vbookURL.Port - 10

	// Get local IP for the interface
	interfacePrefix := vt.netUtils.GetInterfacePrefix(vbookURL)
	localIP, err := vt.netUtils.GetLocalIP(interfacePrefix, serverPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get local IP: %w", err)
	}

	// Start local server
	server := NewLocalServer(serverPort, projectRoot)
	go func() {
		if err := server.Start(); err != nil {
			// Server stopped or failed to start
		}
	}()
	defer server.Stop()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Read script content
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read script file: %w", err)
	}

	// Get extension name from project root
	extName := filepath.Base(projectRoot)

	// Prepare test request data
	testData := &TestRequest{
		IP:       localIP,
		Root:     extName + "/src",
		Language: "javascript",
		Script:   string(scriptContent),
		Input:    vt.prepareInputParams(params),
	}

	// Send test request
	response, err := vt.client.SendTestRequest(vbookURL.Base, testData)
	if err != nil {
		return nil, fmt.Errorf("test request failed: %w", err)
	}

	return response, nil
}

// StartLocalServer starts a local HTTP server for serving project files
func (vt *VbookTester) StartLocalServer(port int, projectPath string) (*LocalServer, error) {
	server := NewLocalServer(port, projectPath)

	go func() {
		if err := server.Start(); err != nil {
			// Log error but don't fail the operation
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	return server, nil
}

// prepareInputParams converts parameter strings to the expected format
func (vt *VbookTester) prepareInputParams(params []string) [][]string {
	if len(params) == 0 {
		return [][]string{{}}
	}

	var result [][]string
	for _, param := range params {
		if strings.Contains(param, ",") {
			// Split comma-separated values
			parts := strings.Split(param, ",")
			var trimmedParts []string
			for _, part := range parts {
				trimmedParts = append(trimmedParts, strings.TrimSpace(part))
			}
			result = append(result, trimmedParts)
		} else {
			result = append(result, []string{strings.TrimSpace(param)})
		}
	}

	return result
}
