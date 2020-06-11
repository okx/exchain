package okex

/*
 OKEX ws api websocket test & sample
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"testing"
	"time"
)

func TestOKWSAgent_AllInOne(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/ticker successfully.
	agent.Subscribe(CHNL_SWAP_TICKER, "BTC-USD-SWAP", nil)
	agent.Subscribe(CHNL_SWAP_TICKER, "BTC-USD-SWAP", DefaultDataCallBack)

	// Step2.1: Subscribe private Channel swap/position before login, so it would be a fail.
	agent.Subscribe(CHNL_SWAP_POSITION, "BTC-USD-SWAP", DefaultDataCallBack)

	// Step3: Wait for the ws server's pushed table responses.
	time.Sleep(60 * time.Second)

	// Step4. Unsubscribe public Channel swap/ticker
	agent.UnSubscribe(CHNL_SWAP_TICKER, "BTC-USD-SWAP")
	time.Sleep(1 * time.Second)

	// Step5. Login
	agent.Login(config.ApiKey, config.Passphrase)
	time.Sleep(1 * time.Second)

	// Step6. Subscribe private Channel swap/position after login, so it would be a success.
	agent.Subscribe(CHNL_SWAP_POSITION, "BTC-USD-SWAP", DefaultDataCallBack)
	time.Sleep(120 * time.Second)

	// Step7. Stop all the go routine run in background.
	agent.Stop()
	time.Sleep(1 * time.Second)
}

func TestOKWSAgent_Depths(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/depths successfully.
	agent.Subscribe(CHNL_SWAP_DEPTH, "BTC-USD-SWAP", DefaultDataCallBack)

	// Step3: Client receive depths from websocket server.
	// Step3.0: Receive partial depths
	// Step3.1: Receive update depths (It may take a very long time to see Update Event.)

	time.Sleep(60 * time.Second)

	// Step4. Stop all the go routine run in background.
	agent.Stop()
	time.Sleep(1 * time.Second)
}

func TestOKWSAgent_mergeDepths(t *testing.T) {
	oldDepths := [][4]interface{}{
		{"5088.59", "34000", 0, 1},
		{"7200", "1", 0, 1},
		{"7300", "1", 0, 1},
	}

	// Case1.
	newDepths1 := [][4]interface{}{
		{"5088.59", "32000", 0, 1},
	}
	expectedMerged1 := [][4]interface{}{
		{"5088.59", "32000", 0, 1},
		{"7200", "1", 0, 1},
		{"7300", "1", 0, 1},
	}

	m1, e1 := mergeDepths(oldDepths, newDepths1)
	assert.True(t, e1 == nil)
	assert.True(t, len(*m1) == len(expectedMerged1) && (*m1)[0][1] == expectedMerged1[0][1] && (*m1)[0][1] == "32000")

	// Case2.
	newDepths2 := [][4]interface{}{
		{"7200", "0", 0, 1},
	}
	expectedMerged2 := [][4]interface{}{
		{"5088.59", "34000", 0, 1},
		{"7300", "1", 0, 1},
	}
	m2, e2 := mergeDepths(oldDepths, newDepths2)
	assert.True(t, e2 == nil)
	assert.True(t, len(*m2) == len(expectedMerged2) && (*m2)[0][1] == expectedMerged2[0][1] && (*m2)[0][1] == "34000")

	// Case3.
	newDepths3 := [][4]interface{}{
		{"5000", "1", 0, 1},
		{"7400", "1", 0, 1},
	}
	expectedMerged3 := [][4]interface{}{
		{"5000", "1", 0, 1},
		{"5088.59", "34000", 0, 1},
		{"7200", "1", 0, 1},
		{"7300", "1", 0, 1},
		{"7400", "1", 0, 1},
	}
	m3, e3 := mergeDepths(oldDepths, newDepths3)
	assert.True(t, e3 == nil)
	assert.True(t, len(*m3) == len(expectedMerged3) && (*m3)[0][1] == expectedMerged3[0][1] && (*m3)[0][1] == "1")

}

func TestOKWSAgent_calCrc32(t *testing.T) {

	askDepths := [][4]interface{}{
		{"5088.59", "34000", 0, 1},
		{"7200", "1", 0, 1},
		{"7300", "1", 0, 1},
	}

	bidDepths1 := [][4]interface{}{
		{"3850", "1", 0, 1},
		{"3800", "1", 0, 1},
		{"3500", "1", 0, 1},
		{"3000", "1", 0, 1},
	}

	crcBuf1, caled1 := calCrc32(&askDepths, &bidDepths1)
	assert.True(t, caled1 != 0 && crcBuf1.String() == "3850:1:3800:1:3500:1:3000:1:5088.59:34000:7200:1:7300:1")

	bidDepths2 := [][4]interface{}{
		{"3850", "1", 0, 1},
		{"3800", "1", 0, 1},
		{"3500", "1", 0, 1},
	}

	crcBuf2, caled2 := calCrc32(&askDepths, &bidDepths2)
	assert.True(t, caled2 != 0 && crcBuf2.String() == "3850:1:5088.59:34000:3800:1:7200:1:3500:1:7300:1")
}

func TestArray(t *testing.T) {

	t1 := [4]int{1, 2, 3, 4}
	t2 := [][4]int{
		{1, 2, 3, 4},
	}
	t3 := [4][]int{
		{1, 2, 3, 4},
	}

	r1, _ := Struct2JsonString(t1)
	r2, _ := Struct2JsonString(t2)
	r3, _ := Struct2JsonString(t3)

	println(len(t1), r1)
	println(len(t2), r2)
	println(len(t3), r3)

	fmt.Printf("%+v\n", t1[0:len(t1)-1])
}

func TestCrc32(t *testing.T) {
	str1 := "3366.1:7:3366.8:9:3366:6:3368:8"
	r := crc32.ChecksumIEEE([]byte(str1))
	println(r)
	assert.True(t, int32(r) == -1881014294)

	str2 := "3366.1:7:3366.8:9:3368:8:3372:8"
	r = crc32.ChecksumIEEE([]byte(str2))
	println(r)
	assert.True(t, int32(r) == 831078360)
}

func TestFmtSprintf(t *testing.T) {
	a := [][]interface{}{
		{"199", "10"},
		{199.0, 10.0},
	}

	for _, v := range a {
		s1 := fmt.Sprintf("%v:%v", v[0], v[1])
		s2 := fmt.Sprintf("%s:%s", v[0], v[1])
		println(s1)
		println(s2)
		assert.True(t, s1 != "" && s2 != "")
	}

}

func TestOKWSAgent_Futures_AllInOne(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()
	publicChannels := []string{
		CHNL_FUTURES_CANDLE60S,
		CHNL_FUTURES_CANDLE180S,
		CHNL_FUTURES_CANDLE300S,
		CHNL_FUTURES_CANDLE900S,
		CHNL_FUTURES_CANDLE1800S,
		CHNL_FUTURES_CANDLE3600S,
		CHNL_FUTURES_CANDLE7200S,
		CHNL_FUTURES_CANDLE14400S,
		CHNL_FUTURES_DEPTH,
		CHNL_FUTURES_DEPTH5,
		CHNL_FUTURES_ESTIMATED_PRICE,
		CHNL_FUTURES_MARK_PRICE,
		CHNL_FUTURES_PRICE_RANGE,
		CHNL_FUTURES_TICKER,
		CHNL_FUTURES_TRADE,
	}

	privateChannels := []string{
		CHNL_FUTURES_ACCOUNT,
		CHNL_FUTURES_ORDER,
		CHNL_FUTURES_POSITION,
	}
	filter := "BTC-USD-170310"

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Login
	agent.Login(config.ApiKey, config.Passphrase)
	time.Sleep(1 * time.Second)

	// Step3: Subscribe privateChannels
	for _, c := range privateChannels {
		agent.Subscribe(c, filter, DefaultDataCallBack)
	}

	// Step4: Subscribe publicChannels
	for _, c := range publicChannels {
		agent.Subscribe(c, filter, DefaultDataCallBack)
	}
	time.Sleep(time.Second * 2)

	// Step5: unsubscribe privateChannels
	for _, c := range privateChannels {
		agent.UnSubscribe(c, filter)
	}

	agent.Stop()
}

func TestOKWSAgent_Spots_AllInOne(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()
	publicChannels := []string{
		CHNL_SPOT_CANDLE60S,
		CHNL_SPOT_CANDLE180S,
		CHNL_SPOT_CANDLE300S,
		CHNL_SPOT_CANDLE900S,
		CHNL_SPOT_CANDLE1800S,
		CHNL_SPOT_CANDLE3600S,
		CHNL_SPOT_CANDLE7200S,
		CHNL_SPOT_CANDLE14400S,
		CHNL_SPOT_DEPTH,
		CHNL_SPOT_DEPTH5,
		CHNL_SPOT_TICKER,
		CHNL_SPOT_TRADE,
	}

	privateChannels := []string{
		CHNL_SPOT_ACCOUNT,
		CHNL_SPOT_MARGIN_ACCOUNT,
		CHNL_SPOT_ORDER,
	}
	filter := "ETH-USDT"

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Login
	agent.Login(config.ApiKey, config.Passphrase)
	time.Sleep(1 * time.Second)

	// Step3: Subscribe privateChannels
	for _, c := range privateChannels {
		agent.Subscribe(c, filter, DefaultDataCallBack)
	}

	// Step4: Subscribe publicChannels
	for _, c := range publicChannels {
		agent.Subscribe(c, filter, DefaultDataCallBack)
	}
	time.Sleep(time.Second * 2)

	// Step5: unsubscribe privateChannels
	for _, c := range privateChannels {
		agent.UnSubscribe(c, filter)
	}

	agent.Stop()
}
