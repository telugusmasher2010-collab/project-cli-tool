package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("proj-init v%s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", Date)
		fmt.Printf("Go: %s\n", GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
