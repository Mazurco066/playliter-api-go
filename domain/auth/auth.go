package auth

import (
	"time"

	"github.com/mazurco066/playliter-api-go/domain/account"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	ID                 uint `gorm:"primaryKey"`
	Token              string
	ResetPasswordToken *string
	AccountID          uint
	Account            account.Account `gorm:"foreignKey:AccountID"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
