package types

type QueryParamsWithReverse struct {
	Page    int  `json:"page,omitempty"`
	Limit   int  `json:"limit,omitempty"`
	Reverse bool `json:"reverse,omitempty"`
}

func NewQueryParamsWithReverse(page, limit int, reverse bool) QueryParamsWithReverse {
	return QueryParamsWithReverse{page, limit, reverse}
}
