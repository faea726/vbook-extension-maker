package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"vbook-cli/internal/utils"
	"vbook-cli/internal/vbook"

	"github.com/spf13/cobra"
)

var (
	appURL string
	params string
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test <script-path>",
	Short: "Test a Vbook extension script",
	Long: `Test a Vbook extension script against a running Vbook app.
The script path should be relative to the project root or an absolute path.

Examples:
  vbook-cli test src/detail.js --app-url http://192.168.1.100:8080
  vbook-cli test src/search.js --app-url http://192.168.1.100:8080 --params "query,page"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := args[0]

		// Convert to absolute path if relative
		if !filepath.IsAbs(scriptPath) {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			scriptPath = filepath.Join(cwd, scriptPath)
		}

		if appURL == "" {
			return fmt.Errorf("app URL is required (use --app-url flag)")
		}

		// Standardize the app URL
		netUtils := utils.NewNetworkUtils()
		normalizedURL := netUtils.NormalizeVbookURL(appURL)
		if normalizedURL == "" {
			return fmt.Errorf("invalid app URL format: %s (expected formats: http://IP:PORT, https://IP:PORT, http://IP, https://IP, or IP)", appURL)
		}
		appURL = normalizedURL

		tester := vbook.NewVbookTester()

		var paramList []string
		if params != "" {
			// Split params by comma
			paramList = []string{params}
		}

		logger.Info("Starting test with script: %s", scriptPath)
		logger.Info("App URL: %s", appURL)
		if len(paramList) > 0 {
			logger.Info("Parameters: %v", paramList)
		}

		result, err := tester.TestScript(scriptPath, appURL, paramList)
		if err != nil {
			return fmt.Errorf("failed to test script: %w", err)
		}

		logger.Info("Test completed successfully")

		// Display detailed response like the original TypeScript version
		fmt.Println("\nResponse:")

		if result.Status != "" {
			fmt.Printf("\nstatus: %s", result.Status)
		}

		if result.Output != "" {
			fmt.Printf("\noutput: %s", result.Output)
		}

		if result.Error != "" {
			fmt.Printf("\nerror: %s", result.Error)
		}

		// Show any additional details from the response
		if len(result.Details) > 0 {
			for key, value := range result.Details {
				if key == "status" || key == "output" || key == "error" {
					continue // Already shown above
				}

				if value != nil {
					switch v := value.(type) {
					case string:
						if v != "" {
							fmt.Printf("\n%s: %s", key, v)
						}
					case map[string]any, []any:
						fmt.Printf("\n%s: %+v", key, v)
					default:
						fmt.Printf("\n%s: %v", key, v)
					}
				}
			}
		}

		fmt.Println("\n\nDone")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVar(&appURL, "app-url", "", "Vbook app URL (required)")
	testCmd.Flags().StringVar(&params, "params", "", "Test parameters (comma-separated)")
	testCmd.MarkFlagRequired("app-url")
}
