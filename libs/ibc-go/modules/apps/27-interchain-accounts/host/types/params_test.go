package types_test

import (
	"testing"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/types"
	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())
	require.NoError(t, types.NewParams(false, []string{}).Validate())
}
