package types

import (
	"math/big"
	"strings"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/suite"
)

const (
	expectedTitle                                           = "default title"
	expectedDescription                                     = "default description"
	expectedManageContractDeploymentWhitelistProposalString = `ManageContractDeploymentWhitelistProposal:
 Title:					default title
 Description:        	default description
 Type:                	ManageContractDeploymentWhitelist
 IsAdded:				true
 DistributorAddrs:
						ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqm2k6w2
						ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxuz0nc`
	expectedManageContractBlockedListProposalString = `ManageContractBlockedListProposal:
 Title:					default title
 Description:        	default description
 Type:                	ManageContractBlockedList
 IsAdded:				true
 ContractAddrs:
						ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqm2k6w2
						ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxuz0nc`
	expectedManageContractMethodBlockedListProposalString = `ManageContractMethodBlockedListProposal:
 Title:					default title
 Description:        	default description
 Type:                	ManageContractMethodBlockedList
 IsAdded:				true
 ContractList:
						Address: ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqm2k6w2
Method List:
Sign: 0x11111111Extra: TEST1
Sign: 0x22222222Extra: TEST2
						Address: ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxuz0nc
Method List:
Sign: 0x33333333Extra: TEST3
Sign: 0x44444444Extra: TEST4`
)

type ProposalTestSuite struct {
	suite.Suite
	strBuilder          strings.Builder
	addrs               AddressList
	blockedContractList BlockedContractList
}

func TestProposalTestSuite(t *testing.T) {
	proposalTestSuite := ProposalTestSuite{
		addrs: AddressList{
			ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
			ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		},
		blockedContractList: BlockedContractList{
			BlockedContract{
				Address: ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
				BlockMethods: ContractMethods{
					ContractMethod{Sign: "0x11111111", Extra: "TEST1"},
					ContractMethod{Sign: "0x22222222", Extra: "TEST2"},
				},
			},
			BlockedContract{
				Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
				BlockMethods: ContractMethods{
					ContractMethod{Sign: "0x33333333", Extra: "TEST3"},
					ContractMethod{Sign: "0x44444444", Extra: "TEST4"},
				},
			},
		},
	}
	suite.Run(t, &proposalTestSuite)
}

func (suite *ProposalTestSuite) TestProposal_ManageContractDeploymentWhitelistProposal() {
	proposal := NewManageContractDeploymentWhitelistProposal(
		expectedTitle,
		expectedDescription,
		suite.addrs,
		true,
	)

	suite.Require().Equal(expectedTitle, proposal.GetTitle())
	suite.Require().Equal(expectedDescription, proposal.GetDescription())
	suite.Require().Equal(RouterKey, proposal.ProposalRoute())
	suite.Require().Equal(proposalTypeManageContractDeploymentWhitelist, proposal.ProposalType())
	suite.Require().Equal(expectedManageContractDeploymentWhitelistProposalString, proposal.String())

	testCases := []struct {
		msg           string
		prepare       func()
		expectedError bool
	}{
		{
			"pass",
			func() {},
			false,
		},
		{
			"empty title",
			func() {
				proposal.Title = ""
			},
			true,
		},
		{
			"overlong title",
			func() {
				for i := 0; i < govtypes.MaxTitleLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Title = suite.strBuilder.String()
			},
			true,
		},
		{
			"empty description",
			func() {
				proposal.Description = ""
				proposal.Title = expectedTitle
			},
			true,
		},
		{
			"overlong description",
			func() {
				suite.strBuilder.Reset()
				for i := 0; i < govtypes.MaxDescriptionLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Description = suite.strBuilder.String()
			},
			true,
		},
		{
			"duplicated distributor addresses",
			func() {
				// add a duplicated address into DistributorAddrs
				proposal.DistributorAddrs = append(proposal.DistributorAddrs, proposal.DistributorAddrs[0])
				proposal.Description = expectedDescription
			},
			true,
		},
		{
			"empty distributor addresses",
			func() {
				proposal.DistributorAddrs = nil
			},
			true,
		},
		{
			"oversize distributor addresses",
			func() {
				for i := int64(0); i <= maxAddressListLength; i++ {
					testAddr := ethcmn.BigToAddress(big.NewInt(i)).Bytes()
					proposal.DistributorAddrs = append(proposal.DistributorAddrs, testAddr)
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()

			err := proposal.ValidateBasic()

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *ProposalTestSuite) TestProposal_ManageContractBlockedListProposal() {
	proposal := NewManageContractBlockedListProposal(
		expectedTitle,
		expectedDescription,
		suite.addrs,
		true,
	)

	suite.Require().Equal(expectedTitle, proposal.GetTitle())
	suite.Require().Equal(expectedDescription, proposal.GetDescription())
	suite.Require().Equal(RouterKey, proposal.ProposalRoute())
	suite.Require().Equal(proposalTypeManageContractBlockedList, proposal.ProposalType())
	suite.Require().Equal(expectedManageContractBlockedListProposalString, proposal.String())

	testCases := []struct {
		msg           string
		prepare       func()
		expectedError bool
	}{
		{
			"pass",
			func() {},
			false,
		},
		{
			"empty title",
			func() {
				proposal.Title = ""
			},
			true,
		},
		{
			"overlong title",
			func() {
				for i := 0; i < govtypes.MaxTitleLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Title = suite.strBuilder.String()
			},
			true,
		},
		{
			"empty description",
			func() {
				proposal.Description = ""
				proposal.Title = expectedTitle
			},
			true,
		},
		{
			"overlong description",
			func() {
				suite.strBuilder.Reset()
				for i := 0; i < govtypes.MaxDescriptionLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Description = suite.strBuilder.String()
			},
			true,
		},
		{
			"duplicated contract addresses",
			func() {
				// add a duplicated address into ContractAddrs
				proposal.ContractAddrs = append(proposal.ContractAddrs, proposal.ContractAddrs[0])
				proposal.Description = expectedDescription
			},
			true,
		},
		{
			"empty contract addresses",
			func() {
				proposal.ContractAddrs = nil
			},
			true,
		},
		{
			"oversize contract addresses",
			func() {
				for i := int64(0); i <= maxAddressListLength; i++ {
					testAddr := ethcmn.BigToAddress(big.NewInt(i)).Bytes()
					proposal.ContractAddrs = append(proposal.ContractAddrs, testAddr)
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()

			err := proposal.ValidateBasic()

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *ProposalTestSuite) TestProposal_ManageContractMethodBlockedListProposal() {
	proposal := NewManageContractMethodBlockedListProposal(
		expectedTitle,
		expectedDescription,
		suite.blockedContractList,
		true,
	)

	suite.Require().Equal(expectedTitle, proposal.GetTitle())
	suite.Require().Equal(expectedDescription, proposal.GetDescription())
	suite.Require().Equal(RouterKey, proposal.ProposalRoute())
	suite.Require().Equal(proposalTypeManageContractMethodBlockedList, proposal.ProposalType())
	suite.Require().Equal(expectedManageContractMethodBlockedListProposalString, proposal.String())

	testCases := []struct {
		msg           string
		prepare       func()
		expectedError bool
	}{
		{
			"pass",
			func() {},
			false,
		},
		{
			"empty title",
			func() {
				proposal.Title = ""
			},
			true,
		},
		{
			"overlong title",
			func() {
				for i := 0; i < govtypes.MaxTitleLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Title = suite.strBuilder.String()
			},
			true,
		},
		{
			"empty description",
			func() {
				proposal.Description = ""
				proposal.Title = expectedTitle
			},
			true,
		},
		{
			"overlong description",
			func() {
				suite.strBuilder.Reset()
				for i := 0; i < govtypes.MaxDescriptionLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Description = suite.strBuilder.String()
			},
			true,
		},
		{
			"duplicated contract addresses",
			func() {
				// add a duplicated address into ContractAddrs
				proposal.ContractList = append(proposal.ContractList, proposal.ContractList[0])
				proposal.Description = expectedDescription
			},
			true,
		},
		{
			"empty contract addresses",
			func() {
				proposal.ContractList = nil
			},
			true,
		},
		{
			"oversize contract addresses",
			func() {
				for i := int64(0); i <= maxAddressListLength; i++ {
					testAddr := ethcmn.BigToAddress(big.NewInt(i)).Bytes()
					testbc := BlockedContract{Address: testAddr, BlockMethods: nil}
					proposal.ContractList = append(proposal.ContractList, testbc)
				}
			},
			true,
		},
		{
			"empty contract addresses",
			func() {
				for i := int64(0); i <= maxAddressListLength; i++ {
					testbc := BlockedContract{Address: nil, BlockMethods: nil}
					proposal.ContractList = append(proposal.ContractList, testbc)
				}
			},
			true,
		},
		{
			"duplicated contract method",
			func() {
				proposal.ContractList[0].BlockMethods = append(proposal.ContractList[0].BlockMethods, proposal.ContractList[0].BlockMethods...)
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()

			err := proposal.ValidateBasic()

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
