// utils.go
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// -------------------- Plugin.json check --------------------
func PluginJsonExist(scriptPath string) bool {
	rootPath := filepath.Join(scriptPath, "../../")
	pluginJsonPath := filepath.Join(rootPath, "plugin.json")

	info, err := os.Stat(pluginJsonPath)
	if err != nil || info.IsDir() {
		return false
	}
	return filepath.Base(filepath.Dir(scriptPath)) == "src"
}

// -------------------- Temporary Data --------------------
var tempData map[string]any = make(map[string]any)

func SetValue(key string, value any, scriptPath string) error {
	if !PluginJsonExist(scriptPath) {
		return fmt.Errorf("invalid workspace")
	}

	rootPath := filepath.Join(scriptPath, "../../")
	tempPath := filepath.Join(rootPath, "test.json")

	tempData[key] = value
	data, err := json.MarshalIndent(tempData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tempPath, data, 0o644)
}

func GetValue(key string, scriptPath string) any {
	if !PluginJsonExist(scriptPath) {
		return nil
	}

	rootPath := filepath.Join(scriptPath, "../../")
	tempPath := filepath.Join(rootPath, "test.json")

	data, err := os.ReadFile(tempPath)
	if err != nil {
		return nil
	}

	json.Unmarshal(data, &tempData)
	return tempData[key]
}

// -------------------- Normalize host --------------------
func NormalizeHost(input string) string {
	patterns := map[string]*regexp.Regexp{
		"http_host_port":  regexp.MustCompile(`\bhttp://(?:\d{1,3}\.){3}\d{1,3}:\d+\b`),
		"https_host_port": regexp.MustCompile(`\bhttps://(?:\d{1,3}\.){3}\d{1,3}:\d+\b`),
		"https_host":      regexp.MustCompile(`\bhttps://(?:\d{1,3}\.){3}\d{1,3}\b`),
		"http_host":       regexp.MustCompile(`\bhttp://(?:\d{1,3}\.){3}\d{1,3}\b`),
		"host_only":       regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
	}

	for caseType, re := range patterns {
		if re.MatchString(input) {
			host := re.FindString(input)
			switch caseType {
			case "http_host_port", "https_host_port":
				return host
			case "http_host", "https_host":
				return host + ":8080"
			case "host_only":
				return "http://" + host + ":8080"
			}
		}
	}
	Log("Invalid IP")
	os.Exit(1)
	return ""
}

// -------------------- Local File Server --------------------
func RunLocalServer(port int, scriptPath string) *http.Server {
	srcPath := filepath.Join(scriptPath, "../../../")

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			file := query.Get("file")
			root := query.Get("root")

			if file == "" || root == "" {
				http.Error(w, "Missing required query parameters: file and root", http.StatusBadRequest)
				return
			}

			filePath := filepath.Join(srcPath, root, file)
			data, err := os.ReadFile(filePath)
			if err != nil {
				http.Error(w, "Error reading the file", http.StatusInternalServerError)
				return
			}

			encoded := base64.StdEncoding.EncodeToString(data)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", fmt.Sprint(len(encoded)))
			w.Write([]byte(encoded))
		}),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	fmt.Println("vbook-ext: Server listening on port", port)
	return server
}

// -------------------- HTTP Response Parser --------------------
type ParsedResponse struct {
	StatusLine string            `json:"statusLine"`
	Headers    map[string]string `json:"headers"`
	Body       any               `json:"body"`
}

func ParseHttpResponse(raw string) HttpResponse {
	resp := HttpResponse{
		Header: make(map[string]string),
		Body:   "",
	}

	// Normalize line endings first
	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	parts := strings.SplitN(raw, "\n\n", 2)
	head := parts[0]
	if len(parts) > 1 {
		body := parts[1]
		var js any
		if err := json.Unmarshal([]byte(body), &js); err == nil {
			resp.Body = DeepParseJSON(body)
		} else {
			resp.Body = body
		}
	}

	lines := strings.Split(head, "\n")
	if len(lines) > 0 {
		resp.Status = lines[0]
	}
	for _, line := range lines[1:] {
		if kv := strings.SplitN(line, ":", 2); len(kv) == 2 {
			resp.Header[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return resp
}

func DeepParseJSON(input string) any {
	var parsed any
	if err := json.Unmarshal([]byte(input), &parsed); err == nil {
		switch v := parsed.(type) {
		case map[string]any:
			for k, val := range v {
				if s, ok := val.(string); ok {
					v[k] = DeepParseJSON(s)
				}
			}
			return v
		case []any:
			for i, val := range v {
				if s, ok := val.(string); ok {
					v[i] = DeepParseJSON(s)
				}
			}
			return v
		default:
			return v
		}
	}
	return input
}

// -------------------- Pretty Print (recursive) --------------------
func PrettyPrintJSON(v any, indent string) {
	switch val := v.(type) {
	case map[string]any:
		for k, sub := range val {
			fmt.Printf("%s%s:\n", indent, k)
			PrettyPrintJSON(sub, indent+"  ")
		}
	case []any:
		for i, sub := range val {
			fmt.Printf("%s[%d]:\n", indent, i)
			PrettyPrintJSON(sub, indent+"  ")
		}
	default:
		fmt.Printf("%s%v\n", indent, val)
	}
}

// -------------------- Simple Logger --------------------
func Log(args ...any) {
	fmt.Println(args...)
}

// -------------------- Helper: Get input from user --------------------
func Prompt(prompt string) string {
	fmt.Print(prompt + ": ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
