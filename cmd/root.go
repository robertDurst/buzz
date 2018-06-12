package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "buzz",
	Short: "A CLI tool to generate USD volume data for a given Stellar account.",
	Long:  `A CLI tool to generate USD volume data for a given Stellar account. This will be used by Lightyear partners and the Lightyear partnership team to measure volume for particular accounts of interest.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Probably not what you are looking for!
			
Consider running the --help command. 
					
To infinity and beyond!`)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
