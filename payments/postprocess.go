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

type Aggregate interface {
	CreateCSV()
	OutputData()
}

type Raw struct {
	Data     [][]TruncatedPayment
	FileName string
}

func (r Raw) CreateCSV() {
	file, err := os.Create(r.FileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)

	// Add a header row to csv
	strV = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
	for _, v := range r.Data {
		for _, p := range v {
			strV = append(strV, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
		}
	}
	writer.WriteAll(strV)
}

func (r Raw) OutputData() {
	table := tablewriter.NewWriter(os.Stdout)

	// Add a header row to csv
	table.SetHeader([]string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
	for _, v := range r.Data {
		for _, p := range v {
			table.Append([]string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
		}
	}

	table.Render()
}

type ByDate struct {
	Data     [][]TruncatedPayment
	FileName string
}

func (b ByDate) CreateCSV() {
	file, err := os.Create(b.FileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)

	// Add a header row to csv
	strV = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD
			date = p.FormattedDate
		}
		strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}
	writer.WriteAll(strV)
}

func (b ByDate) OutputData() {
	table := tablewriter.NewWriter(os.Stdout)

	// Add a header row to csv
	table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD
			date = p.FormattedDate
		}
		table.Append([]string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}

	table.Render()
}

type ByMonth struct {
	Data     [][]TruncatedPayment
	FileName string
}

func (b ByMonth) CreateCSV() {
	file, err := os.Create(b.FileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)

	// Add a header row to csv
	strV = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD

			s := strings.Split(p.FormattedDate, "-")
			date = fmt.Sprintf("%s-%s", s[0], s[1])
		}
		strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}
	writer.WriteAll(strV)
}

func (b ByMonth) OutputData() {
	table := tablewriter.NewWriter(os.Stdout)

	// Add a header row to csv
	table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD

			s := strings.Split(p.FormattedDate, "-")
			date = fmt.Sprintf("%s-%s", s[0], s[1])
		}
		table.Append([]string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}

	table.Render()
}
