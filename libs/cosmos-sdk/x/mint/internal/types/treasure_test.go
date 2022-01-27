package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTreasure(t *testing.T) {
	treasure := NewTreasure(nil, sdk.NewDecWithPrec(1, 2))
	t.Log("1", treasure)
	b, err := ModuleCdc.MarshalBinaryLengthPrefixed(treasure)
	require.NoError(t, err)
	treasure = &Treasure{}
	err = ModuleCdc.UnmarshalBinaryLengthPrefixed(b, treasure)
	require.NoError(t, err)
	t.Log("2", treasure)

	b, err = ModuleCdc.MarshalBinaryBare(treasure)

}
