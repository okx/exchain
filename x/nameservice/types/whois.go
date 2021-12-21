/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/19 8:44 上午
# @File : whois.go
# @Description :
# @Attention :
*/
package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type Whois struct {
	Creator sdk.AccAddress `json:"creator" yaml:"creator"`
	ID      string         `json:"id" yaml:"id"`
	Value   string         `json:"value" yaml:"value"`
	Price   sdk.Coins      `json:"price" yaml:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}
