package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	addr           = "okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas"
	expectedOutput = `Address List:
okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas
okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas`
)

func TestAddressList_String(t *testing.T) {
	accAddr, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	addrList := AddressList{accAddr, accAddr}
	require.Equal(t, expectedOutput, addrList.String())
}
