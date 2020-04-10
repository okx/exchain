package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 50
)

type QueryDexInfoParams struct {
	Owner   string
	Page    int
	PerPage int
}

func NewQueryDexInfoParams(owner string, page, perPage int) (queryDexInfoParams QueryDexInfoParams, err error) {
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
	q.Page = page
	q.PerPage = perPage
	return nil
}
