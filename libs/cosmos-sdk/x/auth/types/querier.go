package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address sdk.AccAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}

type CustomGetTxsEventResponse struct {
	// txs is the list of queried transactions.
	Txs []sdk.Tx `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	// tx_responses is the list of queried TxResponses.
	TxResponses []sdk.TxResponse `protobuf:"bytes,2,rep,name=tx_responses,json=txResponses,proto3" json:"tx_responses,omitempty"`
	// pagination defines an pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,3,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
