package cmd

import (
	"fmt"
	"os"

	"github.com/axllent/ghru/v2"
	"github.com/spf13/cobra"
)

var (
	// Version is the default application version, updated on release
	Version = "dev"

	// Repo on Github for updater
	Repo = "axllent/adguard-home-bg"

	// RepoBinaryName on Github for updater
	RepoBinaryName = "adguard-home-bg"

	ghruConf = ghru.Config{
		Repo:           "axllent/adguard-home-bg",
		ArchiveName:    "adguard-home-bg_{{.OS}}_{{.Arch}}",
		BinaryName:     "adguard-home-bg",
		CurrentVersion: Version,
	}
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version & update information",
	Long:  `Displays the current version & update information.`,
	Run: func(cmd *cobra.Command, args []string) {

		update, _ := cmd.Flags().GetBool("update")

		if update {
			rel, err := ghruConf.SelfUpdate()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("Updated %s to version %s\n", os.Args[0], rel.Tag)
			os.Exit(0)
		}

		fmt.Printf("Version: %s\n", Version)

		release, err := ghruConf.Latest()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// The latest version is the same version
		if release.Tag == Version {
			os.Exit(0)
		}

		// A newer release is available
		fmt.Printf(
			"Update available: %s\nRun `%s version -u` to update (requires read/write access to install directory).\n",
			release.Tag,
			os.Args[0],
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().
		BoolP("update", "u", false, "update to latest version")
}
