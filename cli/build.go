package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// -------------------- Build Extension (zip) --------------------
func BuildExtension(scriptPath string) error {
	Log("\nvbook-ext: buildExtension")

	rootPath := filepath.Join(scriptPath, "../../")

	pluginJSONPath := filepath.Join(rootPath, "plugin.json")
	iconPath := filepath.Join(rootPath, "icon.png")
	srcPath := filepath.Join(rootPath, "src")
	outputPath := filepath.Join(rootPath, "plugin.zip")

	// Validate existence
	if !fileExists(pluginJSONPath) || !fileExists(iconPath) || !dirExists(srcPath) {
		return fmt.Errorf("files not found")
	}

	// Create output zip file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Add plugin.json
	if err := addFileToZip(zipWriter, pluginJSONPath, "plugin.json"); err != nil {
		return err
	}

	// Add icon.png
	if err := addFileToZip(zipWriter, iconPath, "icon.png"); err != nil {
		return err
	}

	// Add src/ directory recursively
	if err := addDirToZip(zipWriter, srcPath, "src"); err != nil {
		return err
	}

	Log(fmt.Sprintf("vbook-ext: plugin.zip created: %s", outputPath))
	return nil
}

// -------------------- Helpers --------------------
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func addFileToZip(zipWriter *zip.Writer, filePath, zipName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := zipWriter.Create(zipName)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

func addDirToZip(zipWriter *zip.Writer, dirPath, baseInZip string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Relative path inside the zip
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		zipName := filepath.Join(baseInZip, relPath)

		return addFileToZip(zipWriter, path, zipName)
	})
}
