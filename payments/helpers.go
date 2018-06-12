// Various helper functions... most deal with the formatting of date/time.

package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func stringToDateSortFormat(t string) int {
	s := strings.Split(t, "-")
	year, month, date := s[0], s[1], s[2]
	if len(date) == 1 {
		date = fmt.Sprintf("0%v", date)
	}
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("%v %v %v 12:00 MST", date, month[:3], year[2:]))
	return int(formattedTime.Unix())
}

func stringToDateSortFormat2(t string) int {
	s := strings.Split(t, "-")
	year, month := s[0], s[1]
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("01 %v %v 12:00 MST", month[:3], year[2:]))
	return int(formattedTime.Unix())
}

func stringToDateLumenFormat(t string) string {
	month, date, year := strings.Fields(t)[0], strings.Fields(t)[1][:len(strings.Fields(t)[1])-1], strings.Fields(t)[2]
	formattedTime, _ := time.Parse(time.RFC822, fmt.Sprintf("%v %v %v 12:00 MST", date, month, year[2:]))
	y, w, d := formattedTime.Date()
	return fmt.Sprintf("%v-%v-%v", y, w, d)
}

func stringToDateCurrencylayerFormat(t string) string {
	date, _ := time.Parse(time.RFC3339, t)
	y, w, d := date.Date()
	return fmt.Sprintf("%v-%v-%v", y, w, d)
}
