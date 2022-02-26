package types

import (
	"errors"
	"fmt"
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

const IGNORE_HEIGHT_CHECKING = -1

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
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

		for _, f := range []decodeFunc{
			evmDecoder,
			ubruDecoder,
			ubDecoder,
			byteTx,
			relayTx,
		} {
			if tx, err = f(cdc, txBytes, height); err == nil {
				return tx, nil
			}
		}

		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

var byteTx decodeFunc = func(c *codec.Codec, bytes []byte, i int64) (sdk.Tx, error) {
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

var relayTx decodeFunc = func(c *codec.Codec, bytes []byte, i int64) (sdk.Tx, error) {
	wp := &sdk.RelayMsgWrapper{}
	err := wp.UnMarshal(bytes)
	if nil != err {
		return nil, err
	}
	msgs := make([]sdk.Msg, 0)
	addr, _ := sdk.AccAddressFromBech32ByPrefix("ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9", "ex")
	for _, v := range wp.Msgs {
		msgs = append(msgs, v)
		v.Singers[0] = addr
	}

	sis := make([]authtypes.StdSignature, 1)
	ret := authtypes.StdTx{
		Msgs:       msgs,
		Fee:        authtypes.StdFee{},
		Signatures: sis,
		Memo:       "okt",
	}
	return ret, nil
}

type decodeFunc func(*codec.Codec, []byte, int64) (sdk.Tx, error)

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {

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
func ubruDecoder(cdc *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err != nil {
		return nil, err
	}
	return sanityCheck(v.(sdk.Tx), height)
}

// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
// 3. the original amino way, decode by reflection.
func ubDecoder(cdc *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {
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
