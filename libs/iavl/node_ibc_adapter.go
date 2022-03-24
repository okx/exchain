package iavl

import (
	"bytes"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

type traversal struct {
	tree         *ImmutableTree
	start, end   []byte        // iteration domain
	ascending    bool          // ascending traversal
	inclusive    bool          // end key inclusiveness
	post         bool          // postorder traversal
	delayedNodes *delayedNodes // delayed nodes to be traversed
}

// delayedNode represents the delayed iteration on the nodes.
// When delayed is set to true, the delayedNode should be expanded, and their
// children should be traversed. When delayed is set to false, the delayedNode is
// already have expanded, and it could be immediately returned.
type delayedNode struct {
	node    *Node
	delayed bool
}

type delayedNodes []delayedNode

func (node *Node) newTraversal(tree *ImmutableTree, start, end []byte, ascending bool, inclusive bool, post bool) *traversal {
	return &traversal{
		tree:         tree,
		start:        start,
		end:          end,
		ascending:    ascending,
		inclusive:    inclusive,
		post:         post,
		delayedNodes: &delayedNodes{{node, true}}, // set initial traverse to the node
	}
}
func (t *ImmutableTree) GetWithProof2(key []byte) (value []byte, proof *RangeProof, err error) {
	proof, _, values, err := t.getRangeProof2(key, cpIncr(key), 2)
	if err != nil {
		return nil, nil, errors.Wrap(err, "constructing range proof")
	}
	if len(values) > 0 && bytes.Equal(proof.Leaves[0].Key, key) {
		return values[0], proof, nil
	}
	return nil, proof, nil
}

func (nodes *delayedNodes) pop() (*Node, bool) {
	node := (*nodes)[len(*nodes)-1]
	*nodes = (*nodes)[:len(*nodes)-1]
	return node.node, node.delayed
}

func (nodes *delayedNodes) push(node *Node, delayed bool) {
	*nodes = append(*nodes, delayedNode{node, delayed})
}

func (nodes *delayedNodes) length() int {
	return len(*nodes)
}

func (node *Node) traverseInRange2(tree *ImmutableTree, start, end []byte, ascending bool, inclusive bool, post bool, cb func(*Node) bool) bool {
	stop := false
	t := node.newTraversal(tree, start, end, ascending, inclusive, post)
	for node2 := t.next(); node2 != nil; node2 = t.next() {
		stop = cb(node2)
		if stop {
			return stop
		}
	}
	return stop
}

// Only used in testing...
func (node *Node) lmd(t *ImmutableTree) *Node {
	if node.isLeaf() {
		return node
	}
	return node.getLeftNode(t).lmd(t)
}

func (t *traversal) next() *Node {
	// End of traversal.
	if t.delayedNodes.length() == 0 {
		return nil
	}

	node, delayed := t.delayedNodes.pop()

	// Already expanded, immediately return.
	if !delayed || node == nil {
		return node
	}

	afterStart := t.start == nil || bytes.Compare(t.start, node.key) < 0
	startOrAfter := afterStart || bytes.Equal(t.start, node.key)
	beforeEnd := t.end == nil || bytes.Compare(node.key, t.end) < 0
	if t.inclusive {
		beforeEnd = beforeEnd || bytes.Equal(node.key, t.end)
	}

	// case of postorder. A-1 and B-1
	// Recursively process left sub-tree, then right-subtree, then node itself.
	if t.post && (!node.isLeaf() || (startOrAfter && beforeEnd)) {
		t.delayedNodes.push(node, false)
	}

	// case of branch node, traversing children. A-2.
	if !node.isLeaf() {
		// if node is a branch node and the order is ascending,
		// We traverse through the left subtree, then the right subtree.
		if t.ascending {
			if beforeEnd {
				// push the delayed traversal for the right nodes,
				t.delayedNodes.push(node.getRightNode(t.tree), true)
			}
			if afterStart {
				// push the delayed traversal for the left nodes,
				t.delayedNodes.push(node.getLeftNode(t.tree), true)
			}
		} else {
			// if node is a branch node and the order is not ascending
			// We traverse through the right subtree, then the left subtree.
			if afterStart {
				// push the delayed traversal for the left nodes,
				t.delayedNodes.push(node.getLeftNode(t.tree), true)
			}
			if beforeEnd {
				// push the delayed traversal for the right nodes,
				t.delayedNodes.push(node.getRightNode(t.tree), true)
			}
		}
	}

	// case of preorder traversal. A-3 and B-2.
	// Process root then (recursively) processing left child, then process right child
	if !t.post && (!node.isLeaf() || (startOrAfter && beforeEnd)) {
		return node
	}

	// Keep traversing and expanding the remaning delayed nodes. A-4.
	return t.next()
}
