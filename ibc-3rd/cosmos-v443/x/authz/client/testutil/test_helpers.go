package testutil

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil"
	clitestutil "github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil/cli"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil/network"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/authz/client/cli"
)

func ExecGrant(val *network.Validator, args []string) (testutil.BufferWriter, error) {
	cmd := cli.NewCmdGrantAuthorization()
	clientCtx := val.ClientCtx
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
}
