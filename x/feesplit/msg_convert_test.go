package feesplit

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/feesplit/types"
	"github.com/stretchr/testify/require"
)

func testMustAccAddressFromBech32(addr string) sdk.AccAddress {
	re, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}
	return re
}

func TestConvertRegisterFeeSplitMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)

	contractAddr := common.HexToAddress("0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414")

	testcases := []struct {
		msgstr  string
		res     types.MsgRegisterFeeSplit
		fnCheck func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit)
	}{
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\",\"withdrawer_address\":\"%s\",\"nonces\":[2,4]}", addr.String(), addr.String()),
			res: types.NewMsgRegisterFeeSplit(contractAddr,
				testMustAccAddressFromBech32(addr.String()),
				testMustAccAddressFromBech32(addr.String()),
				[]uint64{2, 4},
			),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgRegisterFeeSplit), res)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\",\"withdrawer_address\":\"%s\",\"nonces\":[2]}", addr.String(), addr.String()),
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   addr.String(),
				WithdrawerAddress: addr.String(),
				Nonces:            []uint64{2}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgRegisterFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","nonces":[1,2,4]}`,
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				WithdrawerAddress: "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				Nonces:            []uint64{1, 2, 4}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgRegisterFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","nonces":[1,2,4]}`,
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				Nonces:            []uint64{1, 2, 4}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgRegisterFeeSplit), res)
			},
		},
		// error
		{
			msgstr: "123",
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				Nonces:            []uint64{1, 2, 4}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","nonces":[]}`,
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				Nonces:            []uint64{1, 2, 4}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","nonces":[1,2,4]}`,
			res: types.MsgRegisterFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				Nonces:            []uint64{1, 2, 4}},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgRegisterFeeSplit) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertRegisterFeeSplitMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertUpdateFeeSplitMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)

	contractAddr := common.HexToAddress("0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414")

	testcases := []struct {
		msgstr  string
		res     types.MsgUpdateFeeSplit
		fnCheck func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit)
	}{
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\",\"withdrawer_address\":\"%s\"}", addr.String(), addr.String()),
			res: types.NewMsgUpdateFeeSplit(contractAddr,
				testMustAccAddressFromBech32(addr.String()),
				testMustAccAddressFromBech32(addr.String()),
			),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgUpdateFeeSplit), res)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\",\"withdrawer_address\":\"%s\"}", addr.String(), addr.String()),
			res: types.MsgUpdateFeeSplit{
				ContractAddress:   "A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   addr.String(),
				WithdrawerAddress: addr.String(),
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgUpdateFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgUpdateFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				WithdrawerAddress: "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgUpdateFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgUpdateFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgUpdateFeeSplit), res)
			},
		},
		// error
		{
			msgstr: "123",
			res: types.MsgUpdateFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9","withdrawer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgUpdateFeeSplit{
				ContractAddress:   "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress:   "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
				WithdrawerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgUpdateFeeSplit) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertUpdateFeeSplitMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertCancelFeeSplitMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)

	contractAddr := common.HexToAddress("0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414")

	testcases := []struct {
		msgstr  string
		res     types.MsgCancelFeeSplit
		fnCheck func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit)
	}{
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\"}", addr.String()),
			res: types.NewMsgCancelFeeSplit(
				contractAddr,
				testMustAccAddressFromBech32(addr.String()),
			),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgCancelFeeSplit), res)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"contract_address\":\"A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414\",\"deployer_address\":\"%s\"}", addr.String()),
			res: types.MsgCancelFeeSplit{
				ContractAddress: "A4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress: addr.String(),
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgCancelFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgCancelFeeSplit{
				ContractAddress: "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress: "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgCancelFeeSplit), res)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgCancelFeeSplit{
				ContractAddress: "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress: "B2910E22Bb23D129C02d122B77B462ceB0E89Db9",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgCancelFeeSplit), res)
			},
		},
		// error
		{
			msgstr: "123",
			res: types.MsgCancelFeeSplit{
				ContractAddress: "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress: "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: `{"contract_address":"0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414","deployer_address":"B2910E22Bb23D129C02d122B77B462ceB0E89Db9"}`,
			res: types.MsgCancelFeeSplit{
				ContractAddress: "0xA4FFCda536CC8fF1eeFe32D32EE943b9B4e70414",
				DeployerAddress: "889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E",
			},
			fnCheck: func(msg sdk.Msg, err error, res types.MsgCancelFeeSplit) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertCancelFeeSplitMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}
