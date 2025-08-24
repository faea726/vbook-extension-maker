package utils

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// VbookURL represents a parsed Vbook app URL
type VbookURL struct {
	Host string
	Port int
	Base string
}

type NetworkUtils struct{}

func NewNetworkUtils() *NetworkUtils {
	return &NetworkUtils{}
}

// GetLocalIP finds the local IP address that matches the given interface prefix
func (nu *NetworkUtils) GetLocalIP(interfacePrefix string, port int) (string, error) {
	// Handle emulator cases
	if strings.HasPrefix(interfacePrefix, "172.") || strings.HasPrefix(interfacePrefix, "10.") {
		interfacePrefix = "192.168." // Emulator fallback
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Helper function to check if interface should be skipped
	shouldSkipInterface := func(name string) bool {
		name = strings.ToLower(name)
		skipPatterns := []string{"vethernet", "virtual", "hyper-v", "wsl", "vmware", "virtualbox", "docker"}
		for _, pattern := range skipPatterns {
			if strings.Contains(name, pattern) {
				return true
			}
		}
		return false
	}

	// Helper function to extract IP from address
	extractIP := func(addr net.Addr) net.IP {
		switch v := addr.(type) {
		case *net.IPNet:
			return v.IP
		case *net.IPAddr:
			return v.IP
		}
		return nil
	}

	// First pass: look for matching prefix
	for _, iface := range interfaces {
		if shouldSkipInterface(iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			if strings.HasPrefix(ip.String(), interfacePrefix) {
				return fmt.Sprintf("http://%s:%d", ip.String(), port), nil
			}
		}
	}

	// Second pass: fallback to any valid IPv4
	for _, iface := range interfaces {
		if shouldSkipInterface(iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			return fmt.Sprintf("http://%s:%d", ip.String(), port), nil
		}
	}

	return "", fmt.Errorf("no suitable network interface found")
}

// GetLocalIPAuto automatically detects the best local IP address (adapted from TypeScript version)
// Prioritizes private LAN IPs and falls back to any non-internal IPv4
func (nu *NetworkUtils) GetLocalIPAuto(port int) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Helper function to check if IP is private LAN
	isPrivate := func(ip string) bool {
		if strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.168.") {
			return true
		}
		if strings.HasPrefix(ip, "172.") {
			parts := strings.Split(ip, ".")
			if len(parts) >= 2 {
				if second, err := strconv.Atoi(parts[1]); err == nil {
					return second >= 16 && second <= 31
				}
			}
		}
		return false
	}

	// Helper function to extract IP from address
	extractIP := func(addr net.Addr) net.IP {
		switch v := addr.(type) {
		case *net.IPNet:
			return v.IP
		case *net.IPAddr:
			return v.IP
		}
		return nil
	}

	// First pass: prefer private LAN IPs
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Check if it's IPv4 and private
			if ipv4 := ip.To4(); ipv4 != nil && isPrivate(ip.String()) {
				return fmt.Sprintf("http://%s:%d", ip.String(), port), nil
			}
		}
	}

	// Second pass: fallback to any non-internal IPv4
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			return fmt.Sprintf("http://%s:%d", ip.String(), port), nil
		}
	}

	return "", fmt.Errorf("no suitable network interface found")
}

// ValidateURL validates if the given string is a valid URL
func (nu *NetworkUtils) ValidateURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if scheme is present
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include a scheme (http:// or https://)")
	}

	// Check if scheme is supported
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s (only http and https are supported)", parsedURL.Scheme)
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

// ParseVbookURL parses and normalizes a Vbook app URL
func (nu *NetworkUtils) ParseVbookURL(urlStr string) (*VbookURL, error) {
	normalizedURL := nu.normalizeVbookURL(urlStr)

	// If normalization returned empty string, the URL format is invalid
	if normalizedURL == "" {
		return nil, fmt.Errorf("invalid URL format: %s (expected formats: http://IP:PORT, https://IP:PORT, http://IP, https://IP, or IP)", urlStr)
	}

	if err := nu.ValidateURL(normalizedURL); err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse normalized URL: %w", err)
	}

	host := parsedURL.Hostname()
	port := 8080 // default port

	if parsedURL.Port() != "" {
		port, err = strconv.Atoi(parsedURL.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %w", err)
		}
	}

	return &VbookURL{
		Host: host,
		Port: port,
		Base: normalizedURL,
	}, nil
}

// normalizeVbookURL normalizes various URL formats to a standard format
// This replicates the normalizeHost function from the original TypeScript version
func (nu *NetworkUtils) normalizeVbookURL(input string) string {
	input = strings.TrimSpace(input)

	// Define patterns matching the TypeScript version
	// Order matters - more specific patterns first
	patterns := []struct {
		name    string
		pattern *regexp.Regexp
	}{
		{"http_host_port", regexp.MustCompile(`^http://(?:\d{1,3}\.){3}\d{1,3}:\d+$`)},
		{"https_host_port", regexp.MustCompile(`^https://(?:\d{1,3}\.){3}\d{1,3}:\d+$`)},
		{"https_host", regexp.MustCompile(`^https://(?:\d{1,3}\.){3}\d{1,3}$`)},
		{"http_host", regexp.MustCompile(`^http://(?:\d{1,3}\.){3}\d{1,3}$`)},
		{"host_port", regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}:\d+$`)},
		{"host_only", regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)},
	}

	var matched string
	var caseType string

	// Find the first matching pattern
	for _, p := range patterns {
		if p.pattern.MatchString(input) {
			matched = input
			caseType = p.name
			break
		}
	}

	if matched == "" {
		return "" // No valid pattern found
	}

	// Apply normalization rules based on the matched pattern
	switch caseType {
	case "http_host_port", "https_host_port":
		return matched
	case "http_host", "https_host":
		return matched + ":8080"
	case "host_port":
		return "http://" + matched
	case "host_only":
		return "http://" + matched + ":8080"
	default:
		return ""
	}
}

// GetInterfacePrefix extracts the interface prefix from a Vbook URL
func (nu *NetworkUtils) GetInterfacePrefix(vbookURL *VbookURL) string {
	hostParts := strings.Split(vbookURL.Host, ".")
	if len(hostParts) >= 2 {
		return hostParts[0] + "." + hostParts[1] + "."
	}
	return "192.168." // fallback
}

// NormalizeVbookURL exposes the URL normalization functionality publicly
// This replicates the normalizeHost function from the original TypeScript version
func (nu *NetworkUtils) NormalizeVbookURL(input string) string {
	return nu.normalizeVbookURL(input)
}
