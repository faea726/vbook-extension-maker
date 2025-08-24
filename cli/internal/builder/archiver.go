package builder

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Archiver struct{}

func NewArchiver() *Archiver {
	return &Archiver{}
}

// CreateZip creates a ZIP archive from the source directory
func (a *Archiver) CreateZip(sourceDir, outputPath string, includePatterns []string) error {
	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create ZIP writer with maximum compression
	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	// Walk through source directory
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Check if file matches include patterns
		if !a.matchesPatterns(relPath, includePatterns) {
			return nil
		}

		// Add file to ZIP
		return a.AddFileToZip(zipWriter, path, relPath)
	})
	if err != nil {
		return fmt.Errorf("failed to create ZIP archive: %w", err)
	}

	return nil
}

// AddFileToZip adds a single file to the ZIP archive
func (a *Archiver) AddFileToZip(writer *zip.Writer, filePath, zipPath string) error {
	// Open source file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create ZIP file header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create ZIP header: %w", err)
	}

	// Set the name and compression method
	header.Name = filepath.ToSlash(zipPath) // Use forward slashes for ZIP paths
	header.Method = zip.Deflate             // Use maximum compression

	// Create writer for this file
	zipFile, err := writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file entry: %w", err)
	}

	// Copy file content to ZIP
	_, err = io.Copy(zipFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// matchesPatterns checks if a file path matches any of the include patterns
func (a *Archiver) matchesPatterns(filePath string, patterns []string) bool {
	if len(patterns) == 0 {
		return true // Include all files if no patterns specified
	}

	for _, pattern := range patterns {
		// Direct pattern match
		matched, err := filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}

		// Check if pattern ends with /* (directory pattern)
		if strings.HasSuffix(pattern, "/*") {
			dirPattern := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(filePath, dirPattern+"/") {
				return true
			}
		}

		// Check if the file is in a directory that matches the pattern
		dir := filepath.Dir(filePath)
		if strings.HasPrefix(filePath, pattern) || strings.HasPrefix(dir, pattern) {
			return true
		}
	}

	return false
}

// CreateVbookExtensionZip creates a ZIP archive specifically for Vbook extensions
func (a *Archiver) CreateVbookExtensionZip(projectDir, outputPath string) error {
	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create ZIP writer with maximum compression
	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	// Add plugin.json
	pluginJsonPath := filepath.Join(projectDir, "plugin.json")
	if err := a.AddFileToZip(zipWriter, pluginJsonPath, "plugin.json"); err != nil {
		return fmt.Errorf("failed to add plugin.json: %w", err)
	}

	// Add icon.png
	iconPath := filepath.Join(projectDir, "icon.png")
	if err := a.AddFileToZip(zipWriter, iconPath, "icon.png"); err != nil {
		return fmt.Errorf("failed to add icon.png: %w", err)
	}

	// Add src directory
	srcDir := filepath.Join(projectDir, "src")
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from project directory
		relPath, err := filepath.Rel(projectDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Add file to ZIP
		return a.AddFileToZip(zipWriter, path, relPath)
	})
	if err != nil {
		return fmt.Errorf("failed to add src directory: %w", err)
	}

	return nil
}
