// nolint
package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// Verify interface at compile time
var _ = &MsgWithdrawDelegatorReward{}

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegatorReward struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) MsgWithdrawDelegatorReward {
	return MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
}

func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }
func (msg MsgWithdrawDelegatorReward) Type() string  { return "withdraw_delegator_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr()
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr()
	}

	return nil
}

// Verify interface at compile time
var _ = &MsgWithdrawDelegatorAllRewards{}

// msg struct for delegation withdraw all rewards from all validator
type MsgWithdrawDelegatorAllRewards struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

func NewMsgWithdrawDelegatorAllRewards(delAddr sdk.AccAddress) MsgWithdrawDelegatorAllRewards {
	return MsgWithdrawDelegatorAllRewards{
		DelegatorAddress: delAddr,
	}
}

func (msg MsgWithdrawDelegatorAllRewards) Route() string { return ModuleName }
func (msg MsgWithdrawDelegatorAllRewards) Type() string  { return "withdraw_delegator_all_rewards" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorAllRewards) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorAllRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorAllRewards) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr()
	}

	return nil
}
