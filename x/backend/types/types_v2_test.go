package types

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex"
	"github.com/stretchr/testify/require"
)

func TestDefaultTickerV2(t *testing.T) {
	instrumentID := "101"
	defaultTicker := DefaultTickerV2(instrumentID)
	require.Equal(t, defaultTicker.InstrumentID, instrumentID)
}

func TestConvertOrderToOrderV2(t *testing.T) {
	orderV1 := Order{
		OrderID:        "test-order_id",
		Timestamp:      time.Now().Unix(),
		Quantity:       "100",
		RemainQuantity: "10",
		FilledAvgPrice: "5",
	}
	orderV2 := ConvertOrderToOrderV2(orderV1)
	require.Equal(t, orderV1.OrderID, orderV2.OrderID)

	timeStamp := time.Unix(orderV1.Timestamp, 0).UTC().Format("2006-01-02T15:04:05.000Z")
	require.Equal(t, timeStamp, orderV2.Timestamp)

	filledSizeDec := sdk.MustNewDecFromStr(orderV1.Quantity).Sub(sdk.MustNewDecFromStr(orderV1.RemainQuantity))
	require.Equal(t, filledSizeDec.String(), orderV2.FilledSize)

	filledNotionalDec := filledSizeDec.Mul(sdk.MustNewDecFromStr(orderV1.FilledAvgPrice))
	require.Equal(t, filledNotionalDec.String(), orderV2.FilledNotional)
}

func TestConvertTokenPairToInstrumentV2(t *testing.T) {
	tokenPair := dex.GetBuiltInTokenPair()
	instrumentV2 := ConvertTokenPairToInstrumentV2(tokenPair)
	instrumentID := tokenPair.Name()
	require.Equal(t, instrumentID, instrumentV2.InstrumentID)

	fSizeIncrement := 1 / math.Pow10(int(tokenPair.MaxQuantityDigit))
	sizeIncrement := strings.TrimRight(fmt.Sprintf("%.10f", fSizeIncrement), "0")
	require.Equal(t, sizeIncrement, instrumentV2.SizeIncrement)

	fTickSize := 1 / math.Pow10(int(tokenPair.MaxPriceDigit))
	tickSize := strings.TrimRight(fmt.Sprintf("%.10f", fTickSize), "0")
	require.Equal(t, tickSize, instrumentV2.TickSize)
}
