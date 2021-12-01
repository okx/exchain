package iavl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetNodeFromPool(t *testing.T) {
	n := GetNodeFromPool()
	require.Nil(t, n)
	node := NewNode([]byte("node"), []byte("this is node"), 1)
	SetNodeToPool(node)
	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.Nil(t, n)
}

func TestSetNodeToPool(t *testing.T) {
	node := NewNode([]byte("node"), []byte("this is node"), 1)
	SetNodeToPool(node)
	n := GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)

	SetNodeToPool(node)
	SetNodeToPool(node)

	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.Nil(t, n)
}

func TestNode_Reset(t *testing.T) {
	leafNode1 := NewNode([]byte("leafNode1"), []byte("this is leafNode1"), 1)
	leafNode2 := NewNode([]byte("leafNode2"), []byte("this is leafNode2"), 1)
	leafNode1.Reset(leafNode2.key, leafNode2.value, nil, nil, nil, leafNode2.version, leafNode2.size, nil, nil, leafNode2.height, false, false)

	require.Equal(t, leafNode2, leafNode1)
}
