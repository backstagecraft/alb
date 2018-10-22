# bvs-cosmos
BVS demo using cosmos-sdk

# work logs
1. Copied basecoin example from cosmos-sdk
1. Rebranding: Basecoin &rarr; BVS or Bvs
1. Added test-purpose genesis.json with pre-allocated assets
1. Added Travis CI configuration (build only for now, no deploy)
1. Added several types for BVS application

# build and install
1. install golang (see [cosmos doc](https://cosmos.network/docs/getting-started/installation.html))
1. install dep command `go get dep`
1. clone this repo into `$GOPATH/src/github.com/dcgraph/bvs-cosmos`
1. cd to `$GOPATH/src/github.com/dcgraph/bvs-cosmos`
1. run `dep ensure -v`
1. run `make` (build and test)
1. run `make install` (build and install)

# bvsd test run
1. run `bvsd init`
1. copy `testdata/genesis.json` to `$HOME/.bvsd/config/genesis.json` or simply run `make install` again
1. run `bvsd start`
1. wait until the first block generated and committed
1. press `^C` to stop bvsd
1. run `bvsd export` to see internal state of bvsd

# cli test
In one terminal:
1. run `bvsd start` (leave this terminal)

In another terminal:
1. run `bvscli status`
1. run `bvscli account ...`
1. run `bvscli codex ...`
1. run `bvscli voucher ...`
1. run `bvscli send ...`
1. ...

See [Gaia CLI document](https://cosmos.network/docs/sdk/clients.html#gaia-cli) for more information.

# notes on codes
- `BvsApp` in `app/app.go` is a fork of `Basecoin` example, which implements ABCI application in cosmos terms.
- types:
    - `UserAccount` means every user account associated with a private key.
        - `UserAccount` is an extension of `BaseAccount`. So, in order to handle genesis state this type needs `GenesisAccount` which has flat member fields.
        - `UserAccount` has `Address`(type `sdk.AccAddress`) member to handle operations associated with digital signatures.
        - `UserAccount` also has `Id`(type `string`) member to handle BVS-related operations.
    - `Codex` and `Voucher` are system account without private key associated.
        - They are simple structs and do not need `GenesisSomething` struct to handle genesis state.
        - `Codex` and `Voucher` have `Id` member to handle BVS-related operations.
- data stores:
    - `Account` KVStore stores user accounts with the help of `AccountMapper`.
    - `Codex` KVStore stores codex accounts with the help of `CodexMapper`.
    - `Voucher` KVStore stores vouchers with the help of `VoucherMapper`.
- BVS asset
    - `BvsAsset` in `types/voucher.go` represent arbitrary asset in BVS environment.
    - `BvsAsset` shall be moved into `types/asset.go`.
    - `UserAccount` and `Codex` can hold coins, which are handled as `sdk.Coins`.
    - `UserAccount` and `Codex` use `cosmos-sdk/x/bank` module in order to handle `sdk.Coins`.
