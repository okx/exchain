package keeper

import (
	"fmt"
	"sort"

	"github.com/okex/okchain/x/dex/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
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
		default:
			return nil, sdk.ErrUnknownRequest("unknown dex query endpoint")
		}
	}
}

func queryProduct(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {

	var params types.QueryDexInfoParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", errUnmarshal.Error()))
	}

	var tokenPairs []*types.TokenPair
	if params.Owner != "" {
		ownerAddr, err := sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(fmt.Sprintf("invalid address：%s", params.Owner))
		}

		tokenPairs = keeper.GetUserTokenPairs(ctx, ownerAddr)
	} else {
		tokenPairs = keeper.GetTokenPairs(ctx)
	}

	// sort tokenPairs
	sort.SliceStable(tokenPairs, func(i, j int) bool {
		return tokenPairs[i].ID < tokenPairs[j].ID
	})

	offset, limit := common.GetPage(params.Page, params.PerPage)

	switch {
	case len(tokenPairs) < offset:
		tokenPairs = tokenPairs[0:0]
	case len(tokenPairs) < offset+limit:
		tokenPairs = tokenPairs[offset:]
	default:
		tokenPairs = tokenPairs[offset : offset+limit]
	}

	res, errMarshal := codec.MarshalJSONIndent(types.ModuleCdc, tokenPairs)
	if errMarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to  marshal result to JSON", errMarshal.Error()))
	}
	return res, nil

}

type depositsData struct {
	ProductName     string      `json:"product"`
	ProductDeposits sdk.DecCoin `json:"deposits"`
}

func queryDeposits(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {

	var params types.QueryDexInfoParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", errUnmarshal.Error()))
	}

	var tokenPairs []*types.TokenPair
	if params.Owner != "" {
		ownerAddr, err := sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(fmt.Sprintf("invalid address：%s", params.Owner))
		}

		tokenPairs = keeper.GetUserTokenPairs(ctx, ownerAddr)
	} else {
		tokenPairs = keeper.GetTokenPairs(ctx)
	}

	var deposits []depositsData
	for _, product := range tokenPairs {
		if product == nil {
			return res, types.ErrNilPointer()
		}
		if product.Owner.String() == params.Owner {
			deposits = append(deposits, depositsData{fmt.Sprintf("%s_%s", product.BaseAssetSymbol, product.QuoteAssetSymbol), product.Deposits})
		}
	}

	offset, limit := common.GetPage(params.Page, params.PerPage)

	switch {
	case len(deposits) < offset:
		deposits = deposits[0:0]
	case len(deposits) < offset+limit:
		deposits = deposits[offset:]
	default:
		deposits = deposits[offset : offset+limit]
	}

	sort.SliceStable(deposits, func(i, j int) bool {
		return deposits[i].ProductDeposits.IsLT(deposits[j].ProductDeposits)
	})

	res, errMarshal := codec.MarshalJSONIndent(types.ModuleCdc, deposits)
	if errMarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to  marshal result to JSON", errMarshal.Error()))
	}
	return res, nil
}

func queryMatchOrder(ctx sdk.Context, req abci.RequestQuery, keeper IKeeper) (res []byte, err sdk.Error) {

	var params types.QueryDexInfoParams
	errUnmarshal := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if errUnmarshal != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", errUnmarshal.Error()))
	}

	tokenPairs := keeper.GetTokenPairsOrdered(ctx)

	var products []string

	for _, tokenPair := range tokenPairs {
		if tokenPair == nil {
			return res, types.ErrNilPointer()
		}
		products = append(products, fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol))
	}

	offset, limit := common.GetPage(params.Page, params.PerPage)

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
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
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
			return res, types.ErrNilPointer()
		}
		if tokenPairs[i].Delisting {
			tokenPairNames = append(tokenPairNames, fmt.Sprintf("%s_%s", tokenPairs[i].BaseAssetSymbol, tokenPairs[i].QuoteAssetSymbol))
		}
	}

	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, tokenPairNames)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to  marshal result to JSON", errUnmarshal.Error()))
	}

	return res, nil

}
