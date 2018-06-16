// Post Process contains methods for sorting and outputing data.
package payments

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func CreateCSV(data [][]TruncatedPayment, fileName string, aggregate string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)

	switch aggregate {
	case "day":
		// Add a header row to csv
		strV = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
		for _, v := range data {
			volume := 0.0
			date := ""
			for _, p := range v {
				volume += p.VolumeUSD
				date = p.FormattedDate
			}
			strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
		}
		break
	case "month":
		// Add a header row to csv
		strV = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
		for _, v := range data {
			volume := 0.0
			date := ""
			for _, p := range v {
				volume += p.VolumeUSD

				s := strings.Split(p.FormattedDate, "-")
				date = fmt.Sprintf("%s-%s", s[0], s[1])
			}
			strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
		}
		break
	default:
		// Add a header row to csv
		strV = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
		for _, v := range data {
			for _, p := range v {
				strV = append(strV, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
			}
		}
	}

	writer.WriteAll(strV)
}

func OutputData(data [][]TruncatedPayment, aggregate string) {
	table := tablewriter.NewWriter(os.Stdout)
	switch aggregate {
	case "day":
		table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
		for _, v := range data {
			volume := 0.0
			date := ""
			for _, p := range v {
				volume += p.VolumeUSD
				date = p.FormattedDate
			}
			table.Append([]string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
		}
		break
	case "month":
		// Add a header row to csv
		table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
		for _, v := range data {
			volume := 0.0
			date := ""
			for _, p := range v {
				volume += p.VolumeUSD

				s := strings.Split(p.FormattedDate, "-")
				date = fmt.Sprintf("%s-%s", s[0], s[1])
			}
			table.Append([]string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
		}
		break
	default:
		// Add a header row to csv
		table.SetHeader([]string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
		for _, v := range data {
			for _, p := range v {
				table.Append([]string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
			}
		}
	}

	table.Render()
}
