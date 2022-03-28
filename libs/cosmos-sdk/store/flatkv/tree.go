package flatkv

type Tree interface {
	Get(key []byte) (index int64, value []byte)
	ShouldPersist(version int64) bool
}
