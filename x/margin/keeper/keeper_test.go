package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAccountInvariant(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("okchain123")
	fmt.Print(addr)
	require.Nil(t, err)
}
