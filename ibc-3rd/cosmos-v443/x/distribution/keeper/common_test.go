package keeper_test

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	authtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/distribution/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)

	valConsPk1 = PKS[0]
	valConsPk2 = PKS[1]
	valConsPk3 = PKS[2]

	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())

	distrAcc = authtypes.NewEmptyModuleAccount(types.ModuleName)
)
