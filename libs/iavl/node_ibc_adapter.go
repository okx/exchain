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