package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

var (
	delPk1       = ed25519.GenPrivKey().PubKey()
	delPk2       = ed25519.GenPrivKey().PubKey()
	delAddr1     = sdk.AccAddress(delPk1.Address())
	delAddr2     = sdk.AccAddress(delPk2.Address())
	emptyDelAddr sdk.AccAddress

	valPk1       = ed25519.GenPrivKey().PubKey()
	valAddr1     = sdk.ValAddress(valPk1.Address())
	emptyValAddr sdk.ValAddress
)

// TestNewMsgSetWithdrawAddress test ValidateBasic for NewMsgSetWithdrawAddress
func TestNewMsgSetWithdrawAddress(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(delAddr1, delAddr2)
	bz := ModuleCdc.MustMarshalJSON(msg)
	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, "set_withdraw_address", msg.Type())
	require.Equal(t, []sdk.AccAddress{delAddr1}, msg.GetSigners())
	require.Equal(t, sdk.MustSortJSON(bz), msg.GetSignBytes())
	require.NoError(t, msg.ValidateBasic())
}

// TestNewMsgWithdrawValidatorCommission test ValidateBasic for MsgWithdrawValidatorCommission
func TestNewMsgWithdrawValidatorCommission(t *testing.T) {
	msg := NewMsgWithdrawValidatorCommission(valAddr1)
	bz := ModuleCdc.MustMarshalJSON(msg)
	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, "withdraw_validator_commission", msg.Type())
	require.Equal(t, []sdk.AccAddress{valAddr1.Bytes()}, msg.GetSigners())
	require.Equal(t, sdk.MustSortJSON(bz), msg.GetSignBytes())
	require.NoError(t, msg.ValidateBasic())
}

// TestMsgSetWithdrawAddress test ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddress(t *testing.T) {
	tests := []struct {
		delegatorAddr sdk.AccAddress
		withdrawAddr  sdk.AccAddress
		expectPass    bool
	}{
		{delAddr1, delAddr2, true},
		{delAddr1, delAddr1, true},
		{emptyDelAddr, delAddr1, false},
		{delAddr1, emptyDelAddr, false},
		{emptyDelAddr, emptyDelAddr, false},
	}

	for i, tc := range tests {
		msg := NewMsgSetWithdrawAddress(tc.delegatorAddr, tc.withdrawAddr)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test index: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test index: %v", i)
		}
	}
}

// TestMsgWithdrawValidatorCommission test ValidateBasic for MsgWithdrawValidatorCommission
func TestMsgWithdrawValidatorCommission(t *testing.T) {
	tests := []struct {
		validatorAddr sdk.ValAddress
		expectPass    bool
	}{
		{valAddr1, true},
		{emptyValAddr, false},
	}
	for i, tc := range tests {
		msg := NewMsgWithdrawValidatorCommission(tc.validatorAddr)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test index: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test index: %v", i)
		}
	}
}
