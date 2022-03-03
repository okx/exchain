package types

import (
	"errors"
	"fmt"

	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	typestx "github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

const IGNORE_HEIGHT_CHECKING = -1

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec, proxy ...*codec.CodecProxy) sdk.TxDecoder {
	return func(txBytes []byte, heights ...int64) (sdk.Tx, error) {
		if len(heights) > 1 {
			return nil, fmt.Errorf("to many height parameters")
		}
		var tx sdk.Tx
		var err error
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		var height int64
		if len(heights) == 1 {
			height = heights[0]
		} else {
			height = global.GetGlobalHeight()
		}
		proxyCodec := &codec.CodecProxy{}
		if len(proxy) > 0 {
			proxyCodec = proxy[0]
		}

		for _, f := range []decodeFunc{
			evmDecoder,
			ubruDecoder,
			ubDecoder,
			byteTx,
			relayTx,
		} {
			if tx, err = f(cdc, proxyCodec, txBytes, height); err == nil {
				return tx, nil
			}
		}

		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

var byteTx decodeFunc = func(c *codec.Codec, proxy *codec.CodecProxy, bytes []byte, i int64) (sdk.Tx, error) {
	bw := new(sdk.BytesWrapper)
	txBytes, err := bw.UnmarshalToTx(bytes)
	if nil != err {
		return nil, err
	}
	tt := new(auth.StdTx)
	err = c.UnmarshalJSON(txBytes, &tt)
	if len(tt.GetMsgs()) == 0 {
		return nil, errors.New("asd")
	}
	logrusplugin.Info("tx", "coins", fmt.Sprintf("%s", tt.GetFee()))
	//err = c.UnmarshalJSON(txBytes, &tt)
	return *tt, err
}

var relayTx decodeFunc = func(c *codec.Codec, proxy *codec.CodecProxy, bytes []byte, i int64) (sdk.Tx, error) {
	// wp := &sdk.RelayMsgWrapper{}
	// err := wp.UnMarshal(bytes)
	// if nil != err {
	// 	return nil, err
	// }
	// msgs := make([]sdk.Msg, 0)
	// addr, _ := sdk.AccAddressFromBech32ByPrefix("ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9", "ex")
	// for _, v := range wp.Msgs {
	// 	msgs = append(msgs, v)
	// 	v.Singers[0] = addr
	// }

	// sis := make([]authtypes.StdSignature, 1)
	// ret := authtypes.StdTx{
	// 	Msgs:       msgs,
	// 	Fee:        authtypes.StdFee{},
	// 	Signatures: sis,
	// 	Memo:       "okt",
	// }
	// return ret, nil
	simReq := &typestx.SimulateRequest{}
	var err error
	_ = &typestx.SimulateRequest{}
	//err := simReq.Unmarshal(bytes)
	// err := protoCodec.Unmarshal(req.Data, i)
	// if err != nil {
	// 	return err
	// }
	// if qrt.interfaceRegistry != nil {
	// 	return codectypes.UnpackInterfaces(i, qrt.interfaceRegistry)
	// }
	// 	return nil

	if err != nil {
		return authtypes.StdTx{}, err
	}

	txBytes := simReq.TxBytes
	if txBytes == nil && simReq.Tx != nil {
		// This block is for backwards-compatibility.
		// We used to support passing a `Tx` in req. But if we do that, sig
		// verification might not pass, because the .Marshal() below might not
		// be the same marshaling done by the client.
		var err error
		//txBytes, err = proto.Marshal(simReq.Tx)
		if err != nil {
			return nil, nil //status.Errorf(codes.InvalidArgument, "invalid tx; %v", err)
		}
	}

	if txBytes == nil {
		return nil, nil //status.Errorf(codes.InvalidArgument, "empty txBytes is not allowed")
	}
	//amino := ibccodec.NewLegacyAmino()
	// interfaceRegistry := codestypes.NewInterfaceRegistry() //codectypes.NewInterfaceRegistry()
	//marshaler := ibccodec.NewProtoCodec(interfaceRegistry)
	//decoder := authtx.DefaultTxDecoder(marshaler)
	if proxy != nil {
		// marshaler := proxy.GetProtocMarshal()
		// //evmtxtypes.DefaultTxDecoder(marshaler)
		// decode := sdk.DefaultTxDecoder(marshaler)
		// decode(txBytes)
	}

	////authtx.DefaultTxDecoder(marshaler)
	// ec := EncodingConfig{
	// 	InterfaceRegistry: interfaceRegistry,
	// 	Marshaler:         marshaler,
	// 	//Amino:             amino,
	// }

	// ibcstd.RegisterLegacyAminoCodec(ec.Amino)
	//ibcstd.RegisterInterfaces(ec.InterfaceRegistry)
	// ModuleBasics.RegisterLegacyAminoCodec(ec.Amino)
	// ModuleBasics.RegisterInterfaces(ec.InterfaceRegistry)

	// txz, err := decoder(txBytes)
	// if err != nil {
	// 	return authtypes.StdTx{}, nil
	// }
	// //wrappered tx
	// fmt.Println(txz)
	// msgs := txz.GetMsgs()
	// if err := validateBasicTxMsgs(msgs); err != nil {
	// 	return authtypes.StdTx{}, err
	// }

	return convertTx(simReq.Tx), nil
}

// func validateBasicTxMsgs(msgs []ibcsdk.Msg) error {
// 	if len(msgs) == 0 {
// 		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must contain at least one message")
// 	}

// 	for _, msg := range msgs {
// 		err := msg.ValidateBasic()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func convertTx(tx *ibctx.Tx) authtypes.StdTx {
	amount := tx.AuthInfo.Fee.Amount[0].Amount.BigInt()

	fee := authtypes.StdFee{
		Amount: []sdk.DecCoin{
			sdk.DecCoin{
				Denom:  tx.AuthInfo.Fee.Amount[0].Denom,
				Amount: sdk.NewDecFromBigInt(amount),
			},
		},
	}
	signature := []authtypes.StdSignature{}
	for _, s := range tx.Signatures {
		signature = append(signature, authtypes.StdSignature{Signature: s})
	}
	//tx.Body.Messages
	signers := []sdk.AccAddress{}
	ss := tx.GetSigners()
	for _, v := range ss {
		signers = append(signers, v.Bytes())
	}

	m := sdk.RelayMsg{
		TypeStr: tx.Body.Messages[0].TypeUrl,
		Bytes:   tx.Body.Messages[0].Value,
		//RouterStr:
		Singers: signers,
	}
	ms := []sdk.RelayMsg{m}

	msgs := []sdk.Msg{}
	for _, v := range ms {
		msgs = append(msgs, &v)
	}

	return authtypes.StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: signature,
		Memo:       tx.GetBody().GetMemo(),
	}
}

type decodeFunc func(*codec.Codec, *codec.CodecProxy, []byte, int64) (sdk.Tx, error)

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {

	// bypass height checking in case of a negative number
	if height >= 0 && !types.HigherThanVenus(height) {
		err = fmt.Errorf("lower than Venus")
		return
	}

	var ethTx MsgEthereumTx
	if err = authtypes.EthereumTxDecode(txBytes, &ethTx); err == nil {
		tx = ethTx
	}
	return
}

// 2. try customized unmarshalling implemented by UnmarshalFromAmino. higher performance!
func ubruDecoder(cdc *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err != nil {
		return nil, err
	}
	return sanityCheck(v.(sdk.Tx), height)
}

// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
// 3. the original amino way, decode by reflection.
func ubDecoder(cdc *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	return sanityCheck(tx, height)
}

func sanityCheck(tx sdk.Tx, height int64) (sdk.Tx, error) {
	if _, ok := tx.(MsgEthereumTx); ok && types.HigherThanVenus(height) {
		return nil, fmt.Errorf("amino decode is not allowed for MsgEthereumTx")
	}
	return tx, nil
}
