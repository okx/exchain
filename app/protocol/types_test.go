package protocol

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewExportedAccount(t *testing.T) {
	// prepare
	testPubKey := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	testAccAddr := sdk.AccAddress(testPubKey.Address())
	testCoins := sdk.NewCoins(sdk.NewCoin("btc", sdk.NewInt(1024)), sdk.NewCoin("eth", sdk.NewInt(2048)))
	vestAcc := auth.NewContinuousVestingAccount(auth.NewBaseAccount(testAccAddr, testCoins, testPubKey, 2, 4), time.Now().UnixNano(), time.Now().UnixNano()+1024)

	// check the field member
	exportedAccount := NewExportedAccount(vestAcc)
	require.NotNil(t, exportedAccount.StartTime)
	require.NotNil(t, exportedAccount.EndTime)
	require.Equal(t, testAccAddr, exportedAccount.Address)
	require.Equal(t, uint64(2), exportedAccount.AccountNumber)
	require.Equal(t, uint64(4), exportedAccount.Sequence)
	require.Equal(t, 2, len(exportedAccount.Coins))
}
