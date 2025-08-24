package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type PlatformUtils struct{}

func NewPlatformUtils() *PlatformUtils {
	return &PlatformUtils{}
}

// GetPlatform returns the current platform information
func (pu *PlatformUtils) GetPlatform() string {
	return runtime.GOOS
}

// GetArchitecture returns the current architecture
func (pu *PlatformUtils) GetArchitecture() string {
	return runtime.GOARCH
}

// IsWindows checks if the current platform is Windows
func (pu *PlatformUtils) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux checks if the current platform is Linux
func (pu *PlatformUtils) IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMacOS checks if the current platform is macOS
func (pu *PlatformUtils) IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// GetPathSeparator returns the platform-appropriate path separator
func (pu *PlatformUtils) GetPathSeparator() string {
	if pu.IsWindows() {
		return "\\"
	}
	return "/"
}

// GetExecutableExtension returns the executable extension for the current platform
func (pu *PlatformUtils) GetExecutableExtension() string {
	if pu.IsWindows() {
		return ".exe"
	}
	return ""
}

// GetHomeDirectory returns the user's home directory
func (pu *PlatformUtils) GetHomeDirectory() (string, error) {
	if pu.IsWindows() {
		home := os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
		if home == "" {
			return "", fmt.Errorf("unable to determine home directory on Windows")
		}
		return home, nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("unable to determine home directory")
	}
	return home, nil
}

// GetConfigDirectory returns the platform-appropriate configuration directory
func (pu *PlatformUtils) GetConfigDirectory(appName string) (string, error) {
	home, err := pu.GetHomeDirectory()
	if err != nil {
		return "", err
	}

	if pu.IsWindows() {
		// Use AppData\Roaming on Windows
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return fmt.Sprintf("%s\\%s", appData, appName), nil
		}
		return fmt.Sprintf("%s\\AppData\\Roaming\\%s", home, appName), nil
	}

	if pu.IsMacOS() {
		// Use ~/Library/Application Support on macOS
		return fmt.Sprintf("%s/Library/Application Support/%s", home, appName), nil
	}

	// Use ~/.config on Linux and other Unix-like systems
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		return fmt.Sprintf("%s/%s", configHome, strings.ToLower(appName)), nil
	}
	return fmt.Sprintf("%s/.config/%s", home, strings.ToLower(appName)), nil
}

// GetCacheDirectory returns the platform-appropriate cache directory
func (pu *PlatformUtils) GetCacheDirectory(appName string) (string, error) {
	home, err := pu.GetHomeDirectory()
	if err != nil {
		return "", err
	}

	if pu.IsWindows() {
		// Use AppData\Local on Windows
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			return fmt.Sprintf("%s\\%s", localAppData, appName), nil
		}
		return fmt.Sprintf("%s\\AppData\\Local\\%s", home, appName), nil
	}

	if pu.IsMacOS() {
		// Use ~/Library/Caches on macOS
		return fmt.Sprintf("%s/Library/Caches/%s", home, appName), nil
	}

	// Use ~/.cache on Linux and other Unix-like systems
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome != "" {
		return fmt.Sprintf("%s/%s", cacheHome, strings.ToLower(appName)), nil
	}
	return fmt.Sprintf("%s/.cache/%s", home, strings.ToLower(appName)), nil
}

// NormalizeLineEndings converts line endings to the platform-appropriate format
func (pu *PlatformUtils) NormalizeLineEndings(text string) string {
	if pu.IsWindows() {
		// Convert to CRLF on Windows
		text = strings.ReplaceAll(text, "\r\n", "\n") // Normalize first
		text = strings.ReplaceAll(text, "\n", "\r\n")
		return text
	}
	// Convert to LF on Unix-like systems
	return strings.ReplaceAll(text, "\r\n", "\n")
}

// GetEnvironmentVariable gets an environment variable with platform-specific handling
func (pu *PlatformUtils) GetEnvironmentVariable(name string) string {
	value := os.Getenv(name)

	// On Windows, also try the uppercase version
	if value == "" && pu.IsWindows() {
		value = os.Getenv(strings.ToUpper(name))
	}

	return value
}

// SetEnvironmentVariable sets an environment variable
func (pu *PlatformUtils) SetEnvironmentVariable(name, value string) error {
	return os.Setenv(name, value)
}

// GetPlatformInfo returns detailed platform information
func (pu *PlatformUtils) GetPlatformInfo() map[string]string {
	return map[string]string{
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"go_version":    runtime.Version(),
		"num_cpu":       fmt.Sprintf("%d", runtime.NumCPU()),
		"num_goroutine": fmt.Sprintf("%d", runtime.NumGoroutine()),
	}
}
