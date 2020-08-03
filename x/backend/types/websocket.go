package types

import (
	"fmt"
)

const WebsocketChanCapacity = 2048

type IWebsocket interface {
	GetChannelInfo() (channel, filter string, err error)
	GetFullChannel() string
	FormatResult() interface{}
	GetTimestamp() int64
}

// FakeWSEvent won't result in sending out tendermint event
type FakeWSEvent struct {
	channel string
	filter  string
	ts      int64
}

func NewFakeWSEvent(channel, filter string, ts int64) *FakeWSEvent {
	return &FakeWSEvent{
		channel: channel,
		filter:  filter,
		ts:      ts,
	}
}

func (f *FakeWSEvent) GetChannelInfo() (channel, filter string, err error) {
	return f.channel, f.filter, nil
}

func (f *FakeWSEvent) GetFullChannel() string {
	if f.filter == "" {
		return f.channel
	} else {
		return f.channel + ":" + f.filter
	}
}

func (f *FakeWSEvent) FormatResult() interface{} {
	return nil
}

func (f *FakeWSEvent) GetTimestamp() int64 {
	return f.ts
}

type MergedTickersEvent struct {
	freq    int
	ts      int64
	tickers []interface{}
}

func (m *MergedTickersEvent) GetFullChannel() string {
	return fmt.Sprintf("dex_spot/all_ticker_%ds", m.freq)
}

func (m *MergedTickersEvent) GetChannelInfo() (channel, filter string, err error) {
	return m.GetFullChannel(), "", nil
}

func NewMergedTickersEvent(ts int64, freq int, tickers []interface{}) *MergedTickersEvent {
	return &MergedTickersEvent{
		freq:    freq,
		ts:      ts,
		tickers: tickers,
	}
}

func (m *MergedTickersEvent) FormatResult() interface{} {
	r := []interface{}{}
	for _, ticker := range m.tickers {
		origin := ticker.(map[string]string)
		result := map[string]string{
			"id": origin["product"],
			"o":  origin["open"],
			"h":  origin["high"],
			"l":  origin["low"],
			"v":  origin["volume"],
			"p":  origin["price"],
		}

		r = append(r, result)
	}

	return r
}

func (m *MergedTickersEvent) GetTimestamp() int64 {
	return m.ts
}
