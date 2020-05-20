package token

import "encoding/json"

const (
	// InitFeeDetailsCap default fee detail list cap
	InitFeeDetailsCap = 2000
)

// Cache for detail
type Cache struct {
	FeeDetails []*FeeDetail `json:"fee_details"`
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

// nolint
func (c *Cache) Clone() *Cache {
	cache := &Cache{}
	bytes, _ := json.Marshal(c)
	err := json.Unmarshal(bytes, cache)
	if err != nil {
		return c.DepthCopy()
	}

	return cache
}

// nolint
func (c *Cache) DepthCopy() *Cache {
	cpFeeDetails := make([]*FeeDetail, 0, len(c.FeeDetails))
	for _, v := range c.FeeDetails {
		cpFeeDetails = append(cpFeeDetails, &FeeDetail{
			Address: v.Address,
			Fee: v.Fee,
			FeeType: v.FeeType,
			Timestamp: v.Timestamp,
		})
	}

	return &Cache{FeeDetails: cpFeeDetails}
}