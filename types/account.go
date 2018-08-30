package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*BvsAccount)(nil)

// BvsAccount is a custom extension for this application. It is an example of
// extending auth.BaseAccount with custom fields. It is compatible with the
// stock auth.AccountStore, since auth.AccountStore uses the flexible go-amino
// library.
type BvsAccount struct {
	auth.BaseAccount

	Name string `json:"name"`
}

// nolint
func (acc BvsAccount) GetName() string      { return acc.Name }
func (acc *BvsAccount) SetName(name string) { acc.Name = name }

// NewBvsAccount returns a reference to a new BvsAccount given a name and an
// auth.BaseAccount.
func NewBvsAccount(name string, baseAcct auth.BaseAccount) *BvsAccount {
	return &BvsAccount{BaseAccount: baseAcct, Name: name}
}

// GetAccountDecoder returns the AccountDecoder function for the custom
// BvsAccount.
func GetAccountDecoder(cdc *wire.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (auth.Account, error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}

		acct := new(BvsAccount)
		err := cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}

		return acct, err
	}
}

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount reflects a genesis account the application expects in it's
// genesis state.
type GenesisAccount struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns a reference to a new GenesisAccount given an
// BvsAccount.
func NewGenesisAccount(aa *BvsAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:    aa.Name,
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// ToBvsAccount converts a GenesisAccount to an BvsAccount.
func (ga *GenesisAccount) ToBvsAccount() (acc *BvsAccount, err error) {
	return &BvsAccount{
		Name: ga.Name,
		BaseAccount: auth.BaseAccount{
			Address: ga.Address,
			Coins:   ga.Coins.Sort(),
		},
	}, nil
}
