package types

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDepthBookInsertAndRemove(t *testing.T) {
	depthBook := &DepthBook{}
	// Test insert order
	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.5"), depthBook.Items[0].Price)

	order2 := MockOrder("", TestTokenPair, SellOrder, "0.6", "2.1")
	depthBook.InsertOrder(order2)
	depthBook.InsertOrder(order2)
	require.EqualValues(t, 2, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("4.2"), depthBook.Items[0].SellQuantity)

	order3 := MockOrder("", TestTokenPair, BuyOrder, "0.4", "1.5")
	depthBook.InsertOrder(order3)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.4"), depthBook.Items[2].Price)
	require.EqualValues(t, sdk.MustNewDecFromStr("1.5"), depthBook.Items[2].BuyQuantity)

	depthBook.InsertOrder(order1)
	require.EqualValues(t, sdk.MustNewDecFromStr("2.2"), depthBook.Items[1].BuyQuantity)

	// Test remove order
	depthBook.RemoveOrder(order2)
	require.EqualValues(t, 3, len(depthBook.Items))
	depthBook.RemoveOrder(order1)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.0"), depthBook.Items[0].BuyQuantity)

	depthBook.RemoveOrder(order1)
	depthBook.RemoveOrder(order2)
	depthBook.RemoveOrder(order3)
	require.EqualValues(t, 0, len(depthBook.Items))
}

func TestInsertOrder(t *testing.T) {
	depthBook := &DepthBook{}

	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.5"), depthBook.Items[0].Price)
}

func TestRemoveOrder(t *testing.T) {
	depthBook := &DepthBook{}

	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)

	depthBook.RemoveOrder(order1)
	require.EqualValues(t, 0, len(depthBook.Items))
}

func TestSub(t *testing.T) {
	depthBook := &DepthBook{}

	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)

	depthBook.Sub(0, sdk.MustNewDecFromStr("1.0"), BuyOrder)
	require.EqualValues(t, sdk.MustNewDecFromStr("0.1"), depthBook.Items[0].BuyQuantity)

	depthBook.Sub(0, sdk.MustNewDecFromStr("0.5"), SellOrder)
	require.EqualValues(t, sdk.MustNewDecFromStr("-0.5"), depthBook.Items[0].SellQuantity)
}

func TestRemoveIfEmpty(t *testing.T) {
	depthBook := &DepthBook{}

	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)

	depthBook.Sub(0, sdk.MustNewDecFromStr("1.1"), BuyOrder)
	depthBook.RemoveIfEmpty(0)

	require.EqualValues(t, 0, len(depthBook.Items))
}

func TestCopy(t *testing.T) {
	depthBook := &DepthBook{}

	order1 := MockOrder("", TestTokenPair, BuyOrder, "0.5", "1.1")
	depthBook.InsertOrder(order1)

	bookCopy := depthBook.Copy()
	bookCopy.Sub(0, sdk.MustNewDecFromStr("1.1"), BuyOrder)
	bookCopy.RemoveIfEmpty(0)

	require.EqualValues(t, 1, len(depthBook.Items))
	require.EqualValues(t, sdk.MustNewDecFromStr("0.5"), depthBook.Items[0].Price)
}
