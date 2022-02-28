package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	codectypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/codec/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/std"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil/testdata"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
