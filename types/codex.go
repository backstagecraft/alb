package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*Codex)(nil)

// A Codex is a service account which is responsible for issueing a new voucher
// according to the value description.
type Codex struct {
	auth.BaseAccount

	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Value       string `json:"value"`
	UnitPrice   int    `json:"unit-price"`
	SaleType    string `json:"sale-type"`
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
	SaleType    string `json:"sale-type"`
	ExpireAfter int    `json:"expire-after"`
	Deposit     int    `json:"deposit"`
	CountTotal  int    `json:"count-total"`
}

//////////////////////////////////////////////////////////////////
// Handling genesis.json

type GenesisCodex struct {
	Id          string         `json:"id"`
	Owner       string         `json:"owner"`
	Value       string         `json:"value"`
	UnitPrice   int            `json:"unit-price"`
	SaleType    string         `json:"sale-type"`
	ExpireAfter int            `json:"expire-after"`
	Deposit     int            `json:"deposit"`
	CountAvail  int            `json:"count-avail"`
	CountLive   int            `json:"count-live"`
	Address     sdk.AccAddress `json:"address"`
	Coins       sdk.Coins      `json:"coins"`
}

func NewGenesisCodex(cod *Codex) *GenesisCodex {
	return &GenesisCodex{
		Owner:       cod.Owner,
		Value:       cod.Value,
		UnitPrice:   cod.UnitPrice,
		SaleType:    cod.SaleType,
		ExpireAfter: cod.ExpireAfter,
		Deposit:     cod.Deposit,
		CountAvail:  cod.CountAvail,
		CountLive:   cod.CountLive,
		Address:     cod.Address,
		Coins:       cod.Coins.Sort(),
	}
}

func (gcod *GenesisCodex) ToCodex() (cod *Codex, err error) {
	return &Codex{
		Id:          gcod.Id,
		Owner:       gcod.Owner,
		Value:       gcod.Value,
		UnitPrice:   gcod.UnitPrice,
		SaleType:    gcod.SaleType,
		ExpireAfter: gcod.ExpireAfter,
		Deposit:     gcod.Deposit,
		CountAvail:  gcod.CountAvail,
		CountLive:   gcod.CountLive,
		BaseAccount: auth.BaseAccount{
			Address: gcod.Address,
			Coins:   gcod.Coins.Sort(),
		},
	}, nil
}
