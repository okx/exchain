package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	addr           = "okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas"
	expectedOutput = `Contract Deployment Whitelist:
okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas
okexchain1qj5c07sm6jetjz8f509qtrxgh4psxkv32x0qas`
)

func TestContractDeploymentWhitelist_String(t *testing.T) {
	accAddr, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	whitelist := ContractDeploymentWhitelist{accAddr, accAddr}
	require.Equal(t, expectedOutput, whitelist.String())
}
