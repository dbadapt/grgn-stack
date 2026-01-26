package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grgn",
	Short: "GRGN Stack CLI - Development tools for the GRGN stack",
	Long: `GRGN CLI provides development tools for managing the GRGN stack:
  - Migration management (up, down, status)
  - Code generation orchestration
  - App scaffolding (future)
  - Architecture validation (future)`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	// Note: seedCmd is registered in seed.go init()
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grgn v0.1.0")
	},
}
