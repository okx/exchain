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

func (ndb *nodeDB) updateBranchConcurrency(node *Node, savedNodes map[string]*Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	nodeCh := make(chan *Node, 1024)
	leftHashCh := make(chan []byte, 1)
	rightHashCh := make(chan []byte, 1)
	wg := &sync.WaitGroup{}

	var needNilNodeNum = 0
	if node.leftNode != nil {
		needNilNodeNum += 1
		go updateBranchRoutine(node.leftNode, nodeCh, leftHashCh)
	}
	if node.rightNode != nil {
		needNilNodeNum += 1
		go updateBranchRoutine(node.rightNode, nodeCh, rightHashCh)
	}

	if needNilNodeNum > 0 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, needNilNodeNum int, savedNodes map[string]*Node, ndb *nodeDB, nodeCh <-chan *Node) {
			getNodeNil := 0
			for n := range nodeCh {
				if n == nil {
					getNodeNil += 1
					if getNodeNil == needNilNodeNum {
						wg.Done()
						return
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
		}(wg, needNilNodeNum, savedNodes, ndb, nodeCh)
	}

	if node.leftNode != nil {
		node.leftHash = <-leftHashCh
		close(leftHashCh)
	}

	if node.rightNode != nil {
		node.rightHash = <-rightHashCh
		close(rightHashCh)
	}
	node._hash()

	wg.Wait()
	close(nodeCh)
	ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil

	// TODO: handle magic number
	if savedNodes != nil {
		savedNodes[string(node.hash)] = node
	}

	return node.hash
}

func updateBranchRoutine(node *Node, saveNodesCh chan<- *Node, result chan<- []byte) {
	if node.persisted || node.prePersisted {
		saveNodesCh <- nil
		result <- node.hash
		return
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

	result <- node.hash
	return
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
