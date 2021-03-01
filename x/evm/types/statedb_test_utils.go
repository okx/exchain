package types

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func CompareCommitStateDB(t *testing.T, db1, db2 *CommitStateDB) {

	require.True(t, reflect.DeepEqual(db1.storeKey, db2.storeKey))
	require.True(t, reflect.DeepEqual(db1.paramSpace, db2.paramSpace))
	require.True(t, reflect.DeepEqual(db1.ctx, db2.ctx))

	//require.True(t, reflect.DeepEqual(db1.stateObjects, db2.stateObjects))
	require.True(t, reflect.DeepEqual(db1.stateObjectsDirty, db2.stateObjectsDirty))
	require.True(t, reflect.DeepEqual(db1.addressToObjectIndex, db2.addressToObjectIndex))

	require.True(t, reflect.DeepEqual(db1.logSize, db2.logSize))
	require.True(t, reflect.DeepEqual(db1.refund, db2.refund))

	require.True(t, reflect.DeepEqual(db1.preimages, db2.preimages))
	require.True(t, reflect.DeepEqual(db1.hashToPreimageIndex, db2.hashToPreimageIndex))

	require.True(t, reflect.DeepEqual(db1.journal, db2.journal))
	require.True(t, reflect.DeepEqual(db1.thash, db2.thash))
	require.True(t, reflect.DeepEqual(db1.bhash, db2.bhash))
	require.True(t, reflect.DeepEqual(db1.txIndex, db2.txIndex))

	require.True(t, reflect.DeepEqual(db1.validRevisions, db2.validRevisions))

	require.True(t, reflect.DeepEqual(db1.nextRevisionID, db2.nextRevisionID))
	require.True(t, reflect.DeepEqual(db1.accessList.addresses, db2.accessList.addresses))
	require.Equal(t, len(db1.accessList.slots), len(db2.accessList.slots))
	for i, v := range db1.accessList.slots {
		require.True(t, reflect.DeepEqual(v, db2.accessList.slots[i]))

	}

	require.True(t, reflect.DeepEqual(db1.dbErr, db2.dbErr))


	for i, obj := range db1.stateObjects {
		obj2 := db2.stateObjects[i]
		require.Equal(t, obj.address , obj2.address)
		require.Equal(t, obj.stateObject.address, obj2.stateObject.address)
		require.True(t, reflect.DeepEqual(obj.stateObject.deleted, obj2.stateObject.deleted))
		require.True(t, reflect.DeepEqual(obj.stateObject.dbErr, obj2.stateObject.dbErr))
		require.True(t, reflect.DeepEqual(obj.stateObject.code, obj2.stateObject.code))
		require.True(t, reflect.DeepEqual(obj.stateObject.address, obj2.stateObject.address))
		require.True(t, reflect.DeepEqual(obj.stateObject.account, obj2.stateObject.account))
		require.True(t, reflect.DeepEqual(obj.stateObject.dirtyCode, obj2.stateObject.dirtyCode))
		require.True(t, reflect.DeepEqual(obj.stateObject.dirtyStorage, obj2.stateObject.dirtyStorage))
		require.True(t, reflect.DeepEqual(obj.stateObject.keyToDirtyStorageIndex, obj2.stateObject.keyToDirtyStorageIndex))
		require.True(t, reflect.DeepEqual(obj.stateObject.keyToOriginStorageIndex, obj2.stateObject.keyToOriginStorageIndex))
		require.True(t, reflect.DeepEqual(obj.stateObject.originStorage, obj2.stateObject.originStorage))
		require.True(t, reflect.DeepEqual(obj.stateObject.suicided, obj2.stateObject.suicided))
		require.True(t, obj.stateObject.stateDB == db1)
		require.True(t, obj2.stateObject.stateDB == db2)
	}

}