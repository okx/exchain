package types

type QueryParamsWithReverse struct {
	Page    int  `json:"page,omitempty"`
	Limit   int  `json:"limit,omitempty"`
	Reverse bool `json:"reverse,omitempty"`
}

func NewQueryParamsWithReverse(page, limit int, reverse bool) QueryParamsWithReverse {
	return QueryParamsWithReverse{page, limit, reverse}
}

type QueryCodeInfoResponse struct {
	Data  []CodeInfoResponse `json:"data,omitempty"`
	Total int                `json:"total,omitempty"`
}

func NewQueryCodeInfoResponse(data []CodeInfoResponse, total int) QueryCodeInfoResponse {
	return QueryCodeInfoResponse{data, total}
}

type QueryContractListByCodeResponse struct {
	Data  []string `json:"data"`
	Total int      `json:"total,omitempty"`
}

func NewQueryContractListByCodeResponse(data []string, total int) QueryContractListByCodeResponse {
	return QueryContractListByCodeResponse{data, total}
}

type QueryContractStateAllResponse struct {
	Data  []Model `json:"data,omitempty"`
	Total int     `json:"total,omitempty"`
}

func NewQueryContractStateAllResponse(data []Model, total int) QueryContractStateAllResponse {
	return QueryContractStateAllResponse{data, total}
}

type QueryContractCodeHistoryResponse struct {
	Data  []ContractCodeHistoryEntry `json:"data"`
	Total int                        `json:"total,omitempty"`
}

func NewQueryContractCodeHistoryResponse(data []ContractCodeHistoryEntry, total int) QueryContractCodeHistoryResponse {
	return QueryContractCodeHistoryResponse{data, total}
}
