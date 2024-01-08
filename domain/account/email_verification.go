package account

import (
	"gorm.io/gorm"
)

type EmailVerification struct {
	gorm.Model
	Code      string
	AccountID uint
	Account   Account `gorm:"foreignKey:AccountID"`
}
