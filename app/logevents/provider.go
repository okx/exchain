package logevents


type provider struct {
	eventChan chan string
	ipaddr string
}

func newProvider() *provider {
	return &provider{
		eventChan: make(chan string, 1000),
	}
}

func (p* provider) init()  {

}

func (p* provider) AddEvent(event string)  {
	p.eventChan <- event
}

func (p* provider) eventRoutine()  {
	for event := range p.eventChan {
		p.eventHandler(event)
	}
}

func (p* provider) eventHandler(string)  {

}