package v0_11

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/token/legacy/v0_10"
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
		Description         string         `json:"description" v2:"description"`                     // e.g. "OK Group Global Utility Token"
		Symbol              string         `json:"symbol" v2:"symbol"`                               // e.g. "okt"
		OriginalSymbol      string         `json:"original_symbol" v2:"original_symbol"`             // e.g. "OKT"
		WholeName           string         `json:"whole_name" v2:"whole_name"`                       // e.g. "OKT"
		OriginalTotalSupply sdk.Dec        `json:"original_total_supply" v2:"original_total_supply"` // e.g. 1000000000.00000000
		Owner               sdk.AccAddress `json:"owner" v2:"owner"`                                 // e.g. ex1rf9wr069pt64e58f2w3mjs9w72g8vemzw26658
		Mintable            bool           `json:"mintable" v2:"mintable"`                           // e.g. false
	}
)
