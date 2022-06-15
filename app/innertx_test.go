package app

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/distribution/keeper"
	"github.com/okex/exchain/x/evm"
	evm_types "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/gov/types"
	"github.com/okex/exchain/x/staking"
	staking_keeper "github.com/okex/exchain/x/staking/keeper"
	staking_types "github.com/okex/exchain/x/staking/types"

	"github.com/stretchr/testify/suite"
)

var (
	coin10  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	coin20  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 20)
	coin30  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 30)
	coin40  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 40)
	coin50  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)
	coin60  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 60)
	coin70  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 70)
	coin80  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 80)
	coin90  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 90)
	coin100 = sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)
	fees    = auth.NewStdFee(21000, sdk.NewCoins(coin10))
)

type InnerTxTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *OKExChainApp
	stateDB *evm_types.CommitStateDB
	codec   *codec.Codec

	handler sdk.Handler
}

func (suite *InnerTxTestSuite) SetupTest() {
	checkTx := false
	chain_id := "ethermint-3"
	key, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	innerTxPath := TestDataPath + "/innerTx-" + key.PubKey().Address().String()
	viper.Set(cli.HomeFlag, innerTxPath)
	viper.Set("db_backend", "goleveldb")

	suite.app = Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: chain_id, Time: time.Now().UTC()})
	suite.ctx.SetDeliver()
	suite.stateDB = evm_types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)
	suite.codec = suite.app.Codec()

	err = ethermint.SetChainId(chain_id)
	suite.Nil(err)

	params := evm_types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
}

func TestInnerTxTestSuite(t *testing.T) {
	suite.Run(t, new(InnerTxTestSuite))
}

func (suite *InnerTxTestSuite) TestMsgSend() {
	defer func() {
		err := os.RemoveAll(TestDataPath)
		suite.Require().NoError(err)
	}()
	var (
		tx          sdk.Tx
		privFrom, _ = ethsecp256k1.GenerateKey()
		ethFrom     = common.HexToAddress(privFrom.PubKey().Address().String())
		cmFrom      = sdk.AccAddress(privFrom.PubKey().Address())
		privTo      = secp256k1.GenPrivKeySecp256k1([]byte("private key to"))
		ethTo       = common.HexToAddress(privTo.PubKey().Address().String())
		cmTo        = sdk.AccAddress(privTo.PubKey().Address())

		valPriv      = ed25519.GenPrivKeyFromSecret([]byte("ed25519 private key"))
		valpub       = valPriv.PubKey()
		valopaddress = sdk.ValAddress(valpub.Address())
		valcmaddress = sdk.AccAddress(valpub.Address())

		privFrom1 = secp256k1.GenPrivKeySecp256k1([]byte("from1"))
		cmFrom1   = sdk.AccAddress(privFrom1.PubKey().Address())
		privTo1   = secp256k1.GenPrivKeySecp256k1([]byte("to1"))
		cmTo1     = sdk.AccAddress(privTo1.PubKey().Address())

		exceptedInnerTx       = make([]vm.InnerTx, 0)
		normalExceptedInnerTx = make([]vm.InnerTx, 0)
	)
	hash := tmhash.Sum([]byte("hash-1"))
	hexHash := hexutil.Encode(hash)
	votes := []abci.VoteInfo{
		{Validator: abci.Validator{Address: valpub.Address(), Power: 1}, SignedLastBlock: true},
	}
	bHeader := abci.Header{Height: 0, ChainID: "ethermint-3", Time: time.Now().UTC(), ProposerAddress: sdk.ConsAddress(valpub.Address())}
	suite.ctx.SetBlockHeader(bHeader)
	req := abci.RequestBeginBlock{Header: bHeader, Hash: hash}

	normal := func() {
		err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(coin100))
		suite.Require().NoError(err)
	}

	normalExcepted1 := innertx.CreateInnerTx(0, "ex1m3h30wlvsf8llruxtpukdvsy0km2kum8qc3awh", "ex17xpfvakm2amg962yls6f84z3kell8c5lcs49z2", innertx.CosmosCallType, innertx.SendCallName, big.NewInt(500000000000000000), nil)
	normalExcepted2 := innertx.CreateInnerTx(0, "ex1m3h30wlvsf8llruxtpukdvsy0km2kum8qc3awh", "ex1u0dcdmkqk5pc22pf2fku3n0up8y7y3xfkyzh7w", innertx.CosmosCallType, innertx.SendCallName, big.NewInt(500000000000000000), nil)
	normalExceptedInnerTx = append(normalExceptedInnerTx, *normalExcepted1, *normalExcepted2)
	testCases := []struct {
		msg        string
		prepare    func()
		expPass    bool
		expectfunc func()
	}{
		{
			"send msg(bank)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = bank.NewHandler(suite.app.BankKeeper)

				msg := bank.NewMsgSend(cmFrom, cmTo, sdk.NewCoins(coin10))
				tx = auth.NewStdTx([]sdk.Msg{msg}, fees, nil, "")
				excepted := innertx.CreateInnerTx(0, cmFrom.String(), cmTo.String(), innertx.CosmosCallType, innertx.SendCallName, coin10.Amount.Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin90))))

				toBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo).GetCoins()
				suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin10))))
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"send msgs(bank)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = bank.NewHandler(suite.app.BankKeeper)

				msg := bank.NewMsgSend(cmFrom, cmTo, sdk.NewCoins(coin10))
				tx = auth.NewStdTx([]sdk.Msg{msg, msg}, fees, nil, "")

				excepted := innertx.CreateInnerTx(0, cmFrom.String(), cmTo.String(), innertx.CosmosCallType, innertx.SendCallName, coin10.Amount.Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted, *excepted)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin80))))

				toBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo).GetCoins()
				suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin20))))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"multi msg(bank)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = bank.NewHandler(suite.app.BankKeeper)
				suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(coin100))
				suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom1, sdk.NewCoins(coin100))
				inputCoin1 := sdk.NewCoins(coin20)
				inputCoin2 := sdk.NewCoins(coin10)
				outputCoin1 := sdk.NewCoins(coin10)
				outputCoin2 := sdk.NewCoins(coin20)
				input1 := bank.NewInput(cmFrom, inputCoin1)
				input2 := bank.NewInput(cmFrom1, inputCoin2)
				output1 := bank.NewOutput(cmTo, outputCoin1)
				output2 := bank.NewOutput(cmTo1, outputCoin2)

				msg := bank.NewMsgMultiSend([]bank.Input{input1, input2}, []bank.Output{output1, output2})
				tx = auth.NewStdTx([]sdk.Msg{msg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, cmFrom.String(), sdk.AccAddress{}.String(), innertx.CosmosCallType, innertx.MultiCallName, coin20.Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom1.String(), sdk.AccAddress{}.String(), innertx.CosmosCallType, innertx.MultiCallName, coin10.Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, sdk.AccAddress{}.String(), cmTo.String(), innertx.CosmosCallType, innertx.MultiCallName, coin10.Amount.Int, nil)
				excepted4 := innertx.CreateInnerTx(0, sdk.AccAddress{}.String(), cmTo1.String(), innertx.CosmosCallType, innertx.MultiCallName, coin20.Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3, *excepted4)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin80))))
				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom1).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin90))))

				toBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo).GetCoins()
				suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin10))))
				toBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo1).GetCoins()
				suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin20))))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"evm send msg(evm)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = evm.NewHandler(suite.app.EvmKeeper)
				tx = evm_types.NewMsgEthereumTx(0, &ethTo, coin10.Amount.BigInt(), 3000000, big.NewInt(0), nil)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(suite.ctx.ChainID())
				suite.Require().NoError(err)

				// sign transaction
				ethTx, ok := tx.(*evm_types.MsgEthereumTx)
				suite.Require().True(ok)

				err = ethTx.Sign(chainID, privFrom.ToECDSA())
				suite.Require().NoError(err)
				tx = ethTx

				excepted1 := innertx.CreateInnerTx(0, ethFrom.String(), ethTx.To().String(), innertx.CosmosCallType, innertx.EvmCallName, coin10.Amount.BigInt(), nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin90))))

				toBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo).GetCoins()
				suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin10))))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"create validator(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)
				tx = auth.NewStdTx([]sdk.Msg{msg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))

				suite.app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.ctx)
				val, ok := suite.app.StakingKeeper.GetValidator(suite.ctx, valopaddress)
				suite.Require().True(ok)
				suite.Require().Equal(valopaddress, val.OperatorAddress)
				suite.Require().True(val.MinSelfDelegation.Equal(sdk.NewDec(10000)))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"destroy validator(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				destroyValMsg := staking_types.NewMsgDestroyValidator([]byte(valopaddress))
				tx = auth.NewStdTx([]sdk.Msg{msg, destroyValMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", "ex1tygms3xhhs3yv487phx3dw4a95jn7t7lfjrmx5", innertx.CosmosCallType, innertx.SendCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, "ex1tygms3xhhs3yv487phx3dw4a95jn7t7lfjrmx5", valcmaddress.String(), innertx.CosmosCallType, innertx.UndelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))

				suite.app.EndBlocker(suite.ctx.WithBlockTime(time.Now().Add(staking_types.DefaultUnbondingTime)), abci.RequestEndBlock{Height: 2})
				_, ok := suite.app.StakingKeeper.GetValidator(suite.ctx, valopaddress)
				suite.Require().False(ok)
				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"deposit msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)
				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"withdraw msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)
				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))

				withdrawMsg := staking_types.NewMsgWithdraw(cmFrom, keeper.NewTestSysCoin(10000, 0))
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, withdrawMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", "ex1tygms3xhhs3yv487phx3dw4a95jn7t7lfjrmx5", innertx.CosmosCallType, innertx.SendCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted4 := innertx.CreateInnerTx(0, "ex1tygms3xhhs3yv487phx3dw4a95jn7t7lfjrmx5", cmFrom.String(), innertx.CosmosCallType, innertx.UndelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3, *excepted4)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))
				suite.app.EndBlocker(suite.ctx.WithBlockTime(time.Now().Add(staking_types.DefaultUnbondingTime)), abci.RequestEndBlock{Height: 2})
				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))

				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"addshare msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)
				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				addShareMsg := staking_types.NewMsgAddShares(cmFrom, []sdk.ValAddress{valopaddress})
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, addShareMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"proxy reg msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)
				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				regMsg := staking_types.NewMsgRegProxy(cmFrom, true)
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, regMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"proxy unreg msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)
				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)

				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				regMsg := staking_types.NewMsgRegProxy(cmFrom, true)
				unregMsg := staking_types.NewMsgRegProxy(cmFrom, false)
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, regMsg, unregMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"proxy bind msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom1, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)
				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				regMsg := staking_types.NewMsgRegProxy(cmFrom, true)
				depositMsg1 := staking_types.NewMsgDeposit(cmFrom1, keeper.NewTestSysCoin(10000, 0))
				bindMsg := staking_types.NewMsgBindProxy(cmFrom1, cmFrom)

				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, regMsg, depositMsg1, bindMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom1).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"proxy unbind msg(staking)",
			func() {
				suite.app.BeginBlocker(suite.ctx, req)
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom1, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)
				depositMsg := staking_types.NewMsgDeposit(cmFrom, keeper.NewTestSysCoin(10000, 0))
				regMsg := staking_types.NewMsgRegProxy(cmFrom, true)
				depositMsg1 := staking_types.NewMsgDeposit(cmFrom1, keeper.NewTestSysCoin(10000, 0))
				bindMsg := staking_types.NewMsgBindProxy(cmFrom1, cmFrom)
				ubindMsg := staking_types.NewMsgUnbindProxy(cmFrom1)
				tx = auth.NewStdTx([]sdk.Msg{msg, depositMsg, regMsg, depositMsg1, bindMsg, ubindMsg}, fees, nil, "")

				excepted1 := innertx.CreateInnerTx(0, valcmaddress.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000).Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, cmFrom.String(), "ex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3ajl2sq", innertx.CosmosCallType, innertx.DelegateCallName, keeper.NewTestSysCoin(10000, 0).Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom1).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"withdraw validator(distr)",
			func() {
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)
				_, err = suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))
				suite.app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.ctx)
				val, ok := suite.app.StakingKeeper.GetValidator(suite.ctx, valopaddress)
				suite.Require().True(ok)
				suite.Require().Equal(valopaddress, val.OperatorAddress)
				suite.Require().True(val.MinSelfDelegation.Equal(sdk.NewDec(10000)))

				suite.app.Commit(abci.RequestCommit{})
				for i := 0; i < 100; i++ {
					header := abci.Header{Height: int64(i + 2), ProposerAddress: sdk.ConsAddress(valpub.Address())}
					req := abci.RequestBeginBlock{Header: header,
						LastCommitInfo: abci.LastCommitInfo{Votes: votes}}
					suite.ctx.SetBlockHeader(header)
					suite.app.BeginBlocker(suite.ctx, req)
					suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				}
				commision := suite.app.DistrKeeper.GetValidatorAccumulatedCommission(suite.ctx, valopaddress)
				suite.Require().True(commision.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(49, 0)))))

				suite.handler = distr.NewHandler(suite.app.DistrKeeper)
				withdrawMsg := distr.NewMsgWithdrawValidatorCommission(valopaddress)
				tx = auth.NewStdTx([]sdk.Msg{withdrawMsg}, fees, nil, "")

				suite.app.BeginBlocker(suite.ctx, req)
				excepted1 := innertx.CreateInnerTx(0, "ex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd80kjeqg", valcmaddress.String(), innertx.CosmosCallType, innertx.SendCallName, sdk.NewDecWithPrec(49, 0).Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				expectCommision := sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(49, 0)))
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))).Add2(expectCommision)))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				normalExceptedInnerTxTemp := make([]vm.InnerTx, 0)
				temp := innertx.CreateInnerTx(0, "ex17xpfvakm2amg962yls6f84z3kell8c5lcs49z2", "ex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd80kjeqg", innertx.CosmosCallType, innertx.SendCallName, big.NewInt(500000000000000000), nil)
				normalExceptedInnerTxTemp = append(normalExceptedInnerTxTemp, normalExceptedInnerTx...)
				normalExceptedInnerTxTemp = append(normalExceptedInnerTxTemp, *temp)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTxTemp))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"set withdraw address(distr)",
			func() {
				suite.handler = staking.NewHandler(suite.app.StakingKeeper)

				err := suite.app.BankKeeper.SetCoins(suite.ctx, valcmaddress, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000)))
				suite.Require().NoError(err)

				msg := staking_keeper.NewTestMsgCreateValidator(valopaddress, valpub, coin10.Amount)
				_, err = suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))
				suite.app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.ctx)
				val, ok := suite.app.StakingKeeper.GetValidator(suite.ctx, valopaddress)
				suite.Require().True(ok)
				suite.Require().Equal(valopaddress, val.OperatorAddress)
				suite.Require().True(val.MinSelfDelegation.Equal(sdk.NewDec(10000)))

				suite.app.Commit(abci.RequestCommit{})
				votes := []abci.VoteInfo{
					{Validator: abci.Validator{Address: valpub.Address(), Power: 1}, SignedLastBlock: true},
				}
				for i := 0; i < 100; i++ {
					header := abci.Header{Height: int64(i + 2), ProposerAddress: sdk.ConsAddress(valpub.Address())}
					req := abci.RequestBeginBlock{Header: header,
						LastCommitInfo: abci.LastCommitInfo{Votes: votes}}
					suite.ctx.SetBlockHeader(header)
					suite.app.BeginBlocker(suite.ctx, req)
					suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				}
				commision := suite.app.DistrKeeper.GetValidatorAccumulatedCommission(suite.ctx, valopaddress)
				suite.Require().True(commision.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(49, 0)))))

				suite.handler = distr.NewHandler(suite.app.DistrKeeper)
				setwithdrawMsg := distr.NewMsgSetWithdrawAddress(valcmaddress, cmFrom1)
				withdrawMsg := distr.NewMsgWithdrawValidatorCommission(valopaddress)
				tx = auth.NewStdTx([]sdk.Msg{setwithdrawMsg, withdrawMsg}, fees, nil, "")
				suite.app.BeginBlocker(suite.ctx, req)
				excepted1 := innertx.CreateInnerTx(0, "ex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd80kjeqg", cmFrom1.String(), innertx.CosmosCallType, innertx.SendCallName, sdk.NewDecWithPrec(49, 0).Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, valcmaddress).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)))))
				fromBalance = suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom1).GetCoins()
				expectCommision := sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(49, 0)))
				suite.Require().True(fromBalance.IsEqual(expectCommision))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				normalExceptedInnerTxTemp := make([]vm.InnerTx, 0)
				temp := innertx.CreateInnerTx(0, "ex17xpfvakm2amg962yls6f84z3kell8c5lcs49z2", "ex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd80kjeqg", innertx.CosmosCallType, innertx.SendCallName, big.NewInt(500000000000000000), nil)
				normalExceptedInnerTxTemp = append(normalExceptedInnerTxTemp, normalExceptedInnerTx...)
				normalExceptedInnerTxTemp = append(normalExceptedInnerTxTemp, *temp)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTxTemp))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"submit proposal(gov)",
			func() {
				suite.handler = gov.NewHandler(suite.app.GovKeeper)

				content := gov.NewTextProposal("Test", "description")
				newProposalMsg := gov.NewMsgSubmitProposal(content, keeper.NewTestSysCoins(100, 0), cmFrom)
				tx = auth.NewStdTx([]sdk.Msg{newProposalMsg}, fees, nil, "")

				suite.app.BeginBlocker(suite.ctx, req)
				excepted := innertx.CreateInnerTx(0, cmFrom.String(), "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", innertx.CosmosCallType, innertx.SendCallName, coin100.Amount.Int, nil)
				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsZero())

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"deposit proposal(gov)",
			func() {
				suite.handler = gov.NewHandler(suite.app.GovKeeper)

				content := gov.NewTextProposal("Test", "description")
				newProposalMsg := gov.NewMsgSubmitProposal(content, keeper.NewTestSysCoins(10, 0), cmFrom)
				depositMsg := gov.NewMsgDeposit(cmFrom, 1, keeper.NewTestSysCoins(90, 0))
				tx = auth.NewStdTx([]sdk.Msg{newProposalMsg, depositMsg}, fees, nil, "")

				suite.app.BeginBlocker(suite.ctx, req)
				excepted1 := innertx.CreateInnerTx(0, cmFrom.String(), "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", innertx.CosmosCallType, innertx.SendCallName, coin10.Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", innertx.CosmosCallType, innertx.SendCallName, coin90.Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsZero())

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
		{
			"vote proposal(gov)",
			func() {
				suite.handler = gov.NewHandler(suite.app.GovKeeper)

				validator := staking_types.NewValidator(valopaddress, valpub, staking_types.NewDescription("test description", "", "", ""), staking_types.DefaultMinDelegation)
				suite.app.StakingKeeper.SetValidator(suite.ctx, validator)
				content := gov.NewTextProposal("Test", "description")
				newProposalMsg := gov.NewMsgSubmitProposal(content, keeper.NewTestSysCoins(10, 0), cmFrom)
				depositMsg := gov.NewMsgDeposit(cmFrom, 1, keeper.NewTestSysCoins(90, 0))
				voteMsg := gov.NewMsgVote(valcmaddress, 1, types.OptionYes)
				tx = auth.NewStdTx([]sdk.Msg{newProposalMsg, depositMsg, voteMsg}, fees, nil, "")

				suite.app.BeginBlocker(suite.ctx, req)
				excepted1 := innertx.CreateInnerTx(0, cmFrom.String(), "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", innertx.CosmosCallType, innertx.SendCallName, coin10.Amount.Int, nil)
				excepted2 := innertx.CreateInnerTx(0, cmFrom.String(), "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", innertx.CosmosCallType, innertx.SendCallName, coin90.Amount.Int, nil)
				excepted3 := innertx.CreateInnerTx(0, "ex10d07y265gmmuvt4z0w9aw880jnsr700jjt9qly", cmFrom.String(), innertx.CosmosCallType, innertx.SendCallName, coin100.Amount.Int, nil)

				exceptedInnerTx = make([]vm.InnerTx, 0)
				exceptedInnerTx = append(exceptedInnerTx, *excepted1, *excepted2, *excepted3)
				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				suite.ctx.SetTxBytes(txBytes)
			},
			true,
			func() {
				fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
				suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(coin100))))

				suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
				innerBlock := GetBlockInternalTransactions(hexHash)
				innerTxs, ok := innerBlock[hexHash]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, normalExceptedInnerTx))

				txBytes, err := utils.GetTxEncoder(suite.codec)(tx)
				suite.Require().NoError(err)
				txHash := tmtypes.Tx(txBytes).Hash(suite.ctx.BlockHeight())
				txStr := common.BytesToHash(txHash).Hex()
				innerTxs, ok = innerBlock[txStr]
				suite.Require().True(ok)
				suite.Require().True(InnerTxsEqual(innerTxs, exceptedInnerTx))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset
			normal()
			//nolint
			tc.prepare()
			suite.ctx.SetGasMeter(sdk.NewInfiniteGasMeter())
			msgs := tx.GetMsgs()
			for _, msg := range msgs {
				_, err := suite.handler(suite.ctx, msg)

				//nolint
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			}
			tc.expectfunc()

		})
	}
}

func GetBlockInternalTransactions(blockHash string) map[string][]vm.InnerTx {
	var rtn = make(map[string][]vm.InnerTx)
	txHashes := vm.GetBlockDB(blockHash)
	if txHashes != nil {
		for _, txHash := range txHashes {
			inners := vm.GetFromDB(txHash)
			rtn[txHash] = inners
		}
	} else {
		rtn = nil
	}
	return rtn
}

func InnerTxsEqual(actual []vm.InnerTx, excepted []vm.InnerTx) (isEqual bool) {
	defer func() {
		if !isEqual {
			fmt.Println("excepted", excepted)
			fmt.Println("actual", actual)
		}
	}()
	if len(actual) != len(excepted) {
		return false
	}
	actualMap := make(map[string]vm.InnerTx, 0)
	for i, _ := range actual {
		bytes, err := json.Marshal(actual[i])
		if err != nil {
			panic(err)
		}
		actualMap[hexutil.Encode(bytes)] = actual[i]
	}

	for i, _ := range excepted {
		bytes, err := json.Marshal(excepted[i])
		if err != nil {
			panic(err)
		}
		_, ok := actualMap[hexutil.Encode(bytes)]
		if !ok {
			return false
		}
	}
	return true
}
