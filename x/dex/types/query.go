package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultPage defines default number of page
	DefaultPage = 1
	// DefaultPerPage defines default number per page
	DefaultPerPage = 50
)

// QueryDexInfoParams defines query params of dex info
type QueryDexInfoParams struct {
	Owner   string
	Page    uint
	PerPage uint
}

// NewQueryDexInfoParams creates query params of dex info
func NewQueryDexInfoParams(owner string, page, perPage uint) (queryDexInfoParams QueryDexInfoParams, err error) {
	if len(owner) == 0 {
		owner = ""
	} else {
		_, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return QueryDexInfoParams{}, fmt.Errorf("invalid address：%s", owner)
		}
	}

	if page <= 0 {
		return QueryDexInfoParams{}, fmt.Errorf("invalid page：%d", page)
	}
	if perPage <= 0 {
		return QueryDexInfoParams{}, fmt.Errorf("invalid per-page：%d", perPage)
	}
	return QueryDexInfoParams{
		Owner:   owner,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// SetPageAndPerPage handles params of page
func (q *QueryDexInfoParams) SetPageAndPerPage(owner, pageStr, perPageStr string) (err error) {

	if len(owner) == 0 {
		owner = ""
	} else {
		_, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return fmt.Errorf("invalid address：%s", owner)
		}
	}
	var page, perPage = DefaultPage, DefaultPerPage

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
		if page <= 0 {
			return fmt.Errorf("invalid page：%s", pageStr)
		}
	}
	if perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			return err
		}
		if perPage <= 0 {
			return fmt.Errorf("invalid per-page：%s", perPageStr)
		}
	}

	q.Owner = owner
	q.Page = uint(page)
	q.PerPage = uint(perPage)
	return nil
}
