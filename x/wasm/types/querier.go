package types

type QueryParamsWithReverse struct {
	Page, Limit int
	Reverse     bool
}

func NewQueryParamsWithReverse(page, limit int, reverse bool) QueryParamsWithReverse {
	return QueryParamsWithReverse{page, limit, reverse}
}
