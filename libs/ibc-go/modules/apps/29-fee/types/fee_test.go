package types_test

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/stretchr/testify/require"
)

var (
	// defaultRecvFee is the default packet receive fee used for testing purposes
	defaultRecvFee = sdk.CoinAdapters{sdk.CoinAdapter{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(100)}}

	// defaultAckFee is the default packet acknowledgement fee used for testing purposes
	defaultAckFee = sdk.CoinAdapters{sdk.CoinAdapter{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(200)}}

	// defaultTimeoutFee is the default packet timeout fee used for testing purposes
	defaultTimeoutFee = sdk.CoinAdapters{sdk.CoinAdapter{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(300)}}

	// invalidFee is an invalid coin set used to trigger error cases for testing purposes
	invalidFee = sdk.CoinAdapters{sdk.CoinAdapter{Denom: "invalid-denom", Amount: sdk.NewInt(-2)}}

	// defaultAccAddress is the default account used for testing purposes
	defaultAccAddress = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
)

func TestFeeTotal(t *testing.T) {
	fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)

	total := fee.Total()
	require.Equal(t, sdk.NewInt(600), total.AmountOf(sdk.DefaultBondDenom))
}

func TestPacketFeeValidation(t *testing.T) {
	var packetFee types.PacketFee

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
		{
			"success with empty slice for Relayers",
			func() {
				packetFee.Relayers = []string{}
			},
			true,
		},
		{
			"should fail when refund address is invalid",
			func() {
				packetFee.RefundAddress = "invalid-address"
			},
			false,
		},
		{
			"should fail when all fees are invalid",
			func() {
				packetFee.Fee.AckFee = invalidFee
				packetFee.Fee.RecvFee = invalidFee
				packetFee.Fee.TimeoutFee = invalidFee
			},
			false,
		},
		{
			"should fail with single invalid fee",
			func() {
				packetFee.Fee.AckFee = invalidFee
			},
			false,
		},
		{
			"should fail with two invalid fees",
			func() {
				packetFee.Fee.TimeoutFee = invalidFee
				packetFee.Fee.AckFee = invalidFee
			},
			false,
		},
		{
			"should pass with two empty fees",
			func() {
				packetFee.Fee.TimeoutFee = sdk.CoinAdapters{}
				packetFee.Fee.AckFee = sdk.CoinAdapters{}
			},
			true,
		},
		{
			"should pass with one empty fee",
			func() {
				packetFee.Fee.TimeoutFee = sdk.CoinAdapters{}
			},
			true,
		},
		{
			"should fail if all fees are empty",
			func() {
				packetFee.Fee.AckFee = sdk.CoinAdapters{}
				packetFee.Fee.RecvFee = sdk.CoinAdapters{}
				packetFee.Fee.TimeoutFee = sdk.CoinAdapters{}
			},
			false,
		},
		{
			"should fail with non empty Relayers",
			func() {
				packetFee.Relayers = []string{"relayer"}
			},
			false,
		},
	}

	for _, tc := range testCases {
		fee := types.NewFee(defaultRecvFee, defaultAckFee, defaultTimeoutFee)
		packetFee = types.NewPacketFee(fee, defaultAccAddress, nil)

		tc.malleate() // malleate mutates test data

		err := packetFee.Validate()

		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
