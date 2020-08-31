package pulsarclient

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestPb(t *testing.T) {
	name := "xkb_oxb"
	marketId := marketIdMap[name]
	size := rand.Float64()
	price := rand.Float64()
	timestamp := time.Now().Unix()
	i32 := int32(1)
	i64 := int64(1)
	ui64 := uint64(1)

	matchResultMsg := MatchResultMsg{
		BizType:          &bizType,
		MarketId:         &marketId,
		MarketType:       &marketType,
		Size:             &size,
		Price:            &price,
		CreatedTime:      &timestamp,
		InstrumentId:     &marketId,
		InstrumentName:   &name,
		IsCalc:           &iscalc,
		UserId:           &i64,
		BrokerId:         &i32,
		OrderId:          &ui64,
		TradeId:          &i64,
		Tradeside:        &i32,
		OppositeUserId:   &i64,
		OppositeBrokerId: &i32,
		OppositeOrderId:  &ui64,
		OrderSide:        &i32,
		EventId:          &i64,
		EventType:        &i32,
	}
	getAllmethod(&matchResultMsg)
	matchResultMsg.Descriptor()
	matchResultMsg.XXX_DiscardUnknown()
	matchResultMsg.String()

	msg, err := proto.Marshal(&matchResultMsg)
	require.NoError(t, err)
	err = matchResultMsg.XXX_Unmarshal(msg)
	require.NoError(t, err)

	matchResultMsg.Reset()
	getAllmethod(&matchResultMsg)
}

func getAllmethod(matchResultMsg *MatchResultMsg) {
	matchResultMsg.GetInstrumentId()
	matchResultMsg.GetSize()
	matchResultMsg.GetPrice()
	matchResultMsg.GetCreatedTime()
	matchResultMsg.GetMarketType()
	matchResultMsg.GetBrokerId()
	matchResultMsg.GetOrderId()
	matchResultMsg.GetBizType()
	matchResultMsg.GetInstrumentName()
	matchResultMsg.GetTradeId()
	matchResultMsg.GetOppositeBrokerId()
	matchResultMsg.GetUserId()
	matchResultMsg.GetOppositeOrderId()
	matchResultMsg.GetEventId()
	matchResultMsg.GetEventType()
	matchResultMsg.GetIsCalc()
	matchResultMsg.GetOrderSide()
	matchResultMsg.GetTradeside()
	matchResultMsg.GetMarketId()
	matchResultMsg.GetOppositeUserId()
}
