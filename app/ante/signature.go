package ante

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

var (
	serverConfigOnce   = sync.Once{}
	currentNodeKeyOnce = sync.Once{}
	currentNodePub     crypto.PubKey
	currentNodePriv    crypto.PrivKey
	serverConfig       *cfg.Config
	effectiveHeight    int64
)

// CreateAppCallback return the struct carry the callbacks
func CreateAppCallback(cdc *codec.Codec) server.AppCallback {
	return server.AppCallback{
		MempoolTxSignatureNodeKeysSetter: SetCurrentNodeKeys,
		ServerConfigCallback:             SetServerConfig,
	}
}

// SetCurrentNodeKeys used in the BaseApp to set the node keys
func SetCurrentNodeKeys(pub crypto.PubKey, priv crypto.PrivKey) {
	currentNodeKeyOnce.Do(func() {
		currentNodePriv = priv
		currentNodePub = pub
	})
}

// SetServerConfig use the callback to set the server config reference
func SetServerConfig(cfg *cfg.Config) {
	serverConfigOnce.Do(func() {
		serverConfig = cfg
	})
}

// SetServerConfigTest only used for test
func SetServerConfigTest(cfg *cfg.Config) {
	serverConfig = cfg
}

// SetWrappedTxEffectiveHeight set the effective height
func SetWrappedTxEffectiveHeight(height int64) {
	effectiveHeight = height
}

// use current config to verify the signature with the tx bytes
func VerifyConfidentTx(message, signature, pub []byte) (confident bool, err error) {
	pubKey := ed25519.PubKeyEd25519{}
	err = pubKey.UnmarshalFromAmino(pub)
	if err != nil {
		return
	}
	if pubKey.VerifyBytes(message, signature) {
		confidents := getConfidntNodeKeys()
		for _, v := range confidents {
			if v.Equals(pubKey) {
				confident = true
				return
			}
		}
	} else {
		err = errors.New("can not verify the signature")
	}
	return
}

// init and return current node keys
func getCurrentNodeKey() (crypto.PrivKey, crypto.PubKey) {
	return currentNodePriv, currentNodePub
}

// get the confident keys from the config
func getConfidntNodeKeys() []ed25519.PubKeyEd25519 {
	keys, _ := serverConfig.Mempool.GetCondifentNodeKeys()
	res := []ed25519.PubKeyEd25519{}
	for _, v := range keys {
		slice, e := hexutil.Decode(v)
		if e != nil {
			continue
		}
		key := ed25519.PubKeyEd25519{}
		e = key.UnmarshalFromAmino(slice)
		if e != nil {
			continue
		}
		res = append(res, key)
	}
	return res
}

// return if skip the wrapped logic
func isSkipWrapped(height int64) bool {
	if height > effectiveHeight {
		return len(serverConfig.Mempool.ConfidentNodeKeys) <= 0
	}
	return false
}

// StdTxSignatureWrapperDecorator for replacing the wrapped tx with current tx
type StdTxSignatureWrapperDecorator struct {
	cdc *codec.Codec
}

func NewStdTxSignatureWrapperDecorator(cdc *codec.Codec) StdTxSignatureWrapperDecorator {
	return StdTxSignatureWrapperDecorator{
		cdc: cdc,
	}
}

func (decorator StdTxSignatureWrapperDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !isSkipWrapped(ctx.BlockHeight()) {
		wrapped := app.NewWrappedTx(tx, app.StdTransaction)
		priv, pub := getCurrentNodeKey()
		signature, err := priv.Sign(ctx.TxBytes())
		if err != nil {
			return next(ctx, tx, simulate)
		}
		wrapped = wrapped.WithSignature(signature, pub.Bytes())
		result, err := decorator.cdc.MarshalBinaryLengthPrefixed(wrapped)
		if err != nil {
			return next(ctx, tx, simulate)
		}
		ctx = ctx.WithReplaceTx(result)
		return next(ctx, tx, simulate)
	}
	return next(ctx, tx, simulate)
}

// EthereumTxSignatureWrapperDecorator for replacing the wrapped tx with current tx
type EthereumTxSignatureWrapperDecorator struct {
	cdc *codec.Codec
}

func NewEthereumTxSignatureWrapperDecorator(cdc *codec.Codec) EthereumTxSignatureWrapperDecorator {
	return EthereumTxSignatureWrapperDecorator{
		cdc: cdc,
	}
}

func (decorator EthereumTxSignatureWrapperDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !isSkipWrapped(ctx.BlockHeight()) {
		wrapped := app.NewWrappedTx(tx, app.EthereumTransaction)
		priv, pub := getCurrentNodeKey()
		signature, err := priv.Sign(ctx.TxBytes())
		if err != nil {
			return next(ctx, tx, simulate)
		}
		wrapped = wrapped.WithSignature(signature, pub.Bytes())
		result, err := decorator.cdc.MarshalBinaryLengthPrefixed(wrapped)
		if err != nil {
			return next(ctx, tx, simulate)
		}
		ctx = ctx.WithReplaceTx(result)
		return next(ctx, tx, simulate)
	}
	return next(ctx, tx, simulate)
}

// AnonymousDecorator used to wrapp raw ante handler to decorator
type AnonymousDecorator struct {
	handler sdk.AnteHandler
}

func NewAnonymousDecorator(handler sdk.AnteHandler) AnonymousDecorator {
	return AnonymousDecorator{
		handler: handler,
	}
}

func (decorator AnonymousDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	netCtx, err := decorator.handler(ctx, tx, simulate)
	if err != nil {
		return netCtx, err
	}
	return next(netCtx, tx, simulate)
}

// WrappedTxVerifyDecorator verify the wrapped tx
// this antehandler should raise before derive decorator
type WrappedTxVerifyDecorator struct {
	cdc              *codec.Codec
	confidentHandler sdk.AnteHandler
	commonHandler    sdk.AnteHandler
}

func NewWrappedTxVerifyDecorator(cdc *codec.Codec, confidentHandler, commonAnteHandler sdk.AnteHandler) WrappedTxVerifyDecorator {
	return WrappedTxVerifyDecorator{
		cdc:              cdc,
		confidentHandler: confidentHandler,
		commonHandler:    commonAnteHandler,
	}
}

func (decorator WrappedTxVerifyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	defaultChain := NewAnonymousDecorator(next)
	handler := sdk.ChainAnteDecorators(NewAnonymousDecorator(decorator.commonHandler), defaultChain)
	wrappedTx, ok := tx.(app.WrappedTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}
	if isSkipWrapped(ctx.BlockHeight()) {
		//return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "don't support wrapped tx currenly: %T", tx)
		return handler(ctx, wrappedTx.GetOriginTx(), simulate) // SKIP ?
	}

	message, err := decorator.cdc.MarshalBinaryLengthPrefixed(wrappedTx.GetOriginTx())
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid wrapped origin tx: %v", wrappedTx.GetOriginTx())
	}

	confident, err := VerifyConfidentTx(message, wrappedTx.Signature, wrappedTx.NodeKey)
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid wrapped signature : %v", wrappedTx)
	}
	ctx = ctx.WithConfident(confident, wrappedTx.Type)
	origin := wrappedTx.GetOriginTx() // extract the origin tx to pass in the next chain
	if confident {
		handler = sdk.ChainAnteDecorators(NewAnonymousDecorator(decorator.confidentHandler), defaultChain)
	}

	return handler(ctx, origin, simulate)
}

// WrappedTxDeriveFromOriginDecorator as the final ante handler to wrap the origin tx or
// no confident tx to wrapped tx
type WrappedTxDeriveFromOriginDecorator struct {
	cdc *codec.Codec
}

func NewWrappedTxDeriveFromOriginDecorator(cdc *codec.Codec) WrappedTxDeriveFromOriginDecorator {
	return WrappedTxDeriveFromOriginDecorator{
		cdc: cdc,
	}
}

func (decorator WrappedTxDeriveFromOriginDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// NOTE: now this tx is the origin tx
	if ctx.Confident() {
		// verified
		newCtx = ctx.WithReplaceTx(nil) // keep safe
		return next(newCtx, tx, simulate)
	} else {
		if !isSkipWrapped(ctx.BlockHeight()) {
			wrappedTx := app.NewWrappedTx(tx, ctx.OriginTxType())
			priv, pub := getCurrentNodeKey()
			message, _ := decorator.cdc.MarshalBinaryLengthPrefixed(tx)
			signature, err := priv.Sign(message)
			if err == nil {
				wrappedTx.Signature = signature
				wrappedTx.NodeKey = pub.Bytes()
				message, _ := decorator.cdc.MarshalBinaryLengthPrefixed(wrappedTx)
				newCtx = ctx.WithReplaceTx(message)
				return next(newCtx, tx, simulate)
			}
		}
		return next(ctx, tx, simulate)
	}
}
