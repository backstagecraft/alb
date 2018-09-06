package types

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
	Codices  []*GenesisCodex   `json:"codices"`
}
