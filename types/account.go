package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*UserAccount)(nil)

// UserAccount is a custom extension for this application. It is an example of
// extending auth.BaseAccount with custom fields. It is compatible with the
// stock auth.AccountStore, since auth.AccountStore uses the flexible go-amino
// library.
type UserAccount struct {
	auth.BaseAccount

	Id string `json:"id"`
	//Vouchers []string // from voucher list
}

func (acc UserAccount) GetId() string    { return acc.Id }
func (acc *UserAccount) SetId(id string) { acc.Id = id }

// NewUserAccount returns a reference to a new UserAccount given an id and an
// auth.BaseAccount.
func NewUserAccount(id string, baseAcct auth.BaseAccount) *UserAccount {
	return &UserAccount{BaseAccount: baseAcct, Id: id}
}

// GetAccountDecoder returns the AccountDecoder function for the custom
// UserAccount.
func GetAccountDecoder(cdc *wire.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (auth.Account, error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}

		acct := new(UserAccount)
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
	Id      string         `json:"id"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns a reference to a new GenesisAccount given an
// UserAccount.
func NewGenesisAccount(aa *UserAccount) *GenesisAccount {
	return &GenesisAccount{
		Id:      aa.Id,
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// ToUserAccount converts a GenesisAccount to an UserAccount.
func (ga *GenesisAccount) ToUserAccount() (acc *UserAccount, err error) {
	return &UserAccount{
		Id: ga.Id,
		BaseAccount: auth.BaseAccount{
			Address: ga.Address,
			Coins:   ga.Coins.Sort(),
		},
	}, nil
}
