package context

type TxRequest interface {
	GetData() []byte
	GetModeDetail() int32
}

type TxResponse interface {
	HandleResponse(data interface{}) interface{}
}
