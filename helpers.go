package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Aggregate(vs []TruncatedPayment) (sortedByDate map[Date][]TruncatedPayment, dates []Date, assetString string) {
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

func decodeResponse(resp *http.Response, object interface{}) (err error) {
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &object)
	if err != nil {
		return
	}
	return
}

func FillInVolumePerPayment(payments []TruncatedPayment) map[Date][]TruncatedPayment {
	// Aggregate the Data. Returns the following:
	// sortedByDate -- map[Date][]TruncatedPayment
	// dates -- []Date
	// assets -- string (ASSET1,ASSET2,ASSET3)
	all_data, dates, assetString := Aggregate(payments)
	new_data := make(map[Date][]TruncatedPayment)

	// Get lumens prices for all time
	lumenPrices := GetStellarHistoricalData()

	for _, v := range dates {
		url := fmt.Sprintf("http://apilayer.net/api/live?access_key=%s&currencies=%s&date=%s", "2338486eeca99d18437949e73d6b9027", assetString, v)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept-Encoding", "gzip, deflate")
		res, _ := client.Do(req)

		defer res.Body.Close()

		var t CurrencyExchangeResponse
		decodeResponse(res, &t)

		new_data[v] = updateVolumeForDate(all_data[v], v, t, lumenPrices[v])
	}

	return new_data
}

func Filter(vs []Payment, f func(Payment) bool) []TruncatedPayment {
	vsf := make([]TruncatedPayment, 0)
	for _, v := range vs {
		if f(v) {

			var asset string
			if v.AssetCode == "EURT" {
				asset = "EUR"
			} else {
				asset = v.AssetCode
			}

			vsf = append(vsf, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: StringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     asset,
				Amount:        v.Amount,
				Volume_USD:    0,
			})
		}
	}
	return vsf
}

func StringToDateSortFormat(t string) int {
	s := strings.Split(t, "-")
	year, month, date := s[0], s[1], s[2]
	if len(date) == 1 {
		date = fmt.Sprintf("0%v", date)
	}
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("%v %v %v 12:00 MST", date, month[:3], year[2:]))
	return int(formattedTime.Unix())
}

func StringToDateLumenFormat(t string) string {
	month, date, year := strings.Fields(t)[0], strings.Fields(t)[1][:len(strings.Fields(t)[1])-1], strings.Fields(t)[2]
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("%v %v %v 12:00 MST", date, month, year[2:]))
	y, w, d := formattedTime.Date()
	return fmt.Sprintf("%v-%v-%v", y, w, d)
}

func StringToDateCurrencylayerFormat(t string) string {
	date, _ := time.Parse(time.RFC3339, t)
	y, w, d := date.Date()
	return fmt.Sprintf("%v-%v-%v", y, w, d)
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
				FormattedDate: StringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     v.AssetCode,
				Amount:        v.Amount,
				Volume_USD:    volume,
			})
		} else {
			vsfz = append(vsfz, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: StringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     "XLM",
				Amount:        v.Amount,
				Volume_USD:    amt * lp,
			})
		}
	}
	return vsfz
}
