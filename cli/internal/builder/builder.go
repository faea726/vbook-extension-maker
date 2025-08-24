package builder

import (
	"fmt"
	"path/filepath"

	"vbook-cli/internal/project"
)

type ExtensionBuilder struct {
	validator *project.ProjectValidator
	archiver  *Archiver
}

func NewExtensionBuilder() *ExtensionBuilder {
	return &ExtensionBuilder{
		validator: project.NewProjectValidator(),
		archiver:  NewArchiver(),
	}
}

// BuildExtension builds a Vbook extension into a distributable ZIP file
func (eb *ExtensionBuilder) BuildExtension(projectPath string) (string, error) {
	// Validate project structure
	if err := eb.ValidateBuildFiles(projectPath); err != nil {
		return "", fmt.Errorf("build validation failed: %w", err)
	}

	// Create output path
	outputPath := filepath.Join(projectPath, "plugin.zip")

	// Create archive
	if err := eb.CreateArchive(projectPath, outputPath); err != nil {
		return "", fmt.Errorf("failed to create archive: %w", err)
	}

	return outputPath, nil
}

// ValidateBuildFiles validates that all required files exist for building
func (eb *ExtensionBuilder) ValidateBuildFiles(projectPath string) error {
	// Use project validator to check required files
	if err := eb.validator.ValidateProject(projectPath); err != nil {
		return err
	}

	return nil
}

// CreateArchive creates the ZIP archive for the extension
func (eb *ExtensionBuilder) CreateArchive(projectPath, outputPath string) error {
	// Use the specialized Vbook extension ZIP creation method
	return eb.archiver.CreateVbookExtensionZip(projectPath, outputPath)
}
