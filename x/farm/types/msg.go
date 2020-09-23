package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type MsgCreatePool struct {

}

var _ sdk.Msg = MsgCreatePool{}

func (m MsgCreatePool) Route() string {
	panic("implement me")
}

func (m MsgCreatePool) Type() string {
	panic("implement me")
}

func (m MsgCreatePool) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgCreatePool) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

