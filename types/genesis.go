package types

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
	Codices  []*Codex          `json:"codices"`
	Vouchers []*Voucher        `json:"vouchers"`
}
