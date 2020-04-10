package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// just 4 test
/////////////////////////////////////////////////////////////
// MsgCreateValidator - struct for bonding transactions
type MsgUpgradeConfig struct {
	ProposalID uint64         `json:"proposal_id"`
	Version    uint64         `json:"version"`
	Height     uint64         `json:"height"`
	Software   string         `json:"software"`
	Owner      sdk.AccAddress `json:"owner"`
}

// Default way to create validator. Delegator address and validator address are the same
func NewMsgUpgradeConfig(proposalID, version, height uint64, software string, owner sdk.AccAddress) MsgUpgradeConfig {

	return MsgUpgradeConfig{
		ProposalID: proposalID,
		Version:    version,
		Height:     height,
		Software:   software,
		Owner:      owner,
	}
}

// nolint
func (msg MsgUpgradeConfig) Route() string { return RouterKey }
func (msg MsgUpgradeConfig) Type() string  { return "upgrade_config" }

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgUpgradeConfig) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgUpgradeConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUpgradeConfig) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}
