package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// -------------------- Install Extension --------------------
func InstallExtension(scriptPath string) error {
	Log("\nvbook-ext: installExtension")

	// Check plugin.json exists
	if !PluginJsonExist(scriptPath) {
		return fmt.Errorf("invalid workspace")
	}
	rootPath := filepath.Join(scriptPath, "../../")

	// Prepare plugin data
	data, err := PreparePluginData(rootPath)
	if err != nil {
		return err
	}

	// Ask user for IP (simulate vscode.showInputBox)
	// Ask user for app IP
	appIP := GetValue("appIP", scriptPath)
	if appIP != nil {
		fmt.Printf("Current app IP: %s\n", appIP)
		choice := Prompt("Press Enter to reuse, or type a new IP")
		if choice != "" {
			appIP = NormalizeHost(choice)
			SetValue("appIP", appIP, scriptPath)
		}
	} else {
		ip := Prompt("Enter Vbook app IP (http://192.168.1.7:8080)")
		appIP = NormalizeHost(ip)
		SetValue("appIP", appIP, scriptPath)
	}
	appIPStr := fmt.Sprint(appIP)

	// Send install request
	Log("vbook-ext: Installing to:", appIPStr)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	headerData := base64.StdEncoding.EncodeToString(jsonBytes)

	req, err := http.NewRequest("GET", appIPStr+"/install", nil)
	if err != nil {
		return err
	}
	req.Header.Set("data", headerData)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		Log("vbook-ext: done installation process.")
		return err
	}

	Log("vbook-ext: Installation request sent.")
	return nil
}

// -------------------- Prepare Plugin Data --------------------
func PreparePluginData(pluginDir string) (map[string]interface{}, error) {
	pluginDetailPath := filepath.Join(pluginDir, "plugin.json")
	iconPath := filepath.Join(pluginDir, "icon.png")

	if _, err := os.Stat(pluginDetailPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid plugin: missing plugin.json")
	}
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid plugin: missing icon.png")
	}

	// Read plugin.json
	raw, err := os.ReadFile(pluginDetailPath)
	if err != nil {
		return nil, err
	}
	var pluginDetail map[string]interface{}
	if err := json.Unmarshal(raw, &pluginDetail); err != nil {
		return nil, err
	}

	metadata, _ := pluginDetail["metadata"].(map[string]interface{})
	script, _ := pluginDetail["script"].(map[string]interface{})

	// Merge metadata + script
	data := make(map[string]interface{})
	for k, v := range metadata {
		data[k] = v
	}
	for k, v := range script {
		data[k] = v
	}

	// Assign extra fields
	if source, ok := metadata["source"].(string); ok {
		data["id"] = "debug-" + source
	} else {
		data["id"] = "debug-unknown"
	}

	iconBuffer, err := os.ReadFile(iconPath)
	if err != nil {
		return nil, err
	}
	data["icon"] = "data:image/*;base64," + base64.StdEncoding.EncodeToString(iconBuffer)
	data["enabled"] = true
	data["debug"] = true
	data["data"] = map[string]string{}

	// Collect .js files in src/
	srcDir := filepath.Join(pluginDir, "src")
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	jsFiles := make(map[string]string)
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".js" {
			scriptPath := filepath.Join(srcDir, f.Name())
			content, err := os.ReadFile(scriptPath)
			if err == nil {
				jsFiles[f.Name()] = string(content)
			}
		}
	}
	// Store scripts as JSON string (like TS version)
	jsFilesJSON, _ := json.Marshal(jsFiles)
	data["data"] = string(jsFilesJSON)

	return data, nil
}
