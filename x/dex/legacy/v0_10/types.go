package v0_10

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "dex"

type (
	// GenesisState - all slashing state that must be provided at genesis
	GenesisState struct {
		Params        Params         `json:"params"`
		TokenPairs    []*TokenPair   `json:"token_pairs"`
		WithdrawInfos WithdrawInfos  `json:"withdraw_infos"`
		ProductLocks  ProductLockMap `json:"product_locks"`
	}

	// Params defines param object
	Params struct {
		ListFee              sdk.DecCoin `json:"list_fee"`
		TransferOwnershipFee sdk.DecCoin `json:"transfer_ownership_fee"`

		//  maximum period for okt holders to deposit on a dex delist proposal
		DelistMaxDepositPeriod time.Duration `json:"delist_max_deposit_period"`
		//  minimum deposit for a critical dex delist proposal to enter voting period
		DelistMinDeposit sdk.DecCoins `json:"delist_min_deposit"`
		//  length of the critical voting period for dex delist proposal
		DelistVotingPeriod time.Duration `json:"delist_voting_period"`

		WithdrawPeriod time.Duration `json:"withdraw_period"`
	}

	// TokenPair represents token pair object
	TokenPair struct {
		BaseAssetSymbol  string         `json:"base_asset_symbol"`
		QuoteAssetSymbol string         `json:"quote_asset_symbol"`
		InitPrice        sdk.Dec        `json:"price"`
		MaxPriceDigit    int64          `json:"max_price_digit"`
		MaxQuantityDigit int64          `json:"max_size_digit"`
		MinQuantity      sdk.Dec        `json:"min_trade_size"`
		ID               uint64         `json:"token_pair_id"`
		Delisting        bool           `json:"delisting"`
		Owner            sdk.AccAddress `json:"owner"`
		Deposits         sdk.DecCoin    `json:"deposits"`
		BlockHeight      int64          `json:"block_height"`
	}

	// WithdrawInfo represents infos for withdrawing
	WithdrawInfo struct {
		Owner        sdk.AccAddress `json:"owner"`
		Deposits     sdk.DecCoin    `json:"deposits"`
		CompleteTime time.Time      `json:"complete_time"`
	}

	// WithdrawInfos defines list of WithdrawInfo
	WithdrawInfos []WithdrawInfo

	ProductLock struct {
		BlockHeight  int64
		Price        sdk.Dec
		Quantity     sdk.Dec
		BuyExecuted  sdk.Dec
		SellExecuted sdk.Dec
	}

	ProductLockMap struct {
		Data map[string]*ProductLock
	}
)
