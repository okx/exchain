package types

// nolint
type FeeDetail struct {
	Address   string `gorm:"index;type:varchar(80)" json:"address" v2:"address"`
	Receiver  string `gorm:"index;type:varchar(80)" json:"receiver" v2:"receiver"` // added for opendex
	Fee       string `gorm:"type:varchar(40)" json:"fee" v2:"fee"`
	FeeType   string `gorm:"index;type:varchar(20)" json:"fee_type" v2:"fee_type"` // defined in order/types/const.go
	Timestamp int64  `gorm:"type:bigint" json:"timestamp" v2:"timestamp"`
}
