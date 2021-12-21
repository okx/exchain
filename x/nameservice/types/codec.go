/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/19 8:56 上午
# @File : codec.go
# @Description :
# @Attention :
*/
package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// this line is used by starport scaffolding # 1
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgDeleteName{}, "nameservice/DeleteName", nil)
}
