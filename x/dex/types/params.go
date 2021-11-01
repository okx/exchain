package types

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/params"
)

var (
	keyDexListFee             = []byte("DexListFee")
	keyTransferOwnershipFee   = []byte("TransferOwnershipFee")
	keyRegisterOperatorFee    = []byte("RegisterOperatorFee")
	keyDelistMaxDepositPeriod = []byte("DelistMaxDepositPeriod")
	keyDelistMinDeposit       = []byte("DelistMinDeposit")
	keyDelistVotingPeriod     = []byte("DelistVotingPeriod")
	keyWithdrawPeriod         = []byte("WithdrawPeriod")
	keyOwnershipConfirmWindow = []byte("OwnershipConfirmWindow")
)

// Params defines param object
type Params struct {
	ListFee              sdk.SysCoin `json:"list_fee"`
	TransferOwnershipFee sdk.SysCoin `json:"transfer_ownership_fee"`
	RegisterOperatorFee  sdk.SysCoin `json:"register_operator_fee"`

	//  maximum period for okt holders to deposit on a dex delist proposal
	DelistMaxDepositPeriod time.Duration `json:"delist_max_deposit_period"`
	//  minimum deposit for a critical dex delist proposal to enter voting period
	DelistMinDeposit sdk.SysCoins `json:"delist_min_deposit"`
	//  length of the critical voting period for dex delist proposal
	DelistVotingPeriod time.Duration `json:"delist_voting_period"`

	WithdrawPeriod         time.Duration `json:"withdraw_period"`
	OwnershipConfirmWindow time.Duration `json:"ownership_confirm_window"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: keyDexListFee, Value: &p.ListFee, ValidatorFn: common.ValidateSysCoin("list fee")},
		{Key: keyTransferOwnershipFee, Value: &p.TransferOwnershipFee, ValidatorFn: common.ValidateSysCoin("transfer ownership fee")},
		{Key: keyRegisterOperatorFee, Value: &p.RegisterOperatorFee, ValidatorFn: common.ValidateSysCoin("register operator fee")},
		{Key: keyDelistMaxDepositPeriod, Value: &p.DelistMaxDepositPeriod, ValidatorFn: common.ValidateDurationPositive("delist max deposit period")},
		{Key: keyDelistMinDeposit, Value: &p.DelistMinDeposit, ValidatorFn: common.ValidateSysCoins("delist min deposit")},
		{Key: keyDelistVotingPeriod, Value: &p.DelistVotingPeriod, ValidatorFn: common.ValidateDurationPositive("delist voting period")},
		{Key: keyWithdrawPeriod, Value: &p.WithdrawPeriod, ValidatorFn: common.ValidateDurationPositive("withdraw period")},
		{Key: keyOwnershipConfirmWindow, Value: &p.OwnershipConfirmWindow, ValidatorFn: common.ValidateDurationPositive("ownership confirm window")},
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() *Params {
	defaultListFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultFeeList))
	defaultTransferOwnershipFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultFeeTransferOwnership))
	defaultDelistMinDeposit := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultDelistMinDeposit))
	return &Params{
		ListFee:                defaultListFee,
		TransferOwnershipFee:   defaultTransferOwnershipFee,
		RegisterOperatorFee:    sdk.NewDecCoinFromDec(common.NativeToken, sdk.ZeroDec()),
		DelistMaxDepositPeriod: time.Hour * 24,
		DelistMinDeposit:       sdk.SysCoins{defaultDelistMinDeposit},
		DelistVotingPeriod:     time.Hour * 72,
		WithdrawPeriod:         DefaultWithdrawPeriod,
		OwnershipConfirmWindow: DefaultOwnershipConfirmWindow,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	return fmt.Sprintf("Params: \nDexListFee:%s\nTransferOwnershipFee:%s\nRegisterOperatorFee:%s\nDelistMaxDepositPeriod:%s\n"+
		"DelistMinDeposit:%s\nDelistVotingPeriod:%s\nWithdrawPeriod:%d\nOwnershipConfirmWindow: %s\n",
		p.ListFee, p.TransferOwnershipFee, p.RegisterOperatorFee, p.DelistMaxDepositPeriod, p.DelistMinDeposit, p.DelistVotingPeriod, p.WithdrawPeriod, p.OwnershipConfirmWindow)
}
