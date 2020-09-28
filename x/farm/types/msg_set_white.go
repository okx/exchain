package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


type MsgSetWhite struct {
	PoolName string         `json:"pool_name" yaml:"pool_name"`
	Address sdk.AccAddress         `json:"address" yaml:"address"`
}

func NewMsgSetWhite(poolName string, address sdk.AccAddress) MsgSetWhite {
	return MsgSetWhite{
		PoolName: poolName,
		Address:  address,
	}
}

var _ sdk.Msg = MsgSetWhite{}

func (m MsgSetWhite) Route() string {
	return RouterKey
}

func (m MsgSetWhite) Type() string {
	return "setwhite"
}

func (m MsgSetWhite) ValidateBasic() sdk.Error {
	if len(m.PoolName) > 0  {
		return ErrNilAddress(DefaultCodespace)
	}
	return nil
}

func (m MsgSetWhite) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgSetWhite) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Address}
}
