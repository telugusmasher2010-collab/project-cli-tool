package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available templates:")
		fmt.Println()
		for _, t := range templates {
			fmt.Printf("  %s\n", t.Name)
			fmt.Printf("    %s\n", t.Description)
			fmt.Printf("    Stack: %s\n", t.Stack)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
