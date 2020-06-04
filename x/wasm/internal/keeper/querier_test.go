package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/wasm/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryContractState(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, keepers := CreateTestInput(t, false, tempDir, SupportedFeatures, nil, nil)
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit))
	anyAddr := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.Create(ctx, creator, wasmCode, "", "")
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    anyAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.Instantiate(ctx, contractID, creator, initMsgBz, "demo contract to query", deposit)
	require.NoError(t, err)

	contractModel := []types.Model{
		{Key: []byte("foo"), Value: []byte(`"bar"`)},
		{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)},
	}
	keeper.setContractState(ctx, addr, contractModel)

	// this gets us full error, not redacted sdk.Error
	q := NewQuerier(keeper)
	specs := map[string]struct {
		srcPath []string
		srcReq  abci.RequestQuery
		// smart queries return raw bytes from contract not []types.Model
		// if this is set, then we just compare - (should be json encoded string)
		expSmartRes string
		// if success and expSmartRes is not set, we parse into []types.Model and compare
		expModelLen      int
		expModelContains []types.Model
		expErr           sdk.Error
	}{
		"query all": {
			srcPath:     []string{QueryGetContractState, addr.String(), QueryMethodContractStateAll},
			expModelLen: 3,
			expModelContains: []types.Model{
				{Key: []byte("foo"), Value: []byte(`"bar"`)},
				{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)},
			},
		},
		"query raw key": {
			srcPath:          []string{QueryGetContractState, addr.String(), QueryMethodContractStateRaw},
			srcReq:           abci.RequestQuery{Data: []byte("foo")},
			expModelLen:      1,
			expModelContains: []types.Model{{Key: []byte("foo"), Value: []byte(`"bar"`)}},
		},
		"query raw binary key": {
			srcPath:          []string{QueryGetContractState, addr.String(), QueryMethodContractStateRaw},
			srcReq:           abci.RequestQuery{Data: []byte{0x0, 0x1}},
			expModelLen:      1,
			expModelContains: []types.Model{{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)}},
		},
		"query smart": {
			srcPath:     []string{QueryGetContractState, addr.String(), QueryMethodContractStateSmart},
			srcReq:      abci.RequestQuery{Data: []byte(`{"verifier":{}}`)},
			expSmartRes: fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()),
		},
		"query smart invalid request": {
			srcPath: []string{QueryGetContractState, addr.String(), QueryMethodContractStateSmart},
			srcReq:  abci.RequestQuery{Data: []byte(`{"raw":{"key":"config"}}`)},
			expErr:  types.ErrQueryFailed(""),
		},
		"query unknown raw key": {
			srcPath:     []string{QueryGetContractState, addr.String(), QueryMethodContractStateRaw},
			srcReq:      abci.RequestQuery{Data: []byte("unknown")},
			expModelLen: 0,
		},
		"query empty raw key": {
			srcPath:     []string{QueryGetContractState, addr.String(), QueryMethodContractStateRaw},
			expModelLen: 0,
		},
		"query raw with unknown address": {
			srcPath:     []string{QueryGetContractState, anyAddr.String(), QueryMethodContractStateRaw},
			expModelLen: 0,
		},
		"query all with unknown address": {
			srcPath:     []string{QueryGetContractState, anyAddr.String(), QueryMethodContractStateAll},
			expModelLen: 0,
		},
		"query smart with unknown address": {
			srcPath:     []string{QueryGetContractState, anyAddr.String(), QueryMethodContractStateSmart},
			expModelLen: 0,
			expErr:      types.ErrNotFound(""),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			binResult, _ := q(ctx, spec.srcPath, spec.srcReq)
			// require.True(t, spec.expErr.Is(err), "unexpected error")
			//require.True(t, spec.expErr.Is(err), err)

			// if smart query, check custom response
			if spec.expSmartRes != "" {
				require.Equal(t, spec.expSmartRes, string(binResult))
				return
			}

			// otherwise, check returned models
			var r []types.Model
			if spec.expErr == nil {
				require.NoError(t, json.Unmarshal(binResult, &r))
				require.NotNil(t, r)
			}
			require.Len(t, r, spec.expModelLen)
			// and in result set
			for _, v := range spec.expModelContains {
				assert.Contains(t, r, v)
			}
		})
	}
}

func TestListContractByCodeOrdering(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, keepers := CreateTestInput(t, false, tempDir, SupportedFeatures, nil, nil)
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 500))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)
	anyAddr := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	codeID, err := keeper.Create(ctx, creator, wasmCode, "", "")
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    anyAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	// manage some realistic block settings
	var h int64 = 10
	setBlock := func(ctx sdk.Context, height int64) sdk.Context {
		ctx = ctx.WithBlockHeight(height)
		meter := sdk.NewGasMeter(1000000)
		ctx = ctx.WithGasMeter(meter)
		ctx = ctx.WithBlockGasMeter(meter)
		return ctx
	}

	// create 10 contracts with real block/gas setup
	for i := range [10]int{} {
		// 3 tx per block, so we ensure both comparisons work
		if i%3 == 0 {
			ctx = setBlock(ctx, h)
			h++
		}
		_, err = keeper.Instantiate(ctx, codeID, creator, initMsgBz, fmt.Sprintf("contract %d", i), topUp)
		require.NoError(t, err)
	}

	// query and check the results are properly sorted
	q := NewQuerier(keeper)
	query := []string{QueryListContractByCode, fmt.Sprintf("%d", codeID)}
	data := abci.RequestQuery{}
	res, err := q(ctx, query, data)
	require.NoError(t, err)

	var contracts []ContractInfoWithAddress
	err = json.Unmarshal(res, &contracts)
	require.NoError(t, err)

	require.Equal(t, 10, len(contracts))

	for i, contract := range contracts {
		assert.Equal(t, fmt.Sprintf("contract %d", i), contract.Label)
		assert.NotEmpty(t, contract.Address)
		// ensure these are not shown
		assert.Nil(t, contract.InitMsg)
		assert.Nil(t, contract.Created)
	}
}
