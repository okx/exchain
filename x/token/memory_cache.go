package token

const (
	// InitFeeDetailsCap default fee detail list cap
	InitFeeDetailsCap = 2000
)

// Cache for detail
type Cache struct {
	FeeDetails []*FeeDetail
}

// NewCache new cache for detail
func NewCache() *Cache {
	return &Cache{
		FeeDetails: []*FeeDetail{},
	}
}

// Reset reset cache
func (c *Cache) Reset() {
	feeDetails := make([]*FeeDetail, 0, InitFeeDetailsCap)
	c.FeeDetails = feeDetails
}

// AddFeeDetail add fee to cache
func (c *Cache) AddFeeDetail(feeDetail *FeeDetail) {
	c.FeeDetails = append(c.FeeDetails, feeDetail)
}

// GetFeeDetailList gets fee detail from cache
func (c *Cache) GetFeeDetailList() []*FeeDetail {
	return c.FeeDetails
}
