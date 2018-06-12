package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "usd_volume_by_account",
	Short: "A CLI tool to generate USD volume data for a given Stellar account.",
	Long:  `A CLI tool to generate USD volume data for a given Stellar account. This will be used by Lightyear partners and the Lightyear partnership team to measure volume for particular accounts of interest.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Works")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
