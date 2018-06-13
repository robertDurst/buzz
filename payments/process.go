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

	"github.com/fatih/color"
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

	// Capture unmatched assets
	priceCapture := make(map[string][]float64)

	for _, v := range dates {
		url := fmt.Sprintf("http://apilayer.net/api/historical?access_key=%s&currencies=%s&date=%s", apikey, assetString, v)
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

		for k, v := range t.Quotes {
			if _, ok := priceCapture[k]; !ok {
				priceCapture[k] = []float64{v}
			} else {
				priceCapture[k] = append(priceCapture[k], v)
			}
		}

		returnData[v] = updateVolumeForDate(allData[v], v, t, lumenPrices[v])

	}

	resultMsg, unmatchedAssets := currencylayerIntegrityCheck(priceCapture, len(dates))

	switch resultMsg {
	case "api bad or all non-fiat assets":
		color.Red("Either bad API key, or this account only has non-native asset payments in payment history.")
		color.Blue("The following assets were not matched: %s", assetString)
		break
	case "api bad part way through":
		color.Red("Part way through the API hit its limit")
		break
	case "ok some non-natives":
		color.Blue("The following assets were not matched: %s", unmatchedAssets)
		break
	case "ok":
		color.Green("Success!")
		break
	}

	return returnData
}

func currencylayerIntegrityCheck(m map[string][]float64, expectedLength int) (string, string) {
	var buffer bytes.Buffer
	totalZero := make([]string, 0)
	partialZero := make([]string, 0)

	z := 0

	for asset, prices := range m {
		sum := 0.0
		hasZero := false
		z += len(prices)
		for _, v := range prices {
			if v == 0 {
				hasZero = true
			}
			sum += v
		}
		if sum == 0.0 {
			totalZero = append(totalZero, asset)
			buffer.WriteString(asset)
			buffer.WriteString(" ")
		} else if hasZero {
			partialZero = append(partialZero, asset)
		}
	}

	// All zero meaning API totally failed and/or all fiat assets
	if len(m) == len(totalZero) {
		return "api bad or all non-fiat assets", buffer.String()
	}

	// Some partial zeros meaning assets that worked to start eventually failed
	// due to api hitting the limit. This may happen in combination with non -
	// native, unmatched assets.
	if expectedLength > z {
		return "api bad part way through", buffer.String()
	}

	// All fiat tokens worked, however still had some non-native assets
	if len(totalZero) > 0 {
		return "ok some non-native", buffer.String()
	}

	// All fiat assets and everything worked perfectly
	return "ok", buffer.String()
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
				SentRecv:      v.SentRecv,
				FromTo:        v.FromTo,
				Price:         1 / exchange,
			})
		} else {
			vsfz = append(vsfz, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: stringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     "XLM",
				Amount:        v.Amount,
				Volume_USD:    amt * lp,
				SentRecv:      v.SentRecv,
				FromTo:        v.FromTo,
				Price:         lp,
			})
		}
	}
	return vsfz
}
