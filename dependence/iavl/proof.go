package iavl

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	amino "github.com/tendermint/go-amino"
	cmn "github.com/tendermint/iavl/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	// ErrInvalidProof is returned by Verify when a proof cannot be validated.
	ErrInvalidProof = fmt.Errorf("invalid proof")

	// ErrInvalidInputs is returned when the inputs passed to the function are invalid.
	ErrInvalidInputs = fmt.Errorf("invalid inputs")

	// ErrInvalidRoot is returned when the root passed in does not match the proof's.
	ErrInvalidRoot = fmt.Errorf("invalid root")
)

//----------------------------------------

type ProofInnerNode struct {
	Height  int8   `json:"height"`
	Size    int64  `json:"size"`
	Version int64  `json:"version"`
	Left    []byte `json:"left"`
	Right   []byte `json:"right"`
}

func (pin ProofInnerNode) String() string {
	return pin.stringIndented("")
}

func (pin ProofInnerNode) stringIndented(indent string) string {
	return fmt.Sprintf(`ProofInnerNode{
%s  Height:  %v
%s  Size:    %v
%s  Version: %v
%s  Left:    %X
%s  Right:   %X
%s}`,
		indent, pin.Height,
		indent, pin.Size,
		indent, pin.Version,
		indent, pin.Left,
		indent, pin.Right,
		indent)
}

func (pin ProofInnerNode) Hash(childHash []byte) []byte {
	hasher := tmhash.New()
	buf := new(bytes.Buffer)

	err := amino.EncodeInt8(buf, pin.Height)
	if err == nil {
		err = amino.EncodeVarint(buf, pin.Size)
	}
	if err == nil {
		err = amino.EncodeVarint(buf, pin.Version)
	}

	if len(pin.Left) == 0 {
		if err == nil {
			err = amino.EncodeByteSlice(buf, childHash)
		}
		if err == nil {
			err = amino.EncodeByteSlice(buf, pin.Right)
		}
	} else {
		if err == nil {
			err = amino.EncodeByteSlice(buf, pin.Left)
		}
		if err == nil {
			err = amino.EncodeByteSlice(buf, childHash)
		}
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to hash ProofInnerNode: %v", err))
	}

	_, err = hasher.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}

//----------------------------------------

type ProofLeafNode struct {
	Key       cmn.HexBytes `json:"key"`
	ValueHash cmn.HexBytes `json:"value"`
	Version   int64        `json:"version"`
}

func (pln ProofLeafNode) String() string {
	return pln.stringIndented("")
}

func (pln ProofLeafNode) stringIndented(indent string) string {
	return fmt.Sprintf(`ProofLeafNode{
%s  Key:       %v
%s  ValueHash: %X
%s  Version:   %v
%s}`,
		indent, pln.Key,
		indent, pln.ValueHash,
		indent, pln.Version,
		indent)
}

func (pln ProofLeafNode) Hash() []byte {
	hasher := tmhash.New()
	buf := new(bytes.Buffer)

	err := amino.EncodeInt8(buf, 0)
	if err == nil {
		err = amino.EncodeVarint(buf, 1)
	}
	if err == nil {
		err = amino.EncodeVarint(buf, pln.Version)
	}
	if err == nil {
		err = amino.EncodeByteSlice(buf, pln.Key)
	}
	if err == nil {
		err = amino.EncodeByteSlice(buf, pln.ValueHash)
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to hash ProofLeafNode: %v", err))
	}
	_, err = hasher.Write(buf.Bytes())
	if err != nil {
		panic(err)

	}

	return hasher.Sum(nil)
}

//----------------------------------------

// If the key does not exist, returns the path to the next leaf left of key (w/
// path), except when key is less than the least item, in which case it returns
// a path to the least item.
func (node *Node) PathToLeaf(t *ImmutableTree, key []byte) (PathToLeaf, *Node, error) {
	path := new(PathToLeaf)
	val, err := node.pathToLeaf(t, key, path)
	return *path, val, err
}

// pathToLeaf is a helper which recursively constructs the PathToLeaf.
// As an optimization the already constructed path is passed in as an argument
// and is shared among recursive calls.
func (node *Node) pathToLeaf(t *ImmutableTree, key []byte, path *PathToLeaf) (*Node, error) {
	if node.height == 0 {
		if bytes.Equal(node.key, key) {
			return node, nil
		}
		return node, errors.New("key does not exist")
	}

	// Note that we do not store the left child in the ProofInnerNode when we're going to add the
	// left node as part of the path, similarly we don't store the right child info when going down
	// the right child node. This is done as an optimization since the child info is going to be
	// already stored in the next ProofInnerNode in PathToLeaf.
	if bytes.Compare(key, node.key) < 0 {
		// left side
		pin := ProofInnerNode{
			Height:  node.height,
			Size:    node.size,
			Version: node.version,
			Left:    nil,
			Right:   node.getRightNode(t).hash,
		}
		*path = append(*path, pin)
		n, err := node.getLeftNode(t).pathToLeaf(t, key, path)
		return n, err
	}
	// right side
	pin := ProofInnerNode{
		Height:  node.height,
		Size:    node.size,
		Version: node.version,
		Left:    node.getLeftNode(t).hash,
		Right:   nil,
	}
	*path = append(*path, pin)
	n, err := node.getRightNode(t).pathToLeaf(t, key, path)
	return n, err
}
