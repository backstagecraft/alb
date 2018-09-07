package shop

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dcgraph/bvs-cosmos/types"
)

// Msg

type MsgBvs struct {
	SenderAccount sdk.AccAddress `json:"sender-account"`
	Sender        string         `json:"sender"`
	Recipient     string         `json:"recipient"`
	Asset         types.BvsAsset `json:"asset"`
}

var _ sdk.Msg = MsgBvs{}

// Implementw sdk.Msg
func (msg MsgBvs) Type() string { return "bvs" }

// Implementw sdk.Msg
func (msg MsgBvs) ValidateBasic() sdk.Error { return nil }

// Implementw sdk.Msg
func (msg MsgBvs) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		return []byte{}
	}
	return b
}

// Implementw sdk.Msg
func (msg MsgBvs) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SenderAccount}
}

// build the sendTx msg
func BuildBvsMsg(senderAccount sdk.AccAddress, sender string, recp string, asset *types.BvsAsset) sdk.Msg {
	return &MsgBvs{
		SenderAccount: senderAccount,
		Sender:        sender,
		Recipient:     recp,
		Asset:         *asset,
	}
}

func NewHandler(mapper types.CodexMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// TODO
		return sdk.Result{}
	}
}
