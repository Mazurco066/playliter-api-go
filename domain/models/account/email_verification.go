package account

import (
	"gorm.io/gorm"
)

type EmailVerification struct {
	gorm.Model
	Code      string  `json:"code"`
	AccountID uint    `json:"account_id"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account"`
}
