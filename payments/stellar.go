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

	// Paging token used to increment through the payments sequentially
	pagingToken := "0"

	for lens == 200 {
		url := fmt.Sprintf("https://horizon.stellar.org/accounts/%s/payments?limit=200&order=asc&cursor=%s", account, pagingToken)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var t Page
		err = decodeResponse(resp, &t)
		if err != nil {
			log.Fatal(err)
		}

		lens := len(t.Embedded.Records)
		if lens < 1 {
			break
		}
		pagingToken = t.Embedded.Records[lens-1].PagingToken

		// Filter out merge account and create account methods
		filterFunction := func(p Payment) bool {
			return (p.Account == "")
		}

		p := filterPayments(t.Embedded.Records, filterFunction, account)

		payments = append(payments, p...)
	}
	return payments
}

func filterPayments(p []Payment, f func(Payment) bool, account string) []TruncatedPayment {
	payments := make([]TruncatedPayment, 0)
	for _, v := range p {
		if f(v) {

			// Here insert mappings for tokens that don't quite
			// reflect the FIAT ticker.
			// Ex: EURT --> EUR
			var asset string
			if v.AssetCode == "EURT" {
				asset = "EUR"
			} else {
				asset = v.AssetCode
			}

			// Here capture two metadata variables for the raw output:
			// 	1. Whether the payment was sent or received
			// 	2. The counter party address (if sent, the recipient, if received, the sender)
			var sentrevc string
			var fromto string
			if v.From == account {
				sentrevc = "Sent"
				fromto = v.To
			} else {
				sentrevc = "Received"
				fromto = v.From
			}

			payments = append(payments, TruncatedPayment{
				CreatedAt:     v.CreatedAt,
				FormattedDate: stringToDateCurrencylayerFormat(v.CreatedAt),
				AssetCode:     asset,
				Amount:        v.Amount,
				VolumeUSD:     0,
				SentRecv:      sentrevc,
				FromTo:        fromto,
			})
		}
	}
	return payments
}
