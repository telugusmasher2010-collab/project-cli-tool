package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	Commit    = "none"
	Date      = "unknown"
	GoVersion = "unknown"
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "proj-init",
	Short: "Scaffold any project stack in 3 seconds",
	Long: `proj-init — A developer CLI tool that scaffolds a complete,
working project stack in seconds.

Choose from curated templates and get a ready-to-code folder
with src/, CI/CD, Dockerfile, README, .gitignore, and LICENSE.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.SetVersionTemplate(fmt.Sprintf("proj-init v%s\nCommit: %s\nBuilt: %s\nGo: %s\n", Version, Commit, Date, GoVersion))
	rootCmd.Version = Version
}

func initConfig() {
}
