package staking

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
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

	testcases := []struct {
		msgstr string
		res    types.MsgDeposit
	}{
		{
			msgstr: `{"delegator_address": "cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1.5"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(15, 1)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "0.5"}}`,
			res:    NewMsgDeposit(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertDepositMsg([]byte(ts.msgstr), ts.res.GetSigners())
		require.NoError(t, err)
		require.Equal(t, msg.(types.MsgDeposit), ts.res)
	}
}

func TestConvertWithdrawMsg(t *testing.T) {
	testcases := []struct {
		msgstr string
		res    types.MsgWithdraw
	}{
		{
			msgstr: `{"delegator_address": "cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("0xB2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1000"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(1000, 0)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "1.5"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(15, 1)),
		},
		{
			msgstr: `{"delegator_address": "B2910E22Bb23D129C02d122B77B462ceB0E89Db9","quantity": {"denom": "okt","amount": "0.5"}}`,
			res:    NewMsgWithdraw(testMustAccAddressFromBech32("B2910E22Bb23D129C02d122B77B462ceB0E89Db9"), newTestSysCoin(5, 1)),
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertWithdrawMsg([]byte(ts.msgstr), ts.res.GetSigners())
		require.NoError(t, err)
		require.Equal(t, msg.(types.MsgWithdraw), ts.res)
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

func testMustValAddressFromHex(addrs ...string) []sdk.ValAddress {
	var results []sdk.ValAddress
	for _, addr := range addrs {
		re, err := sdk.ValAddressFromHex(addr)
		if err != nil {
			panic(err)
		}
		results = append(results, re)
	}
	return results
}

func TestConvertAddSharesMsg(t *testing.T) {
	testcases := []struct {
		msgstr string
		res    types.MsgAddShares
	}{
		{
			msgstr: `{"delegator_address": "0xb2910e22bb23d129c02d122b77b462ceb0e89db9","validator_addresses": ["cosmosvaloper1q7380u26f7ntke3facjmynajs4umlr32qchq7w","cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj"]}`,
			res: NewMsgAddShares(testMustAccAddressFromBech32("0xb2910e22bb23d129c02d122b77b462ceb0e89db9"),
				testMustValAddressFromBech32("cosmosvaloper1q7380u26f7ntke3facjmynajs4umlr32qchq7w", "cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj")),
		},
		{
			msgstr: `{"delegator_address": "0xb2910e22bb23d129c02d122b77b462ceb0e89db9","validator_addresses": ["cosmosvaloper1q7380u26f7ntke3facjmynajs4umlr32qchq7w"]}`,
			res: NewMsgAddShares(testMustAccAddressFromBech32("0xb2910e22bb23d129c02d122b77b462ceb0e89db9"),
				testMustValAddressFromBech32("cosmosvaloper1q7380u26f7ntke3facjmynajs4umlr32qchq7w")),
		},
	}

	for _, ts := range testcases {
		msg, err := ConvertAddSharesMsg([]byte(ts.msgstr), ts.res.GetSigners())
		require.NoError(t, err)
		require.Equal(t, msg.(types.MsgAddShares), ts.res)
	}
}
