package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

type MsgProvide struct {

}

var _ sdk.Msg = MsgProvide{}

func (m MsgProvide) Route() string {
	panic("implement me")
}

func (m MsgProvide) Type() string {
	panic("implement me")
}

func (m MsgProvide) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgProvide) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgProvide) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

type MsgClaim struct {

}

var _ sdk.Msg = MsgClaim{}

func (m MsgClaim) Route() string {
	panic("implement me")
}

func (m MsgClaim) Type() string {
	panic("implement me")
}

func (m MsgClaim) ValidateBasic() sdk.Error {
	panic("implement me")
}

func (m MsgClaim) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgClaim) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

type MsgLock struct {
	PoolName string         `json:"pool_name" yaml:"pool_name"`
	Address  sdk.AccAddress `json:"address" yaml:address"`
	Amount   sdk.DecCoin    `json:"amount" yaml:"amount"`
}

func NewMsgLock(poolName string, address sdk.AccAddress, amount sdk.DecCoin) MsgLock {
	return MsgLock{
		PoolName: poolName,
		Address:  address,
		Amount:   amount,
	}
}

var _ sdk.Msg = MsgLock{}

func (m MsgLock) Route() string {
	return RouterKey
}

func (m MsgLock) Type() string {
	return "lock"
}

func (m MsgLock) ValidateBasic() sdk.Error {
	if m.Amount.Amount.LTE(sdk.ZeroDec()) || !m.Amount.IsValid() {
		return nil
	}
	return nil
}

func (m MsgLock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Address}
}

type MsgUnlock struct {
	PoolName string         `json:"pool_name" yaml:"pool_name"`
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Amount   sdk.DecCoin    `json:"amount" yaml:"amount"`
}

func NewMsgUnlock(poolName string, Address sdk.AccAddress, amount sdk.DecCoin) MsgUnlock {
	return MsgUnlock{
		PoolName: poolName,
		Address:  Address,
		Amount:   amount,
	}
}

var _ sdk.Msg = MsgUnlock{}

func (m MsgUnlock) Route() string {
	return RouterKey
}

func (m MsgUnlock) Type() string {
	return "unlock"
}

func (m MsgUnlock) ValidateBasic() sdk.Error {
	if m.Amount.Amount.LTE(sdk.ZeroDec()) || !m.Amount.IsValid() {
		return nil
	}
	return nil
}

func (m MsgUnlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUnlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Address}
}
