package types

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	addr           = "ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc"
	expectedOutput = `Address List:
ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc
ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc`
)

func TestAddressList_String(t *testing.T) {
	accAddr, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	addrList := AddressList{accAddr, accAddr}
	require.Equal(t, expectedOutput, addrList.String())
}
