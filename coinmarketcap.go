package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// GetStellarHistoricalData retreives all historical data of XLM by date.
// Note: here we use the closing price as there is not an average price for the day.
func GetStellarHistoricalData() map[string]float64 {
	resp, err := http.Get("https://coinmarketcap.com/currencies/stellar/historical-data/?start=20130428&end=20180612")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data = make(map[string]float64)
	var price float64
	var date string

	// Parse HTML document
	doc.Find(".table-responsive .table tbody tr").Each(func(_ int, s *goquery.Selection) {
		s.Find("td").Each(func(q int, z *goquery.Selection) {
			// Capture the date, in Month Date, Year format
			if q == 0 {
				date = z.Text()
			}

			// Capture the close price
			if q == 4 {
				price, err = strconv.ParseFloat(z.Text(), 64)
				if err != nil {
					log.Fatal(err)
				}
			}
		})
		data[StringToDateLumenFormat(date)] = price
	})

	return data
}
