package feesplit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/okx/okbchain/app"
	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/feesplit"
	"github.com/okx/okbchain/x/feesplit/types"
	"github.com/stretchr/testify/suite"
)

type FeeSplitTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	app     *app.OKBChainApp
}

func TestFeeSplitTestSuite(t *testing.T) {
	suite.Run(t, new(FeeSplitTestSuite))
}

func (suite *FeeSplitTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 2, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.handler = feesplit.NewHandler(suite.app.FeeSplitKeeper)
	params := types.DefaultParams()
	params.EnableFeeSplit = true
	suite.app.FeeSplitKeeper.SetParams(suite.ctx, params)
}

func (suite *FeeSplitTestSuite) TestRegisterFeeSplit() {
	deployer := ethsecp256k1.GenerateAddress()
	fakeDeployer := ethsecp256k1.GenerateAddress()
	contract1 := crypto.CreateAddress(deployer, 1)
	factory1 := contract1
	factory2 := crypto.CreateAddress(factory1, 0)
	codeHash := common.Hex2Bytes("fa98cd094c09bb300de0037ba34e94f569b145ce8baa36ed863a08d7b7433f8d")

	contractBaseAcc := authtypes.NewBaseAccountWithAddress(contract1.Bytes())
	contractAccount := ethermint.EthAccount{
		BaseAccount: &contractBaseAcc,
		CodeHash:    codeHash,
	}
	deployerAccount := authtypes.NewBaseAccountWithAddress(deployer.Bytes())
	fakeDeployerAccount := authtypes.NewBaseAccountWithAddress(fakeDeployer.Bytes())

	testCases := []struct {
		name         string
		deployer     sdk.AccAddress
		withdraw     sdk.AccAddress
		contract     common.Address
		nonces       []uint64
		malleate     func()
		expPass      bool
		errorMessage string
	}{
		{
			"ok - contract deployed by EOA",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set deployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			true,
			"",
		},
		{
			"ok - contract deployed by factory in factory",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			crypto.CreateAddress(factory2, 1),
			[]uint64{1, 0, 1},
			func() {
				// set deployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				factoryContractBaseAcc := authtypes.NewBaseAccountWithAddress(crypto.CreateAddress(factory2, 1).Bytes())
				factoryContractAccount := ethermint.EthAccount{
					BaseAccount: &factoryContractBaseAcc,
					CodeHash:    codeHash,
				}
				suite.app.AccountKeeper.SetAccount(suite.ctx, factoryContractAccount)
			},
			true,
			"",
		},
		{
			"ok - omit withdraw address, it is stored as deployer string",
			sdk.AccAddress(deployer.Bytes()),
			nil,
			contract1,
			[]uint64{1},
			func() {
				// set deployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			true,
			"",
		},
		{
			"ok - deployer == withdraw, withdraw is stored as empty string",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set deployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			true,
			"",
		},
		{
			"not ok - deployer account not found",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set only contract account
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"deployer account not found",
		},
		{
			"not ok - deployer cannot be a contract",
			sdk.AccAddress(contract1.Bytes()),
			sdk.AccAddress(contract1.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set contract account
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"deployer cannot be a contract",
		},
		{
			"not ok - contract is already registered",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set deployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				msg := types.NewMsgRegisterFeeSplit(
					contract1,
					sdk.AccAddress(deployer.Bytes()),
					sdk.AccAddress(deployer.Bytes()),
					[]uint64{1},
				)
				suite.handler(suite.ctx, msg)
			},
			false,
			types.ErrFeeSplitAlreadyRegistered.Error(),
		},
		{
			"not ok - not contract deployer",
			sdk.AccAddress(fakeDeployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set deployer, fakeDeployer and contract accounts
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, &fakeDeployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"not contract deployer",
		},
		{
			"not ok - contract not deployed",
			sdk.AccAddress(deployer.Bytes()),
			sdk.AccAddress(deployer.Bytes()),
			contract1,
			[]uint64{1},
			func() {
				// set deployer account
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
			},
			false,
			"no contract code found at address",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()

			msg := types.NewMsgRegisterFeeSplit(tc.contract, tc.deployer, tc.withdraw, tc.nonces)

			res, err := suite.handler(suite.ctx, msg)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res)

				feeSplit, ok := suite.app.FeeSplitKeeper.GetFeeSplit(suite.ctx, tc.contract)
				suite.Require().True(ok, "unregistered feeSplit")
				suite.Require().Equal(tc.contract, feeSplit.ContractAddress, "wrong contract")
				suite.Require().Equal(tc.deployer, feeSplit.DeployerAddress, "wrong deployer")
				if tc.withdraw.String() != tc.deployer.String() && !tc.withdraw.Empty() {
					suite.Require().Equal(tc.withdraw, feeSplit.WithdrawerAddress, "wrong withdraw address")
				} else {
					suite.Require().Equal(tc.deployer, feeSplit.WithdrawerAddress, "wrong withdraw address")
				}
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().Contains(err.Error(), tc.errorMessage)
			}
		})
	}
}

func (suite *FeeSplitTestSuite) TestUpdateFeeSplit() {
	deployer := ethsecp256k1.GenerateAddress()
	deployerAddr := sdk.AccAddress(deployer.Bytes())
	withdrawer := sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	newWithdrawer := sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	contract1 := crypto.CreateAddress(deployer, 1)
	codeHash := common.Hex2Bytes("fa98cd094c09bb300de0037ba34e94f569b145ce8baa36ed863a08d7b7433f8d")

	contractBaseAcc := authtypes.NewBaseAccountWithAddress(contract1.Bytes())
	contractAccount := ethermint.EthAccount{
		BaseAccount: &contractBaseAcc,
		CodeHash:    codeHash,
	}
	deployerAccount := authtypes.NewBaseAccountWithAddress(deployer.Bytes())

	testCases := []struct {
		name          string
		deployer      sdk.AccAddress
		withdraw      sdk.AccAddress
		newWithdrawer sdk.AccAddress
		contract      common.Address
		nonces        []uint64
		malleate      func()
		expPass       bool
		errorMessage  string
	}{
		{
			"ok - change withdrawer to deployer",
			deployerAddr,
			withdrawer,
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// Prepare
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			true,
			"",
		},
		{
			"ok - change withdrawer to newWithdrawer",
			deployerAddr,
			withdrawer,
			newWithdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// Prepare
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			true,
			"",
		},
		{
			"fail - feesplit disabled",
			deployerAddr,
			withdrawer,
			newWithdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)

				params := types.DefaultParams()
				params.EnableFeeSplit = false
				suite.app.FeeSplitKeeper.SetParams(suite.ctx, params)
			},
			false,
			"",
		},
		{
			"fail - contract not registered",
			deployerAddr,
			withdrawer,
			newWithdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"",
		},
		{
			"fail - deployer not the one registered",
			newWithdrawer,
			withdrawer,
			newWithdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			false,
			"",
		},
		{
			"fail - everything is the same",
			deployerAddr,
			withdrawer,
			withdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			false,
			"",
		},
		{
			"fail - previously cancelled contract",
			deployerAddr,
			withdrawer,
			withdrawer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)

				msgCancel := types.NewMsgCancelFeeSplit(contract1, deployerAddr)
				_, err = suite.handler(suite.ctx, msgCancel)
				suite.Require().NoError(err)
			},
			false,
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			tc.malleate()

			msgUpdate := types.NewMsgUpdateFeeSplit(tc.contract, tc.deployer, tc.newWithdrawer)

			res, err := suite.handler(suite.ctx, msgUpdate)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res)

				feeSplit, ok := suite.app.FeeSplitKeeper.GetFeeSplit(suite.ctx, tc.contract)
				suite.Require().True(ok, "unregistered feeSplit")
				suite.Require().Equal(tc.contract, feeSplit.ContractAddress, "wrong contract")
				suite.Require().Equal(tc.deployer, feeSplit.DeployerAddress, "wrong deployer")

				found := suite.app.FeeSplitKeeper.IsWithdrawerMapSet(suite.ctx, tc.withdraw, tc.contract)
				suite.Require().False(found)
				if tc.newWithdrawer.String() != tc.deployer.String() {
					suite.Require().Equal(tc.newWithdrawer, feeSplit.WithdrawerAddress, "wrong withdraw address")
					found := suite.app.FeeSplitKeeper.IsWithdrawerMapSet(suite.ctx, tc.newWithdrawer, tc.contract)
					suite.Require().True(found)
				} else {
					suite.Require().Equal(tc.deployer, feeSplit.WithdrawerAddress, "wrong withdraw address")
					found := suite.app.FeeSplitKeeper.IsWithdrawerMapSet(suite.ctx, tc.newWithdrawer, tc.contract)
					suite.Require().True(found)
				}
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().Contains(err.Error(), tc.errorMessage)
			}
		})
	}
}

func (suite *FeeSplitTestSuite) TestCancelFeeSplit() {
	deployer := ethsecp256k1.GenerateAddress()
	deployerAddr := sdk.AccAddress(deployer.Bytes())
	withdrawer := sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	fakeDeployer := sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	contract1 := crypto.CreateAddress(deployer, 1)
	codeHash := common.Hex2Bytes("fa98cd094c09bb300de0037ba34e94f569b145ce8baa36ed863a08d7b7433f8d")

	contractBaseAcc := authtypes.NewBaseAccountWithAddress(contract1.Bytes())
	contractAccount := ethermint.EthAccount{
		BaseAccount: &contractBaseAcc,
		CodeHash:    codeHash,
	}
	deployerAccount := authtypes.NewBaseAccountWithAddress(deployerAddr)

	testCases := []struct {
		name         string
		deployer     sdk.AccAddress
		contract     common.Address
		nonces       []uint64
		malleate     func()
		expPass      bool
		errorMessage string
	}{
		{
			"ok - cancelled",
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// Prepare
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})

				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			true,
			"",
		},
		{
			"ok - cancelled - no withdrawer",
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// Prepare
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, deployerAddr, []uint64{1})

				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			true,
			"",
		},
		{
			"fail - feesplit disabled",
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)

				params := types.DefaultParams()
				params.EnableFeeSplit = false
				suite.app.FeeSplitKeeper.SetParams(suite.ctx, params)
			},
			false,
			"",
		},
		{
			"fail - contract not registered",
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"",
		},
		{
			"fail - deployer not the one registered",
			fakeDeployer,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)

				// register contract
				msg := types.NewMsgRegisterFeeSplit(contract1, deployerAddr, withdrawer, []uint64{1})
				_, err := suite.handler(suite.ctx, msg)
				suite.Require().NoError(err)
			},
			false,
			"",
		},
		{
			"fail - everything is the same",
			deployerAddr,
			contract1,
			[]uint64{1},
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx, &deployerAccount)
				suite.app.AccountKeeper.SetAccount(suite.ctx, contractAccount)
			},
			false,
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			tc.malleate()

			msgCancel := types.NewMsgCancelFeeSplit(tc.contract, tc.deployer)
			res, err := suite.handler(suite.ctx, msgCancel)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res)

				_, ok := suite.app.FeeSplitKeeper.GetFeeSplit(suite.ctx, tc.contract)
				suite.Require().False(ok, "registered feeSplit")

				found := suite.app.FeeSplitKeeper.IsWithdrawerMapSet(suite.ctx, withdrawer, tc.contract)
				suite.Require().False(found)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().Contains(err.Error(), tc.errorMessage)
			}
		})
	}
}
