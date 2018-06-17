// Various helper functions... most deal with the formatting of date/time.

package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func decodeResponse(resp *http.Response, object interface{}) (err error) {
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &object)
	if err != nil {
		return
	}
	return
}

func stringToUnixTimeStamp(t string) int {
	var year, month, date string
	s := strings.Split(t, "-")
	if len(s) == 3 {
		year, month, date = s[0], s[1], s[2]
		if len(date) == 1 {
			date = fmt.Sprintf("0%v", date)
		}
	} else {
		year, month, date = s[0], s[1], "01"
	}

	formattedTime, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v-%v-%vT15:04:05Z", year, month, date))
	return int(formattedTime.Unix())
}

func stringToDateLumenFormat(ti string) string {
	t := strings.Fields(ti)
	mon, date, year := t[0], t[1][:len(t[1])-1], t[2]
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("%v %v %v 12:00 MST", date, mon, year[2:]))

	return timeToString(formattedTime)
}

func stringToDateCurrencylayerFormat(t string) string {
	formattedTime, _ := time.Parse(time.RFC3339, t)
	return timeToString(formattedTime)
}

func timeToString(formattedDate time.Time) string {
	y, month, day := formattedDate.Date()
	d := strconv.Itoa(day)
	m := monthStringToNumber(month)
	if len(d) == 1 {
		d = fmt.Sprintf("0%v", d)
	}

	return fmt.Sprintf("%v-%v-%v", y, m, d)
}

func monthStringToNumber(month time.Month) string {
	m := ""
	switch month.String()[:3] {
	case "Jan":
		m = "01"
		break
	case "Feb":
		m = "02"
		break
	case "Mar":
		m = "03"
		break
	case "Apr":
		m = "04"
		break
	case "May":
		m = "05"
		break
	case "Jun":
		m = "06"
		break
	case "Jul":
		m = "07"
		break
	case "Aug":
		m = "08"
		break
	case "Sep":
		m = "09"
		break
	case "Oct":
		m = "10"
		break
	case "Nov":
		m = "11"
		break
	case "Dec":
		m = "12"
		break
	}

	return m
}

// OrderDataByDate sorts data by date
func OrderDataByDate(data map[string][]TruncatedPayment) [][]TruncatedPayment {
	var keys []int
	keyToString := make(map[int]string)
	for k := range data {
		keys = append(keys, stringToUnixTimeStamp(k))
		keyToString[stringToUnixTimeStamp(k)] = k
	}

	sort.Ints(keys)

	orderedData := make([][]TruncatedPayment, 0)

	for _, k := range keys {
		orderedData = append(orderedData, data[keyToString[k]])
	}

	return orderedData
}

// OrderDataByMonth sorts data by month
func OrderDataByMonth(data map[string][]TruncatedPayment) [][]TruncatedPayment {
	m := make(map[string][]TruncatedPayment)
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
		keys = append(keys, stringToUnixTimeStamp(k))
		keyToString[stringToUnixTimeStamp(k)] = k
	}

	sort.Ints(keys)

	orderedData := make([][]TruncatedPayment, 0)

	for _, k := range keys {
		orderedData = append(orderedData, m[keyToString[k]])
	}

	return orderedData
}
