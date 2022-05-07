package client

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	signingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
)

type (
	// TxEncodingConfig defines an interface that contains transaction
	// encoders and decoders
	TxEncodingConfig interface {
		TxEncoder() ibcmsg.IBCTxEncoder
		TxDecoder() ibcmsg.IbcTxDecoder
		TxJSONEncoder() ibcmsg.IBCTxEncoder
		//TxJSONDecoder() sdk.TxDecoder
		MarshalSignatureJSON([]signing.SignatureV2) ([]byte, error)

		UnmarshalSignatureJSON([]byte) ([]signing.SignatureV2, error)
	}

	// TxConfig defines an interface a client can utilize to generate an
	// application-defined concrete transaction type. The type returned must
	// implement TxBuilder.
	TxConfig interface {
		TxEncodingConfig

		NewTxBuilder() TxBuilder
		WrapTxBuilder(tx ibcmsg.Tx) (TxBuilder, error)
		SignModeHandler() signingtypes.SignModeHandler
	}

	// TxBuilder defines an interface which an application-defined concrete transaction
	// type must implement. Namely, it must be able to set messages, generate
	// signatures, and provide canonical bytes to sign over. The transaction must
	// also know how to encode itself.
	TxBuilder interface {
		GetTx() signingtypes.Tx

		SetMsgs(msgs ...ibcmsg.Msg) error
		SetSignatures(signatures ...signing.SignatureV2) error
		SetMemo(memo string)
		SetFeeAmount(amount sdk.CoinAdapters)
		SetGasLimit(limit uint64)
		SetTimeoutHeight(height uint64)
		SetFeeGranter(feeGranter sdk.AccAddress)
	}
)
