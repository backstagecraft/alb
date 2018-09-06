package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

//////////////////////////////////////////////////////////////////
// CodexMapper

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

func Id2StoreKey(prefix string, id string) []byte {
	return append([]byte(prefix), []byte(id)...)
}

func (cm CodexMapper) GetCodex(ctx sdk.Context, id string) *Codex {
	store := ctx.KVStore(cm.key)
	bz := store.Get(Id2StoreKey("codex:", id))
	if bz == nil {
		return nil
	}
	cod := cm.decodeCodex(bz)
	return cod
}

func (cm CodexMapper) SetCodex(ctx sdk.Context, cod *Codex) {
	store := ctx.KVStore(cm.key)
	bz := cm.encodeCodex(cod)
	store.Set(Id2StoreKey("codex:", cod.Id), bz)
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

//////////////////////////////////////////////////////////////////
// VoucherMapper

type VoucherMapper struct {
	key   sdk.StoreKey
	proto func() *Voucher
	cdc   *wire.Codec
}

func NewVoucherMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() *Voucher) VoucherMapper {
	return VoucherMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

func (vm VoucherMapper) GetVoucher(ctx sdk.Context, id string) *Voucher {
	store := ctx.KVStore(vm.key)
	bz := store.Get(Id2StoreKey("voucher:", id))
	if bz == nil {
		return nil
	}
	voucher := vm.decodeVoucher(bz)
	return voucher
}

func (vm VoucherMapper) SetVoucher(ctx sdk.Context, voucher *Voucher) {
	store := ctx.KVStore(vm.key)
	bz := vm.encodeVoucher(voucher)
	store.Set(Id2StoreKey("voucher:", voucher.Id), bz)
}

func (vm VoucherMapper) IterateVouchers(ctx sdk.Context, process func(*Voucher) (stop bool)) {
	store := ctx.KVStore(vm.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("voucher:"))
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		voucher := vm.decodeVoucher(val)
		if process(voucher) {
			return
		}
		iter.Next()
	}
}

func (vm VoucherMapper) encodeVoucher(voucher *Voucher) []byte {
	bz, err := vm.cdc.MarshalBinaryBare(voucher)
	if err != nil {
		panic(err)
	}
	return bz
}

func (vm VoucherMapper) decodeVoucher(bz []byte) (voucher *Voucher) {
	voucher = &Voucher{}
	err := vm.cdc.UnmarshalBinaryBare(bz, voucher)
	if err != nil {
		panic(err)
	}
	return
}
