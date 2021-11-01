package types

import (
	"encoding/json"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type Token struct {
	Description         string         `json:"description" v2:"description"`                     // e.g. "OK Group Global Utility Token"
	Symbol              string         `json:"symbol" v2:"symbol"`                               // e.g. "okt"
	OriginalSymbol      string         `json:"original_symbol" v2:"original_symbol"`             // e.g. "OKT"
	WholeName           string         `json:"whole_name" v2:"whole_name"`                       // e.g. "OKT"
	OriginalTotalSupply sdk.Dec        `json:"original_total_supply" v2:"original_total_supply"` // e.g. 1000000000.00000000
	Type                int            `json:"type"`                                             //e.g. 1 common token, 2 interest token
	Owner               sdk.AccAddress `json:"owner" v2:"owner"`                                 // e.g. ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
	Mintable            bool           `json:"mintable" v2:"mintable"`                           // e.g. false
}

func (token Token) String() string {
	b, err := json.Marshal(token)
	if err != nil {
		return "{}"
	}
	return string(b)
}

type TokenResp struct {
	Description         string         `json:"description" v2:"description"`
	Symbol              string         `json:"symbol" v2:"symbol"`
	OriginalSymbol      string         `json:"original_symbol" v2:"original_symbol"`
	WholeName           string         `json:"whole_name" v2:"whole_name"`
	OriginalTotalSupply sdk.Dec        `json:"original_total_supply" v2:"original_total_supply"`
	Type                int            `json:"type"`
	Owner               sdk.AccAddress `json:"owner" v2:"owner"`
	Mintable            bool           `json:"mintable" v2:"mintable"`
	TotalSupply         sdk.Dec        `json:"total_supply" v2:"total_supply"`
}

func (token TokenResp) String() string {
	b, err := json.Marshal(token)
	if err != nil {
		return "{}"
	}
	return string(b)
}

type Tokens []TokenResp

func (tokens Tokens) String() string {
	b, err := json.Marshal(tokens)
	if err != nil {
		return "[{}]"
	}
	return string(b)
}

type Currency struct {
	Description string  `json:"description"`
	Symbol      string  `json:"symbol"`
	TotalSupply sdk.Dec `json:"total_supply"`
}

func (currency Currency) String() string {
	b, err := json.Marshal(currency)
	if err != nil {
		return "[{}]"
	}
	return string(b)
}

type Transfer struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type TransferUnit struct {
	To    sdk.AccAddress `json:"to"`
	Coins sdk.SysCoins   `json:"coins"`
}

type CoinsInfo []CoinInfo

func (d CoinsInfo) Len() int           { return len(d) }
func (d CoinsInfo) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d CoinsInfo) Less(i, j int) bool { return d[i].Symbol < d[j].Symbol }

type AccountResponse struct {
	Address    string    `json:"address"`
	Currencies CoinsInfo `json:"currencies"`
}

func NewAccountResponse(addr string) AccountResponse {
	var accountResponse AccountResponse
	accountResponse.Address = addr
	accountResponse.Currencies = []CoinInfo{}
	return accountResponse
}

type CoinInfo struct {
	Symbol    string `json:"symbol" v2:"currency"`
	Available string `json:"available" v2:"available"`
	Locked    string `json:"locked" v2:"locked"`
}

func NewCoinInfo(symbol, available, locked string) *CoinInfo {
	return &CoinInfo{
		Symbol:    symbol,
		Available: available,
		Locked:    locked,
	}
}

type AccountParam struct {
	Symbol string `json:"symbol"`
	Show   string `json:"show"`
}

type AccountParamV2 struct {
	Currency string `json:"currency"`
	HideZero string `json:"hide_zero"`
}

type AccCoins struct {
	Acc   sdk.AccAddress `json:"address"`
	Coins sdk.SysCoins   `json:"coins"`
}
