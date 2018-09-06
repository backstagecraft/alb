package types

// A Voucher is an asset representing a value guaranteed by a voucher issuer.
type Voucher struct {
	Id       string `json:"id"`
	Origin   string `json:"origin"` // may be a Codex or a Dealer
	Holder   string `json:"holder"`
	ExpireOn int    `json:"expire-on"`
}
