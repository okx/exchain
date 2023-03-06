package benchmarks

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/okx/okbchain/libs/system"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/okbchain/app"
	types3 "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/simapp"
	"github.com/okx/okbchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	banktypes "github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto"
	"github.com/okx/okbchain/libs/tendermint/crypto/secp256k1"
	"github.com/okx/okbchain/libs/tendermint/global"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/libs/tendermint/types"
	dbm "github.com/okx/okbchain/libs/tm-db"
	types2 "github.com/okx/okbchain/x/evm/types"
	wasmtypes "github.com/okx/okbchain/x/wasm/types"
	"github.com/stretchr/testify/require"
)

func TestTxSending(t *testing.T) {
	db := dbm.NewMemDB()
	defer db.Close()
	appInfo := InitializeOKXApp(t, db, 50)
	height := int64(2)
	global.SetGlobalHeight(height - 1)
	appInfo.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})
	txs := GenSequenceOfTxs(&appInfo, bankSendMsg, 100)
	for _, tx := range txs {
		res := appInfo.App.DeliverTx(abci.RequestDeliverTx{Tx: tx})
		require.True(t, res.IsOK())
	}

	appInfo.App.EndBlock(abci.RequestEndBlock{Height: height})
	appInfo.App.Commit(abci.RequestCommit{})
}

func TestOip20TxSending(t *testing.T) {
	db := dbm.NewMemDB()
	defer db.Close()
	appInfo := InitializeOKXApp(t, db, 50)
	err := deployOip20(&appInfo)
	require.NoError(t, err)
	global.SetGlobalHeight(appInfo.height)
	height := appInfo.height + 1
	appInfo.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})
	txs := buildOip20Transfer(100, &appInfo)
	for _, tx := range txs {
		res := appInfo.App.DeliverTx(abci.RequestDeliverTx{Tx: tx})
		require.True(t, res.IsOK())
	}

	appInfo.App.EndBlock(abci.RequestEndBlock{Height: height})
	appInfo.App.Commit(abci.RequestCommit{})
}

func TestCw20TxSending(t *testing.T) {
	db := dbm.NewMemDB()
	defer db.Close()
	appInfo := InitializeOKXApp(t, db, 50)

	emptyBlock(&appInfo)
	err := deployCw20(&appInfo)
	require.NoError(t, err)

	height := appInfo.height + 1
	appInfo.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})
	txs := buildTxFromMsg(cw20TransferMsg)(100, &appInfo)
	for _, tx := range txs {
		res := appInfo.App.DeliverTx(abci.RequestDeliverTx{Tx: tx})
		require.True(t, res.IsOK())
	}

	appInfo.App.EndBlock(abci.RequestEndBlock{Height: height})
	appInfo.App.Commit(abci.RequestCommit{})
}

type AppInfo struct {
	height int64

	App              *app.OKExChainApp
	evmMintKey       *ecdsa.PrivateKey
	evmMintAddr      sdk.AccAddress
	MinterKey        crypto.PrivKey
	MinterAddr       sdk.AccAddress
	ContractAddr     ethcmn.Address
	Cw20CodeID       uint64
	Cw20ContractAddr string
	Cw1CodeID        uint64
	Cw1ContractAddr  string
	Denom            string
	AccNum           uint64
	SeqNum           uint64
	Nonce            uint64
}

func InitializeOKXApp(b testing.TB, db dbm.DB, numAccounts int) AppInfo {
	types.UnittestOnlySetMilestoneEarthHeight(1)
	evmMinter, _ := ethcrypto.HexToECDSA(PrivateKey)
	evmMinterAddr := sdk.AccAddress(ethcrypto.PubkeyToAddress(evmMinter.PublicKey).Bytes())

	// constants
	minter := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(minter.PubKey().Address())
	denom := system.Currency

	// genesis setup (with a bunch of random accounts)
	genAccs := make([]authexported.GenesisAccount, numAccounts+2)
	genAccs[0] = &types3.EthAccount{
		BaseAccount: &authtypes.BaseAccount{
			Address: evmMinterAddr,
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(denom, 1<<60)),
		},
	}
	genAccs[1] = &authtypes.BaseAccount{
		Address: addr,
		Coins:   sdk.NewCoins(sdk.NewInt64Coin(denom, 1<<60)),
		PubKey:  minter.PubKey(),
	}

	for i := 2; i <= numAccounts+1; i++ {
		priv := secp256k1.GenPrivKey()
		genAccs[i] = &authtypes.BaseAccount{
			Address: sdk.AccAddress(priv.PubKey().Address()),
			Coins:   sdk.NewCoins(sdk.NewInt64Coin(denom, 100000000000)),
			PubKey:  priv.PubKey(),
		}
	}
	okxApp := SetupWithGenesisAccounts(b, db, genAccs)

	info := AppInfo{
		height:      1,
		App:         okxApp,
		evmMintKey:  evmMinter,
		evmMintAddr: evmMinterAddr,
		MinterKey:   minter,
		MinterAddr:  addr,
		Denom:       denom,
		AccNum:      1,
		SeqNum:      0,
		Nonce:       0,
	}

	return info
}

func setup(db dbm.DB, withGenesis bool, invCheckPeriod uint) (*app.OKExChainApp, simapp.GenesisState) {
	okxApp := app.NewOKExChainApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, invCheckPeriod)
	if withGenesis {
		return okxApp, app.NewDefaultGenesisState()
	}
	return okxApp, simapp.GenesisState{}
}

// SetupWithGenesisAccounts initializes a new OKExChainApp with the provided genesis
// accounts and possible balances.
func SetupWithGenesisAccounts(b testing.TB, db dbm.DB, genAccs []authexported.GenesisAccount) *app.OKExChainApp {
	okxApp, genesisState := setup(db, true, 0)
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	appCodec := okxApp.Codec()

	genesisState[authtypes.ModuleName] = appCodec.MustMarshalJSON(authGenesis)

	bankGenesis := banktypes.NewGenesisState(true)
	genesisState[banktypes.ModuleName] = appCodec.MustMarshalJSON(bankGenesis)
	evmGenesis := types2.DefaultGenesisState()
	evmGenesis.Params.EnableCreate = true
	evmGenesis.Params.EnableCall = true
	evmGenesis.Params.MaxGasLimitPerTx = GasLimit * 2
	genesisState[types2.ModuleName] = appCodec.MustMarshalJSON(evmGenesis)

	genesisState[wasmtypes.ModuleName] = appCodec.MustMarshalJSON(
		wasmtypes.GenesisState{
			Params: wasmtypes.DefaultParams(),
		})

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	okxApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: types.TM2PB.ConsensusParams(types.DefaultConsensusParams()),
			AppStateBytes:   stateBytes,
		},
	)

	okxApp.Commit(abci.RequestCommit{})
	okxApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: okxApp.LastBlockHeight() + 1}})

	return okxApp
}

func deployOip20(info *AppInfo) error {
	// add oip20 contract
	global.SetGlobalHeight(info.height)
	height := info.height + 1
	info.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})

	// deploy oip20
	OipBytes, err := hex.DecodeString(Oip20Bin)
	if err != nil {
		return err
	}
	oip20DeployTx := types2.NewMsgEthereumTxContract(0, nil, GasLimit, big.NewInt(GasPrice), OipBytes)
	if err = oip20DeployTx.Sign(big.NewInt(ChainId), info.evmMintKey); err != nil {
		return err
	}
	signedOipBytes, err := rlp.EncodeToBytes(oip20DeployTx)
	if err != nil {
		return err
	}
	res := info.App.DeliverTx(abci.RequestDeliverTx{Tx: signedOipBytes})
	info.Nonce++

	// TODO: parse contract address better
	i := strings.Index(res.Log, "contract address")
	info.ContractAddr = ethcmn.HexToAddress(res.Log[i+17 : i+17+42])

	info.App.EndBlock(abci.RequestEndBlock{Height: height})
	info.App.Commit(abci.RequestCommit{})

	info.height++
	return nil
}

func deployCw20(info *AppInfo) error {
	// add cw20 contract
	global.SetGlobalHeight(info.height)
	height := info.height + 1
	info.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})

	// upload cw20
	txs := buildTxFromMsg(cw20StoreMsg)(1, info)
	res := info.App.DeliverTx(abci.RequestDeliverTx{Tx: txs[0]})
	if !res.IsOK() {
		return errors.New("deliver tx error")
	}
	codeID, err := strconv.Atoi(string(res.Events[2].Attributes[0].Value))
	if err != nil {
		return err
	}
	info.Cw20CodeID = uint64(codeID)

	// instantiate cw20
	txs = buildTxFromMsg(cw20InstantiateMsg)(1, info)
	res = info.App.DeliverTx(abci.RequestDeliverTx{Tx: txs[0]})
	if !res.IsOK() {
		return errors.New("deliver tx error")
	}
	info.Cw20ContractAddr = string(res.Events[2].Attributes[0].Value)

	info.App.EndBlock(abci.RequestEndBlock{Height: height})
	info.App.Commit(abci.RequestCommit{})

	info.height++
	return nil
}

func emptyBlock(info *AppInfo) {
	global.SetGlobalHeight(info.height)
	height := info.height + 1
	info.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})
	info.App.EndBlock(abci.RequestEndBlock{Height: height})
	info.App.Commit(abci.RequestCommit{})

	info.height++
}

func GenSequenceOfTxs(info *AppInfo, msgGen func(*AppInfo) ([]sdk.Msg, error), numToGenerate int) []types.Tx {
	fees := sdk.Coins{sdk.NewInt64Coin(info.Denom, 1)}
	txs := make([]types.Tx, numToGenerate)

	for i := 0; i < numToGenerate; i++ {
		msgs, err := msgGen(info)
		if err != nil {
			panic(err)
		}
		tx := helpers.GenTx(
			msgs,
			fees,
			1e8,
			"exchain-67",
			[]uint64{info.AccNum},
			[]uint64{info.SeqNum},
			info.MinterKey,
		)
		txs[i] = info.App.Codec().MustMarshalBinaryLengthPrefixed(tx)

		info.SeqNum += 1
	}

	return txs
}
