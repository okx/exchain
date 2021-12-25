package delta

type DeltaBroker interface {
	GetLocker() bool
	ReleaseLocker()
	ResetLatestHeightAfterUpload(height int64, uploadBytes []byte) bool
	SetDeltas(height int64, bytes []byte) error
	GetDeltas(height int64) ([]byte, error)
}
