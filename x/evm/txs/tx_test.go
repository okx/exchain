package txs

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
	"reflect"
	"testing"
	"time"
)

var sdkResult sdk.Result

type EmptyAccount struct{}

func (e EmptyAccount) Copy() interface{}                            { return nil }
func (e EmptyAccount) GetAddress() sdk.AccAddress                   { return sdk.AccAddress{} }
func (e EmptyAccount) SetAddress(address sdk.AccAddress) error      { return nil }
func (e EmptyAccount) GetPubKey() crypto.PubKey                     { return nil }
func (e EmptyAccount) SetPubKey(key crypto.PubKey) error            { return nil }
func (e EmptyAccount) GetAccountNumber() uint64                     { return 0 }
func (e EmptyAccount) SetAccountNumber(u uint64) error              { return nil }
func (e EmptyAccount) GetSequence() uint64                          { return 0 }
func (e EmptyAccount) SetSequence(u uint64) error                   { return nil }
func (e EmptyAccount) GetCoins() sdk.Coins                          { return sdk.Coins{} }
func (e EmptyAccount) SetCoins(coins sdk.Coins) error               { return nil }
func (e EmptyAccount) SpendableCoins(blockTime time.Time) sdk.Coins { return sdk.Coins{} }
func (e EmptyAccount) String() string                               { return "ut" }

type EmptyTx struct {
	PrepareFail        bool
	GetChainConfigFail bool
	TransitionFail     bool
	DecorateResultFail bool
}

func (e EmptyTx) Prepare(msg *types.MsgEthereumTx) (err error) {
	if e.PrepareFail {
		return fmt.Errorf("prepare error")
	}
	return nil
}
func (e EmptyTx) SaveTx(msg *types.MsgEthereumTx) {}
func (e EmptyTx) GetChainConfig() (*types.ChainConfig, bool) {
	if e.GetChainConfigFail {
		return &types.ChainConfig{}, false
	}
	return &types.ChainConfig{}, true
}
func (e EmptyTx) GetSenderAccount() authexported.Account                                          { return EmptyAccount{} }
func (e EmptyTx) ResetWatcher(account authexported.Account)                                       {}
func (e EmptyTx) RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int) {}
func (e EmptyTx) Transition(config *types.ChainConfig) (result base.Result, err error) {
	if e.TransitionFail {
		err = fmt.Errorf("transition error")
		return
	}
	var execResult types.ExecutionResult
	execResult.Result = &sdkResult
	result.ExecResult = &execResult
	return
}

func (e EmptyTx) DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error) {
	if e.DecorateResultFail {
		return nil, fmt.Errorf("decorate result error")
	}
	if inErr != nil {
		return nil, inErr
	}

	return &sdkResult, nil
}

func (e EmptyTx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {}
func (e EmptyTx) Commit(msg *types.MsgEthereumTx, result *base.Result)      {}
func (e EmptyTx) EmitEvent(msg *types.MsgEthereumTx, result *base.Result)   {}
func (e EmptyTx) FinalizeWatcher(account authexported.Account, err error)   {}
func (e EmptyTx) AnalyzeStart(tag string)                                   {}
func (e EmptyTx) AnalyzeStop(tag string)                                    {}

func TestTransitionEvmTx(t *testing.T) {
	privateKey, _ := ethsecp256k1.GenerateKey()
	sender := ethcmn.HexToAddress(privateKey.PubKey().Address().String())
	msg := types.NewMsgEthereumTx(0, &sender, big.NewInt(100), 3000000, big.NewInt(1), nil)
	type args struct {
		tx  Tx
		msg *types.MsgEthereumTx
	}
	tests := []struct {
		name       string
		args       args
		wantResult *sdk.Result
		wantErr    bool
	}{
		{"1. none error", args{tx: EmptyTx{}, msg: &msg}, &sdkResult, false},
		{"2. prepare error", args{tx: EmptyTx{PrepareFail: true}, msg: &msg}, nil, true},
		{"3. transition error", args{tx: EmptyTx{TransitionFail: true}, msg: &msg}, nil, true},
		{"4. decorate result error", args{tx: EmptyTx{DecorateResultFail: true}, msg: &msg}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := TransitionEvmTx(tt.args.tx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionEvmTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("TransitionEvmTx() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
