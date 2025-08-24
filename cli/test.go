// test.go
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// -------------------- Test Script --------------------
func TestScript(filePath string, fileContent string) error {
	Log("\nvbook-ext: testScript")

	if !PluginJsonExist(filePath) {
		return fmt.Errorf("invalid workspace")
	}

	// Ask user for app IP
	appIP := GetValue("appIP", filePath)
	if appIP != nil {
		fmt.Printf("Current app IP: %s\n", appIP)
		choice := Prompt("Press Enter to reuse, or type a new IP")
		if choice != "" {
			appIP = NormalizeHost(choice)
			SetValue("appIP", appIP, filePath)
		}
	} else {
		ip := Prompt("Enter Vbook app IP (http://192.168.1.7:8080)")
		appIP = NormalizeHost(ip)
		SetValue("appIP", appIP, filePath)
	}

	appIPStr := fmt.Sprint(appIP)

	// Parse appIP into URL
	_url, err := url.Parse(appIPStr)
	if err != nil {
		Log("vbook-ext: Invalid App IP:", appIPStr)
		return err
	}

	// Server port = app port - 10
	port := 8080
	if _url.Port() != "" {
		fmt.Sscanf(_url.Port(), "%d", &port)
	}
	serverPort := port - 10

	// Params
	scriptName := filepath.Base(filePath) // e.g. "myScript.js"

	savedParams := GetValue(scriptName, filePath)
	var params string

	if savedParams != nil {
		fmt.Printf("Current params for %s: %s\n", scriptName, savedParams)
		choice := Prompt("Press Enter to reuse, or type new params (comma-separated)")
		if choice != "" {
			params = choice
			SetValue(scriptName, params, filePath)
		} else {
			params = fmt.Sprint(savedParams)
		}
	} else {
		params = Prompt("Enter params (comma-separated)")
		SetValue(scriptName, params, filePath)
	}
	Log("vbook-ext: Params:", params)

	extName := filepath.Base(filepath.Join(filePath, "../../"))

	data := map[string]any{
		"ip":       GetLocalIP(serverPort),
		"root":     fmt.Sprintf("%s/src", extName),
		"language": "javascript",
		"script":   fileContent,
	}

	if strings.Contains(params, ",") {
		parts := strings.Split(params, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		data["input"] = []any{parts}
	} else if strings.TrimSpace(params) != "" {
		data["input"] = []any{strings.TrimSpace(params)}
	} else {
		data["input"] = []any{}
	}

	// Build HTTP request manually
	jsonBytes, _ := json.Marshal(data)
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)

	request := strings.Join([]string{
		"GET /test HTTP/1.1",
		"Host: " + _url.Hostname(),
		"Connection: close",
		"data: " + encoded,
		"", "",
	}, "\r\n")

	// Run local server
	server := RunLocalServer(serverPort, filePath)

	// Connect to app via TCP
	conn, err := net.Dial("tcp", net.JoinHostPort(_url.Hostname(), fmt.Sprint(port)))
	if err != nil {
		server.Close()
		return fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	Log("vbook-ext: Connected to vbook:", _url.Hostname(), ":", port)

	// Send request
	_, err = conn.Write([]byte(request))
	if err != nil {
		server.Close()
		return err
	}

	// Read response
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		server.Close()
		return err
	}

	server.Close()
	Log("vbook-ext: Disconnected from server")

	rspStr := buf.String()

	parsed := ParseHttpResponse(rspStr)

	Log("\nResponse:")
	PrettyPrintJSON(parsed.Body, "")

	Log("\nvbook-ext: Done")
	return nil
}

// -------------------- Local IP Helper --------------------
func GetLocalIP(port int) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	isPrivate := func(ip string) bool {
		return strings.HasPrefix(ip, "10.") ||
			strings.HasPrefix(ip, "192.168.") ||
			(strings.HasPrefix(ip, "172.") && func() bool {
				parts := strings.Split(ip, ".")
				if len(parts) < 2 {
					return false
				}
				second := 0
				fmt.Sscanf(parts[1], "%d", &second)
				return second >= 16 && second <= 31
			}())
	}

	// Prefer private LAN
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip.To4() != nil && !ip.IsLoopback() && isPrivate(ip.String()) {
				localIp := fmt.Sprintf("http://%s:%d", ip.String(), port)
				Log("vbook-ext: Local IP (LAN):", localIp)
				return localIp
			}
		}
	}

	// Fallback: any non-loopback IPv4
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip.To4() != nil && !ip.IsLoopback() {
				localIp := fmt.Sprintf("http://%s:%d", ip.String(), port)
				Log("vbook-ext: Local IP (fallback):", localIp)
				return localIp
			}
		}
	}

	return ""
}

// -------------------- Get Opening File Content --------------------
func GetFileContent(filePath string) (name string, content string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}
	return filepath.Base(filePath), string(data), nil
}

// -------------------- HTTP Parser --------------------
type HttpResponse struct {
	Status string
	Header map[string]string
	Body   any
}
