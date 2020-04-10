package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/params"
)

var (
	KeyDexListFee             = []byte("DexListFee")
	KeyDexDelistFee           = []byte("DexDeListFee")
	KeyTransferOwnershipFee   = []byte("TransferOwnershipFee")
	KeyDelistMaxDepositPeriod = []byte("DelistMaxDepositPeriod")
	KeyDelistMinDeposit       = []byte("DelistMinDeposit")
	KeyDelistVotingPeriod     = []byte("DelistVotingPeriod")
	KeyWithdrawPeriod         = []byte("WithdrawPeriod")
)

type Params struct {
	ListFee              sdk.DecCoin `json:"list_fee"`
	TransferOwnershipFee sdk.DecCoin `json:"transfer_ownership_fee"`
	//DelistFee            sdk.DecCoins `json:"delist_fee"`

	//  maximum period for okt holders to deposit on a dex delist proposal
	DelistMaxDepositPeriod time.Duration `json:"delist_max_deposit_period"`
	//  minimum deposit for a critical dex delist proposal to enter voting period
	DelistMinDeposit sdk.DecCoins `json:"delist_min_deposit"`
	//  length of the critical voting period for dex delist proposal
	DelistVotingPeriod time.Duration `json:"delist_voting_period"`

	WithdrawPeriod time.Duration `json:"withdraw_period"`
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyDexListFee, Value: &p.ListFee},
		{Key: KeyTransferOwnershipFee, Value: &p.TransferOwnershipFee},
		{Key: KeyDelistMaxDepositPeriod, Value: &p.DelistMaxDepositPeriod},
		{Key: KeyDelistMinDeposit, Value: &p.DelistMinDeposit},
		{Key: KeyDelistVotingPeriod, Value: &p.DelistVotingPeriod},
		{Key: KeyWithdrawPeriod, Value: &p.WithdrawPeriod},
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() *Params {
	var defaultListFee = sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeList))
	//var defaultDeListFee = sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeDelist))
	var defaultTransferOwnershipFee = sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeTransferOwnership))
	var defaultDelistMinDeposit = sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultDelistMinDeposit))
	return &Params{
		ListFee:                defaultListFee,
		TransferOwnershipFee:   defaultTransferOwnershipFee,
		DelistMaxDepositPeriod: time.Hour * 24,
		DelistMinDeposit:       sdk.DecCoins{defaultDelistMinDeposit},
		DelistVotingPeriod:     time.Hour * 72,
		WithdrawPeriod:         DefaultWithdrawPeriod,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("DexListFee:%s\n", p.ListFee))
	//sb.WriteString(fmt.Sprintf("DexDelistFee:%s\n", p.DelistFee))
	sb.WriteString(fmt.Sprintf("TransferOwnershipFee:%s\n", p.TransferOwnershipFee))
	sb.WriteString(fmt.Sprintf("DelistMaxDepositPeriod:%s\n", p.DelistMaxDepositPeriod))
	sb.WriteString(fmt.Sprintf("DelistMinDeposit:%s\n", p.DelistMinDeposit))
	sb.WriteString(fmt.Sprintf("DelistVotingPeriod:%s\n", p.DelistMaxDepositPeriod))
	sb.WriteString(fmt.Sprintf("WithdrawPeriod:%d\n", p.WithdrawPeriod))
	return sb.String()
}
