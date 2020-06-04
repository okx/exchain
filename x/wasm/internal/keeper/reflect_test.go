package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/okex/okchain/x/wasm/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// MaskInitMsg is {}

// MaskHandleMsg is used to encode handle messages
type MaskHandleMsg struct {
	Reflect *reflectPayload `json:"reflect_msg,omitempty"`
	Change  *ownerPayload   `json:"change_owner,omitempty"`
}

type ownerPayload struct {
	Owner sdk.Address `json:"owner"`
}

type reflectPayload struct {
	Msgs []wasmTypes.CosmosMsg `json:"msgs"`
}

// MaskQueryMsg is used to encode query messages
type MaskQueryMsg struct {
	Owner         *struct{} `json:"owner,omitempty"`
	ReflectCustom *Text     `json:"reflect_custom,omitempty"`
}

type Text struct {
	Text string `json:"text"`
}

type OwnerResponse struct {
	Owner string `json:"owner,omitempty"`
}

const MaskFeatures = "staking,mask"

func TestMaskReflectContractSend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, keepers := CreateTestInput(t, false, tempDir, MaskFeatures, maskEncoders(MakeTestCodec()), nil)
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)
	_, _, bob := keyPubAddr()

	// upload mask code
	maskCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	maskID, err := keeper.Create(ctx, creator, maskCode, "", "")
	require.NoError(t, err)
	require.Equal(t, uint64(1), maskID)

	// upload hackatom escrow code
	escrowCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)
	escrowID, err := keeper.Create(ctx, creator, escrowCode, "", "")
	require.NoError(t, err)
	require.Equal(t, uint64(2), escrowID)

	// creator instantiates a contract and gives it tokens
	maskStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	maskAddr, err := keeper.Instantiate(ctx, maskID, creator, []byte("{}"), "mask contract 2", maskStart)
	require.NoError(t, err)
	require.NotEmpty(t, maskAddr)

	// now we set contract as verifier of an escrow
	initMsg := InitMsg{
		Verifier:    maskAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	escrowStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 25000))
	escrowAddr, err := keeper.Instantiate(ctx, escrowID, creator, initMsgBz, "escrow contract 2", escrowStart)
	require.NoError(t, err)
	require.NotEmpty(t, escrowAddr)

	// let's make sure all balances make sense
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 35000))) // 100k - 40k - 25k
	checkAccount(t, ctx, accKeeper, maskAddr, maskStart)
	checkAccount(t, ctx, accKeeper, escrowAddr, escrowStart)
	checkAccount(t, ctx, accKeeper, bob, nil)

	// now for the trick.... we reflect a message through the mask to call the escrow
	// we also send an additional 14k tokens there.
	// this should reduce the mask balance by 14k (to 26k)
	// this 14k is added to the escrow, then the entire balance is sent to bob (total: 39k)
	approveMsg := []byte(`{"release":{}}`)
	msgs := []wasmTypes.CosmosMsg{{
		Wasm: &wasmTypes.WasmMsg{
			Execute: &wasmTypes.ExecuteMsg{
				ContractAddr: escrowAddr.String(),
				Msg:          approveMsg,
				Send: []wasmTypes.Coin{{
					Denom:  "denom",
					Amount: "14000",
				}},
			},
		},
	}}
	reflectSend := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: msgs,
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keeper.Execute(ctx, maskAddr, creator, reflectSendBz, nil)
	require.NoError(t, err)

	// did this work???
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 35000)))  // same as before
	checkAccount(t, ctx, accKeeper, maskAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 26000))) // 40k - 14k (from send)
	checkAccount(t, ctx, accKeeper, escrowAddr, sdk.Coins{})                                  // emptied reserved
	checkAccount(t, ctx, accKeeper, bob, sdk.NewCoins(sdk.NewInt64Coin("denom", 39000)))      // all escrow of 25k + 14k

}

func TestMaskReflectCustomMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, keepers := CreateTestInput(t, false, tempDir, MaskFeatures, maskEncoders(MakeTestCodec()), maskPlugins())
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)
	bob := createFakeFundedAccount(ctx, accKeeper, deposit)
	_, _, fred := keyPubAddr()

	// upload code
	maskCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keeper.Create(ctx, creator, maskCode, "", "")
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	contractAddr, err := keeper.Instantiate(ctx, codeID, creator, []byte("{}"), "mask contract 1", contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// set owner to bob
	transfer := MaskHandleMsg{
		Change: &ownerPayload{
			Owner: bob,
		},
	}
	transferBz, err := json.Marshal(transfer)
	require.NoError(t, err)
	_, err = keeper.Execute(ctx, contractAddr, creator, transferBz, nil)
	require.NoError(t, err)

	// check some account values
	checkAccount(t, ctx, accKeeper, contractAddr, contractStart)
	checkAccount(t, ctx, accKeeper, bob, deposit)
	checkAccount(t, ctx, accKeeper, fred, nil)

	// bob can send contract's tokens to fred (using SendMsg)
	msgs := []wasmTypes.CosmosMsg{{
		Bank: &wasmTypes.BankMsg{
			Send: &wasmTypes.SendMsg{
				FromAddress: contractAddr.String(),
				ToAddress:   fred.String(),
				Amount: []wasmTypes.Coin{{
					Denom:  "denom",
					Amount: "15000",
				}},
			},
		},
	}}
	reflectSend := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: msgs,
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keeper.Execute(ctx, contractAddr, bob, reflectSendBz, nil)
	require.NoError(t, err)

	// fred got coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin("denom", 15000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 25000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)

	// construct an opaque message
	var sdkSendMsg sdk.Msg = &bank.MsgSend{
		FromAddress: contractAddr,
		ToAddress:   fred,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("denom", 23000)),
	}
	opaque, err := toMaskRawMsg(keeper.cdc, sdkSendMsg)
	require.NoError(t, err)
	reflectOpaque := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: []wasmTypes.CosmosMsg{opaque},
		},
	}
	reflectOpaqueBz, err := json.Marshal(reflectOpaque)
	require.NoError(t, err)

	_, err = keeper.Execute(ctx, contractAddr, bob, reflectOpaqueBz, nil)
	require.NoError(t, err)

	// fred got more coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin("denom", 38000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 2000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)
}

func TestMaskReflectCustomQuery(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, keepers := CreateTestInput(t, false, tempDir, MaskFeatures, maskEncoders(MakeTestCodec()), maskPlugins())
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	// upload code
	maskCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keeper.Create(ctx, creator, maskCode, "", "")
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	contractAddr, err := keeper.Instantiate(ctx, codeID, creator, []byte("{}"), "mask contract 1", contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// let's perform a normal query of state
	ownerQuery := MaskQueryMsg{
		Owner: &struct{}{},
	}
	ownerQueryBz, err := json.Marshal(ownerQuery)
	require.NoError(t, err)
	ownerRes, err := keeper.QuerySmart(ctx, contractAddr, ownerQueryBz)
	require.NoError(t, err)
	var res OwnerResponse
	err = json.Unmarshal(ownerRes, &res)
	require.NoError(t, err)
	assert.Equal(t, res.Owner, creator.String())

	// and now making use of the custom querier callbacks
	customQuery := MaskQueryMsg{
		ReflectCustom: &Text{
			Text: "all Caps noW",
		},
	}
	customQueryBz, err := json.Marshal(customQuery)
	require.NoError(t, err)
	custom, err := keeper.QuerySmart(ctx, contractAddr, customQueryBz)
	require.NoError(t, err)
	var resp customQueryResponse
	err = json.Unmarshal(custom, &resp)
	require.NoError(t, err)
	assert.Equal(t, resp.Msg, "ALL CAPS NOW")
}

func checkAccount(t *testing.T, ctx sdk.Context, accKeeper auth.AccountKeeper, addr sdk.AccAddress, expected sdk.Coins) {
	acct := accKeeper.GetAccount(ctx, addr)
	if expected == nil {
		assert.Nil(t, acct)
	} else {
		assert.NotNil(t, acct)
		if expected.Empty() {
			// there is confusion between nil and empty slice... let's just treat them the same
			assert.True(t, acct.GetCoins().Empty())
		} else {
			assert.Equal(t, acct.GetCoins(), expected)
		}
	}
}

/**** Code to support custom messages *****/

type maskCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toMaskRawMsg encodes an sdk msg using amino json encoding.
// Then wraps it as an opaque message
func toMaskRawMsg(cdc *codec.Codec, msg sdk.Msg) (wasmTypes.CosmosMsg, error) {
	rawBz, err := cdc.MarshalJSON(msg)
	if err != nil {
		return wasmTypes.CosmosMsg{}, sdk.ErrInternal(err.Error())
	}
	customMsg, err := json.Marshal(maskCustomMsg{
		Raw: rawBz,
	})
	res := wasmTypes.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}

// maskEncoders needs to be registered in test setup to handle custom message callbacks
func maskEncoders(cdc *codec.Codec) *MessageEncoders {
	return &MessageEncoders{
		Custom: fromMaskRawMsg(cdc),
	}
}

// fromMaskRawMsg decodes msg.Data to an sdk.Msg using amino json encoding.
// this needs to be registered on the Encoders
func fromMaskRawMsg(cdc *codec.Codec) CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var custom maskCustomMsg
		err := json.Unmarshal(msg, &custom)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		if custom.Raw != nil {
			var sdkMsg sdk.Msg
			err := cdc.UnmarshalJSON(custom.Raw, &sdkMsg)
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			return []sdk.Msg{sdkMsg}, nil
		}
		if custom.Debug != "" {
			msg := fmt.Sprintf("Custom Debug: %s", custom.Debug)
			return nil, types.ErrInvalidMsg(msg)
		}
		return nil, types.ErrInvalidMsg("Unknown Custom message variant")
	}
}

type maskCustomQuery struct {
	Ping    *struct{} `json:"ping,omitempty"`
	Capital *Text     `json:"capital,omitempty"`
}

type customQueryResponse struct {
	Msg string `json:"msg"`
}

// maskPlugins needs to be registered in test setup to handle custom query callbacks
func maskPlugins() *QueryPlugins {
	return &QueryPlugins{
		Custom: performCustomQuery,
	}
}

func performCustomQuery(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom maskCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	if custom.Capital != nil {
		msg := strings.ToUpper(custom.Capital.Text)
		return json.Marshal(customQueryResponse{Msg: msg})
	}
	if custom.Ping != nil {
		return json.Marshal(customQueryResponse{Msg: "pong"})
	}
	return nil, types.ErrInvalidMsg("Unknown Custom query variant")
}
