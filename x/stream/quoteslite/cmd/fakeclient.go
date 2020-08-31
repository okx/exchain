package main

import (
	"fmt"

	okex "github.com/okex/okchain/x/stream/quoteslite/okwebsocket"
)

func main() {
	fmt.Println("start fakeclient!")
	agent := okex.OKWSAgent{}
	config := okex.GetDefaultConfig()
	config.WSEndpoint = "ws://127.0.0.1:26661/"
	config.IsPrint = true
	agent.SetConfig(config)

	// Step1: Start agent.
	agent.Start(config)

	address := "okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya"
	agent.DexLogin(address)
	product := "eos-774_okt"
	// Step2: Subscribe Channel
	// okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya
	argsMap := map[string]string{
		// dex_spot/account:okt
		okex.DexSpotAccount: "okt",
		// dex_spot/order:xxb_okt
		okex.DexSpotOrder: product,
		// dex_spot/matches:xxb_okt
		okex.DexSpotMatch: product,
		// dex_spot/optimized_depth:xxb_okt
		okex.DexSpotDepthBook: product,

		okex.DexSpotAllTicker3s: "",
		okex.DexSpotTicker:      product,
	}
	for channel, filter := range argsMap {
		agent.Subscribe(channel, filter, okex.DefaultDataCallBack)
	}

	// Step3 stop
	select {}
	agent.Stop()
	fmt.Println("stop fakeclient!")
}
