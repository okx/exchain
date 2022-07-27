package iavl

type fastNodeChanges struct {
	additions map[string]*FastNode
	removals  map[string]interface{}
}

func newFastNodeChanges() *fastNodeChanges {
	return &fastNodeChanges{
		additions: make(map[string]*FastNode),
		removals:  make(map[string]interface{}),
	}
}

func (fnc *fastNodeChanges) Get(key []byte) (*FastNode, bool) {
	if node, ok := fnc.additions[string(key)]; ok {
		return node, true
	}
	if _, ok := fnc.removals[string(key)]; ok {
		return nil, true
	}

	return nil, false
}

func (fnc *fastNodeChanges) Add(key string, fastNode *FastNode) {
	fnc.additions[key] = fastNode
	delete(fnc.removals, key)
}

func (fnc *fastNodeChanges) Remove(key string, value interface{}) {
	fnc.removals[key] = value
	delete(fnc.additions, key)
}

func (fnc *fastNodeChanges) Reset() {
	for k := range fnc.additions {
		delete(fnc.additions, k)
	}
	for k := range fnc.removals {
		delete(fnc.removals, k)
	}
}
