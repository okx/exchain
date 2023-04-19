package wasm

import (
	"encoding/json"
	"fmt"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/types"
)

var (
	addr, _ = sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
)

func testMustAccAddressFromBech32(addr string) sdk.AccAddress {
	re, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}
	return re
}

func newTestSysCoin(i int64, precison int64) sdk.SysCoin {
	return sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(i, precison))
}

func TestMsgStoreCode(t *testing.T) {
	msg := types.MsgStoreCode{
		Sender:                "0x67582AB2adb08a8583A181b7745762B53710e9B1",
		WASMByteCode:          []byte("hello"),
		InstantiatePermission: &types.AccessConfig{3, "0x67582AB2adb08a8583A181b7745762B53710e9B1"},
	}
	d, err := json.Marshal(msg)
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestMsgInstantiateContract(t *testing.T) {
	msg := types.MsgInstantiateContract{
		Sender: "0x67582AB2adb08a8583A181b7745762B53710e9B1",
		Admin:  "0x67582AB2adb08a8583A181b7745762B53710e9B2",
		CodeID: 2,
		Label:  "hello",
		Msg:    []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
		Funds:  sdk.CoinsToCoinAdapters([]sdk.DecCoin{sdk.NewDecCoin("mytoken", sdk.NewInt(10))}),
	}
	d, err := json.Marshal(msg)
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestMsgExecuteContract(t *testing.T) {
	msg := types.MsgExecuteContract{
		Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
		Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
		Msg:      []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
		Funds:    sdk.CoinsToCoinAdapters([]sdk.DecCoin{sdk.NewDecCoin("mytoken", sdk.NewInt(10))}),
	}
	d, err := json.Marshal(msg)
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestMsgMigrateContract(t *testing.T) {
	msg := types.MsgMigrateContract{
		Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
		Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
		CodeID:   1,
		Msg:      []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
	}
	d, err := json.Marshal(msg)
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestMsgUpdateAdmin(t *testing.T) {
	msg := types.MsgUpdateAdmin{
		Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
		NewAdmin: "0x67582AB2adb08a8583A181b7745762B53710e9B3",
		Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
	}
	d, err := json.Marshal(msg)
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestConvertMsgStoreCode(t *testing.T) {
	//addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	//require.NoError(t, err)

	testcases := []struct {
		msgstr  string
		res     types.MsgStoreCode
		fnCheck func(msg sdk.Msg, err error, res types.MsgStoreCode)
	}{
		{
			msgstr: "{\"sender\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\",\"wasm_byte_code\":\"aGVsbG8=\",\"instantiate_permission\":{\"permission\":\"OnlyAddress\",\"address\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\"}}",
			res: types.MsgStoreCode{
				Sender:                "0x67582AB2adb08a8583A181b7745762B53710e9B1",
				WASMByteCode:          []byte("hello"),
				InstantiatePermission: &types.AccessConfig{2, "0x67582AB2adb08a8583A181b7745762B53710e9B1"},
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgStoreCode) {
				require.NoError(t, err)
				require.Equal(t, *msg.(*types.MsgStoreCode), res)
			},
		},
	}

	tmtypes.InitMilestoneVenus6Height(1)
	for _, ts := range testcases {
		msg, err := ConvertMsgStoreCode([]byte(ts.msgstr), ts.res.GetSigners(), 2)
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertMsgInstantiateContract(t *testing.T) {
	//addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	//require.NoError(t, err)

	testcases := []struct {
		msgstr  string
		res     types.MsgInstantiateContract
		fnCheck func(msg sdk.Msg, err error, res types.MsgInstantiateContract)
	}{
		{
			msgstr: "{\"sender\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\",\"admin\":\"0x67582AB2adb08a8583A181b7745762B53710e9B2\",\"code_id\":2,\"label\":\"hello\",\"msg\":{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}},\"funds\":[{\"denom\":\"mytoken\",\"amount\":\"10000000000000000000\"}]}",
			res: types.MsgInstantiateContract{
				Sender: "0x67582AB2adb08a8583A181b7745762B53710e9B1",
				Admin:  "0x67582AB2adb08a8583A181b7745762B53710e9B2",
				CodeID: 2,
				Label:  "hello",
				Msg:    []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
				Funds:  sdk.CoinsToCoinAdapters([]sdk.DecCoin{sdk.NewDecCoin("mytoken", sdk.NewInt(10))}),
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgInstantiateContract) {
				require.NoError(t, err)
				require.Equal(t, *msg.(*types.MsgInstantiateContract), res)
			},
		},
	}

	tmtypes.InitMilestoneVenus6Height(1)
	for _, ts := range testcases {
		msg, err := ConvertMsgInstantiateContract([]byte(ts.msgstr), ts.res.GetSigners(), 2)
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertMsgExecuteContract(t *testing.T) {
	testcases := []struct {
		msgstr  string
		res     types.MsgExecuteContract
		fnCheck func(msg sdk.Msg, err error, res types.MsgExecuteContract)
	}{
		{
			msgstr: "{\"sender\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\",\"contract\":\"0x67582AB2adb08a8583A181b7745762B53710e9B2\",\"msg\":{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}},\"funds\":[{\"denom\":\"mytoken\",\"amount\":\"10000000000000000000\"}]}",
			res: types.MsgExecuteContract{
				Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
				Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
				Msg:      []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
				Funds:    sdk.CoinsToCoinAdapters([]sdk.DecCoin{sdk.NewDecCoin("mytoken", sdk.NewInt(10))}),
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgExecuteContract) {
				require.NoError(t, err)
				require.Equal(t, *msg.(*types.MsgExecuteContract), res)
			},
		},
	}

	tmtypes.InitMilestoneVenus6Height(1)
	for _, ts := range testcases {
		msg, err := ConvertMsgExecuteContract([]byte(ts.msgstr), ts.res.GetSigners(), 2)
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertMsgMigrateContract(t *testing.T) {
	testcases := []struct {
		msgstr  string
		res     types.MsgMigrateContract
		fnCheck func(msg sdk.Msg, err error, res types.MsgMigrateContract)
	}{
		{
			msgstr: "{\"sender\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\",\"contract\":\"0x67582AB2adb08a8583A181b7745762B53710e9B2\",\"code_id\":1,\"msg\":{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}}",
			res: types.MsgMigrateContract{
				Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
				Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
				CodeID:   1,
				Msg:      []byte("{\"balance\":{\"address\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"),
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgMigrateContract) {
				require.NoError(t, err)
				require.Equal(t, *msg.(*types.MsgMigrateContract), res)
			},
		},
	}
	tmtypes.InitMilestoneVenus6Height(1)
	for _, ts := range testcases {
		msg, err := ConvertMsgMigrateContract([]byte(ts.msgstr), ts.res.GetSigners(), 2)
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertMsgUpdateAdmin(t *testing.T) {
	testcases := []struct {
		msgstr  string
		res     types.MsgUpdateAdmin
		fnCheck func(msg sdk.Msg, err error, res types.MsgUpdateAdmin)
	}{
		{
			msgstr: "{\"sender\":\"0x67582AB2adb08a8583A181b7745762B53710e9B1\",\"new_admin\":\"0x67582AB2adb08a8583A181b7745762B53710e9B3\",\"contract\":\"0x67582AB2adb08a8583A181b7745762B53710e9B2\"}",
			res: types.MsgUpdateAdmin{
				Sender:   "0x67582AB2adb08a8583A181b7745762B53710e9B1",
				NewAdmin: "0x67582AB2adb08a8583A181b7745762B53710e9B3",
				Contract: "0x67582AB2adb08a8583A181b7745762B53710e9B2",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateAdmin) {
				require.NoError(t, err)
				require.Equal(t, *msg.(*types.MsgUpdateAdmin), res)
			},
		},
	}
	tmtypes.InitMilestoneVenus6Height(1)
	for _, ts := range testcases {
		msg, err := ConvertMsgUpdateAdmin([]byte(ts.msgstr), ts.res.GetSigners(), 2)
		ts.fnCheck(msg, err, ts.res)
	}
}
