package payments

// CurrencyExchangeResponse is used to
// capture the currencylayer responses.
type CurrencyExchangeResponse struct {
	Quotes map[string]float64 `json:"quotes"`
}

// Payment is a data structure copy-pasta
// from the Horizon Client. It is not great
// since it just lump sums all payments
// together. However it is fine for our
// purpose here since we just use the
// Account field to filter out create
// account and merge account operations.
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

// PaymentsPage removes the Link variable
// since it is not necessary.
type Page struct {
	Embedded struct {
		Records []Payment `json:"records"`
	} `json:"_embedded"`
}

// TruncatedPayment is a custom data struct that represents
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
