// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"

	"github.com/robertdurst/buzz/payments"
	"github.com/spf13/cobra"
	"github.com/stellar/go/strkey"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [stellar address] [currencylayer api key]",
	Short: "Query USD volume for a specific Stellar Address",
	Long: `Queries for all payments to and from a specific account. Then
		   calculates total USD volume per day via Currencylayer API
		   and data scraped from Coinmarketcap. This data may be returned
		   in numerous ways, including raw and csv.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires stellar address and currencylayer api key")
		}

		_, err := strkey.Decode(strkey.VersionByteAccountID, args[0])
		if err != nil {
			return errors.New("invalid stellar address")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var agg payments.Aggregate

		// capture arguments
		account, apikey := args[0], args[1]

		// capture flags
		aggregate := cmd.Flag("aggregate").Value.String()
		output := cmd.Flag("output").Value.String()
		filename := cmd.Flag("filename").Value.String()

		// get all payments for account from horizon
		p := payments.PaymentsForAccount(account)

		// utilize coinmarketcap and currencylayer to
		// gather and fill in the usd payments volumes
		rawdata := payments.FillInVolumePerPayment(p, apikey)

		switch aggregate {
		case "day":
			agg = payments.ByDate{Data: payments.OrderDataByDate(rawdata)}
			break
		case "month":
			agg = payments.ByMonth{Data: payments.OrderDataByMonth(rawdata)}
			break
		default:
			agg = payments.Raw{Data: payments.OrderDataByDate(rawdata)}
		}

		switch output {
		case "csv":
			payments.CreateCSV(agg, filename)
			break
		case "markdown":
			payments.CreateMarkdown(agg)
			break
		default:
			payments.CreateTable(agg)
			break
		}
	},
}

func init() {
	rootCmd.PersistentFlags().String("aggregate", "none", "aggregate data by time interval (accepted inputs: none, day, month)")
	rootCmd.PersistentFlags().String("output", "terminal", "output type (accepted inputs: terminal, csv, markdown)")
	rootCmd.PersistentFlags().String("filename", "results", "csv output file name")
	rootCmd.AddCommand(queryCmd)
}
