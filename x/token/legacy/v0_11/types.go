package v0_11

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/token/legacy/v0_10"
)

const ModuleName = "token"

type (
	// all state that must be provided in genesis file
	GenesisState struct {
		Params       v0_10.Params     `json:"params"`
		Tokens       []Token          `json:"tokens"`
		LockedAssets []v0_10.AccCoins `json:"locked_assets"`
		LockedFees   []v0_10.AccCoins `json:"locked_fees"`
	}

	Token struct {
		Description         string         `json:"description" v2:"description"`                     // e.g. "The utility token of the OKX ecosystem"
		Symbol              string         `json:"symbol" v2:"symbol"`                               // e.g. system.Currency
		OriginalSymbol      string         `json:"original_symbol" v2:"original_symbol"`             // e.g. "OKB"
		WholeName           string         `json:"whole_name" v2:"whole_name"`                       // e.g. "OKB"
		OriginalTotalSupply sdk.Dec        `json:"original_total_supply" v2:"original_total_supply"` // e.g. 1000000000.00000000
		Owner               sdk.AccAddress `json:"owner" v2:"owner"`                                 // e.g. ex1rf9wr069pt64e58f2w3mjs9w72g8vemzw26658
		Mintable            bool           `json:"mintable" v2:"mintable"`                           // e.g. false
	}
)
