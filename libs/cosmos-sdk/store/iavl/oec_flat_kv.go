package iavl

import "fmt"

func (st *Store) getFlatKV(key []byte) []byte {
	if st.flatKVStore == nil {
		return nil
	}
	return st.flatKVStore.Get(key)
}

func (st *Store) setFlatKV(key, value []byte) {
	if st.flatKVStore == nil {
		return
	}
	st.flatKVStore.Set(key, value)
}

func (st *Store) commitFlatKV(version int64) {
	if st.flatKVStore == nil {
		return
	}
	st.flatKVStore.Commit(version)
}

func (st *Store) hasFlatKV(key []byte) bool {
	if st.flatKVStore == nil {
		return false
	}
	return st.flatKVStore.Has(key)
}

func (st *Store) deleteFlatKV(key []byte) {
	if st.flatKVStore == nil {
		return
	}
	st.flatKVStore.Delete(key)
}

func (st *Store) resetFlatKVCount() {
	if st.flatKVStore == nil {
		return
	}
	st.flatKVStore.ResetCount()
}

func (st *Store) GetFlatKVReadTime() int {
	if st.flatKVStore == nil {
		return 0
	}
	return st.flatKVStore.GetDBReadTime()
}

func (st *Store) GetFlatKVWriteTime() int {
	if st.flatKVStore == nil {
		return 0
	}
	return st.flatKVStore.GetDBWriteTime()
}

func (st *Store) GetFlatKVReadCount() int {
	if st.flatKVStore == nil {
		return 0
	}
	return st.flatKVStore.GetDBReadCount()
}

func (st *Store) GetFlatKVWriteCount() int {
	if st.flatKVStore == nil {
		return 0
	}
	return st.flatKVStore.GetDBWriteCount()
}

func (st *Store) ValidateFlatVersion() error {
	if !st.flatKVStore.Enable() {
		return nil
	}

	treeVersion := st.tree.Version()
	flatVersion := st.flatKVStore.GetLatestVersion()
	if flatVersion != 0 && flatVersion != treeVersion {
		return fmt.Errorf("the version of flat db(%d) does not match the version of iavl tree(%d), you can delete flat.db and restart node",
			flatVersion, treeVersion)
	}
	return nil
}
