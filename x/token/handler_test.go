package token_test

import (
	"testing"

	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/mock"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	okexchain "github.com/okex/exchain/app"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/token"
	"github.com/okex/exchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestHandlerBlockedContractAddrSend(t *testing.T) {
	okexapp := initApp(true)
	ctx := okexapp.BaseApp.NewContext(true, abci.Header{Height: 1})
	gAcc := CreateEthAccounts(3, sdk.SysCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(10000)),
	})
	okexapp.AccountKeeper.SetAccount(ctx, gAcc[0])
	okexapp.AccountKeeper.SetAccount(ctx, gAcc[1])
	gAcc[2].CodeHash = []byte("contract code hash")
	okexapp.AccountKeeper.SetAccount(ctx, gAcc[2])

	// multi send
	multiSendStr := `[{"to":"` + gAcc[1].Address.String() + `","amount":" 10` + common.NativeToken + `"}]`
	transfers, err := types.StrToTransfers(multiSendStr)
	require.Nil(t, err)
	multiSendStr2 := `[{"to":"` + gAcc[2].Address.String() + `","amount":" 10` + common.NativeToken + `"}]`
	transfers2, err := types.StrToTransfers(multiSendStr2)
	require.Nil(t, err)

	successfulSendMsg := types.NewMsgTokenSend(gAcc[0].Address, gAcc[1].Address, sdk.SysCoins{sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1))})
	sendToContractMsg := types.NewMsgTokenSend(gAcc[0].Address, gAcc[2].Address, sdk.SysCoins{sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1))})
	successfulMultiSendMsg := types.NewMsgMultiSend(gAcc[0].Address, transfers)
	multiSendToContractMsg := types.NewMsgMultiSend(gAcc[0].Address, transfers2)
	handler := token.NewTokenHandler(okexapp.TokenKeeper, version.CurrentProtocolVersion)
	okexapp.BankKeeper.SetSendEnabled(ctx, true)
	TestSets := []struct {
		description string
		balance     string
		msg         sdk.Msg
		account     app.EthAccount
	}{
		// 0.01okt as fixed fee in each stdTx
		{"success to send", "9999.000000000000000000okt", successfulSendMsg, gAcc[0]},
		{"success to multi-send", "9989.000000000000000000okt", successfulMultiSendMsg, gAcc[0]},
		{"success to send", "9988.000000000000000000okt", successfulSendMsg, gAcc[0]},
		{"success to multi-send", "9978.000000000000000000okt", successfulMultiSendMsg, gAcc[0]},
		//{"fail to send to contract", "9978.000000000000000000okt", failedSendMsg, gAcc[0]},
		//{"fail to multi-send to contract", "9978.000000000000000000okt", failedMultiSendMsg, gAcc[0]},
		{"fail to send to contract", "9978.000000000000000000okt", sendToContractMsg, gAcc[0]},
		{"fail to multi-send to contract", "9978.000000000000000000okt", multiSendToContractMsg, gAcc[0]},
	}
	for i, tt := range TestSets {
		t.Run(tt.description, func(t *testing.T) {
			handler(ctx, TestSets[i].msg)
			acc := okexapp.AccountKeeper.GetAccount(ctx, tt.account.Address)
			acc.GetCoins().String()
			require.Equal(t, acc.GetCoins().String(), tt.balance)
		})
	}
}

// Setup initializes a new OKExChainApp. A Nop logger is set in OKExChainApp.
func initApp(isCheckTx bool) *okexchain.OKExChainApp {
	db := dbm.NewMemDB()
	app := okexchain.NewOKExChainApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := okexchain.NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
		app.EndBlock(abci.RequestEndBlock{})
		app.Commit()
	}

	return app
}

func CreateEthAccounts(numAccs int, genCoins sdk.SysCoins) (genAccs []app.EthAccount) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		ak := mock.NewAddrKeys(addr, pubKey, privKey)
		testAccount := app.EthAccount{
			BaseAccount: &auth.BaseAccount{
				Address: ak.Address,
				Coins:   genCoins,
			},
			CodeHash: ethcrypto.Keccak256(nil),
		}
		genAccs = append(genAccs, testAccount)
	}
	return
}
