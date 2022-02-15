package cmd

import (
	"fmt"

	"github.com/axllent/adguard-home-bg/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the HTTP server",
	Long: `Run the HTTP server. This will allow you to access the web UI from 
a browser, and provide a downloadable URL for generated blocklists.`,
	Run: func(cmd *cobra.Command, args []string) {

		server.Listen = cmd.Flag("listen").Value.String()

		if err := server.Start(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("listen", "l", "0.0.0.0:8080", "Help message for toggle")
}
