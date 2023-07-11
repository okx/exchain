package keeper

import (
	"encoding/json"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetWasmCallInfo(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures)
	keeper := keepers.ContractKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit...)

	codeID, err := keeper.Create(ctx, creator, hackatomWasm, nil)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	em := sdk.NewEventManager()
	// create with no balance is also legal
	ctx.SetEventManager(em)
	gotContractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, nil, initMsgBz, "demo contract 1", nil)
	require.NoError(t, err)

	// 1. contractAddress is equal to storeAddress
	_, _, _, _, err = getCallerInfo(ctx, gotContractAddr.String(), gotContractAddr.String())
	require.NoError(t, err)

	// 2. contractAddress is not exist
	_, _, _, _, err = getCallerInfo(ctx, "0xE70e7466a2f18FAd8C97c45Ba8fEc57d90F3435E", "0xE70e7466a2f18FAd8C97c45Ba8fEc57d90F3435E")
	require.NotNil(t, err)

	// 3. storeAddress is not exist
	_, _, _, _, err = getCallerInfo(ctx, gotContractAddr.String(), "0xE70e7466a2f18FAd8C97c45Ba8fEc57d90F3435E")
	require.NotNil(t, err)

	// 4. contractAddress is not equal to storeAddress
	gotContractAddr2, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, nil, initMsgBz, "demo contract 1", nil)
	_, kvs, q, _, err := getCallerInfo(ctx, gotContractAddr.String(), gotContractAddr2.String())
	require.NoError(t, err)
	require.NotNil(t, kvs)
	require.NotNil(t, q)
}
