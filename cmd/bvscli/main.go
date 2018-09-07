package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/wire"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	//bank "github.com/cosmos/cosmos-sdk/x/bank"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	//bankclient "github.com/cosmos/cosmos-sdk/x/bank/client"
	ibccli "github.com/cosmos/cosmos-sdk/x/ibc/client/cli"
	stakecli "github.com/cosmos/cosmos-sdk/x/stake/client/cli"

	"github.com/dcgraph/bvs-cosmos/app"
	"github.com/dcgraph/bvs-cosmos/types"
	"github.com/dcgraph/bvs-cosmos/x/shop"
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
			GetVoucherCmd("voucher", cdc),
		)...)
	rootCmd.AddCommand(client.LineBreak)

	rootCmd.AddCommand(
		client.PostCommands(
			bankcli.SendTxCmd(cdc), // TODO
			//BvsSendCmd(cdc), // TODO
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
			key := types.Id2StoreKey("codex:", id)
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

func GetVoucherCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "voucher [id]",
		Short: "Query voucher status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			key := types.Id2StoreKey("voucher:", id)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("No voucher found with the id %s", id)
			}

			voucher := &types.Voucher{}
			err = cdc.UnmarshalBinaryBare(res, voucher)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, voucher)
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			return nil
		},
	}
}

func BvsSendCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send assets to an account",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcli.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			accAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			sender := accAddress.String()

			// TODO: check if the recipient exists
			recp := viper.GetString("recp")

			// parse coins trying to be sent
			assetStr := viper.GetString("asset")
			asset := types.ParseBvsAsset(assetStr)
			if !types.IsOwner(sender, asset) {
				return errors.Errorf("Can't send asset. Invalid ownership.")
			}
			msg := shop.BuildBvsMsg(accAddress, sender, recp, asset)

			if len(asset.Coins) > 0 {
				// ensure account has enough coins
				account, err := cliCtx.GetAccount(accAddress)
				if err != nil {
					return errors.Errorf("Failed to get account with the address %s.", accAddress)
				}
				if !account.GetCoins().IsGTE(asset.Coins) {
					return errors.Errorf("Address %s insufficient assets to send.", accAddress)
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String("recp", "", "Recipient address")
	cmd.Flags().String("asset", "", "List of assets to send")

	return cmd
}
