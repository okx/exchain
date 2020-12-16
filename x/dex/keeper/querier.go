package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/okex/okexchain/x/dex/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper IKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryProducts:
			return queryProduct(ctx, req, keeper)
		case types.QueryDeposits:
			return queryDeposits(ctx, req, keeper)
		case types.QueryMatchOrder:
			return queryMatchOrder(ctx, req, keeper)
		case types.QueryParameters:
			return queryParams(ctx, req, keeper)
		case types.QueryProductsDelisting:
			return queryProductsDelisting(ctx, keeper)
		case types.QueryOperator:
			return queryOperator(ctx, req, keeper)
		case types.QueryOperators:
			return queryOperators(ctx, keeper)
		default:
			return nil, types.ErrUnknownQueryType()
		}
	}
}

func queryProduct(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {
	var params types.QueryDexInfoParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, common.ErrUnMarshalJSONFailed(errUnmarshal.Error())
	}

	offset, limit := common.GetPage(int(params.Page), int(params.PerPage))

	if offset < 0 || limit < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}

	var tokenPairs []*types.TokenPair
	if params.Owner != "" {
		ownerAddr, err := sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, common.ErrCreateAddrFromBech32Failed(params.Owner)
		}

		tokenPairs = keeper.GetUserTokenPairs(ctx, ownerAddr)
	} else {
		tokenPairs = keeper.GetTokenPairs(ctx)
	}

	// sort tokenPairs
	sort.SliceStable(tokenPairs, func(i, j int) bool {
		return tokenPairs[i].ID < tokenPairs[j].ID
	})

	total := len(tokenPairs)
	switch {
	case total < offset:
		tokenPairs = tokenPairs[0:0]
	case total < offset+limit:
		tokenPairs = tokenPairs[offset:]
	default:
		tokenPairs = tokenPairs[offset : offset+limit]
	}

	var response *common.ListResponse
	if len(tokenPairs) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, tokenPairs)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}

	res, errMarshal := json.MarshalIndent(response, "", "  ")
	if errMarshal != nil {
		return nil, common.ErrMarshalJSONFailed(errMarshal.Error())
	}
	return res, nil

}

type depositsData struct {
	ProductName     string         `json:"product"`
	ProductDeposits sdk.SysCoin    `json:"deposits"`
	Rank            int            `json:"rank"`
	BlockHeight     int64          `json:"block_height"`
	Owner           sdk.AccAddress `json:"owner"`
}

func queryDeposits(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {
	var params types.QueryDepositParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, common.ErrUnMarshalJSONFailed(errUnmarshal.Error())
	}

	if params.Address == "" && params.BaseAsset == "" && params.QuoteAsset == "" {
		return nil, types.ErrAddrAndProductAllRequired()
	}

	offset, limit := common.GetPage(int(params.Page), int(params.PerPage))
	if offset < 0 || limit < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}

	tokenPairs := keeper.GetTokenPairsOrdered(ctx)

	var deposits []depositsData
	for i, tokenPair := range tokenPairs {
		if tokenPair == nil {
			return nil, types.ErrIsNil()
		}
		// filter address
		if params.Address != "" && tokenPair.Owner.String() != params.Address {
			continue
		}
		// filter base asset
		if params.BaseAsset != "" && !strings.Contains(tokenPair.BaseAssetSymbol, params.BaseAsset) {
			continue
		}
		// filter quote asset
		if params.QuoteAsset != "" && !strings.Contains(tokenPair.QuoteAssetSymbol, params.QuoteAsset) {
			continue
		}
		deposits = append(deposits, depositsData{fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol), tokenPair.Deposits, i + 1, tokenPair.BlockHeight, tokenPair.Owner})
	}
	total := len(deposits)

	switch {
	case total < offset:
		deposits = deposits[0:0]
	case total < offset+limit:
		deposits = deposits[offset:]
	default:
		deposits = deposits[offset : offset+limit]
	}

	sort.SliceStable(deposits, func(i, j int) bool {
		return deposits[i].ProductDeposits.IsLT(deposits[j].ProductDeposits)
	})

	var response *common.ListResponse
	if total > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, deposits)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}

	res, errMarshal := json.MarshalIndent(response, "", "  ")
	if errMarshal != nil {
		return nil, common.ErrMarshalJSONFailed(errMarshal.Error())
	}

	return res, nil
}

func queryMatchOrder(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {

	var params types.QueryDexInfoParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, common.ErrUnMarshalJSONFailed(errUnmarshal.Error())
	}
	offset, limit := common.GetPage(int(params.Page), int(params.PerPage))

	if offset < 0 || limit < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}
	tokenPairs := keeper.GetTokenPairsOrdered(ctx)

	var products []string

	for _, tokenPair := range tokenPairs {
		if tokenPair == nil {
			panic("the nil pointer is not expected")
		}
		products = append(products, fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol))
	}

	switch {
	case len(products) < offset:
		products = products[0:0]
	case len(products) < offset+limit:
		products = products[offset:]
	default:
		products = products[offset : offset+limit]
	}

	res, errMarshal := codec.MarshalJSONIndent(types.ModuleCdc, products)

	if errMarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to  marshal result to JSON", errMarshal.Error()))
	}
	return res, nil

}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)
	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if errUnmarshal != nil {
		return nil, common.ErrMarshalJSONFailed(errUnmarshal.Error())
	}
	return res, nil
}

//queryProductsDelisting query the tokenpair name under dex delisting
func queryProductsDelisting(ctx sdk.Context, keeper IKeeper) (res []byte, err sdk.Error) {
	var tokenPairNames []string
	tokenPairs := keeper.GetTokenPairs(ctx)
	tokenPairLen := len(tokenPairs)
	for i := 0; i < tokenPairLen; i++ {
		if tokenPairs[i] == nil {
			return nil, types.ErrTokenPairIsInvalid()
		}
		if tokenPairs[i].Delisting {
			tokenPairNames = append(tokenPairNames, fmt.Sprintf("%s_%s", tokenPairs[i].BaseAssetSymbol, tokenPairs[i].QuoteAssetSymbol))
		}
	}

	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, tokenPairNames)
	if errUnmarshal != nil {
		return nil, common.ErrMarshalJSONFailed(errUnmarshal.Error())
	}

	return res, nil
}

// nolint
func queryOperator(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) ([]byte, sdk.Error) {
	var params types.QueryDexOperatorParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	operator, isExist := keeper.GetOperator(ctx, params.Addr)
	if !isExist {
		return nil, types.ErrUnknownOperator(params.Addr)
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, operator)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// nolint
func queryOperators(ctx sdk.Context, keeper IKeeper) ([]byte, sdk.Error) {
	var operators types.DEXOperators
	keeper.IterateOperators(ctx, func(operator types.DEXOperator) bool {
		//info.HandlingFees = keeper.GetBankKeeper().GetCoins(ctx, info.HandlingFeeAddress).String()
		operators = append(operators, operator)
		return false
	})

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, operators)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}
