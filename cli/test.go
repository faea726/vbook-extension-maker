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
	"strconv"
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

	// get first 3 prefix from appIP
	var interfaceIP string
	if _url.Hostname() != "" {
		// test if hostname is ipv4
		if net.ParseIP(_url.Hostname()).To4() != nil {
			interfaceIP = ""
		}
		interfaceIP = strings.Split(_url.Hostname(), ".")[0] + "." + strings.Split(_url.Hostname(), ".")[1] + "." + strings.Split(_url.Hostname(), ".")[2]
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
		"ip":       GetLocalIP(serverPort, interfaceIP),
		"root":     fmt.Sprintf("%s/src", extName),
		"language": "javascript",
		"script":   fileContent,
	}

	input := []any{}

	for p := range strings.SplitSeq(params, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			input = append(input, p)
		}
	}

	data["input"] = input

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
func GetLocalIP(port int, interfaceIP ...string) string {
	var prefix string
	if len(interfaceIP) > 0 {
		prefix = strings.TrimSpace(interfaceIP[0])
	}

	isPrivate := func(ip string) bool {
		if strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.168.") {
			return true
		}
		if strings.HasPrefix(ip, "172.") {
			parts := strings.Split(ip, ".")
			if len(parts) < 2 {
				return false
			}
			second, _ := strconv.Atoi(parts[1])
			return second >= 16 && second <= 31
		}
		return false
	}

	matchPrefix := func(ip string) bool {
		if prefix == "" {
			return false
		}
		ipParts := strings.Split(ip, ".")
		prefParts := strings.Split(prefix, ".")
		for i, p := range prefParts {
			if p == "" {
				continue
			}
			if ipParts[i] != p {
				return false
			}
		}
		return true
	}

	type candidate struct {
		ip    string
		score int
	}

	var candidates []candidate

	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip == nil || ip.To4() == nil || ip.IsLoopback() {
				continue
			}

			ipStr := ip.String()
			score := 0
			if isPrivate(ipStr) {
				score += 10
			}
			if matchPrefix(ipStr) {
				score += 5
			}
			candidates = append(candidates, candidate{ip: ipStr, score: score})
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	// pick highest score
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.score > best.score {
			best = c
		}
	}

	localIp := fmt.Sprintf("http://%s:%d", best.ip, port)
	fmt.Println("vbook-ext: Local IP (chosen):", localIp)
	return localIp
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
