package v0_11

import (
	"time"

	"github.com/okex/exchain/x/dex/legacy/v0_10"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	ModuleName = "dex"
)

type (
	// GenesisState - all dex state that must be provided at genesis
	GenesisState struct {
		Params         Params               `json:"params"`
		TokenPairs     []*v0_10.TokenPair   `json:"token_pairs"`
		WithdrawInfos  v0_10.WithdrawInfos  `json:"withdraw_infos"`
		ProductLocks   v0_10.ProductLockMap `json:"product_locks"`
		Operators      DEXOperators         `json:"operators"`
		MaxTokenPairID uint64               `json:"max_token_pair_id" yaml:"max_token_pair_id"`
	}

	// Params defines param object
	Params struct {
		ListFee              sdk.SysCoin `json:"list_fee"`
		TransferOwnershipFee sdk.SysCoin `json:"transfer_ownership_fee"`
		RegisterOperatorFee  sdk.SysCoin `json:"register_operator_fee"`

		//  maximum period for okt holders to deposit on a dex delist proposal
		DelistMaxDepositPeriod time.Duration `json:"delist_max_deposit_period"`
		//  minimum deposit for a critical dex delist proposal to enter voting period
		DelistMinDeposit sdk.SysCoins `json:"delist_min_deposit"`
		//  length of the critical voting period for dex delist proposal
		DelistVotingPeriod time.Duration `json:"delist_voting_period"`

		WithdrawPeriod time.Duration `json:"withdraw_period"`
	}

	// OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"`
	DEXOperator struct {
		Address            sdk.AccAddress `json:"address"`
		HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
		Website            string         `json:"website"`
		InitHeight         int64          `json:"init_height"`
		TxHash             string         `json:"tx_hash"`
	}

	// nolint
	DEXOperators []DEXOperator
)
