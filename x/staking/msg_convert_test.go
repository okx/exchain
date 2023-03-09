package staking

import (
	"fmt"
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/staking/types"
	"github.com/stretchr/testify/require"
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

func TestConvertDepositMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)

	testcases := []struct {
		msgstr  string
		res     types.MsgDeposit
		fnCheck func(msg sdk.Msg, err error, res types.MsgDeposit)
	}{
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"%s\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"1000\"}}", addr.String()),
			res:    NewMsgDeposit(testMustAccAddressFromBech32(addr.String()), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgDeposit), res)
			},
		},
		{
			msgstr: `{"delegator_address": "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1000"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgDeposit), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1000"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgDeposit), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1.5"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(15, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgDeposit), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "0.5"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgDeposit), res)
			},
		},
		// error
		{
			msgstr: "123",
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"0.5\"}}"),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"0.5\"}}"),
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgDeposit) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertDepositMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}

func TestConvertWithdrawMsg(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("B2910E22Bb23D129C02d122B77B462ceB0E89Db9")
	require.NoError(t, err)
	testcases := []struct {
		msgstr  string
		res     types.MsgWithdraw
		fnCheck func(msg sdk.Msg, err error, res types.MsgWithdraw)
	}{
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"%s\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"1000\"}}", addr.String()),
			res:    NewMsgWithdraw(testMustAccAddressFromBech32(addr.String()), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdraw), res)
			},
		},
		{
			msgstr: `{"delegator_address": "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1000"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdraw), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1000"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdraw), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "1.5"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(15, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdraw), res)
			},
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okb","amount": "0.5"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgWithdraw), res)
			},
		},
		// error
		{
			msgstr: "123",
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"0.5\"}}"),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E\",\"quantity\": {\"denom\": \"okb\",\"amount\": \"0.5\"}}"),
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgWithdraw) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertWithdrawMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}

func testMustValAddressFromBech32(addrs ...string) []sdk.ValAddress {
	var results []sdk.ValAddress
	for _, addr := range addrs {
		re, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		results = append(results, re)
	}
	return results
}

func TestConvertAddSharesMsg(t *testing.T) {
	valAddr1, err := sdk.ValAddressFromHex("07a277f15a4fa6bb6629ee25b24fb28579bf8e2a")
	require.NoError(t, err)
	valAddr2, err := sdk.ValAddressFromHex("422f2e2e38c34fd23c4de0a5aaddc3ca926817ed")
	require.NoError(t, err)

	testcases := []struct {
		msgstr  string
		res     types.MsgAddShares
		fnCheck func(msg sdk.Msg, err error, res types.MsgAddShares)
	}{
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"0xb2910e22bb23d129c02d122b77b462ceb0e89db9\",\"validator_addresses\": [\"%s\",\"%s\"]}", valAddr1.String(), valAddr2.String()),
			res: NewMsgAddShares(testMustAccAddressFromBech32("0xb2910e22bb23d129c02d122b77b462ceb0e89db9"),
				testMustValAddressFromBech32(valAddr1.String(), valAddr2.String())),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgAddShares) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgAddShares), res)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"0xb2910e22bb23d129c02d122b77b462ceb0e89db9\",\"validator_addresses\": [\"%s\"]}", valAddr1.String()),
			res: NewMsgAddShares(testMustAccAddressFromBech32("0xb2910e22bb23d129c02d122b77b462ceb0e89db9"),
				testMustValAddressFromBech32(valAddr1.String())),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgAddShares) {
				require.NoError(t, err)
				require.Equal(t, msg.(types.MsgAddShares), res)
			},
		},
		// error
		{
			msgstr: "123",
			fnCheck: func(msg sdk.Msg, err error, res types.MsgAddShares) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"\",\"validator_addresses\": [\"%s\",\"%s\"]}", valAddr1.String(), valAddr2.String()),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgAddShares) {
				require.Error(t, err)
				require.Nil(t, msg)
			},
		},
		{
			msgstr: fmt.Sprintf("{\"delegator_address\": \"0xb2910e22bb23d129c02d122b77b462ceb0e89db9\",\"validator_addresses\": [\"%s\",\"%s\"]}", valAddr1.String(), valAddr2.String()),
			res: NewMsgAddShares(testMustAccAddressFromBech32("0x889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E"),
				testMustValAddressFromBech32(valAddr1.String(), valAddr2.String())),
			fnCheck: func(msg sdk.Msg, err error, res types.MsgAddShares) {
				require.Equal(t, ErrCheckSignerFail, err)
				require.Nil(t, msg)
			},
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertAddSharesMsg([]byte(ts.msgstr), ts.res.GetSigners())
		ts.fnCheck(msg, err, ts.res)
	}
}
