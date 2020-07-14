package types

type IWebsocket interface {
	GetChannelInfo() (channel, filter string, err error)
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

func (f *FakeWSEvent) FormatResult() interface{} {
	return nil
}

func (f *FakeWSEvent) GetTimestamp() int64 {
	return f.ts
}
