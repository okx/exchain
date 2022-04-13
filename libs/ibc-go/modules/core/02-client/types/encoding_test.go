package types_test

import (
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	ibctmtypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
)

func (suite *TypesTestSuite) TestMarshalHeader() {

	cdc := suite.chainA.App.AppCodec()
	h := &ibctmtypes.Header{
		TrustedHeight: types.NewHeight(4, 100),
	}

	// marshal header
	bz, err := types.MarshalHeader(cdc, h)
	suite.Require().NoError(err)

	// unmarshal header
	// newHeader, err := types.UnmarshalHeader(cdc, bz)
	newHeader, err := types.UnmarshalHeader(cdc.GetProtocMarshal(), bz)
	suite.Require().NoError(err)

	suite.Require().Equal(h, newHeader)

	// use invalid bytes
	// invalidHeader, err := types.UnmarshalHeader(cdc, []byte("invalid bytes"))
	invalidHeader, err := types.UnmarshalHeader(cdc.GetProtocMarshal(), []byte("invalid bytes"))
	suite.Require().Error(err)
	suite.Require().Nil(invalidHeader)

}
