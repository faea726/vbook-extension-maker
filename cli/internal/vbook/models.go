package vbook

import "fmt"

// TestRequest represents the data sent to Vbook app for testing
type TestRequest struct {
	IP       string     `json:"ip"`
	Root     string     `json:"root"`
	Language string     `json:"language"`
	Script   string     `json:"script"`
	Input    [][]string `json:"input"`
}

// TestResponse represents the response from Vbook app
type TestResponse struct {
	Status  string                 `json:"status"`
	Output  string                 `json:"output"`
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PluginData represents the data structure for plugin installation
type PluginData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Regexp      string `json:"regexp"`
	Locale      string `json:"locale"`
	Tag         string `json:"tag"`
	Type        string `json:"type"`
	Home        string `json:"home"`
	Genre       string `json:"genre"`
	Detail      string `json:"detail"`
	Search      string `json:"search"`
	Page        string `json:"page"`
	Toc         string `json:"toc"`
	Chap        string `json:"chap"`
	Icon        string `json:"icon"` // base64 encoded
	Enabled     bool   `json:"enabled"`
	Debug       bool   `json:"debug"`
	Data        string `json:"data"` // JSON string of script files
}

// PluginConfig represents the structure of plugin.json
type PluginConfig struct {
	Metadata struct {
		Name        string `json:"name"`
		Author      string `json:"author"`
		Version     int    `json:"version"`
		Source      string `json:"source"`
		Regexp      string `json:"regexp"`
		Description string `json:"description"`
		Locale      string `json:"locale"`
		Tag         string `json:"tag"`
		Type        string `json:"type"`
	} `json:"metadata"`
	Script struct {
		Home   string `json:"home"`
		Genre  string `json:"genre"`
		Detail string `json:"detail"`
		Search string `json:"search"`
		Page   string `json:"page"`
		Toc    string `json:"toc"`
		Chap   string `json:"chap"`
	} `json:"script"`
}

// Validate validates the PluginConfig structure
func (pc *PluginConfig) Validate() error {
	if pc.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}
	if pc.Metadata.Author == "" {
		return fmt.Errorf("metadata.author is required")
	}
	if pc.Metadata.Source == "" {
		return fmt.Errorf("metadata.source is required")
	}
	if pc.Metadata.Version <= 0 {
		return fmt.Errorf("metadata.version must be greater than 0")
	}
	return nil
}

// Validate validates the TestRequest structure
func (tr *TestRequest) Validate() error {
	if tr.IP == "" {
		return fmt.Errorf("IP is required")
	}
	if tr.Root == "" {
		return fmt.Errorf("root is required")
	}
	if tr.Language == "" {
		return fmt.Errorf("language is required")
	}
	if tr.Script == "" {
		return fmt.Errorf("script content is required")
	}
	return nil
}

// Validate validates the PluginData structure
func (pd *PluginData) Validate() error {
	if pd.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if pd.Name == "" {
		return fmt.Errorf("name is required")
	}
	if pd.Author == "" {
		return fmt.Errorf("author is required")
	}
	if pd.Version == "" {
		return fmt.Errorf("version is required")
	}
	if pd.Source == "" {
		return fmt.Errorf("source is required")
	}
	if pd.Icon == "" {
		return fmt.Errorf("icon is required")
	}
	if pd.Data == "" {
		return fmt.Errorf("data is required")
	}
	return nil
}
