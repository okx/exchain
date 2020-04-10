package v0_8

import (
	"time"

	"github.com/okex/okchain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	Params struct {
		ListAsset              sdk.Dec       `json:"list_asset"`                // Initial Coin Offering,Initial value:100000tokt
		IssueAsset             sdk.Dec       `json:"issue_asset"`               // Issue token,Initial value:20000tokt
		MintAsset              sdk.Dec       `json:"mint_asset"`                // Mint token,Initial value:2000tokt
		BurnAsset              sdk.Dec       `json:"burn_asset"`                // Burn token,Initial value:10tokt
		Transfer               sdk.Dec       `json:"transfer"`                  // Transfer,Initial value:0.0125tokt
		FreezeAsset            sdk.Dec       `json:"freeze_asset"`              // Freeze,Initial value:0.1tokt
		UnfreezeAsset          sdk.Dec       `json:"unfreeze_asset"`            // Unfreeze,Initial value:0.1tokt
		ListPeriod             time.Duration `json:"list_period"`               // Initial Coin Offering window,Initial value:24hours
		ListProposalMinDeposit sdk.Dec       `json:"list_proposal_min_deposit"` // Initial Coin Offering Min Deposit,Initial value:20000tokt
	}

	TokenPair struct {
		BaseAssetSymbol  string  `json:"baseAssetSymbol"`
		QuoteAssetSymbol string  `json:"quoteAssetSymbol"`
		InitPrice        sdk.Dec `json:"price"`
		MaxPriceDigit    int64   `json:"maxPriceDigit"`
		MaxQuantityDigit int64   `json:"maxSizeDigit"`
		MinQuantity      sdk.Dec `json:"minTradeSize"`
		ID               uint64  `json:"tokenPairId"`
	}

	Token struct {
		Desc           string         `json:"desc"`
		Symbol         string         `json:"symbol"`
		OriginalSymbol string         `json:"originalSymbol"`
		WholeName      string         `json:"wholeName"`
		TotalSupply    int64          `json:"totalSupply"`
		Owner          sdk.AccAddress `json:"owner"`
		Mintable       bool           `json:"mintable"`
	}

	GenesisState struct {
		Params     Params           `json:"params"`
		Info       []Token          `json:"info"`
		LockCoins  []types.AccCoins `json:"lock_coins"`
		TokenPairs []*TokenPair     `json:"token_pairs"`
	}
)
