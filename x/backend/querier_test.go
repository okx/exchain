package backend

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	ordertypes "github.com/okex/okchain/x/order/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func mockQuerier(t *testing.T) (*MockApp, sdk.Context, sdk.Querier, []*ordertypes.Order) {

	mapp, orders := FireEndBlockerPeriodicMatch(t, true)
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	time.Sleep(time.Second)
	querier := NewQuerier(mapp.backendKeeper)

	return mapp, ctx, querier, orders
}

func TestQuerier_queryTickerList(t *testing.T) {
	mapp, ctx, querier, _ := mockQuerier(t)
	path := []string{types.QueryTickerList}

	params := types.QueryTickerParams{
		Sort:  true,
		Count: 100,
	}
	request := abci.RequestQuery{}

	// 1. Invalid request
	invalidRequest := request
	_, err := querier(ctx, path, invalidRequest)
	assert.True(t, err != nil)

	// 2. No product request.
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)
	bytesBuffer, err := querier(ctx, path, request)
	require.Nil(t, err)
	finalResult := &map[string]interface{}{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))

	data := (*finalResult)["data"]
	assert.True(t, data != nil)

	// 3. With product request.
	params2 := params
	params2.Product = types.TestTokenPair
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params2)
	bytesBuffer, err = querier(ctx, path, request)
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	data = (*finalResult)["data"]
	assert.True(t, data != nil)
}

func TestQuerier_queryMatchList(t *testing.T) {

	mapp, ctx, querier, _ := mockQuerier(t)
	params := types.NewQueryMatchParams(types.TestTokenPair, 0, 0, 1, 100)
	request := abci.RequestQuery{}
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)

	path := []string{types.QueryMatchResults}

	bytesBuffer, err := querier(ctx, path, request)
	finalResult := &map[string]interface{}{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	assert.True(t, err == nil)
}

func TestQuerier_queryDealsList(t *testing.T) {

	mapp, ctx, querier, _ := mockQuerier(t)
	params := types.NewQueryDealsParams("NotExists", types.TestTokenPair, 0, 0, 1, 100, "")
	request := abci.RequestQuery{}
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)

	path := []string{types.QueryDealList}

	bytesBuffer, err := querier(ctx, path, request)
	finalResult := &map[string]interface{}{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	assert.True(t, err == nil)

	params = types.NewQueryDealsParams("NotExists", types.TestTokenPair, 0, 0, 1, 100, types.BuyOrder)
	request = abci.RequestQuery{}
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)
	bytesBuffer, err = querier(ctx, path, request)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	assert.True(t, err == nil)
}

func TestQuerier_queryCandleList(t *testing.T) {

	mapp, ctx, querier, _ := mockQuerier(t)
	//time.Sleep(time.Second * 60)

	params := types.NewQueryKlinesParams(types.TestTokenPair, 60, 100)
	request := abci.RequestQuery{}
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)

	path := []string{types.QueryCandleList}
	bytesBuffer, err := querier(ctx, path, request)
	finalResult := &map[string]interface{}{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	assert.True(t, err == nil)
	data := (*finalResult)["data"]
	assert.True(t, data != nil && reflect.TypeOf(data).Kind() == reflect.Slice)
}

func TestQuerier_QueryCandlesByMarketService(t *testing.T) {

}

func TestQuerier_QueryFeeDetails(t *testing.T) {
	mapp, ctx, querier, _ := mockQuerier(t)

	params := types.NewQueryFeeDetailsParams("NotExists", 1, 10)
	request := abci.RequestQuery{}
	request.Data, _ = mapp.Cdc.MarshalJSON(params)

	path := []string{types.QueryFeeDetails}
	_, err := querier(ctx, path, request)
	require.NotNil(t, err)
}

func TestQuerier_QueryOrderList(t *testing.T) {
	mapp, ctx, querier, orders := mockQuerier(t)

	params := types.NewQueryOrderListParams("NotExists", types.TestTokenPair, "", 1, 10, 0, 0, false)
	request := abci.RequestQuery{}
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params)
	path := []string{types.QueryOrderList, "open"}

	bytesBuffer, err := querier(ctx, path, request)
	finalResult := &common.ListResponse{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	data := (*finalResult).Data
	assert.True(t, err != nil && data.Data == nil)

	params2 := params
	params2.Address = orders[0].Sender.String()
	request.Data, _ = mapp.backendKeeper.MarshalJSON(params2)
	bytesBuffer, err = querier(ctx, path, request)
	finalResult = &common.ListResponse{}
	json.Unmarshal(bytesBuffer, finalResult)
	fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	assert.True(t, err == nil)

}

func TestQuerier_QueryTxList(t *testing.T) {
	mapp, ctx, querier, orders := mockQuerier(t)
	params := types.NewQueryTxListParams(orders[0].Sender.String(), 1, 0, time.Now().Unix(), 1, 10)
	path := []string{types.QueryTxList}
	request := abci.RequestQuery{}

	for i := 1; i <= 3; i++ {
		params.TxType = int64(i)
		request.Data, _ = mapp.backendKeeper.MarshalJSON(params)
		bytesBuffer, err := querier(ctx, path, request)
		assert.True(t, err == nil)
		finalResult := &common.ListResponse{}
		json.Unmarshal(bytesBuffer, finalResult)
		fmt.Println(fmt.Sprintf("finalResult: %+v, bytes: %s", finalResult, bytesBuffer))
	}

}
