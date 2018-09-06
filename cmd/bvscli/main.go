package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/wire"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	ibccli "github.com/cosmos/cosmos-sdk/x/ibc/client/cli"
	stakecli "github.com/cosmos/cosmos-sdk/x/stake/client/cli"

	"github.com/dcgraph/bvs-cosmos/app"
	"github.com/dcgraph/bvs-cosmos/types"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "bvscli",
		Short: "BVS light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// add standard rpc, and tx commands
	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			stakecli.GetCmdQueryValidator("stake", cdc),
			stakecli.GetCmdQueryValidators("stake", cdc),
			stakecli.GetCmdQueryDelegation("stake", cdc),
			stakecli.GetCmdQueryDelegations("stake", cdc),
			authcli.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
			GetCodexCmd("codex", cdc),
		)...)
	rootCmd.AddCommand(client.LineBreak)

	rootCmd.AddCommand(
		client.PostCommands(
			bankcli.SendTxCmd(cdc),
			ibccli.IBCTransferCmd(cdc),
			ibccli.IBCRelayCmd(cdc),
			stakecli.GetCmdCreateValidator(cdc),
			stakecli.GetCmdEditValidator(cdc),
			stakecli.GetCmdDelegate(cdc),
			stakecli.GetCmdUnbond("stake", cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.bvscli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func GetCodexCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "codex [id]",
		Short: "Query codex status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			key := types.Id2StoreKey(id)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("No codex found with the id %s", id)
			}

			codex := &types.Codex{}
			err = cdc.UnmarshalBinaryBare(res, codex)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, codex)
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			return nil
		},
	}
}
