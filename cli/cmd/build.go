package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"vbook-cli/internal/builder"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [project-path]",
	Short: "Build a Vbook extension into a distributable package",
	Long: `Build a Vbook extension project into a plugin.zip file for distribution.
If no project path is specified, the current directory will be used.

Examples:
  vbook-cli build
  vbook-cli build ./my-extension
  vbook-cli build /path/to/my-extension`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectPath string
		if len(args) > 0 {
			projectPath = args[0]
		} else {
			var err error
			projectPath, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
		}

		// Convert to absolute path if relative
		if !filepath.IsAbs(projectPath) {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = filepath.Join(cwd, projectPath)
		}

		extensionBuilder := builder.NewExtensionBuilder()

		outputPath, err := extensionBuilder.BuildExtension(projectPath)
		if err != nil {
			return fmt.Errorf("failed to build extension: %w", err)
		}

		// Get file size
		fileInfo, err := os.Stat(outputPath)
		if err != nil {
			logger.Info("Extension built successfully: %s", outputPath)
		} else {
			logger.Info("Extension built successfully: %s (%d bytes)", outputPath, fileInfo.Size())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
