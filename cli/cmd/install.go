package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"vbook-cli/internal/utils"
	"vbook-cli/internal/vbook"

	"github.com/spf13/cobra"
)

var installAppURL string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [project-path]",
	Short: "Install a Vbook extension to the Vbook app",
	Long: `Install a Vbook extension directly to a running Vbook app.
If no project path is specified, the current directory will be used.

Examples:
  vbook-cli install --app-url http://192.168.1.100:8080
  vbook-cli install ./my-extension --app-url http://192.168.1.100:8080
  vbook-cli install /path/to/my-extension --app-url http://192.168.1.100:8080`,
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

		if installAppURL == "" {
			return fmt.Errorf("app URL is required (use --app-url flag)")
		}

		// Standardize the app URL
		netUtils := utils.NewNetworkUtils()
		normalizedURL := netUtils.NormalizeVbookURL(installAppURL)
		if normalizedURL == "" {
			return fmt.Errorf("invalid app URL format: %s (expected formats: http://IP:PORT, https://IP:PORT, http://IP, https://IP, or IP)", installAppURL)
		}
		installAppURL = normalizedURL

		installer := vbook.NewVbookInstaller()

		if err := installer.InstallExtension(projectPath, installAppURL); err != nil {
			return fmt.Errorf("failed to install extension: %w", err)
		}

		logger.Info("Extension installed successfully to %s", installAppURL)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installAppURL, "app-url", "", "Vbook app URL (required)")
	installCmd.MarkFlagRequired("app-url")
}
