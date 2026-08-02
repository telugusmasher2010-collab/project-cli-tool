package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/output"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/updater"
)

var updateCheckOnly bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update proj-init to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := updater.NewClient()
		info, err := updater.Check(Version, client)
		if err != nil {
			output.Cross(fmt.Sprintf("Failed to check for updates: %v", err))
			return err
		}

		if !info.Updateable {
			output.Check(fmt.Sprintf("Already up to date (v%s)", info.Current))
			return nil
		}

		output.Infof("New version available: v%s (current v%s)", info.Latest, info.Current)
		if updateCheckOnly {
			return nil
		}

		s := output.NewSpinner("Downloading update...")
		s.Start()
		err = updater.Apply(info, client)
		s.Stop()
		if err != nil {
			output.Cross(fmt.Sprintf("Update failed: %v", err))
			return err
		}
		output.Check(fmt.Sprintf("Updated to v%s. Restart proj-init to use the new version.", info.Latest))
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateCheckOnly, "check", false, "only check if an update is available")
	rootCmd.AddCommand(updateCmd)
}
