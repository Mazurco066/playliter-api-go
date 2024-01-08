package account

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerification struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Code      string
	AccountID uint
	Account   Account `gorm:"foreignKey:AccountID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
