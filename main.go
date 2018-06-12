package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Get all payments as truncated payments structs
	// Created At
	// Asset Code -- will be "" if XLM
	// Amount
	payments := PaymentsForAccount()
	data := FillInVolumePerPayment(payments)

	createCSV(data)
}

func createCSV(data map[Date][]TruncatedPayment) {
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	strV := make([][]string, 0)

	strV = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Volume in USD"})
	// strV = append(strV, []string{"Created At (Pretty)", "Volume in USD"})

	var keys []int
	keyToString := make(map[int]string)
	for k := range data {
		keys = append(keys, StringToDateSortFormat(k))
		keyToString[StringToDateSortFormat(k)] = k
	}

	sort.Ints(keys)

	// To perform the opertion you want
	for _, k := range keys {
		// volume := 0.0
		// date := ""
		for _, p := range data[keyToString[k]] {
			// volume += p.Volume_USD
			// date = p.FormattedDate
			strV = append(strV, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Volume_USD, 'f', 6, 64)})
		}

		// strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
	}

	writer.WriteAll(strV)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
