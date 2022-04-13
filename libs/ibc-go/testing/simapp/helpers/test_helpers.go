package helpers

import (
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"math/rand"
	"time"

	tmsecp256k1 "github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	//cryptotypes "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
)

// SimAppChainID hardcoded chainID for simulation
const (
	DefaultGenTxGas = 1000000
	SimAppChainID   = "simulation-app"
)

// GenTx generates a signed mock transaction.
func GenTx( /*gen client.TxConfig,*/ msgs []sdk.Msg, feeAmt sdk.Coins, gas uint64, chainID string, accNums, accSeqs []uint64, priv ...tmsecp256k1.PrivKeySecp256k1) (sdk.Tx, error) {
	sigs := make([]authtypes.StdSignature, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	//	for i, p := range priv {
	//		sigs[i] = signing.SignatureV2{
	//			PubKey: p.PubKey(),
	//			Data: &signing.SingleSignatureData{
	//				SignMode: signMode,
	//			},
	//			Sequence: accSeqs[i],
	//		}
	//	}
	fee := authtypes.StdFee{
		Amount: feeAmt,
		Gas:    gas,
	}
	for i, p := range priv {
		sig, err := p.Sign(authtypes.StdSignBytes(chainID, accNums[i], accSeqs[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}
		sigs[i] = authtypes.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return authtypes.NewStdTx(msgs, fee, sigs, memo), nil
	//	tx := gen.NewTxBuilder()
	//	err := tx.SetMsgs(msgs...)
	//	if err != nil {
	//		return nil, err
	//	}
	//	err = tx.SetSignatures(sigs...)
	//	if err != nil {
	//		return nil, err
	//	}
	//	tx.SetMemo(memo)
	//	tx.SetFeeAmount(feeAmt)
	//	tx.SetGasLimit(gas)
	//
	//	// 2nd round: once all signer infos are set, every signer can sign.
	//	for i, p := range priv {
	//		signerData := authsign.SignerData{
	//			ChainID:       chainID,
	//			AccountNumber: accNums[i],
	//			Sequence:      accSeqs[i],
	//		}
	//		signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
	//		if err != nil {
	//			panic(err)
	//		}
	//		sig, err := p.Sign(signBytes)
	//		if err != nil {
	//			panic(err)
	//		}
	//		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
	//		err = tx.SetSignatures(sigs...)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//
	//	return tx.GetTx(), nil
}
