package cmd

import (
	"fmt"
	"os"

	"github.com/axllent/adguard-home-bg/updater"
	"github.com/spf13/cobra"
)

var (
	// Version is the default application version, updated on release
	Version = "dev"

	// Repo on Github for updater
	Repo = "axllent/adguard-home-bg"

	// RepoBinaryName on Github for updater
	RepoBinaryName = "adguard-home-bg"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version & update information",
	Long:  `Displays the current version & update information.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		update, _ := cmd.Flags().GetBool("update")

		if update {
			return updateApp()
		}

		fmt.Printf("Version: %s\n", Version)
		latest, _, _, err := updater.GithubLatest(Repo, RepoBinaryName)
		if err == nil && updater.GreaterThan(latest, Version) {
			fmt.Printf(
				"Update available: %s\nRun `%s version -u` to update (requires read/write access to install directory).\n",
				latest,
				os.Args[0],
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().
		BoolP("update", "u", false, "update to latest version")
}

func updateApp() error {
	rel, err := updater.GithubUpdate(Repo, RepoBinaryName, Version)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %s to version %s\n", os.Args[0], rel)
	return nil
}
