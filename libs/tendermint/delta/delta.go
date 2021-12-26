package delta

type DeltaBroker interface {
	GetLocker() bool
	ReleaseLocker()
	ResetLatestHeightAfterUpload(height int64, getBytes func() ([]byte, bool)) bool
	SetDeltas(height int64, bytes []byte) error
	GetDeltas(height int64) ([]byte, error)
}
