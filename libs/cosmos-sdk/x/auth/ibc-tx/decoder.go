package ibc_tx

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/unknownproto"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	"google.golang.org/protobuf/encoding/protowire"
	//"github.com/okex/exchain/libs/cosmos-sdk/codec/unknownproto"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	ibckey "github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/ibc-key"
	tx "github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
)

// DefaultTxDecoder returns a default protobuf TxDecoder using the provided Marshaler.
//func IbcTxDecoder(cdc codec.ProtoCodecMarshaler) ibcadapter.TxDecoder {
func IbcTxDecoder(cdc codec.ProtoCodecMarshaler) ibctx.IbcTxDecoder {
	return func(txBytes []byte) (*authtypes.IbcTx, error) {
		// Make sure txBytes follow ADR-027.
		err := rejectNonADR027TxRaw(txBytes)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		var raw tx.TxRaw

		// reject all unknown proto fields in the root TxRaw
		err = unknownproto.RejectUnknownFieldsStrict(txBytes, &raw, cdc.InterfaceRegistry())
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		err = cdc.UnmarshalBinaryBare(txBytes, &raw)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		var body tx.TxBody
		// allow non-critical unknown fields in TxBody
		// txBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true, cdc.InterfaceRegistry())
		// if err != nil {
		// 	//Ywmet todo couldnot decode
		// 	//return authtypes.StdTx{}, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		// }

		err = cdc.UnmarshalBinaryBare(raw.BodyBytes, &body)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		var authInfo tx.AuthInfo

		// reject all unknown proto fields in AuthInfo
		err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo, cdc.InterfaceRegistry())
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		err = cdc.UnmarshalBinaryBare(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		ibcTx := &tx.Tx{
			Body:       &body,
			AuthInfo:   &authInfo,
			Signatures: raw.Signatures,
		}
		fee, signFee, err := convertFee(authInfo)
		if err != nil {
			return nil, err
		}

		signatures := convertSignature(cdc, ibcTx)

		// construct Msg
		stdMsgs, signMsgs, err := constructMsgs(ibcTx)
		if err != nil {
			return nil, err
		}

		var modeInfo *tx.ModeInfo_Single_
		var ok bool
		if len(authInfo.SignerInfos) > 0 {
			modeInfo, ok = authInfo.SignerInfos[0].ModeInfo.Sum.(*tx.ModeInfo_Single_)
			if !ok {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "only support ModeInfo_Single")
			}
		}
		var signMode signing.SignMode
		if modeInfo != nil && modeInfo.Single != nil {
			signMode = modeInfo.Single.Mode
		}

		stx := authtypes.IbcTx{
			&authtypes.StdTx{
				Msgs:       stdMsgs,
				Fee:        fee,
				Signatures: signatures,
				Memo:       ibcTx.Body.Memo,
			},
			raw.AuthInfoBytes,
			raw.BodyBytes,
			signMode,
			signFee,
			signMsgs,
		}

		return &stx, nil
	}
}

func constructMsgs(ibcTx *tx.Tx) ([]sdk.Msg, []sdk.Msg, error) {
	var err error
	stdMsgs, signMsgs := []sdk.Msg{}, []sdk.Msg{}
	for _, ibcMsg := range ibcTx.Body.Messages {
		m, ok := ibcMsg.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil, nil, sdkerrors.Wrap(
				sdkerrors.ErrInternal, "messages in ibcTx.Body not implement sdk.Msg",
			)
		}
		var newMsg sdk.Msg
		switch msg := m.(type) {
		case DenomAdapterMsg:
			// ibc transfer okt is not allowed,should do filter
			newMsg, err = msg.RulesFilter()
			if err != nil {
				return nil, nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "ibc tx decoder not support okt amount")
			}
		default:
			newMsg = m
		}
		stdMsgs = append(stdMsgs, newMsg)
		signMsgs = append(signMsgs, m)
	}
	return stdMsgs, signMsgs, nil
}

func convertSignature(cdc codec.ProtoCodecMarshaler, ibcTx *tx.Tx) []authtypes.StdSignature {
	signatures := []authtypes.StdSignature{}
	for i, s := range ibcTx.Signatures {
		pk := &ibckey.PubKey{}
		if ibcTx.AuthInfo.SignerInfos != nil {
			cdc.UnmarshalBinaryBare(ibcTx.AuthInfo.SignerInfos[i].PublicKey.Value, pk)
		}

		//convert crypto pubkey to tm pubkey
		tmPubKey := tmtypes.PubKeySecp256k1{}
		copy(tmPubKey[:], pk.Bytes())
		signatures = append(signatures,
			authtypes.StdSignature{
				Signature: s,
				PubKey:    tmPubKey,
			},
		)
	}
	return signatures
}

func convertFee(authInfo tx.AuthInfo) (authtypes.StdFee, authtypes.IbcFee, error) {

	gaslimit := uint64(0)
	var decCoins sdk.DecCoins
	var err error
	// for verify signature
	var signFee authtypes.IbcFee
	if authInfo.Fee != nil {
		decCoins, err = feeDenomFilter(authInfo.Fee.Amount)
		if err != nil {
			return authtypes.StdFee{}, authtypes.IbcFee{}, err
		}
		gaslimit = authInfo.Fee.GasLimit
		signFee = authtypes.IbcFee{
			authInfo.Fee.Amount,
			authInfo.Fee.GasLimit,
		}
	}

	return authtypes.StdFee{
		Amount: decCoins,
		Gas:    gaslimit,
	}, signFee, nil
}

func feeDenomFilter(coins sdk.CoinAdapters) (sdk.DecCoins, error) {
	decCoins := sdk.DecCoins{}

	if coins != nil {
		for _, fee := range coins {
			amount := fee.Amount.BigInt()
			denom := fee.Denom
			// convert ibc denom to DefaultBondDenom
			if denom == sdk.DefaultIbcWei {
				decCoins = append(decCoins, sdk.DecCoin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewDecFromIntWithPrec(sdk.NewIntFromBigInt(amount), sdk.Precision),
				})
			} else {
				// not suport other denom fee
				return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "ibc tx decoder only support wei fee")
			}
		}
	}
	return decCoins, nil
}

// DefaultJSONTxDecoder returns a default protobuf JSON TxDecoder using the provided Marshaler.
//func DefaultJSONTxDecoder(cdc codec.ProtoCodecMarshaler) sdk.TxDecoder {
//	return func(txBytes []byte) (sdk.Tx, error) {
//		var theTx tx.Tx
//		err := cdc.UnmarshalJSON(txBytes, &theTx)
//		if err != nil {
//			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
//		}
//
//		return &wrapper{
//			tx: &theTx,
//		}, nil
//	}
//}

// rejectNonADR027TxRaw rejects txBytes that do not follow ADR-027. This is NOT
// a generic ADR-027 checker, it only applies decoding TxRaw. Specifically, it
// only checks that:
// - field numbers are in ascending order (1, 2, and potentially multiple 3s),
// - and varints are as short as possible.
// All other ADR-027 edge cases (e.g. default values) are not applicable with
// TxRaw.
func rejectNonADR027TxRaw(txBytes []byte) error {
	// Make sure all fields are ordered in ascending order with this variable.
	prevTagNum := protowire.Number(0)

	for len(txBytes) > 0 {
		tagNum, wireType, m := protowire.ConsumeTag(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// TxRaw only has bytes fields.
		if wireType != protowire.BytesType {
			return fmt.Errorf("expected %d wire type, got %d", protowire.BytesType, wireType)
		}
		// Make sure fields are ordered in ascending order.
		if tagNum < prevTagNum {
			return fmt.Errorf("txRaw must follow ADR-027, got tagNum %d after tagNum %d", tagNum, prevTagNum)
		}
		prevTagNum = tagNum

		// All 3 fields of TxRaw have wireType == 2, so their next component
		// is a varint, so we can safely call ConsumeVarint here.
		// Byte structure: <varint of bytes length><bytes sequence>
		// Inner  fields are verified in `DefaultTxDecoder`
		lengthPrefix, m := protowire.ConsumeVarint(txBytes[m:])
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// We make sure that this varint is as short as possible.
		n := varintMinLength(lengthPrefix)
		if n != m {
			return fmt.Errorf("length prefix varint for tagNum %d is not as short as possible, read %d, only need %d", tagNum, m, n)
		}

		// Skip over the bytes that store fieldNumber and wireType bytes.
		_, _, m = protowire.ConsumeField(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		txBytes = txBytes[m:]
	}

	return nil
}

// varintMinLength returns the minimum number of bytes necessary to encode an
// uint using varint encoding.
func varintMinLength(n uint64) int {
	switch {
	// Note: 1<<N == 2**N.
	case n < 1<<(7):
		return 1
	case n < 1<<(7*2):
		return 2
	case n < 1<<(7*3):
		return 3
	case n < 1<<(7*4):
		return 4
	case n < 1<<(7*5):
		return 5
	case n < 1<<(7*6):
		return 6
	case n < 1<<(7*7):
		return 7
	case n < 1<<(7*8):
		return 8
	case n < 1<<(7*9):
		return 9
	default:
		return 10
	}
}
