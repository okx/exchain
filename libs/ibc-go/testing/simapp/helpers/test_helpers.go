package helpers

import (
	"math/rand"
	"time"

	ica "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"

	ibcfee "github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee"

	okexchaincodec "github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	capabilityModule "github.com/okex/exchain/libs/cosmos-sdk/x/capability"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	ibctransfer "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/x/icamauth"
)

// SimAppChainID hardcoded chainID for simulation
const (
	DefaultGenTxGas = 1000000
	SimAppChainID   = "simulation-app"
)

// GenTx generates a signed mock transaction.
func GenTx(gen client.TxConfig, msgs []ibcmsg.Msg, feeAmt sdk.CoinAdapters, gas uint64, chainID string, accNums, accSeqs []uint64, smode int, priv ...crypto.PrivKey) (sdk.Tx, error) {
	txBytes, err := GenTxBytes(gen, msgs, feeAmt, gas, chainID, accNums, accSeqs, smode, priv...)
	if nil != err {
		panic(err)
	}
	cdcProxy := newProxyDecoder()

	ibcTx, err := ibc_tx.IbcTxDecoder(cdcProxy.GetProtocMarshal())(txBytes)

	return ibcTx, err
}

func GenTxBytes(gen client.TxConfig, msgs []ibcmsg.Msg, feeAmt sdk.CoinAdapters, gas uint64, chainID string, accNums, accSeqs []uint64, smode int, priv ...crypto.PrivKey) ([]byte, error) {
	sigs := make([]signing.SignatureV2, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode := gen.SignModeHandler().DefaultMode()
	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	// 1 mode single
	// 2 mode multi
	switch smode {
	case 1:
		for i, p := range priv {
			pubKey := ibc_tx.LagacyKey2PbKey(p.PubKey())
			sigs[i] = signing.SignatureV2{
				PubKey: pubKey,
				Data: &signing.SingleSignatureData{
					SignMode: gen.SignModeHandler().DefaultMode(),
				},
				Sequence: accSeqs[i],
			}
		}
	case 2:
		//only support for ut
		keyLen := 10
		for i, p := range priv {
			pubKey := ibc_tx.LagacyKey2PbKey(p.PubKey())
			sigs[i] = signing.SignatureV2{
				PubKey: pubKey,
				Data: &signing.MultiSignatureData{
					BitArray:   types.NewCompactBitArray(keyLen),
					Signatures: make([]signing.SignatureData, 0, keyLen),
				},
				Sequence: accSeqs[i],
			}
		}
	}

	tx := gen.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = tx.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range priv {
		signerData := signing2.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
		if err != nil {
			panic(err)
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = tx.SetSignatures(sigs...)
		if err != nil {
			panic(err)
		}
	}
	return gen.TxEncoder()(tx.GetTx())
}

func newProxyDecoder() *codec.CodecProxy {
	ModuleBasics := module.NewBasicManager(
		bank.AppModuleBasic{},
		capabilityModule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		icamauth.AppModuleBasic{},
	)
	cdc := okexchaincodec.MakeCodec(ModuleBasics)
	interfaceReg := okexchaincodec.MakeIBC(ModuleBasics)
	protoCodec := codec.NewProtoCodec(interfaceReg)
	codecProxy := codec.NewCodecProxy(protoCodec, cdc)
	return codecProxy
}
