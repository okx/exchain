package types

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/okex/okchain/x/common"

	"github.com/stretchr/testify/require"
)

const (
	testQuantity = "1"
	testPrice    = "0.1"
	testOrderID  = "abc"
)

func TestMsgNewOrder(t *testing.T) {

	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg := NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	require.Nil(t, orderMsg.ValidateBasic())
	require.Equal(t, "order", orderMsg.Route())
	require.Equal(t, "new", orderMsg.Type())

	bytesMsg := orderMsg.GetSignBytes()
	resOrderMsg := &MsgNewOrder{}
	err = json.Unmarshal(bytesMsg, resOrderMsg)
	require.Nil(t, err)
	resAddr := orderMsg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
}

func TestMsgNewOrderInvalid(t *testing.T) {
	//nil sender
	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg := NewMsgNewOrder(nil, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//empty sender
	orderMsg = NewMsgNewOrder([]byte{}, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//empty product
	orderMsg = NewMsgNewOrder(addr, "", BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//invalid product
	orderMsg = NewMsgNewOrder(addr, "btc"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//okb in the left
	orderMsg = NewMsgNewOrder(addr, common.NativeToken+"_btc", BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	//invalid side
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, "abc", testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//zero price
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, "0", testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//zero quantity
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, "0")
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//negative price
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, "-1", testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//negative quantity
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, "-1")
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//same product
	orderMsg = NewMsgNewOrder(addr, common.TestToken+"_"+common.TestToken, BuyOrder, testPrice, "-1")
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)
}

func TestMsgCancelOrder(t *testing.T) {
	orderID := testOrderID
	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg := NewMsgCancelOrder(addr, orderID)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)
	require.Equal(t, "order", orderMsg.Route())
	require.Equal(t, "cancel", orderMsg.Type())

	bytesMsg := orderMsg.GetSignBytes()
	resOrderMsg := &MsgNewOrder{}
	err = json.Unmarshal(bytesMsg, resOrderMsg)
	require.Nil(t, err)
	resAddr := orderMsg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
}

func TestMsgCancelOrderInvalid(t *testing.T) {
	//empty sender
	orderID := testOrderID
	addr, err := hex.DecodeString("")
	require.Nil(t, err)
	orderMsg := NewMsgCancelOrder(addr, orderID)
	sdkErr := orderMsg.ValidateBasic()
	require.NotNil(t, sdkErr)

	// empty orderID
	orderID = ""
	addr, err = hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg = NewMsgCancelOrder(addr, orderID)
	sdkErr = orderMsg.ValidateBasic()
	require.NotNil(t, sdkErr)
}

func TestMsgMultiNewOrder(t *testing.T) {
	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderItems := NewOrderItem("btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	orderMsg := NewMsgNewOrders(addr, []OrderItem{orderItems})
	require.Nil(t, orderMsg.ValidateBasic())
	require.Equal(t, "order", orderMsg.Route())
	require.Equal(t, "new", orderMsg.Type())

	bytesMsg := orderMsg.GetSignBytes()
	resOrderMsg := &OrderItem{}
	err = json.Unmarshal(bytesMsg, resOrderMsg)
	require.Nil(t, err)
	resAddr := orderMsg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
}

func TestMsgMultiNewOrderInvalid(t *testing.T) {
	//nil sender
	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg := NewMsgNewOrder(nil, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//empty sender
	orderMsg = NewMsgNewOrder([]byte{}, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//empty product
	orderMsg = NewMsgNewOrder(addr, "", BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//invalid product
	orderMsg = NewMsgNewOrder(addr, "btc"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//okb in the left
	orderMsg = NewMsgNewOrder(addr, common.NativeToken+"_btc", BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	//invalid side
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, "abc", testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//zero price
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, "0", testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//zero quantity
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, "0")
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//negative price
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, "-1", testQuantity)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	//negative quantity
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, "-1")
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// true
	orderMsg = NewMsgNewOrder(addr, "btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	// empty MultiNewOrderItems
	orderMsg.OrderItems = []OrderItem{}
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// nil MultiNewOrderItems
	orderMsg.OrderItems = nil
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// max limit
	orderMsg.OrderItems = []OrderItem{}
	item := NewOrderItem("btc_"+common.NativeToken, BuyOrder, testPrice, testQuantity)
	for i := 0; i < OrderItemLimit; i++ {
		orderMsg.OrderItems = append(orderMsg.OrderItems, item)
	}
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	// over limit
	orderMsg.OrderItems = append(orderMsg.OrderItems, item)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)
}

func TestMsgMultiCancelOrder(t *testing.T) {
	orderID := testOrderID
	addr, err := hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, err)
	orderMsg := NewMsgCancelOrder(addr, orderID)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)
	require.Equal(t, "order", orderMsg.Route())
	require.Equal(t, "cancel", orderMsg.Type())

	bytesMsg := orderMsg.GetSignBytes()
	resOrderMsg := &MsgNewOrder{}
	err = json.Unmarshal(bytesMsg, resOrderMsg)
	require.Nil(t, err)
	resAddr := orderMsg.GetSigners()[0]
	require.EqualValues(t, addr, resAddr)
}

func TestMsgMultiCancelOrderInvalid(t *testing.T) {
	//empty sender
	orderID := testOrderID
	addr, decodeErr := hex.DecodeString("")
	require.Nil(t, decodeErr)
	orderMsg := NewMsgCancelOrder(addr, orderID)
	err := orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// empty orderID
	orderID = ""
	addr, decodeErr = hex.DecodeString("1212121212121212123412121212121212121234")
	require.Nil(t, decodeErr)
	orderMsg = NewMsgCancelOrder(addr, orderID)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// true
	orderID = testOrderID
	orderMsg = NewMsgCancelOrder(addr, orderID)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	// empty OrderIDs
	orderMsg.OrderIDs = []string{}
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// nil OrderIDs
	orderMsg.OrderIDs = nil
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// max limit
	var orderIDItems []string
	for i := 0; i < MultiCancelOrderItemLimit; i++ {
		orderIDItems = append(orderIDItems, testOrderID+strconv.Itoa(i))
	}
	orderMsg = NewMsgCancelOrders(addr, orderIDItems)
	err = orderMsg.ValidateBasic()
	require.Nil(t, err)

	// over limit
	orderIDItems = append(orderIDItems, testOrderID)
	orderMsg = NewMsgCancelOrders(addr, orderIDItems)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)

	// Duplicated Id
	orderIDItems = []string{}
	orderIDItems = append(orderIDItems, testOrderID)
	orderIDItems = append(orderIDItems, testOrderID)
	orderMsg = NewMsgCancelOrders(addr, orderIDItems)
	err = orderMsg.ValidateBasic()
	require.NotNil(t, err)
}

func TestHasDuplicatedID(t *testing.T) {
	ids1 := []string{"1", "2", "3", "4", "5"}
	result1 := hasDuplicatedID(ids1)
	require.EqualValues(t, false, result1)

	ids2 := []string{"1", "3", "3", "4", "4"}
	result2 := hasDuplicatedID(ids2)
	require.EqualValues(t, true, result2)
}
