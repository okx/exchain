package iavl

import (
	"bytes"
	"crypto/sha256"
)

func (t *ImmutableTree) getRangeProof2(keyStart, keyEnd []byte, limit int) (proof *RangeProof, keys, values [][]byte, err error) {
	if keyStart != nil && keyEnd != nil && bytes.Compare(keyStart, keyEnd) >= 0 {
		panic("if keyStart and keyEnd are present, need keyStart < keyEnd.")
	}
	if limit < 0 {
		panic("limit must be greater or equal to 0 -- 0 means no limit")
	}
	if t.root == nil {
		return nil, nil, nil, nil
	}
	t.root.hashWithCount() // Ensure that all hashes are calculated.

	// Get the first key/value pair proof, which provides us with the left key.
	path, left, err := t.root.PathToLeaf(t, keyStart)
	if err != nil {
		// Key doesn't exist, but instead we got the prev leaf (or the
		// first or last leaf), which provides proof of absence).
		err = nil
	}
	startOK := keyStart == nil || bytes.Compare(keyStart, left.key) <= 0
	endOK := keyEnd == nil || bytes.Compare(left.key, keyEnd) < 0
	// If left.key is in range, add it to key/values.
	if startOK && endOK {
		keys = append(keys, left.key) // == keyStart
		values = append(values, left.value)
	}

	h := sha256.Sum256(left.value)
	var leaves = []ProofLeafNode{
		{
			Key:       left.key,
			ValueHash: h[:],
			Version:   left.version,
		},
	}

	// 1: Special case if limit is 1.
	// 2: Special case if keyEnd is left.key+1.
	_stop := false
	if limit == 1 {
		_stop = true // case 1
	} else if keyEnd != nil && bytes.Compare(cpIncr(left.key), keyEnd) >= 0 {
		_stop = true // case 2
	}
	if _stop {
		return &RangeProof{
			LeftPath: path,
			Leaves:   leaves,
		}, keys, values, nil
	}

	// Get the key after left.key to iterate from.
	afterLeft := cpIncr(left.key)

	// Traverse starting from afterLeft, until keyEnd or the next leaf
	// after keyEnd.
	var allPathToLeafs = []PathToLeaf(nil)
	var currentPathToLeaf = PathToLeaf(nil)
	var leafCount = 1 // from left above.
	var pathCount = 0

	t.root.traverseInRange2(t, afterLeft, nil, true, false, false,
		func(node *Node) (stop bool) {

			// Track when we diverge from path, or when we've exhausted path,
			// since the first allPathToLeafs shouldn't include it.
			if pathCount != -1 {
				if len(path) <= pathCount {
					// We're done with path counting.
					pathCount = -1
				} else {
					pn := path[pathCount]
					if pn.Height != node.height ||
						pn.Left != nil && !bytes.Equal(pn.Left, node.leftHash) ||
						pn.Right != nil && !bytes.Equal(pn.Right, node.rightHash) {

						// We've diverged, so start appending to allPathToLeaf.
						pathCount = -1
					} else {
						pathCount++
					}
				}
			}

			if node.height == 0 { // Leaf node
				// Append all paths that we tracked so far to get to this leaf node.
				allPathToLeafs = append(allPathToLeafs, currentPathToLeaf)
				// Start a new one to track as we traverse the tree.
				currentPathToLeaf = PathToLeaf(nil)

				h := sha256.Sum256(node.value)
				leaves = append(leaves, ProofLeafNode{
					Key:       node.key,
					ValueHash: h[:],
					Version:   node.version,
				})

				leafCount++

				// Maybe terminate because we found enough leaves.
				if limit > 0 && limit <= leafCount {
					return true
				}

				// Terminate if we've found keyEnd or after.
				if keyEnd != nil && bytes.Compare(node.key, keyEnd) >= 0 {
					return true
				}

				// Value is in range, append to keys and values.
				keys = append(keys, node.key)
				values = append(values, node.value)

				// Terminate if we've found keyEnd-1 or after.
				// We don't want to fetch any leaves for it.
				if keyEnd != nil && bytes.Compare(cpIncr(node.key), keyEnd) >= 0 {
					return true
				}

			} else if pathCount < 0 { // Inner node.
				// Only store if the node is not stored in currentPathToLeaf already. We track if we are
				// still going through PathToLeaf using pathCount. When pathCount goes to -1, we
				// start storing the other paths we took to get to the leaf nodes. Also we skip
				// storing the left node, since we are traversing the tree starting from the left
				// and don't need to store unnecessary info as we only need to go down the right
				// path.
				currentPathToLeaf = append(currentPathToLeaf, ProofInnerNode{
					Height:  node.height,
					Size:    node.size,
					Version: node.version,
					Left:    nil,
					Right:   node.rightHash,
				})
			}
			return false
		},
	)

	return &RangeProof{
		LeftPath:   path,
		InnerNodes: allPathToLeafs,
		Leaves:     leaves,
	}, keys, values, nil
}
