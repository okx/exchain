package delta

type DeltaBroker interface {
	GetLocker() bool
	ReleaseLocker()
	ResetLatestHeightAfterUpload(height int64, upload func() bool) bool
	SetDeltas(height int64, bytes []byte) error
	GetDeltas(height int64) ([]byte, error, int64)
}
