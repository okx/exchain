package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func GetMarginAccount(name string) (marginAcc sdk.AccAddress) {
	return sdk.AccAddress(crypto.AddressHash([]byte(name)))
}

func GetMarginAccountBySwap(name string) (marginAcc sdk.Address, err error) {

	var buffer bytes.Buffer

	byteAcc, err := sdk.AccAddressFromBech32(name)
	if err != nil {
		return marginAcc, err
	}
	byteSwap := append(byteAcc[10:20], byteAcc[0:10]...)
	hexStr := fmt.Sprintf("%x", byteSwap)

	buffer.WriteString(hexStr)

	marginAcc, err = sdk.AccAddressFromHex(buffer.String())
	return marginAcc, err
}
