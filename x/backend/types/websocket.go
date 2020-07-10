package types



type IWebsocket interface {
	GetChannelInfo() (channel, filter string, err error)
	FormatResult() interface{}
}