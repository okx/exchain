package coretypes

import abci "github.com/okex/exchain/libs/tendermint/abci/types"

// Query abci msg
type CM39ResultABCIQuery struct {
	Response abci.CM39ResponseQuery `json:"response"`
}

func (c CM39ResultABCIQuery) ToResultABCIQuery() *ResultABCIQuery {
	ret := &ResultABCIQuery{
		Response: abci.ResponseQuery{
			Code:      c.Response.Code,
			Log:       c.Response.Log,
			Info:      c.Response.Info,
			Index:     c.Response.Index,
			Key:       c.Response.Key,
			Value:     c.Response.Value,
			Proof:     c.Response.Proof,
			Height:    c.Response.Height,
			Codespace: c.Response.Codespace,
		},
	}

	return ret
}
