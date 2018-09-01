package types

// A Codex is a service account which is responsible for issueing a new voucher
// according to the value description.
type Codex struct {
	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Value       string `json:"value"`
	UnitPrice   int    `json:"unit-price"`
	SaleType    int    `json:"sale-type"`
	ExpireAfter int    `json:"expire-after"`
	Deposit     int    `json:"deposit"`
	CountAvail  int    `json:"count-avail"`
	CountLive   int    `json:"count-live"`
}

// CodexDef is a definition of a new codex to be created.
type CodexDef struct {
	Owner       string `json:"owner"`
	Value       string `json:"value"`
	UnitPrice   int    `json:"unit-price"`
	SaleType    int    `json:"sale-type"`
	ExpireAfter int    `json:"expire-after"`
	Deposit     int    `json:"deposit"`
	CountTotal  int    `json:"count-total"`
}
