// Stellar contains a method to grab all payments for a Stellar account.
package payments

import (
	"fmt"
	"log"
	"net/http"
)

func PaymentsForAccount(account string) []TruncatedPayment {
	payments := make([]TruncatedPayment, 0)

	lens := 200
	pagingToken := "0"

	for lens == 200 {
		url := fmt.Sprintf("https://horizon.stellar.org/accounts/%s/payments?limit=200&order=asc&cursor=%s", account, pagingToken)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var t PaymentsPage
		err = decodeResponse(resp, &t)
		if err != nil {
			log.Fatal(err)
		}

		lens := len(t.Embedded.Records)
		if lens < 1 {
			break
		}
		pagingToken = t.Embedded.Records[lens-1].PagingToken

		filterFunction := func(p Payment) bool {
			return (p.Account == "")
		}

		p := filterPayments(t.Embedded.Records, filterFunction)

		payments = append(payments, p...)
	}
	return payments
}

func filterPayments(vs []Payment, f func(Payment) bool) []TruncatedPayment {
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
				FormattedDate: stringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     asset,
				Amount:        v.Amount,
				Volume_USD:    0,
			})
		}
	}
	return vsf
}
