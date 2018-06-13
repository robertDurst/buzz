// Post Process contains methods for sorting and outputing data.
package payments

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
				volume += p.Volume_USD
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
				volume += p.Volume_USD

				s := strings.Split(p.FormattedDate, "-")
				date = fmt.Sprintf("%s-%s", s[0], s[1])
			}
			strV = append(strV, []string{date, strconv.FormatFloat(volume, 'f', 6, 64)})
		}
		break
	default:
		// Add a header row to csv
		strV = append(strV, []string{"Created At (Raw)", "Created At (Pretty)", "Asset Code", "Amount", "Volume in USD"})
		for _, v := range data {
			for _, p := range v {
				strV = append(strV, []string{p.CreatedAt, p.FormattedDate, p.AssetCode, p.Amount, strconv.FormatFloat(p.Volume_USD, 'f', 6, 64)})
			}
		}
	}

	writer.WriteAll(strV)
}

// Sort data by date
func OrderDataByDate(data map[Date][]TruncatedPayment) [][]TruncatedPayment {
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

// Sort data by date
func OrderDataByMonth(data map[Date][]TruncatedPayment) [][]TruncatedPayment {
	m := make(map[Date][]TruncatedPayment)
	for k, v := range data {
		s := strings.Split(k, "-")
		y := fmt.Sprintf("%s-%s", s[0], s[1])

		if _, ok := m[y]; !ok {
			m[y] = v
		} else {
			m[y] = append(m[y], v...)
		}
	}

	var keys []int
	keyToString := make(map[int]string)
	for k := range m {
		keys = append(keys, stringToDateSortFormat2(k))
		keyToString[stringToDateSortFormat2(k)] = k
	}

	sort.Ints(keys)

	orderedData := make([][]TruncatedPayment, 0)

	for _, k := range keys {
		orderedData = append(orderedData, m[keyToString[k]])
	}

	return orderedData
}
