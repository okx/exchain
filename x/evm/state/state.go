package state

type State struct {
	prefix []byte
	store  *stateStore
}

func NewState(prefix []byte) *State {
	return &State{prefix: prefix, store: InstanceOfStateStore()}
}

func (s State) genPrefix(tail []byte) (res []byte) {
	res = make([]byte, len(s.prefix)+len(tail))
	copy(res, s.prefix)
	copy(res[len(s.prefix):], tail)
	return
}

func (s State) Set(key, value []byte) error {
	return nil
}

func (s State) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (s *State) Delete(key []byte) error {
	return nil
}
