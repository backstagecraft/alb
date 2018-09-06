package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

type CodexMapper struct {
	key   sdk.StoreKey
	proto func() *Codex
	cdc   *wire.Codec
}

func NewCodexMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() *Codex) CodexMapper {
	return CodexMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

func Id2StoreKey(id string) []byte {
	return append([]byte("codex:"), []byte(id)...)
}

func (cm CodexMapper) GetCodex(ctx sdk.Context, id string) *Codex {
	store := ctx.KVStore(cm.key)
	bz := store.Get(Id2StoreKey(id))
	if bz == nil {
		return nil
	}
	cod := cm.decodeCodex(bz)
	return cod
}

func (cm CodexMapper) SetCodex(ctx sdk.Context, cod *Codex) {
	store := ctx.KVStore(cm.key)
	bz := cm.encodeCodex(cod)
	store.Set(Id2StoreKey(cod.Id), bz)
}

func (cm CodexMapper) IterateCodices(ctx sdk.Context, process func(*Codex) (stop bool)) {
	store := ctx.KVStore(cm.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("codex:"))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		cod := cm.decodeCodex(val)
		if process(cod) {
			return
		}
		iter.Next()
	}
}

func (cm CodexMapper) encodeCodex(cod *Codex) []byte {
	bz, err := cm.cdc.MarshalBinaryBare(cod)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cm CodexMapper) decodeCodex(bz []byte) (cod *Codex) {
	cod = &Codex{}
	err := cm.cdc.UnmarshalBinaryBare(bz, cod)
	if err != nil {
		panic(err)
	}
	return
}
