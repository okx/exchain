package token

const (
	// InitFeeDetailsCap default fee detail list cap
	InitFeeDetailsCap = 2000
)

// Cache for detail
type Cache struct {
	FeeDetails []*FeeDetail
}

// NewCache news cache for detail
func NewCache() *Cache {
	return &Cache{
		FeeDetails: []*FeeDetail{},
	}
}

func (c *Cache) reset() {
	feeDetails := make([]*FeeDetail, 0, InitFeeDetailsCap)
	c.FeeDetails = feeDetails
}

func (c *Cache) addFeeDetail(feeDetail *FeeDetail) {
	c.FeeDetails = append(c.FeeDetails, feeDetail)
}

func (c *Cache) getFeeDetailList() []*FeeDetail {
	return c.FeeDetails
}
