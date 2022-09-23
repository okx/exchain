package iavl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type ExceptedKey struct {
	key      string
	fastNode interface{}
}

func TestUnsavedFastIterator(t *testing.T) {
	start := []byte("test01")
	end := []byte("test100")

	cases := []struct {
		name                     string
		unsavedfastNodeAdditions map[string]*FastNode
		unsavedfastNodeRemovals  map[string]interface{}
		expected                 []ExceptedKey
		init                     func() *nodeDB
	}{
		{
			"interator with additions and ndb success",
			map[string]*FastNode{
				"test01": NewFastNode([]byte("test01"), []byte("test01valuenew"), 11),
				"test02": NewFastNode([]byte("test02"), []byte("test02valuenew"), 11),
			},
			make(map[string]interface{}),
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valuenew"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valuenew"), 11)},
				ExceptedKey{"test03", NewFastNode([]byte("test03"), []byte("test03valueold"), 11)},
				ExceptedKey{"test04", NewFastNode([]byte("test04"), []byte("test04valueold"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.set([]byte("test01"), []byte("test01valueold"))
				tree.set([]byte("test02"), []byte("test02valueold"))
				tree.set([]byte("test03"), []byte("test03valueold"))
				tree.set([]byte("test04"), []byte("test04valueold"))
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
		{
			"interator with removals and ndb success",
			map[string]*FastNode{},
			map[string]interface{}{
				"test03": NewFastNode([]byte("test03"), []byte("test03value"), 11),
				"test04": NewFastNode([]byte("test04"), []byte("test04value"), 11),
			},
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valueold"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valueold"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.set([]byte("test01"), []byte("test01valueold"))
				tree.set([]byte("test02"), []byte("test02valueold"))
				tree.set([]byte("test03"), []byte("test03valueold"))
				tree.set([]byte("test04"), []byte("test04valueold"))
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
		{
			"interator fastnode and ndb interator with new fastnode additions same with ndb without removals success",
			map[string]*FastNode{
				"test01": NewFastNode([]byte("test01"), []byte("test01valuenew"), 11),
				"test02": NewFastNode([]byte("test02"), []byte("test02valuenew"), 11),
			},
			make(map[string]interface{}),
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valuenew"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valuenew"), 11)},
				ExceptedKey{"test03", NewFastNode([]byte("test03"), []byte("test03valueold"), 11)},
				ExceptedKey{"test04", NewFastNode([]byte("test04"), []byte("test04valueold"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.set([]byte("test01"), []byte("test01valueold"))
				tree.set([]byte("test02"), []byte("test02valueold"))
				tree.set([]byte("test03"), []byte("test03valueold"))
				tree.set([]byte("test04"), []byte("test04valueold"))
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
		{
			"interator fastnode and ndb with deletions success",
			map[string]*FastNode{
				"test01": NewFastNode([]byte("test01"), []byte("test01valuenew"), 11),
				"test02": NewFastNode([]byte("test02"), []byte("test02valuenew"), 11),
			},
			map[string]interface{}{
				"test03": NewFastNode([]byte("test03"), []byte("test03value"), 11),
				"test04": NewFastNode([]byte("test04"), []byte("test04value"), 11),
			},
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valuenew"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valuenew"), 11)},
				ExceptedKey{"test05", NewFastNode([]byte("test05"), []byte("test05valueold"), 11)},
				ExceptedKey{"test06", NewFastNode([]byte("test06"), []byte("test06valueold"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.set([]byte("test01"), []byte("test01valueold"))
				tree.set([]byte("test02"), []byte("test02valueold"))
				tree.set([]byte("test03"), []byte("test03valueold"))
				tree.set([]byte("test04"), []byte("test04valueold"))
				tree.set([]byte("test05"), []byte("test05valueold"))
				tree.set([]byte("test06"), []byte("test06valueold"))
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
		{
			"interator fastnode without ndb success",
			map[string]*FastNode{
				"test01": NewFastNode([]byte("test01"), []byte("test01valuenew"), 11),
				"test02": NewFastNode([]byte("test02"), []byte("test02valuenew"), 11),
			},
			map[string]interface{}{},
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valuenew"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valuenew"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
		{
			"interator only with ndb success",
			map[string]*FastNode{},
			map[string]interface{}{},
			[]ExceptedKey{
				ExceptedKey{"test01", NewFastNode([]byte("test01"), []byte("test01valueold"), 11)},
				ExceptedKey{"test02", NewFastNode([]byte("test02"), []byte("test02valueold"), 11)},
			},
			func() *nodeDB {
				tree, err := getRandDBNameTestTree(0)
				require.NoError(t, err)
				tree.set([]byte("test01"), []byte("test01valueold"))
				tree.set([]byte("test02"), []byte("test02valueold"))
				tree.SaveVersion(false)
				tree.SaveVersion(false)
				return tree.ndb
			},
		},
	}
	for _, c := range cases {
		ndb := c.init()
		iter := newUnsavedFastIterator(start, end, true, ndb, c.unsavedfastNodeAdditions, c.unsavedfastNodeRemovals)
		for i := 0; iter.Valid(); iter.Next() {
			expectedKey := c.expected[i].key
			expectedValue := c.expected[i].fastNode.(*FastNode).value
			require.EqualValues(t, iter.Key(), expectedKey)
			require.EqualValues(t, iter.Value(), expectedValue)
			i++
		}
	}
}
