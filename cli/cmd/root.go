package cmd

import (
	"vbook-cli/internal/utils"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	version = "1.0.0"
	logger  *utils.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vbook-cli",
	Short: "A CLI tool for managing Vbook extensions",
	Long: `Vbook CLI is a cross-platform command-line tool for testing, 
building, and installing Vbook extensions. It provides the core functionality 
for managing existing Vbook extensions as a standalone CLI application.

Examples:
  vbook-cli test src/detail.js --app-url http://192.168.1.100:8080
  vbook-cli build ./my-extension
  vbook-cli install ./my-extension --app-url http://192.168.1.100:8080`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger = utils.NewLogger(verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	// Set custom help template
	rootCmd.SetHelpTemplate(`{{.Long}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Flags:
{{.LocalFlags.FlagUsages}}{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages}}{{end}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
}
