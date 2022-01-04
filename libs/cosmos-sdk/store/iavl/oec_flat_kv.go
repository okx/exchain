package iavl

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

func (st *Store) commitFlatKV() {
	if st.flatKVStore == nil {
		return
	}
	st.flatKVStore.Commit()
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
	return st.flatKVStore.GetDBReadTime()
}

func (st *Store) GetFlatKVReadCount() int {
	return st.flatKVStore.GetDBReadCount()
}

func (st *Store) GetFlatKVWriteCount() int {
	return st.flatKVStore.GetDBWriteCount()
}
