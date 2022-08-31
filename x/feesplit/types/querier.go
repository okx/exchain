package types

import (
	_ "github.com/gogo/protobuf/gogoproto"
	query "github.com/okex/exchain/libs/cosmos-sdk/types/query"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// QueryFeeSplitsRequest is the request type for the Query/FeeSplits.
type QueryFeeSplitsRequest struct {
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

// QueryFeeSplitsResponse is the response type for the Query/FeeSplits.
type QueryFeeSplitsResponse struct {
	FeeSplits []FeeSplit `json:"fee_splits"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `json:"pagination,omitempty"`
}

// QueryFeeSplitRequest is the request type for the Query/FeeSplit.
type QueryFeeSplitRequest struct {
	// contract identifier is the hex contract address of a contract
	ContractAddress string `json:"contract_address,omitempty"`
}

// QueryFeeSplitResponse is the response type for the Query/FeeSplit.
type QueryFeeSplitResponse struct {
	FeeSplit FeeSplit `json:"fee_split"`
}

// QueryParamsRequest is the request type for the Query/Params.
type QueryParamsRequest struct {
}

// QueryParamsResponse is the response type for the Query/Params.
type QueryParamsResponse struct {
	Params Params `json:"params"`
}

// QueryDeployerFeeSplitsRequest is the request type for the
// Query/DeployerFeeSplits.
type QueryDeployerFeeSplitsRequest struct {
	// deployer bech32 address
	DeployerAddress string `json:"deployer_address,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

// QueryDeployerFeeSplitsResponse is the response type for the
// Query/DeployerFeeSplits.
type QueryDeployerFeeSplitsResponse struct {
	ContractAddresses []string `json:"contract_addresses,omitempty"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `json:"pagination,omitempty"`
}

// QueryWithdrawerFeeSplitsRequest is the request type for the
// Query/WithdrawerFeeSplits.
type QueryWithdrawerFeeSplitsRequest struct {
	// withdrawer bech32 address
	WithdrawerAddress string `json:"withdrawer_address,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

// QueryWithdrawerFeeSplitsResponse is the response type for the
// Query/WithdrawerFeeSplits.
type QueryWithdrawerFeeSplitsResponse struct {
	ContractAddresses []string `json:"contract_addresses,omitempty"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `json:"pagination,omitempty"`
}
