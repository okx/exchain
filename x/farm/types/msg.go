package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

const (
	MaxPoolNameLength = 128

	createPoolMsgType  = "create_pool"
	destroyPoolMsgType = "destroy_pool"
	provideMsgType     = "provide"
	lockMsgType        = "lock"
	unlockMsgType      = "unlock"
	claimMsgType       = "claim"
)

type MsgCreatePool struct {
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
	PoolName      string         `json:"pool_name" yaml:"pool_name"`
	MinLockAmount sdk.SysCoin    `json:"min_lock_amount" yaml:"min_lock_amount"`
	YieldedSymbol string         `json:"yielded_symbol"  yaml:"yielded_symbol"`
}

var _ sdk.Msg = MsgCreatePool{}

func NewMsgCreatePool(address sdk.AccAddress, poolName string, minLockAmount sdk.SysCoin, yieldedSymbol string) MsgCreatePool {
	return MsgCreatePool{
		Owner:         address,
		PoolName:      poolName,
		MinLockAmount: minLockAmount,
		YieldedSymbol: yieldedSymbol,
	}
}

func (m MsgCreatePool) Route() string {
	return RouterKey
}

func (m MsgCreatePool) Type() string {
	return createPoolMsgType
}

func (m MsgCreatePool) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return ErrNilAddress()
	}
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrPoolNameLength(m.PoolName, len(m.PoolName), MaxPoolNameLength)
	}
	if m.MinLockAmount.Amount.LT(sdk.ZeroDec()) || !m.MinLockAmount.IsValid() {
		return ErrInvalidInputAmount(m.MinLockAmount.String())
	}
	if m.YieldedSymbol == "" {
		return ErrInvalidInput("yielded symbol is empty")
	}
	return nil
}

func (m MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

type MsgDestroyPool struct {
	Owner    sdk.AccAddress `json:"owner" yaml:"owner"`
	PoolName string         `json:"pool_name" yaml:"pool_name"`
}

var _ sdk.Msg = MsgDestroyPool{}

func NewMsgDestroyPool(address sdk.AccAddress, poolName string) MsgDestroyPool {
	return MsgDestroyPool{
		Owner:    address,
		PoolName: poolName,
	}
}

func (m MsgDestroyPool) Route() string {
	return RouterKey
}

func (m MsgDestroyPool) Type() string {
	return destroyPoolMsgType
}

func (m MsgDestroyPool) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return ErrNilAddress()
	}
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrPoolNameLength(m.PoolName, len(m.PoolName), MaxPoolNameLength)
	}
	return nil
}

func (m MsgDestroyPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgDestroyPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

type MsgProvide struct {
	PoolName              string         `json:"pool_name" yaml:"pool_name"`
	Address               sdk.AccAddress `json:"address" yaml:"address"`
	Amount                sdk.SysCoin    `json:"amount" yaml:"amount"`
	AmountYieldedPerBlock sdk.Dec        `json:"amount_yielded_per_block" yaml:"amount_yielded_per_block"`
	StartHeightToYield    int64          `json:"start_height_to_yield" yaml:"start_height_to_yield"`
}

func NewMsgProvide(poolName string, address sdk.AccAddress, amount sdk.SysCoin,
	amountYieldedPerBlock sdk.Dec, startHeightToYield int64) MsgProvide {
	return MsgProvide{
		PoolName:              poolName,
		Address:               address,
		Amount:                amount,
		AmountYieldedPerBlock: amountYieldedPerBlock,
		StartHeightToYield:    startHeightToYield,
	}
}

var _ sdk.Msg = MsgProvide{}

func (m MsgProvide) Route() string {
	return RouterKey
}

func (m MsgProvide) Type() string {
	return provideMsgType
}

func (m MsgProvide) ValidateBasic() sdk.Error {
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrInvalidInput(m.PoolName)
	}
	if m.Address.Empty() {
		return ErrNilAddress()
	}
	if m.Amount.Amount.LTE(sdk.ZeroDec()) || !m.Amount.IsValid() {
		return ErrInvalidInputAmount(m.Amount.String())
	}
	if m.AmountYieldedPerBlock.LTE(sdk.ZeroDec()) {
		return ErrInvalidInput("amount yielded per block must be > 0")
	}
	if m.Amount.Amount.LT(m.AmountYieldedPerBlock) {
		return ErrInvalidInput("provided amount must be bigger than amount_yielded_per_block")
	}
	if m.StartHeightToYield <= 0 {
		return ErrInvalidInput("start height to yield must be > 0")
	}
	return nil
}

func (m MsgProvide) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgProvide) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Address}
}

type MsgLock struct {
	PoolName string         `json:"pool_name" yaml:"pool_name"`
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Amount   sdk.SysCoin    `json:"amount" yaml:"amount"`
}

func NewMsgLock(poolName string, address sdk.AccAddress, amount sdk.SysCoin) MsgLock {
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
	return lockMsgType
}

func (m MsgLock) ValidateBasic() sdk.Error {
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrInvalidInput(m.PoolName)
	}
	if m.Address.Empty() {
		return ErrNilAddress()
	}
	if m.Amount.Amount.LTE(sdk.ZeroDec()) || !m.Amount.IsValid() {
		return ErrInvalidInputAmount(m.Amount.Amount.String())
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
	Amount   sdk.SysCoin    `json:"amount" yaml:"amount"`
}

func NewMsgUnlock(poolName string, address sdk.AccAddress, amount sdk.SysCoin) MsgUnlock {
	return MsgUnlock{
		PoolName: poolName,
		Address:  address,
		Amount:   amount,
	}
}

var _ sdk.Msg = MsgUnlock{}

func (m MsgUnlock) Route() string {
	return RouterKey
}

func (m MsgUnlock) Type() string {
	return unlockMsgType
}

func (m MsgUnlock) ValidateBasic() sdk.Error {
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrInvalidInput(m.PoolName)
	}
	if m.Address.Empty() {
		return ErrNilAddress()
	}
	if m.Amount.Amount.LTE(sdk.ZeroDec()) || !m.Amount.IsValid() {
		return ErrInvalidInputAmount(m.Amount.Amount.String())
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

type MsgClaim struct {
	PoolName string         `json:"pool_name" yaml:"pool_name"`
	Address  sdk.AccAddress `json:"address" yaml:"address"`
}

func NewMsgClaim(poolName string, address sdk.AccAddress) MsgClaim {
	return MsgClaim{
		PoolName: poolName,
		Address:  address,
	}
}

var _ sdk.Msg = MsgClaim{}

func (m MsgClaim) Route() string {
	return RouterKey
}

func (m MsgClaim) Type() string {
	return claimMsgType
}

func (m MsgClaim) ValidateBasic() sdk.Error {
	if m.PoolName == "" || len(m.PoolName) > MaxPoolNameLength {
		return ErrInvalidInput(m.PoolName)
	}
	if m.Address.Empty() {
		return ErrNilAddress()
	}
	return nil
}

func (m MsgClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Address}
}
