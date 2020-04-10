package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var (
	coinPos  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	coinZero = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
)

// test ValidateBasic for MsgCreateValidator
func TestMsgCreateValidator(t *testing.T) {
	addr1 := valAddr1
	tests := []struct {
		name, moniker, identity, website, details string
		minSelfDelegation                         sdk.Int
		validatorAddr                             sdk.ValAddress
		delegatorAddr							  sdk.AccAddress
		pubkey                                    crypto.PubKey
		bond                                      sdk.Coin
		expectPass                                bool
	}{
		{"empty bond", "a", "b", "c", "d", sdk.OneInt(), addr1, dlgAddr1,pk1, coinZero, true},
		{"zero min self delegation", "a", "b", "c", "d", sdk.ZeroInt(), addr1, dlgAddr1, pk1,  coinPos, false},
		{"basic good", "a", "b", "c", "d", sdk.OneInt(), addr1, dlgAddr1, pk1,  coinPos, true},
		{"partial description", "", "", "c", "", sdk.OneInt(), addr1, dlgAddr1, pk1,  coinPos, true},
		{"empty description", "", "", "", "", sdk.OneInt(), addr1, dlgAddr1, pk1,  coinPos, false},
		{"empty address1", "a", "b", "c", "d", sdk.OneInt(), emptyAddr, dlgAddr1, pk1, coinPos, false},
		{"empty address2", "a", "b", "c", "d", sdk.OneInt(), nil, nil,pk1, coinPos, false},
		{"valAddr dlgAddr not equals", "a", "b", "c", "d", sdk.OneInt(), addr1, dlgAddr2,pk1, coinPos, false},
		{"empty pubkey", "a", "b", "c", "d", sdk.OneInt(), addr1, dlgAddr1, emptyPubkey, coinPos, true},
		//{"negative min self delegation", "a", "b", "c", "d", commission1, sdk.NewInt(-1), addr1, pk1, coinPos, false},
		//{"delegation less than min self delegation", "a", "b", "c", "d", commission1, coinPos.Amount.Add(sdk.OneInt()), addr1, pk1, coinPos, false},
	}

	for _, tc := range tests {
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.details)
		coin := sdk.NewDecCoin(sdk.DefaultBondDenom, tc.minSelfDelegation)

		msg := MsgCreateValidator{
			Description:       description,
			DelegatorAddress:  tc.delegatorAddr,
			ValidatorAddress:  tc.validatorAddr,
			PubKey:            tc.pubkey,
			MinSelfDelegation: coin,
		}

		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgCreateValidator_Smoke(t *testing.T) {

	msd := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(2000))

	msg := NewMsgCreateValidator(valAddr1, pk1,
		NewDescription("my moniker", "my identity", "my website", "my details"), msd,
	)
	require.Contains(t, msg.Route(), RouterKey)
	require.Contains(t, msg.Type(), "create_validator")
	require.True(t, len(msg.GetSigners()) == 1, msg)
	require.True(t, len(msg.GetSignBytes()) > 0, msg)

	bz, err := msg.MarshalJSON()
	require.Nil(t, err)
	newMsg := MsgCreateValidator{}
	err2 := newMsg.UnmarshalJSON(bz)
	require.Nil(t, err2)
	require.Equal(t, msg, newMsg)

	err3 := newMsg.UnmarshalJSON(nil)
	require.NotNil(t, err3)
}

// test ValidateBasic for MsgEditValidator
func TestMsgEditValidator(t *testing.T) {
	tests := []struct {
		name, moniker, identity, website, details string
		validatorAddr                             sdk.ValAddress
		expectPass                                bool
	}{
		{"basic good", "a", "b", "c", "d", valAddr1, true},
		{"partial description", "", "", "c", "", valAddr1, true},
		{"empty description", "", "", "", "", valAddr1, false},
		{"empty address", "a", "b", "c", "d", emptyAddr, false},
	}

	for _, tc := range tests {
		description := NewDescription(tc.moniker, tc.identity, tc.website, tc.details)
		msg := NewMsgEditValidator(tc.validatorAddr, description)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgEditValidator_Smoke(t *testing.T) {

	desc := NewDescription("my moniker", "my identity", "my website", "my details")
	msg := NewMsgEditValidator(valAddr1, desc)

	require.Contains(t, msg.Route(), RouterKey)
	require.Contains(t, msg.Type(), "edit_validator")
	require.True(t, len(msg.GetSigners()) == 1, msg)
	require.True(t, len(msg.GetSignBytes()) > 0, msg)
}


// test ValidateBasic for MsgDelegate
func TestMsgDelegate(t *testing.T) {

	coinPos := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(1000))
	coinZero := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.ZeroDec())

	tests := []struct {
		name          string
		delegatorAddr sdk.AccAddress
		amount		  sdk.DecCoin
		expectPass    bool
	}{
		{"basic good", dlgAddr1, coinPos, true},
		{"empty delegator", sdk.AccAddress(emptyAddr), coinPos, false},
		{"empty bond", sdk.AccAddress(addr1), coinZero, false},
	}

	for _, tc := range tests {
		msg := NewMsgDelegate(tc.delegatorAddr, tc.amount)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
		}
	}
}

func TestMsgDelegate_Smoke(t *testing.T) {

	coinZero := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.ZeroDec())
	msg := NewMsgDelegate(dlgAddr1, coinZero)

	require.Contains(t, msg.Route(), RouterKey)
	require.Contains(t, msg.Type(), "delegate")
	require.True(t, len(msg.GetSigners()) == 1, msg)
	require.True(t, len(msg.GetSignBytes()) > 0, msg)
}

//// test ValidateBasic for MsgUnbond
//func TestMsgBeginRedelegate(t *testing.T) {
//	tests := []struct {
//		name             string
//		delegatorAddr    sdk.AccAddress
//		validatorSrcAddr sdk.ValAddress
//		validatorDstAddr sdk.ValAddress
//		amount           sdk.Coin
//		expectPass       bool
//	}{
//		{"regular", sdk.AccAddress(addr1), addr2, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
//		{"zero amount", sdk.AccAddress(addr1), addr2, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
//		{"empty delegator", sdk.AccAddress(emptyAddr), addr1, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
//		{"empty source validator", sdk.AccAddress(addr1), emptyAddr, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
//		{"empty destination validator", sdk.AccAddress(addr1), addr2, emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
//	}
//
//	for _, tc := range tests {
//		msg := NewMsgBeginRedelegate(tc.delegatorAddr, tc.validatorSrcAddr, tc.validatorDstAddr, tc.amount)
//		if tc.expectPass {
//			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
//		} else {
//			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
//		}
//	}
//}
//
//// test ValidateBasic for MsgUnbond
//func TestMsgUndelegate(t *testing.T) {
//	tests := []struct {
//		name          string
//		delegatorAddr sdk.AccAddress
//		validatorAddr sdk.ValAddress
//		amount        sdk.Coin
//		expectPass    bool
//	}{
//		{"regular", sdk.AccAddress(addr1), addr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), true},
//		{"zero amount", sdk.AccAddress(addr1), addr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), false},
//		{"empty delegator", sdk.AccAddress(emptyAddr), addr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
//		{"empty validator", sdk.AccAddress(addr1), emptyAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1), false},
//	}
//
//	for _, tc := range tests {
//		msg := NewMsgUndelegate(tc.delegatorAddr, tc.validatorAddr, tc.amount)
//		if tc.expectPass {
//			require.Nil(t, msg.ValidateBasic(), "test: %v", tc.name)
//		} else {
//			require.NotNil(t, msg.ValidateBasic(), "test: %v", tc.name)
//		}
//	}
//}
