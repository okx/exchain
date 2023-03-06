package keeper_test

import (
	"fmt"

	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/query"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/feesplit/types"
)

func (suite *KeeperTestSuite) TestFeeSplits() {
	var (
		req    *types.QueryFeeSplitsRequest
		expRes *types.QueryFeeSplitsResponse
	)

	testCases := []struct {
		name     string
		path     []string
		malleate func()
		expPass  bool
	}{
		{
			"no fee infos registered",
			[]string{types.QueryFeeSplits},
			func() {
				req = &types.QueryFeeSplitsRequest{}
				expRes = &types.QueryFeeSplitsResponse{Pagination: &query.PageResponse{}}
			},
			true,
		},
		{
			"1 fee infos registered w/pagination",
			[]string{types.QueryFeeSplits},
			func() {
				req = &types.QueryFeeSplitsRequest{
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true},
				}
				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)

				expRes = &types.QueryFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 1},
					FeeSplits: []types.FeeSplitWithShare{
						{
							ContractAddress:   contract.Hex(),
							DeployerAddress:   deployer.String(),
							WithdrawerAddress: withdraw.String(),
							Share:             suite.app.FeeSplitKeeper.GetParams(suite.ctx).DeveloperShares,
						},
					},
				}
			},
			true,
		},
		{
			"2 fee infos registered wo/pagination",
			[]string{types.QueryFeeSplits},
			func() {
				req = &types.QueryFeeSplitsRequest{}
				contract2 := ethsecp256k1.GenerateAddress()
				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				feeSplit2 := types.NewFeeSplit(contract2, deployer, nil)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit2)

				expRes = &types.QueryFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 2},
					FeeSplits: []types.FeeSplitWithShare{
						{
							ContractAddress:   contract.Hex(),
							DeployerAddress:   deployer.String(),
							WithdrawerAddress: withdraw.String(),
							Share:             suite.app.FeeSplitKeeper.GetParams(suite.ctx).DeveloperShares,
						},
						{
							ContractAddress:   contract2.Hex(),
							DeployerAddress:   deployer.String(),
							WithdrawerAddress: deployer.String(),
							Share:             suite.app.FeeSplitKeeper.GetParams(suite.ctx).DeveloperShares,
						},
					},
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			data, err := suite.app.Codec().MarshalJSON(req)
			suite.Require().NoError(err)
			res, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{Data: data})
			if tc.expPass {
				suite.Require().NoError(err)

				var resp types.QueryFeeSplitsResponse
				err = suite.app.Codec().UnmarshalJSON(res, &resp)
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Pagination, resp.Pagination)
				suite.Require().ElementsMatch(expRes.FeeSplits, resp.FeeSplits)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestFee() {
	var (
		req    *types.QueryFeeSplitRequest
		expRes *types.QueryFeeSplitResponse
	)

	testCases := []struct {
		name     string
		path     []string
		malleate func()
		expPass  bool
	}{
		{
			"empty contract address",
			[]string{types.QueryFeeSplit},
			func() {
				req = &types.QueryFeeSplitRequest{}
				expRes = &types.QueryFeeSplitResponse{}
			},
			false,
		},
		{
			"invalid contract address",
			[]string{types.QueryFeeSplit},
			func() {
				req = &types.QueryFeeSplitRequest{
					ContractAddress: "1234",
				}
				expRes = &types.QueryFeeSplitResponse{}
			},
			false,
		},
		{
			"fee info not found",
			[]string{types.QueryFeeSplit},
			func() {
				req = &types.QueryFeeSplitRequest{
					ContractAddress: contract.String(),
				}
				expRes = &types.QueryFeeSplitResponse{}
			},
			false,
		},
		{
			"fee info found",
			[]string{types.QueryFeeSplit},
			func() {
				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)

				req = &types.QueryFeeSplitRequest{
					ContractAddress: contract.Hex(),
				}
				expRes = &types.QueryFeeSplitResponse{FeeSplit: types.FeeSplitWithShare{
					ContractAddress:   contract.Hex(),
					DeployerAddress:   deployer.String(),
					WithdrawerAddress: withdraw.String(),
					Share:             suite.app.FeeSplitKeeper.GetParams(suite.ctx).DeveloperShares,
				}}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			data, err := suite.app.Codec().MarshalJSON(req)
			suite.Require().NoError(err)
			res, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{Data: data})
			if tc.expPass {
				suite.Require().NoError(err)

				var resp types.QueryFeeSplitResponse
				err = suite.app.Codec().UnmarshalJSON(res, &resp)
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, &resp)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDeployerFees() {
	var (
		req    *types.QueryDeployerFeeSplitsRequest
		expRes *types.QueryDeployerFeeSplitsResponse
	)

	testCases := []struct {
		name     string
		path     []string
		malleate func()
		expPass  bool
	}{
		{
			"no contract registered",
			[]string{types.QueryDeployerFeeSplits},
			func() {
				req = &types.QueryDeployerFeeSplitsRequest{}
				expRes = &types.QueryDeployerFeeSplitsResponse{Pagination: &query.PageResponse{}}
			},
			false,
		},
		{
			"invalid deployer address",
			[]string{types.QueryDeployerFeeSplits},
			func() {
				req = &types.QueryDeployerFeeSplitsRequest{
					DeployerAddress: "123",
				}
				expRes = &types.QueryDeployerFeeSplitsResponse{Pagination: &query.PageResponse{}}
			},
			false,
		},
		{
			"1 fee registered w/pagination",
			[]string{types.QueryDeployerFeeSplits},
			func() {
				req = &types.QueryDeployerFeeSplitsRequest{
					Pagination:      &query.PageRequest{Limit: 10, CountTotal: true},
					DeployerAddress: deployer.String(),
				}

				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer, contract)
				suite.app.FeeSplitKeeper.SetWithdrawerMap(suite.ctx, withdraw, contract)

				expRes = &types.QueryDeployerFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 1},
					ContractAddresses: []string{
						contract.Hex(),
					},
				}
			},
			true,
		},
		{
			"2 fee infos registered for one contract wo/pagination",
			[]string{types.QueryDeployerFeeSplits},
			func() {
				req = &types.QueryDeployerFeeSplitsRequest{
					DeployerAddress: deployer.String(),
				}
				contract2 := ethsecp256k1.GenerateAddress()
				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer, contract)
				suite.app.FeeSplitKeeper.SetWithdrawerMap(suite.ctx, withdraw, contract)

				feeSplit2 := types.NewFeeSplit(contract2, deployer, nil)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit2)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer, contract2)

				expRes = &types.QueryDeployerFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 2},
					ContractAddresses: []string{
						contract.Hex(),
						contract2.Hex(),
					},
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			data, err := suite.app.Codec().MarshalJSON(req)
			suite.Require().NoError(err)
			res, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{Data: data})
			if tc.expPass {
				suite.Require().NoError(err)

				var resp types.QueryDeployerFeeSplitsResponse
				err = suite.app.Codec().UnmarshalJSON(res, &resp)
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Pagination, resp.Pagination)
				suite.Require().ElementsMatch(expRes.ContractAddresses, resp.ContractAddresses)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestWithdrawerFeeSplits() {
	var (
		req    *types.QueryWithdrawerFeeSplitsRequest
		expRes *types.QueryWithdrawerFeeSplitsResponse
	)

	testCases := []struct {
		name     string
		path     []string
		malleate func()
		expPass  bool
	}{
		{
			"no contract registered",
			[]string{types.QueryWithdrawerFeeSplits},
			func() {
				req = &types.QueryWithdrawerFeeSplitsRequest{}
				expRes = &types.QueryWithdrawerFeeSplitsResponse{Pagination: &query.PageResponse{}}
			},
			false,
		},
		{
			"invalid withdraw address",
			[]string{types.QueryWithdrawerFeeSplits},
			func() {
				req = &types.QueryWithdrawerFeeSplitsRequest{
					WithdrawerAddress: "123",
				}
				expRes = &types.QueryWithdrawerFeeSplitsResponse{Pagination: &query.PageResponse{}}
			},
			false,
		},
		{
			"1 fee registered w/pagination",
			[]string{types.QueryWithdrawerFeeSplits},
			func() {
				req = &types.QueryWithdrawerFeeSplitsRequest{
					Pagination:        &query.PageRequest{Limit: 10, CountTotal: true},
					WithdrawerAddress: withdraw.String(),
				}

				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer, contract)
				suite.app.FeeSplitKeeper.SetWithdrawerMap(suite.ctx, withdraw, contract)

				expRes = &types.QueryWithdrawerFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 1},
					ContractAddresses: []string{
						contract.Hex(),
					},
				}
			},
			true,
		},
		{
			"2 fees registered for one withdraw address wo/pagination",
			[]string{types.QueryWithdrawerFeeSplits},
			func() {
				req = &types.QueryWithdrawerFeeSplitsRequest{
					WithdrawerAddress: withdraw.String(),
				}
				contract2 := ethsecp256k1.GenerateAddress()
				deployer2 := sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())

				feeSplit := types.NewFeeSplit(contract, deployer, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer, contract)
				suite.app.FeeSplitKeeper.SetWithdrawerMap(suite.ctx, withdraw, contract)

				feeSplit2 := types.NewFeeSplit(contract2, deployer2, withdraw)
				suite.app.FeeSplitKeeper.SetFeeSplit(suite.ctx, feeSplit2)
				suite.app.FeeSplitKeeper.SetDeployerMap(suite.ctx, deployer2, contract2)
				suite.app.FeeSplitKeeper.SetWithdrawerMap(suite.ctx, withdraw, contract2)

				expRes = &types.QueryWithdrawerFeeSplitsResponse{
					Pagination: &query.PageResponse{Total: 2},
					ContractAddresses: []string{
						contract.Hex(),
						contract2.Hex(),
					},
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			data, err := suite.app.Codec().MarshalJSON(req)
			suite.Require().NoError(err)
			res, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{Data: data})
			if tc.expPass {
				suite.Require().NoError(err)

				var resp types.QueryWithdrawerFeeSplitsResponse
				err = suite.app.Codec().UnmarshalJSON(res, &resp)
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Pagination, resp.Pagination)
				suite.Require().ElementsMatch(expRes.ContractAddresses, resp.ContractAddresses)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryParams() {
	expParams := types.DefaultParams()

	res, err := suite.querier(suite.ctx, []string{types.QueryParameters}, abci.RequestQuery{Data: nil})
	suite.Require().NoError(err)
	var resp types.QueryParamsResponse
	err = suite.app.Codec().UnmarshalJSON(res, &resp)
	suite.Require().NoError(err)
	suite.Require().Equal(expParams, resp.Params)
}
