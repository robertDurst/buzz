package payments

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guptarohit/asciigraph"
	"github.com/olekukonko/tablewriter"
)

type Aggregate interface {
	FormatCSV(strV [][]string) (csv [][]string)
	FormatTable(table *tablewriter.Table)
	FormatMarkdown(table *tablewriter.Table)
	FormatGraph()
}

func CreateCSV(a Aggregate, filename string) {
	file, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)

	// Write formatted data to csv
	writer.WriteAll(a.FormatCSV(strV))
}

func CreateTable(a Aggregate) {
	t := tablewriter.NewWriter(os.Stdout)
	a.FormatTable(t)
}

func CreateMarkdown(a Aggregate) {
	t := tablewriter.NewWriter(os.Stdout)
	a.FormatMarkdown(t)
}

func CreateGraph(a Aggregate) {
	a.FormatGraph()
}

type Raw struct {
	Data [][]TruncatedPayment
}

func (r Raw) FormatCSV(strV [][]string) (csv [][]string) {
	// Add a header row to csv
	csv = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
	for _, v := range r.Data {
		for _, p := range v {
			csv = append(csv, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
		}
	}
	return
}

func (r Raw) FormatTable(table *tablewriter.Table) {
	// Add a header row to table
	table.SetHeader([]string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
	for _, v := range r.Data {
		for _, p := range v {
			table.Append([]string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
		}
	}

	table.Render()
}

func (r Raw) FormatMarkdown(table *tablewriter.Table) {
	// Add a header row to table
	table.SetHeader([]string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Price", "Volume in USD", "Sent or Received", "From/To"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, v := range r.Data {
		for _, p := range v {
			table.Append([]string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Price, 'f', 6, 64), strconv.FormatFloat(p.VolumeUSD, 'f', 6, 64), p.SentRecv, p.FromTo})
		}
	}

	table.Render()
}

func (r Raw) FormatGraph() {
	data := make([]float64, 0)
	for _, v := range r.Data {
		for _, p := range v {
			data = append(data, p.VolumeUSD)
		}
	}

	graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50), asciigraph.Caption("Volume vs. Time (Raw -- no Aggregation)"))
	fmt.Println(graph)
}

type ByDate struct {
	Data [][]TruncatedPayment
}

func (b ByDate) FormatCSV(strV [][]string) (csv [][]string) {
	// Add a header row to csv
	csv = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD
			date = p.FormattedDate
		}
		csv = append(csv, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}
	return
}

func (b ByDate) FormatTable(table *tablewriter.Table) {
	// Add a header row to table
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

func (b ByDate) FormatMarkdown(table *tablewriter.Table) {
	// Add a header row to table
	table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
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

func (b ByDate) FormatGraph() {
	data := make([]float64, 0)
	for _, v := range b.Data {
		volume := 0.0
		for _, p := range v {
			volume += p.VolumeUSD
		}
		data = append(data, volume)
	}

	graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50), asciigraph.Caption("Volume vs. Time (Aggregated by Day)"))
	fmt.Println(graph)
}

type ByMonth struct {
	Data [][]TruncatedPayment
}

func (b ByMonth) FormatCSV(strV [][]string) (csv [][]string) {
	// Add a header row to csv
	csv = append(strV, []string{"Created At (Pretty)", "Volume in USD"})
	for _, v := range b.Data {
		volume := 0.0
		date := ""
		for _, p := range v {
			volume += p.VolumeUSD

			s := strings.Split(p.FormattedDate, "-")
			date = fmt.Sprintf("%s-%s", s[0], s[1])
		}
		csv = append(csv, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}

	return
}

func (b ByMonth) FormatTable(table *tablewriter.Table) {
	// Add a header row to table
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

func (b ByMonth) FormatMarkdown(table *tablewriter.Table) {
	// Add a header row to table
	table.SetHeader([]string{"Created At (Pretty)", "Volume in USD"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
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

func (b ByMonth) FormatGraph() {
	data := make([]float64, 0)
	for _, v := range b.Data {
		volume := 0.0
		for _, p := range v {
			volume += p.VolumeUSD
		}
		data = append(data, volume)
	}

	graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50), asciigraph.Caption("Volume vs. Time (Aggregated by Month)"))
	fmt.Println(graph)
}
