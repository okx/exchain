package delta

type DeltaBroker interface {
	SetLatestHeight(height int64) bool
	SetBlock(height int64, bytes []byte) error
	SetDeltas(height int64, bytes []byte) error
	GetBlock(height int64) ([]byte, error)
	GetDeltas(height int64) ([]byte, error)
}
