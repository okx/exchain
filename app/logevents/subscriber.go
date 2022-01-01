package logevents

type Subscriber interface {
	Init(urls string, topic string)
	Run()
}

func NewSubscriber() Subscriber {
	return &subscriber{}
}

type subscriber struct {
	
}

func (s* subscriber) Init(urls string, topic string)  {

}
func (s* subscriber) Run()  {

}




