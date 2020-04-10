package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewMsgSetWithdrawAddress(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(DelAddr1, DelAddr2)
	bz := ModuleCdc.MustMarshalJSON(msg)
	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, "set_withdraw_address", msg.Type())
	require.Equal(t, []sdk.AccAddress{DelAddr1}, msg.GetSigners())
	require.Equal(t, sdk.MustSortJSON(bz), msg.GetSignBytes())
	require.NoError(t, msg.ValidateBasic())
}

func TestNewMsgWithdrawValidatorCommission(t *testing.T) {
	msg := NewMsgWithdrawValidatorCommission(ValAddr1)
	bz := ModuleCdc.MustMarshalJSON(msg)
	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, "withdraw_validator_commission", msg.Type())
	require.Equal(t, []sdk.AccAddress{ValAddr1.Bytes()}, msg.GetSigners())
	require.Equal(t, sdk.MustSortJSON(bz), msg.GetSignBytes())
	require.NoError(t, msg.ValidateBasic())
}

// test ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddress(t *testing.T) {
	tests := []struct {
		delegatorAddr sdk.AccAddress
		withdrawAddr  sdk.AccAddress
		expectPass    bool
	}{
		{DelAddr1, DelAddr2, true},
		{DelAddr1, DelAddr1, true},
		{EmptyDelAddr, DelAddr1, false},
		{DelAddr1, EmptyDelAddr, false},
		{EmptyDelAddr, EmptyDelAddr, false},
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

// test ValidateBasic for MsgWithdrawValidatorCommission
func TestMsgWithdrawValidatorCommission(t *testing.T) {
	tests := []struct {
		validatorAddr sdk.ValAddress
		expectPass    bool
	}{
		{ValAddr1, true},
		{EmptyValAddr, false},
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
