package types

import (
	"math/big"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types3 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	types4 "github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/testing/mock"
	helpers2 "github.com/okex/exchain/libs/ibc-go/testing/simapp/helpers"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types/testdata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

const (
	TransferPort   = "transfer"
	FirstChannelId = "channel-0"
)

var (
	priv = ed25519.GenPrivKey()
)

func TestIbcTxDecoderSignMode(t *testing.T) {

	// keys and addresses
	priv, _, addr := types.KeyTestPubAddr()

	//addrs := []sdk.AccAddress{addr}
	packet := channeltypes.NewPacket([]byte(mock.MockPacketData), 1,
		"transfer", "channel-0",
		"transfer", "channel-1",
		clienttypes.NewHeight(1, 0), 0)
	msgs := []ibcmsg.Msg{channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), addr.String())}

	interfaceRegistry := types3.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := ibc_tx.NewTxConfig(marshaler, ibc_tx.DefaultSignModes)

	// multi mode should no error
	require.Panics(t, func() {
		helpers2.GenTx(
			txConfig,
			msgs,
			sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
			helpers.DefaultGenTxGas,
			"exchain-101",
			[]uint64{0}, //[]uint64{acc.GetAccountNumber()},
			[]uint64{0}, //[]uint64{acc.GetSequence()},
			2,
			priv,
		)
	})
	// single mode should no error
	_, err := helpers2.GenTx(
		txConfig,
		msgs,
		sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
		helpers.DefaultGenTxGas,
		"exchain-101",
		[]uint64{0}, //[]uint64{acc.GetAccountNumber()},
		[]uint64{0}, //[]uint64{acc.GetSequence()},
		1,
		priv,
	)
	require.NoError(t, err)
}

// TestTxDecode decode ibc tx with unkown field
func TestIbcDecodeUnknownFields(t *testing.T) {
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
	cdcProxy := newProxyDecoder()
	decoder := TxDecoder(cdcProxy)

	tests := []struct {
		name           string
		body           *testdata.TestUpdatedTxBody
		authInfo       *testdata.TestUpdatedAuthInfo
		shouldErr      bool
		shouldAminoErr string
	}{
		{
			name: "no new fields should pass",
			body: &testdata.TestUpdatedTxBody{
				Memo: "foo",
			},
			authInfo:  &testdata.TestUpdatedAuthInfo{},
			shouldErr: false,
		},
		{
			name: "critical fields in AuthInfo should error on decode",
			body: &testdata.TestUpdatedTxBody{
				Memo: "foo",
			},
			authInfo: &testdata.TestUpdatedAuthInfo{
				NewField_3: []byte("xyz"),
			},
			shouldErr: true,
		},
		{
			name: "non-critical fields in AuthInfo should error on decode",
			body: &testdata.TestUpdatedTxBody{
				Memo: "foo",
			},
			authInfo: &testdata.TestUpdatedAuthInfo{
				NewField_1024: []byte("xyz"),
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bodyBz, err := tt.body.Marshal()
			require.NoError(t, err)

			authInfoBz, err := tt.authInfo.Marshal()
			require.NoError(t, err)

			txRaw := &tx.TxRaw{
				BodyBytes:     bodyBz,
				AuthInfoBytes: authInfoBz,
			}
			txBz, err := txRaw.Marshal()
			require.NoError(t, err)

			_, err = decoder(txBz)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Log("test TxRaw no new fields, should succeed")
	txRaw := &testdata.TestUpdatedTxRaw{
		BodyBytes: []byte("1"),
	}
	txBz, err := txRaw.Marshal()
	require.NoError(t, err)
	_, err = decoder(txBz)
	require.Error(t, err)

	t.Log("new field in TxRaw should fail")
	txRaw = &testdata.TestUpdatedTxRaw{
		NewField_5: []byte("abc"),
	}
	txBz, err = txRaw.Marshal()
	require.NoError(t, err)
	_, err = decoder(txBz)
	require.Error(t, err)

	//
	t.Log("new \"non-critical\" field in TxRaw should fail")
	txRaw = &testdata.TestUpdatedTxRaw{
		NewField_1024: []byte("abc"),
	}
	txBz, err = txRaw.Marshal()
	require.NoError(t, err)
	_, err = decoder(txBz)
	require.Error(t, err)
}

func TestRejectNonADR027(t *testing.T) {
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
	cdcProxy := newProxyDecoder()
	decoder := TxDecoder(cdcProxy)

	body := &testdata.TestUpdatedTxBody{Memo: "AAA"} // Look for "65 65 65" when debugging the bytes stream.
	bodyBz, err := body.Marshal()
	require.NoError(t, err)
	authInfo := &testdata.TestUpdatedAuthInfo{Fee: &tx.Fee{GasLimit: 127}} // Look for "127" when debugging the bytes stream.
	authInfoBz, err := authInfo.Marshal()
	txRaw := &tx.TxRaw{
		BodyBytes:     bodyBz,
		AuthInfoBytes: authInfoBz,
		Signatures:    [][]byte{{41}, {42}, {43}}, // Look for "42" when debugging the bytes stream.
	}

	// We know these bytes are ADR-027-compliant.
	txBz, err := txRaw.Marshal()

	// From the `txBz`, we extract the 3 components:
	// bodyBz, authInfoBz, sigsBz.
	// In our tests, we will try to decode txs with those 3 components in all
	// possible orders.
	//
	// Consume "BodyBytes" field.
	_, _, m := protowire.ConsumeField(txBz)
	bodyBz = append([]byte{}, txBz[:m]...)
	txBz = txBz[m:] // Skip over "BodyBytes" bytes.
	// Consume "AuthInfoBytes" field.
	_, _, m = protowire.ConsumeField(txBz)
	authInfoBz = append([]byte{}, txBz[:m]...)
	txBz = txBz[m:] // Skip over "AuthInfoBytes" bytes.
	// Consume "Signature" field, it's the remaining bytes.
	sigsBz := append([]byte{}, txBz...)

	// bodyBz's length prefix is 5, with `5` as varint encoding. We also try a
	// longer varint encoding for 5: `133 00`.
	longVarintBodyBz := append(append([]byte{bodyBz[0]}, byte(133), byte(00)), bodyBz[2:]...)

	tests := []struct {
		name      string
		txBz      []byte
		shouldErr bool
	}{
		{
			"authInfo, body, sigs",
			append(append(authInfoBz, bodyBz...), sigsBz...),
			true,
		},
		{
			"authInfo, sigs, body",
			append(append(authInfoBz, sigsBz...), bodyBz...),
			true,
		},
		{
			"sigs, body, authInfo",
			append(append(sigsBz, bodyBz...), authInfoBz...),
			true,
		},
		{
			"sigs, authInfo, body",
			append(append(sigsBz, authInfoBz...), bodyBz...),
			true,
		},
		{
			"body, sigs, authInfo",
			append(append(bodyBz, sigsBz...), authInfoBz...),
			true,
		},
		{
			"body, authInfo, sigs (valid txRaw)",
			append(append(bodyBz, authInfoBz...), sigsBz...),
			false,
		},
		{
			"longer varint than needed",
			append(append(longVarintBodyBz, authInfoBz...), sigsBz...),
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err = decoder(tt.txBz)
			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHeightSensitive(t *testing.T) {
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
	cdcProxy := newProxyDecoder()
	decoder := TxDecoder(cdcProxy)

	msg := types4.NewMsgRegisterPayee("port_Id", "channel_id", "mock", "mock2")
	msgs := []ibcmsg.Msg{msg}

	interfaceRegistry := types3.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := ibc_tx.NewTxConfig(marshaler, ibc_tx.DefaultSignModes)

	priv, _, _ := types.KeyTestPubAddr()
	txBytes, err := helpers2.GenTxBytes(
		txConfig,
		msgs,
		sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
		helpers.DefaultGenTxGas,
		"exchain-101",
		[]uint64{0}, //[]uint64{acc.GetAccountNumber()},
		[]uint64{0}, //[]uint64{acc.GetSequence()},
		1,
		priv,
	)
	require.NoError(t, err)
	_, err = decoder(txBytes)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not support before")
}
