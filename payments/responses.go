package payments

type CurrencyExchangeResponse struct {
	Quotes map[string]float64 `json:"quotes"`
}

type Payment struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	PagingToken string `json:"paging_token"`

	Links struct {
		Transaction struct {
			Href string `json:"href"`
		} `json:"transaction"`
	} `json:"_links"`

	TransactionHash string `json:"transaction_hash"`
	SourceAccount   string `json:"source_account"`
	CreatedAt       string `json:"created_at"`

	// create_account and account_merge field
	Account string `json:"account"`

	// create_account fields
	Funder          string `json:"funder"`
	StartingBalance string `json:"starting_balance"`

	// account_merge fields
	Into string `json:into"`

	// payment/path_payment fields
	From        string `json:"from"`
	To          string `json:"to"`
	AssetType   string `json:"asset_type"`
	AssetCode   string `json:"asset_code"`
	AssetIssuer string `json:"asset_issuer"`
	Amount      string `json:"amount"`

	// transaction fields
	Memo struct {
		Type  string `json:"memo_type"`
		Value string `json:"memo"`
	}
}

// Notice I removed the Link variable
// since I don't think it is necessary.
type PaymentsPage struct {
	Embedded struct {
		Records []Payment `json:"records"`
	} `json:"_embedded"`
}

// Custom data struct that represents
// a subset of the payments fields
// that I have deemed interesting.
type TruncatedPayment struct {
	CreatedAt     string
	FormattedDate string
	AssetCode     string
	Amount        string
	VolumeUSD     float64
	SentRecv      string
	FromTo        string
	Price         float64
}
