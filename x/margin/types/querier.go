package types

import "github.com/okex/okchain/x/token"

// Query endpoints supported by the margin querier
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryMarginAccount = "margin-account"

	QueryParameters = "params"
)

type QueryMarginAccountParams struct {
	SpotAccount string
}

type AccountResponse struct {
	Address    string         `json:"address"`
	Currencies token.CoinInfo `json:"currencies"`
}

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/
