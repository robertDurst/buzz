// Process is the code used to process all the received payment data.
// It requests exchange prices for each fiat to USD and then updates
// the USD Volume per payment (includes XLM-USD as well).
package payments

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

func aggregateData(vs []TruncatedPayment) (sortedByDate map[Date][]TruncatedPayment, dates []Date, assetString string) {
	var buffer bytes.Buffer
	sortedByDate = make(map[Date][]TruncatedPayment)
	assets := make(map[string]bool)
	dates = make([]Date, 0)
	for _, v := range vs {
		if _, ok := sortedByDate[v.FormattedDate]; ok {
			sortedByDate[v.FormattedDate] = append(sortedByDate[v.FormattedDate], v)
		} else {
			sortedByDate[v.FormattedDate] = []TruncatedPayment{v}
			dates = append(dates, v.FormattedDate)

		}
		if _, ok := assets[v.AssetCode]; !ok {
			assets[v.AssetCode] = true
			buffer.WriteString(v.AssetCode)
			buffer.WriteString(",")
		}
	}

	assetString = buffer.String()
	if assetString[0] == ',' {
		assetString = assetString[1 : len(assetString)-2]
	} else {
		assetString = assetString[:len(assetString)-2]
	}

	return
}

func FillInVolumePerPayment(payments []TruncatedPayment, apikey string) map[Date][]TruncatedPayment {
	// Aggregate the Data. Returns the following:
	// sortedByDate -- map[Date][]TruncatedPayment
	// dates -- []Date
	// assets -- string (ASSET1,ASSET2,ASSET3)
	allData, dates, assetString := aggregateData(payments)
	returnData := make(map[Date][]TruncatedPayment)

	// Get lumens prices for all time
	lumenPrices := getStellarHistoricalData()

	for _, v := range dates {
		url := fmt.Sprintf("http://apilayer.net/api/live?access_key=%s&currencies=%s&date=%s", apikey, assetString, v)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept-Encoding", "gzip, deflate")
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		var t CurrencyExchangeResponse
		decodeResponse(res, &t)

		returnData[v] = updateVolumeForDate(allData[v], v, t, lumenPrices[v])
	}

	return returnData
}

func updateVolumeForDate(vsf []TruncatedPayment, date Date, cer CurrencyExchangeResponse, lp float64) []TruncatedPayment {
	vsfz := make([]TruncatedPayment, 0)
	for _, v := range vsf {
		amt, _ := strconv.ParseFloat(v.Amount, 64)

		if v.AssetCode != "" {
			name := fmt.Sprintf("USD%s", v.AssetCode)
			exchange := cer.Quotes[name]
			volume := amt / exchange
			if math.IsInf(volume, 1) {
				volume = 0
			}

			vsfz = append(vsfz, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: stringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     v.AssetCode,
				Amount:        v.Amount,
				Volume_USD:    volume,
			})
		} else {
			vsfz = append(vsfz, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: stringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     "XLM",
				Amount:        v.Amount,
				Volume_USD:    amt * lp,
			})
		}
	}
	return vsfz
}
