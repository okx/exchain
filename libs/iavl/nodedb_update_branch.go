package iavl

import "sync"

func (ndb *nodeDB) updateBranch(node *Node, savedNodes map[string]*Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = ndb.updateBranch(node.leftNode, savedNodes)
	}
	if node.rightNode != nil {
		node.rightHash = ndb.updateBranch(node.rightNode, savedNodes)
	}

	node._hash()
	ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil

	// TODO: handle magic number
	savedNodes[string(node.hash)] = node

	return node.hash
}

func (ndb *nodeDB) updateBranchConcurrency(node *Node, savedNodes map[string]*Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	nodeCh := make(chan *Node, 1024)
	wg := &sync.WaitGroup{}

	var needNilNodeNum = 0
	if node.leftNode != nil {
		needNilNodeNum += 1
		wg.Add(1)
		go func(node *Node, wg *sync.WaitGroup, nodeCh chan *Node) {
			node.leftHash = updateBranchRoutine(node.leftNode, nodeCh)
			wg.Done()
		}(node, wg, nodeCh)
	}
	if node.rightNode != nil {
		needNilNodeNum += 1
		wg.Add(1)
		go func(node *Node, wg *sync.WaitGroup, nodeCh chan *Node) {
			node.rightHash = updateBranchRoutine(node.rightNode, nodeCh)
			wg.Done()
		}(node, wg, nodeCh)
	}

	if needNilNodeNum > 0 {
		getNodeNil := 0
		for n := range nodeCh {
			if n == nil {
				getNodeNil += 1
				if getNodeNil == needNilNodeNum {
					break
				}
			} else {
				ndb.saveNodeToPrePersistCache(n)
				n.leftNode = nil
				n.rightNode = nil
				if savedNodes != nil {
					savedNodes[string(n.hash)] = n
				}
			}
		}
	}

	close(nodeCh)
	wg.Wait()

	node._hash()

	ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil

	// TODO: handle magic number
	if savedNodes != nil {
		savedNodes[string(node.hash)] = node
	}

	return node.hash
}

func updateBranchRoutine(node *Node, saveNodesCh chan<- *Node) []byte {
	if node.persisted || node.prePersisted {
		saveNodesCh <- nil
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = updateBranchAndSaveNodeToChan(node.leftNode, saveNodesCh)
	}
	if node.rightNode != nil {
		node.rightHash = updateBranchAndSaveNodeToChan(node.rightNode, saveNodesCh)
	}

	node._hash()

	saveNodesCh <- node
	saveNodesCh <- nil

	return node.hash
}

func updateBranchAndSaveNodeToChan(node *Node, saveNodesCh chan<- *Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = updateBranchAndSaveNodeToChan(node.leftNode, saveNodesCh)
	}
	if node.rightNode != nil {
		node.rightHash = updateBranchAndSaveNodeToChan(node.rightNode, saveNodesCh)
	}

	node._hash()

	saveNodesCh <- node

	return node.hash
}

func (ndb *nodeDB) updateBranchForFastNode(fnc *fastNodeChanges) {
	ndb.mtx.Lock()
	ndb.prePersistFastNode.mergeLater(fnc)
	ndb.mtx.Unlock()
}

func (ndb *nodeDB) updateBranchMoreConcurrency(node *Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	wg := &sync.WaitGroup{}

	if node.leftNode != nil {
		wg.Add(1)
		go func(node *Node, wg *sync.WaitGroup) {
			node.leftHash = ndb.updateBranchConcurrency(node.leftNode, nil)
			wg.Done()
		}(node, wg)
	}
	if node.rightNode != nil {
		wg.Add(1)
		go func(node *Node, wg *sync.WaitGroup) {
			node.rightHash = ndb.updateBranchConcurrency(node.rightNode, nil)
			wg.Done()
		}(node, wg)
	}

	wg.Wait()

	node._hash()
	ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil

	return node.hash
}
