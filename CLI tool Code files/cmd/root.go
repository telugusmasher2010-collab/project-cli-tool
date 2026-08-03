package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/config"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/logger"
)

var (
	Version    = "1.0.0"
	Commit    = "none"
	Date      = "unknown"
	GoVersion = "unknown"
	verbose   bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "proj-init",
	Short: "Scaffold any project stack in 3 seconds",
	Long: `proj-init — A developer CLI tool that scaffolds a complete,
working project stack in seconds.

Choose from curated templates and get a ready-to-code folder
with src/, CI/CD, Dockerfile, README, .gitignore, and LICENSE.

Find more info at: https://github.com/telugusmasher2010-collab/project-cli-tool`,
	Example: `  # List the available templates
  proj-init list

  # Scaffold a new project interactively
  proj-init init

  # Scaffold into a specific directory
  proj-init init --output ./my-app

  # Check the current version
  proj-init version

  # Update proj-init to the latest release
  proj-init update`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path (default: ~/.proj-init/config.yml)")
	rootCmd.SetVersionTemplate(fmt.Sprintf("proj-init v%s\nCommit: %s\nBuilt: %s\nGo: %s\n", Version, Commit, Date, GoVersion))
	rootCmd.Version = Version
}

func initConfig() {
	if verbose {
		logger.SetLevel(logger.LevelDebug)
	}
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
	if cfg, err := config.Load(); err != nil {
		logger.Debug("config not loaded: %v", err)
	} else if cfg != nil && cfg.DefaultTemplate != "" {
		logger.Debug("default template from config: %s", cfg.DefaultTemplate)
	}
	logger.Debug("verbose mode enabled")
}
