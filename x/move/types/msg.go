//nolint
package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// Verify interface at compile time
var _ sdk.Msg = &MsgPublishMove{}

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgPublishMove struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	MovePath         string         `json:"move_path" yaml:"move_path"`
}

func NewMsgPublishMove(delAddr sdk.AccAddress, movePath string) MsgPublishMove {
	return MsgPublishMove{
		DelegatorAddress: delAddr,
		MovePath:         movePath,
	}
}

func (msg MsgPublishMove) Route() string { return ModuleName }
func (msg MsgPublishMove) Type() string  { return "publish_move" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgPublishMove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgPublishMove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgPublishMove) ValidateBasic() sdk.Error {
	return nil
}
