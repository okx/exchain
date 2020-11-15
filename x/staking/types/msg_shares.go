package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = (*MsgAddShares)(nil)
	_ sdk.Msg = (*MsgDestroyValidator)(nil)
)

// MsgDestroyValidator - struct for transactions to deregister a validator
type MsgDestroyValidator struct {
	DelAddr sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// NewMsgDestroyValidator creates a msg of destroy-validator
func NewMsgDestroyValidator(delAddr sdk.AccAddress) MsgDestroyValidator {
	return MsgDestroyValidator{
		DelAddr: delAddr,
	}
}

// nolint
func (MsgDestroyValidator) Route() string { return RouterKey }
func (MsgDestroyValidator) Type() string  { return "destroy_validator" }
func (msg MsgDestroyValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelAddr}
}

// ValidateBasic gives a quick validity check
func (msg MsgDestroyValidator) ValidateBasic() error {
	if msg.DelAddr.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}

	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgDestroyValidator) GetSignBytes() []byte {
	bytes := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bytes)
}

// MsgUnbindProxy - structure for unbinding proxy relationship between the delegator and the proxy
type MsgUnbindProxy struct {
	DelAddr sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// NewMsgUnbindProxy creates a msg of unbinding proxy
func NewMsgUnbindProxy(delAddr sdk.AccAddress) MsgUnbindProxy {
	return MsgUnbindProxy{
		DelAddr: delAddr,
	}
}

// nolint
func (MsgUnbindProxy) Route() string { return RouterKey }
func (MsgUnbindProxy) Type() string  { return "unbind_proxy" }
func (msg MsgUnbindProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelAddr}
}

// ValidateBasic gives a quick validity check
func (msg MsgUnbindProxy) ValidateBasic() error {
	if msg.DelAddr.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgUnbindProxy) GetSignBytes() []byte {
	bytes := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bytes)
}

// MsgRegProxy - register delegator as proxy or unregister proxy to delegator
// if Reg == true, action is reg, otherwise action is unreg
type MsgRegProxy struct {
	ProxyAddress sdk.AccAddress `json:"proxy_address" yaml:"proxy_address"`
	Reg          bool           `json:"reg" yaml:"reg"`
}

// NewMsgRegProxy creates a msg of registering proxy
func NewMsgRegProxy(proxyAddress sdk.AccAddress, reg bool) MsgRegProxy {
	return MsgRegProxy{
		ProxyAddress: proxyAddress,
		Reg:          reg,
	}
}

// nolint
func (MsgRegProxy) Route() string { return RouterKey }
func (MsgRegProxy) Type() string  { return "reg_or_unreg_proxy" }
func (msg MsgRegProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ProxyAddress}
}

// ValidateBasic gives a quick validity check
func (msg MsgRegProxy) ValidateBasic() error {
	if msg.ProxyAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgRegProxy) GetSignBytes() []byte {
	bytes := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bytes)
}

// MsgBindProxy - structure for bind proxy relationship between the delegator and the proxy
type MsgBindProxy struct {
	DelAddr      sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ProxyAddress sdk.AccAddress `json:"proxy_address" yaml:"proxy_address"`
}

// NewMsgBindProxy creates a msg of binding proxy
func NewMsgBindProxy(delAddr sdk.AccAddress, ProxyDelAddr sdk.AccAddress) MsgBindProxy {
	return MsgBindProxy{
		DelAddr:      delAddr,
		ProxyAddress: ProxyDelAddr,
	}
}

// nolint
func (MsgBindProxy) Route() string { return RouterKey }
func (MsgBindProxy) Type() string  { return "bind_proxy" }
func (msg MsgBindProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelAddr}
}

// ValidateBasic gives a quick validity check
func (msg MsgBindProxy) ValidateBasic() error {
	if msg.DelAddr.Empty() || msg.ProxyAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}

	if msg.DelAddr.Equals(msg.ProxyAddress) {
		return ErrWrongOperationAddr(DefaultCodespace,
			fmt.Sprintf("ProxyAddress: %s eqauls to DelegatorAddress: %s",
				msg.ProxyAddress.String(), msg.DelAddr.String()))
	}

	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgBindProxy) GetSignBytes() []byte {
	bytes := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bytes)
}

// MsgAddShares - struct for adding-shares transaction
type MsgAddShares struct {
	DelAddr  sdk.AccAddress   `json:"delegator_address" yaml:"delegator_address"`
	ValAddrs []sdk.ValAddress `json:"validator_addresses" yaml:"validator_addresses"`
}

// NewMsgAddShares creates a msg of adding shares to vals
func NewMsgAddShares(delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) MsgAddShares {
	return MsgAddShares{
		DelAddr:  delAddr,
		ValAddrs: valAddrs,
	}
}

// nolint
func (MsgAddShares) Route() string { return RouterKey }
func (MsgAddShares) Type() string  { return "add_shares_to_validators" }
func (msg MsgAddShares) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelAddr}
}

// ValidateBasic gives a quick validity check
func (msg MsgAddShares) ValidateBasic() error {
	if msg.DelAddr.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}

	if msg.ValAddrs == nil || len(msg.ValAddrs) == 0 {
		return ErrWrongOperationAddr(DefaultCodespace, "ValAddrs is empty")
	}

	if isValsDuplicate(msg.ValAddrs) {
		return ErrTargetValsDuplicate(DefaultCodespace)
	}

	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgAddShares) GetSignBytes() []byte {
	bytes := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bytes)
}

func isValsDuplicate(valAddrs []sdk.ValAddress) bool {
	lenAddrs := len(valAddrs)
	filter := make(map[string]struct{}, lenAddrs)
	for i := 0; i < lenAddrs; i++ {
		key := valAddrs[i].String()
		if _, ok := filter[key]; ok {
			return true
		}
		filter[key] = struct{}{}
	}

	return false
}

// MsgDeposit - structure for depositing to the delegator account
type MsgDeposit struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	Amount           sdk.SysCoin    `json:"quantity" yaml:"quantity"`
}

// NewMsgDeposit creates a new instance of MsgDeposit
func NewMsgDeposit(delAddr sdk.AccAddress, amount sdk.SysCoin) MsgDeposit {
	return MsgDeposit{
		DelegatorAddress: delAddr,
		Amount:           amount,
	}
}

// nolint
func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return "deposit" }
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// ValidateBasic gives a quick validity check
func (msg MsgDeposit) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroDec()) || !msg.Amount.IsValid() {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgWithdraw - structure for withdrawing okt and the corresponding shares from all validators
type MsgWithdraw struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	Amount           sdk.SysCoin    `json:"quantity" yaml:"quantity"`
}

// NewMsgWithdraw creates a new instance of MsgWithdraw
func NewMsgWithdraw(delAddr sdk.AccAddress, amount sdk.SysCoin) MsgWithdraw {
	return MsgWithdraw{
		DelegatorAddress: delAddr,
		Amount:           amount,
	}
}

// nolint
func (msg MsgWithdraw) Route() string { return RouterKey }
func (msg MsgWithdraw) Type() string  { return "withdraw" }
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// ValidateBasic gives a quick validity check
func (msg MsgWithdraw) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if !msg.Amount.IsValid() {
		return ErrBadUnDelegationAmount(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
