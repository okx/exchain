package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	signingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/client/input"
	"github.com/okex/exchain/libs/cosmos-sdk/client/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate" yaml:"gas_estimate"`
}

func (gr GasEstimateResponse) String() string {
	return fmt.Sprintf("gas estimate: %d", gr.GasEstimate)
}

// GenerateOrBroadcastMsgs creates a StdTx given a series of messages. If
// the provided context has generate-only enabled, the tx will only be printed
// to STDOUT in a fully offline manner. Otherwise, the tx will be signed and
// broadcasted.
func GenerateOrBroadcastMsgs(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
}

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTxCLI(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txConfig := NewPbTxConfig(cliCtx.InterfaceRegistry)
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil
	}
	txBytes := []byte{}
	pbtxMsgs, isPbTxMsg := convertIfPbTx(msgs)
	if !cliCtx.SkipConfirm {
		var signData interface{}
		var json []byte
		if isPbTxMsg {

			tx, err := buildUnsignedPbTx(txBldr, txConfig, pbtxMsgs...)
			if err != nil {
				return err
			}
			json, err = txConfig.TxJSONEncoder()(tx.GetTx())
			if err != nil {
				panic(err)
			}
		} else {
			signData, err = txBldr.BuildSignMsg(msgs)
			if err != nil {
				return err
			}

			if viper.GetBool(flags.FlagIndentResponse) {
				json, err = cliCtx.Codec.MarshalJSONIndent(signData, "", "  ")
				if err != nil {
					panic(err)
				}
			} else {
				json = cliCtx.Codec.MustMarshalJSON(signData)
			}
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}

	}

	if isPbTxMsg {
		txBytes, err = PbTxBuildAndSign(cliCtx, txConfig, txBldr, keys.DefaultKeyPass, pbtxMsgs)
		if err != nil {
			panic(err)
		}
	} else {
		// build and sign the transaction
		txBytes, err = txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
		if err != nil {
			return err
		}
	}
	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

func buildUnsignedPbTx(txf authtypes.TxBuilder, txConfig client.TxConfig, msgs ...txmsg.Msg) (client.TxBuilder, error) {
	if txf.ChainID() == "" {
		return nil, fmt.Errorf("chain ID required but not specified")
	}

	fees := txf.Fees()

	if !txf.GasPrices().IsZero() {
		if !fees.IsZero() {
			return nil, errors.New("cannot provide both fees and gas prices")
		}

		glDec := sdk.NewDec(int64(txf.Gas()))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		fees = make(sdk.Coins, len(txf.GasPrices()))

		for i, gp := range txf.GasPrices() {
			fee := gp.Amount.Mul(glDec)
			fees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}
	}

	tx := txConfig.NewTxBuilder()

	if err := tx.SetMsgs(msgs...); err != nil {
		return nil, err
	}
	tx.SetMemo(txf.Memo())
	coins := []sdk.CoinAdapter{}
	for _, fee := range txf.Fees() {
		prec := newCoinFromDec()

		am := sdk.NewIntFromBigInt(fee.Amount.BigInt().Div(fee.Amount.BigInt(), prec))

		coins = append(coins, sdk.NewCoinAdapter(fee.Denom, am))
	}
	tx.SetFeeAmount(coins)

	tx.SetGasLimit(txf.Gas())
	//tx.SetTimeoutHeight(txf.TimeoutHeight())

	return tx, nil
}

func newCoinFromDec() *big.Int {
	n := big.Int{}
	prec, ok := n.SetString(sdk.DefaultDecStr, 10)
	if !ok {
		panic(errors.New("newCoinFromDec setstring error"))
	}
	return prec
}

func PbTxBuildAndSign(clientCtx context.CLIContext, txConfig client.TxConfig, txbld authtypes.TxBuilder, passphrase string, msgs []txmsg.Msg) ([]byte, error) {
	//txb := txConfig.NewTxBuilder()
	txb, err := buildUnsignedPbTx(txbld, txConfig, msgs...)
	if err != nil {
		return nil, err
	}
	if !clientCtx.SkipConfirm {
		out, err := txConfig.TxJSONEncoder()(txb.GetTx())
		if err != nil {
			return nil, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	err = signPbTx(txConfig, txbld, clientCtx.GetFromName(), passphrase, &txb, true)
	if err != nil {
		return nil, err
	}

	return txConfig.TxEncoder()(txb.GetTx())
}

func signPbTx(txConfig client.TxConfig, txf authtypes.TxBuilder, name string, passwd string, pbTxBld *client.TxBuilder, overwriteSig bool) error {
	if txf.Keybase() == nil {
		return errors.New("keybase must be set prior to signing a transaction")
	}
	signMode := txConfig.SignModeHandler().DefaultMode()
	privKey, err := txf.Keybase().ExportPrivateKeyObject(name, passwd)
	if err != nil {
		return err
	}

	pubKeyPB := ibctx.LagacyKey2PbKey(privKey.PubKey())

	signerData := signingtypes.SignerData{
		ChainID:       txf.ChainID(),
		AccountNumber: txf.AccountNumber(),
		Sequence:      txf.Sequence(),
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}

	sig := signing.SignatureV2{
		PubKey:   pubKeyPB,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}
	var prevSignatures []signing.SignatureV2
	if !overwriteSig {
		prevSignatures, err = (*pbTxBld).GetTx().GetSignaturesV2()
		if err != nil {
			return err
		}
	}
	if err := (*pbTxBld).SetSignatures(sig); err != nil {
		return err
	}

	// Generate the bytes to be signed.
	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, (*pbTxBld).GetTx())
	if err != nil {
		return err
	}

	sigBytes, err := privKey.Sign(bytesToSign)
	if err != nil {
		panic(err)
	}
	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubKeyPB,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}

	if overwriteSig {
		return (*pbTxBld).SetSignatures(sig)
	}
	prevSignatures = append(prevSignatures, sig)

	return (*pbTxBld).SetSignatures(prevSignatures...)
}

func convertIfPbTx(msgs []sdk.Msg) ([]txmsg.Msg, bool) {
	retmsg := []txmsg.Msg{}
	for _, msg := range msgs {
		if m, ok := msg.(txmsg.Msg); ok {
			retmsg = append(retmsg, m)
		}
	}

	if len(retmsg) > 0 {
		return retmsg, true
	}
	return nil, false
}

// EnrichWithGas calculates the gas estimate that would be consumed by the
// transaction and set the transaction's respective value accordingly.
func EnrichWithGas(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (authtypes.TxBuilder, error) {
	_, adjusted, err := simulateMsgs(txBldr, cliCtx, msgs)
	if err != nil {
		return txBldr, err
	}

	return txBldr.WithGas(adjusted), nil
}

// CalculateGas simulates the execution of a transaction and returns
// the simulation response obtained by the query and the adjusted gas amount.
func CalculateGas(
	queryFunc func(string, []byte) ([]byte, int64, error), cdc *codec.Codec,
	txBytes []byte, adjustment float64,
) (sdk.SimulationResponse, uint64, error) {

	// run a simulation (via /app/simulate query) to
	// estimate gas and update TxBuilder accordingly
	rawRes, _, err := queryFunc("/app/simulate", txBytes)
	if err != nil {
		return sdk.SimulationResponse{}, 0, err
	}

	simRes, err := parseQueryResponse(cdc, rawRes)
	if err != nil {
		return sdk.SimulationResponse{}, 0, err
	}

	adjusted := adjustGasEstimate(simRes.GasUsed, adjustment)
	return simRes, adjusted, nil
}
func NewPbTxConfig(reg types2.InterfaceRegistry) client.TxConfig {
	marshaler := codec.NewProtoCodec(reg)
	return ibctx.NewTxConfig(marshaler, ibctx.DefaultSignModes)
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
func PrintUnsignedStdTx(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	stdTx, err := buildUnsignedStdTxOffline(txBldr, cliCtx, msgs)
	if err != nil {
		return err
	}

	var json []byte
	pbTxMsgs, isPbTxMsg := convertIfPbTx(msgs)

	if isPbTxMsg {
		txConfig := NewPbTxConfig(cliCtx.InterfaceRegistry)
		tx, err := buildUnsignedPbTx(txBldr, txConfig, pbTxMsgs...)
		if err != nil {
			return err
		}
		json, err = txConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return err
		}
	} else {
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdTx, "", "  ")
		} else {
			json, err = cliCtx.Codec.MarshalJSON(stdTx)
		}
		if err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintf(cliCtx.Output, "%s\n", json)
	return nil
}

// SignStdTx appends a signature to a StdTx and returns a copy of it. If appendSig
// is false, it replaces the signatures already attached with the new signature.
// Don't perform online validation or lookups if offline is true.
func SignStdTx(
	txBldr authtypes.TxBuilder, cliCtx context.CLIContext, name string,
	stdTx *authtypes.StdTx, appendSig bool, offline bool,
) (*authtypes.StdTx, error) {

	info, err := txBldr.Keybase().Get(name)
	if err != nil {
		return nil, err
	}

	addr := info.GetPubKey().Address()

	// check whether the address is a signer
	if !isTxSigner(sdk.AccAddress(addr), stdTx.GetSigners()) {
		return nil, fmt.Errorf("%s: %s", errInvalidSigner, name)
	}

	if !offline {
		txBldr, err = populateAccountFromState(txBldr, cliCtx, sdk.AccAddress(addr))
		if err != nil {
			return nil, err
		}
	}

	return txBldr.SignStdTx(name, keys.DefaultKeyPass, stdTx, appendSig)
}

// SignStdTxWithSignerAddress attaches a signature to a StdTx and returns a copy of a it.
// Don't perform online validation or lookups if offline is true, else
// populate account and sequence numbers from a foreign account.
func SignStdTxWithSignerAddress(txBldr authtypes.TxBuilder, cliCtx context.CLIContext,
	addr sdk.AccAddress, name string, stdTx *authtypes.StdTx,
	offline bool) (signedStdTx *authtypes.StdTx, err error) {

	// check whether the address is a signer
	if !isTxSigner(addr, stdTx.GetSigners()) {
		return signedStdTx, fmt.Errorf("%s: %s", errInvalidSigner, name)
	}

	if !offline {
		txBldr, err = populateAccountFromState(txBldr, cliCtx, addr)
		if err != nil {
			return signedStdTx, err
		}
	}

	return txBldr.SignStdTx(name, keys.DefaultKeyPass, stdTx, false)
}

// Read and decode a StdTx from the given filename.  Can pass "-" to read from stdin.
func ReadStdTxFromFile(cdc *codec.Codec, filename string) (*authtypes.StdTx, error) {
	var bytes []byte
	var tx authtypes.StdTx
	var err error

	if filename == "-" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		return nil, err
	}

	if err = cdc.UnmarshalJSON(bytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

func populateAccountFromState(
	txBldr authtypes.TxBuilder, cliCtx context.CLIContext, addr sdk.AccAddress,
) (authtypes.TxBuilder, error) {

	num, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(addr)
	if err != nil {
		return txBldr, err
	}

	return txBldr.WithAccountNumber(num).WithSequence(seq), nil
}

type txEncoderConfig struct {
	isEthereumTx bool
}

type Option func(config *txEncoderConfig)

func WithEthereumTx() Option {
	return func(cfg *txEncoderConfig) {
		cfg.isEthereumTx = true
	}
}

// GetTxEncoder return tx encoder from global sdk configuration if ones is defined.
// Otherwise returns encoder with default logic.
func GetTxEncoder(cdc *codec.Codec, options ...Option) (encoder sdk.TxEncoder) {
	encoder = sdk.GetConfig().GetTxEncoder()
	if encoder == nil {
		var cfg txEncoderConfig
		for _, op := range options {
			op(&cfg)
		}
		if cfg.isEthereumTx {
			encoder = authtypes.EthereumTxEncoder(cdc)
		} else {
			encoder = authtypes.DefaultTxEncoder(cdc)
		}
	}

	return
}

// simulateMsgs simulates the transaction and returns the simulation response and
// the adjusted gas value.
func simulateMsgs(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (sdk.SimulationResponse, uint64, error) {
	txBytes, err := txBldr.BuildTxForSim(msgs)
	if err != nil {
		return sdk.SimulationResponse{}, 0, err
	}

	return CalculateGas(cliCtx.QueryWithData, cliCtx.Codec, txBytes, txBldr.GasAdjustment())
}

func adjustGasEstimate(estimate uint64, adjustment float64) uint64 {
	return uint64(adjustment * float64(estimate))
}

func parseQueryResponse(cdc *codec.Codec, rawRes []byte) (sdk.SimulationResponse, error) {
	var simRes sdk.SimulationResponse
	if err := cdc.UnmarshalBinaryBare(rawRes, &simRes); err != nil {
		return sdk.SimulationResponse{}, err
	}

	return simRes, nil
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilder(txBldr authtypes.TxBuilder, cliCtx context.CLIContext) (authtypes.TxBuilder, error) {
	from := cliCtx.GetFromAddress()

	accGetter := authtypes.NewAccountRetriever(cliCtx)
	if err := accGetter.EnsureExists(from); err != nil {
		return txBldr, err
	}

	txbldrAccNum, txbldrAccSeq := txBldr.AccountNumber(), txBldr.Sequence()
	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txbldrAccNum == 0 || txbldrAccSeq == 0 {
		num, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
		if err != nil {
			return txBldr, err
		}

		if txbldrAccNum == 0 {
			txBldr = txBldr.WithAccountNumber(num)
		}
		if txbldrAccSeq == 0 {
			txBldr = txBldr.WithSequence(seq)
		}
	}

	return txBldr, nil
}

func buildUnsignedStdTxOffline(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx *authtypes.StdTx, err error) {
	if txBldr.SimulateAndExecute() {
		if cliCtx.GenerateOnly {
			return stdTx, errors.New("cannot estimate gas with generate-only")
		}

		txBldr, err = EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return stdTx, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txBldr.Gas())
	}

	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return stdTx, err
	}

	return authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
}

func isTxSigner(user sdk.AccAddress, signers []sdk.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}

	return false
}

func CliConvertCoinToCoinAdapters(coins sdk.Coins) sdk.CoinAdapters {
	ret := make(sdk.CoinAdapters, 0)
	for _, v := range coins {
		ret = append(ret, CliConvertCoinToCoinAdapter(v))
	}
	return ret
}

func CliConvertCoinToCoinAdapter(coin sdk.Coin) sdk.CoinAdapter {
	prec := newCoinFromDec()

	am := sdk.NewIntFromBigInt(coin.Amount.BigInt().Div(coin.Amount.BigInt(), prec))

	return sdk.NewCoinAdapter(coin.Denom, am)
}
