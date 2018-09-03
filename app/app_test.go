package app

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dcgraph/bvs-cosmos/types"
)

func setGenesis(bvsApp *BvsApp, accounts ...*types.BvsAccount) (types.GenesisState, error) {
	genAccts := make([]*types.GenesisAccount, len(accounts))
	for i, appAct := range accounts {
		genAccts[i] = types.NewGenesisAccount(appAct)
	}

	genesisState := types.GenesisState{Accounts: genAccts}
	stateBytes, err := wire.MarshalJSONIndent(bvsApp.cdc, genesisState)
	if err != nil {
		return types.GenesisState{}, err
	}

	// initialize and commit the chain
	bvsApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})
	bvsApp.Commit()

	return genesisState, nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	bvsApp := NewBvsApp(logger, db)

	// construct a pubkey and an address for the test account
	pubkey := ed25519.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	// construct some test coins
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)

	// create an auth.BaseAccount for the given test account and set it's coins
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	err = baseAcct.SetCoins(coins)
	require.Nil(t, err)

	// create a new test BvsAccount with the given auth.BaseAccount
	bvsAcct := types.NewBvsAccount("foobar", baseAcct)
	genState, err := setGenesis(bvsApp, bvsAcct)
	require.Nil(t, err)

	// create a context for the BaseApp
	ctx := bvsApp.BaseApp.NewContext(true, abci.Header{})
	res := bvsApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, bvsAcct, res)

	// reload app and ensure the account is still there
	bvsApp = NewBvsApp(logger, db)

	stateBytes, err := wire.MarshalJSONIndent(bvsApp.cdc, genState)
	require.Nil(t, err)

	// initialize the chain with the expected genesis state
	bvsApp.InitChain(abci.RequestInitChain{
		Validators: []abci.Validator{}, AppStateBytes: stateBytes,
	})

	ctx = bvsApp.BaseApp.NewContext(true, abci.Header{})
	res = bvsApp.accountMapper.GetAccount(ctx, baseAcct.Address)
	require.Equal(t, bvsAcct, res)
}
