package types_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/tests/mocks"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func TestChainAnteDecorators(t *testing.T) {
	t.Parallel()
	// test panic
	require.Nil(t, sdk.ChainAnteDecorators([]sdk.AnteDecorator{}...))

	ctx, tx := sdk.Context{}, sdk.Tx(nil)
	mockCtrl := gomock.NewController(t)
	mockAnteDecorator1 := mocks.NewMockAnteDecorator(mockCtrl)
	mockAnteDecorator1.EXPECT().AnteHandle(gomock.Eq(ctx), gomock.Eq(tx), true, gomock.Any()).AnyTimes()
	sdk.ChainAnteDecorators(mockAnteDecorator1)(ctx, tx, true)

	mockAnteDecorator2 := mocks.NewMockAnteDecorator(mockCtrl)
	mockAnteDecorator1.EXPECT().AnteHandle(gomock.Eq(ctx), gomock.Eq(tx), true, mockAnteDecorator2).AnyTimes()
	mockAnteDecorator2.EXPECT().AnteHandle(gomock.Eq(ctx), gomock.Eq(tx), true, nil).AnyTimes()
	sdk.ChainAnteDecorators(mockAnteDecorator1, mockAnteDecorator2)
}
