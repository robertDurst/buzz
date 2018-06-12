// Post Process contains methods for sorting and outputing data.

package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
)

// Create a CSV from a list of payments
func createCSV(data [][]TruncatedPayment, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Slice for capturing data to append to csv
	strV := make([][]string, 0)
	// Add a header row to csv
	strV = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Volume in USD"})
	for _, v := range data {
		for _, p := range v {
			strV = append(strV, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Volume_USD, 'f', 6, 64)})
		}
	}

	writer.WriteAll(strV)
}

// Sort data by date
func orderData(data map[Date][]TruncatedPayment) [][]TruncatedPayment {
	var keys []int
	keyToString := make(map[int]string)
	for k := range data {
		keys = append(keys, stringToDateSortFormat(k))
		keyToString[stringToDateSortFormat(k)] = k
	}

	sort.Ints(keys)

	orderedData := make([][]TruncatedPayment, 0)

	for _, k := range keys {
		orderedData = append(orderedData, data[keyToString[k]])
	}

	return orderedData
}
