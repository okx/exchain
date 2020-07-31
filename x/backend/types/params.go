package types

import "time"

// nolint
const (
	DefaultPage    = 1
	DefaultPerPage = 50
)

// nolint
type QueryDealsParams struct {
	Address string
	Product string
	Start   int64
	End     int64
	Page    int
	PerPage int
	Side    string
}

// NewQueryDealsParams creates a new instance of QueryDealsParams
func NewQueryDealsParams(addr, product string, start, end int64, page, perPage int, side string) QueryDealsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryDealsParams{
		Address: addr,
		Product: product,
		Start:   start,
		End:     end,
		Page:    page,
		PerPage: perPage,
		Side:    side,
	}
}

// nolint
type QueryMatchParams struct {
	Product string
	Start   int64
	End     int64
	Page    int
	PerPage int
}

// NewQueryMatchParams creates a new instance of QueryMatchParams
func NewQueryMatchParams(product string, start, end int64, page, perPage int) QueryMatchParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryMatchParams{
		Product: product,
		Start:   start,
		End:     end,
		Page:    page,
		PerPage: perPage,
	}
}

// nolint
type QueryFeeDetailsParams struct {
	Address string
	Page    int
	PerPage int
}

// NewQueryFeeDetailsParams creates a new instance of QueryFeeDetailsParams
func NewQueryFeeDetailsParams(addr string, page, perPage int) QueryFeeDetailsParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryFeeDetailsParams{
		Address: addr,
		Page:    page,
		PerPage: perPage,
	}
}

// nolint
type QueryKlinesParams struct {
	Product     string
	Granularity int
	Size        int
}

// NewQueryKlinesParams creates a new instance of QueryKlinesParams
func NewQueryKlinesParams(product string, granularity, size int) QueryKlinesParams {
	return QueryKlinesParams{
		product,
		granularity,
		size,
	}
}

// nolint
type QueryTickerParams struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
	Sort    bool   `json:"sort"`
}

// nolint
type QueryOrderListParams struct {
	Address    string
	Product    string
	Page       int
	PerPage    int
	Start      int64
	End        int64
	Side       string
	HideNoFill bool
}

// NewQueryOrderListParams creates  a new instance of QueryOrderListParams
func NewQueryOrderListParams(addr, product, side string, page, perPage int, start, end int64,
	hideNoFill bool) QueryOrderListParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	if start == 0 && end == 0 {
		end = time.Now().Unix()
	}
	return QueryOrderListParams{
		Address:    addr,
		Product:    product,
		Page:       page,
		PerPage:    perPage,
		Start:      start,
		End:        end,
		Side:       side,
		HideNoFill: hideNoFill,
	}
}

// nolint
type QueryTxListParams struct {
	Address   string
	TxType    int64
	StartTime int64
	EndTime   int64
	Page      int
	PerPage   int
}

// NewQueryTxListParams creates a new instance of QueryTxListParams
func NewQueryTxListParams(addr string, txType, startTime, endTime int64, page, perPage int) QueryTxListParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryTxListParams{
		Address:   addr,
		TxType:    txType,
		StartTime: startTime,
		EndTime:   endTime,
		Page:      page,
		PerPage:   perPage,
	}
}

// nolint
type QueryDexFeesParams struct {
	DexHandlingAddr string
	BaseAsset       string
	QuoteAsset      string
	Page            int
	PerPage         int
}

// NewQueryDexFeesParams creates a new instance of QueryDexFeesParams
func NewQueryDexFeesParams(dexHandlingAddr, baseAsset, quoteAsset string, page, perPage int) QueryDexFeesParams {
	if page == 0 && perPage == 0 {
		page = DefaultPage
		perPage = DefaultPerPage
	}
	return QueryDexFeesParams{
		DexHandlingAddr: dexHandlingAddr,
		BaseAsset:       baseAsset,
		QuoteAsset:      quoteAsset,
		Page:            page,
		PerPage:         perPage,
	}
}
