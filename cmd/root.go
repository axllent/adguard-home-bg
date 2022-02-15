package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/axllent/adguard-home-bg/parser"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adguard-home-bg <blocklist-url>",
	Short: "AdGuard Home Blocklist Generator",
	Long: `AdGuard Home Blocklist Generator.

A CLI and web utility to convert a public blocklist into a formatted
AdGuard Home blocklist, including adding additional rule modifiers.

Project & documentation: https://github.com/axllent/adguard-home-bg`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// config := parser.modifiersFromArgs

		domains, err := parser.Config.URLToBlocklist(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		currentTime := time.Now()
		date := currentTime.Format(time.UnixDate)
		sec := "! --------------------------------------------------\n"

		output := fmt.Sprintf("%s! Source:  %s\n! Domains: %d\n! Updated: %s\n%s", sec, args[0], len(domains), date, sec) + domains

		if cmd.Flag("output").Value.String() != "" {
			if err := os.WriteFile(cmd.Flag("output").Value.String(), []byte(output), 0644); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(domains)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Flags().SortFlags = false
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.Flags().StringVarP(&parser.Config.Ctag, "ctag", "t", "", "[value1,value2,...]")
	rootCmd.Flags().StringVarP(&parser.Config.Client, "client", "c", "", "[value1,value2,...]")
	rootCmd.Flags().StringVar(&parser.Config.DenyAllow, "denyallow", "", "[value1,value2,...]")
	rootCmd.Flags().StringVar(&parser.Config.DNSType, "dnstype", "", "[value1,value2,...]")
	rootCmd.Flags().StringVar(&parser.Config.DNSRewrite, "dnsrewrite", "", "[value1,value2,...]")
	rootCmd.Flags().BoolVarP(&parser.Config.Important, "important", "i", false, "increase priority over other rules without modifier")
	rootCmd.Flags().BoolVarP(&parser.Config.BadFilter, "badfilter", "b", false, "disable all rules")
	rootCmd.Flags().StringP("output", "o", "", "save output to file (default: stdout)")
}
