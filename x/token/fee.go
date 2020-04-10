package token

// FeeDetail fee detail
type FeeDetail struct {
	Address   string `gorm:"index;type:varchar(80)" json:"address" v2:"address"`
	Fee       string `gorm:"type:varchar(40)" json:"fee" v2:"fee"`
	FeeType   string `gorm:"index;type:varchar(20)" json:"fee_type" v2:"fee_type"` // transfer, deal, etc. see common/const.go
	Timestamp int64  `gorm:"index;bigint" json:"timestamp" v2:"timestamp"`
}
