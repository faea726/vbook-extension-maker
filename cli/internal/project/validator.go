package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ProjectValidator struct{}

func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{}
}

// ValidateProject validates that the given path is a valid Vbook extension project
func (pv *ProjectValidator) ValidateProject(projectPath string) error {
	// Check if path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Check if it's a directory
	fileInfo, err := os.Stat(projectPath)
	if err != nil {
		return fmt.Errorf("failed to check project path: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", projectPath)
	}

	// Check for required files
	if err := pv.checkRequiredFilesInternal(projectPath); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	return nil
}

// ValidateForTesting validates that the project has the minimum files needed for testing
func (pv *ProjectValidator) ValidateForTesting(projectPath string) error {
	// For testing, we only need plugin.json to exist
	pluginJsonPath := filepath.Join(projectPath, "plugin.json")
	if err := pv.validateFile(pluginJsonPath, "plugin.json"); err != nil {
		return err
	}

	// Validate plugin.json content
	if err := pv.validatePluginJson(pluginJsonPath); err != nil {
		return fmt.Errorf("invalid plugin.json: %w", err)
	}

	return nil
}

// ValidateForBuilding validates that the project has all files needed for building
func (pv *ProjectValidator) ValidateForBuilding(projectPath string) error {
	// For building, we need plugin.json, icon.png, and src directory
	requiredFiles := []string{"plugin.json", "icon.png"}
	if err := pv.CheckRequiredFiles(projectPath, requiredFiles); err != nil {
		return err
	}

	// Check for src directory
	srcPath := filepath.Join(projectPath, "src")
	if err := pv.validateDirectory(srcPath, "src"); err != nil {
		return err
	}

	// Check if src directory contains at least one .js file
	if err := pv.validateSrcDirectory(srcPath); err != nil {
		return err
	}

	return nil
}

// ValidateForInstall validates that the project has all files needed for installation
func (pv *ProjectValidator) ValidateForInstall(projectPath string) error {
	// For installation, we need plugin.json and icon.png
	requiredFiles := []string{"plugin.json", "icon.png"}
	if err := pv.CheckRequiredFiles(projectPath, requiredFiles); err != nil {
		return err
	}

	// Check for src directory and JavaScript files
	srcPath := filepath.Join(projectPath, "src")
	if err := pv.validateDirectory(srcPath, "src"); err != nil {
		return err
	}

	if err := pv.validateSrcDirectory(srcPath); err != nil {
		return err
	}

	return nil
}

// ResolveProjectPath resolves relative paths and defaults to current directory if empty
func (pv *ProjectValidator) ResolveProjectPath(path string) (string, error) {
	if path == "" {
		// Default to current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		return cwd, nil
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		path = filepath.Join(cwd, path)
	}

	// Clean the path
	path = filepath.Clean(path)

	// Verify the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("path does not exist: %s", path)
	}

	return path, nil
}

// CheckRequiredFiles checks if specific files exist in the project (updated signature)
func (pv *ProjectValidator) CheckRequiredFiles(projectPath string, files []string) error {
	for _, fileName := range files {
		filePath := filepath.Join(projectPath, fileName)
		if err := pv.validateFile(filePath, fileName); err != nil {
			return err
		}

		// Special validation for plugin.json
		if fileName == "plugin.json" {
			if err := pv.validatePluginJson(filePath); err != nil {
				return fmt.Errorf("invalid plugin.json: %w", err)
			}
		}
	}
	return nil
}

// checkRequiredFilesInternal checks if all required files exist in the project (internal method)
func (pv *ProjectValidator) checkRequiredFilesInternal(projectPath string) error {
	// Check for plugin.json
	pluginJsonPath := filepath.Join(projectPath, "plugin.json")
	if err := pv.validateFile(pluginJsonPath, "plugin.json"); err != nil {
		return err
	}

	// Validate plugin.json content
	if err := pv.validatePluginJson(pluginJsonPath); err != nil {
		return fmt.Errorf("invalid plugin.json: %w", err)
	}

	// Check for icon.png
	iconPath := filepath.Join(projectPath, "icon.png")
	if err := pv.validateFile(iconPath, "icon.png"); err != nil {
		return err
	}

	// Check for src directory
	srcPath := filepath.Join(projectPath, "src")
	if err := pv.validateDirectory(srcPath, "src"); err != nil {
		return err
	}

	// Check if src directory contains at least one .js file
	if err := pv.validateSrcDirectory(srcPath); err != nil {
		return err
	}

	return nil
}

// validateFile checks if a file exists and is readable
func (pv *ProjectValidator) validateFile(filePath, fileName string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("required file '%s' not found", fileName)
	}
	if err != nil {
		return fmt.Errorf("failed to check file '%s': %w", fileName, err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("'%s' should be a file, not a directory", fileName)
	}
	return nil
}

// validateDirectory checks if a directory exists
func (pv *ProjectValidator) validateDirectory(dirPath, dirName string) error {
	fileInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("required directory '%s' not found", dirName)
	}
	if err != nil {
		return fmt.Errorf("failed to check directory '%s': %w", dirName, err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("'%s' should be a directory, not a file", dirName)
	}
	return nil
}

// validatePluginJson validates the structure of plugin.json
func (pv *ProjectValidator) validatePluginJson(pluginJsonPath string) error {
	data, err := os.ReadFile(pluginJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin.json: %w", err)
	}

	var pluginConfig map[string]interface{}
	if err := json.Unmarshal(data, &pluginConfig); err != nil {
		return fmt.Errorf("plugin.json is not valid JSON: %w", err)
	}

	// Check for required sections
	if _, exists := pluginConfig["metadata"]; !exists {
		return fmt.Errorf("plugin.json missing required 'metadata' section")
	}

	if _, exists := pluginConfig["script"]; !exists {
		return fmt.Errorf("plugin.json missing required 'script' section")
	}

	return nil
}

// validateSrcDirectory checks if src directory contains JavaScript files
func (pv *ProjectValidator) validateSrcDirectory(srcPath string) error {
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read src directory: %w", err)
	}

	hasJsFiles := false
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".js" {
			hasJsFiles = true
			break
		}
	}

	if !hasJsFiles {
		return fmt.Errorf("src directory must contain at least one JavaScript (.js) file")
	}

	return nil
}

// CheckPluginJsonExists checks if plugin.json exists in the given path or its parent directories
func (pv *ProjectValidator) CheckPluginJsonExists(scriptPath string) (string, error) {
	// Start from the script's directory
	currentDir := filepath.Dir(scriptPath)

	// Look for plugin.json in current directory and parent directories
	for {
		pluginJsonPath := filepath.Join(currentDir, "plugin.json")
		if _, err := os.Stat(pluginJsonPath); err == nil {
			return currentDir, nil
		}

		// Move to parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root directory
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("plugin.json not found in script path or parent directories")
}
