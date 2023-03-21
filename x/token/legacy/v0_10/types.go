package v0_10

import sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

const ModuleName = "token"

type (
	// all state that must be provided in genesis file
	GenesisState struct {
		Params       Params     `json:"params"`
		Tokens       []Token    `json:"tokens"`
		LockedAssets []AccCoins `json:"locked_assets"`
		LockedFees   []AccCoins `json:"locked_fees"`
	}

	Params struct {
		FeeIssue  sdk.SysCoin `json:"issue_fee"`
		FeeMint   sdk.SysCoin `json:"mint_fee"`
		FeeBurn   sdk.SysCoin `json:"burn_fee"`
		FeeModify sdk.SysCoin `json:"modify_fee"`
		FeeChown  sdk.SysCoin `json:"transfer_ownership_fee"`
	}

	Token struct {
		Description         string         `json:"description" v2:"description"`                     // e.g. "The utility token of the OKX ecosystem"
		Symbol              string         `json:"symbol" v2:"symbol"`                               // e.g. system.Currency
		OriginalSymbol      string         `json:"original_symbol" v2:"original_symbol"`             // e.g. "OKB"
		WholeName           string         `json:"whole_name" v2:"whole_name"`                       // e.g. "OKB"
		OriginalTotalSupply sdk.Dec        `json:"original_total_supply" v2:"original_total_supply"` // e.g. 1000000000.00000000
		TotalSupply         sdk.Dec        `json:"total_supply" v2:"total_supply"`                   // e.g. 1000000000.00000000
		Owner               sdk.AccAddress `json:"owner" v2:"owner"`                                 // e.g. ex1rf9wr069pt64e58f2w3mjs9w72g8vemzw26658
		Mintable            bool           `json:"mintable" v2:"mintable"`                           // e.g. false
	}

	AccCoins struct {
		Acc   sdk.AccAddress `json:"address"`
		Coins sdk.SysCoins   `json:"coins"`
	}
)
