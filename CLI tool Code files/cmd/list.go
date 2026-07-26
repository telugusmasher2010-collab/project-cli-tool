package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Run: func(cmd *cobra.Command, args []string) {
		available := templates.List()
		fmt.Println("Available templates:")
		fmt.Println()
		for _, t := range available {
			fmt.Printf("  %s\n", t.Name)
			fmt.Printf("    %s\n", t.Description)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
