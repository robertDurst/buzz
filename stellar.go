package main

import (
	"fmt"
	"log"
	"net/http"
)

func PaymentsForAccount() []TruncatedPayment {
	payments := make([]TruncatedPayment, 0)
	account := "GCKX3XVTPVNFXQWLQCIBZX6OOPOIUT7FOAZVNOFCNEIXEZFRFSPNZKZT"

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
		if lens < 200 {
			break
		}
		pagingToken = t.Embedded.Records[lens-1].PagingToken

		filterFunction := func(p Payment) bool {
			return (p.Account == "")
		}

		p := Filter(t.Embedded.Records, filterFunction)

		payments = append(payments, p...)
	}

	return payments
}
