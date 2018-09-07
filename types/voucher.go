package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// A Voucher is an asset representing a value guaranteed by a voucher issuer.
type Voucher struct {
	Id       string `json:"id"`
	Origin   string `json:"origin"` // may be a Codex or a Dealer
	Holder   string `json:"holder"`
	ExpireOn int    `json:"expire-on"`
}

type BvsAsset struct {
	Silver   sdk.Int   `json:"silver"`
	Gold     sdk.Int   `json:"gold"`
	Vouchers []string  `json:"vouchers"`
	Coins    sdk.Coins // advisory member to work with cosmos-sdk
}

func ParseBvsAsset(str string) (asset *BvsAsset) {
	// TODO: careful about overflow
	asset = &BvsAsset{}
	var silver int64 = 0
	var gold int64 = 0
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return nil
	}
	chunks := strings.Split(str, ",")
	for _, chunk := range chunks {
		if strings.Contains(chunk, "bvs") {
			coin, err := sdk.ParseCoin(chunk)
			if err != nil {
				continue
			}
			silver += coin.Amount.Int64()
		} else if strings.Contains(chunk, "bvg") {
			coin, err := sdk.ParseCoin(chunk)
			if err != nil {
				continue
			}
			gold += coin.Amount.Int64()
		} else {
			asset.Vouchers = append(asset.Vouchers, chunk)
		}
	}

	if silver > 0 {
		asset.Silver.AddRaw(silver)
		asset.Coins = append(asset.Coins, sdk.NewCoin("bvs", asset.Silver))
	}
	if gold > 0 {
		asset.Gold.AddRaw(silver)
		asset.Coins = append(asset.Coins, sdk.NewCoin("bvg", asset.Silver))
	}

	return
}

func IsOwner(owner string, asset *BvsAsset) bool {
	return true
}
