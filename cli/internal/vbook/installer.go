package vbook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"vbook-cli/internal/project"
	"vbook-cli/internal/utils"
)

type VbookInstaller struct {
	validator *project.ProjectValidator
	client    *VbookClient
	netUtils  *utils.NetworkUtils
}

func NewVbookInstaller() *VbookInstaller {
	return &VbookInstaller{
		validator: project.NewProjectValidator(),
		client:    NewVbookClient(),
		netUtils:  utils.NewNetworkUtils(),
	}
}

// InstallExtension installs a Vbook extension to the specified Vbook app
func (vi *VbookInstaller) InstallExtension(projectPath, appURL string) error {
	// Validate project structure
	if err := vi.validator.ValidateProject(projectPath); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	// Prepare plugin data
	pluginData, err := vi.PreparePluginData(projectPath)
	if err != nil {
		return fmt.Errorf("failed to prepare plugin data: %w", err)
	}

	// Send installation request
	if err := vi.client.SendInstallRequest(appURL, pluginData); err != nil {
		return fmt.Errorf("installation request failed: %w", err)
	}

	return nil
}

// PreparePluginData prepares the plugin data structure for installation
func (vi *VbookInstaller) PreparePluginData(projectPath string) (*PluginData, error) {
	// Read plugin.json
	pluginConfig, err := vi.readPluginConfig(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin config: %w", err)
	}

	// Read and encode icon
	iconBase64, err := vi.readAndEncodeIcon(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read icon: %w", err)
	}

	// Read script files
	scriptData, err := vi.readScriptFiles(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read script files: %w", err)
	}

	// Convert script data to JSON string
	scriptDataJSON, err := json.Marshal(scriptData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal script data: %w", err)
	}

	// Create plugin data structure matching TypeScript implementation
	pluginData := &PluginData{
		ID:          "debug-" + pluginConfig.Metadata.Source,
		Name:        pluginConfig.Metadata.Name,
		Author:      pluginConfig.Metadata.Author,
		Version:     strconv.Itoa(pluginConfig.Metadata.Version),
		Description: pluginConfig.Metadata.Description,
		Source:      pluginConfig.Metadata.Source,
		Regexp:      pluginConfig.Metadata.Regexp,
		Locale:      pluginConfig.Metadata.Locale,
		Tag:         pluginConfig.Metadata.Tag,
		Type:        pluginConfig.Metadata.Type,
		Home:        pluginConfig.Script.Home,
		Genre:       pluginConfig.Script.Genre,
		Detail:      pluginConfig.Script.Detail,
		Search:      pluginConfig.Script.Search,
		Page:        pluginConfig.Script.Page,
		Toc:         pluginConfig.Script.Toc,
		Chap:        pluginConfig.Script.Chap,
		Icon:        iconBase64,
		Enabled:     true,
		Debug:       true,
		Data:        string(scriptDataJSON),
	}

	return pluginData, nil
}

// readPluginConfig reads and parses the plugin.json file
func (vi *VbookInstaller) readPluginConfig(projectPath string) (*PluginConfig, error) {
	pluginJsonPath := filepath.Join(projectPath, "plugin.json")

	data, err := os.ReadFile(pluginJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin.json: %w", err)
	}

	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse plugin.json: %w", err)
	}

	return &config, nil
}

// readAndEncodeIcon reads the icon file and encodes it as base64
func (vi *VbookInstaller) readAndEncodeIcon(projectPath string) (string, error) {
	iconPath := filepath.Join(projectPath, "icon.png")

	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		return "", fmt.Errorf("failed to read icon.png: %w", err)
	}

	// Encode as base64 with data URL prefix (matching TypeScript implementation)
	iconBase64 := "data:image/*;base64," + base64.StdEncoding.EncodeToString(iconData)

	return iconBase64, nil
}

// readScriptFiles reads all JavaScript files from the src directory
func (vi *VbookInstaller) readScriptFiles(projectPath string) (map[string]string, error) {
	srcPath := filepath.Join(projectPath, "src")

	scriptData := make(map[string]string)

	// Read all .js files in src directory
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read src directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .js files
		if filepath.Ext(entry.Name()) != ".js" {
			continue
		}

		scriptPath := filepath.Join(srcPath, entry.Name())
		scriptContent, err := os.ReadFile(scriptPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read script file %s: %w", entry.Name(), err)
		}

		scriptData[entry.Name()] = string(scriptContent)
	}

	if len(scriptData) == 0 {
		return nil, fmt.Errorf("no JavaScript files found in src directory")
	}

	return scriptData, nil
}
