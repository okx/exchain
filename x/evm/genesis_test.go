package evm_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/app"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
)

func (suite *EvmTestSuite) TestExportImport() {
	var genState types.GenesisState
	suite.Require().NotPanics(func() {
		genState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
	})

	_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, genState)
}

func (suite *EvmTestSuite) TestInitGenesis() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	testCases := []struct {
		name     string
		malleate func()
		genState types.GenesisState
		expPanic bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesisState(),
			false,
		},
		{
			"valid account",
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
				suite.Require().NotNil(acc)
				err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
				suite.Require().NoError(err)
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
						Storage: types.Storage{
							{Key: common.BytesToHash([]byte("key")), Value: common.BytesToHash([]byte("value"))},
						},
					},
				},
			},
			false,
		},
		{
			"account not found",
			func() {},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
		{
			"invalid account type",
			func() {
				acc := authtypes.NewBaseAccountWithAddress(address.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, &acc)
			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values

			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(
					func() {
						_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, tc.genState)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, tc.genState)
					},
				)
			}
		})
	}
}

func (suite *EvmTestSuite) TestInit() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	testCases := []struct {
		name     string
		malleate func(genesisState *simapp.GenesisState)
		genState types.GenesisState
		expPanic bool
	}{
		{
			"valid account",
			func(genesisState *simapp.GenesisState) {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
				suite.Require().NotNil(acc)
				err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
				suite.Require().NoError(err)
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
				authGenesisState := auth.ExportGenesis(suite.ctx, suite.app.AccountKeeper)
				(*genesisState)["auth"] = authtypes.ModuleCdc.MustMarshalJSON(authGenesisState)

			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
				TxsLogs:     []types.TransactionLogs{},
				ChainConfig: types.DefaultChainConfig(),
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values

			db := dbm.NewMemDB()
			chain := app.NewOKExChainApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
			genesisState := app.NewDefaultGenesisState()

			tc.malleate(&genesisState)

			genesisState["evm"] = types.ModuleCdc.MustMarshalJSON(tc.genState)
 			stateBytes, err := codec.MarshalJSONIndent(chain.Codec(), genesisState)
			if err != nil {
				panic(err)
			}

			if tc.expPanic {
				suite.Require().Panics(
					func() {
						chain.InitChain(
							abci.RequestInitChain{
								Validators:    []abci.ValidatorUpdate{},
								AppStateBytes: stateBytes,
							},
						)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						chain.InitChain(
							abci.RequestInitChain{
								Validators:    []abci.ValidatorUpdate{},
								AppStateBytes: stateBytes,
							},
						)
					},
				)
			}
		})
	}
}

func (suite *EvmTestSuite) TestExport() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
	suite.Require().NotNil(acc)
	err = acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
	suite.Require().NoError(err)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	initGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		Accounts: []types.GenesisAccount{
			{
				Address: address.String(),
				Storage: types.Storage{
					{Key: common.BytesToHash([]byte("key")), Value: common.BytesToHash([]byte("value"))},
				},
			},
		},
		TxsLogs: []types.TransactionLogs{
			{
				Hash: common.BytesToHash([]byte("tx_hash")),
				Logs: []*ethtypes.Log{
					{
						Address:     address,
						Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
						Data:        []byte("data"),
						BlockNumber: 1,
						TxHash:      ethcmn.BytesToHash([]byte("tx_hash")),
						TxIndex:     1,
						BlockHash:   ethcmn.BytesToHash([]byte("block_hash")),
						Index:       1,
						Removed:     false,
					},
				},
			},
		},
	}
	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)

	suite.Require().NotPanics(func() {
		evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
	})
}
