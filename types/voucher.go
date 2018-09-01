package types

// A Voucher is an asset representing a value guaranteed by a voucher issuer.
type Voucher struct {
	Id      string `json:"id"`
	Creator string `json:"creator"` // may be a Codex or a Dealer
	Holder  string `json:"holder"`
	Value   string // from codex
	Price   string // from codex
}
