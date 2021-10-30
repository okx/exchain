package iavl

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_aminoSize(t *testing.T) {
	node := &Node{
		key:       randBytes(10),
		value:     randBytes(10),
		version:   1,
		height:    0,
		size:      100,
		hash:      randBytes(20),
		leftHash:  randBytes(20),
		leftNode:  nil,
		rightHash: randBytes(20),
		rightNode: nil,
		persisted: false,
	}

	// leaf node
	require.Equal(t, 26, node.aminoSize())

	// non-leaf node
	node.height = 1
	require.Equal(t, 57, node.aminoSize())
}

func TestNode_validate(t *testing.T) {
	k := []byte("key")
	v := []byte("value")
	h := []byte{1, 2, 3}
	c := &Node{key: []byte("child"), value: []byte("x"), version: 1, size: 1}

	testcases := map[string]struct {
		node  *Node
		valid bool
	}{
		"nil node":               {nil, false},
		"leaf":                   {&Node{key: k, value: v, version: 1, size: 1}, true},
		"leaf with nil key":      {&Node{key: nil, value: v, version: 1, size: 1}, false},
		"leaf with empty key":    {&Node{key: []byte{}, value: v, version: 1, size: 1}, true},
		"leaf with nil value":    {&Node{key: k, value: nil, version: 1, size: 1}, false},
		"leaf with empty value":  {&Node{key: k, value: []byte{}, version: 1, size: 1}, true},
		"leaf with version 0":    {&Node{key: k, value: v, version: 0, size: 1}, false},
		"leaf with version -1":   {&Node{key: k, value: v, version: -1, size: 1}, false},
		"leaf with size 0":       {&Node{key: k, value: v, version: 1, size: 0}, false},
		"leaf with size 2":       {&Node{key: k, value: v, version: 1, size: 2}, false},
		"leaf with size -1":      {&Node{key: k, value: v, version: 1, size: -1}, false},
		"leaf with left hash":    {&Node{key: k, value: v, version: 1, size: 1, leftHash: h}, false},
		"leaf with left child":   {&Node{key: k, value: v, version: 1, size: 1, leftNode: c}, false},
		"leaf with right hash":   {&Node{key: k, value: v, version: 1, size: 1, rightNode: c}, false},
		"leaf with right child":  {&Node{key: k, value: v, version: 1, size: 1, rightNode: c}, false},
		"inner":                  {&Node{key: k, version: 1, size: 1, height: 1, leftHash: h, rightHash: h}, true},
		"inner with nil key":     {&Node{key: nil, value: v, version: 1, size: 1, height: 1, leftHash: h, rightHash: h}, false},
		"inner with value":       {&Node{key: k, value: v, version: 1, size: 1, height: 1, leftHash: h, rightHash: h}, false},
		"inner with empty value": {&Node{key: k, value: []byte{}, version: 1, size: 1, height: 1, leftHash: h, rightHash: h}, false},
		"inner with left child":  {&Node{key: k, version: 1, size: 1, height: 1, leftHash: h}, true},
		"inner with right child": {&Node{key: k, version: 1, size: 1, height: 1, rightHash: h}, true},
		"inner with no child":    {&Node{key: k, version: 1, size: 1, height: 1}, false},
		"inner with height 0":    {&Node{key: k, version: 1, size: 1, height: 0, leftHash: h, rightHash: h}, false},
	}

	for desc, tc := range testcases {
		tc := tc // appease scopelint
		t.Run(desc, func(t *testing.T) {
			err := tc.node.validate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func BenchmarkNode_aminoSize(b *testing.B) {
	node := &Node{
		key:       randBytes(25),
		value:     randBytes(100),
		version:   rand.Int63n(10000000),
		height:    1,
		size:      rand.Int63n(10000000),
		leftHash:  randBytes(20),
		rightHash: randBytes(20),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.aminoSize()
	}
}

func BenchmarkNode_WriteBytes(b *testing.B) {
	node := &Node{
		key:       randBytes(25),
		value:     randBytes(100),
		version:   rand.Int63n(10000000),
		height:    1,
		size:      rand.Int63n(10000000),
		leftHash:  randBytes(20),
		rightHash: randBytes(20),
	}
	b.ResetTimer()
	b.Run("NoPreAllocate", func(sub *testing.B) {
		sub.ReportAllocs()
		for i := 0; i < sub.N; i++ {
			var buf bytes.Buffer
			buf.Reset()
			_ = node.writeBytes(&buf)
		}
	})
	b.Run("PreAllocate", func(sub *testing.B) {
		sub.ReportAllocs()
		for i := 0; i < sub.N; i++ {
			var buf bytes.Buffer
			buf.Grow(node.aminoSize())
			_ = node.writeBytes(&buf)
		}
	})
}
