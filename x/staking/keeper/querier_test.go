package keeper

import (
	"testing"

	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryValidators(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)

	getParams := func(status string) types.QueryValidatorsParams {
		basicParams := types.NewQueryValidatorsParams(1, 100, status)
		return basicParams
	}

	checkPairs := make(map[string]int)
	checkPairs[types2.BondStatusUnbonded] = 1
	checkPairs[types2.BondStatusBonded] = 0
	checkPairs[types2.BondStatusUnbonding] = 0

	querior := NewQuerier(keeper)
	// queryValidators
	for status, expectedCnt := range checkPairs {
		bz, _ := amino.MarshalJSON(getParams(status))
		path := types.QueryValidators
		data, err := querior(ctx, []string{path}, abci.RequestQuery{Data: bz})
		require.True(t, err == nil, err)

		validators := types.Validators{}
		e := amino.UnmarshalJSON(data, &validators)
		require.True(t, e == nil, e)
		if expectedCnt == 1 {
			require.True(t, validators != nil && len(validators) == expectedCnt, validators, status, expectedCnt)

		} else {
			require.True(t, validators == nil, status, expectedCnt)
		}

	}
	data, err := querior(ctx, []string{types.QueryValidators}, abci.RequestQuery{Data: nil})
	require.Error(t, err)
	require.Nil(t, data)
}

func TestQueryValidator(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)

	getParams := func(acc types2.ValAddress) types.QueryValidatorParams {
		basicParams := types.NewQueryValidatorParams(acc)
		return basicParams
	}

	checkAddress := []types2.ValAddress{addrVals[0], addrVals[1]}
	expectedExist := []bool{true, false}

	querior := NewQuerier(keeper)

	// queryValidators
	for i := 0; i < len(checkAddress); i++ {

		bz, _ := amino.MarshalJSON(getParams(checkAddress[i]))
		path := types.QueryValidator
		data, qerr := querior(ctx, []string{path}, abci.RequestQuery{Data: bz})

		validator := types.Validator{}
		e := amino.UnmarshalJSON(data, &validator)
		if expectedExist[i] {
			require.True(t, qerr == nil, qerr)
			require.True(t, e == nil, e)
			require.True(t, validator.OperatorAddress.Equals(checkAddress[i]))

		} else {
			require.True(t, qerr != nil, qerr)
			require.True(t, e != nil, e)
			require.True(t, validator.OperatorAddress == nil || len(validator.OperatorAddress.Bytes()) == 0)
		}
	}
	data, err := querior(ctx, []string{types.QueryValidator}, abci.RequestQuery{Data: nil})
	require.Error(t, err)
	require.Nil(t, data)
}

func TestQueryParams(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper
	querior := NewQuerier(keeper)
	data, err := querior(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.True(t, err == nil)

	params := types.Params{}
	_ = amino.UnmarshalJSON(data, &params)
	require.True(t, params.BondDenom == "okt", params)

}

func TestQueryAddress(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)

	querior := NewQuerier(keeper)
	data, err := querior(ctx, []string{types.QueryAddress}, abci.RequestQuery{})
	require.True(t, err == nil)
	require.NotNil(t, data)

	data, err = querior(ctx, []string{"wrong path"}, abci.RequestQuery{})
	require.Error(t, err)
	require.Nil(t, data)

}

func TestQueryForAddress(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)
	querior := NewQuerier(keeper)

	// successful case
	data, err := querior(ctx, []string{types.QueryForAddress}, abci.RequestQuery{Data: []byte(fakeValidator.GetConsPubKey().Address().String())})
	require.True(t, err == nil)
	require.NotNil(t, data)

	// error case : the request length of param != 40
	data, err = querior(ctx, []string{types.QueryForAddress}, abci.RequestQuery{Data: []byte("wrong address")})
	require.Error(t, err)
	require.Nil(t, data)

	// error case : the request param is not address
	data, err = querior(ctx, []string{types.QueryForAddress}, abci.RequestQuery{Data: []byte("A58856F0FD53BF058B4909A21AEC019107BA6012")})
	require.Error(t, err)
	require.Nil(t, data)
}

func TestQueryForAccAddress(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)
	querior := NewQuerier(keeper)

	// successful case
	data, err := querior(ctx, []string{types.QueryForAccAddress}, abci.RequestQuery{Data: []byte(addrVals[0].String())})
	require.True(t, err == nil)
	require.NotNil(t, data)

	// error case : request param is not address
	data, err = querior(ctx, []string{types.QueryForAccAddress}, abci.RequestQuery{Data: []byte("string")})
	require.Error(t, err)
	require.Nil(t, data)
}

func TestQueriorQueryProxy(t *testing.T) {
	ctx, _, mockKeeper := CreateTestInput(t, false, SufficientInitBalance)
	keeper := mockKeeper.Keeper

	fakeValidator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
	keeper.SetValidator(ctx, fakeValidator)
	querior := NewQuerier(keeper)

	// error case
	badParams := "bad"
	bz, _ := types.ModuleCdc.MarshalJSON(badParams)
	data, err := querior(ctx, []string{types.QueryProxy}, abci.RequestQuery{Data: bz})
	require.True(t, err != nil)
	require.True(t, data == nil)

	// success, but no response
	goodParams := types.NewQueryDelegatorParams(addrDels[0])
	bz, _ = types.ModuleCdc.MarshalJSON(goodParams)
	data, err = querior(ctx, []string{types.QueryProxy}, abci.RequestQuery{Data: bz})
	require.True(t, err == nil)
	require.True(t, data != nil)

}

func TestParams(t *testing.T) {
	ctx, _, keeper := CreateTestInput(t, false, 0)
	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.SetParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}
